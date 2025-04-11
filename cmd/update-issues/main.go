// Package main contains the main function for the updating github issue templates.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Paths
const (
	docsPath     = "./internal/data/docs"
	templatePath = ".github/ISSUE_TEMPLATE/ideology_request.md"
)

// Frontmatter represents the YAML frontmatter in markdown files
type Frontmatter struct {
	Tags []string `yaml:"tags"`
}

func main() {
	fmt.Println("Extracting tags from markdown files...")
	tags, err := extractTags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting tags: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d unique tags.\n", len(tags))

	fmt.Println("Updating ideology request template...")
	if err := updateTemplate(tags); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating template: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Process completed successfully.")
}

// extractTags finds all markdown files and extracts unique tags from their frontmatter
func extractTags() ([]string, error) {
	tagSet := make(map[string]struct{})

	err := filepath.WalkDir(docsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		// Read the file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		// Extract frontmatter
		frontmatter, err := extractFrontmatter(string(content))
		if err != nil {
			fmt.Printf("Warning: Error parsing frontmatter in %s: %v\n", path, err)
			return nil
		}

		// Add tags to set
		for _, tag := range frontmatter.Tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tagSet[tag] = struct{}{}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	// Convert map to sorted slice
	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	return tags, nil
}

// extractFrontmatter parses YAML frontmatter from markdown content
func extractFrontmatter(content string) (Frontmatter, error) {
	var frontmatter Frontmatter

	// Regular expression to match frontmatter
	re := regexp.MustCompile(`(?s)^---\n(.*?)\n---`)
	match := re.FindStringSubmatch(content)

	if len(match) < 2 {
		return frontmatter, fmt.Errorf("no frontmatter found")
	}

	// Parse YAML
	err := yaml.Unmarshal([]byte(match[1]), &frontmatter)
	if err != nil {
		return frontmatter, fmt.Errorf("error parsing YAML: %w", err)
	}

	return frontmatter, nil
}

// updateTemplate updates the template file with the extracted tags
func updateTemplate(tags []string) error {
	// Read the template file
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	// Create the tags section content
	var tagsSection string
	if len(tags) > 0 {
		tagLines := make([]string, len(tags))
		for i, tag := range tags {
			tagLines[i] = "- " + tag
		}
		tagsSection = strings.Join(tagLines, "\n")
	} else {
		tagsSection = "No existing tags found."
	}

	// Replace the existing tags section
	re := regexp.MustCompile(`(?s)<existing-tags>.*?</existing-tags>`)
	updatedContent := re.ReplaceAllString(string(content), "<existing-tags>\n"+tagsSection+"\n</existing-tags>")

	// Write the updated content back to the file
	err = os.WriteFile(templatePath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing updated template: %w", err)
	}

	fmt.Printf("Successfully updated template with %d tags.\n", len(tags))
	fmt.Println("Tags:", tags)

	return nil
}
