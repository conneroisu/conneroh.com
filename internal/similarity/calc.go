package similarity

import (
	"sync"

	"gonum.org/v1/gonum/mat"
)

// Calculator is a struct that can be used to calculate similarity between two strings.
type Calculator struct {
	biFuncs []biFuncCoefficient
	workers int
}

// Option is a function that can be used to configure a Calculator.
type Option func(c *Calculator)

type handler func(queryVec *mat.VecDense, indexVec *mat.VecDense) (float64, error)

type biFuncCoefficient struct {
	handler     handler
	coefficient float64
}

// NewCalculator creates a new Calculator with the given options.
func NewCalculator(opts ...Option) *Calculator {
	c := &Calculator{}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Similarity calculates the similarity between two strings.
func (c *Calculator) Similarity(qVec, iVec []float64) (sim float64, err error) {
	qVecMat := mat.NewVecDense(len(qVec), qVec)
	iVecMat := mat.NewVecDense(len(iVec), iVec)

	// If workers is not set or is set to 1, use sequential processing
	if c.workers <= 1 {
		var val float64
		for _, biFunc := range c.biFuncs {
			val, err = biFunc.handler(qVecMat, iVecMat)
			if err != nil {
				return 0, err
			}
			sim += biFunc.coefficient * val
		}

		return sim, nil
	}

	// Use parallel processing with workers
	results := make(chan float64, len(c.biFuncs))
	errCh := make(chan error, len(c.biFuncs))

	// Create a work queue
	jobs := make(chan int, len(c.biFuncs))
	for i := range c.biFuncs {
		jobs <- i
	}
	close(jobs)

	// Create worker pool
	var wg sync.WaitGroup
	numWorkers := min(c.workers, len(c.biFuncs))
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for jobIndex := range jobs {
				biFunc := c.biFuncs[jobIndex]
				val, err := biFunc.handler(qVecMat, iVecMat)
				if err != nil {
					errCh <- err

					return
				}
				results <- biFunc.coefficient * val
			}
		}()
	}

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
		close(errCh)
	}()

	// Check for errors
	select {
	case err := <-errCh:
		if err != nil {
			return 0, err
		}
	default:
		// No errors
	}

	// Sum up the results
	for res := range results {
		sim += res
	}

	return sim, nil
}

// WithWorkers sets the number of workers to use for the calculation.
func WithWorkers(workers int) Option {
	return func(c *Calculator) {
		c.workers = workers
	}
}

// WithSimilarityDotMatrix sets the similarity function to use with a
// coefficient.
//
// $$a \cdot b=\sum_{i=1}^{n} a_{i} b_{i}$$
//
// It adds the similarity dot matrix to the comparison functions with the given
// coefficient.
func WithSimilarityDotMatrix(coefficient float64) Option {
	return func(c *Calculator) {
		c.biFuncs = append(c.biFuncs, biFuncCoefficient{
			handler:     similarityDotMatrix,
			coefficient: coefficient,
		})
	}
}

// WithEuclideanDistance sets the EuclideanDistance function with a coefficient.
//
// $$d(x, y) = \sqrt{\sum_{i=1}^{n}(x_i - y_i)^2}$$.
func WithEuclideanDistance(coefficient float64) Option {
	return func(c *Calculator) {
		c.biFuncs = append(c.biFuncs, biFuncCoefficient{
			handler:     euclideanDistance,
			coefficient: coefficient,
		})
	}
}

// WithManhattanDistance sets the ManhattanDistance function with a coefficient.
//
// $$d(x, y) = |x_1 - y_1| + |x_2 - y_2| + ... + |x_n - y_n|$$
//
// It adds the manhatten distance to the comparison functions with the given
// coefficient.
func WithManhattanDistance(coefficient float64) Option {
	return func(c *Calculator) {
		c.biFuncs = append(c.biFuncs, biFuncCoefficient{
			handler:     manhattanDistance,
			coefficient: coefficient,
		})
	}
}

// WithJaccardSimilarity sets the JaccardSimilarity function with a coefficient.
//
// $$J(A, B)=\frac{|A \cap B|}{|A \cup B|}$$
//
// It adds the jaccard similarity to the comparison functions with the given
// coefficient.
func WithJaccardSimilarity(coefficient float64) Option {
	return func(c *Calculator) {
		c.biFuncs = append(c.biFuncs, biFuncCoefficient{
			handler:     jaccardSimilarity,
			coefficient: coefficient,
		})
	}
}

// WithPearsonCorrelation sets the PearsonCorrelation function with a
// coefficient.
//
// $$r=\frac{\sum\left(x_{i}-\bar{x}\right)\left(y_{i}-\bar{y}\right)}{\sqrt{\sum\left(x_{i}-\bar{x}\right)^{2} \sum\left(y_{i}-\bar{y}\right)^{2}}}$$
//
// It adds the pearson correlation to the comparison functions with the given
// coefficient.
func WithPearsonCorrelation(coefficient float64) Option {
	return func(c *Calculator) {
		c.biFuncs = append(c.biFuncs, biFuncCoefficient{
			handler:     pearsonCorrelation,
			coefficient: coefficient,
		})
	}
}

// WithHammingDistance sets the HammingDistance function with a coefficient.
func WithHammingDistance(coefficient float64) Option {
	return func(c *Calculator) {
		c.biFuncs = append(c.biFuncs, biFuncCoefficient{
			handler:     hammingDistance,
			coefficient: coefficient,
		})
	}
}
