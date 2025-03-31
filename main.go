// Package main contains the main function for the conneroh.com website.
package main

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/conneroisu/conneroh.com/cmd/conneroh"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
)

func main() {
	var post *gen.Post
	var project *gen.Project
	var tag *gen.Tag
	for _, tag = range gen.AllTags {
		for _, post = range gen.AllPosts {
			// for each post that has this tag add it the tag's posts
			if slices.Contains(post.TagSlugs, tag.Slug) {
				tag.Posts = append(tag.Posts, post)
			}
		}
		for _, project = range gen.AllProjects {
			// for each project that has this tag add it the tag's projects
			if slices.Contains(project.TagSlugs, tag.Slug) {
				tag.Projects = append(tag.Projects, project)
			}
		}
	}
	if err := conneroh.Run(
		context.Background(),
		os.Getenv,
	); err != nil {
		fmt.Println(err)
		return
	}
}
