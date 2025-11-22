package main

import (
	"fmt"

	"lab/games"
)

func saddlePoint(m [][]int) (bool, int, [][2]int) {
	// empty matrix
	if len(m) == 0 || len(m[0]) == 0 {
		return false, 0, nil
	}

	rows := len(m)
	cols := len(m[0])

	// compute row minima
	rowMins := make([]int, rows)
	for i := 0; i < rows; i++ {
		min := m[i][0]
		for j := 1; j < cols; j++ {
			if m[i][j] < min {
				min = m[i][j]
			}
		}
		rowMins[i] = min
	}

	// compute max of row minima (maximin)
	maxOfRowMins := rowMins[0]
	for i := 1; i < rows; i++ {
		if rowMins[i] > maxOfRowMins {
			maxOfRowMins = rowMins[i]
		}
	}

	// compute column maxima (fixed loop)
	colMaxs := make([]int, cols)
	for j := 0; j < cols; j++ {
		max := m[0][j]
		for i := 1; i < rows; i++ {
			if m[i][j] > max {
				max = m[i][j]
			}
		}
		colMaxs[j] = max
	}

	// compute min of column maxima (minimax)
	minOfColMaxs := colMaxs[0]
	for j := 1; j < cols; j++ {
		if colMaxs[j] < minOfColMaxs {
			minOfColMaxs = colMaxs[j]
		}
	}

	// saddle exists only if maximin == minimax
	if maxOfRowMins != minOfColMaxs {
		return false, 0, nil
	}

	value := maxOfRowMins
	positions := make([][2]int, 0)

	// collect all positions that are equilibrium (value, row-min and col-max)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if m[i][j] == value && m[i][j] == rowMins[i] && m[i][j] == colMaxs[j] {
				positions = append(positions, [2]int{i, j})
			}
		}
	}

	return true, value, positions
}

func printResult(name string, m [][]int) {
	if len(m) == 0 || len(m[0]) == 0 {
		fmt.Println("No valid matrix provided.")
		return
	}

	ok, val, pos := saddlePoint(m)
	fmt.Println("----", name, "----")
	fmt.Printf("Matrix %dx%d\n", len(m), len(m[0]))

	if !ok {
		fmt.Println("No saddle point (no pure-strategy equilibrium).")
	} else {
		fmt.Printf("Saddle point(s) found. Game value (payoff to row player) = %d\n", val)
		fmt.Println("Equilibrium positions (row, column) — zero-based indices and payoffs:")
		for _, p := range pos {
			i, j := p[0], p[1]
			fmt.Printf("(%d, %d)  payoff = %d\n", i, j, m[i][j])
		}
	}
	fmt.Println()
}

func main() {
	printResult("1. Coin simple (H/T) game", games.CoinGame())
	printResult("2. s,t in {-1,0,1}", games.Game2_s_vals())
	printResult("3. Rock-Paper-Scissors", games.RPS())
	printResult("4. Morra Two-Finger (4x4)", games.Morra(2))
	printResult("5. Morra Three-Finger (9x9)", games.Morra(3))
	printResult("6. Integers 1..k (k=4)", games.Game6(4))
	printResult("7. Colonel Blotto (attacker=3, defender=3)", games.Blotto(3, 3))
	printResult("8. Seller problem (k=5,a=10,b=4,alpha=0,beta=5)", games.SellerProblem(5, 10, 4, 0, 5))
	printResult("custom 1. Prisoners Dilemma", games.PrisonersDilemma())
	printResult("custom 2. Chicken(Hawk-Dove)", games.ChickenDefault())

	printResult("Test Matrix with 2 saddle points", [][]int{{1, 2}, {1, 2}})
}
