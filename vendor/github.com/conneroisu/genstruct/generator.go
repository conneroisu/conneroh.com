package genstruct

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dave/jennifer/jen"
)

// Generator is responsible for generating code for static struct arrays
type Generator struct {
	// Primary configuration options
	PackageName      string
	TypeName         string
	ConstantIdent    string
	VarPrefix        string
	OutputFile       string
	IdentifierFields []string
	CustomVarNameFn  func(structValue reflect.Value) string
	Logger           *slog.Logger

	// Internal state
	Data   any            // The primary array of structs to generate code for
	Refs   map[string]any // Additional arrays that can be referenced
	File   *jen.File
}

// Option is a functional option for customizing the generator.
type Option func(g *Generator)

// WithPackageName sets the package name for the generated code.
// If not specified, the package name is inferred from the output file directory.
func WithPackageName(name string) Option {
	return func(g *Generator) { g.PackageName = name }
}

// WithTypeName sets the type name for the generated code.
// If not specified, the type name is inferred from the data struct type.
func WithTypeName(name string) Option {
	return func(g *Generator) { g.TypeName = name }
}

// WithConstantIdent sets the prefix for generated constants.
// For example, with prefix "Animal", constants will be named "AnimalLionID", etc.
// If not specified, defaults to the TypeName.
func WithConstantIdent(name string) Option {
	return func(g *Generator) { g.ConstantIdent = name }
}

// WithVarPrefix sets the prefix for generated variables.
// For example, with prefix "Animal", variables will be named "AnimalLion", etc.
// If not specified, defaults to the TypeName.
func WithVarPrefix(name string) Option {
	return func(g *Generator) { g.VarPrefix = name }
}

// WithOutputFile sets the output file path for the generated code.
// The path can include directories.
// Export mode is automatically determined based on this path - if it contains
// directory separators, qualified imports will be used for external types.
// If not specified, defaults to lowercase(typename_generated.go).
func WithOutputFile(path string) Option {
	return func(g *Generator) { g.OutputFile = path }
}

// WithIdentifierFields sets the fields to use for variable naming.
// These fields are checked in order until a non-empty string field is found.
// If not specified, defaults to ["ID", "Name", "Slug", "Title", "Key", "Code"].
func WithIdentifierFields(fields []string) Option {
	return func(g *Generator) { g.IdentifierFields = fields }
}

// WithCustomVarNameFn sets a custom function to control variable naming.
// This takes precedence over IdentifierFields if provided.
// The function receives a reflect.Value of the struct and should return a string
// to be used as the base name for the variable.
func WithCustomVarNameFn(fn func(structValue reflect.Value) string) Option {
	return func(g *Generator) { g.CustomVarNameFn = fn }
}

// WithLogger sets a custom slog.Logger instance for logging during generation.
// If not specified, the default logger is used.
func WithLogger(logger *slog.Logger) Option {
	return func(g *Generator) { g.Logger = logger }
}

//

// NewGenerator creates a new generator instance with the specified options.
//
// Example usage:
//
//	generator := genstruct.NewGenerator(
//	    genstruct.WithPackageName("mypackage"),
//	    genstruct.WithOutputFile("output_file.go"),
//	)
//	
//	// Generate code for posts with tags references
//	err := generator.Generate(posts, tags)
//
// Configuration is handled through functional options, and many values are automatically
// inferred if not specified:
//   - TypeName: Inferred from the struct type in the data provided to Generate()
//   - ConstantIdent: Defaults to TypeName if not specified
//   - VarPrefix: Defaults to TypeName if not specified
//   - OutputFile: Defaults to lowercase(typename_generated.go) if not specified
//   - IdentifierFields: Uses default fields if not specified
//   - Logger: Uses the default logger if not specified
//
// Export mode (referencing types from other packages) is automatically determined
// based on the output file path. If the path contains directory separators,
// it will use qualified imports when referencing types from other packages.
func NewGenerator(opts ...Option) *Generator {
	// Create a new generator with default values
	g := &Generator{
		Refs: make(map[string]any),
		IdentifierFields: []string{
			"ID",
			"Name",
			"Slug",
			"Title",
			"Key",
			"Code",
		},
		Logger: GetLogger(),
	}

	// Apply options
	for _, opt := range opts {
		opt(g)
	}

	return g
}

// inferConfig analyzes the data and fills in any missing configuration values.
// This automatically determines struct type names and other configurable values
// that haven't been explicitly set through functional options.
func (g *Generator) inferConfig(data any) error {
	// Get the element type from the data
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice && dataValue.Kind() != reflect.Array {
		// Can't determine type from non-slice/array
		return InvalidTypeError{Kind: dataValue.Kind()}
	}

	// Make sure we have at least one element to analyze
	if dataValue.Len() == 0 {
		// Can't determine type from empty slice
		return EmptyError{}
	}

	firstElem := dataValue.Index(0)
	var structType reflect.Type

	// Support both direct struct slices and pointer slices
	if firstElem.Kind() == reflect.Struct {
		structType = firstElem.Type()
	} else if firstElem.Kind() == reflect.Pointer && firstElem.Elem().Kind() == reflect.Struct {
		structType = firstElem.Elem().Type()
	} else {
		// Only struct or struct pointer slices are supported
		return InvalidTypeError{Kind: firstElem.Kind()}
	}

	typeName := structType.Name()

	// Infer TypeName if not specified
	if g.TypeName == "" {
		g.TypeName = typeName
	}

	// Infer ConstantIdent if not specified
	if g.ConstantIdent == "" {
		g.ConstantIdent = g.TypeName
	}

	// Infer VarPrefix if not specified
	if g.VarPrefix == "" {
		g.VarPrefix = g.TypeName
	}

	// Infer OutputFile if not specified
	if g.OutputFile == "" {
		g.OutputFile = strings.ToLower(g.TypeName) + "_generated.go"
	}

	// If PackageName is not specified, use the directory name from the output file
	if g.PackageName == "" {
		g.PackageName = GetPackageNameFromPath(g.OutputFile)
	}

	// Log the configuration
	g.Logger.Debug(
		"Configuration inferred",
		slog.String("typeName", g.TypeName),
		slog.String("packageName", g.PackageName),
		slog.String("outputFile", g.OutputFile),
	)

	return nil
}

// GetPackageNameFromPath extracts the containing folder name from a file path
// This can be used to determine the package name for a given Go file
// Example: "./out/penguin/gen.go" would return "penguin"
func GetPackageNameFromPath(filePath string) string {
	// Clean the path to handle any OS-specific separators and normalize it
	cleanPath := filepath.Clean(filePath)

	// Get the directory containing the file
	dir := filepath.Dir(cleanPath)

	// Split the directory path into components
	components := strings.Split(dir, string(filepath.Separator))

	// The package name is the last component of the directory path
	// If the path ends with a separator, the last component will be empty
	if len(components) > 0 {
		lastComponent := components[len(components)-1]
		if lastComponent != "" {
			return lastComponent
		}

		// If the last component is empty, try the second-to-last one
		if len(components) > 1 {
			return components[len(components)-2]
		}
	}

	// Default to "main" if we couldn't extract a package name
	return "main"
}

// Generate performs the code generation for both primary data and reference data.
//
// Parameters:
//   - data: The primary array of structs to generate code for
//   - refs: Optional additional arrays that can be referenced by the primary data
//
// The refs parameters enable struct references via the `structgen` tag. For example,
// a Post struct with a TagSlugs field can reference Tag structs:
//
//     type Post struct {
//         ID       string
//         TagSlugs []string  // Contains identifiers
//         Tags     []*Tag    `structgen:"TagSlugs"` // Will be populated automatically
//     }
//
// Reference fields can be either direct structs or pointers to structs:
//   - []Tag `structgen:"TagSlugs"` - Direct struct references
//   - []*Tag `structgen:"TagSlugs"` - Pointer-based struct references (recommended)
//
// This method generates:
// 1. Constants for the primary data's IDs
// 2. Variables for each item in the primary data
// 3. A slice containing all primary data items
// 4. Constants for each reference data set's IDs
// 5. Variables for each item in each reference data set
// 6. A slice for each reference data set
// 7. Creates references between primary data and reference data as specified by structgen tags
//
// All generated code is written to a single output file specified in the OutputFile field.
//
// Returns an error if:
//   - The data is not a slice or array
//   - The data is empty (no elements to analyze)
//   - The data elements are not structs
//   - Required fields couldn't be inferred
func (g *Generator) Generate(data any, refs ...any) error {
	// Store the data for processing
	g.Data = data
	
	// Create a map of reference datasets
	g.Refs = make(map[string]any)
	for i, ref := range refs {
		// Get type name for this dataset
		refType := reflect.TypeOf(ref)
		if refType.Kind() == reflect.Slice || refType.Kind() == reflect.Array {
			elemType := refType.Elem()
			if elemType.Kind() == reflect.Struct {
				g.Refs[elemType.Name()] = ref
			} else if elemType.Kind() == reflect.Pointer &&
				elemType.Elem().Kind() == reflect.Struct {
				// Handle pointer slice ([]*Type)
				g.Refs[elemType.Elem().Name()] = ref
			} else {
				g.Refs[fmt.Sprintf("Ref%d", i)] = ref
			}
		}
	}
	
	// Infer config options based on the data
	if err := g.inferConfig(data); err != nil {
		return err
	}
	
	// Initialize the file with the package name
	g.File = jen.NewFile(g.PackageName)
	
	g.Logger.Info(
		"Starting code generation",
		slog.String("package", g.PackageName),
		slog.String("type", g.TypeName),
		slog.String("output", g.OutputFile),
	)
	
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("failed to read build info for version number")
	}
	
	// Find github.com/conneroisu/genstruct dep
	var dep *debug.Module
	for _, d := range bi.Deps {
		if d.Path == "github.com/conneroisu/genstruct" {
			dep = d
			break
		}
	}
	if dep == nil {
		dep = &debug.Module{
			Path:    "github.com/conneroisu/genstruct",
			Version: "Unknown",
		}
	}
	
	g.File.PackageComment(fmt.Sprintf(
		"// Code generated by genstruct. DO NOT EDIT.\n// Package %s contains auto-generated %s data\n//\n// genstruct Version: %s\n//",
		g.PackageName,
		g.TypeName,
		dep.Version,
	))

	// Validate that we have an array or slice
	dataValue := reflect.ValueOf(g.Data)
	if dataValue.Kind() != reflect.Slice &&
		dataValue.Kind() != reflect.Array {
		g.Logger.Error(
			"Invalid data type",
			"expected",
			"slice or array",
			"got",
			dataValue.Kind().String(),
		)
		return NonSliceOrArrayError{dataValue.Kind()}
	}

	// Make sure we have at least one element to analyze the type
	if dataValue.Len() == 0 {
		g.Logger.Error("Empty data slice", "type", g.TypeName)
		return EmptyError{}
	}

	// Get the type of the first element
	firstElem := dataValue.Index(0)
	// Support both direct struct slices and pointer slices
	if firstElem.Kind() != reflect.Struct &&
		(firstElem.Kind() != reflect.Pointer ||
			firstElem.Elem().Kind() != reflect.Struct) {
		g.Logger.Error(
			"Invalid element type",
			slog.String("expected", "struct or pointer to struct"),
			slog.String("got", firstElem.Kind().String()),
		)
		return InvalidTypeError{firstElem.Kind()}
	}

	// Generate constants for IDs if there's an ID field
	g.Logger.Debug(
		"Generating constants",
		"type",
		g.TypeName,
	)
	g.generateConstants(dataValue)

	// Generate variables for each struct
	g.Logger.Debug(
		"Generating variables",
		"type",
		g.TypeName,
		"count",
		dataValue.Len(),
	)
	g.generateVariables(dataValue)

	// Generate a slice with all structs
	g.Logger.Debug(
		"Generating slice",
		"type",
		g.TypeName,
	)
	g.generateSlice(dataValue)

	// Process reference datasets to generate their constants and variables
	// This ensures that all referenced types (like Tag in Post.Tags) are properly defined
	// in the generated code, making the references fully usable.
	g.Logger.Debug(
		"Processing reference datasets",
		slog.Int("count", len(g.Refs)),
	)
	for typeName, refDataObj := range g.Refs {
		g.Logger.Debug(
			"Processing reference dataset",
			slog.String("type", typeName),
		)
		refDataValue := reflect.ValueOf(refDataObj)
		if refDataValue.Kind() == reflect.Slice ||
			refDataValue.Kind() == reflect.Array {
			if refDataValue.Len() > 0 {
				refElem := refDataValue.Index(0)
				// Support both direct structs and pointer-to-structs
				if refElem.Kind() == reflect.Struct ||
					(refElem.Kind() == reflect.Pointer &&
						refElem.Elem().Kind() == reflect.Struct) {
					// Store original config values so we can restore them after
					// processing this reference type
					originalTypeName := g.TypeName
					originalVarPrefix := g.VarPrefix
					originalConstantIdent := g.ConstantIdent

					// Temporarily set config values for the reference type
					// This ensures that constants and variables are named correctly
					// (e.g., TagGoProgramming instead of PostGoProgramming)
					g.TypeName = typeName
					g.VarPrefix = typeName
					g.ConstantIdent = typeName

					// Generate constants, variables, and slice for this reference dataset
					// using the same generation methods as for the primary dataset
					g.generateConstants(refDataValue)
					g.generateVariables(refDataValue)
					g.generateSlice(refDataValue)

					// Restore original config values for processing the next reference dataset
					g.TypeName = originalTypeName
					g.VarPrefix = originalVarPrefix
					g.ConstantIdent = originalConstantIdent
				}
			}
		}
	}

	// We no longer need to export embedded types, as they should always use the
	// pkg-qualified name (e.g., pkg.Embedded) when in export mode
	// This block is kept empty as a placeholder comment to explain the change

	// Generate the code as a string
	g.Logger.Debug("Rendering generated code")
	buf := &bytes.Buffer{}
	if err := g.File.Render(buf); err != nil {
		g.Logger.Error("Failed to render code", "error", err)
		return err
	}

	// Save the formatted code to file
	g.Logger.Debug(
		"Writing generated code to file",
		slog.String("file", g.OutputFile),
	)
	return os.WriteFile(g.OutputFile, buf.Bytes(), 0644)
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

// getStructIdentifier returns a string to identify this struct instance
func (g *Generator) getStructIdentifier(structValue reflect.Value) string {
	// Handle pointer to struct case
	if structValue.Kind() == reflect.Pointer {
		structValue = structValue.Elem()
	}

	// If a custom name function is provided, use it
	if g.CustomVarNameFn != nil {
		return g.CustomVarNameFn(structValue)
	}

	// Try all configured identifier fields
	for _, fieldName := range g.IdentifierFields {
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
	return fmt.Sprintf("%s-%d", g.TypeName, time.Now().UnixNano())
}
