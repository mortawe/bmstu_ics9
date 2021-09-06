package main

import "github.com/dlclark/regexp2"

type Lexer struct {
	Code   string
	Pos    int
	Line   int
	Cursor int
}

type Token struct {
	StartPos int
	EndPos   int
	Value    string
	Type     TokenType
}

type TokenType string

func NewLexer(code string) Lexer {
	return Lexer{
		Code:   code,
		Pos:    0,
		Line:   0,
		Cursor: 0}
}

var (
	ws = ` \r\t\f\v\n`
	nonTerminalRegex = regexp2.MustCompile(`[A-Z]+[0-9]*`, 0)
	terminalRegex = regexp2.MustCompile(`\'.\'|[a-z]`, 0)
	assignRegex = regexp2.MustCompile(`::=`, 0)
	nonTerminalKw = `non-terminal`
	terminalKw = `terminal`
	axiomKw = `axiom`
	epsilonKw = `epsilon`
	semicolon = `;`
	verticalBarRegex = `/\|/`
	regex = regexp2.MustCompile("^(?<nonTerminalKeyword>" + ")|^(?<terminalKeyword>${terminalKeywordRegex.source})|^(?<axiomKeyword>${axiomKeywordRegex.source})|^(?<epsilonKeyword>${epsilonKeywordRegex.source})|^(?<semicolon>${semicolonRegex.source})|^(?<comma>${commaRegex.source})|^(?<verticalBar>${verticalBarRegex.source})|^(?<assign>${assignRegex.source})|^(?<nonTerminal>${nonTerminalRegex.source})|^(?<terminal>${terminalRegex.source})", 0)

)



func (l *Lexer) NextToken() (Token, error) {

}
