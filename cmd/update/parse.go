package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/conneroisu/conneroh.com/internal/data/docs"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

// pathParse parses the markdown files in the given path.
func pathParse(fsPath string, embedFs embed.FS) ([]*Markdown, error) {
	var parseds []*Markdown
	err := fs.WalkDir(
		embedFs,
		fsPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("failed to walk projects: %w", err)
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			parsed, err := parse(path, embedFs)
			if err != nil {
				return fmt.Errorf("failed to parse project %s: %w", path, err)
			}
			if parsed == nil {
				return nil
			}
			parseds = append(parseds, parsed)
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update projects: %v", err)
	}
	return parseds, nil
}

// parse parses the markdown file at the given path.
func parse(fsPath string, embedFs embed.FS) (*Markdown, error) {
	if filepath.Ext(fsPath) != ".md" {
		return nil, nil
	}

	var (
		ctx      = parser.NewContext()
		fm       Markdown
		metadata = frontmatter.Get(ctx)
		body     []byte
		buf      = bytes.NewBufferString("")
		err      error
	)

	body, err = embedFs.ReadFile(fsPath)
	if err != nil {
		return nil, err
	}
	err = md.Convert(body, buf, parser.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	metadata = frontmatter.Get(ctx)
	if metadata == nil {
		return nil, fmt.Errorf("frontmatter is nil for %s", fsPath)
	}
	err = metadata.Decode(&fm)
	if err != nil {
		return nil, err
	}

	switch embedFs {
	case docs.Posts:
		fsPath = strings.Replace(fsPath, "posts/", "", 1)
	case docs.Tags:
		fsPath = strings.Replace(fsPath, "tags/", "", 1)
	case docs.Projects:
		fsPath = strings.Replace(fsPath, "projects/", "", 1)
	default:
		return nil, fmt.Errorf("unknown embedFs %v", embedFs)
	}
	fsPath = strings.TrimSuffix(fsPath, filepath.Ext(fsPath))
	fm.Slug = fsPath
	fm.RenderContent = buf.String()
	if fm.Description == "" {
		return nil, fmt.Errorf("description is empty for %s", fsPath)
	}
	fm.RawContent = string(body)

	if fm.Icon == "" {
		fm.Icon = "tag"
	}

	return &fm, nil
}
