package main

import (
	"decision-theory/graph"
	"fmt"
	"math"
	"math/rand"
	"slices"
)

const (
	N = 12
)

func NewCoordinates(len uint32) [][2]int {
	coords := make([][2]int, len)
	for i := range len {
		x := rand.Intn(2*(N+3)) - (N + 3)
		y := rand.Intn(2*(N+3)) - (N + 3)

		coords[i] = [2]int{x, y}
	}
	return coords
}

func LinearConvolution(coords [][2]int) (int, []float64) {
	weights := [2]float64{0.3, 0.7}
	sums := make([]float64, len(coords))

	bestSum := 0.0
	bestIdx := -1

	for i, coord := range coords {
		sum := float64(coord[0])*weights[0] + float64(coord[1])*weights[1]
		sums[i] = sum
		if sum > bestSum {
			bestSum = sum
			bestIdx = i
		}
	}

	return bestIdx, sums
}

func IdealPoint(coords [][2]int) (int, []float64) {
	dists := make([]float64, len(coords))
	max_x, max_y := -1, -1

	for _, coord := range coords {
		if coord[0] > max_x {
			max_x = coord[0]
		}
		if coord[1] > max_y {
			max_y = coord[1]
		}
	}

	best_dist := math.Inf(1)
	best_idx := -1

	for i, coord := range coords {
		if max_x == coord[0] && max_y == coord[1] {
			fmt.Println("Ideal point found!")
		}

		diff_x := max_x - coord[0]
		diff_y := max_y - coord[1]

		dists[i] = math.Sqrt(float64(diff_x*diff_x) + float64(diff_y*diff_y))

		if dists[i] < best_dist {
			best_dist = dists[i]
			best_idx = i
		}
	}

	return best_idx, dists
}

func DecoupleCoords(coords [][2]int) ([]float64, []float64) {
	x := make([]float64, len(coords))
	y := make([]float64, len(coords))

	for i, coord := range coords {
		x[i] = float64(coord[0])
		y[i] = float64(coord[1])
	}

	return x, y
}

func DominatedByPareto(alt1, alt2 [2]int) bool {
	alt1_x, alt1_y := alt1[0], alt1[1]
	alt2_x, alt2_y := alt2[0], alt2[1]

	better_or_equal := (alt2_x >= alt1_x) && (alt2_y >= alt1_y)

	strictly_better := (alt2_x > alt1_x) || (alt2_y > alt1_y)

	return better_or_equal && strictly_better
}

func FindparetoOptimal(alts [][2]int) []uint {
	optimals := make([]uint, 0)

	for i, alt1 := range alts {
		is_pareto_optimal := true
		for j, alt2 := range alts {
			if i == j || !DominatedByPareto(alt1, alt2) {
				continue
			}

			is_pareto_optimal = false
		}

		if is_pareto_optimal {
			optimals = append(optimals, uint(i))
		}
	}

	return optimals
}

func PrintTable(alts [][2]int, sums, dists []float64) {
	optimal := FindparetoOptimal(alts)

	header := "Alternative:"
	x_coords := fmt.Sprintf("%-*s", len(header), "X:")
	y_coords := fmt.Sprintf("%-*s", len(header), "Y:")
	convolution_values := fmt.Sprintf("%-*s", len(header), "LCS:")
	ideal_point_dist := fmt.Sprintf("%-*s", len(header), "diff:")
	is_optimal := fmt.Sprintf("%-*s", len(header), "is PO:")

	for i := range alts {
		header += fmt.Sprintf(" |  %-5s", fmt.Sprintf("A%d", i+1))
		x_coords += fmt.Sprintf(" | %-6d", alts[i][0])
		y_coords += fmt.Sprintf(" | %-6d", alts[i][1])
		convolution_values += fmt.Sprintf(" | %-6.2f", sums[i])
		ideal_point_dist += fmt.Sprintf(" | %-6.2f", dists[i])

		if slices.Contains(optimal, uint(i)) {
			is_optimal += " |   po  "
		} else {
			is_optimal += " |       "
		}
	}

	fmt.Println(header)
	fmt.Println(x_coords)
	fmt.Println(y_coords)
	fmt.Println(convolution_values)
	fmt.Println(ideal_point_dist)
	fmt.Println(is_optimal)
}

func main() {
	alternatives := NewCoordinates(10)
	labels := make([]string, len(alternatives))
	for i := range len(alternatives) {
		labels[i] = fmt.Sprintf("A%d", i+1)
	}

	lc_best_idx, sums := LinearConvolution(alternatives)
	if lc_best_idx < 0 {
		panic("no solution found with linear convolution")
	}

	ip_best_idx, dists := IdealPoint(alternatives)
	if ip_best_idx < 0 {
		panic("no solution found with ideat point")
	}

	PrintTable(alternatives, sums, dists)

	fmt.Printf("Best alternative according to linear convolution: A%d\n", lc_best_idx+1)
	fmt.Printf("Best alternative according to ideal point: A%d", ip_best_idx+1)

	g := graph.NewGraph()
	ls := graph.NewLS()
	ls.Dots(4)

	x, y := DecoupleCoords(alternatives)

	g.Plot(x, y, ls, labels)

	if err := g.Draw(); err != nil {
		panic(err)
	}

	if err := g.SavePNG("images/alternatives.png"); err != nil {
		panic(err)
	}
}
