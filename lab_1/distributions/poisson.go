package distributions

import "math"

func Poisson(lambda float64, x []float64) []float64 {
	if lambda <= 0 {
		panic("Poisson: lambda must be positive")
	}

	res := make([]float64, len(x))
	for i, v := range x {
		if v < 0 || math.Abs(v-math.Round(v)) > 1e-9 {
			res[i] = 0
			continue
		}
		res[i] = math.Pow(lambda, v) / float64(factorial(int(v))) * math.Exp(-lambda)
	}
	return res
}
