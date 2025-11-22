package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
)

var (
	RandomConsistencyIndex = map[int]float64{
		1:  0.00,
		2:  0.00,
		3:  0.58,
		4:  0.90,
		5:  1.12,
		6:  1.24,
		7:  1.32,
		8:  1.41,
		9:  1.45,
		10: 1.49,
	}

	ExampleMatrixPWCriteria = [][]float64{
		{1, 4, 0.2, 5, 1.0 / 3.0},
		{1.0 / 4.0, 1, 1.0 / 7.0, 1.0 / 2.0, 1.0 / 3.0},
		{5, 7, 1, 6, 3},
		{1.0 / 5.0, 2, 1.0 / 6.0, 1, 1.0 / 5.0},
		{3, 3, 1.0 / 3.0, 5, 1},
	}

	ExampleMatricesPWAlternatives = [][][]float64{
		{
			{1, 2, 1.0 / 3.0, 2},
			{1.0 / 2.0, 1, 1.0 / 5.0, 2},
			{3, 5, 1, 2},
			{1.0 / 2.0, 1.0 / 2.0, 1.0 / 2.0, 1},
		},
		{
			{1, 1.0 / 7.0, 1.0 / 5.0, 1.0 / 2.0},
			{7, 1, 6, 2},
			{5, 1.0 / 6.0, 1, 1.0 / 3.0},
			{2, 1.0 / 2.0, 3, 1},
		},
		{
			{1, 3, 4, 1},
			{1.0 / 3.0, 1, 3, 2},
			{1.0 / 4.0, 1.0 / 3.0, 1, 1.0 / 5.0},
			{1, 1.0 / 2.0, 5, 1},
		},
		{
			{1, 1.0 / 5.0, 4, 1.0 / 3.0},
			{5, 1, 6, 2},
			{1.0 / 4.0, 1.0 / 6.0, 1, 1.0 / 7.0},
			{3, 1.0 / 2.0, 7, 1},
		},
		{
			{1, 3, 1.0 / 5.0, 1.0 / 7.0},
			{1.0 / 3.0, 1, 1.0 / 3.0, 1.0 / 8.0},
			{5, 3, 1, 1.0 / 7.0},
			{7, 8, 7, 1},
		},
	}

	// alternatives:
	// A1 - HP EliteBook 840 G7: 470, 53 W/h, Intel Core i5-10310U, 16GB RAM, 512GB SSD, 14" FHD, almost new, 1.34 kg
	// A2 - Lenovo ThinkPad E14 Gen 2: 420, i5-1135G7, 16GB RAM, 512GB NVMe, 14" FHD, used, 1.59 kg
	// A3 - Acer Aspire Lite: 270, 50 W/h, Intel Core i3-1334U, 16GB RAM, 512 SSD, 16" FHD, new, 1.53 kg
	// A4 - Dell 15 DC15255: 510, 41 W/h, AMD Ryzen 3 7320U, 8GB RAM, 512GB SSD, 15.6" FHD, new, 1.9 kg
	// A5 - Asus VivoBook 15 X515JA: 590, 42 W/h, Intel Core i3-1315U, 12GB RAM, 512GB SSD, 15.6" FHD, new, 1.7 kg

	// criteria (order):
	// price - battery life - cpu - ram - storage - size - state - weight
	NoteBookMatrixPWCriteria = [][]float64{
		// price, battery, cpu, ram, storage, size, state, weight
		{1, 3, 5, 7, 9, 9, 9, 9},
		{1.0 / 3.0, 1, 3, 5, 7, 9, 9, 9},
		{1.0 / 5.0, 1.0 / 3.0, 1, 3, 5, 7, 9, 9},
		{1.0 / 7.0, 1.0 / 5.0, 1.0 / 3.0, 1, 3, 5, 7, 9},
		{1.0 / 9.0, 1.0 / 7.0, 1.0 / 5.0, 1.0 / 3.0, 1, 3, 5, 7},
		{1.0 / 9.0, 1.0 / 9.0, 1.0 / 7.0, 1.0 / 5.0, 1.0 / 3.0, 1, 3, 5},
		{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0, 1.0 / 7.0, 1.0 / 5.0, 1.0 / 3.0, 1, 3},
		{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0, 1.0 / 7.0, 1.0 / 5.0, 1.0 / 3.0, 1},
	}
	NoteBookMatricesPWAlternatives = [][][]float64{
		// 1) price (lower better): A1,A2,A3,A4,A5
		{
			{1, 1.0 / 3.0, 1.0 / 4.0, 3, 5},
			{3, 1, 1.0 / 3.0, 5, 7},
			{4, 3, 1, 5, 6},
			{1.0 / 3.0, 1.0 / 5.0, 1.0 / 5.0, 1, 3},
			{1.0 / 5.0, 1.0 / 7.0, 1.0 / 6.0, 1.0 / 3.0, 1},
		},
		// 2) battery life (higher better): A1,A2,A3,A4,A5  (A2≈A3)
		{
			{1, 3, 3, 7, 5},
			{1.0 / 3.0, 1, 1, 5, 3},
			{1.0 / 3.0, 1, 1, 5, 3},
			{1.0 / 7.0, 1.0 / 5.0, 1.0 / 5.0, 1, 1.0 / 3.0},
			{1.0 / 5.0, 1.0 / 3.0, 1.0 / 3.0, 3, 1},
		},
		// 3) CPU (A3 > A2 > A1 > A4 > A5)
		{
			{1, 1.0 / 3.0, 1.0 / 5.0, 3, 5},
			{3, 1, 1.0 / 3.0, 5, 7},
			{5, 3, 1, 7, 9},
			{1.0 / 3.0, 1.0 / 5.0, 1.0 / 7.0, 1, 3},
			{1.0 / 5.0, 1.0 / 7.0, 1.0 / 9.0, 1.0 / 3.0, 1},
		},
		// 4) RAM (A1=A2=A3 > A5 > A4)
		{
			{1, 1, 1, 5, 7},
			{1, 1, 1, 5, 7},
			{1, 1, 1, 5, 7},
			{1.0 / 5.0, 1.0 / 5.0, 1.0 / 5.0, 1, 3},
			{1.0 / 7.0, 1.0 / 7.0, 1.0 / 7.0, 1.0 / 3.0, 1},
		},
		// 5) Storage (all 512GB -> equal)
		{
			{1, 1, 1, 1, 1},
			{1, 1, 1, 1, 1},
			{1, 1, 1, 1, 1},
			{1, 1, 1, 1, 1},
			{1, 1, 1, 1, 1},
		},
		// 6) Size (smaller better: A1=A2 (14") > A4=A5 (15.6") > A3 (16"))
		{
			{1, 1, 5, 3, 3},
			{1, 1, 5, 3, 3},
			{1.0 / 5.0, 1.0 / 5.0, 1, 1.0 / 3.0, 1.0 / 3.0},
			{1.0 / 3.0, 1.0 / 3.0, 3, 1, 1},
			{1.0 / 3.0, 1.0 / 3.0, 3, 1, 1},
		},
		// 7) State (condition: new (A3,A4,A5) > almost new (A1) > used (A2))
		{
			{1, 5, 1.0 / 3.0, 1.0 / 3.0, 1.0 / 3.0},
			{1.0 / 5.0, 1, 1.0 / 7.0, 1.0 / 7.0, 1.0 / 7.0},
			{3, 7, 1, 1, 1},
			{3, 7, 1, 1, 1},
			{3, 7, 1, 1, 1},
		},
		// 8) Weight (lighter better: A1 > A3 > A2 > A5 > A4)
		{
			{1, 5, 3, 9, 7},
			{1.0 / 5.0, 1, 1.0 / 3.0, 5, 3},
			{1.0 / 3.0, 3, 1, 7, 5},
			{1.0 / 9.0, 1.0 / 5.0, 1.0 / 7.0, 1, 1.0 / 3.0},
			{1.0 / 7.0, 1.0 / 3.0, 1.0 / 5.0, 3, 1},
		},
	}
)

const EPS = 1e-9

func GeometricMean(x []float64) float64 {
	sum := 0.0
	for _, v := range x {
		sum += math.Log(v)
	}
	return math.Exp(sum / float64(len(x)))
}

func LocalPriorityScores(matrix [][]float64) []float64 {
	n := len(matrix)
	m := len(matrix[0])
	means := make([]float64, n)

	for i := range n {
		means[i] = GeometricMean(matrix[i])
	}

	normalization_factor := 0.0
	for _, v := range means {
		normalization_factor += v
	}

	for i := range m {
		if math.Abs(normalization_factor) < EPS {
			means[i] = 0
		} else {
			means[i] /= normalization_factor
		}
	}

	return means
}

func RandPairwiseQualityScore(n int) [][]float64 {
	result := make([][]float64, n)

	for i := range n {
		result[i] = make([]float64, n)
	}

	for i := range n {
		for j := i; j < n; j++ {
			if i == j {
				result[i][j] = 1.0
				continue
			}

			to_mirror := rand.Intn(2) == 1
			score := float64(rand.Intn(9) + 1)
			if to_mirror {
				score = 1 / score
			}

			result[i][j] = score
			result[j][i] = 1 / score
		}
	}

	return result
}

func MaxSelfValue(matrix [][]float64, scores []float64) (val float64) {
	n := len(matrix)
	v := make([]float64, n)
	for i := range n {
		row := matrix[i]
		sum := 0.0
		for j := 0; j < len(row) && j < len(scores); j++ {
			sum += row[j] * scores[j]
		}
		v[i] = sum
	}
	// lambda = (1/n) * sum_i (v_i / w_i)
	for i := range n {
		w := scores[i]
		if math.Abs(w) < EPS {
			w = EPS
		}
		val += v[i] / w
	}
	return val / float64(n)
}

func ConsistencyIndex(self_max_value, n float64) float64 {
	return (self_max_value - n) / (n - 1.0)
}

func Composition(m [][]float64, v []float64) []float64 {
	n := len(m)
	result := make([]float64, n)
	for i := range n {
		sum := 0.0
		for j := 0; j < len(m[i]) && j < len(v); j++ {
			sum += m[i][j] * v[j]
		}
		result[i] = sum
	}
	return result
}

type Pair struct {
	Index int
	Value float64
}

func Arrange(unarranged []float64) []Pair {
	result := make([]Pair, len(unarranged))

	for i, v := range unarranged {
		result[i] = Pair{Index: i, Value: v}
	}

	slices.SortFunc(result, func(a, b Pair) int {
		if a.Value < b.Value {
			return 1
		} else if a.Value > b.Value {
			return -1
		}
		return 0
	})

	return result
}

// DisplayHierarchy prints the decision hierarchy: goal -> criteria (weights) -> alternatives (local scores)
func DisplayHierarchy(goal string, criteriaNames []string, altNames []string, criteriaWeights []float64, localScores [][]float64) {
	fmt.Println("\nDecision Hierarchy:")
	fmt.Printf("Goal: %s\n", goal)

	for i, cname := range criteriaNames {
		w := 0.0
		if i < len(criteriaWeights) {
			w = criteriaWeights[i]
		}
		fmt.Printf("- Criterion %d: %s (weight: %.4f)\n", i+1, cname, w)
		// list alternatives and their scores for this criterion
		for j, aname := range altNames {
			score := 0.0
			if j < len(localScores) && i < len(localScores[j]) {
				score = localScores[j][i]
			}
			fmt.Printf("    - A%d: %s => local score: %.4f\n", j+1, aname, score)
		}
	}
}

func PrintMatrix(name string, m [][]float64) {
	fmt.Println(name)
	for i := range m {
		fmt.Print("  ")
		for j := range m[i] {
			fmt.Printf("%8.4f ", m[i][j])
		}
		fmt.Println()
	}
}

func PrintVector(name string, v []float64) {
	fmt.Print(name + ":")
	for _, val := range v {
		fmt.Printf(" %.6f", val)
	}
	fmt.Println()
}

func main() {
	ALTS := len(NoteBookMatricesPWAlternatives[0])
	CRIT := len(NoteBookMatrixPWCriteria)

	matrix := NoteBookMatrixPWCriteria
	scores := LocalPriorityScores(matrix)

	// print criteria pairwise matrix and criteria weight vector
	PrintMatrix("Criteria pairwise matrix:", matrix)
	PrintVector("Criteria weights (local priority vector):", scores)

	self_max_value := MaxSelfValue(matrix, scores)
	consistency_index := ConsistencyIndex(self_max_value, float64(len(matrix)))

	consistency_ratio := consistency_index / RandomConsistencyIndex[len(matrix)]

	if consistency_ratio <= 0.1 {
		fmt.Println("Saati consistent with CR:", consistency_ratio)
	} else if consistency_ratio <= 0.2 {
		fmt.Println("Saati barely consistent with CR:", consistency_ratio)
	} else {
		fmt.Println("Saati inconsistent with CR:", consistency_ratio)
	}

	// human-readable names for criteria and alternatives
	criteriaNames := []string{"price", "battery", "cpu", "ram", "storage", "size", "state", "weight"}
	altNames := []string{"HP EliteBook 840 G7", "Lenovo ThinkPad E14 Gen 2", "Acer Aspire Lite", "Dell 15 DC15255", "Asus VivoBook 15 X515JA"}

	local_priority_scores := make([][]float64, ALTS)
	for i := range ALTS {
		local_priority_scores[i] = make([]float64, CRIT)
	}

	for i := range CRIT {
		criterion_matrix := NoteBookMatricesPWAlternatives[i]

		// print pairwise matrix for this criterion
		if i < len(criteriaNames) {
			PrintMatrix(fmt.Sprintf("Pairwise matrix for criterion %d (%s):", i+1, criteriaNames[i]), criterion_matrix)
		} else {
			PrintMatrix(fmt.Sprintf("Pairwise matrix for criterion %d:", i+1), criterion_matrix)
		}

		criterion_scores := LocalPriorityScores(criterion_matrix)

		// print local priority vector for this criterion
		PrintVector(fmt.Sprintf("Local priority vector for criterion %d:", i+1), criterion_scores)

		for j := range ALTS {
			local_priority_scores[j][i] = criterion_scores[j]
		}

		self_max_value := MaxSelfValue(criterion_matrix, criterion_scores)
		consistency_index := ConsistencyIndex(self_max_value, float64(len(criterion_matrix)))

		consistency_ratio := consistency_index / RandomConsistencyIndex[len(criterion_matrix)]

		if consistency_ratio <= 0.1 {
			fmt.Println("Criterion", i+1, "Saati consistent with CR:", consistency_ratio)
		} else if consistency_ratio <= 0.2 {
			fmt.Println("Criterion", i+1, "Saati barely consistent with CR:", consistency_ratio)
		} else {
			fmt.Println("Criterion", i+1, "Saati inconsistent with CR:", consistency_ratio)
		}
	}

	// Display the decision hierarchy (goal -> criteria -> alternatives)
	DisplayHierarchy("Select best notebook", criteriaNames, altNames, scores, local_priority_scores)

	global_scores := Composition(local_priority_scores, scores)
	PrintVector("Global priority vector (unsorted)", global_scores)

	fmt.Println("Local priority scores:")
	for i := range local_priority_scores {
		fmt.Printf("A%d:", i+1)
		for j := range local_priority_scores[i] {
			fmt.Printf(" %.3f", local_priority_scores[i][j])
		}
		fmt.Println()
	}

	arranged := Arrange(global_scores)

	fmt.Println("Global priority scores:")
	for _, pair := range arranged {
		fmt.Printf("A%d: %.3f\n", pair.Index+1, pair.Value)
	}
}
