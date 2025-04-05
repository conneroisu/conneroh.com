package assets

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"mime"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/credited"
)

// Asset is a struct for embedding assets.
type Asset struct {
	Slug string
	Path string
	Data []byte
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
	client *credited.AWSClient,
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
