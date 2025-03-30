package markdown

import (
	"context"
	"math"
	"time"

	"github.com/conneroisu/conneroh.com/internal/data/gen"
	ollama "github.com/prathyushnallamothu/ollamago"
	"gonum.org/v1/gonum/mat"
)

var client = ollama.NewClient(ollama.WithTimeout(time.Minute * 5))

// TextEmbeddingCreate creates an embedding for the given text.
func TextEmbeddingCreate(
	ctx context.Context,
	client *ollama.Client,
	input string,
) ([gen.EmbedLength]float64, float64, float64, float64, error) {
	resp, err := client.Embeddings(ctx, ollama.EmbeddingsRequest{
		Model:  "nomic-embed-text",
		Prompt: input,
	})
	if err != nil {
		return [gen.EmbedLength]float64{}, 0, 0, 0, err
	}
	proj := projectionMatrixCreate(gen.EmbedLength, 3)
	x, y, z := projectTo3D(resp.Embedding, proj)
	embs := [gen.EmbedLength]float64{}
	for i := range embs {
		embs[i] = resp.Embedding[i]
	}
	return embs, x, y, z, nil
}

// Generate a deterministic projection matrix for dimensionality reduction
func projectionMatrixCreate(inputDim, outputDim int) *mat.Dense {
	// Create a matrix with dimensions [outputDim x inputDim]
	// This will multiply directly with the input vector without transposition
	matrix := mat.NewDense(outputDim, inputDim, nil)

	// For each output dimension
	for i := range outputDim {
		// Calculate a center point for this dimension's projection focus
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
