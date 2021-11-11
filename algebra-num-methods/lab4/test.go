package lab4

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/wcharczuk/go-chart/v2"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"

	"algebra-num-methods/lab1"
	"algebra-num-methods/lab2/lib"
)

func eucDiffNorm(n int, m1, m2 []float64) float64 {
	norm := 0.0
	for i := 0; i < n; i++ {
		norm += (m1[i] - m2[i]) * (m1[i] - m2[i])
	}
	norm = math.Sqrt(norm)
	return norm
}

func eucNorm(n int, m2 []float64) float64 {
	norm := 0.0
	for i := 0; i < n; i++ {
		norm += m2[i] * m2[i]
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
			m[i][j] = float64(rand.Int() % 100)  / 100000
			m[j][i] = m[i][j]
			// if i == j {
			// 	m[i][i] = math.Abs(m[i][i]) * 100.0 * float64(n+1)
			// }
		}
		b[i] = rand.Float64()
	}

	for i := 0; i < n; i++ {
		sum := 0.0
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			sum += math.Abs(m[i][j])
		}
		m[i][i] = sum * 1.01
		// m[i][i] =  float64(rand.Int() % 100)  / 100000
	}
	return m, b
}

func IterMethod1(n int, m [][]float64, b []float64, EPS float64) (int, []float64) {
	err := 10.0
	prev := make([]float64, n)
	copy(prev, b)
	steps := 0
	// g := SubsMatrix(n, GenUnitMatrix(n), m)
	for ; ; steps++ {
		// fmt.Println(prev)
		current := make([]float64, n)
		for i := 0; i < n; i++ {
			current[i] = prev[i] + b[i]
			for j := 0; j < n; j++ {
				current[i] += -m[i][j] * prev[j]
			}
		}
		err = eucDiffNorm(n, current, prev)
		if math.Abs(err)*mt <= EPS {
			return steps, prev
		}
		if math.IsNaN(err) {
			return -1, prev
		}
		prev = current
	}
	return steps, prev
}
func IterMethod2(n int, m [][]float64, b []float64, EPS float64) (int, []float64) {
	err := 10.0
	prev := make([]float64, n)
	copy(prev, b)
	steps := 0
	// g := SubsMatrix(n, GenUnitMatrix(n), m)
	for ; ; steps++ {
		// fmt.Println(prev)
		current := make([]float64, n)
		for i := 0; i < n; i++ {
			current[i] = b[i]
			for j := 0; j < n; j++ {
				current[i] += m[i][j] * prev[j]
			}
		}
		err = eucDiffNorm(n, current, prev)
		if math.Abs(err)*mt <= EPS {
			return steps, prev
		}
		if math.IsNaN(err) {
			return -1, prev
		}
		prev = current
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

var mt = 0.0

func GenUnitMatrix(n int) [][]float64 {
	m := make([][]float64, n)
	for i := 0; i < n; i++ {
		m[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if i == j {
				m[i][j] = 1
			} else {
				m[i][j] = 0
			}

		}
	}
	return m
}
func Test() {
	n := 5
	mGen, b := GenMatrix(n)
	// n, mGen, b := lab3.ReadMatrix("lab4/test/1")
	a := mat.NewSymDense(n, lab1.Flat(mGen))
	fmt.Println(mGen)
	var eigsym mat.EigenSym
	_ = eigsym.Factorize(a, true)
	eigs := eigsym.Values(nil)
	tL, tR, tOpt := 0.0, 2/floats.Max(eigs), 2/(floats.Max(eigs)+floats.Min(eigs))

	fmt.Println(tL, tR, tOpt)
	step := (tR - tL) / 100
	xs, ys, zs := []float64{}, []float64{}, []float64{}
	eps := 0.001
	res, _ := lib.GaussPartial(mGen, b)
	fmt.Println(res)

	// _, res1 := IterMethod(n, mGen, b, eps)

	g := make([]float64, n)
	for j := 0; j < n; j++ {
		g[j] = b[j] * tOpt
	}
	// _, res1 := IterMethod(n,  MultConstMatrix(n, mGen, tOpt), g, eps)
	// fmt.Println(res, res1)
	for i := tL + step; i <= tR + step; i += step {
	// for i := tOpt - step; i <= tR; i += step {
		t := i
		if tOpt-i < step {
			t = tOpt
		}
		m := SubsMatrix(n, GenUnitMatrix(n), MultConstMatrix(n, mGen, t))
		g := make([]float64, n)
		for j := 0; j < n; j++ {
			g[j] = b[j] * t
		}
		mDet := mat.NewSymDense(n, lab1.Flat(m))
		mt = mat.Norm(mDet, 1)
		fmt.Println(mt)
		if mt > 1 {
			// panic("norm > 1")
		}
		mt = math.Abs(mt / (1 - mt))
		steps, res1 := IterMethod2(n, m, g, eps)
		if steps == -1 {
			continue
		}
		xs = append(xs, float64(steps))
		ys = append(ys, t)
		zs = append(zs, eucDiffNorm(n, res, res1))
	}
	fmt.Println(len(xs), floats.Max(xs), floats.Min(xs))
	fmt.Println(len(ys), floats.Max(ys), floats.Min(ys))
	fmt.Println(ys[floats.MinIdx(zs)], floats.Max(ys))

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
			chart.ContinuousSeries{
				Name:    "Params",
				XValues: ys,
				YValues: xs,
			},
			chart.AnnotationSeries{
				Annotations: []chart.Value2{
					{XValue: tOpt, YValue: xs[0], Label: "Optimal"}},
			}},
	}
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}

	f, _ := os.Create("steps.png")
	defer f.Close()
	graph.Render(chart.PNG, f)

	graphErr := chart.Chart{
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
			Name: "Error",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%v", v)
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "Errors",
				XValues: ys,
				YValues: zs,
			},
			chart.AnnotationSeries{
				Annotations: []chart.Value2{
					{XValue: tOpt, YValue: zs[0], Label: "Optimal"}},
			}},
	}
	graphErr.Elements = []chart.Renderable{
		chart.LegendLeft(&graphErr),
	}

	fErr, _ := os.Create("errors.png")
	defer fErr.Close()
	graphErr.Render(chart.PNG, fErr)

}
