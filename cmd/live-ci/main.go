// Package main is the main package for the live-ci command.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/playwright-community/playwright-go"
	"github.com/rotisserie/eris"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	"golang.org/x/sync/errgroup"

	_ "modernc.org/sqlite"
)

var (
	url         = flag.String("url", "http://localhost:8080", "Base URL to visit")
	workers     = flag.Int("workers", 32, "Number of workers to use")
	lowestBytes = 5000
)

func main() {
	flag.Parse()
	ctx := context.Background()
	if err := run(ctx, *url, *workers); err != nil {
		log.Fatal(err)
	}
}

func run(
	ctx context.Context,
	base string,
	workers int,
) error {
	eg, _ := errgroup.WithContext(ctx)
	eg.SetLimit(workers)
	// Set skip browser installation since we're using Nix-provided browsers
	runOption := &playwright.RunOptions{SkipInstallBrowsers: true}

	// Install playwright driver (will use PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1)
	err := playwright.Install(runOption)
	if err != nil {
		return err
	}

	// Initialize playwright
	pw, err := playwright.Run()
	if err != nil {
		return err
	}
	// Launch the browser (using Nix-provided browser)
	option := playwright.BrowserTypeLaunchOptions{
		ExecutablePath: playwright.String(os.Getenv("PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH") + "/chrome-linux/chrome"),
		Headless:       playwright.Bool(false),
	}
	browser, err := pw.Chromium.Launch(option)
	if err != nil {
		return err
	}
	// Don't use defer with log.Fatal as it skips deferred functions
	defer func() {
		if closeErr := browser.Close(); closeErr != nil {
			err = eris.Wrap(closeErr, "could not close browser")
		}
	}()

	// Create a new browser context
	bCtx, err := browser.NewContext()
	if err != nil {
		return err
	}

	p, err := bCtx.NewPage()
	if err != nil {
		return err
	}
	sqlDB, err := sql.Open("sqlite", assets.DBName())
	if err != nil {
		return eris.Wrap(err, "error opening database")
	}
	defer sqlDB.Close()
	db := bun.NewDB(sqlDB, sqlitedialect.New())
	assets.RegisterModels(db)
	if os.Getenv("DEBUG") == "true" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	var (
		// Instance Caches
		allPosts    = []*assets.Post{}
		allProjects = []*assets.Project{}
		allTags     = []*assets.Tag{}
	)

	err = db.NewSelect().Model(&allPosts).
		Order("updated_at").
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Scan(ctx)
	if err != nil {
		return eris.Wrap(
			err,
			"failed to scan posts for home page",
		)
	}
	err = db.NewSelect().Model(&allProjects).
		Order("updated_at").
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Scan(ctx)
	if err != nil {
		return eris.Wrap(
			err,
			"failed to scan projects for home page",
		)
	}
	err = db.NewSelect().Model(&allTags).
		Order("updated_at").
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Scan(ctx)
	if err != nil {
		return eris.Wrap(
			err,
			"failed to scan tags for home page",
		)
	}

	for _, post := range allPosts {
		url := base + post.PagePath()
		eg.Go(func() error {
			var bdy string
			page, pErr := bCtx.NewPage()
			if pErr != nil {
				return pErr
			}
			resp, pErr := page.Goto(url, playwright.PageGotoOptions{})
			if pErr != nil {
				return pErr
			}
			bdy, pErr = resp.Text()
			if pErr != nil {
				return pErr
			}
			if len(bdy) < lowestBytes {
				return fmt.Errorf("body length too small: %d", len(bdy))
			}

			return page.Close()
		})
	}
	for _, project := range allProjects {
		url := base + project.PagePath()
		eg.Go(func() error {
			var bdy string
			page, pErr := bCtx.NewPage()
			if pErr != nil {
				return pErr
			}
			resp, pErr := page.Goto(url, playwright.PageGotoOptions{})
			if pErr != nil {
				return pErr
			}
			bdy, pErr = resp.Text()
			if pErr != nil {
				return pErr
			}
			if len(bdy) < lowestBytes {
				return fmt.Errorf("body length too small: %d", len(bdy))
			}

			return page.Close()
		})
	}
	for _, tag := range allTags {
		url := base + tag.PagePath()
		eg.Go(func() error {
			var bdy string
			page, pErr := bCtx.NewPage()
			if pErr != nil {
				return pErr
			}
			resp, pErr := page.Goto(url, playwright.PageGotoOptions{})
			if pErr != nil {
				return pErr
			}
			bdy, pErr = resp.Text()
			if pErr != nil {
				return pErr
			}
			if len(bdy) < lowestBytes {
				return fmt.Errorf("body length too small: %d", len(bdy))
			}

			return page.Close()
		})
	}

	err = eg.Wait()
	if err != nil {
		return err
	}

	if closeErr := p.Close(); closeErr != nil {
		return eris.Wrap(closeErr, "could not close page")
	}

	if closeErr := bCtx.Close(); closeErr != nil {
		return eris.Wrap(closeErr, "could not close browser context")
	}

	return pw.Stop()
}
