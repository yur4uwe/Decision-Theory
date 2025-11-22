package simplex

import (
	"fmt"
	"lab/matrix"
	"strings"
)

type Coefficients = []float64

func RunSimplex(c Coefficients, A matrix.Matrix, b Coefficients, basis []int, show_tables bool) (Coefficients, matrix.Matrix, Coefficients, []int, int, error) {
	// detect infeasible initial RHS (need Phase I) and return a clear error
	for i := range b {
		if b[i] < 0 {
			return nil, nil, nil, nil, 0, fmt.Errorf("infeasible initial RHS: b[%d]=%v (needs Phase I / artificials)", i, b[i])
		}
	}

	max_iterations := 1000
	if basis == nil {
		c, A, b, basis = InitializeBasis(c, A, b)
	}
	is_found := false

	// simplex tableau initialization
	for iters := range max_iterations {
		c, A, b, basis, is_found = RunSimplexIteration(c, A, b, basis, show_tables)
		if is_found {
			res := make([]float64, len(c))
			for i := 0; i < len(basis); i++ {
				res[basis[i]] = b[i]
			}
			return res, A, b, basis, iters + 1, nil
		}
	}

	return c, A, b, basis, max_iterations, fmt.Errorf("max iterations exceeded")
}

func RunSimplexIteration(c Coefficients, A matrix.Matrix, b Coefficients, basis []int, show_tables bool) (Coefficients, matrix.Matrix, Coefficients, []int, bool) {
	if basis == nil {
		c, A, b, basis = InitializeBasis(c, A, b)
	}

	net_eval := CalculateNetEvaluationRow(c, A, b, basis)

	if show_tables {
		printTableau(c, A, b, net_eval, basis, 0)
	}

	pivot_col_idx := getMaxNetEvaluationIndex(net_eval)

	if net_eval[pivot_col_idx] <= 0 {
		res := make([]float64, len(c))
		for i := 0; i < len(basis); i++ {
			res[basis[i]] = b[i]
		}
		return res, A, b, basis, true // Optimal solution found
	}

	// get pivot column
	col, err := A.Col(pivot_col_idx)
	if err != nil {
		panic(err) // Handle column retrieval error
	}

	// get pivot row index
	pivot_row_idx := findPivotRow(col, b)

	// update basis
	basis[pivot_row_idx] = pivot_col_idx

	// update the tableau
	A, b = updateTableau(A, b, pivot_row_idx, pivot_col_idx)

	return c, A, b, basis, false // Continue iteration
}

func InitializeBasis(c Coefficients, A matrix.Matrix, b Coefficients) (Coefficients, matrix.Matrix, Coefficients, []int) {
	new_A := matrix.NewFromShape(len(A), len(A[0])+len(b), 0)
	for i := 0; i < len(A); i++ {
		for j := 0; j < len(A[0]); j++ {
			new_A[i][j] = A[i][j]
		}
		new_A[i][len(A[0])+i] = 1 // slack variable
	}
	new_c := make(Coefficients, len(c)+len(b))
	new_b := make(Coefficients, len(b))

	copy(new_c, c)
	copy(new_b, b)

	for i := range b {
		new_c[len(c)+i] = 0
	}

	basis := make([]int, len(b))
	for i := range basis {
		basis[i] = len(c) + i
	}

	return new_c, new_A, new_b, basis
}

func CalculateNetEvaluationRow(c []float64, A matrix.Matrix, b Coefficients, basis []int) []float64 {
	num_of_vars := len(c)

	z := make([]float64, num_of_vars+1)

	// calculate profit row
	for i := 0; i < num_of_vars; i++ {
		z[i] = 0
		for j := 0; j < len(basis); j++ {
			z[i] += c[basis[j]] * A[j][i]
		}
	}

	for i := 0; i < len(b); i++ {
		z[num_of_vars] += c[basis[i]] * b[i]
	}

	// subtract cost row
	for i := 0; i < num_of_vars; i++ {
		z[i] = c[i] - z[i]
	}

	return z
}

func getMaxNetEvaluationIndex(z []float64) int {
	max_index := 0
	max := 0.0
	// Exclude the last element in the row that is the profit
	for i := 0; i < len(z)-1; i++ {
		if z[i] > max {
			max = z[i]
			max_index = i
		}
	}
	return max_index
}

func findPivotRow(max_net_col []float64, b Coefficients) int {
	min := 1e9 // Use a large value for initialization
	row_index := -1

	for i := 0; i < len(b); i++ {
		if max_net_col[i] <= 0 {
			continue
		}

		ratio := b[i] / max_net_col[i]
		if ratio < min {
			min = ratio
			row_index = i
		}
	}

	if row_index == -1 {
		panic("No valid pivot row found")
	}

	return row_index
}

func updateTableau(A matrix.Matrix, b []float64, pivot_row_idx, pivot_col_idx int) (matrix.Matrix, []float64) {
	// Update pivot row
	pivot_row, err := A.Row(pivot_row_idx)
	if err != nil {
		panic(err) // Handle row retrieval error
	}

	pivot_element, err := A.Element(pivot_row_idx, pivot_col_idx)
	if err != nil {
		panic(err) // Handle element retrieval error
	}

	for i := 0; i < len(pivot_row); i++ {
		pivot_row[i] /= pivot_element
	}

	b[pivot_row_idx] /= pivot_element

	A.SetRow(pivot_row_idx, pivot_row)

	// update the tableau
	for i := 0; i < A.Rows(); i++ {
		if i == pivot_row_idx {
			continue
		}

		row, err := A.Row(i)
		if err != nil {
			panic(err) // Handle row retrieval error
		}

		row_pivot_element, err := A.Element(i, pivot_col_idx)
		if err != nil {
			panic(err) // Handle element retrieval error
		}

		for j := 0; j < len(row); j++ {
			row[j] -= row_pivot_element * pivot_row[j]
		}

		b[i] -= b[pivot_row_idx] * row_pivot_element

		A.SetRow(i, row)
	}

	return A, b
}

// iter 1 | x1 x2 x3 x4 s1 s2 s3 s4 |
//  basis | 8  10 12 18 0  0  0  0  | b
// ---------------------------------------
// x5 | 0 | 2  4  6  8  1  0  0  0  | 1260
// x6 | 0 | 2  2  0  6  0  1  0  0  | 900
// x7 | 0 | 0  1  1  2  0  0  1  0  | 530
// x8 | 0 | 1  0  1  6  0  0  0  1  | 210
// ---------------------------------------
// z      | 0  0  0  0  0  0  0  0  | 0
// c - z  | 8  10 12 18 0  0  0  0  | 0

func printTableau(c Coefficients, A matrix.Matrix, b []float64, netEval []float64, basis []int, iteration int) {
	numVars := len(c)
	nSlack := A.Rows()
	nOriginal := numVars - nSlack

	// Header line
	fmt.Printf("\n     iter %d  |", iteration)
	for i := 0; i < numVars; i++ {
		if i < nOriginal {
			fmt.Printf("    x%d    ", i+1)
		} else {
			fmt.Printf("    s%d    ", i-nOriginal+1)
		}
	}
	fmt.Printf("|   b\n")

	// Basis coefficients
	fmt.Printf("      basis  |")
	for i := 0; i < numVars; i++ {
		fmt.Printf(" %7.3f ", c[i])
	}
	fmt.Printf("|%9s\n", "b")

	// Separator
	sep := "----+--------" + strings.Repeat("+---------", numVars) + "+----------"
	fmt.Println(sep)

	// Basis rows
	for i := 0; i < A.Rows(); i++ {
		varName := "x"
		if basis[i] >= nOriginal {
			varName = fmt.Sprintf("s%d", basis[i]-nOriginal+1)
		} else {
			varName = fmt.Sprintf("x%d", basis[i]+1)
		}

		fmt.Printf("%-4s|%7.3f |", varName, c[basis[i]])
		row, _ := A.Row(i)
		for _, val := range row {
			fmt.Printf(" %7.3f ", val)
		}
		fmt.Printf("|%9.3f\n", b[i])
	}

	// Bottom separator
	fmt.Println(sep)

	// z row
	fmt.Printf("           z |")
	for i := 0; i < numVars; i++ {
		fmt.Printf(" %7.3f ", c[i]-netEval[i])
	}
	fmt.Printf("|%9.3f\n", netEval[numVars])

	// c - z row
	fmt.Printf("       c - z |")
	for i := 0; i < numVars; i++ {
		fmt.Printf(" %7.3f ", netEval[i])
	}
	fmt.Printf("|\n\n")
}
