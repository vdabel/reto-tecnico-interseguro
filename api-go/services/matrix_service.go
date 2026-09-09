package services

import (
	"errors"
	"math"

	"gonum.org/v1/gonum/mat"
)

var (
	ErrEmptyMatrix    = errors.New("la matriz no puede estar vacía")
	ErrJaggedMatrix   = errors.New("todas las filas de la matriz deben tener la misma longitud")
	ErrQRInvalidShape = errors.New("para la factorización QR, el número de filas debe ser mayor o igual al número de columnas (m >= n)")
)

// MatrixService encapsula la lógica matemática y de transformación de matrices.
type MatrixService struct{}

// NewMatrixService retorna una nueva instancia de MatrixService.
func NewMatrixService() *MatrixService {
	return &MatrixService{}
}

// ValidateMatrix valida que la matriz sea rectangular y no vacía.
func (s *MatrixService) ValidateMatrix(matrix [][]float64) error {
	if len(matrix) == 0 {
		return ErrEmptyMatrix
	}

	cols := len(matrix[0])
	if cols == 0 {
		return ErrEmptyMatrix
	}

	for _, row := range matrix {
		if len(row) != cols {
			return ErrJaggedMatrix
		}
	}

	return nil
}

// FactorizeQR realiza la descomposición QR (A = Q * R) utilizando gonum/mat.
// Q es una matriz ortogonal (m x m) y R es una matriz triangular superior (m x n).
func (s *MatrixService) FactorizeQR(matrix [][]float64) ([][]float64, [][]float64, error) {
	if err := s.ValidateMatrix(matrix); err != nil {
		return nil, nil, err
	}

	rows := len(matrix)
	cols := len(matrix[0])

	if rows < cols {
		return nil, nil, ErrQRInvalidShape
	}

	// Aplanar datos para gonum Dense
	flatData := make([]float64, rows*cols)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			flatData[i*cols+j] = matrix[i][j]
		}
	}

	denseA := mat.NewDense(rows, cols, flatData)

	var qr mat.QR
	qr.Factorize(denseA)

	var qDense mat.Dense
	qr.QTo(&qDense)

	var rDense mat.Dense
	qr.RTo(&rDense)

	qSlice := s.denseToCleanSlice(&qDense)
	rSlice := s.denseToCleanSlice(&rDense)

	return qSlice, rSlice, nil
}

// Rotate90 realiza una rotación de 90 grados en sentido horario a la matriz rectangular.
// Una matriz de m x n pasa a ser de n x m.
func (s *MatrixService) Rotate90(matrix [][]float64) ([][]float64, error) {
	if err := s.ValidateMatrix(matrix); err != nil {
		return nil, err
	}

	rows := len(matrix)
	cols := len(matrix[0])

	rotated := make([][]float64, cols)
	for j := 0; j < cols; j++ {
		rotated[j] = make([]float64, rows)
		for i := 0; i < rows; i++ {
			rotated[j][rows-1-i] = matrix[i][j]
		}
	}

	return rotated, nil
}

// denseToCleanSlice convierte una matriz de Gonum a [][]float64 limpiando residuos de precisión infinitesimal.
func (s *MatrixService) denseToCleanSlice(m mat.Matrix) [][]float64 {
	rows, cols := m.Dims()
	result := make([][]float64, rows)

	const epsilon = 1e-12

	for i := 0; i < rows; i++ {
		result[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			val := m.At(i, j)
			if math.Abs(val) < epsilon {
				val = 0.0
			}
			result[i][j] = math.Round(val*1e6) / 1e6 // Redondeo a 6 decimales para precisión limpia
		}
	}

	return result
}
