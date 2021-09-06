package main

import "fmt"

type Token struct {
	StartPos int
	EndPos   int
	Line     int
	Value    string
	TokenType TokenType
}

func (t *Token) String() string {
	result := fmt.Sprintf("<%s> : (%d, %d)", t.TokenType, t.Line, t.StartPos)
	if t.StartPos != t.EndPos {
		result += fmt.Sprintf(" - (%d, %d)", t.Line, t.EndPos)
	}
	if t.Value != "" {
		result += fmt.Sprintf(" - %s ", t.Value)
	}
	return result
}