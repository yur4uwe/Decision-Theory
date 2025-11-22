package distributions

import "math"

func Gamma(z float64) float64 {
	// for g=7, n=9
	p := [9]float64{
		0.99999999999980993,
		676.5203681218851,
		-1259.1392167224028,
		771.32342877765313,
		-176.61502916214059,
		12.507343278686905,
		-0.13857109526572012,
		9.9843695780195716e-6,
		1.5056327351493116e-7,
	}
	if z < 0.5 {
		return math.Pi / (math.Sin(math.Pi*z) * Gamma(1-z))
	}
	z -= 1
	x := p[0]
	for i := 1; i < len(p); i++ {
		x += p[i] / (z + float64(i))
	}
	t := z + 7 + 0.5
	return math.Sqrt(2*math.Pi) * math.Pow(t, z+0.5) * math.Exp(-t) * x
}

func Students(alpha float64, x []float64) []float64 {
	if alpha <= 0 {
		panic("Student's t: degrees of freedom (nu) must be positive")
	}

	res := make([]float64, len(x))
	for i, v := range x {
		res[i] = (Gamma((alpha+1)/2) / (math.Sqrt(alpha*math.Pi) * Gamma(alpha/2))) * math.Pow(1+v*v/alpha, -(alpha+1)/2)
	}
	return res
}
