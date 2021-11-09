package main

import (
	"fmt"
	"math"
	"os"

	"github.com/wcharczuk/go-chart/v2"
)

var (
	StepsN = 100000
)

func VDP(x, y, m float64) (float64, float64) {
	xDot := y
	yDot := -m*(x*x-1)*y - x
	return xDot, yDot
}

func ChartVDP() {
	ms := []float64{0.01, 0.1, 0.5, 1, 1.5, 2, 3, 3.5, 4}
	chartSeries := []chart.Series{}
	for _, m := range ms {
		xs := []float64{0.0}
		ys := []float64{2.0}
		dt := 0.01

		for i := 0; i < 2000; i++ {
			x, y := VDP(xs[i], ys[i], m)
			xs = append(xs, xs[i]+x*dt)
			ys = append(ys, ys[i]+y*dt)
		}
		chartSeries = append(chartSeries,
			chart.ContinuousSeries{
				Name:    fmt.Sprintf("m=%.3f", m),
				XValues: xs,
				YValues: ys,
			},
		)
	}
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
				return fmt.Sprintf("%.10e", v)
			},
		},
		YAxis: chart.YAxis{
			Name: "y",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%.10e", v)
			},
		},
		Series: chartSeries,
	}
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}
	f, _ := os.Create("vdp.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func ChartDuffing() {
	ms := []float64{0.01, 0.1, 0.5, 1, 1.5, 2, 3, 3.5, 4}
	chartSeries := []chart.Series{}
	for _, m := range ms {
		xs := []float64{0.0}
		ys := []float64{2.0}
		dt := 0.01

		for i := 0; i < 2000; i++ {
			x, y := VDP(xs[i], ys[i], m)
			xs = append(xs, xs[i]+x*dt)
			ys = append(ys, ys[i]+y*dt)
		}
		chartSeries = append(chartSeries,
			chart.ContinuousSeries{
				Name:    fmt.Sprintf("m=%.3f", m),
				XValues: xs,
				YValues: ys,
			},
		)
	}
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
				return fmt.Sprintf("%.10e", v)
			},
		},
		YAxis: chart.YAxis{
			Name: "y",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%.10e", v)
			},
		},
		Series: chartSeries,
	}
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}
	f, _ := os.Create("duffing.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
func main() {
	ChartVDP()

	ChartDuffing()
}
