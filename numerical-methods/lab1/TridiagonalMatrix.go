package main

import (
	"errors"
	"fmt"
	"log"

	"gonum.org/v1/gonum/mat"
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
func (m *TridiagonalMatrix) getMatrix() *mat.Dense {
	matrixF := make([]float64, m.dimension*m.dimension)
	for i := 0; i < m.dimension; i++ {
		for j := 0; j < m.dimension; j++ {
			switch i {
			case j:
				matrixF[i*m.dimension+j] = m.b[j]
			case j + 1:
				matrixF[i*m.dimension+j] = m.c[j]
			case j - 1:
				matrixF[i*m.dimension+j] = m.a[j]
			default:
				matrixF[i*m.dimension+j] = 0
			}
		}
	}
	matrix := mat.NewDense(m.dimension, m.dimension, matrixF)
	return matrix
}

func (m *TridiagonalMatrix) CalcError(x []float64, d []float64) {
	mDense := m.getMatrix()
	xDense := mat.NewDense(m.dimension, 1, x)
	var dCalc mat.Dense
	dCalc.Mul(mDense, xDense)

	dDense := mat.NewDense(m.dimension, 1, d)

	var rDense mat.Dense

	rDense.Sub(&dCalc, dDense)

	var revMatrix mat.Dense
	err := revMatrix.Inverse(mDense)
	if err != nil {
		log.Println(err)
		return
	}
	var resMatrix mat.Dense

	resMatrix.Mul(&revMatrix, &rDense)

	var xCor mat.Dense
	xCor.Sub(xDense, &resMatrix)
	fmt.Println("x corrected:", xCor)
	fmt.Println("error:", resMatrix)
}
