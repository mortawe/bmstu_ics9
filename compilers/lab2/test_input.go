package main

import (
	"fmt"
)

func main() {
	fmt.Print("rfsddsge")
	go func() {
		InitRoutes()
	}()

	go InitRoutes()
}

func InitRoutes() {
}
