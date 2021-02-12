package main

import (
	"fmt"
	"log"
)

func main() {
	matrix, err := NewTridiagonalMatrix(4, []float64{0, 1, 1, 1}, []float64{4, 4, 4, 4}, []float64{1, 1, 1, 0})

	if err != nil {
		log.Fatal(err)
	}
	solution, err := matrix.Count([]float64{5, 6, 6, 5})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(solution)
}
