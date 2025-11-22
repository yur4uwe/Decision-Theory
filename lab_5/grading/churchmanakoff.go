package grading

import (
	"fmt"
	"slices"
)

// CheckTransitivity checks if the pairwise comparison matrix is transitive.
// Returns true if transitive, false otherwise.
func CheckTransitivity(matrix [][]int, n int) bool {
	for i := range n {
		for j := range n {
			for k := range n {
				if matrix[i][j] == 1 && matrix[j][k] == 1 && matrix[i][k] != 1 {
					return false
				}
			}
		}
	}
	return true
}

// ChurchmanAckoff performs pairwise comparisons and ensures transitivity.
func ChurchmanAckoff(n int) {
	alts := alternatives(n)

	matrix := NewSquareMatrix[int](n)

	fmt.Println("\nPerforming pairwise comparisons:")
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			fmt.Printf("Which alternative is better? (1 for '%d', 2 for '%d'): ", alts[i], alts[j])
			var choice int
			fmt.Scan(&choice)
			switch choice {
			case 1:
				matrix[i][j] = 1
				matrix[j][i] = 0
			case 2:
				matrix[i][j] = 0
				matrix[j][i] = 1
			default:
				fmt.Println("Invalid input. Please enter 1 or 2.")
				j--
			}
		}
	}

	for !CheckTransitivity(matrix, n) {
		fmt.Println("\nThe matrix is not transitive. Please reevaluate your judgments.")
		for i := range n {
			for j := i + 1; j < n; j++ {
				for k := range n {
					if matrix[i][j] == 1 && matrix[j][k] == 1 && matrix[i][k] != 1 {
						fmt.Printf("Reevaluate: If '%d' > '%d' and '%d' > '%d', should '%d' > '%d'? (1 for Yes, 2 for No): ",
							alts[i], alts[j], alts[j], alts[k], alts[i], alts[k])
						var choice int
						fmt.Scan(&choice)
						switch choice {
						case 1:
							matrix[i][k] = 1
							matrix[k][i] = 0
						case 2:
							matrix[i][k] = 0
							matrix[k][i] = 1
						default:
							fmt.Println("Invalid input. Please enter 1 or 2.")
							k--
						}
					}
				}
			}
		}
	}

	scores := make([]float32, n)
	for i := range n {
		for j := range n {
			scores[i] += float32(matrix[i][j])
		}
	}

	rankings := make([]asp, n)
	for i := range n {
		rankings[i] = asp{alt: alts[i], score: scores[i]}
	}

	slices.SortFunc(rankings, DescendingSort())

	fmt.Println("\nPairwise Comparison Matrix:")
	for i := range n {
		fmt.Println(matrix[i])
	}

	fmt.Println("\nFinal Rankings:")
	for i, r := range rankings {
		fmt.Printf("%d. %d (Score: %.2f)\n", i+1, r.alt, r.score)
	}
}
