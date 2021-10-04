package main

import (
	"flag"
	"log"

	"algebra-num-methods/lab1"
	"algebra-num-methods/lab2"
)

const LAB1 = "lab1"
const LAB2 = "lab2"

func main() {
	labName := flag.String("lab", "", "lab to execute")
	isTest := flag.Int("test", 0, "is test mode")

	flag.Parse()
	if *labName == "" {
		log.Fatal(flag.ErrHelp)
	}

	switch *labName {
	case LAB1:
		lab1.Main()
	case LAB2:
		lab2.Main(*isTest)
	}

}
