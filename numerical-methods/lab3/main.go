package main

import (
	"fmt"
	"math"
)

var lastError float64

func simpsonMethod(a float64, b float64, n int) float64 {
	h := (b - a) / float64(n)
	x1 := make([]float64, n)
	x2 := make([]float64, n-1)

	for i := 0; i < n; i++ {
		x1[i] = a + (float64(i)+1)*h - h/2
		if i == n-1 {
			break
		}
		x2[i] = a + (float64(i)+1)*h
	}
	seriesOfSum := 0.0

	for i := 0; i < n; i++ {
		seriesOfSum += 4 * math.Exp(x1[i])
		if i < n-1 {
			seriesOfSum += 2 * math.Exp(x2[i])
		}
	}

	lastError = math.Pow(b-a, 5) / 2880 * math.Abs(math.Exp(b))
	return h / 6 * (math.Exp(a) + math.Exp(b) + seriesOfSum)

}

func rectanglesMethod(a float64, b float64, n int) float64 {
	h := (b - a) / float64(n)
	result := float64(0)
	for i := 0; i < n; i++ {
		result += h * math.Exp(a+h*float64(i)+h/2)
	}
	return result
}

func trapeziaMethod(a float64, b float64, n int) float64 {
	h := (b - a) / float64(n)
	x := make([]float64, n-1)
	for i := 0; i < n-1; i++ {
		x[i] = a + (float64(i)+1)*h
	}
	sum := 0.0
	for i := 0; i < n-1; i++ {
		sum += math.Exp(x[i])
	}
	return h * ((math.Exp(a)+math.Exp(b))/2 + sum)
}

func calcIntegral(eps float64, method func(float64, float64, int) float64, k int) {
	n := 1
	richardson := math.MaxFloat64
	n1 := float64(0)
	i := 0
	for i = 0; math.Abs(richardson) >= eps; i++ {
		n *= 2
		n2 := n1
		n1 = method(0, 1, n)

		richardson = (n1 - n2) / (math.Pow(2, float64(k)) - 1)

	}
	fmt.Println("		iterations : ", i)
	fmt.Println("		result : ", n1)
	fmt.Println("		result with richardson : ", n1+richardson)
}

func main() {
	eps := 0.000001
	fmt.Println("eps : ", eps)
	fmt.Println("rectangles method : ")
	calcIntegral(eps, rectanglesMethod, 2)

	fmt.Println("trapezia method : ")
	calcIntegral(eps, trapeziaMethod, 2)

	fmt.Println("simpsons method : ")
	calcIntegral(eps, simpsonMethod, 2)
	fmt.Println("		error : ", lastError)

}
