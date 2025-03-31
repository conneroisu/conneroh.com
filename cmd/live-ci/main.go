// Package main is the main package for the live-ci command.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/playwright-community/playwright-go"
)

func main() {
	// Set skip browser installation since we're using Nix-provided browsers
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: true,
	}

	// Install playwright driver (will use PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1)
	err := playwright.Install(runOption)
	if err != nil {
		log.Fatalf("could not install playwright dependencies: %v", err)
	}

	// Initialize playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	defer pw.Stop()

	// Launch the browser (using Nix-provided browser)
	option := playwright.BrowserTypeLaunchOptions{
		ExecutablePath: playwright.String(os.Getenv("PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH") + "/chrome-linux/chrome"),
		Headless:       playwright.Bool(false),
	}
	browser, err := pw.Chromium.Launch(option)
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	// Create a new browser context
	ctx, err := browser.NewContext()
	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}
	defer ctx.Close()

	// Create a new page in a browser context
	page, err := ctx.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	// Navigate to Go's package website
	if _, err = page.Goto("https://pkg.go.dev/"); err != nil {
		log.Fatalf("could not navigate to pkg.go.dev: %v", err)
	}

	// Take a screenshot of a specific element
	element, err := page.QuerySelector("img.Homepage-logo")
	if err != nil {
		log.Fatalf("could not get logo element: %v", err)
	}

	_, err = element.Screenshot(playwright.ElementHandleScreenshotOptions{
		Path: playwright.String("elementScreenshot.png"),
	})
	if err != nil {
		log.Fatalf("could not take element screenshot: %v", err)
	}

	// Navigate to another site for full page screenshot
	if _, err = page.Goto("https://brank.as/"); err != nil {
		log.Fatalf("could not navigate to brank.as: %v", err)
	}

	// Take a full page screenshot
	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String("fullScreenshot.png"),
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not take full screenshot: %v", err)
	}

	// Check if the files were created
	if _, err := os.Stat("elementScreenshot.png"); err == nil {
		fmt.Println("Successfully created elementScreenshot.png")
	}

	if _, err := os.Stat("fullScreenshot.png"); err == nil {
		fmt.Println("Successfully created fullScreenshot.png")
	}
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
