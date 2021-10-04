package lab2

import (
	"fmt"
	"os"
	"sync"

	chart "github.com/wcharczuk/go-chart/v2"
)

func Main(test int) {
	switch test {
	case 0:
		Test(test, 0)
		return
	case 1:
		Test(test, 100)
	}
	xVal1 := []float64{}
	yVal1 := []float64{}
	yVal2 := []float64{}
	yVal3 := []float64{}
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i := 10; i < 200; i += 10 {
		wg.Add(1)
		go func(i  int) {
			res1, res2, res3 := Test(test, i)
			mu.Lock()
			xVal1 = append(xVal1, float64(i))
			yVal1 = append(yVal1, res1)
			yVal2 = append(yVal2, res2)
			yVal3 = append(yVal3, res3)
			mu.Unlock()
			wg.Done()
		}(i)
	}

	wg.Wait()
	graph1 := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%f", v)
			},
		},
		YAxis: chart.YAxis{
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%e", v)
			},
		},

		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "Gauss",
				XValues: xVal1,
				YValues: yVal1,
			},
			chart.ContinuousSeries{
				Name: "Col Row Gauss",
				XValues: xVal1,
				YValues: yVal2,
			},
			chart.ContinuousSeries{
				Name: "Row Gauss",

				XValues: xVal1,
				YValues: yVal3,
			},
		},
	}
	graph1.Elements = []chart.Renderable{
		chart.Legend(&graph1),
	}
	f, _ := os.Create("output.png")
	defer f.Close()
	graph1.Render(chart.PNG, f)
}
