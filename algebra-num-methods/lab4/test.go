package lab4

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/wcharczuk/go-chart/v2"
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
			m[i][j] = float64(1000 + rand.Int()%1000)
			m[j][i] = m[i][j]
			// if i == j {
			// 	m[i][i] = math.Abs(m[i][i]) * 100.0 * float64(n+1)
			// }
		}
		b[i] = float64(1000 + rand.Int()%1000)
	}

	for i := 0; i < n; i++ {
		sum := 0.0
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			sum += math.Abs(m[i][j])
		}
		m[i][i] = sum * math.Abs(rand.Float64()) * 100000
	}
	return m, b
}

func IterMethod1(n int, m [][]float64, b []float64, EPS float64) (int, []float64) {
	err := 1.0
	prev := make([]float64, n)
	copy(prev, b)
	steps := 0
	for ; err >= EPS && steps < 1000; steps++ {
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
		copy(current, prev)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				current[i] += m[i][j] * prev[j]
			}
			current[i] -= b[i]
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
	// n := 30
	// mGen, b := GenMatrix(n)
	n, mGen, b := lab3.ReadMatrix("lab4/test/1")
	a := mat.NewSymDense(n, lab1.Flat(mGen))
	det := mat.Det(a)
	if det < 0 {
		log.Fatal("det < 0")
	}
	var eigsym mat.EigenSym
	_ = eigsym.Factorize(a, true)
	eigs := eigsym.Values(nil)
	tL, tR, tOpt := 0.0, 2/floats.Max(eigs), 2/(floats.Max(eigs)+floats.Min(eigs))
	fmt.Println(tL, tR, tOpt)
	step := (tR - tL) / 100
	xs, ys, zs := []float64{}, []float64{}, []float64{}
	_, res := IterMethod1(n, mGen, b, 0.001)

	for i := tL + step; i <= tR; i += step {
		m := SubsMatrix(n, GenUnitMatrix(n), MultConstMatrix(n, mGen, i))
		g := make([]float64, n)
		for j := 0; j < n; j++ {
			g[j] = -b[j] * i
		}
		// fmt.Println(steps1)

		steps, res1 := IterMethod1(n, m, g, 0.001)
		// fmt.Println(steps)
		xs = append(xs, float64(steps))
		ys = append(ys, i)
		if math.IsNaN(eucDiffNorm(n, res, res1)) {
			zs = append(zs, 0.0)
		} else {
			zs = append(zs, eucDiffNorm(n, res, res1))
		}
	}
	fmt.Println(len(xs), len(ys), len(zs))
	// m := AddMatrix(n, GenUnitMatrix(n), MultConstMatrix(n, mGen, tOpt))
	// g := make([]float64, n)
	// for j := 0; j < n; j++ {
	// 	g[j] = -b[j] * tOpt
	// }
	// steps, _ := IterMethod1(n, m, g, 0.001)

	graph := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 100,
			},
		},
		XAxis: chart.XAxis{
			Name: "Param",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%v", v)
			},
		},
		YAxis: chart.YAxis{
			Name: "Iter Number",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%v", v)
			},
		},
		Series: []chart.Series{
			// chart.ContinuousSeries{
			// Name:    "Params",
			// XValues: ys,
			// YValues: xs,
		// },
			chart.ContinuousSeries{
				Name:    "Errors",
				XValues: ys,
				YValues: zs,
			// },
			// chart.AnnotationSeries{
			// 	Annotations: []chart.Value2{
			// 		{XValue: tOpt, YValue: float64(steps), Label: "Optimal"}},
			}},
	}
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}

	f, _ := os.Create("duffing.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
