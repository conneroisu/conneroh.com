package genstruct

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dave/jennifer/jen"
)

// getValueStatement generates code for a value based on its type
func (g *Generator) getValueStatement(value reflect.Value) *jen.Statement {
	switch value.Kind() {
	case reflect.Bool:
		return jen.Lit(value.Bool())
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return jen.Lit(value.Int())
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return jen.Lit(value.Uint())
	case reflect.Float32, reflect.Float64:
		return jen.Lit(value.Float())
	case reflect.Complex64, reflect.Complex128:
		return jen.Lit(value.Complex())
	case reflect.Array:
		// Handle arrays properly with their type and dimensions
		elemType := g.getTypeStatement(value.Type().Elem())
		dimensions := value.Len()

		// Create array type with dimensions
		arrayType := jen.Index(jen.Lit(dimensions)).Add(elemType)

		// Create values inside the array
		return arrayType.ValuesFunc(func(group *jen.Group) {
			for i := range value.Len() {
				group.Add(g.getValueStatement(value.Index(i)))
			}
		})
	case reflect.Slice:
		// Create a slice with proper syntax
		return jen.Index().Add(
			g.getTypeStatement(value.Type().Elem()),
		).ValuesFunc(func(group *jen.Group) {
			for i := range value.Len() {
				group.Add(g.getValueStatement(value.Index(i)))
			}
		})
	case reflect.Map:
		return g.getMapStatement(value)
	case reflect.String:
		return jen.Lit(value.String())
	case reflect.Struct:
		// Special case for time.Time
		if value.Type().String() == "time.Time" {
			t := value.Interface().(time.Time)
			return jen.Qual("time", "Date").Call(
				jen.Lit(t.Year()),
				jen.Qual("time", t.Month().String()),
				jen.Lit(t.Day()),
				jen.Lit(t.Hour()),
				jen.Lit(t.Minute()),
				jen.Lit(t.Second()),
				jen.Lit(t.Nanosecond()),
				jen.Qual("time", "UTC"),
			)
		}

		// Check if this struct is from another package in export mode
		isExportMode := strings.Contains(g.OutputFile, "/")
		pkgPath := value.Type().PkgPath()

		if isExportMode && pkgPath != "" && pkgPath != "main" && pkgPath != g.PackageName {
			// For structs from another package, use fully qualified names
			return jen.Qual(pkgPath, value.Type().Name()).ValuesFunc(func(group *jen.Group) {
				g.generateStructValues(group, value)
			})
		}

		// For other structs, create a new values block with the struct fields
		return jen.Id(
			value.Type().Name(),
		).ValuesFunc(func(group *jen.Group) {
			g.generateStructValues(group, value)
		})
	case reflect.Pointer:
		if value.IsNil() {
			return jen.Nil()
		}
		return jen.Op("&").Add(g.getValueStatement(value.Elem()))
	case reflect.Interface:
		if value.IsNil() {
			return jen.Nil()
		}
		return g.getValueStatement(value.Elem())
	default:
		// For complex cases, fallback to string representation
		return jen.Lit(fmt.Sprintf("%v", value.Interface()))
	}
}

// getMapStatement generates code for a map
func (g *Generator) getMapStatement(mapValue reflect.Value) *jen.Statement {
	// Return empty map if there are no entries
	if mapValue.Len() == 0 {
		return jen.Map(
			g.getTypeStatement(mapValue.Type().Key()),
		).Add(
			g.getTypeStatement(mapValue.Type().Elem()),
		).Values()
	}

	// Use ValuesFunc for populated maps
	return jen.Map(
		g.getTypeStatement(mapValue.Type().Key()),
	).Add(
		g.getTypeStatement(mapValue.Type().Elem()),
	).ValuesFunc(func(group *jen.Group) {
		var (
			dict = jen.Dict{}
			key  reflect.Value
		)

		// Add all key-value pairs to the Dict
		for _, key = range mapValue.MapKeys() {
			var stmt = g.getValueStatement(mapValue.MapIndex(key))
			dict[g.getValueStatement(key)] = stmt
		}

		// Add dict to group
		group.Add(dict)
	})
}

// generateStructValues adds values for a struct to a Dict
func (g *Generator) generateStructValues(group *jen.Group, structValue reflect.Value) {
	if structValue.Kind() == reflect.Pointer {
		structValue = structValue.Elem()
	}

	structType := structValue.Type()

	dict := jen.Dict{}

	// Track fields that need to be processed in a second pass (with structgen tag)
	type deferredField struct {
		fieldIndex int
		fieldType  reflect.StructField
		srcField   string
	}
	var deferredFields []deferredField

	// First pass: process all regular fields
	for i := range structValue.NumField() {
		var (
			field     = structValue.Field(i)
			fieldType = structType.Field(i)
		)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		// Check if this field has a structgen tag
		structgenVal, hasStructgenTag := fieldType.Tag.Lookup("structgen")

		if hasStructgenTag && structgenVal != "" {
			// Add to deferred fields for second pass
			deferredFields = append(deferredFields, deferredField{
				fieldIndex: i,
				fieldType:  fieldType,
				srcField:   structgenVal,
			})
			continue
		}

		// Handle embedded fields specially in export mode
		isExportMode := strings.Contains(g.OutputFile, "/")
		if fieldType.Anonymous && isExportMode {
			// For embedded fields in export mode, check if it comes from another package
			embeddedType := fieldType.Type
			pkgPath := embeddedType.PkgPath()

			if pkgPath != "" && pkgPath != "main" && pkgPath != g.PackageName {
				// Use qualified package reference for embedded fields from other packages
				// but still generate all the fields inside it
				dict[jen.Id(fieldType.Name)] = jen.Qual(pkgPath, embeddedType.Name()).ValuesFunc(func(embGroup *jen.Group) {
					// Generate inner struct values
					innerDict := jen.Dict{}

					for j := range field.NumField() {
						innerField := field.Field(j)
						innerFieldType := field.Type().Field(j)

						// Skip unexported fields
						if !innerFieldType.IsExported() {
							continue
						}

						// Add each field with its value
						innerDict[jen.Id(innerFieldType.Name)] = g.getValueStatement(innerField)
					}

					embGroup.Add(innerDict)
				})
			} else {
				// Use regular reference for embedded fields from same package
				dict[jen.Id(fieldType.Name)] = g.getValueStatement(field)
			}
		} else {
			// Regular field
			dict[jen.Id(fieldType.Name)] = g.getValueStatement(field)
		}
	}

	// Second pass: process fields with structgen tag
	for _, df := range deferredFields {
		value := g.generateStructGenField(structValue, df.srcField, df.fieldType)
		if value != nil {
			dict[jen.Id(df.fieldType.Name)] = value
		}
	}

	// Add all fields to the group
	group.Add(dict)
}

// generateStructGenField generates a value for a field with the structgen tag
//
// The structgen tag enables automatic population of struct fields from reference datasets.
// It takes the source field name as a value, which should contain identifiers (strings or string slices)
// that can be used to look up matching structs in the reference datasets.
//
// Supported reference patterns:
//   - String to Struct: A string field (e.g., "AuthorID") referencing a single struct or struct pointer (*T)
//   - String Slice to Struct Slice: A slice of strings (e.g., "TagSlugs") referencing a slice of structs ([]T) or struct pointers ([]*T)
//
// Parameters:
//   - structValue: The struct instance being processed
//   - srcFieldName: The name of the source field (from the tag value)
//   - targetField: The field to populate with references
func (g *Generator) generateStructGenField(
	structValue reflect.Value,
	srcFieldName string,
	targetField reflect.StructField,
) *jen.Statement {
	structType := structValue.Type()

	// Find the source field
	srcField, found := structType.FieldByName(srcFieldName)
	if !found {
		// Source field not found
		return nil
	}

	// Get the source field's value
	srcValue := structValue.FieldByName(srcFieldName)
	if !srcValue.IsValid() {
		return nil
	}

	// Determine the target type
	targetType := targetField.Type

	// Check for slice of structs or struct pointers referencing a string slice
	if targetType.Kind() == reflect.Slice &&
		((targetType.Elem().Kind() == reflect.Struct) ||
			(targetType.Elem().Kind() == reflect.Pointer && targetType.Elem().Elem().Kind() == reflect.Struct)) &&
		srcField.Type.Kind() == reflect.Slice &&
		srcField.Type.Elem().Kind() == reflect.String {

		// We need to look up structs by ID or another field
		return g.generateReferenceSlice(srcValue, targetType)
	}

	// Check for single struct or struct pointer referencing a string
	if (targetType.Kind() == reflect.Struct ||
		(targetType.Kind() == reflect.Pointer && targetType.Elem().Kind() == reflect.Struct)) &&
		srcField.Type.Kind() == reflect.String {

		// We need to look up one struct by ID or another field
		return g.generateReferenceSingle(srcValue, targetType)
	}

	// Unsupported reference type
	return nil
}

// generateReferenceSlice generates a slice of referenced structs for string slice to struct slice references
//
// This method handles the case where a field contains a slice of strings (e.g., ["tag1", "tag2"])
// and needs to generate a slice of structs (e.g., []Tag or []*Tag) by looking up each string in a reference dataset.
//
// Parameters:
//   - srcValue: The source field value (slice of strings)
//   - targetType: The target field type (slice of structs or struct pointers)
func (g *Generator) generateReferenceSlice(srcValue reflect.Value, targetType reflect.Type) *jen.Statement {
	// Determine if we're dealing with a pointer slice ([]*T) or struct slice ([]T)
	isPointerSlice := targetType.Elem().Kind() == reflect.Pointer

	// Get the target struct type name
	var structTypeName string
	if isPointerSlice {
		structTypeName = targetType.Elem().Elem().Name()
	} else {
		structTypeName = targetType.Elem().Name()
	}

	// Check if we need to use fully qualified type references
	isExportMode := strings.Contains(g.OutputFile, "/")
	refType := targetType.Elem()
	if isPointerSlice {
		refType = refType.Elem()
	}
	pkgPath := refType.PkgPath()
	useQualified := isExportMode && pkgPath != "" && pkgPath != "main" && pkgPath != g.PackageName

	// Check if we have this reference type
	refDataObj, hasRef := g.Refs[structTypeName]
	if !hasRef {
		// We don't have this reference data
		if isPointerSlice {
			if useQualified {
				return jen.Index().Add(jen.Op("*").Qual(pkgPath, structTypeName)).Values()
			}
			return jen.Index().Add(jen.Op("*").Id(structTypeName)).Values()
		}
		if useQualified {
			return jen.Index().Add(jen.Qual(pkgPath, structTypeName)).Values()
		}
		return jen.Index().Add(jen.Id(structTypeName)).Values()
	}

	// Convert to reflect.Value
	refData := reflect.ValueOf(refDataObj)
	if refData.Kind() != reflect.Slice && refData.Kind() != reflect.Array {
		// Reference isn't a slice/array
		if isPointerSlice {
			if useQualified {
				return jen.Index().Add(jen.Op("*").Qual(pkgPath, structTypeName)).Values()
			}
			return jen.Index().Add(jen.Op("*").Id(structTypeName)).Values()
		}
		if useQualified {
			return jen.Index().Add(jen.Qual(pkgPath, structTypeName)).Values()
		}
		return jen.Index().Add(jen.Id(structTypeName)).Values()
	}

	// Create a statement for the appropriate slice type
	var sliceStmt *jen.Statement

	// Use the qualified type if needed
	if useQualified {
		if isPointerSlice {
			// For []*pkg.T
			sliceStmt = jen.Index().Add(jen.Op("*").Qual(pkgPath, structTypeName))
		} else {
			// For []pkg.T
			sliceStmt = jen.Index().Add(jen.Qual(pkgPath, structTypeName))
		}
	} else {
		// Regular non-exported mode
		if isPointerSlice {
			// For []*T
			sliceStmt = jen.Index().Add(jen.Op("*").Id(structTypeName))
		} else {
			// For []T
			sliceStmt = jen.Index().Add(jen.Id(structTypeName))
		}
	}

	// Now create a slice with all matching references
	return sliceStmt.ValuesFunc(func(group *jen.Group) {
		// For each source ID
		for i := range srcValue.Len() {
			idValue := srcValue.Index(i).String()

			// Try to find a matching reference struct
			for j := range refData.Len() {
				refStruct := refData.Index(j)

				// Handle pointer to struct case
				if refStruct.Kind() == reflect.Pointer {
					refStruct = refStruct.Elem()
				}

				// Try each possible identifier field
				for _, idField := range g.IdentifierFields {
					refIDField := refStruct.FieldByName(idField)

					if refIDField.IsValid() &&
						refIDField.Kind() == reflect.String &&
						refIDField.String() == idValue {

						// Found a matching reference
						// Get a name for the referenced variable
						identValue := g.getStructIdentifier(refStruct)
						refVarName := structTypeName + slugToIdentifier(identValue)

						// Use a direct reference to the variable (e.g., TagGoProgramming)
						// For pointer slices, add the & operator
						if isPointerSlice {
							group.Add(jen.Op("&").Id(refVarName))
						} else {
							group.Add(jen.Id(refVarName))
						}
						break
					}
				}
			}
		}
	})
}

// generateReferenceSingle generates a single referenced struct for string to struct references
//
// This method handles the case where a field contains a string (e.g., "author-1")
// and needs to generate a struct (e.g., Author or *Author) by looking up the string in a reference dataset.
//
// Parameters:
//   - srcValue: The source field value (string)
//   - targetType: The target field type (struct or pointer to struct)
func (g *Generator) generateReferenceSingle(srcValue reflect.Value, targetType reflect.Type) *jen.Statement {
	// Determine if we're dealing with a pointer (*T) or struct (T)
	isPointer := targetType.Kind() == reflect.Pointer

	// Get the target struct type name
	var structTypeName string
	if isPointer {
		structTypeName = targetType.Elem().Name()
	} else {
		structTypeName = targetType.Name()
	}

	// Check if we have this reference type
	refDataObj, hasRef := g.Refs[structTypeName]
	if !hasRef {
		// We don't have this reference data
		if isPointer {
			return jen.Op("&").Id(structTypeName).Values()
		}
		return jen.Id(structTypeName).Values()
	}

	// Convert to reflect.Value
	refData := reflect.ValueOf(refDataObj)
	if refData.Kind() != reflect.Slice && refData.Kind() != reflect.Array {
		// Reference isn't a slice/array
		if isPointer {
			return jen.Op("&").Id(structTypeName).Values()
		}
		return jen.Id(structTypeName).Values()
	}

	// Get ID value from source
	idValue := srcValue.String()

	// Try to find a matching reference struct
	for j := range refData.Len() {
		refStruct := refData.Index(j)

		// Handle pointer to struct case
		if refStruct.Kind() == reflect.Pointer {
			refStruct = refStruct.Elem()
		}

		// Try each possible identifier field
		for _, idField := range g.IdentifierFields {
			refIDField := refStruct.FieldByName(idField)

			if refIDField.IsValid() &&
				refIDField.Kind() == reflect.String &&
				refIDField.String() == idValue {

				// Found match - get a name for the referenced variable
				identValue := g.getStructIdentifier(refStruct)
				refVarName := structTypeName + slugToIdentifier(identValue)

				// For pointer types, just return a pointer to the existing variable
				if isPointer {
					return jen.Op("&").Id(refVarName)
				}
				// For non-pointer types, return the variable directly
				return jen.Id(refVarName)
			}
		}
	}

	// No match found
	if isPointer {
		return jen.Op("&").Id(structTypeName).Values()
	}
	return jen.Id(structTypeName).Values()
}
