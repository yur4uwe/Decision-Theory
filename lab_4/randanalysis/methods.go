package randanalysis

import (
	"math"
	"slices"
)

type function func(float64) float64

func ColhomorovMean(data []float64, f, reversedf function) float64 {
	if len(data) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range data {
		sum += f(v)
	}

	return reversedf(sum / float64(len(data)))
}

func Mean(data []float64) float64 {
	f := func(x float64) float64 {
		return x
	}

	reversedf := f

	return ColhomorovMean(data, f, reversedf)
}

func Mode(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	freq := make(map[float64]int)

	for _, v := range data {

		freq[v]++
	}

	mode := data[0]
	maxCount := 0
	for k, v := range freq {
		if v > maxCount {
			maxCount = v
			mode = k
		}
	}

	return mode
}

func Median(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	slices.Sort(data)

	mid := len(data) / 2

	if len(data)%2 == 0 {
		return (data[mid-1] + data[mid]) / 2
	}
	return data[mid]
}

func RMS(data []float64) float64 {
	f := func(x float64) float64 {
		return x * x
	}

	reversedf := func(x float64) float64 {
		return math.Sqrt(x)
	}

	return ColhomorovMean(data, f, reversedf)
}

func GeometricMean(data []float64) float64 {
	f := func(x float64) float64 {
		if x <= 0 {
			return 0
		}
		return math.Log(x)
	}

	reversedf := func(x float64) float64 {
		return math.Exp(x)
	}

	return ColhomorovMean(data, f, reversedf)
}

func HarmonicMean(data []float64) float64 {
	f := func(x float64) float64 {
		if x == 0 {
			return 0
		}
		return 1 / x
	}

	reversedf := f

	return ColhomorovMean(data, f, reversedf)
}
