package main

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	. "github.com/dave/jennifer/jen"
)

// GenConfig holds the configuration for code generation
type GenConfig struct {
	PackageName      string   // Target package name
	TypeName         string   // The name of the struct type to generate
	ConstantIdent    string   // Prefix for constants (e.g., "Post" for "PostMyPostID")
	VarPrefix        string   // Prefix for variables (e.g., "Post" for "PostMyPost")
	OutputFile       string   // Output file name
	IdentifierFields []string // Fields to try using for naming, in priority order (optional)
	// Custom function to generate variable names (optional)
	// If provided, this takes precedence over IdentifierFields
	CustomVarNameFn func(structValue reflect.Value) string
}

// Generator is responsible for generating code for static struct arrays
type Generator struct {
	Config GenConfig
	Data   any // The array of structs to generate code for
	File   *File
}

// NewGenerator creates a new generator instance
func NewGenerator(config GenConfig, data any) *Generator {
	// Set default identifier fields if none provided
	if config.IdentifierFields == nil {
		config.IdentifierFields = []string{"ID", "Name", "Title", "Slug", "Key", "Code"}
	}

	return &Generator{
		Config: config,
		Data:   data,
		File:   NewFile(config.PackageName),
	}
}

// slugToIdentifier converts a string to a valid Go identifier
func slugToIdentifier(s string) string {
	// Replace non-alphanumeric characters with spaces
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	processed := reg.ReplaceAllString(s, " ")

	// Title case each word and remove spaces
	words := strings.Fields(processed)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[0:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, "")
}

// Generate performs the code generation
func (g *Generator) Generate() error {
	// Add package comment
	g.File.PackageComment(fmt.Sprintf("// Code generated within %s for %s. DO NOT EDIT.",
		g.Config.PackageName, g.Config.TypeName))

	// Validate that we have an array or slice
	dataValue := reflect.ValueOf(g.Data)
	if dataValue.Kind() != reflect.Slice && dataValue.Kind() != reflect.Array {
		return fmt.Errorf("data must be a slice or array, got %s", dataValue.Kind())
	}

	// Make sure we have at least one element to analyze the type
	if dataValue.Len() == 0 {
		return fmt.Errorf("data must contain at least one element")
	}

	// Get the type of the first element
	firstElem := dataValue.Index(0)
	if firstElem.Kind() != reflect.Struct {
		return fmt.Errorf("data elements must be structs, got %s", firstElem.Kind())
	}

	// Generate constants for IDs if there's an ID field
	g.generateConstants(dataValue)

	// Generate variables for each struct
	g.generateVariables(dataValue)

	// Generate a slice with all structs
	g.generateSlice(dataValue)

	// Save the generated code to file
	return g.File.Save(g.Config.OutputFile)
}

// generateStructType creates the struct type definition
func (g *Generator) generateStructType(structType reflect.Type) {
	// Create struct with fields
	g.File.Type().Id(g.Config.TypeName).StructFunc(func(group *Group) {
		for i := range structType.NumField() {
			field := structType.Field(i)

			// Skip unexported fields
			if !field.IsExported() {
				continue
			}

			// Add the field to the struct
			fieldType := g.getTypeStatement(field.Type)

			// Add json tags if they exist
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" {
				group.Id(field.Name).Add(fieldType).Tag(map[string]string{"json": jsonTag})
			} else {
				group.Id(field.Name).Add(fieldType)
			}
		}
	})
}

// getTypeStatement converts a reflect.Type to a jen.Statement
func (g *Generator) getTypeStatement(t reflect.Type) *Statement {
	switch t.Kind() {
	case reflect.Bool:
		return Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Id(t.String())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return Id(t.String())
	case reflect.Float32, reflect.Float64:
		return Id(t.String())
	case reflect.Complex64, reflect.Complex128:
		return Id(t.String())
	case reflect.Array, reflect.Slice:
		return Index().Add(g.getTypeStatement(t.Elem()))
	case reflect.Map:
		return Map(g.getTypeStatement(t.Key())).Add(g.getTypeStatement(t.Elem()))
	case reflect.String:
		return String()
	case reflect.Struct:
		// Handle special types like time.Time
		if t.String() == "time.Time" {
			return Qual("time", "Time")
		}
		return Id(t.Name())
	case reflect.Pointer:
		return Op("*").Add(g.getTypeStatement(t.Elem()))
	case reflect.Interface:
		if t.NumMethod() == 0 {
			return Interface() // empty interface
		}
		// Complex interfaces would need more handling
		return Interface()
	default:
		return Id(t.String())
	}
}

// generateConstants creates ID constants for each struct if an ID field exists
func (g *Generator) generateConstants(dataValue reflect.Value) {
	// Check if the struct has an ID field
	firstElem := dataValue.Index(0)
	hasIDField := false
	idFieldName := ""

	// Look for an "ID" field (case insensitive)
	for i := 0; i < firstElem.NumField(); i++ {
		fieldName := firstElem.Type().Field(i).Name
		if strings.ToLower(fieldName) == "id" {
			hasIDField = true
			idFieldName = fieldName
			break
		}
	}

	if !hasIDField {
		return // No ID field found
	}

	// Create constants for each ID
	g.File.Const().DefsFunc(func(group *Group) {
		for i := range dataValue.Len() {
			elem := dataValue.Index(i)
			idField := elem.FieldByName(idFieldName)

			// If there's an ID field that's a string, create a constant
			if idField.IsValid() && idField.Kind() == reflect.String {
				idValue := idField.String()
				// If ID is empty, generate one
				if idValue == "" {
					idValue = fmt.Sprintf("%s-%d", strings.ToLower(g.Config.TypeName), i+1)
				}

				// Get a name for the constant based on the struct
				identValue := g.getStructIdentifier(elem)

				constName := g.Config.ConstantIdent + slugToIdentifier(identValue) + "ID"
				group.Id(constName).Op("=").Lit(idValue)
			}
		}
	})
}

// getStructIdentifier returns a string to identify this struct instance
func (g *Generator) getStructIdentifier(structValue reflect.Value) string {
	// If a custom name function is provided, use it
	if g.Config.CustomVarNameFn != nil {
		return g.Config.CustomVarNameFn(structValue)
	}

	// Try all configured identifier fields
	for _, fieldName := range g.Config.IdentifierFields {
		field := structValue.FieldByName(fieldName)
		if field.IsValid() && field.Kind() == reflect.String && field.String() != "" {
			return field.String()
		}
	}

	// Fallback 1: Look for any string field
	for i := range structValue.NumField() {
		field := structValue.Field(i)
		if field.Kind() == reflect.String && field.String() != "" {
			return field.String()
		}
	}

	// Fallback 2: Generate a name based on the type
	return fmt.Sprintf("%s-%d", g.Config.TypeName, time.Now().UnixNano())
}

// generateVariables creates variables for each struct
func (g *Generator) generateVariables(dataValue reflect.Value) {
	// Generate a variable for each struct
	for i := range dataValue.Len() {
		elem := dataValue.Index(i)

		// Determine the variable name using the identifier function
		identValue := g.getStructIdentifier(elem)
		varName := g.Config.VarPrefix + slugToIdentifier(identValue)

		// Create the variable with its value
		g.File.Var().Id(varName).Op("=").Id(g.Config.TypeName).ValuesFunc(func(group *Group) {
			g.generateStructValues(group, elem)
		})
	}
}

// generateStructValues adds values for a struct to a Dict
func (g *Generator) generateStructValues(group *Group, structValue reflect.Value) {
	structType := structValue.Type()

	// Create a Dict for each field in the struct
	dict := Dict{}

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		// Add the field to the dict
		dict[Id(fieldType.Name)] = g.getValueStatement(field)
	}

	// Add all fields to the group
	group.Add(dict)
}

// getValueStatement generates code for a value based on its type
func (g *Generator) getValueStatement(value reflect.Value) *Statement {
	switch value.Kind() {
	case reflect.Bool:
		return Lit(value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Lit(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Lit(value.Uint())
	case reflect.Float32, reflect.Float64:
		return Lit(value.Float())
	case reflect.Complex64, reflect.Complex128:
		return Lit(value.Complex())
	case reflect.Array, reflect.Slice:
		return g.getSliceStatement(value)
	case reflect.Map:
		return g.getMapStatement(value)
	case reflect.String:
		return Lit(value.String())
	case reflect.Struct:
		// Special case for time.Time
		if value.Type().String() == "time.Time" {
			t := value.Interface().(time.Time)
			return Qual("time", "Date").Call(
				Lit(t.Year()),
				Qual("time", t.Month().String()),
				Lit(t.Day()),
				Lit(t.Hour()),
				Lit(t.Minute()),
				Lit(t.Second()),
				Lit(t.Nanosecond()),
				Qual("time", "UTC"),
			)
		}
		// For other structs, create a new values block with the struct fields
		return Id(value.Type().Name()).ValuesFunc(func(group *Group) {
			g.generateStructValues(group, value)
		})
	case reflect.Pointer:
		if value.IsNil() {
			return Nil()
		}
		return Op("&").Add(g.getValueStatement(value.Elem()))
	case reflect.Interface:
		if value.IsNil() {
			return Nil()
		}
		return g.getValueStatement(value.Elem())
	default:
		// For complex cases, fallback to string representation
		return Lit(fmt.Sprintf("%v", value.Interface()))
	}
}

// getSliceStatement generates code for a slice
func (g *Generator) getSliceStatement(sliceValue reflect.Value) *Statement {
	// Create a new statement for the slice values
	stmt := &Statement{}

	// Add opening and closing braces
	stmt.Add(ListFunc(func(group *Group) {
		for i := range sliceValue.Len() {
			group.Add(g.getValueStatement(sliceValue.Index(i)))
		}
	}))

	return stmt
}

// getMapStatement generates code for a map
func (g *Generator) getMapStatement(mapValue reflect.Value) *Statement {
	// Create a Dict for the map entries
	mapDict := Dict{}

	// Add all key-value pairs to the Dict
	for _, key := range mapValue.MapKeys() {
		mapDict[g.getValueStatement(key)] = g.getValueStatement(mapValue.MapIndex(key))
	}

	// Return the Dict inside Values
	return Values(mapDict)
}

// generateSlice creates a slice containing all struct instances
func (g *Generator) generateSlice(dataValue reflect.Value) {
	// Determine the slice name - handle both regular and irregular plurals
	var sliceName string
	if g.Config.TypeName[len(g.Config.TypeName)-1] == 's' ||
		g.Config.TypeName[len(g.Config.TypeName)-1] == 'x' ||
		g.Config.TypeName[len(g.Config.TypeName)-1] == 'z' ||
		strings.HasSuffix(g.Config.TypeName, "sh") ||
		strings.HasSuffix(g.Config.TypeName, "ch") {
		sliceName = fmt.Sprintf("All%ses", g.Config.TypeName)
	} else if g.Config.TypeName[len(g.Config.TypeName)-1] == 'y' {
		sliceName = fmt.Sprintf("All%sies", g.Config.TypeName[:len(g.Config.TypeName)-1])
	} else {
		sliceName = fmt.Sprintf("All%ss", g.Config.TypeName)
	}

	g.File.Var().Id(sliceName).Op("=").Index().Id(g.Config.TypeName).ValuesFunc(func(group *Group) {
		for i := range dataValue.Len() {
			elem := dataValue.Index(i)

			// Get the variable name using the same method as in generateVariables
			identValue := g.getStructIdentifier(elem)
			varName := g.Config.VarPrefix + slugToIdentifier(identValue)

			group.Id(varName)
		}
	})
}

// Person is an example struct
type Person struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birth_date"`
	Address   Address   `json:"address"`
}

// Address is an example struct
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

func main() {

	// Example for Persons
	people := []Person{
		{
			ID:        "person-1",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			BirthDate: time.Date(1985, time.June, 15, 0, 0, 0, 0, time.UTC),
			Address: Address{
				Street:  "123 Main St",
				City:    "Boston",
				State:   "MA",
				ZipCode: "02108",
				Country: "USA",
			},
		},
		{
			ID:        "person-2",
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			BirthDate: time.Date(1990, time.August, 22, 0, 0, 0, 0, time.UTC),
			Address: Address{
				Street:  "456 Oak Ave",
				City:    "San Francisco",
				State:   "CA",
				ZipCode: "94107",
				Country: "USA",
			},
		},
	}

	// Create a generator for people
	personGen := NewGenerator(GenConfig{
		PackageName:   "main",
		TypeName:      "Person",
		ConstantIdent: "Person",
		VarPrefix:     "Person",
		OutputFile:    "people.go",
	}, people)

	// Generate code for people
	err := personGen.Generate()
	if err != nil {
		fmt.Println("Error generating person code:", err)
		os.Exit(1)
	}

	fmt.Println("Code generation completed successfully!")
}
