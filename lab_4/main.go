package main

import (
	"flag"
	"fmt"
	"graph"
	"lab/randanalysis"
	"math/rand"
	"slices"
	"strconv"
)

func randint(min, max int, amount int) []int {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = min + rand.Intn(max-min+1)
	}
	return values
}

func frequency(values []int) ([]int, []int) {
	set := make(map[int]int)
	for _, v := range values {
		set[v]++
	}

	freq := make([]int, 0, len(set))
	vals := make([]int, 0, len(set))
	for k, v := range set {
		vals = append(vals, k)
		freq = append(freq, v)
	}

	return freq, vals
}

func TransformToPillarLabels(vals, freqs []int) []string {
	result := make([]string, len(vals))
	for i, v := range vals {
		result[i] = strconv.Itoa(v) + ": " + strconv.Itoa(freqs[i])
	}
	return result
}

func printfreq(vals, freqs []int) {
	fmt.Println("Number: Frequency")

	for i := 0; i < len(vals); i++ {
		idx := slices.Index(vals, i+1)
		if idx == -1 {
			continue
		}

		fmt.Printf("%d: %d\n", vals[idx], freqs[idx])
	}
}

func main() {
	rangeflag := flag.Int("range", 60, "Range of random numbers")
	amountflag := flag.Int("amount", 3000, "Amount of random numbers")

	flag.Parse()

	// 5 * 12 = 60
	num_range := *rangeflag

	nums := randint(1, num_range, *amountflag)
	numsf := graph.ToFloatSlice(nums)

	fmt.Println("Mean:", randanalysis.Mean(numsf))
	fmt.Println("Mode:", randanalysis.Mode(numsf))
	fmt.Println("Median:", randanalysis.Median(numsf))
	fmt.Println("RMS:", randanalysis.RMS(numsf))
	fmt.Println("Geometric Mean:", randanalysis.GeometricMean(numsf))
	fmt.Println("Harmonic Mean:", randanalysis.HarmonicMean(numsf))

	freq, vals := frequency(nums)
	printfreq(vals, freq)

	freqf := graph.ToFloatSlice(freq)
	valsf := graph.ToFloatSlice(vals)

	g := graph.NewGraph()

	ls := graph.NewLS()
	ls.Pillars(4)

	g.SetLineStyle(ls)

	g.Plot(valsf, freqf, TransformToPillarLabels(vals, freq))

	if err := g.Draw(); err != nil {
		panic(err)
	}

	if err := g.SavePNG("images/output.png"); err != nil {
		panic(err)
	}

	g.Clear()
}
