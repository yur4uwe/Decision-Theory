package grading

import (
	"fmt"
	"math"
	"slices"

	"gonum.org/v1/gonum/stat/distuv"
)

func CalculateAndDisplayProbabilities(scores []float32, norm distuv.Normal) []float32 {
	n := len(scores)

	probMatrix := NewSquareMatrix[float32](n)
	for i := 0; i < n; i++ {
		probMatrix[i] = make([]float32, n)
		for j := 0; j < n; j++ {
			if i == j {
				probMatrix[i][j] = 0.0
			} else {
				z := (float64(scores[i]) - float64(scores[j])) / math.Sqrt(2)
				probMatrix[i][j] = float32(norm.CDF(z))
			}
		}
	}

	overallProbs := make([]float32, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				overallProbs[i] += probMatrix[i][j]
			}
		}
	}

	// Normalize probabilities
	var total float32
	for _, p := range overallProbs {
		total += p
	}
	for i := 0; i < n; i++ {
		overallProbs[i] /= total
	}

	return overallProbs
}

func Thurstone(n, experts int) {
	alts := alternatives(n)

	matrix := NewSquareMatrix[int](n)

	fmt.Println("For each pair (i,j) enter number of experts who prefer i over j.")
	for i := range n {
		for j := i + 1; j < n; j++ {
			var cnt int
			fmt.Printf("How many experts prefer %d over %d? (0..%d): ", alts[i], alts[j], experts)
			_, err := fmt.Scan(&cnt)
			if err != nil || cnt < 0 || cnt > experts {
				fmt.Println("Invalid input, try again.")
				j--
				continue
			}
			matrix[i][j] = cnt
			matrix[j][i] = experts - cnt
		}
	}

	norm := distuv.Normal{Mu: 0, Sigma: 1}

	z := NewSquareMatrix[float32](n)
	for i := range n {
		for j := range n {
			if i == j {
				z[i][j] = 0
				continue
			}

			// Calculate the proportion of experts preferring i over j
			p := (float32(matrix[i][j]) + 0.5) / (float32(experts) + 1.0)
			if p <= 0 {
				p = 1e-12
			} else if p >= 1 {
				p = 1 - 1e-12
			}

			// Get the z-score (strength of preference/standard deviation) for this proportion
			z[i][j] = float32(norm.Quantile(float64(p)))
		}
	}

	s := make([]float32, n)
	for i := range n {
		var sum float32
		for j := range n {
			if i == j {
				continue
			}
			sum += z[i][j]
		}
		s[i] = sum / float32(n-1)
	}

	// center at 0
	var mean float32
	for i := range n {
		mean += s[i]
	}
	mean /= float32(n)
	for i := range n {
		s[i] -= mean
	}

	results := make([]asp, n)
	for i := range n {
		results[i] = asp{alt: alts[i], score: s[i]}
	}

	fmt.Println("\nPairwise Comparison Matrix:")
	for i := range n {
		fmt.Println(matrix[i])
	}

	slices.SortFunc(results, DescendingSort())
	probs := CalculateAndDisplayProbabilities(s, norm)

	fmt.Println("\nThurstone (Case V) scale values (best first):")
	for rank, r := range results {
		fmt.Printf("%d) Alternative %d -> %.4f (%.2f%%)\n", rank+1, r.alt, r.score, probs[r.alt-1]*100)
	}

}
