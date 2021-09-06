package main

import (
	"fmt"
	"log"
)

func main() {
	matrix, err := NewTridiagonalMatrix(4, []float64{0, 1, 1, 1}, []float64{3,3,3,3}, []float64{1, 1, 1, 0})

	if err != nil {
		log.Fatal(err)
	}
	solution, err := matrix.Count([]float64{4,5,5,4})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(solution)

	matrix.CalcError(solution, []float64{4,5,5,4})
}
