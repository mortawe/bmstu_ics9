package doublependulum

import (
	"fmt"
	"math"
)

type DoublePendulumSystem struct {
	m1, m2 float64 // mass
	l1, l2 float64 // length
	a1, a2 float64 // angle
	g      float64 // gravity
	w1, w2 float64 // initial velocity
}

func (s *DoublePendulumSystem) StringName() string {
	return fmt.Sprintf("m1=%.2f, m2=%.2f, l1=%.2f, l2=%.2f", s.m1, s.m2, s.l1, s.l2)
}

func NewPendulumSystem(m1, m2, l1, l2, a1, a2, w1, w2, g float64) *DoublePendulumSystem {
	return &DoublePendulumSystem{
		m1: m1,
		m2: m2,
		l1: l1,
		l2: l2,
		a1: a1,
		a2: a2,
		g:  g,
		w1: w1,
		w2: w2,
	}
}

func (s *DoublePendulumSystem) potentialE() float64 {
	y1 := -s.l1 * math.Cos(s.a1)
	y2 := y1 - s.l2*math.Cos(s.a2)
	return s.m1*s.g*y1 + s.m2*s.g*y2
}

func (s *DoublePendulumSystem) kineticE() float64 {
	k1 := 0.5 * s.m1 * math.Pow(s.l1*s.w1, 2)
	k2 := 0.5*s.m2*math.Pow(s.l1*s.w1, 2) +
		math.Pow(s.l2*s.w2, 2) + 2*s.l1*s.l2*s.w1*s.w2*math.Cos(s.a1-s.a2)

	return k1 + k2
}

func (s *DoublePendulumSystem) mechanicalE() float64 {
	return s.kineticE() + s.potentialE()
}

type Kn struct {
	a1, a2, w1, w2 float64
}

func (s *DoublePendulumSystem) LagrangeRhs(k Kn) Kn {
	a1 := (s.l2 / s.l1) * (s.m2 / (s.m1 + s.m2)) * math.Cos(k.a1-k.a2)
	a2 := (s.l1 / s.l2) * math.Cos(k.a1-k.a2)

	f1 := -(s.l2/s.l1)*(s.m2/(s.m1+s.m2))*math.Pow(k.w2, 2)*math.Sin(k.a1-k.a2) - (s.g/s.l1)*math.Sin(k.a1)
	f2 := (s.l1/s.l2)*math.Pow(k.w1, 2)*math.Sin(k.a1-k.a2) - (s.g/s.l2)*math.Sin(k.a2)

	g1 := (f1 - a1*f2) / (1 - a1*a2)
	g2 := (f2 - a2*f1) / (1 - a1*a2)

	return Kn{k.w1, k.w2, g1, g2}
}

func step23(step float64, k, kn Kn) Kn {
	return Kn{
		a1: k.a1 + step*kn.a1/2,
		a2: k.a2 + step*kn.a2/2,
		w1: k.w1 + step*kn.w1/2,
		w2: k.w2 + step*kn.w2/2,
	}
}

func step4(step float64, k, kn Kn) Kn {
	return Kn{
		a1: k.a1 + step*kn.a1,
		a2: k.a2 + step*kn.a2,
		w1: k.w1 + step*kn.w1,
		w2: k.w2 + step*kn.w2,
	}
}

func resultN(step float64, k1, k2, k3, k4 Kn) Kn {
	mul := 1.0 / 6.0 * step
	return Kn{
		a1: mul * (k1.a1 + 2*k2.a1 + 2*k3.a1 + k4.a1),
		a2: mul * (k1.a2 + 2*k2.a2 + 2*k3.a2 + k4.a2),
		w1: mul * (k1.w1 + 2*k2.w1 + 2*k3.w1 + k4.w1),
		w2: mul * (k1.w2 + 2*k2.w2 + 2*k3.w2 + k4.w2),
	}
}

func (s *DoublePendulumSystem) RK4(dt float64) {
	k := Kn{s.a1, s.a2, s.w1, s.w2}
	k1 := s.LagrangeRhs(k)
	k2 := s.LagrangeRhs(step23(dt, k, k1))
	k3 := s.LagrangeRhs(step23(dt, k, k2))
	k4 := s.LagrangeRhs(step4(dt, k, k3))

	r := resultN(dt, k1, k2, k3, k4)

	s.a1 += r.a1
	s.a2 += r.a2
	s.w1 += r.w1
	s.w2 += r.w2
}

func (s *DoublePendulumSystem) Step(dt float64) {
	s.RK4(dt)
}
