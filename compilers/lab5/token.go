package main

import "fmt"

type TokenType string

type Token struct {
	StartPos  int
	EndPos    int
	Line      int
	Value     string
	TokenType TokenType
}

var (
	IDENT      TokenType = `IDENT`
	NUMBERS    TokenType = `NUMBERS`
	OPERATORS  TokenType = `OPERATORS`
	COMMENTARY TokenType = `COMMENTARY`
	KEYWORDS   TokenType = `KEYWORDS`
	WHITESPACE TokenType = `WHITESPACE`
)

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
