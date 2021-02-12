package main

import (
	"errors"
	"fmt"
)

type TridiagonalMatrix struct {
	a []float64 // upper diagonal
	b []float64 // main diagonal
	c []float64 // lower diagonal

	dimension int
}

func NewTridiagonalMatrix(dimension int, a, b, c []float64) (TridiagonalMatrix, error) {
	if len(a) != dimension || len(b) != dimension || len(c) != dimension {
		return TridiagonalMatrix{}, errors.New("wrong input vectors dimension")
	}
	return TridiagonalMatrix{
		a:         a,
		b:         b,
		c:         c,
		dimension: dimension,
	}, nil
}

// Finds Ax = d solutions (A is m tridiagonal matrix)
func (m *TridiagonalMatrix) Count(d []float64) ([]float64, error) {
	if len(d) != m.dimension {
		return nil, errors.New("wrong d vector dimension")
	}
	alpha := make([]float64, m.dimension)
	beta := make([]float64, m.dimension)
	alpha[0] = - m.c[0] / m.b[0]
	beta[0] = d[0] / m.b[0]
	for i := 1; i < m.dimension; i++ {
		alpha[i] = - m.c[i] / (m.b[i] + alpha[i-1]*m.a[i])
		beta[i] = (d[i] - m.a[i]*beta[i-1]) / (m.a[i]*alpha[i-1] + m.b[i])
	}
	x := make([]float64, m.dimension)
	x[m.dimension-1] = beta[m.dimension-1]
	for i := m.dimension - 2; i >= 0; i-- {
		x[i] = alpha[i]*x[i+1] + beta[i]
	}
	return x, nil
}
