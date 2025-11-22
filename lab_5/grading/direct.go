package grading

import (
	"fmt"
	"slices"
)

func Direct(n int) {
	evaluations := make([]asp, n)
	for i := 0; i < n; i++ {
		fmt.Printf("Enter the evaluation for alternative %d (0-1): ", i+1)
		var eval float32
		fmt.Scan(&eval)
		if eval < 0 || eval > 1 {
			fmt.Println("Evaluation must be between 0 and 1.")
			i--
		} else {
			evaluations[i] = asp{alt: i + 1, score: eval}
		}
	}

	slices.SortFunc(evaluations, AscendingSort())

	fmt.Println("Arranged Evaluations (best -> worst):")
	for _, eval := range evaluations {
		fmt.Printf("Alternative %d: %.2f\n", eval.alt, eval.score)
	}
}
