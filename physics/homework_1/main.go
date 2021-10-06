package main

import (
	"flag"
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var (
	m1 = flag.Float64("m1", 50.0, "top pendulum mass")
	m2 = flag.Float64("m2", 5.0, "bottom pendulum mass")
	l1 = flag.Float64("l1", 10.0, "top pod length")
	l2 = flag.Float64("l2", 10.0, "bottom pod length")
	a1 = flag.Float64("a1", 1.0, "top angle")
	a2 = flag.Float64("a2", 1.0, "bottom angle")
	w1 = flag.Float64("w1", 0.0, "init top velocity")
	w2 = flag.Float64("w2", 0.0, "init bottom velocity")
	dt = flag.Float64("dt", 0.1, "time step")
	g  = flag.Float64("g", 10, "gravity")
)

var (
	darkblue = color.RGBA{
		R: 38,
		G: 70,
		B: 83,
		A: 255,
	}
	green = color.RGBA{
		R: 42,
		G: 157,
		B: 143,
		A: 255,
	}
	yellow = color.RGBA{
		R: 233,
		G: 196,
		B: 106,
		A: 255,
	}
	orange = color.RGBA{
		R: 244,
		G: 162,
		B: 97,
		A: 255,
	}
	red = color.RGBA{
		R: 231,
		G: 111,
		B: 81,
		A: 255,
	}
	)


var (
	Nx = 500.0
	Ny = 500.0
)

func run() {
	flag.Parse()
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, Nx, Ny),
		VSync:       true,
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)
	win.SetMatrix(pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, -1)).Moved(pixel.V(Nx/2, Ny/2)))
	s := NewPendulumSystem(*m1, *m2, *l1, *l2, *a1, *a2, *w1, *w2, *g)

	R1 := 12 * (*m1 / (*m1 + *m2))
	R2 := 12 * (*m2 / (*m1 + *m2))

	P1 := 100 * (*l1 / (*l1 + *l2))
	P2 := 100 * (*l2 / (*l1 + *l2))

	trajectoryA := []pixel.Vec{}
	trajectoryB := []pixel.Vec{}

	for !win.Closed() {
		s.Step(*dt)
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))
		win.Clear(color.White)

		if len(trajectoryA) >= 100 {
			trajectoryA = trajectoryA[len(trajectoryA)-100:]
			trajectoryB = trajectoryB[len(trajectoryB)-100:]
		}

		a := pixel.Vec{
			X: P1 * math.Sin(s.a1),
			Y: P1 * math.Cos(s.a1),
		}
		trajectoryA = append(trajectoryA, a)
		b := a.Add(pixel.Vec{
			X: P2 * math.Sin(s.a2),
			Y: P2 * math.Cos(s.a2),
		})
		trajectoryB = append(trajectoryB, b)

		imd := imdraw.New(nil)

		imd.Color = darkblue
		imd.Push(pixel.ZV, a, b)
		imd.Line(3)

		imd.Color = green
		imd.Push(trajectoryA...)
		imd.Line(1)
		imd.Color = red
		imd.Push(trajectoryB...)
		imd.Line(1)

		imd.Color = green
		imd.Push(a)
		imd.Circle(R1, 0)

		imd.Color = red
		imd.Push(b)
		imd.Circle(R2, 0)

		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
