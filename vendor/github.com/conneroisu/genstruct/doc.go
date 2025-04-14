// Package genstruct provides tooling for generating static Go code from struct data.
//
// The generator takes data in the form of struct slices and generates type-safe constants,
// variables, and slices for use in Go code. It supports referencing between different
// struct types using the structgen tag.
//
// Basic usage:
//
//	generator := genstruct.NewGenerator(
//	    genstruct.WithPackageName("zoo"),
//	    genstruct.WithOutputFile("animals.go"),
//	)
//	err := generator.Generate(animals)
//
// For struct references:
//
//	generator := genstruct.NewGenerator(
//	    genstruct.WithPackageName("blog"),
//	    genstruct.WithOutputFile("blog.go"),
//	)
//	err := generator.Generate(posts, tags)
package genstruct

//go:generate gomarkdoc -o README.md -e .
