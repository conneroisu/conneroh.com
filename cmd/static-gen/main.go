package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conneroisu/conneroh.com/cmd/conneroh"
	gen "github.com/conneroisu/genstruct"
)

func main() {
	ctx := context.Background()
	err := Run(ctx, os.Getenv)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Run parses all markdown files in the database.
func Run(
	ctx context.Context,
	getenv func(string) string,
) error {
	db, err := conneroh.NewDb(getenv)
	if err != nil {
		return err
	}

	// Generate Posts
	posts, err := db.Queries.PostsList(ctx)
	if err != nil {
		return err
	}

	fullPosts, err := db.Queries.FullPostsList(ctx, posts)
	if err != nil {
		return err
	}

	// Create a generator for people
	postGen := gen.NewGenerator(gen.Config{
		PackageName:   "master",
		TypeName:      "FullPost",
		ConstantIdent: "FullPost",
		VarPrefix:     "FullPost",
		OutputFile:    "internal/data/master/posts.go",
	}, *fullPosts)

	// Generate code for people
	err = postGen.Generate()
	if err != nil {
		fmt.Println("Error generating person code:", err)
		os.Exit(1)
	}

	// Generate Projects
	projects, err := db.Queries.ProjectsList(ctx)
	if err != nil {
		return err
	}

	fullProjects, err := db.Queries.FullProjectsList(ctx, projects)
	if err != nil {
		return err
	}

	projectGen := gen.NewGenerator(gen.Config{
		PackageName:   "master",
		TypeName:      "FullProject",
		ConstantIdent: "FullProject",
		VarPrefix:     "FullProject",
		OutputFile:    "internal/data/master/projects.go",
	}, *fullProjects)

	err = projectGen.Generate()
	if err != nil {
		fmt.Println("Error generating project code:", err)
		os.Exit(1)
	}

	// Generate Tags
	tags, err := db.Queries.TagsListAlphabetical(ctx)
	if err != nil {
		return err
	}

	fullTags, err := db.Queries.FullTagsList(ctx, tags)
	if err != nil {
		return err
	}

	tagGen := gen.NewGenerator(gen.Config{
		PackageName:   "master",
		TypeName:      "FullTag",
		ConstantIdent: "FullTag",
		VarPrefix:     "FullTag",
		OutputFile:    "internal/data/master/tags.go",
	}, *fullTags)

	err = tagGen.Generate()
	if err != nil {
		fmt.Println("Error generating tag code:", err)
		os.Exit(1)
	}

	fmt.Println("Code generation completed successfully!")

	return nil
}
