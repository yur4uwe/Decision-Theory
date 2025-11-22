package grading

import "slices"

func alternatives(n int) []int {
	alts := make([]int, n)
	for i := 0; i < n; i++ {
		alts[i] = i + 1
	}
	return alts
}

func removeAlternative(alts []int, alt int) []int {
	newAlts := []int{}
	for _, a := range alts {
		if a != alt {
			newAlts = append(newAlts, a)
		}
	}
	return newAlts
}

func Arrange(alternatives []int, gradingFunc func(int) float32) []int {
	slices.SortFunc(alternatives, func(a, b int) int {
		if gradingFunc(a) < gradingFunc(b) {
			return -1
		} else if gradingFunc(a) > gradingFunc(b) {
			return 1
		}
		return 0
	})

	return alternatives
}

func NewGradingFunc(alternatives []int, grades []float32) func(int) float32 {
	return func(alternative int) float32 {
		return grades[slices.Index(alternatives, alternative)]
	}
}

func NewSquareMatrix[T comparable](n int) [][]T {
	matrix := make([][]T, n)
	for i := range matrix {
		matrix[i] = make([]T, n)
	}
	return matrix
}
