// Package main updates the database with new vault content.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rotisserie/eris"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/cache"
	"github.com/conneroisu/conneroh.com/internal/credited"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/genstruct"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/sync/errgroup"
)

const (
	hashFile    = "config.json"
	vaultLoc    = "internal/data/docs/"
	assetsLoc   = "internal/data/docs/assets/"
	postsLoc    = "internal/data/docs/posts/"
	tagsLoc     = "internal/data/docs/tags/"
	projectsLoc = "internal/data/docs/projects/"

	baseBucketEndpoint = "https://fly.storage.tigris.dev"
)

var (
	locations      = []string{assetsLoc, postsLoc, projectsLoc, tagsLoc}
	workers        = flag.Int("jobs", 20, "number of parallel uploads")
	cwd            = flag.String("cwd", "", "current working directory")
	debug          = flag.Bool("debug", false, "enable debug logging for cache operations")
	errInvalidFlag = eris.New("invalid flag value")
)

// Define parse and actualization result types
type parseResult struct {
	location string
	assets   []assets.Asset
	ignored  []string
	err      error
}

type actualizeResult struct {
	contentType string
	posts       []*gen.Post
	projects    []*gen.Project
	tags        []*gen.Tag
	err         error
}

func main() {
	flag.Parse()
	ctx := context.Background()

	if err := run(ctx, os.Getenv); err != nil {
		slog.Error("error", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string) (err error) {
	start := time.Now()
	if *cwd != "" {
		err = os.Chdir(*cwd)
		if err != nil {
			return eris.Wrapf(
				err,
				"failed to change working directory to %s",
				*cwd,
			)
		}
	}
	awsClient, err := credited.NewAWSClient(getenv)
	if err != nil {
		return err
	}
	ollama, err := credited.NewOllamaClient(getenv)
	if err != nil {
		return err
	}
	objCache, err := cache.LoadCache(hashFile)
	if err != nil {
		return err
	}
	mdParser := assets.NewMDParser(parser.NewContext(), awsClient)

	// Log cache stats
	slog.Info("cache loaded", "entries", len(objCache.Hashes))

	// Make sure cache is written at the end, regardless of how the function exits
	defer func() {
		// Log final cache stats and time
		slog.Info("saving cache",
			"entries", len(objCache.Hashes),
			"duration", time.Since(start).String())
		err = objCache.Close()
		if err != nil {
			err = eris.Wrap(
				err,
				"failed to close cache",
			)
		}
	}()

	// Create channels for communication between goroutines
	parseResultCh := make(chan parseResult, 4)           // 4 content types
	actualizeResultCh := make(chan actualizeResult, 3)   // 3 content types
	parsingComplete := make(chan map[string]parseResult) // Stop Signal

	// Start parsing all content types in parallel
	go func() {
		results := make(map[string]parseResult)

		// Use errgroup for better concurrency handling
		var eg errgroup.Group
		eg.SetLimit(*workers) // Limit concurrent parses

		// Mutex to protect the results map
		var mu sync.Mutex

		// Start parsing each location
		for _, loc := range locations {
			loc := loc // Capture for closure
			eg.Go(func() error {
				assets, ignored, egErr := parse(objCache, loc)

				// Store result safely
				mu.Lock()
				result := parseResult{
					location: loc,
					assets:   assets,
					ignored:  ignored,
					err:      egErr,
				}
				results[loc] = result
				mu.Unlock()

				// Push to channel for immediate processing (1/4 split)
				parseResultCh <- result
				return nil
			})
		}

		// Wait for all parsing to complete and notify via the signal channel
		go func() {
			_ = eg.Wait() // Ignore errors as they're captured in the results
			parsingComplete <- results
		}()
	}()

	// Flag to track when parsing is complete
	parseComplete := false

	// Track counts for actualization
	actualizationStarted := 0
	actualizationComplete := 0

	// Results storage
	var (
		allResults     map[string]parseResult
		parsedPosts    []*gen.Post
		parsedProjects []*gen.Project
		parsedTags     []*gen.Tag
	)

	// Process events until everything is complete
	// Exit condition: parsing is done and all actualizations are complete
	for !parseComplete || actualizationComplete != actualizationStarted || actualizationStarted <= 0 {

		// Process all available events
		select {
		case result := <-parseResultCh:
			// Process a parsing result
			if result.err != nil {
				slog.Error("error parsing location",
					"location", result.location,
					"error", result.err)
				continue
			}

			slog.Info("parsed location",
				"location", result.location,
				"assets", len(result.assets),
				"ignored", len(result.ignored))

			// Start actualization early if not assets
			if result.location != assetsLoc {
				actualizationStarted++
				go startActualization(ctx, ollama, mdParser, result, actualizeResultCh)
			}

		case results := <-parsingComplete:
			// All parsing is complete
			allResults = results
			parseComplete = true

			// Process assets immediately once parsing is complete
			assetResult, exists := results[assetsLoc]
			if exists && assetResult.err == nil {
				err = actualizeAssets(ctx, awsClient, assetResult.assets)
				if err != nil {
					return fmt.Errorf("failed to actualize assets: %w", err)
				}
			}

			// Check for parse errors
			for loc, result := range results {
				if result.err != nil {
					return fmt.Errorf("error parsing %s: %w", loc, result.err)
				}
			}

		case result := <-actualizeResultCh:
			// Process actualization result
			actualizationComplete++

			if result.err != nil {
				return fmt.Errorf("failed to actualize %s: %w", result.contentType, result.err)
			}

			switch result.contentType {
			case "posts":
				parsedPosts = result.posts
				slog.Info("posts actualized", "count", len(parsedPosts))
			case "projects":
				parsedProjects = result.projects
				slog.Info("projects actualized", "count", len(parsedProjects))
			case "tags":
				parsedTags = result.tags
				slog.Info("tags actualized", "count", len(parsedTags))
			}

		case <-time.After(3 * time.Second):
			// Periodic log to show we're still working
			slog.Info("still working...",
				"parseComplete", parseComplete,
				"actualizationStarted", actualizationStarted,
				"actualizationComplete", actualizationComplete)
		}
	}

	// Make sure we have results
	if allResults == nil {
		return fmt.Errorf("parsing completed but no results were collected")
	}

	// Get asset results for final log
	assetResult := allResults[assetsLoc]

	// Log final results
	slog.Info("Processing complete (generating)",
		"len(posts)", len(parsedPosts),
		"len(tags)", len(parsedTags),
		"len(projects)", len(parsedProjects),
		"len(assets)", len(assetResult.assets)+len(assetResult.ignored),
		"duration", time.Since(start).String())

	postGen, err := genstruct.NewGenerator(genstruct.Config{
		PackageName: "gen",
		OutputFile:  "internal/data/gen/generated_data.go",
	}, parsedPosts, parsedTags, parsedProjects)
	if err != nil {
		return err
	}

	err = postGen.Generate()
	if err != nil {
		return err
	}

	slog.Info("Generation complete", "duration", time.Since(start).String())
	return nil
}

// Helper function to start actualization of a content type
func startActualization(
	ctx context.Context,
	ollama *credited.OllamaClient,
	mdParser *assets.MDParser,
	result parseResult,
	resultCh chan<- actualizeResult,
) {
	switch result.location {
	case postsLoc:
		posts, err := actualize[gen.Post](
			ctx,
			ollama,
			mdParser,
			result.assets,
			result.ignored,
		)
		resultCh <- actualizeResult{
			contentType: "posts",
			posts:       posts,
			err:         err,
		}
	case projectsLoc:
		projects, err := actualize[gen.Project](
			ctx,
			ollama,
			mdParser,
			result.assets,
			result.ignored,
		)
		resultCh <- actualizeResult{
			contentType: "projects",
			projects:    projects,
			err:         err,
		}
	case tagsLoc:
		tags, err := actualize[gen.Tag](
			ctx,
			ollama,
			mdParser,
			result.assets,
			result.ignored,
		)
		resultCh <- actualizeResult{
			contentType: "tags",
			tags:        tags,
			err:         err,
		}
	}
}

func actualize[T gen.Post | gen.Tag | gen.Project](
	ctx context.Context,
	ollama *credited.OllamaClient,
	mdParser *assets.MDParser,
	contents []assets.Asset,
	ignored []string,
) ([]*T, error) {
	amount := len(contents)

	if amount == 0 {
		// If no new content, just return remembered content
		remembered := rememberMD[T](ignored)
		return remembered, nil
	}

	// Use errgroup for better error handling and concurrency control
	eg := errgroup.Group{}
	eg.SetLimit(*workers)

	// Create result slice with appropriate capacity
	parsedItems := make([]*T, amount)

	// Process each asset using errgroup
	for i, content := range contents {
		index, asset := i, content // Capture for closure

		eg.Go(func() error {
			realized, err := realizeMD[T](ctx, ollama, mdParser, asset)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to parse asset %d",
					index,
				)
			}

			parsedItems[index] = realized
			return nil
		})
	}

	// Wait for all goroutines to complete and check for errors
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// Include previously parsed items that were ignored this time
	remembered := rememberMD[T](ignored)
	parsedItems = append(parsedItems, remembered...)

	slog.Info(
		"actualization complete",
		"type",
		fmt.Sprintf("%T", *new(T)),
		"total",
		len(parsedItems),
	)
	return parsedItems, nil
}

func parse(
	cacheObj *cache.Cache,
	loc string,
) (parsedAssets []assets.Asset, ignoredSlugs []string, err error) {
	eg := &errgroup.Group{}
	eg.SetLimit(*workers)

	var (
		filePaths []string
		mu        sync.Mutex // Guards access to parsedAssets and ignoredSlugs
		processed int32
		skipped   int32
	)

	// First, collect all file paths
	err = filepath.WalkDir(
		loc,
		func(fPath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			filePaths = append(filePaths, fPath)
			return nil
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to walk fsPath (%s): %w", loc, err)
	}

	slog.Info("found files to process", "location", loc, "files", len(filePaths))

	// Create channels for immediate processing (1/4 split strategy)
	// This allows work to begin being processed before all files are collected
	filesCh := make(chan string, len(filePaths))
	go func() {
		for _, fPath := range filePaths {
			filesCh <- fPath
		}
		close(filesCh)
	}()

	// Process files as they come in
	for range *workers {
		eg.Go(func() error {
			for filePath := range filesCh {
				// Read file
				body, err := os.ReadFile(filePath)
				if err != nil {
					slog.Error("failed to read file", "path", filePath, "error", err)
					continue // Skip failed files but don't abort the entire process
				}

				// Calculate hash
				hash := cache.Hash(body)
				path := pathify(filePath)
				slug := slugify(filePath)

				mu.Lock()
				// Check if file has changed
				cachedHash, exists := cacheObj.Hashes[filePath]
				if exists && cachedHash == hash {
					// File hasn't changed
					ignoredSlugs = append(ignoredSlugs, slug)
					atomic.AddInt32(&skipped, 1)
					if *debug {
						slog.Debug("skipped unchanged file", "path", filePath)
					}
				} else {
					// File is new or has changed
					cacheObj.Hashes[filePath] = hash
					parsedAssets = append(parsedAssets, assets.Asset{
						Path: path,
						Data: body,
						Slug: slug,
					})
					atomic.AddInt32(&processed, 1)
					if *debug {
						if exists {
							slog.Debug("file changed", "path", filePath)
						} else {
							slog.Debug("new file", "path", filePath)
						}
					}
				}
				mu.Unlock()
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}

	slog.Info("parse statistics",
		"location", loc,
		"total", len(filePaths),
		"processed", processed,
		"skipped", skipped)

	return parsedAssets, ignoredSlugs, nil
}

func actualizeAssets(
	ctx context.Context,
	client *credited.AWSClient,
	assets []assets.Asset,
) error {
	amount := len(assets)
	if amount == 0 {
		slog.Info("no assets to upload")
		return nil
	}

	slog.Info("uploading assets", "amount", amount)

	// Use errgroup for better error handling and consistent patterns
	eg := errgroup.Group{}
	eg.SetLimit(*workers) // Limit concurrent uploads

	// Process assets concurrently using errgroup
	for _, asset := range assets {
		asset := asset // Capture for closure

		eg.Go(func() error {
			contentType := mime.TypeByExtension(filepath.Ext(asset.Path))
			slog.Info("uploading asset", "path", asset.Path, "contentType", contentType)

			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket:      aws.String("conneroh"),
				Key:         aws.String(asset.Path),
				Body:        bytes.NewReader(asset.Data),
				ContentType: aws.String(contentType),
			})

			if err != nil {
				// Log error but let errgroup handle it
				slog.Error("asset upload failed", "path", asset.Path, "error", err)
				return fmt.Errorf("failed to upload asset %s: %w", asset.Path, err)
			}

			return nil
		})
	}

	// Wait for all uploads to complete and check for errors
	err := eg.Wait()
	if err != nil {
		return err
	}

	slog.Info("assets upload complete", "total", amount)
	return nil
}

func rememberMD[T gen.Post | gen.Project | gen.Tag](ignored []string) []*T {
	ignoredMap := make(map[string]struct{}, len(ignored))
	for _, slug := range ignored {
		ignoredMap[slug] = struct{}{}
	}

	var typeExample T
	switch any(typeExample).(type) {
	case gen.Post:
		var posts []*gen.Post
		for _, post := range gen.AllPosts {
			if _, ok := ignoredMap[post.Slug]; ok {
				posts = append(posts, post)
			}
		}
		return any(posts).([]*T)
	case gen.Project:
		var projects []*gen.Project
		for _, project := range gen.AllProjects {
			if _, ok := ignoredMap[project.Slug]; ok {
				projects = append(projects, project)
			}
		}
		return any(projects).([]*T)
	case gen.Tag:
		var tags []*gen.Tag
		for _, tag := range gen.AllTags {
			if _, ok := ignoredMap[tag.Slug]; ok {
				tags = append(tags, tag)
			}
		}
		return any(tags).([]*T)
	default:
		return nil
	}
}

func realizeMD[T gen.Post | gen.Project | gen.Tag](
	ctx context.Context,
	ollama *credited.OllamaClient,
	mdParser *assets.MDParser,
	parsed assets.Asset,
) (*T, error) {
	var (
		emb gen.Embedded
		err error
	)

	// Convert markdown to HTML
	emb, err = mdParser.Convert(parsed)
	if err != nil {
		return nil, err
	}

	err = gen.Defaults(&emb)
	if err != nil {
		return nil, err
	}
	err = gen.Validate(&emb)
	if err != nil {
		return nil, err
	}
	err = ollama.Embeddings(ctx, emb.RawContent, &emb)
	if err != nil {
		return nil, err
	}

	return gen.New[T](&emb), nil
}

func slugify(s string) string {
	var path string
	var ok bool
	path, ok = strings.CutPrefix(s, assetsLoc)
	if ok {
		return path
	}
	return strings.TrimSuffix(pathify(s), filepath.Ext(s))
}

func pathify(s string) string {
	var path string
	var ok bool
	path, ok = strings.CutPrefix(s, postsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, projectsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, tagsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, assetsLoc)
	if ok {
		return path
	}
	return s // Return original path instead of panicking
}
