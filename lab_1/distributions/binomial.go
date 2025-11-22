package distributions

import "math"

func factorial(n int) int {
	res := 1

	if n < 0 {
		panic("factorial: negative input not allowed")
	}
	if n > 170 {
		panic("factorial: input too large, would cause overflow")
	}

	if n == 0 || n == 1 {
		return res
	}

	for i := 2; i <= n; i++ {
		res *= i
	}

	return res
}

func combination(n, k int) float64 {
	if k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}

	return float64(factorial(n)) / float64(factorial(k)*factorial(n-k))
}

func Binomial(n int, p float64, x []float64) []float64 {
	if n < 0 {
		panic("Binomial: n must be non-negative")
	}
	if p < 0 || p > 1 {
		panic("Binomial: p must be in [0, 1]")
	}

	res := make([]float64, len(x))
	for i, v := range x {
		if v < 0 || v > float64(n) || math.Abs(v-math.Round(v)) > 1e-9 {
			res[i] = 0
			continue
		}
		res[i] = combination(n, int(v)) * math.Pow(p, v) * math.Pow(1-p, float64(n)-v)
	}
	return res
}
