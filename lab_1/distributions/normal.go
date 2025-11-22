package distributions

import "math"

func Normal(mean, variance float64, x []float64) []float64 {
	if variance <= 0 {
		panic("Normal: variance must be positive")
	}

	deviation := math.Sqrt(variance)

	res := make([]float64, len(x))
	for i, v := range x {
		res[i] = (1 / (deviation * math.Sqrt(2*math.Pi))) * math.Exp(-0.5*math.Pow((v-mean)/deviation, 2))
	}

	return res
}
