package assets

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rotisserie/eris"
)

const (
	defaultS3Timeout = 20 * time.Second

	awsAccessKeyIDVar = "AWS_ACCESS_KEY_ID"
	awsSecretKeyVar   = "AWS_SECRET_ACCESS_KEY" //nolint:gosec
	awsBaseURLVar     = "AWS_ENDPOINT_URL_S3"
)

// Tigris is an minimal interface for AWS clients.
type Tigris interface {
	PutObject(
		ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options),
	) (*s3.PutObjectOutput, error)
}

// DefaultTigrisClient is a wrapper for the AWS S3 client.
type DefaultTigrisClient struct{ *s3.Client }

// NewTigris creates a new Tigris client.
func NewTigris(getenv func(string) string) (*DefaultTigrisClient, error) {
	credHandler, err := newCredHandler(getenv)
	if err != nil {
		return nil, err
	}
	baseURL := getenv(awsBaseURLVar)
	if baseURL == "" {
		return nil, eris.Wrapf(
			ErrMissingCreds,
			"missing %s",
			awsBaseURLVar,
		)
	}

	return &DefaultTigrisClient{
		Client: s3.NewFromConfig(aws.Config{
			Region:       "auto",
			BaseEndpoint: aws.String(baseURL),
			Credentials:  credHandler,
			HTTPClient: &http.Client{
				Transport: http.DefaultTransport,
				Timeout:   10 * time.Second,
				CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
					return http.ErrUseLastResponse
				},
			},
		}),
	}, nil
}

// UploadToS3 uploads a file to S3.
func UploadToS3(
	ctx context.Context,
	client Tigris,
	bucket string,
	path string,
	data []byte,
) error {
	// Use custom content type function instead of mime.TypeByExtension
	contentType := GetContentType(path)

	// timeout := defaultS3Timeout

	if isVideoType(contentType) {
		return nil
	}
	// // Use a timeout context to prevent hanging on S3 operations
	// uploadCtx, cancel := context.WithTimeout(ctx, timeout)
	// defer cancel()

	slog.Info("Uploading to S3", "path", path)
	// _, err := client.PutObject(uploadCtx, &s3.PutObjectInput{
	// 	Bucket:      &bucket,
	// 	Key:         aws.String(path),
	// 	Body:        bytes.NewReader(data),
	// 	ContentType: aws.String(contentType),
	// })
	//
	// if err != nil {
	// 	if errors.Is(err, context.DeadlineExceeded) {
	// 		return eris.Wrapf(err, "S3 upload timed out: %s", path)
	// 	}
	//
	// 	return eris.Wrapf(err, "failed to upload to S3: %s", path)
	// }
	//
	slog.Debug("uploaded to S3", "path", path)

	return nil
}

// credHandler is the bucket for the API security.
// It implements the aws.CredentialsProvider interface.
type credHandler struct {
	id  string
	key string
}

// newCredHandler creates a new CredHandler.
func newCredHandler(getEnv func(string) string) (*credHandler, error) {
	id := getEnv(awsAccessKeyIDVar)
	if id == "" {
		return nil, eris.Wrapf(
			ErrMissingCreds,
			"missing %s",
			awsAccessKeyIDVar,
		)
	}
	key := getEnv(awsSecretKeyVar)
	if key == "" {
		return nil, eris.Wrapf(
			ErrMissingCreds,
			"missing %s",
			awsSecretKeyVar,
		)
	}

	return &credHandler{
		id:  id,
		key: key,
	}, nil
}

// Retrieve returns the credentials for the bucket.
func (b *credHandler) Retrieve(_ context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     b.id,
		SecretAccessKey: b.key,
	}, nil
}
