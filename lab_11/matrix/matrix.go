package matrix

import (
	"fmt"
	"math"
)

type Matrix [][]float64

// Constructor
func NewMatrix(data [][]float64) Matrix {
	m := Matrix(data)
	m.validate()
	return m
}

func NewFromShape(rows, cols int, def_val float64) Matrix {
	if def_val == 0 {
		def_val = 0 // Default value handling
	}
	m := make(Matrix, rows)
	for i := range m {
		m[i] = make([]float64, cols)
		for j := range m[i] {
			m[i][j] = def_val
		}
	}
	return m
}

// Get a row with bounds checking
func (m Matrix) Row(row int) ([]float64, error) {
	if row < 0 || row >= len(m) {
		return nil, fmt.Errorf("row index out of bounds")
	}
	return append([]float64{}, m[row]...), nil
}

// Get a column with bounds checking
func (m Matrix) Col(col int) ([]float64, error) {
	if len(m) == 0 || col < 0 || col >= len(m[0]) {
		return nil, fmt.Errorf("column index out of bounds")
	}
	column := make([]float64, len(m))
	for i := 0; i < len(m); i++ {
		column[i] = m[i][col]
	}
	return column, nil
}

// Set a row with bounds checking
func (m Matrix) SetRow(row int, values []float64) error {
	if row < 0 || row >= len(m) {
		return fmt.Errorf("row index out of bounds")
	}
	if len(values) != len(m[0]) {
		return fmt.Errorf("row length mismatch")
	}
	m[row] = values
	return nil
}

// Set a column with bounds checking
func (m Matrix) SetCol(col int, values []float64) error {
	if len(m) == 0 || col < 0 || col >= len(m[0]) {
		return fmt.Errorf("column index out of bounds")
	}
	if len(values) != len(m) {
		return fmt.Errorf("column length mismatch")
	}
	for i := 0; i < len(m); i++ {
		m[i][col] = values[i]
	}
	return nil
}

// Get an element with bounds checking
func (m Matrix) Element(row, col int) (float64, error) {
	if row < 0 || row >= len(m) || col < 0 || col >= len(m[0]) {
		return 0, fmt.Errorf("index out of bounds")
	}
	return m[row][col], nil
}

// Create a deep copy of the matrix
func (m Matrix) Copy() Matrix {
	copied := make(Matrix, len(m))
	for i := range m {
		copied[i] = append([]float64{}, m[i]...)
	}
	return copied
}

// Validate the matrix structure
func (m Matrix) validate() {
	if len(m) == 0 {
		return
	}
	colLen := len(m[0])
	for i := 1; i < len(m); i++ {
		if len(m[i]) != colLen {
			panic("inconsistent row lengths in matrix")
		}
	}
}

func (m Matrix) Rows() int {
	return len(m)
}

func (m Matrix) Cols() int {
	if len(m) == 0 {
		return 0
	}
	return len(m[0])
}

func (m Matrix) Set(row, col int, value float64) error {
	if row < 0 || row >= len(m) || col < 0 || col >= len(m[0]) {
		return fmt.Errorf("index out of bounds")
	}
	m[row][col] = value
	return nil
}

// Print prints the matrix
func (m Matrix) Print(title string) {
	fmt.Println(title)
	for i := range m {
		fmt.Printf("[")
		for j := range m[i] {
			fmt.Printf("%8.3f ", m[i][j])
		}
		fmt.Printf("]\n")
	}
	fmt.Println()
}

// MakePositive adds a constant to all elements to make the matrix strictly positive
func (m Matrix) MakePositive() (Matrix, float64) {
	minVal := math.Inf(1)
	for i := range m {
		for j := range m[i] {
			if m[i][j] < minVal {
				minVal = m[i][j]
			}
		}
	}

	shift := 0.0
	if minVal <= 0 {
		shift = math.Abs(minVal) + 1
	}

	result := m.Copy()
	if shift > 0 {
		for i := range result {
			for j := range result[i] {
				result[i][j] += shift
			}
		}
	}

	return result, shift
}
