package credited

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rotisserie/eris"
)

const (
	awsAccessKeyIDVar = "AWS_ACCESS_KEY_ID"
	awsSecretKeyVar   = "AWS_SECRET_ACCESS_KEY"
	awsBaseURLVar     = "AWS_ENDPOINT_URL_S3"
)

// AWSClient is a wrapper for the AWS S3 client.
type AWSClient struct{ *s3.Client }

// NewAWSClient creates a new AWSClient.
func NewAWSClient(getenv func(string) string) (*AWSClient, error) {
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
	return &AWSClient{
		Client: s3.NewFromConfig(aws.Config{
			Region:       "auto",
			BaseEndpoint: aws.String(baseURL),
			Credentials:  credHandler,
		}),
	}, nil
}

// credHandler is the bucket for the api security.
// It implements the aws.Authenitcator
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
func (b *credHandler) Retrieve(
	_ context.Context,
) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     b.id,
		SecretAccessKey: b.key,
	}, nil
}
