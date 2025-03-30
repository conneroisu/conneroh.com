// Package views contains the HTML templates for the website.
package views

//go:generate gomarkdoc -o README.md -e .

import (
	"fmt"
)

func readTime(content string) string {
	// Rough estimate - 200 words per minute reading speed
	words := len(content) / 5 // Average word length is 5 characters
	minutes := words / 200

	if minutes < 1 {
		return "1"
	}
	return fmt.Sprintf("%d", minutes)
}
