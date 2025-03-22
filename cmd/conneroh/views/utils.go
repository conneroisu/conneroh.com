// Package views contains the HTML templates for the website.
package views

//go:generate gomarkdoc -o README.md -e .

import (
	"fmt"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/internal/data/master"
	"github.com/conneroisu/conneroh.com/internal/routing"
)

// ListFn returns a fullFn for the list view.
func ListFn(target ListTargets) routing.FullFn {
	return func(
		fullPosts *[]master.FullPost,
		fullProjects *[]master.FullProject,
		fullTags *[]master.FullTag,
		fullPostSlugMap *map[string]master.FullPost,
		fullProjectSlugMap *map[string]master.FullProject,
		fullTagSlugMap *map[string]master.FullTag,
	) templ.Component {
		return List(
			target,
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
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

func formatTags(tags []master.Tag) string {
	// This is a simplified representation - in a real app, you might want to use a JSON library
	if len(tags) == 0 {
		return "[]"
	}

	result := "["
	for i, tag := range tags {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf(`{"id":%d,"name":"%s","slug":"%s"}`, tag.ID, tag.Title, tag.Slug)
	}
	result += "]"
	return result
}

// Helper function to determine the section for a tag
func getTagSection(tagName string) string {
	if strings.Contains(tagName, "/") {
		return strings.Split(tagName, "/")[0]
	}
	return "misc"
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
