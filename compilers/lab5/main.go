package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	content, err := ioutil.ReadFile("./input.txt")
	lexer := NewAutomataLexer(content)
	errors := []error{}
	commnts := []string{}
	for !lexer.isEOF() {
		token, err := lexer.nextToken()
		if err != nil {
			errors = append(errors, fmt.Errorf("%v : %s", err, token.Value))
			continue
		}
		switch token.TokenType {
		case COMMENTARY:
			commnts = append(commnts, token.Value)
		case WHITESPACE:
		default:
			fmt.Println(token.String())
		}
	}

	fmt.Println("COMMENTARIES:")
	for _, c := range commnts {
		fmt.Println(c)
	}

	fmt.Println("ERRORS:")
	for _, c := range errors {
		fmt.Println(c)
	}
	if err != nil {
		log.Fatal(err)
	}
}
