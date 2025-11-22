package main

import (
	"decision-theory/games"
	"decision-theory/lab_11/matrix"
	"fmt"

	"github.com/willauld/lpsimplex"
)

// Vector represents a 1D vector
type Vector []float64

// GameResult contains the solution of a Matrix game
type GameResult struct {
	PlayerXStrategy Vector
	PlayerYStrategy Vector
	GameValue       float64
}

// SolveMatrixGame solves a Matrix game using the generalized algorithm
func SolveMatrixGame(game matrix.Matrix) GameResult {
	fmt.Println("=== Solving Matrix Game ===")
	game.Print("Original Game Matrix:")

	// Step 1: Make Matrix strictly positive
	positiveMatrix, shift := game.MakePositive()
	positiveMatrix.Print("Positive Matrix:")

	rows := len(positiveMatrix)
	cols := len(positiveMatrix[0])

	// Step 2: Solve LP for player 1 (minimization)
	x := solvePlayerOneProblem(positiveMatrix, rows, cols)

	// Step 3: Solve LP for player 2 (maximization)
	y := solvePlayerTwoProblem(positiveMatrix, rows, cols)

	// Step 4: Normalize strategies and calculate game value
	thetaX := sum(x)
	phiY := sum(y)

	xStar := make(Vector, rows)
	for i := range x {
		if thetaX != 0 {
			xStar[i] = x[i] / thetaX
		} else {
			xStar[i] = 0
		}
	}

	yStar := make(Vector, cols)
	for i := range y {
		if phiY != 0 {
			yStar[i] = y[i] / phiY
		} else {
			yStar[i] = 0
		}
	}

	// Game value for positive Matrix
	gameValuePositive := 1.0
	if thetaX != 0 {
		gameValuePositive = 1.0 / thetaX
	}

	// Game value for original Matrix
	gameValue := gameValuePositive - shift

	fmt.Printf("Sum of x (θ): %.6f\n", thetaX)
	fmt.Printf("Sum of y (φ): %.6f\n", phiY)
	fmt.Printf("Game value (positive Matrix): %.6f\n", gameValuePositive)
	fmt.Printf("Game value (original Matrix): %.6f\n\n", gameValue)

	return GameResult{
		PlayerXStrategy: xStar,
		PlayerYStrategy: yStar,
		GameValue:       gameValue,
	}
}

// solvePlayerOneProblem solves the LP for player 1
func solvePlayerOneProblem(A matrix.Matrix, rows, cols int) Vector {
	// Primal (player 1): minimize sum(x_i) subject to A^T * x >= 1, x >= 0
	// Convert constraints into Aub * x <= bub by negation: (-A^T) x <= -1
	n := rows
	m := cols

	c := make([]float64, n)
	for i := 0; i < n; i++ {
		c[i] = 1.0
	}

	// build Aub = -A^T and bub = -1
	Aub := make([][]float64, m)
	bub := make([]float64, m)
	for j := 0; j < m; j++ {
		Aub[j] = make([]float64, n)
		for i := 0; i < n; i++ {
			Aub[j][i] = -A[i][j]
		}
		bub[j] = -1.0
	}

	// Solve (minimize)
	xsol, err := SolveLP(c, Aub, bub, false)
	if err != nil {
		fmt.Printf("Player 1 LP failed: %v\n", err)
		return make(Vector, rows)
	}

	x := make(Vector, rows)
	for i := 0; i < rows; i++ {
		x[i] = xsol[i]
	}
	return x
}

// solvePlayerTwoProblem solves the LP for player 2
func solvePlayerTwoProblem(A matrix.Matrix, rows, cols int) Vector {
	// Dual (player 2): maximize sum(y_j) subject to A * y <= 1, y >= 0
	n := cols
	m := rows

	// objective c (maximize): pass positive c and set maximize=true
	c := make([]float64, n)
	for j := 0; j < n; j++ {
		c[j] = 1.0
	}

	// Aub is A (<= constraints), bub is ones
	Aub := make([][]float64, m)
	bub := make([]float64, m)
	for i := 0; i < m; i++ {
		Aub[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			Aub[i][j] = A[i][j]
		}
		bub[i] = 1.0
	}

	// Solve (maximize)
	ysol, err := SolveLP(c, Aub, bub, true)
	if err != nil {
		fmt.Printf("Player 2 LP failed: %v\n", err)
		return make(Vector, cols)
	}

	y := make(Vector, cols)
	for j := 0; j < cols; j++ {
		y[j] = ysol[j]
	}
	return y
}

// sum calculates the sum of vector elements
func sum(v Vector) float64 {
	s := 0.0
	for _, val := range v {
		s += val
	}
	return s
}

// PrintResult prints the game result
func PrintResult(result GameResult) {
	fmt.Println("=== Game Solution ===")
	fmt.Printf("Player 1 Optimal Strategy: ")
	for i, val := range result.PlayerXStrategy {
		fmt.Printf("x%d=%.4f ", i+1, val)
	}
	fmt.Println()

	fmt.Printf("Player 2 Optimal Strategy: ")
	for i, val := range result.PlayerYStrategy {
		fmt.Printf("y%d=%.4f ", i+1, val)
	}
	fmt.Println()

	fmt.Printf("Game Value: %.4f\n", result.GameValue)
	fmt.Println()
}

func ExampleFromLab() matrix.Matrix {
	return matrix.Matrix{
		{1, 2, -2},
		{-1, 0, 1},
		{1, 1, -1},
	}
}

func ToMatrix(im [][]int) matrix.Matrix {
	m := matrix.NewFromShape(len(im), len(im[0]), 0.0)
	for i := range im {
		for j := range im[i] {
			m[i][j] = float64(im[i][j])
		}
	}
	return m
}

// SolveLP is a simplified wrapper that expects Aub * x <= bub (inequalities only).
// c is the objective coefficients for minimization; set maximize=true to maximize (wrapper negates c).
func SolveLP(c []float64, Aub [][]float64, bub []float64, maximize bool) ([]float64, error) {
	n := len(c)
	if len(Aub) != len(bub) {
		return nil, fmt.Errorf("aub/bub size mismatch")
	}

	// lpsimplex minimizes, so negate objective for maximization
	cc := make([]float64, n)
	copy(cc, c)
	if maximize {
		for j := range n {
			cc[j] = -cc[j]
		}
	}

	maxIter := 1000
	optRes := lpsimplex.LPSimplex(cc, Aub, bub, nil, nil, nil, nil, false, maxIter, 1e-9, false)

	if !optRes.Success {
		return nil, fmt.Errorf("lpsimplex failed to solve LP: %s (status=%d)", optRes.Message, optRes.Status)
	}

	if len(optRes.X) == 0 {
		return nil, fmt.Errorf("lpsimplex returned empty solution")
	}

	x := optRes.X
	if len(x) < n {
		x2 := make([]float64, n)
		copy(x2, x)
		x = x2
	} else if len(x) > n {
		x = x[:n]
	}
	return x, nil
}

func main() {
	fmt.Println("\n### EXAMPLE FROM LAB ###")
	exampleGame := ExampleFromLab()
	result := SolveMatrixGame(exampleGame)
	PrintResult(result)

	fmt.Println("\n### PROBLEM 1: Modified Coin Game ###")
	game1 := ToMatrix(games.CoinGame())
	result1 := SolveMatrixGame(game1)
	PrintResult(result1)

	fmt.Println("\n### PROBLEM 2: Partners Choosing Values ###")
	game2 := ToMatrix(games.Game2_s_vals())
	result2 := SolveMatrixGame(game2)
	PrintResult(result2)

	fmt.Println("\n### PROBLEM 3: Rock, Paper, Scissors ###")
	game3 := ToMatrix(games.RPS())
	result3 := SolveMatrixGame(game3)
	PrintResult(result3)

	fmt.Println("\n### PROBLEM 4: Two-Finger Morra ###")
	game4 := ToMatrix(games.Morra(2))
	result4 := SolveMatrixGame(game4)
	PrintResult(result4)
}
