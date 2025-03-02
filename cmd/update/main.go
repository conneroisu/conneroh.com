// Package main updates the database with new vault content.
package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"strings"

	"github.com/conneroisu/conneroh.com/cmd/conneroh"
	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/docs"
	"github.com/conneroisu/conneroh.com/internal/data/master"
	"github.com/ollama/ollama/api"
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
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return err
	}

	err = fs.WalkDir(
		docs.Tags,
		"tags",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			parsed, err := data.Parse(path, docs.Tags)
			if err != nil {
				slog.Error("failed to parse tag", "path", path, "err", err)
				return err
			}
			tag, err := parsed.UpsertTag(ctx, db, client)
			if err != nil {
				return err
			}
			slog.Info("upserted tag", "tag", tag.Title)
			return nil
		},
	)
	if err != nil {
		return err
	}

	var post master.Post
	err = fs.WalkDir(
		docs.Posts,
		"posts",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			parsed, err := data.Parse(path, docs.Posts)
			if err != nil {
				return err
			}
			post, err = parsed.UpsertPost(ctx, db, client)
			if err != nil {
				return err
			}
			return db.Queries.UpsertPostTags(ctx, parsed.Tags, post.ID)
		},
	)
	if err != nil {
		return err
	}

	err = fs.WalkDir(
		docs.Projects,
		"projects",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			parsed, err := data.Parse(path, docs.Projects)
			if err != nil {
				return err
			}
			project, err := parsed.UpsertProject(ctx, db, client)
			if err != nil {
				return err
			}
			return db.Queries.UpsertProjectTags(ctx, parsed.Tags, project.ID)
		},
	)
	if err != nil {
		return err
	}
	return nil
}
