package lab6

import "math"

type Solver struct {
}

func K(x, y float64) float64 {
	return x + y
}

func U(x float64) float64 {
	return 1 / math.Sqrt(math.Pow(x, 2)+1)
}

func (s *Solver) Solve() {
	b := 3
	a := 0
	n := 10
	h := (b - a) / n
	x := []float64{}
	for i := 0; i < n; i++ {
		x[i] = float64(a + h*(i-1/2))
	}

	res := [][]float64{}

	for i := 0; i < n; i++ {
		sum := 0.0
		for j := 0; j < n; j++ {
			ad := 0.0
			if i == j {
				ad = (1 - K(x[i], x[j])) * 1 / (U(x[i]))
			} else {
				ad = -float64(h) * K(x[i], x[j])
			}
			res[i][j]
		}
	}
}

func solve() {

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
		// realY := math.Pow(x[i], 2) + math.Exp(x[i])
		realY := math.Exp(x[i])
		diff := math.Abs(resY - realY)
		if diff > err {
			err = diff
		}
		fmt.Printf("x: %f  y: %f  real: %f diff: %f\n", x[i], resY, realY, diff)

	}
	fmt.Println("error :", err)

}
