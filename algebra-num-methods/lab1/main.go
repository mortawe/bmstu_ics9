package lab1

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"

	"algebra-num-methods/lab1/entities"
)

func readArray(n int, reader *bufio.Reader) []float64 {
	str, _ := reader.ReadString('\n')
	nums := strings.Split(str, " ")
	log.Println(nums)
	result := make([]float64, n)
	i := 0
	for _, num := range nums {
		numF, err := strconv.ParseFloat(num, 64)
		if err == nil {
			result[i] = numF
			i++
		}
	}
	return result
}

// func readMatrix() *entities.Matrix {
//
// }

func testVector() {
	n := 5
	lVec := []float64{1, 2, 3, 4, 5}
	rVec := []float64{5, 4, 3, 2, 1}
	vector := entities.NewVector(n, lVec)
	result := vector.DotProduct(entities.NewVector(n, rVec))
	libResult := floats.Dot(lVec, rVec)

	if result != libResult {
		log.Println(fmt.Sprintf("my: %v, lib: %v", result, libResult))
	}
}

func flat(x [][]float64) []float64 {
	result := []float64{}
	for _, r := range x {
		for _, c := range r {
			result = append(result, c)
		}
	}
	return result
}
func testMatrix() {
	n := 5
	m := 5
	lMat := [][]float64{
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5},
	}
	rMat := [][]float64{
		{5, 4, 3, 2, 1},
		{5, 4, 3, 2, 1},
		{5, 4, 3, 2, 1},
		{5, 4, 3, 2, 1},
		{5, 4, 3, 2, 1},
	}
	myMatrix := entities.NewMatrix(m, n, lMat)
	result := myMatrix.MultiplyMatrix(entities.NewMatrix(m, n, rMat))
	libResult := mat.NewDense(n, m, flat(lMat))
	libResult.Mul(libResult, mat.NewDense(n, m, flat(rMat)))
	fmt.Println("my:", result, "lib:", libResult)
}

func readInt(reader *bufio.Reader) int {
	n, _ := reader.ReadString('\n')
	nInt, _ := strconv.ParseInt(n, 10, 64)
	return int(nInt)
}

func Main() {
	f, err := os.Open("lab1/test/1")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	vecL := readArray(readInt(reader), reader)
	vecR := readArray(readInt(reader), reader)
	fmt.Println(vecL, vecR)
	// testVector()
	// testMatrix()
}
