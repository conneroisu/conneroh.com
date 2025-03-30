// Package main is the main package for the live-ci command.
package main

import "github.com/conneroisu/conneroh.com/internal/data/gen"

func main() {
	println("hello world")
	for _, post := range gen.AllPosts {
		println(post.Title)
		//TODO: Assert that visiting the post URL returns a 200 OK
	}
	for _, project := range gen.AllProjects {
		println(project.Title)
		// TODO: Assert that visiting the project URL returns a 200 OK
	}
	for _, tag := range gen.AllTags {
		println(tag.Title)
		// TODO: Assert that visiting the tag URL returns a 200 OK
	}
}
