package main

import "fmt"

func cubicSpline(x, y []float64, h float64) (a, b, c, d []float64) {
	aS := make([]float64, len(x)-1)
	bS := make([]float64, len(x)-1)
	cS := make([]float64, len(x)-1)
	dS := make([]float64, len(x)-1)
	aS[0] = 0
	for i := 1; i < len(x)-2; i++ {
		aS[i] = 1
		if i == len(x)-3 {
			cS[i] = 0
		} else {
			cS[i] = 1
		}
		b[i] = 4
		d[i] = (3 / (h * h)) * (y[i+2] - 2*y[i+1] + y[i])
	}

	m, _ := NewTridiagonalMatrix(len(x), aS, bS, cS)
	c, _ = m.Count(dS)
	a = make([]float64, len(x))
	b = make([]float64, len(x))
	d = make([]float64, len(x))

	for i := 1; i < len(x); i++ {
		a[i] = y[i-1]
		b[i] = ((y[i] - y[i-1]) / h) - ((h / 3) * (c[i+1] + 2*c[i]))
		d[i] = (c[i+1] - c[i]) / (3 * h)
	}
	return
}

func main() {
	fmt.Println(cubicSpline([]float64{-3, -2, -1, 0, 1, 2, 3, 4}, []float64{5, 4, 3, 2, -4, -2, -1, 6, 5}, 1))
}
