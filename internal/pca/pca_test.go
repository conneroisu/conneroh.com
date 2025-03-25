package pca

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestPCAWithZeroComponents(t *testing.T) {
	t.Run("Set to pca 0 components with first matrix", func(t *testing.T) {
		X1 := mat.NewDense(3, 7, []float64{6, 5, 4, 3, 8, 2, 9, 5, 1, 10, 2, 3, 8, 7, 5, 14, 2, 3, 6, 3, 2})
		pca := NewPCA(0)
		rows, cols := pca.FitTransform(X1).Dims()
		assert.Equal(t, rows, 3)
		assert.Equal(t, cols, 3)
	})

	t.Run("Set to pca 0 components with second matrix", func(t *testing.T) {
		X1 := mat.NewDense(10, 5, []float64{
			0.52984892, 0.1141001, 0.91599294, 0.9574267, 0.15361222,
			0.07057588, 0.46371013, 0.73091854, 0.84641034, 0.08122213,
			0.96221946, 0.60367214, 0.69851546, 0.91965564, 0.27040597,
			0.03152856, 0.97912403, 0.39487038, 0.12232594, 0.18474705,
			0.77061953, 0.35898551, 0.78684562, 0.11638404, 0.88908044,
			0.35828086, 0.47214831, 0.95781755, 0.74762736, 0.59850757,
			0.07806127, 0.96940955, 0.15751804, 0.00973325, 0.85041635,
			0.02663938, 0.49755131, 0.57984119, 0.12233871, 0.47967853,
			0.63903222, 0.88556565, 0.79797963, 0.13345186, 0.37415535,
			0.60605207, 0.52067165, 0.91217494, 0.57148943, 0.92210331})
		pca := NewPCA(0)
		rows, cols := pca.FitTransform(X1).Dims()
		assert.Equal(t, rows, 10)
		assert.Equal(t, cols, 5)
	})
}

func TestPCAWithNComponents(t *testing.T) {
	t.Run("Set to pca 3 components with 5x5 matrix", func(t *testing.T) {
		X := mat.NewDense(5, 5, []float64{
			0.23030838, 0.05669317, 0.3187813, 0.34455114, 0.98062806,
			0.38995469, 0.2996771, 0.99043575, 0.04443827, 0.99527955,
			0.27266308, 0.14068906, 0.46999473, 0.03296131, 0.90855405,
			0.28360708, 0.8839966, 0.81107014, 0.52673877, 0.59432817,
			0.64107253, 0.56165215, 0.79811756, 0.48845398, 0.20506649})
		pca := NewPCA(3)
		rows, cols := pca.FitTransform(X).Dims()
		assert.Equal(t, rows, 5)
		assert.Equal(t, cols, 3)
	})

	t.Run("Set to pca 2 components with 3x5 matrix", func(t *testing.T) {
		X := mat.NewDense(3, 5, []float64{
			0.12294845, 0.55170713, 0.67572832, 0.60615516, 0.38184551,
			0.93486821, 0.15120374, 0.89760169, 0.74715672, 0.81373931,
			0.42821569, 0.47457753, 0.18960954, 0.42466159, 0.34166049})
		pca := NewPCA(2)
		pca.Fit(X)
		rows, cols := pca.Transform(X).Dims()
		assert.Equal(t, rows, 3)
		assert.Equal(t, cols, 2)
	})
}

func TestPCAFitAndTransformSeparately(t *testing.T) {
	t.Run("Set to pca 3 components with 5x5 matrix", func(t *testing.T) {
		X := mat.NewDense(5, 5, []float64{
			0.23030838, 0.05669317, 0.3187813, 0.34455114, 0.98062806,
			0.38995469, 0.2996771, 0.99043575, 0.04443827, 0.99527955,
			0.27266308, 0.14068906, 0.46999473, 0.03296131, 0.90855405,
			0.28360708, 0.8839966, 0.81107014, 0.52673877, 0.59432817,
			0.64107253, 0.56165215, 0.79811756, 0.48845398, 0.20506649})
		pca := NewPCA(3)
		pca.Fit(X)

		rows, cols := pca.Transform(X).Dims()
		assert.Equal(t, rows, 5)
		assert.Equal(t, cols, 3)
	})
}

func TestPCAWithNilSVD(t *testing.T) {
	t.Run("Transform before fit PCA model", func(t *testing.T) {
		X := mat.NewDense(3, 2, []float64{5, 7, 2, 3, 1, -9})
		pca := NewPCA(2)

		assert.Panics(t, func() { pca.Transform(X) })
	})
}

func TestPCAWithLessThanZeroComponents(t *testing.T) {
	t.Run("Set to pca -2 components", func(t *testing.T) {
		X := mat.NewDense(1, 2, []float64{3, 0})
		pca := NewPCA(-2)

		assert.Panics(t, func() { pca.Fit(X) })
	})
}

func TestMatrixAndVectorMismatchDim(t *testing.T) {
	t.Run("Mismatch dimensions between matrix and vector", func(t *testing.T) {
		X := mat.NewDense(3, 1, []float64{-3, 0, 1})
		v := mat.NewDense(2, 3, []float64{9, 4, 8, 7, 8, 7})

		assert.Panics(t, func() { matrixSubVector(X, v) })
	})
}

func TestPCAComponentBiggerThanFeature(t *testing.T) {
	t.Run("Set to pca 5 components with 3x3 matrix", func(t *testing.T) {
		X := mat.NewDense(3, 3, []float64{-3, 0, 1, 55, 2, 9, -9, 3, 66})
		pca := NewPCA(5)
		rows, cols := pca.FitTransform(X).Dims()
		assert.Equal(t, rows, 3)
		assert.Equal(t, cols, 3)
	})
}
