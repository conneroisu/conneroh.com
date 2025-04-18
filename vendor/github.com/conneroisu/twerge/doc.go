// Package twerge provides a tailwind merger for go-templ with class generation and runtime static hashmap.
//
// It performs four key functions:
// 1. Merges TailwindCSS classes intelligently (resolving conflicts)
// 2. Generates short unique CSS class names from the merged classes
// 3. Creates a mapping from original class strings to generated class names
// 4. Provides a runtime static hashmap for direct class name lookup
//
// Basic Usage:
//
//	import "github.com/conneroisu/twerge"
//
//	// Merge TailwindCSS classes from a space-delimited string
//	merged := twerge.Merge("text-red-500 bg-blue-500 text-blue-700")
//	// Returns: "bg-blue-500 text-blue-700"
//
//	// Generate a short unique class name
//	className := twerge.Generate("text-red-500 bg-blue-500")
//	// Returns something like: "tw-Ab3F5g7"
package twerge

//go:generate gomarkdoc -o README.md -e .
