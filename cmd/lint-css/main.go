package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/conneroisu/twerge"
)

func main() {
	ctx := context.Background()

	// Check if a directory argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <directory>")
		os.Exit(1)
	}

	dirPath := os.Args[1]

	if err := run(ctx, dirPath); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, dirPath string) error {
	var (
		_ = layouts.Page(views.Home(
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)).Render(ctx, io.Discard)
		_ = layouts.Page(views.List(
			routing.PluralTargetPost,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
			"")).Render(ctx, io.Discard)
		_ = layouts.Page(views.List(
			routing.PluralTargetProject,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
			"",
		)).Render(ctx, io.Discard)
		_ = layouts.Page(views.List(
			routing.PluralTargetTag,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
			"",
		)).Render(ctx, io.Discard)
		_ = views.TagControl(
			&gen.Tag{},
		).Render(ctx, io.Discard)
		_ = layouts.Page(views.Post(
			gen.AllPosts[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)).Render(ctx, io.Discard)
		_ = layouts.Page(views.Project(
			gen.AllProjects[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)).Render(ctx, io.Discard)
		_ = layouts.Page(views.Tag(
			gen.AllTags[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)).Render(ctx, io.Discard)
		_ = views.TagControl(
			&gen.Tag{},
		).Render(ctx, io.Discard)
		_ = layouts.Morpher(views.Post(
			gen.AllPosts[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)).Render(ctx, io.Discard)
		_ = layouts.Layout("hello").Render(ctx, io.Discard)
	)

	// Print the original mappings
	fmt.Println("Original class mappings:")
	for gen, genMerged := range twerge.GenClassMergeStr {
		fmt.Printf("%s : %s\n", gen, genMerged)
	}

	fmt.Println("\nClass map strings:")
	for orig, class := range twerge.ClassMapStr {
		fmt.Printf("%s : %s\n", orig, class)
	}

	// Create a replacement map for HTML class replacements
	replacementMap := make(map[string]string)

	// Populate the replacement map with classes that need to be replaced
	fmt.Println("\nClasses to be replaced:")
	for gen, genMerged := range twerge.GenClassMergeStr {
		for orig, class := range twerge.ClassMapStr {
			if gen == class {
				if genMerged != orig {
					// For each original class that has been merged, add an entry to replace it with the merged version
					replacementMap[orig] = genMerged
					fmt.Printf("'%s' -> '%s'\n", orig, genMerged)
				}
			}
		}
	}

	// Count of files processed and modified
	filesProcessed := 0
	filesModified := 0

	// Walk through all files in the directory recursively
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .templ files
		if !strings.HasSuffix(path, ".templ") {
			return nil
		}

		filesProcessed++

		// Read the file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Convert to string for processing
		fileContent := string(content)
		originalContent := fileContent

		// Apply each replacement
		for orig, merged := range replacementMap {
			fileContent = strings.ReplaceAll(fileContent, "\""+orig+"\"", "\""+merged+"\"")
		}

		// Check if the content was modified
		if fileContent != originalContent {
			filesModified++

			// // diff the original and modified content
			// diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			// 	A:        difflib.SplitLines(originalContent),
			// 	B:        difflib.SplitLines(fileContent),
			// 	FromFile: "Original",
			// 	ToFile:   "Modified",
			// 	Context:  3,
			// })
			// if err != nil {
			// 	return err
			// }
			// fmt.Printf("Diff:\n%s\n", diff)

			// Write the modified content back to the file
			err = os.WriteFile(path, []byte(fileContent), info.Mode())
			if err != nil {
				return err
			}

			fmt.Printf("Modified file: %s\n", path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("\nSummary: Processed %d .templ files, modified %d files\n", filesProcessed, filesModified)

	return nil
}
