package main

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"strings"

	"github.com/conneroisu/conneroh.com/cmd/conneroh"
	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/docs"
	"github.com/ollama/ollama/api"
)

func main() {
	ctx := context.Background()
	err := Run(ctx, os.Getenv)
	if err != nil {
		slog.Error("failed to run update", "err", err)
		os.Exit(1)
	}
}

// Run parses all markdown files in the database.
func Run(ctx context.Context, getenv func(string) string) error {
	db, err := conneroh.NewDb(os.Getenv)
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
			slog.Info("parsing tag", "path", path)
			parsed, err := data.Parse(path, true)
			if err != nil {
				return err
			}
			return parsed.UpsertTag(ctx, db, client)
		},
	)
	if err != nil {
		return err
	}

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
			slog.Info("parsing post", "path", path)
			parsed, err := data.Parse(path, false)
			if err != nil {
				return err
			}
			return parsed.UpsertPost(ctx, db, client)
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
			slog.Info("parsing project", "path", path)
			parsed, err := data.Parse(path, false)
			if err != nil {
				return err
			}
			return parsed.UpsertProject(ctx, db, client)
		},
	)
	if err != nil {
		return err
	}
	return nil
}
