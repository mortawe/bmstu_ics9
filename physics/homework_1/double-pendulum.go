package main

import (
	"math"
)

type DoublePendulumSystem struct {
	m1, m2 float64 // mass
	l1, l2 float64 // length
	t1, t2 float64 // angle
	g      float64 // gravity
	w1, w2 float64 // initial velocity
}

func NewPendulumSystem(m1, m2, l1, l2, t1, t2, w1, w2, g float64) *DoublePendulumSystem {
	return &DoublePendulumSystem{
		m1: m1,
		m2: m2,
		l1: l1,
		l2: l2,
		t1: t1,
		t2: t2,
		g:  g,
		w1: w1,
		w2: w2,
	}
}

func (s *DoublePendulumSystem) potentialE() float64 {
	y1 := -s.l1 * math.Cos(s.t1)
	y2 := y1 - s.l2*math.Cos(s.t2)
	return s.m1*s.g*y1 + s.m2*s.g*y2
}

func (s *DoublePendulumSystem) kineticE() float64 {
	k1 := 0.5 * s.m1 * math.Pow(s.l1*s.w1, 2)
	k2 := 0.5*s.m2*math.Pow(s.l1*s.w1, 2) +
		math.Pow(s.l2*s.w2, 2) + 2*s.l1*s.l2*s.w1*s.w2*math.Cos(s.t1-s.t2)

	return k1 + k2
}

func (s *DoublePendulumSystem) mechanicalE() float64 {
	return s.kineticE() + s.potentialE()
}

type Kn struct {
	t1, t2, w1, w2 float64
}

func (s *DoublePendulumSystem) LagrangeRhs(k Kn) Kn {
	a1 := (s.l2 / s.l1) * (s.m2 / (s.m1 + s.m2)) * math.Cos(k.t1-k.t2)
	a2 := (s.l1 / s.l2) * math.Cos(k.t1-k.t2)

	f1 := -(s.l2/s.l1)*(s.m2/(s.m1+s.m2))*math.Pow(k.w2, 2)*math.Sin(k.t1-k.t2) - (s.g/s.l1)*math.Sin(k.t1)
	f2 := (s.l1/s.l2)*math.Pow(k.w1, 2)*math.Sin(k.t1-k.t2) - (s.g/s.l2)*math.Sin(k.t2)

	g1 := (f1 - a1*f2) / (1 - a1*a2)
	g2 := (f2 - a2*f1) / (1 - a1*a2)

	return Kn{k.w1, k.w2, g1, g2}
}

func step123(step float64, k, kn Kn) Kn {
	return Kn{
		t1: k.t1 + step*kn.t1/2,
		t2: k.t2 + step*kn.t2/2,
		w1: k.w1 + step*kn.w1/2,
		w2: k.w2 + step*kn.w2/2,
	}
}

func step4(step float64, k, kn Kn) Kn {
	return Kn{
		t1: k.t1 + step*kn.t1,
		t2: k.t2 + step*kn.t2,
		w1: k.w1 + step*kn.w1,
		w2: k.w2 + step*kn.w2,
	}
}

func resultN(step float64, k1, k2, k3, k4 Kn) Kn {
	mul := 1.0 / 6.0 * step
	return Kn{
		t1: mul * (k1.t1 + 2*k2.t1 + 2*k3.t1 + k4.t1),
		t2: mul * (k1.t2 + 2*k2.t2 + 2*k3.t2 + k4.t2),
		w1: mul * (k1.w1 + 2*k2.w1 + 2*k3.w1 + k4.w1),
		w2: mul * (k1.w2 + 2*k2.w2 + 2*k3.w2 + k4.w2),
	}
}

func (s *DoublePendulumSystem) RK4(dt float64)  {
	k := Kn{s.t1, s.t2, s.w1, s.w2}
	k1 := s.LagrangeRhs(k)
	k2 := s.LagrangeRhs(step123(dt, k, k1))
	k3 := s.LagrangeRhs(step123(dt, k, k2))
	k4 := s.LagrangeRhs(step4(dt, k, k3))

	r := resultN(dt, k1, k2, k3, k4)

	s.t1 += r.t1
	s.t2 += r.t2
	s.w1 += r.w1
	s.w2 += r.w2
}

func (s *DoublePendulumSystem) Step(dt float64) {
	s.RK4(dt)
}