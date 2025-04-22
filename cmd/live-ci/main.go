// Package main is the main package for the live-ci command.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
	"golang.org/x/sync/errgroup"
)

var (
	url     = flag.String("url", "http://localhost:8080", "Base URL to visit")
	workers = flag.Int("workers", 1, "Number of workers to use")
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
	defer func() {
		stopErr := pw.Stop()
		if stopErr != nil {
			log.Fatalf("could not stop playwright: %v", stopErr)
		}
	}()

	// Launch the browser (using Nix-provided browser)
	option := playwright.BrowserTypeLaunchOptions{
		ExecutablePath: playwright.String(os.Getenv("PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH") + "/chrome-linux/chrome"),
		Headless:       playwright.Bool(false),
	}
	browser, err := pw.Chromium.Launch(option)
	if err != nil {
		return err
	}
	defer func() {
		err = browser.Close()
		if err != nil {
			log.Fatalf("could not close browser: %v", err)
		}
	}()

	// Create a new browser context
	bCtx, err := browser.NewContext()
	if err != nil {
		return err
	}
	defer func() {
		err = bCtx.Close()
		if err != nil {
			log.Fatalf("could not close browser context: %v", err)
		}
	}()

	p, err := bCtx.NewPage()
	if err != nil {
		return err
	}
	defer func() {
		err = p.Close()
		if err != nil {
			log.Fatalf("could not close page: %v", err)
		}
	}()

	time.Sleep(time.Second * 40)

	// for _, url := range routing.ComputeAllURLs(base) {
	// 	eg.Go(func() error {
	// 		page, pErr := bCtx.NewPage()
	// 		if pErr != nil {
	// 			return pErr
	// 		}
	// 		defer func() {
	// 			err = page.Close()
	// 			if err != nil {
	// 				log.Fatalf("could not close page: %v", err)
	// 			}
	// 		}()
	// 		if _, pErr = page.Goto(url, playwright.PageGotoOptions{}); pErr != nil {
	// 			return pErr
	// 		}
	// 		time.Sleep(time.Second)
	// 		slog.Info("visited post", "url", url)
	// 		return nil
	// 	})
	// }
	err = eg.Wait()
	if err != nil {
		return err
	}
	return nil
}
