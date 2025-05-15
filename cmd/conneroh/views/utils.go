// Package views contains the HTML templates for the website.
package views

import "strconv"

//go:generate gomarkdoc -o README.md -e .

func readTime(content string) string {
	// Rough estimate - 200 words per minute reading speed
	words := len(content) / 5 // Average word length is 5 characters
	minutes := words / 200

	if minutes < 1 {
		return "1"
	}

	return strconv.Itoa(minutes)
}
