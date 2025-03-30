package markdown

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/transport/http"
)

// CredHandler is the bucket for the api security.
// It implements the aws.Authenitcator
type CredHandler struct {
	Name string `yaml:"name"`
	ID   string
	Key  string
}

// Retrieve returns the credentials for the bucket.
func (b *CredHandler) Retrieve(
	_ context.Context,
) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     b.ID,
		SecretAccessKey: b.Key,
	}, nil
}

// WithHeader adds a header to the s3 client.
func WithHeader(key, value string) func(*s3.Options) {
	return func(options *s3.Options) {
		options.APIOptions = append(
			options.APIOptions,
			http.AddHeaderValue(key, value),
		)
	}
}
