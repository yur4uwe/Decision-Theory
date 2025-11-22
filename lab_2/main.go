package main

import (
	"decision-theory/binrels"
	"flag"
	"fmt"
)

const (
	Fall = iota
	Winter
	Summer
	Spring
)

func set(matrices ...[][]bool) func(int) [][]bool {
	return func(flag int) [][]bool {
		if flag < 0 || flag > len(matrices) {
			return nil
		}
		return matrices[flag]
	}
}

func main() {
	to_print := flag.Bool("p", false, "print matrices")

	intersect := flag.Bool("i", false, "intersect matrices")
	union := flag.Bool("u", false, "union matrices")
	diff := flag.Bool("d", false, "difference matrices")
	symmDiff := flag.Bool("sd", false, "symmetric difference matrices")
	composition := flag.Bool("c", false, "composition matrices")

	reverse := flag.Int("r", 0, "reverse matrix")
	complement := flag.Int("cm", 0, "complement matrix")
	top_int := flag.Int("ti", -1, "top intersection of matrix")
	bottom_int := flag.Int("bi", -1, "bottom intersection of matrix")
	def_dom := flag.Int("dd", -1, "definition domain of matrix")
	mean_dom := flag.Int("md", -1, "meaning domain of matrix")

	runAll := flag.Bool("all", false, "run all operations")

	flag.Parse()

	rels := []string{"О", "З", "Л", "В"}

	is_same := binrels.Zero(len(rels))
	not_colder := binrels.Zero(len(rels))

	is_same[Fall][Fall] = true
	is_same[Winter][Winter] = true
	is_same[Summer][Summer] = true
	is_same[Spring][Spring] = true

	not_colder[Fall][Winter] = true

	not_colder[Spring][Fall] = true
	not_colder[Spring][Winter] = true

	not_colder[Summer][Fall] = true
	not_colder[Summer][Spring] = true
	not_colder[Summer][Winter] = true

	not_colder = binrels.Union(not_colder, is_same)

	getMatrix := set(is_same, not_colder)

	if *runAll || *to_print {
		fmt.Println("Is Same:")
		binrels.PrintWithSource(rels, is_same)
		fmt.Println("\nNot Colder:")
		binrels.PrintWithSource(rels, not_colder)
	}

	if *runAll || *intersect {
		res := binrels.Intersection(is_same, not_colder)
		fmt.Println("\nIntersection:")
		binrels.PrintWithSource(rels, res)
	}
	if *runAll || *union {
		res := binrels.Union(is_same, not_colder)
		fmt.Println("\nUnion:")
		binrels.PrintWithSource(rels, res)
	}
	if *runAll || *diff {
		res := binrels.Diff(is_same, not_colder)
		fmt.Println("\nDifference (is_same - not_colder):")
		binrels.PrintWithSource(rels, res)
	}
	if *runAll || *symmDiff {
		res := binrels.SymmDiff(is_same, not_colder)
		fmt.Println("\nSymmetric Difference:")
		binrels.PrintWithSource(rels, res)
	}
	if *runAll || *composition {
		res := binrels.Composition(is_same, not_colder)
		fmt.Println("\nComposition (is_same o not_colder):")
		binrels.PrintWithSource(rels, res)
	}
	if *runAll || *reverse != 0 {
		var trgt = getMatrix(*reverse)

		if trgt != nil {
			res := binrels.Transpose(trgt)
			fmt.Println("\nTranspose:")
			binrels.PrintWithSource(rels, res)
		}
	}
	if *runAll || *complement != 0 {
		var trgt = getMatrix(*complement)

		if trgt != nil {
			res := binrels.Complement(trgt)
			fmt.Println("\nComplement:")
			binrels.PrintWithSource(rels, res)
		}
	}
	if *runAll || *top_int >= 0 && *top_int < len(rels) {
		same_top_int := binrels.TopIntersection(is_same, *top_int)
		not_colder_top_int := binrels.TopIntersection(not_colder, *top_int)

		fmt.Printf("\nTop Intersection of is_same with %s:\n", rels[*top_int])
		for _, v := range same_top_int {
			fmt.Println(rels[v])
		}
		fmt.Printf("\nTop Intersection of not_colder with %s:\n", rels[*top_int])
		for _, v := range not_colder_top_int {
			fmt.Println(rels[v])
		}
	}
	if *runAll && *bottom_int >= 0 && *bottom_int < len(rels) {
		same_bottom_int := binrels.BottomIntersection(is_same, *bottom_int)
		not_colder_bottom_int := binrels.BottomIntersection(not_colder, *bottom_int)

		fmt.Printf("\nBottom Intersection of is_same with %s:\n", rels[*bottom_int])
		for _, v := range same_bottom_int {
			fmt.Println(rels[v])
		}
		fmt.Printf("\nBottom Intersection of not_colder with %s:\n", rels[*bottom_int])
		for _, v := range not_colder_bottom_int {
			fmt.Println(rels[v])
		}
	}
	if *runAll && *def_dom > -1 && *def_dom < 2 {
		var trgt = getMatrix(*def_dom)
		if trgt != nil {
			res := binrels.DefinitionDomain(trgt)
			fmt.Println("\nDefinition Domain:")
			if len(res) == 0 {
				fmt.Println("∅")
			}
			for _, v := range res {
				fmt.Println(rels[v])
			}
		}
	}
	if *runAll && *mean_dom >= 0 && *mean_dom < 2 {
		var trgt = getMatrix(*mean_dom)

		if trgt != nil {
			res := binrels.MeaningDomain(trgt)
			fmt.Println("\nMeaning Domain:")
			if len(res) == 0 {
				fmt.Println("∅")
			}
			for _, v := range res {
				fmt.Println(rels[v])
			}
		}
	}
}
