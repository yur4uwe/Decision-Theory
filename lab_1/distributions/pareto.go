package distributions

import "math"

func Pareto(x0, alpha float64, x []float64) []float64 {
	if x0 <= 0 {
		panic("Pareto: x0 must be positive")
	}
	if alpha <= 0 {
		panic("Pareto: alpha must be positive")
	}

	if math.Abs(alpha-1) < 1e-9 || math.Abs(alpha-2) < 1e-9 {
		panic("Pareto: alpha cannot be 1 or 2")
	}

	res := make([]float64, len(x))

	for i, v := range x {
		if v < x0 {
			res[i] = 0
		} else {
			res[i] = (alpha * math.Pow(x0, alpha)) / math.Pow(v, alpha+1)
		}
	}

	return res
}
