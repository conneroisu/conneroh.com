package genstruct

import (
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
)

// getTypeStatement converts a reflect.Type to a jen.Statement
func (g *Generator) getTypeStatement(t reflect.Type) *jen.Statement {
	switch t.Kind() {
	case reflect.Bool:
		return jen.Bool()
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return jen.Id(t.String())
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return jen.Id(t.String())
	case reflect.Float32, reflect.Float64:
		return jen.Id(t.String())
	case reflect.Complex64, reflect.Complex128:
		return jen.Id(t.String())
	case reflect.Array, reflect.Slice:
		elemType := t.Elem()
		// Special handling for []*Type pattern
		if elemType.Kind() == reflect.Pointer {
			return jen.Index().Add(jen.Op("*").Add(g.getTypeStatement(elemType.Elem())))
		}
		return jen.Index().Add(g.getTypeStatement(elemType))
	case reflect.Map:
		return jen.Map(
			g.getTypeStatement(t.Key()),
		).Add(g.getTypeStatement(t.Elem()))
	case reflect.String:
		return jen.String()
	case reflect.Struct:
		// Handle special types like time.Time
		if t.String() == "time.Time" {
			return jen.Qual("time", "Time")
		}

		// Check if this is from a different package (has a dot in the name)
		pkgPath := t.PkgPath()
		// Infer ExportDataMode by checking if output file contains package path separator
		isExportMode := strings.Contains(g.OutputFile, "/")
		if pkgPath != "" && pkgPath != "main" && pkgPath != g.PackageName && isExportMode {
			// If the type comes from a different package, reference it with the package name
			pkgName := t.String()
			if lastDot := strings.LastIndex(pkgName, "."); lastDot >= 0 {
				return jen.Qual(pkgPath, t.Name())
			}
		}
		return jen.Id(t.Name())
	case reflect.Pointer:
		return jen.Op("*").Add(g.getTypeStatement(t.Elem()))
	case reflect.Interface:
		if t.NumMethod() == 0 {
			return jen.Interface() // empty interface
		}
		// Complex interfaces would need more handling
		return jen.Interface()
	default:
		return jen.Id(t.String())
	}
}
