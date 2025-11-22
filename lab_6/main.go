package main

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	N = 12
)

// Alternative представляє альтернативу з її критеріальними оцінками
type Alternative struct {
	ID       int
	Criteria []float64
}

// Функція для генерації випадкового числа в діапазоні [min, max]
func randomInRange(min, max int) float64 {
	return float64(min + rand.Intn(max-min+1))
}

// Генерація альтернатив з випадковими критеріальними оцінками
func generateAlternatives(numAlternatives, numCriteria, maxValue int) []Alternative {
	alternatives := make([]Alternative, numAlternatives)
	for i := 0; i < numAlternatives; i++ {
		alternatives[i] = Alternative{
			ID:       i + 1,
			Criteria: make([]float64, numCriteria),
		}
		for j := 0; j < numCriteria; j++ {
			alternatives[i].Criteria[j] = randomInRange(1, maxValue)
		}
	}
	return alternatives
}

// Генерація нормативних значень критеріїв
func generateNormativeValues(numCriteria, maxValue int) []float64 {
	normative := make([]float64, numCriteria)
	for i := 0; i < numCriteria; i++ {
		normative[i] = randomInRange(1, maxValue)
	}
	return normative
}

// Генерація вагових коефіцієнтів (сума = 1)
func generateWeights(numCriteria int) []float64 {
	weights := make([]float64, numCriteria)
	sum := 0.0

	// Генеруємо випадкові числа
	for i := 0; i < numCriteria; i++ {
		weights[i] = rand.Float64()
		sum += weights[i]
	}

	// Нормалізуємо, щоб сума дорівнювала 1
	for i := 0; i < numCriteria; i++ {
		weights[i] /= sum
	}

	return weights
}

// Лінійна згортка (без нормування)
func linearConvolution(alternatives []Alternative, weights []float64) int {
	maxValue := -1.0
	bestAlternative := -1

	fmt.Println("\nLinear Convolution")
	fmt.Printf("Weights: ")
	for i, w := range weights {
		fmt.Printf("λ%d=%.3f ", i+1, w)
	}
	fmt.Println()

	for _, alt := range alternatives {
		sum := 0.0
		for i, criterion := range alt.Criteria {
			sum += weights[i] * criterion
		}

		fmt.Printf("A%d: Q = %.3f\n", alt.ID, sum)

		if sum > maxValue {
			maxValue = sum
			bestAlternative = alt.ID
		}
	}

	return bestAlternative
}

// Нормована лінійна згортка
func normalizedLinearConvolution(alternatives []Alternative, weights []float64) int {
	numCriteria := len(alternatives[0].Criteria)

	// find min and max for each criterion
	minValues := make([]float64, numCriteria)
	maxValues := make([]float64, numCriteria)

	for i := 0; i < numCriteria; i++ {
		minValues[i] = alternatives[0].Criteria[i]
		maxValues[i] = alternatives[0].Criteria[i]

		for _, alt := range alternatives {
			if alt.Criteria[i] < minValues[i] {
				minValues[i] = alt.Criteria[i]
			}
			if alt.Criteria[i] > maxValues[i] {
				maxValues[i] = alt.Criteria[i]
			}
		}
	}

	maxValue := -1.0
	bestAlternative := -1

	fmt.Println("\nNormalized Linear Convolution")
	fmt.Printf("Weights: ")
	for i, w := range weights {
		fmt.Printf("λ%d=%.3f ", i+1, w)
	}
	fmt.Println()

	for _, alt := range alternatives {
		sum := 0.0
		for i, criterion := range alt.Criteria {
			normalized := 0.0
			if maxValues[i] != minValues[i] {
				normalized = (criterion - minValues[i]) / (maxValues[i] - minValues[i])
			}
			sum += weights[i] * normalized
		}

		fmt.Printf("A%d: Q = %.3f\n", alt.ID, sum)

		if sum > maxValue {
			maxValue = sum
			bestAlternative = alt.ID
		}
	}

	return bestAlternative
}

func maximinConvolution(alternatives []Alternative, normative []float64) int {
	maxMinRatio := -1.0
	bestAlternative := -1

	fmt.Println("\nMaximin Convolution")
	fmt.Printf("Normative values: ")
	for i, n := range normative {
		fmt.Printf("Q%d*=%.1f ", i+1, n)
	}
	fmt.Println()

	for _, alt := range alternatives {
		minRatio := 1e9
		for i, criterion := range alt.Criteria {
			ratio := criterion / normative[i]
			if ratio < minRatio {
				minRatio = ratio
			}
		}

		fmt.Printf("A%d: min(Qi/Qi*) = %.3f\n", alt.ID, minRatio)

		if minRatio > maxMinRatio {
			maxMinRatio = minRatio
			bestAlternative = alt.ID
		}
	}

	return bestAlternative
}

// Нормована максимінна згортка
func normalizedMaximinConvolution(alternatives []Alternative, normative []float64, weights []float64) int {
	maxMinValue := -1.0
	bestAlternative := -1

	fmt.Println("\nNormalized Maximin Convolution")
	fmt.Printf("Normative values: ")
	for i, n := range normative {
		fmt.Printf("Q%d*=%.1f ", i+1, n)
	}
	fmt.Println()
	fmt.Printf("Weights: ")
	for i, w := range weights {
		fmt.Printf("λ%d=%.3f ", i+1, w)
	}
	fmt.Println()

	for _, alt := range alternatives {
		minValue := 1e9
		for i, criterion := range alt.Criteria {
			// Зважене відношення: λi * (Qi / Qi*)
			weightedRatio := weights[i] * (criterion / normative[i])
			if weightedRatio < minValue {
				minValue = weightedRatio
			}
		}

		fmt.Printf("A%d: min(λi*Qi/Qi*) = %.3f\n", alt.ID, minValue)

		if minValue > maxMinValue {
			maxMinValue = minValue
			bestAlternative = alt.ID
		}
	}

	return bestAlternative
}

func printAlternativesTable(alternatives []Alternative) {
	numCriteria := len(alternatives[0].Criteria)

	fmt.Printf("%s", strings.Repeat(" ", 10))
	for _, alt := range alternatives {
		fmt.Printf("A%-8d", alt.ID)
	}
	fmt.Println()
	fmt.Println(string(make([]byte, 10+len(alternatives)*9)))

	for i := 0; i < numCriteria; i++ {
		fmt.Printf("Q%-9d", i+1)
		for _, alt := range alternatives {
			fmt.Printf("%-9.1f", alt.Criteria[i])
		}
		fmt.Println()
	}
}

func main() {
	numAlternatives := 3 + N
	numCriteria := 5
	maxValue := 3 + N

	alternatives := generateAlternatives(numAlternatives, numCriteria, maxValue)
	normative := generateNormativeValues(numCriteria, maxValue)
	weights := generateWeights(numCriteria)

	printAlternativesTable(alternatives)

	bestLinear := linearConvolution(alternatives, weights)
	fmt.Printf("\nOptimal: A%d\n", bestLinear)

	bestMaximin := maximinConvolution(alternatives, normative)
	fmt.Printf("\nOptimal: A%d\n", bestMaximin)

	bestNormalizedLinear := normalizedLinearConvolution(alternatives, weights)
	fmt.Printf("\nOptimal: A%d\n", bestNormalizedLinear)

	bestNormalizedMaximin := normalizedMaximinConvolution(alternatives, normative, weights)
	fmt.Printf("\nOptimal: A%d\n", bestNormalizedMaximin)
}
