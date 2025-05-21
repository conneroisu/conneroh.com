package assets

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/tigris"
	"github.com/rotisserie/eris"
)

const (
	defaultS3Timeout = 20 * time.Second
)

// UploadToS3 uploads a file to S3.
func UploadToS3(
	ctx context.Context,
	client tigris.Client,
	bucket string,
	path string,
	data []byte,
) error {
	// Use custom content type function instead of mime.TypeByExtension
	contentType := GetContentType(path)

	timeout := defaultS3Timeout

	if isVideoType(contentType) {
		return nil
	}
	// Use a timeout context to prevent hanging on S3 operations
	uploadCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	slog.Info("Uploading to S3", "path", path)
	_, err := client.PutObject(uploadCtx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         aws.String(path),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return eris.Wrapf(err, "S3 upload timed out: %s", path)
		}

		return eris.Wrapf(err, "failed to upload to S3: %s", path)
	}

	slog.Debug("uploaded to S3", "path", path)

	return nil
}
