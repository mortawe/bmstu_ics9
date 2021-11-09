package entities

type Matrix struct {
	m int
	n int
	x [][]float64
}

func NewMatrix(m int, n int, x [][]float64) *Matrix {
	return &Matrix{
		m: m,
		n: n,
		x: x,
	}
}

func (m *Matrix) ColVector(col int) *Vector {
	v := make([]float64, m.m)
	for i := 0; i < m.m; i++ {
		v[i] = m.x[i][col]
	}
	return NewVector(m.m, v)
}

func (m *Matrix) RowVector(row int) *Vector {
	return NewVector(m.n, m.x[row])
}

func (m *Matrix) MultiplyMatrix(right *Matrix) *Matrix {
	result := NewMatrix(m.m, right.n, make([][]float64, m.m))
	for i := 0; i < m.m; i++ {
		result.x[i] = make([]float64, result.n)
		row := m.RowVector(i)
		for j := 0; j < m.n; j++ {
			col := m.ColVector(j)
			result.x[i][j] = row.DotProduct(col)
		}
	}
	return result
}

func (m *Matrix) MultiplyVector(vector *Vector) *Vector {
	v := NewVector(vector.n, make([]float64, vector.n))
	for i := 0; i < v.n; i++ {
		for j := 0; j < m.n; j++ {
			v.x[i] += m.x[i][j]
		}
	}
	return v
}

func (m *Matrix) Transpose() *Matrix {
	result := NewMatrix(m.n, m.m, make([][]float64, m.n))

	for i := 0; i < m.m; i++ {
		for j := 0; j < m.n; j++ {
			result.x[j][i] = m.x[i][j]
		}
	}
	return result
}

func (m *Matrix) Dimension() int {
	return m.m
}
