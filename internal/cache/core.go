package cache

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/assets"
)

type (
	// OllamaClient is the ollama client.
	OllamaClient interface {
		Embeddings(
			ctx context.Context,
			content string,
			emb *assets.Doc,
		) (err error)
	}

	// TigrisClient is the tigris client.
	TigrisClient interface {
		PutObject(
			ctx context.Context,
			input *s3.PutObjectInput,
			opts ...func(*s3.Options),
		) (*s3.PutObjectOutput, error)
	}
	// WaitGroup is a wait group.
	WaitGroup interface {
		Add(delta int)
		Done()
		Wait()
	}
)
