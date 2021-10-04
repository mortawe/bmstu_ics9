package entities

type Vector struct {
	n int
	x []float64
}

func NewVector(n int, x []float64) *Vector {
	return &Vector{n: n, x: x}
}

func (v *Vector) DotProduct(right *Vector) float64 {
	result := 0.0
	for i := 0; i < v.n; i++ {
		result += v.x[i] * right.x[i]
	}
	return result
}

