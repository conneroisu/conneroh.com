package genstruct

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
)

// generateVariables creates variables for each struct
func (g *Generator) generateVariables(dataValue reflect.Value) {
	// Generate a variable for each struct
	for i := range dataValue.Len() {
		elem := dataValue.Index(i)

		// Determine the variable name using the identifier function
		identValue := g.getStructIdentifier(elem)
		varName := g.VarPrefix + slugToIdentifier(identValue)

		// Get the type to use (may be from another package)
		var typeStmt *jen.Statement

		// Check if this is a struct from another package
		var structType reflect.Type
		if elem.Kind() == reflect.Struct {
			structType = elem.Type()
		} else if elem.Kind() == reflect.Pointer && elem.Elem().Kind() == reflect.Struct {
			structType = elem.Elem().Type()
		}

		// If we have a struct type and it comes from a different package, use qualified name
		if structType != nil {
			pkgPath := structType.PkgPath()
			// Infer ExportDataMode by checking if output file contains package path separator
			isExportMode := strings.Contains(g.OutputFile, "/")
			if isExportMode && pkgPath != "" && pkgPath != "main" && pkgPath != g.PackageName {
				parts := strings.Split(g.TypeName, ".")
				if len(parts) > 1 {
					// If TypeName already has package qualifier (e.g., "pkg.Animal"), use it directly
					typeStmt = jen.Id(g.TypeName)
				} else {
					// Use package qualification
					typeStmt = jen.Qual(pkgPath, structType.Name())
				}
			} else {
				typeStmt = jen.Id(g.TypeName)
			}
		} else {
			typeStmt = jen.Id(g.TypeName)
		}

		// Create the variable with its value
		g.File.Var().Id(varName).Op("=").Add(typeStmt).ValuesFunc(func(group *jen.Group) {
			g.generateStructValues(group, elem)
		})
	}
}

// generateSlice creates a slice containing all struct instances
func (g *Generator) generateSlice(dataValue reflect.Value) {
	// Determine the slice name - handle both regular and irregular plurals
	var sliceName string
	if g.TypeName[len(g.TypeName)-1] == 's' ||
		g.TypeName[len(g.TypeName)-1] == 'x' ||
		g.TypeName[len(g.TypeName)-1] == 'z' ||
		strings.HasSuffix(g.TypeName, "sh") ||
		strings.HasSuffix(g.TypeName, "ch") {
		sliceName = fmt.Sprintf(
			"All%ses",
			g.TypeName,
		)
	} else if g.TypeName[len(g.TypeName)-1] == 'y' {
		sliceName = fmt.Sprintf(
			"All%sies",
			g.TypeName[:len(g.TypeName)-1],
		)
	} else {
		sliceName = fmt.Sprintf("All%ss", g.TypeName)
	}

	// Get the type to use (may be from another package)
	var typeStmt *jen.Statement
	var elemType reflect.Type

	// If we have at least one element, use it to determine the type
	if dataValue.Len() > 0 {
		elem := dataValue.Index(0)

		if elem.Kind() == reflect.Struct {
			elemType = elem.Type()
		} else if elem.Kind() == reflect.Pointer && elem.Elem().Kind() == reflect.Struct {
			elemType = elem.Elem().Type()
		}
	}

	// If we have a struct type and it comes from a different package, use qualified name
	if elemType != nil {
		pkgPath := elemType.PkgPath()
		// Infer ExportDataMode by checking if output file contains package path separator
		isExportMode := strings.Contains(g.OutputFile, "/")
		if isExportMode &&
			pkgPath != "" &&
			pkgPath != "main" &&
			pkgPath != g.PackageName {

			parts := strings.Split(g.TypeName, ".")
			if len(parts) > 1 {
				// If TypeName already has package qualifier (e.g., "pkg.Animal"), use it directly
				typeStmt = jen.Id(g.TypeName)
			} else {
				// Use package qualification
				typeStmt = jen.Qual(pkgPath, elemType.Name())
			}
		} else {
			typeStmt = jen.Id(g.TypeName)
		}
	} else {
		typeStmt = jen.Id(g.TypeName)
	}

	// Generate as pointer slice []*Type with &Var references
	g.File.Var().Id(
		sliceName,
	).Op(
		"=",
	).Index().Op("*").Add(
		typeStmt,
	).ValuesFunc(func(group *jen.Group) {
		for i := range dataValue.Len() {
			elem := dataValue.Index(i)

			// Get the variable name using the same method as in generateVariables
			identValue := g.getStructIdentifier(elem)
			varName := g.VarPrefix + slugToIdentifier(identValue)

			// Add & operator to create pointer references
			group.Op("&").Id(varName)
		}
	})
}
