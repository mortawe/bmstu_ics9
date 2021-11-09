package lab3

import (
	"math"
)

func Seidel(n int, m [][]float64, b []float64, eps float64) ([]float64, int) {
	x := make([]float64, n)
	conv := false
	step := 0
	for ; !conv ; step++ {
		if step >= 1000 {
			return nil, -1
		}
		xNew := make([]float64, n)
		copy(xNew, x)
		for i := 0; i < n; i++ {
			s1, s2 := 0.0, 0.0
			for j := 0; j < i; j++ {
				s1 += m[i][j] * xNew[j]
			}
			for j := i + 1; j < n; j++ {
				s2 += m[i][j] * x[j]
			}
			xNew[i] = (b[i] - s1 - s2) / m[i][i]
		}
		sX := 0.0
		for i := 0; i < n; i++ {
			sX += math.Pow(xNew[i]-x[i], 2)
		}
		if math.IsNaN(sX) {
			return nil, -1
		}
		// fmt.Printf(" iter : %d norm : %.10e \n", step, math.Sqrt(sX))

		conv = math.Sqrt(sX) <= eps
		copy(x, xNew)
	}
	return x, step

}

func MethodSeidel(n int, matrixA [][]float64, vectorF []float64, vectorX []float64,  eps float64) ([]float64, int){
	step := 0
	x := make([]float64, n)
	copy(x, vectorX)
	for true {
		if step >= 100 {
			return nil, -1
		}
		vectorXNew := make([]float64, n)
		copy(vectorXNew, x)

		for i := 0; i < n; i++ {
			s1 := 0.0
			s2 := 0.0
			for j := 0; j < i; j++ {
				s1 += matrixA[i][j] * vectorXNew[j]
			}
			for j := i + 1; j < n; j++ {
				s2 += matrixA[i][j] * vectorX[j]
			}
			vectorXNew[i] = (vectorF[i] - s1 - s2) / matrixA[i][i]

		}
		sum := 0.0
		sum1 := 0.0

		for i := 0; i < n; i++ {
			sum += (vectorXNew[i] - vectorX[i]) * (vectorXNew[i] - vectorX[i])
			sum1 += vectorX[i] * vectorX[i]

			vectorX[i] = vectorXNew[i]
		}

		step++
		if sum1 == 0.0 {
			return nil, step
		}
		converge := math.Sqrt(sum) / math.Sqrt(sum1)
		if converge < eps {
			// fmt.Println("Steps for Seidel: ", step)
			break
		}
	}
	return nil, step
}

func Jacobi(n int, m [][]float64, f []float64, x []float64, eps float64) ([]float64, int) {
	norm := 1.0
	steps := 0
	tempX := make([]float64, n)
	for ; norm > eps; steps++ {
		if steps >= 100 {
			return nil, -1
		}
		for i := 0; i < n; i++ {
			tempX[i] = f[i]
			for j := 0; j < n; j++ {
				if i != j {
					tempX[i] -= m[i][j] * x[j]
				}
			}
			tempX[i] /= m[i][i]
		}
		norm = eucDiffNorm(n, x, tempX)
		for i := 0; i < n; i++ {
			x[i] = tempX[i]
		}
		if math.IsNaN(norm) {
			return nil, -1
		}
		// fmt.Printf(" iter : %d norm : %.32e \n", steps, norm)
	}
	return x, steps
}
func MethodJacobi(n int, matrixA [][]float64, vectorF []float64, vectorX []float64, EPS float64) ([]float64, int) {
	tempVectorX := make([]float64, n)
	norm := math.MaxFloat64
	step := 0

	for norm > EPS {
		if step >= 1000 {
			return nil, -1
		}
		step++
		for i := 0; i < n; i++ {
			tempVectorX[i] = vectorF[i]
			for j := 0; j < n; j++ {
				if i != j {
					tempVectorX[i] -= matrixA[i][j] * vectorX[j]
				}
			}
			tempVectorX[i] /= matrixA[i][i]
		}

		// fmt.Println("Step: ", step)
		// fmt.Println("Vector X: ", tempVectorX)

		normForThisStep := 0.0
		normForPreviousStep := 0.0
		for i := 0; i < n; i++ {
			normForThisStep += (vectorX[i] - tempVectorX[i]) * (vectorX[i] - tempVectorX[i])
			normForPreviousStep += vectorX[i] * vectorX[i]
		}

		norm = math.Sqrt(normForThisStep) / math.Sqrt(normForPreviousStep)
		for i := 0; i < n; i++ {
			if math.Abs(vectorX[i]-tempVectorX[i]) > norm {
				norm = math.Abs(vectorX[i] - tempVectorX[i])
			}
			vectorX[i] = tempVectorX[i]
		}
	}
	return tempVectorX, step
}
