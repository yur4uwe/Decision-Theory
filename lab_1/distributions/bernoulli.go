package distributions

import "math"

const (
	eps = 1e-9
)

func Bernoulli(p float64, x []float64) []float64 {
	if p < 0 || p > 1 {
		panic("Bernoulli: p must be in [0, 1]")
	}

	q := 1 - p

	res := make([]float64, len(x))

	for i, v := range x {
		if math.Abs(v) < eps || math.Abs(v-1) < eps {
			res[i] = math.Pow(p, v) * math.Pow(q, 1-v)
		} else {
			res[i] = 0
		}
	}

	return res
}
