package lab2

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"algebra-num-methods/lab2/lib"
	"algebra-num-methods/lab2/method"
)

func eucNorm(n int, m1, m2 []float64) float64 {
	norm := 0.0
	for i := 0; i < n; i++ {
		norm += (m1[i] - m2[i]) * (m1[i] - m2[i])
	}
	norm = math.Sqrt(norm)
	return norm
}

func Test(test int, preparedN int) (float64, float64, float64) {
	rand.Seed(1)
	n := 0
	m := [][]float64{}
	if test == 1 {
		n = preparedN
		m = make([][]float64, n)
		for i := 0; i < n; i++ {
			m[i] = make([]float64, n+1)
			for j := 0; j < n+1; j++ {
				m[i][j] = rand.Float64()
				if j == i {
					// m[i][j] *= float64(rand.Int())
				}

			}
		}
	} else {
		f, err := os.Open("lab2/test/1")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanWords)
		scanner.Scan()
		n64, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		n = int(n64)
		m = make([][]float64, n)
		for i := 0; i < n; i++ {
			m[i] = make([]float64, n+1)
			for j := 0; j < n+1; j++ {
				scanner.Scan()
				m[i][j], _ = strconv.ParseFloat(scanner.Text(), 64)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	}

	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	res := method.GaussMethod(n, m, order)
	resRow := method.GaussMethod(method.GaussLines(n, m, order))
	resCol := method.GaussMethod(method.GaussColumns(n, m, order))
	resColRow := method.GaussMethod(method.GaussColumns(method.GaussLines(n, m, order)))

	if test == 0 {
		// fmt.Println("lib result", libAnswer)
		fmt.Println("row col gauss ", resColRow)
		fmt.Println("gauss ", res)
		fmt.Println("col gauss ", resCol)
		fmt.Println("row gauss ", resRow)
		fmt.Println("is dominant ", method.IsDiagonalDominant(n, m))
		return 0, 0, 0
	}
	if test == 3 {
		for i := 0; i < n; i++ {
			sum := 0.0
			for j := 0; j < n+1; j++ {
				sum += m[i][j]
			}
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					res[i] += m[i][j] *
				}
			}
		}

	}

	if test == 4 {
		libA := make([][]float64, n)
		libB := make([]float64, n)

		for i := 0; i < n; i++ {
			libA[i] = make([]float64, n)
			for j := 0; j < n; j++ {
				libA[i][j] = m[i][j]
			}
			libB[i] = m[i][n]
		}
		libAnswer, _ := lib.GaussPartial(libA, libB)
		fmt.Println(eucNorm(n, resColRow, libAnswer))

	}

	norm := eucNorm(n, resColRow, res)
	colNorm := eucNorm(n, resColRow, resCol)
	rowNorm := eucNorm(n, resColRow, resRow)
	eucNorm := 0.0
	for i := 0; i < n; i++ {
		eucNorm += resColRow[i] * resColRow[i]
	}
	eucNorm = math.Sqrt(eucNorm)

	return norm * 100 / eucNorm, colNorm * 100 / eucNorm, rowNorm * 100 / eucNorm
}
