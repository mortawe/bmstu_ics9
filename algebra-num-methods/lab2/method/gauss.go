package method

import (
	"math"
)

func copySlice(n int, src [][]float64) [][]float64 {
	res := make([][]float64, n)
	for i := 0; i < n; i++ {
		res[i] = make([]float64, n+1)
		for j := 0; j < n+1; j++ {
			res[i][j] = src[i][j]
		}
	}
	return res
}

func GaussMethod(n int, m1 [][]float64, order1 []int) []float64 {
	m := copySlice(n, m1)
	order := make([]int, n)
	copy(order, order1)
	result := make([][]float64, n)

	for i := 0; i < n; i++ {
		result[i] = make([]float64, n+1)
		for j := 0; j < n+1; j++ {
			result[i][j] = m[i][j]
		}
	}

	for k := 0; k < n; k++ {
		for i := 0; i < n+1; i++ {
			result[k][i] /= m[k][k]
		}
		for i := k + 1; i < n; i++ {
			temp := result[i][k] / result[k][k]
			for j := 0; j < n+1; j++ {
				result[i][j] -= result[k][j] * temp
			}
		}
		for i := 0; i < n; i++ {
			for j := 0; j < n+1; j++ {
				m[i][j] = result[i][j]
			}
		}
		for k := n - 1; k > -1; k-- {
			for i := n; i > -1; i-- {
				result[k][i] /= m[k][k]
			}
			for i := k - 1; i > -1; i-- {
				temp := result[i][k] / result[k][k]
				for j := n; j > -1; j-- {
					result[i][j] -= result[k][j] * temp
				}
			}
		}
	}
	answer := make([]float64, n)
	for i := 0; i < n; i++ {
		answer[i] = result[i][n]
	}
	normAns := make([]float64, n)
	for i := 0; i < n; i++ {
		normAns[order[i]] = answer[i]
	}
	return answer
}
func swapCol(i, j int, m [][]float64) [][]float64 {
	temp := m[i]
	m[i] = m[j]
	m[j] = temp
	return m
}

func swapRow(i, j int, m [][]float64) [][]float64 {
	temp := m[i]
	m[i] = m[j]
	m[j] = temp
	return m
}

func GaussColumns(n int, m1 [][]float64, order1 []int) (int, [][]float64, []int) {
	m := copySlice(n, m1)
	order := make([]int, n)
	copy(order, order1)
	for k := 0; k < n; k++ {
		maxEl := math.Abs(m[k][k])
		maxCol := k
		for i := k + 1; i < n; i++ {
			if math.Abs(m[k][i]) > maxEl {
				maxEl = math.Abs(m[k][i])
				maxCol = i
			}
		}
		m = swapCol(maxCol, k, m)
		temp := order[maxCol]
		order[maxCol] = order[k]
		order[k] = temp
	}
	return n, m, order
}

func GaussLines(n int, m1 [][]float64, order1 []int) (int, [][]float64, []int) {
	m := copySlice(n, m1)
	order := make([]int, n)
	copy(order, order1)

	for k := 0; k < n; k++ {
		maxEl := math.Abs(m[k][k])
		maxRow := k
		for i := k + 1; i < n; i++ {
			if math.Abs(m[i][k]) > maxEl {
				maxEl = math.Abs(m[i][k])
				maxRow = i
			}
		}
		m = swapRow(k, maxRow, m)
	}
	return n, m, order
}

func IsDiagonalDominant(n int, m [][]float64) bool {
	isStrict := false
	for i := 0; i < n; i++ {
		diag := m[i][i]
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			if math.Abs(diag) > math.Abs(m[i][j]) {
				isStrict = true
			}
			if math.Abs(diag) <= math.Abs(m[i][j]) {
				return false
			}
		}
	}
	if isStrict {
		return true
	}
	return false
}
