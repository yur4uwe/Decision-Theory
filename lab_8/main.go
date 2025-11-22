package main

import (
	"fmt"
	"math/rand"
)

const (
	N = 12
)

func linarr(n int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = i
	}

	return arr
}

func LexicographyOptimisationWithPossibleYields(alts [][]int, criterionOrder []int, acceptableYield []int) int {
	if (len(acceptableYield) != 0 && len(acceptableYield) != len(criterionOrder)) || (len(alts[0]) != len(criterionOrder)) {
		panic("Incorrect input")
	}

	if len(acceptableYield) == 0 {
		acceptableYield = make([]int, len(criterionOrder))
	}

	availableAlts := linarr(len(alts))

	for i, criterion := range criterionOrder {
		fmt.Printf("Step %d: available alternatives: %v\n", i, availableAlts)
		if len(availableAlts) <= 1 {
			break
		}

		maxCriterionVal := 0
		for _, remainingAlt := range availableAlts {
			if alts[remainingAlt][criterion] > maxCriterionVal {
				maxCriterionVal = alts[remainingAlt][criterion]
			}
		}

		remainingAlts := []int{}

		for _, remainingAlt := range availableAlts {
			if alts[remainingAlt][criterion] >= maxCriterionVal-acceptableYield[criterion] {
				remainingAlts = append(remainingAlts, remainingAlt)
			}
		}

		availableAlts = remainingAlts
	}

	return availableAlts[0]
}

func CreateAlternatives(n, k int) [][]int {
	alts := make([][]int, n)

	for i := range n {
		alts[i] = make([]int, k)
		for j := range k {
			alts[i][j] = rand.Intn(6)
		}
	}

	return alts
}

func PrintTable(alts [][]int) {
	header := "Alternative:"

	criteria := make([]string, len(alts[0]))

	for i := range alts[0] {
		criteria[i] = fmt.Sprintf("%-*s", len(header), fmt.Sprintf("Q%d", i))
	}

	for i := range alts {
		header += fmt.Sprintf(" |  %-3s", fmt.Sprintf("A%d", i+1))
		for j := range alts[0] {
			criteria[j] += fmt.Sprintf(" | %-4d", alts[i][j])
		}
	}

	fmt.Println(header)
	for i := range alts[0] {
		fmt.Println(criteria[i])
	}
}

func main() {
	alternatives := CreateAlternatives(N+3, 4)
	PrintTable(alternatives)

	order1 := []int{1, 0, 2, 3}
	order2 := []int{3, 0, 1, 2}
	subsequentConcessionsOrder := []int{0, 1, 2, 3}

	yields := []int{1, 2, 3, 0}

	optimal1 := LexicographyOptimisationWithPossibleYields(alternatives, order1, nil)
	fmt.Printf("Lexicographically Optimal for Q order %v is A%d\n", order1, optimal1+1)
	optimal2 := LexicographyOptimisationWithPossibleYields(alternatives, order2, nil)
	fmt.Printf("Lexicographically Optimal for Q order %v is A%d\n", order2, optimal2+1)

	optimalWithYields := LexicographyOptimisationWithPossibleYields(alternatives, subsequentConcessionsOrder, yields)
	fmt.Printf("Optimal with Yields for Q order %v and acceptable yields %v is A%d\n", subsequentConcessionsOrder, yields, optimalWithYields+1)

}
