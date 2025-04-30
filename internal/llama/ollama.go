package llama

import (
	"context"
	"math"
	"net/http"
	"net/url"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/ollama/ollama/api"
	"github.com/rotisserie/eris"
	"gonum.org/v1/gonum/mat"
)

const (
	ollamaURLVar   = "OLLAMA_URL"
	embeddingModel = "nomic-embed-text"
	embeddingDim   = assets.EmbedLength
)

// Client is the minimal interface for the Ollama client.
type Client interface {
	Embeddings(ctx context.Context, req *api.EmbeddingRequest) (*api.EmbeddingResponse, error)
}

// OllamaClient is a wrapper for the Ollama client.
type OllamaClient struct {
	client Client
}

// NewOllamaClient creates a new OllamaClient.
func NewOllamaClient(getenv func(string) string) (*OllamaClient, error) {
	urlStr := getenv(ollamaURLVar)
	if urlStr == "" {
		return nil, eris.Wrapf(ErrMissingCreds, "missing %s", ollamaURLVar)
	}
	urlVal, err := url.Parse(urlStr)
	if err != nil {
		return nil, eris.Wrapf(ErrInvalidCreds, "invalid %s", ollamaURLVar)
	}
	return &OllamaClient{
		client: api.NewClient(urlVal, http.DefaultClient),
	}, nil
}

// Embeddings returns the embeddings for the provided text.
func (c *OllamaClient) Embeddings(
	ctx context.Context,
	content string,
	emb *assets.Doc,
) (err error) {
	resp, err := c.client.Embeddings(ctx, &api.EmbeddingRequest{
		Model:  embeddingModel,
		Prompt: content,
	})
	if err != nil {
		return eris.Wrapf(
			err,
			"embedding generation failed for %s",
			embeddingModel,
		)
	}

	emb.X, emb.Y, emb.Z = projectTo3D(
		resp.Embedding,
		projectionMatrixCreate(assets.EmbedLength, 3),
	)
	// copy(emb.Vec[:], resp.Embedding[:assets.EmbedLength])

	return
}

// Generate a deterministic projection matrix for dimensionality reduction
func projectionMatrixCreate(inputDim, outputDim int) *mat.Dense {
	// Create a matrix with dimensions [outputDim x inputDim]
	// This will multiply directly with the input vector without transposition
	matrix := mat.NewDense(outputDim, inputDim, nil)

	for i := range outputDim { // For each output dimension
		// center point for this dimension's projection focus
		center := (i * inputDim) / outputDim

		for j := range inputDim {
			// Calculate a weight based on distance from the center point
			// Use a Gaussian-like function for smooth projection
			dist := float64(j - center)
			weight := math.Exp(-0.5 * (dist * dist) / float64(inputDim/10))

			// Set the weight in the projection matrix
			matrix.Set(i, j, weight)
		}

		// Normalize the row to unit length
		var rowSum float64
		for j := range inputDim {
			rowSum += matrix.At(i, j) * matrix.At(i, j)
		}
		rowSum = math.Sqrt(rowSum)
		if rowSum > 0 {
			for j := range inputDim {
				matrix.Set(i, j, matrix.At(i, j)/rowSum)
			}
		}
	}

	return matrix
}

// Project a single embedding to 3D using a projection matrix
func projectTo3D(embedding []float64, projectionMatrix *mat.Dense) (x, y, z float64) {
	// Create a vector from the embedding
	embedVec := mat.NewVecDense(len(embedding), embedding)

	// Project to 3D - no need to transpose
	result := mat.NewVecDense(3, nil)
	result.MulVec(projectionMatrix, embedVec)

	// Extract x, y, z coordinates
	x = result.AtVec(0)
	y = result.AtVec(1)
	z = result.AtVec(2)

	return
}
