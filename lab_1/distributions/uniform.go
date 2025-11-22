package distributions

func Uniform(a, b float64, x []float64) []float64 {
	if a >= b {
		panic("Uniform: a must be less than b")
	}

	res := make([]float64, len(x))

	for i, v := range x {
		if v < a || v > b {
			res[i] = 0
		} else {
			res[i] = 1 / (b - a)
		}
	}

	return res
}
