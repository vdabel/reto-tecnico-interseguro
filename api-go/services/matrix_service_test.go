package services

import (
	"math"
	"testing"
)

func TestValidateMatrix(t *testing.T) {
	svc := NewMatrixService()

	t.Run("Valid square matrix", func(t *testing.T) {
		m := [][]float64{
			{1, 2},
			{3, 4},
		}
		if err := svc.ValidateMatrix(m); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Valid rectangular matrix", func(t *testing.T) {
		m := [][]float64{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
			{10, 11, 12},
		}
		if err := svc.ValidateMatrix(m); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Empty matrix", func(t *testing.T) {
		m := [][]float64{}
		if err := svc.ValidateMatrix(m); err != ErrEmptyMatrix {
			t.Errorf("expected ErrEmptyMatrix, got %v", err)
		}
	})

	t.Run("Jagged matrix", func(t *testing.T) {
		m := [][]float64{
			{1, 2, 3},
			{4, 5},
		}
		if err := svc.ValidateMatrix(m); err != ErrJaggedMatrix {
			t.Errorf("expected ErrJaggedMatrix, got %v", err)
		}
	})
}

func TestFactorizeQR(t *testing.T) {
	svc := NewMatrixService()

	t.Run("3x3 Matrix QR Decomposition", func(t *testing.T) {
		input := [][]float64{
			{12, -51, 4},
			{6, 167, -68},
			{-4, 24, -41},
		}

		q, r, err := svc.FactorizeQR(input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Validar dimensiones: Q debe ser 3x3 y R debe ser 3x3
		if len(q) != 3 || len(q[0]) != 3 {
			t.Errorf("unexpected Q dimensions: %dx%d", len(q), len(q[0]))
		}
		if len(r) != 3 || len(r[0]) != 3 {
			t.Errorf("unexpected R dimensions: %dx%d", len(r), len(r[0]))
		}

		// R debe ser triangular superior: r[i][j] == 0 para i > j
		for i := 0; i < len(r); i++ {
			for j := 0; j < i; j++ {
				if math.Abs(r[i][j]) > 1e-4 {
					t.Errorf("R is not upper triangular: r[%d][%d] = %f", i, j, r[i][j])
				}
			}
		}

		// Verificar que Q * R == A (aproximadamente)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				sum := 0.0
				for k := 0; k < 3; k++ {
					sum += q[i][k] * r[k][j]
				}
				if math.Abs(sum-input[i][j]) > 1e-4 {
					t.Errorf("Q*R[%d][%d] = %f, expected %f", i, j, sum, input[i][j])
				}
			}
		}
	})

	t.Run("Rectangular Matrix 4x3 QR Decomposition", func(t *testing.T) {
		input := [][]float64{
			{1, -1, 4},
			{1, 4, -2},
			{1, 4, 2},
			{1, -1, 0},
		}

		q, r, err := svc.FactorizeQR(input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Q es 4x4, R es 4x3
		if len(q) != 4 || len(q[0]) != 4 {
			t.Errorf("unexpected Q dimensions: %dx%d", len(q), len(q[0]))
		}
		if len(r) != 4 || len(r[0]) != 3 {
			t.Errorf("unexpected R dimensions: %dx%d", len(r), len(r[0]))
		}
	})

	t.Run("Matrix with rows < cols fails QR", func(t *testing.T) {
		input := [][]float64{
			{1, 2, 3},
			{4, 5, 6},
		}

		_, _, err := svc.FactorizeQR(input)
		if err != ErrQRInvalidShape {
			t.Errorf("expected ErrQRInvalidShape, got %v", err)
		}
	})
}

func TestRotate90(t *testing.T) {
	svc := NewMatrixService()

	input := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}
	// Matriz 2x3 rotada 90° en sentido horario es 3x2:
	// [4, 1]
	// [5, 2]
	// [6, 3]
	expected := [][]float64{
		{4, 1},
		{5, 2},
		{6, 3},
	}

	rotated, err := svc.Rotate90(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(rotated) != 3 || len(rotated[0]) != 2 {
		t.Fatalf("unexpected dimensions: %dx%d", len(rotated), len(rotated[0]))
	}

	for i := range expected {
		for j := range expected[i] {
			if rotated[i][j] != expected[i][j] {
				t.Errorf("at [%d][%d]: expected %f, got %f", i, j, expected[i][j], rotated[i][j])
			}
		}
	}
}
