package main

import (
	"decision-theory/graph"
	"decision-theory/lab_1/distributions"
	"flag"
	"fmt"
	"math"
)

const (
	N     = 12.0
	inlen = 1000
)

func cont_mean(x, y []float64) float64 {
	dx := x[1] - x[0]
	var mean float64
	for i := range x {
		mean += x[i] * y[i] * dx
	}
	return mean
}

func cont_variance(x, y []float64, mean float64) float64 {
	dx := x[1] - x[0]
	var variance float64
	for i := range x {
		variance += (x[i] - mean) * (x[i] - mean) * y[i] * dx
	}
	return variance
}

func disc_mean(x, y []float64) float64 {
	var mean float64
	for i := range x {
		mean += x[i] * y[i]
	}
	return mean
}

func disc_variance(x, y []float64, mean float64) float64 {
	var variance float64
	for i := range x {
		variance += (x[i] - mean) * (x[i] - mean) * y[i]
	}
	return variance
}

func Print(x, y []float64, expected_mean, expected_variance float64, continious bool) {
	var calculated_mean, calculated_variance float64
	if continious {
		calculated_mean = cont_mean(x, y)
		calculated_variance = cont_variance(x, y, calculated_mean)
	} else {
		calculated_mean = disc_mean(x, y)
		calculated_variance = disc_variance(x, y, calculated_mean)
	}
	fmt.Println("Expected mean:", expected_mean, "Calculated mean:", calculated_mean)
	fmt.Println("Expected variance:", expected_variance, "Calculated variance:", calculated_variance)

	expectations_met := false
	if math.Abs(expected_mean-calculated_mean) < 0.1 && math.Abs(expected_variance-calculated_variance) < 0.1 {
		expectations_met = true
	}

	fmt.Println("Expectations met:", expectations_met)
	fmt.Println()
}

func plotBernoulli(p float64) {
	bernoulli_x := graph.IntLinearArray(0, 2)
	bernoulli_y := distributions.Bernoulli(p, bernoulli_x)
	bline := graph.NewLS()
	bline.Pillars(10)

	expected_mean := p
	expected_variance := p * (1 - p)
	fmt.Println("Bernoulli Distribution (p =", p, ")")
	draw(bernoulli_x, bernoulli_y, expected_mean, expected_variance, false, bline, "images/bernoulli.png")
}

func plotBinomial(n int, p float64) {
	binline := graph.NewLS()
	binline.Dots()
	binline.Solid()
	binomial_x := graph.IntLinearArray(0, n+1)
	binomial_y := distributions.Binomial(n, p, binomial_x)

	expected_mean := float64(n) * p
	expected_variance := float64(n) * p * (1 - p)
	fmt.Println("Binomial Distribution (n =", n, ", p =", p, ")")
	draw(binomial_x, binomial_y, expected_mean, expected_variance, false, binline, "images/binomial.png")
}

func plotPoisson(lambda float64) {
	pline := graph.NewLS()
	pline.Dots()
	pline.Solid()
	kmax := int(N + 5)
	poisson_x := graph.IntLinearArray(0, kmax+1)
	poisson_y := distributions.Poisson(lambda, poisson_x)

	expected_mean := lambda
	expected_variance := lambda
	fmt.Printf("Poisson Distribution (λ = %.2f)\n", lambda)
	draw(poisson_x, poisson_y, expected_mean, expected_variance, false, pline, "images/poisson.png")
}

func plotUniform(a, b float64) {
	uniform_x := graph.LinearArray(a-3, b+3, inlen)
	uniform_y := distributions.Uniform(a, b, uniform_x)
	uline := graph.NewLS()
	uline.Solid()

	expected_mean := (a + b) / 2
	expected_variance := (b - a) * (b - a) / 12
	fmt.Println("Uniform Distribution (a =", a, ", b =", b, ")")
	draw(uniform_x, uniform_y, expected_mean, expected_variance, true, uline, "images/uniform.png")
}

func plotNormal(mu, sigma2 float64) {
	normal_x := graph.LinearArray(mu-4*sigma2, mu+4*sigma2, inlen)
	normal_y := distributions.Normal(mu, sigma2, normal_x)
	nline := graph.NewLS()
	nline.Solid()

	expected_mean := mu
	expected_variance := sigma2
	fmt.Println("Normal Distribution (mean =", mu, ", variance =", sigma2, ")")
	draw(normal_x, normal_y, expected_mean, expected_variance, true, nline, "images/normal.png")
}

func plotPareto(x0, alpha float64) {
	pareto_x := graph.LinearArray(0, 20, inlen)
	pareto_y := distributions.Pareto(x0, alpha, pareto_x)
	parline := graph.NewLS()
	parline.Solid()

	expected_mean := alpha * x0 / (alpha - 1)
	expected_variance := alpha * x0 * x0 / ((alpha - 1) * (alpha - 1) * (alpha - 2))
	fmt.Println("Pareto Distribution (x0 =", x0, ", α =", alpha, ")")
	draw(pareto_x, pareto_y, expected_mean, expected_variance, true, parline, "images/pareto.png")
}

func plotStudents(nu float64) {
	students_x := graph.LinearArray(-5, 5, inlen)
	students_y := distributions.Students(nu, students_x)
	stline := graph.NewLS()
	stline.Solid()

	expected_mean := 0.0
	expected_variance := nu / (nu - 2)
	fmt.Println("Student's t Distribution (ν =", nu, ")")
	draw(students_x, students_y, expected_mean, expected_variance, true, stline, "images/students.png")
}

func draw(x, y []float64, expected_mean, expected_variance float64, continious bool, ls *graph.LineStyle, filename string) {
	g := graph.NewGraph(800, 400)
	g.Plot(x, y, ls)

	if err := g.Draw(); err != nil {
		panic(err)
	}
	if err := g.SavePNG(filename, true); err != nil {
		panic(err)
	}

	Print(x, y, expected_mean, expected_variance, continious)
}

func main() {
	bernFlag := flag.Bool("ber", false, "Plot Bernoulli distribution")
	binomFlag := flag.Bool("bin", false, "Plot Binomial distribution")
	poisFlag := flag.Bool("pois", false, "Plot Poisson distribution")
	unifFlag := flag.Bool("unif", false, "Plot Uniform distribution")
	normFlag := flag.Bool("norm", false, "Plot Normal distribution")
	paretoFlag := flag.Bool("par", false, "Plot Pareto distribution")
	studFlag := flag.Bool("stud", false, "Plot Student's t distribution")
	allFlag := flag.Bool("all", false, "Plot all distributions")

	p := flag.Float64("p", 1/(N+1), "Parameter p (probability) for Bernoulli/Binomial")
	n := flag.Int("n", int(N+2), "Parameter n for Binomial")
	lambda := flag.Float64("l", N+10, "Parameter lambda for Poisson")
	alpha := flag.Float64("al", N, "Parameter alpha for Pareto/Student's t")
	x0 := flag.Float64("x0", N, "Parameter x0 for Pareto")
	a := flag.Float64("a", -N, "Parameter a (lower bound) for Uniform")
	b := flag.Float64("b", N, "Parameter b (upper bound) for Uniform")
	mu := flag.Float64("mu", N, "Parameter mu (mean) for Normal")
	sigma2 := flag.Float64("s2", N/2, "Parameter sigma^2 (variance) for Normal")

	flag.Parse()

	if *allFlag || *bernFlag {
		plotBernoulli(*p)
	}
	if *allFlag || *binomFlag {
		plotBinomial(*n, *p)
	}
	if *allFlag || *poisFlag {
		plotPoisson(*lambda)
	}
	if *allFlag || *unifFlag {
		plotUniform(*a, *b)
	}
	if *allFlag || *normFlag {
		plotNormal(*mu, *sigma2)
	}
	if *allFlag || *paretoFlag {
		plotPareto(*x0, *alpha)
	}
	if *allFlag || *studFlag {
		plotStudents(*alpha)
	}
}
