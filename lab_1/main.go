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
	bernoulli_x := graph.IntLinearArray(0, 1)
	bernoulli_y := distributions.Bernoulli(p, bernoulli_x)

	bline := graph.NewLS()
	bline.Pillars()

	g := graph.NewGraph()
	g.SetLineStyle(bline)
	g.Plot(bernoulli_x, bernoulli_y)

	if err := g.Draw(); err != nil {
		panic(err)
	}
	if err := g.SavePNG("images/bernoulli.png", true); err != nil {
		panic(err)
	}

	expected_mean := p
	expected_variance := p * (1 - p)
	fmt.Println("Bernoulli Distribution (p =", p, ")")
	Print(bernoulli_x, bernoulli_y, expected_mean, expected_variance, false)
}

func plotBinomial(n int, p float64) {
	binline := graph.NewLS()
	binline.Dots()
	binline.Solid()

	g := graph.NewGraph()
	g.SetLineStyle(binline)

	binomial_x := graph.IntLinearArray(0, n)
	binomial_y := distributions.Binomial(n, p, binomial_x)
	g.Plot(binomial_x, binomial_y)

	if err := g.Draw(); err != nil {
		panic(err)
	}
	if err := g.SavePNG("images/binomial.png", true); err != nil {
		panic(err)
	}

	expected_mean := float64(n) * p
	expected_variance := float64(n) * p * (1 - p)
	fmt.Println("Binomial Distribution (n =", n, ", p =", p, ")")
	Print(binomial_x, binomial_y, expected_mean, expected_variance, false)
}

func plotPoisson(lambda float64) {
	pline := graph.NewLS()
	pline.Dots()
	pline.Solid()

	g := graph.NewGraph()
	g.SetLineStyle(pline)

	kmax := int(N + 5)
	poisson_x := graph.IntLinearArray(0, kmax)
	poisson_y := distributions.Poisson(lambda, poisson_x)
	g.Plot(poisson_x, poisson_y)

	if err := g.Draw(); err != nil {
		panic(err)
	}
	if err := g.SavePNG("images/poisson.png", true); err != nil {
		panic(err)
	}

	expected_mean := lambda
	expected_variance := lambda
	fmt.Printf("Poisson Distribution (λ = %.2f)\n", lambda)
	Print(poisson_x, poisson_y, expected_mean, expected_variance, false)
}

func plotUniform(a, b float64) {
	uniform_x := graph.LinearArray(a-3, b+3, inlen)
	uniform_y := distributions.Uniform(a, b, uniform_x)

	uline := graph.NewLS()
	uline.Solid()

	g := graph.NewGraph()
	g.SetLineStyle(uline)
	g.Plot(uniform_x, uniform_y)

	if err := g.Draw(); err != nil {
		panic(err)
	}
	if err := g.SavePNG("images/uniform.png", true); err != nil {
		panic(err)
	}

	expected_mean := (a + b) / 2
	expected_variance := (b - a) * (b - a) / 12
	fmt.Println("Uniform Distribution (a =", a, ", b =", b, ")")
	Print(uniform_x, uniform_y, expected_mean, expected_variance, true)
}

func plotNormal(mu, sigma2 float64) {
	normal_x := graph.LinearArray(mu-4*sigma2, mu+4*sigma2, inlen)
	normal_y := distributions.Normal(mu, sigma2, normal_x)

	nline := graph.NewLS()
	nline.Solid()

	g := graph.NewGraph()
	g.SetLineStyle(nline)
	g.Plot(normal_x, normal_y)

	if err := g.Draw(); err != nil {
		panic(err)
	}
	if err := g.SavePNG("images/normal.png", true); err != nil {
		panic(err)
	}

	expected_mean := mu
	expected_variance := sigma2
	fmt.Println("Normal Distribution (mean =", mu, ", variance =", sigma2, ")")
	Print(normal_x, normal_y, expected_mean, expected_variance, true)
}

func plotPareto(x0, alpha float64) {
	pareto_x := graph.LinearArray(0, 20, inlen)
	pareto_y := distributions.Pareto(x0, alpha, pareto_x)

	parline := graph.NewLS()
	parline.Solid()

	g := graph.NewGraph()
	g.SetLineStyle(parline)
	g.Plot(pareto_x, pareto_y)

	if err := g.Draw(); err != nil {
		panic(err)
	}
	if err := g.SavePNG("images/pareto.png", true); err != nil {
		panic(err)
	}

	expected_mean := alpha * x0 / (alpha - 1)
	expected_variance := alpha * x0 * x0 / ((alpha - 1) * (alpha - 1) * (alpha - 2))
	fmt.Println("Pareto Distribution (x0 =", x0, ", α =", alpha, ")")
	Print(pareto_x, pareto_y, expected_mean, expected_variance, true)
}

func plotStudents(nu float64) {
	students_x := graph.LinearArray(-5, 5, inlen)
	students_y := distributions.Students(nu, students_x)

	stline := graph.NewLS()
	stline.Solid()

	g := graph.NewGraph()
	g.SetLineStyle(stline)
	g.Plot(students_x, students_y)

	if err := g.Draw(); err != nil {
		panic(err)
	}
	if err := g.SavePNG("images/students.png", true); err != nil {
		panic(err)
	}

	expected_mean := 0.0
	expected_variance := nu / (nu - 2)
	fmt.Println("Student's t Distribution (ν =", nu, ")")
	Print(students_x, students_y, expected_mean, expected_variance, true)
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
