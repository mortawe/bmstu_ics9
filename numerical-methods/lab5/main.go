package main

import (
	"fmt"
	"math"
)

var (
	a, b = 0.0, 1.0
	A, B = 1.0, 1.0 + math.E
	// A, B = 1.0,  math.E
	n    = 10
	h    = (b - a) / float64(n)
	y    = [][]float64{{A, A + h}, {0, h}}
	x    = []float64{}
)

// 2(1-x) = y'' - y'
// a = 0; b = 1+e
// y(0) = 1 y'(0) = 1

// solution : x ^ 2 + e^x
// e^x - решение
func f(x float64) float64 {
	return 2 * (1 - x)
	// return 3 * math.Exp(x)
}

func p(x float64) float64 {
	return -1
	// return 1
}

func q(x float64) float64 {
	return 0
	// return 1
}

func getC1() float64 {
	return (B - y[0][n]) / y[1][n]
}

func solveYi(i int) float64 {
	return y[0][i] + getC1()*y[1][i]
}

func main() {

	for i := 0; i < n; i++ {
		x = append(x, a+h*float64(i))
	}

	for i := 1; i < n; i++ {
		y[0] = append(y[0],
			(math.Pow(h, 2)*f(x[i])-(1.0-(h/2)*p(x[i]))*y[0][i-1]-
				(math.Pow(h, 2)*q(x[i])-2)*y[0][i])/
				(1+h/2*p(x[i])))
		y[1] = append(y[1], (-(1-h/2*p(x[i]))*y[1][i-1]-(math.Pow(h, 2)*q(x[i])-2)*y[1][i])/
			(1+h/2*p(x[i])))
	}
	err := 0.0
	for i := 0; i < n; i++ {
		resY := solveYi(i)
		realY := math.Pow(x[i], 2) + math.Exp(x[i])
		// realY := math.Exp(x[i])
		diff := math.Abs(resY - realY)
		if diff > err {
			err = diff
		}
		fmt.Printf("x: %f  y: %f  real: %f diff: %f\n", x[i], resY, realY, diff)

	}
	fmt.Println("error :", err)

}
/* низкая точность результата обусловлена тем, что в функциональном языке го значение функции экспоненты вычисляется через ряды, в результате чего вычисление происходит не точно */