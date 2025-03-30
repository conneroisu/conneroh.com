// Package main updates the database with new vault content.
package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/data/docs"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/markdown"
	"github.com/conneroisu/genstruct"
	"golang.org/x/sync/errgroup"
)

const (
	uploadJobs = 10
)

func main() {
	if err := Run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Run parses all markdown files in the database.
func Run(ctx context.Context) error {
	client := s3.NewFromConfig(aws.Config{
		Region:       "auto",
		BaseEndpoint: aws.String("https://fly.storage.tigris.dev"),
		Credentials: &markdown.CredHandler{
			Name: "conneroh",
			ID:   os.Getenv("AWS_ACCESS_KEY_ID"),
			Key:  os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	})
	assets, err := assetsParse(docs.Assets)
	if err != nil {
		return err
	}
	eg := errgroup.Group{}
	eg.SetLimit(uploadJobs)
	for _, asset := range assets {
		eg.Go(func() error {
			return asset.Upload(ctx, client)
		})
	}
	err = eg.Wait()
	if err != nil {
		return err
	}
	parsedTags, err := pathParse[gen.Tag](ctx, assets, "tags", docs.Tags)
	if err != nil {
		return err
	}
	parsedPosts, err := pathParse[gen.Post](ctx, assets, "posts", docs.Posts)
	if err != nil {
		return err
	}
	parsedProjects, err := pathParse[gen.Project](ctx, assets, "projects", docs.Projects)
	if err != nil {
		return fmt.Errorf("failed to parse projects: %v", err)
	}

	postGen, err := genstruct.NewGenerator(genstruct.Config{
		PackageName: "gen",
		OutputFile:  "internal/data/gen/generated_data.go",
	}, parsedPosts, parsedTags, parsedProjects)
	if err != nil {
		return err
	}

	return postGen.Generate()
}

// pathParse parses the markdown files in the given path.
func pathParse[T gen.Post | gen.Project | gen.Tag](
	ctx context.Context,
	assets []markdown.Asset,
	fsPath string,
	embedFs embed.FS,
) ([]T, error) {
	var parseds []T
	err := fs.WalkDir(
		embedFs,
		fsPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf(
					"failed to walk fsPath (%s): %w",
					fsPath,
					err,
				)
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			parsed, err := markdown.Parse[T](
				ctx,
				path,
				embedFs,
				markdown.NewParser(assets),
			)
			if err != nil {
				return err
			}
			if parsed == nil {
				return nil
			}
			parseds = append(parseds, *parsed)
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to walk fsPath (%s): %w", fsPath, err)
	}
	return parseds, nil
}

func assetsParse(
	embedFs embed.FS,
) ([]markdown.Asset, error) {
	var assets []markdown.Asset
	err := fs.WalkDir(
		embedFs,
		"assets",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf(
					"failed to walk fsPath (%s): %w",
					"assets",
					err,
				)
			}
			if d.IsDir() {
				return nil
			}
			asset, err := embedFs.ReadFile(path)
			if err != nil {
				return err
			}
			assets = append(assets, markdown.Asset{
				Path: strings.TrimPrefix(path, "assets/"),
				Data: asset,
			})
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return assets, nil
}
