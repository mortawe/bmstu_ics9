package lab3

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/wcharczuk/go-chart/v2"
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
		for j := 0; j < n; j++ {
			m[i][j] = float64(rand.Int() % 1000)
			// if i == j {
			// 	m[i][i] = math.Abs(m[i][i]) * 100.0 * float64(n+1)
			// }
		}
		b[i] = float64(rand.Int() % 1000)
	}

	for i := 0; i < n; i++ {
		sum := 0.0
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			sum += math.Abs(m[i][j])
		}
		m[i][i] = sum * math.Abs(rand.Float64()) * 10
	}
	return m, b
}

func ReadMatrix( path string) (int, [][]float64, []float64,) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	n64, _ := strconv.ParseInt(scanner.Text(), 10, 64)
	n := int(n64)
	m := make([][]float64, n)
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		m[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			scanner.Scan()
			m[i][j], _ = strconv.ParseFloat(scanner.Text(), 64)
		}
		scanner.Scan()

		b[i], _ = strconv.ParseFloat(scanner.Text(), 64)

	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return n, m, b
}

func Test(gen bool, test int, preparedN int) (float64, float64, float64) {
	if true {
		testify()
	} else {
		n := 0

		m := [][]float64{}
		b := []float64{}
		if gen {
			m, b = GenMatrix(preparedN)
			n = preparedN
		} else {
			n, m, b = ReadMatrix("lab3/test/1")
		}

		eps := 0.0001

		_, stepsS := Seidel(n, m, b, eps)
		fmt.Println("Seidel steps : ", stepsS)

		bCopy := make([]float64, n)
		copy(bCopy, b)
		_, stepsJ := MethodJacobi(n, m, b, bCopy, eps)
		fmt.Println("Jacobi steps : ", stepsJ)

	}
	return 0, 0, 0

}

func testify() {
	xs := []float64{}
	js := []float64{}
	ss := []float64{}

	for i := 10; i < 1000; i += 10  {
		m := [][]float64{}
		b := []float64{}
		m, b = GenMatrix(i)
		eps := 0.001

		bCopy := make([]float64, i)
		copy(bCopy, b)
		_, stepsJ := Jacobi(i, m, b, bCopy, eps)
		if stepsJ == -1 {
			continue
		}
		b1Copy := make([]float64, i)
		copy(b1Copy, b)
		_, stepsS := MethodSeidel(i, m, b, b1Copy, eps)
		if stepsS == -1 {
			continue
		}
		ss = append(ss, float64(stepsS))
		xs = append(xs, float64(i))
		js = append(js, float64(stepsJ))
	}
	fmt.Println(len(xs), len(ss), len(js))
	graph := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 100,
			},
		},
		XAxis: chart.XAxis{
			Name: "x",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%v", v)
			},
		},
		YAxis: chart.YAxis{
			Name: "v",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%v", v)
			},
		},
		Series: []chart.Series{chart.ContinuousSeries{
				Name:    "Jacobi",
				XValues: xs,
				YValues: js,
			}, chart.ContinuousSeries{
			Name:    "Seidel",
			XValues: xs,
			YValues: ss,
		}},
	}
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}
	f, _ := os.Create("duffing.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
