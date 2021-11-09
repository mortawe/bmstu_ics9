package lab4

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"

	"algebra-num-methods/lab1"
	"algebra-num-methods/lab3"
)

func eucDiffNorm(n int, m1, m2 []float64) float64 {
	norm := 0.0
	for i := 0; i < n; i++ {
		norm += (m1[i] - m2[i]) * (m1[i] - m2[i])
	}
	norm = math.Sqrt(norm)
	return norm
}

func GenMatrix(n int) ([][]float64, []float64) {
	m := make([][]float64, n)
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		m[i] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			m[i][j] = float64(1000 + rand.Int() % 1000)
			m[j][i] = m[i][j]
			// if i == j {
			// 	m[i][i] = math.Abs(m[i][i]) * 100.0 * float64(n+1)
			// }
		}
		b[i] = float64(1000 + rand.Int() % 1000)
	}

	for i := 0; i < n; i++ {
		sum := 0.0
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			sum += math.Abs(m[i][j])
		}
		m[i][i] = sum * math.Abs(rand.Float64()) * 1000
	}
	return m, b
}

func IterMethod1(n int, m [][]float64, b []float64, EPS float64) (int, []float64) {
	err := 1.0
	prev := make([]float64, n)
	copy(prev, b)
	steps := 0
	for ; err >= EPS; steps++ {
		current := make([]float64, n)
		for i := 0; i < n; i++ {
			current[i] = b[i]
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				current[i] -= m[i][j] * prev[j]
			}
			current[i] /= m[i][i]
		}
		err = eucDiffNorm(n, current, prev)
		prev = current
	}
	return steps, prev
}

func IterMethod2(n int, m [][]float64, b []float64, EPS float64) (int, []float64) {
	err := 1.0
	prev := make([]float64, n)
	copy(prev, b)
	steps := 0
	for ; err >= EPS; steps++ {
		current := make([]float64, n)
		copy(current, b)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				current[i] += m[i][j] * prev[j]
			}
		}
		err = eucDiffNorm(n, current, prev)
		copy(prev, current)
	}
	return steps, prev
}

func IterMethod(n int, m [][]float64, b []float64, EPS float64) (int, []float64) {
	err := 1.0
	prev := make([]float64, n)
	copy(prev, b)
	steps := 0
	for ; err >= EPS; steps++ {
		current := make([]float64, n)
		copy(current, b)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				current[i] += m[i][j] * prev[j]
			}
		}
		err = eucDiffNorm(n, current, prev)
		copy(prev, current)
	}
	return steps, prev
}

func Flat(n int, x [][]float64) []float64 {
	result := []float64{}
	for j := 0; j < n; j++ {
		for i := j; i < n; i++ {
			result = append(result, x[j][i])
		}
	}
	return result
}

func AddMatrix(n int, l, r [][]float64) [][]float64 {
	m := make([][]float64, n)
	for i := 0; i < n; i++ {
		m[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			m[i][j] = l[i][j] + r[i][j]
		}
	}
	return m
}

func SubsMatrix(n int, l, r [][]float64) [][]float64 {
	m := make([][]float64, n)
	for i := 0; i < n; i++ {
		m[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			m[i][j] = l[i][j] - r[i][j]
		}
	}
	return m
}

func MultConstMatrix(n int, r [][]float64, c float64) [][]float64 {
	m := make([][]float64, n)
	for i := 0; i < n; i++ {
		m[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			m[i][j] = r[i][j] * c
		}
	}
	return m
}

func GenUnitMatrix(n int) [][]float64 {
	m := make([][]float64, n)
	for i := 0; i < n; i++ {
		m[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			m[i][j] = 1
		}
	}
	return m
}



func Test() {
	// n := 3
	// mGen, b := GenMatrix(n)
	n, mGen, b := lab3.ReadMatrix()
	a := mat.NewSymDense(n, lab1.Flat(mGen))
	det := mat.Det(a)
	if det < 0 {
		log.Fatal("det < 0")
	}
	var eigsym mat.EigenSym
	_ = eigsym.Factorize(a, true)
	eigs := eigsym.Values(nil)
	tL, tR, tOpt := 0.0, 2/floats.Max(eigs), 2/(floats.Max(eigs)+floats.Min(eigs))
	m := SubsMatrix(n,  GenUnitMatrix(n), MultConstMatrix(n, mGen, tOpt))
	fmt.Println(mGen, m)

	fmt.Println(tL, tR, tOpt)
	g := make([]float64, n)
	for i := 0; i < n; i++ {
		g[i] = b[i] * tOpt
	}
	steps1, res1 := IterMethod1(n, mGen, b, 0.0001)
	fmt.Println(steps1, res1)

	steps, res := IterMethod(n, m, g, 0.0001)
	fmt.Println(steps, res)
}
