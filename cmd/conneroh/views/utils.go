// Package views contains the HTML templates for the website.
package views

//go:generate gomarkdoc -o README.md -e .

import (
	"fmt"

	"github.com/conneroisu/conneroh.com/internal/data/gen"
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

// Helper function to check if two tags have common posts or projects
func haveCommonItems(tag1 *gen.Tag, tag2 *gen.Tag) bool {
	// Check for common posts
	for _, post := range tag1.Posts {
		for _, relatedPost := range tag2.Posts {
			if post.Slug == relatedPost.Slug {
				return true
			}
		}
	}

	// Check for common projects
	for _, project := range tag1.Projects {
		for _, relatedProject := range tag2.Projects {
			if project.Slug == relatedProject.Slug {
				return true
			}
		}
	}

	return false
}
