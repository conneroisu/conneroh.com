// Package markdown parses markdown files.
package markdown

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/conneroisu/conneroh.com/internal/data/docs"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

// Parse parses the markdown file at the given path.
func Parse[T gen.Post | gen.Project | gen.Tag](
	ctx context.Context,
	fsPath string,
	embedFs embed.FS,
	md goldmark.Markdown,
) (*T, error) {
	var (
		pCtx     = parser.NewContext()
		fm       gen.Embedded
		metadata = frontmatter.Get(pCtx)
		body     []byte
		buf      = bytes.NewBufferString("")
		err      error
	)

	body, err = embedFs.ReadFile(fsPath)
	if err != nil {
		return nil, err
	}
	err = md.Convert(body, buf, parser.WithContext(pCtx))
	if err != nil {
		return nil, err
	}
	metadata = frontmatter.Get(pCtx)
	if metadata == nil {
		return nil, fmt.Errorf("frontmatter is nil for %s", fsPath)
	}
	err = metadata.Decode(&fm)
	if err != nil {
		return nil, fmt.Errorf("failed to decode frontmatter: %w", err)
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
	fm.Content = buf.String()
	if fm.Description == "" {
		return nil, fmt.Errorf("description is empty for %s", fsPath)
	}
	fm.RawContent = string(body)

	if fm.Icon == "" {
		fm.Icon = "tag"
	}

	if fm.Content == "" {
		return nil, fmt.Errorf("content is empty for %s", fsPath)
	}
	fm.Vec, fm.X, fm.Y, fm.Z, err = TextEmbeddingCreate(
		ctx,
		client,
		fm.RawContent,
	)
	if err != nil {
		return nil, err
	}

	return gen.New[T](fm), nil
}
