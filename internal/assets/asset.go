package assets

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"mime"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/tigris"
)

const (
	hashFile    = "config.json"
	vaultLoc    = "internal/data/docs/"
	assetsLoc   = "internal/data/docs/assets/"
	postsLoc    = "internal/data/docs/posts/"
	tagsLoc     = "internal/data/docs/tags/"
	projectsLoc = "internal/data/docs/projects/"
)

// Asset is a struct for embedding assets.
type Asset struct {
	Slug string
	Path string
	Data []byte
}

// NewAsset creates a new asset.
func NewAsset(path string, data []byte) *Asset {
	return &Asset{
		Slug: slugify(path),
		Path: pathify(path),
		Data: data,
	}
}

// URL returns the url of the asset.
func (a *Asset) URL() string {
	return fmt.Sprintf(
		"https://conneroh.fly.storage.tigris.dev/%s",
		a.Path,
	)
}

// Upload uploads the asset to S3.
func (a *Asset) Upload(
	ctx context.Context,
	client tigris.Client,
) error {
	extension := filepath.Ext(a.Path)
	if extension == "" {
		return fmt.Errorf("failed to get extension for %s", a.Path)
	}

	contentType := mime.TypeByExtension(extension)
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("conneroh"),
		Key:         aws.String(a.Path),
		Body:        bytes.NewReader(a.Data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		// Log error but let errgroup handle it
		slog.Error("asset upload failed", "path", a.Path, "error", err)
		return fmt.Errorf("failed to upload asset %s: %w", a.Path, err)
	}

	return nil
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
