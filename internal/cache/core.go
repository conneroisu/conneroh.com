package cache

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/assets"
)

type (
	// ollamaClient is the ollama client.
	ollamaClient interface {
		Embeddings(
			ctx context.Context,
			content string,
			emb *assets.Doc,
		) (err error)
	}

	// tigrisClient is the tigris client.
	tigrisClient interface {
		PutObject(
			ctx context.Context,
			input *s3.PutObjectInput,
			opts ...func(*s3.Options),
		) (*s3.PutObjectOutput, error)
	}
	// waitGroup is a wait group.
	waitGroup interface {
		Add(delta int)
		Done()
		Wait()
	}
)
