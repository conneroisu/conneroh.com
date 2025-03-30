package markdown

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Asset is a struct for embedding assets.
type Asset struct {
	Path string
	Data []byte
}

// Upload uploads the asset to the s3 bucket.
func (a *Asset) Upload(ctx context.Context, client *s3.Client) error {
	mimeType := mime.TypeByExtension(filepath.Ext(a.Path))
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("BUCKET_NAME")),
		Key:         aws.String(a.Path),
		Body:        bytes.NewReader(a.Data),
		ContentType: aws.String(mimeType),
	})

	return err
}

// URL returns the url of the asset.
func (a *Asset) URL() string {
	return fmt.Sprintf(
		"https://conneroh.fly.storage.tigris.dev/%s",
		a.Path,
	)
}

// Filename returns the filename of the asset.
func (a *Asset) Filename() string {
	str := a.Path
	str = str[strings.LastIndex(str, "/")+1:]
	return str
}

// Resolvable returns true if the asset is resolvable.
func (a *Asset) Resolvable(path string) bool {
	return strings.HasPrefix(path, a.Path)
}
