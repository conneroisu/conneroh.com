// Package views contains the HTML templates for the website.
package views

//go:generate gomarkdoc -o README.md -e .

import (
	"fmt"
	"time"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/internal/data/master"
	"github.com/conneroisu/conneroh.com/internal/routing"
)

func Single(
	posts *[]master.FullPost, projects *[]master.FullProject, tags *[]master.FullTag, postSlugMap *map[string]master.FullPost, projectSlugMap *map[string]master.FullProject, tagSlugMap *map[string]master.FullTag,
) func(
	target routing.SingleTarget,
	id string,
) templ.Component {
	return func(
		target routing.SingleTarget,
		id string,
	) templ.Component {
		return single(
			target,
			id,
			posts,
			projects,
			tags,
			postSlugMap,
			projectSlugMap,
			tagSlugMap,
		)
	}
}

func readTime(content string) string {
	// Rough estimate - 200 words per minute reading speed
	words := len(content) / 5 // Average word length is 5 characters
	minutes := words / 200

	if minutes < 1 {
		return "1"
	}
	return fmt.Sprintf("%d", minutes)
}

// Helper functions for formatting data for Alpine.js
func formatDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("Jan 02, 2006")
}

// Helper function to check if two tags have common posts or projects
func haveCommonItems(tag1 *master.FullTag, tag2 *master.FullTag) bool {
	// Check for common posts
	for _, post := range tag1.Posts {
		for _, relatedPost := range tag2.Posts {
			if post.ID == relatedPost.ID {
				return true
			}
		}
	}

	// Check for common projects
	for _, project := range tag1.Projects {
		for _, relatedProject := range tag2.Projects {
			if project.ID == relatedProject.ID {
				return true
			}
		}
	}

	return false
}
