package genstruct

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
)

// generateConstants creates ID constants for each struct if an ID field exists
func (g *Generator) generateConstants(dataValue reflect.Value) {
	var (
		hasIDField  bool
		idFieldName string
	)

	// Check if the struct has an ID field
	firstElem := dataValue.Index(0)
	// Handle pointer to struct case
	if firstElem.Kind() == reflect.Pointer {
		firstElem = firstElem.Elem()
	}

	// Look for an "ID" field (case insensitive)
	for i := range firstElem.NumField() {
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
	g.File.Const().DefsFunc(func(group *jen.Group) {
		for i := range dataValue.Len() {
			elem := dataValue.Index(i)
			// Handle pointer to struct case
			if elem.Kind() == reflect.Pointer {
				elem = elem.Elem()
			}

			idField := elem.FieldByName(idFieldName)

			// If there's an ID field that's a string, create a constant
			if idField.IsValid() &&
				idField.Kind() == reflect.String {

				idValue := idField.String()
				// If ID is empty, generate one
				if idValue == "" {
					idValue = fmt.Sprintf("%s-%d", strings.ToLower(g.TypeName), i+1)
				}

				// Get a name for the constant based on the struct
				identValue := g.getStructIdentifier(elem)

				constName := g.ConstantIdent + slugToIdentifier(identValue) + "ID"
				group.Id(constName).Op("=").Lit(idValue)
			}
		}
	})
}
