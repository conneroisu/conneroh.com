// Package main updates the database with new vault content.
package main

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/conneroisu/conneroh.com/internal/data/docs"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/genstruct"
	ollama "github.com/prathyushnallamothu/ollamago"
	"gonum.org/v1/gonum/mat"
)

const (
	embeddingSize = 768
)

var client = ollama.NewClient(
	ollama.WithTimeout(time.Minute * 5),
)

// Project a single embedding to 3D using a projection matrix
func projectTo3D(embedding []float64, projectionMatrix *mat.Dense) (x, y, z float64) {
	// Create a vector from the embedding
	embedVec := mat.NewVecDense(len(embedding), embedding)

	// Project to 3D
	result := mat.NewVecDense(3, nil)
	result.MulVec(projectionMatrix.T(), embedVec)

	// Extract x, y, z coordinates
	x = result.AtVec(0)
	y = result.AtVec(1)
	z = result.AtVec(2)

	return
}

// Generate a random projection matrix for demonstration
// In practice, this would be calculated using PCA or another technique
func generateProjectionMatrix(inputDim, outputDim int) *mat.Dense {
	data := make([]float64, inputDim*outputDim)
	for i := range data {
		data[i] = rand.Float64()*2 - 1 // Random values between -1 and 1
	}
	return mat.NewDense(inputDim, outputDim, data)
}

func main() {
	if err := Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Run parses all markdown files in the database.
func Run() error {
	parsedTags, err := pathParse[gen.Tag]("tags", docs.Tags)
	if err != nil {
		return err
	}

	parsedPosts, err := pathParse[gen.Post]("posts", docs.Posts)
	if err != nil {
		return err
	}

	parsedProjects, err := pathParse[gen.Project]("projects", docs.Projects)
	if err != nil {
		return fmt.Errorf("failed to parse projects: %v", err)
	}
	// Create a generator for people
	postGen, err := genstruct.NewGenerator(genstruct.Config{
		PackageName: "gen",
		OutputFile:  fmt.Sprintf("internal/data/gen/generated_%s.go", "data"),
		Logger:      DefaultLogger,
	}, parsedPosts, parsedTags, parsedProjects)
	if err != nil {
		return err
	}

	return postGen.Generate()
}

// DefaultLogger is a default logger.
var DefaultLogger = slog.New(
	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			if a.Key == "level" {
				return slog.Attr{}
			}
			if a.Key == slog.SourceKey {
				str := a.Value.String()
				split := strings.Split(str, "/")
				if len(split) > 2 {
					a.Value = slog.StringValue(
						strings.Join(split[len(split)-2:], "/"),
					)
					a.Value = slog.StringValue(
						strings.ReplaceAll(a.Value.String(), "}", ""),
					)
				}
			}
			return a
		}}),
)
