package main

import (
	"fmt"
	"unicode"
)

type AutomataLexer struct {
	Line int
	Pos  int

	Tokens []Token

	cursor int
	input  []byte
}

func NewAutomataLexer(input []byte) *AutomataLexer {
	return &AutomataLexer{
		Line:   1,
		Pos:    0,
		Tokens: make([]Token, 0),
		cursor: -1,
		input:  input,
	}
}

func (l *AutomataLexer) consume() {
	l.cursor++
	l.Pos++
	if l.isEOF() {
		return
	}
	if l.isNewLine() {
		l.Line++
		// l.cursor++
		l.Pos = 0
	}
}

func (l *AutomataLexer) peek() byte {
	if l.cursor+1 >= len(l.input) {
		return byte(0)
	}
	sym := l.input[l.cursor+1]
	return sym
}

func (l *AutomataLexer) nextToken() (Token, error) {
	curState := ST
	prevState := DomainTag(0)

	startCursor := l.cursor + 1
	token := Token{
		StartPos: l.Pos,
		Line:     l.Line,
	}
	prevState = curState

	for curState != UK && !l.isEOF() {
		curSymbol := l.peek()
		curSymType := getCol(curSymbol)
		if curSymType == Eof {
			if l.cursor < startCursor {
				l.consume()
				return Token{
					TokenType: "EOF",
				}, nil
			}
			break
		}
		curState = transes[curState][curSymType]
		// fmt.Println("symbol ", curSymbol, " type ", curSymType, " state ", curState)
		if curState != UK {
			prevState = curState
			l.consume()
		}
	}
	if l.cursor < startCursor && curState == UK {
		l.consume()
		return token, fmt.Errorf("%s, prev state : %v", token.String(), prevState)

	}
	token.Value = string(l.input[startCursor : l.cursor+1])

	// token.Value = strings.TrimSpace(string(l.input[startCursor : l.cursor+1]))
	token.EndPos = l.Pos

	switch prevState {
	case NU:
		token.TokenType = NUMBERS
	case ID:
		token.TokenType = IDENT
	case CM:
		token.TokenType = COMMENTARY
	case C2:
		token.TokenType = COMMENTARY
	case C3:
		token.TokenType = COMMENTARY
	case OP:
		token.TokenType = OPERATORS
	case KW:
		token.TokenType = KEYWORDS
	case WS:
		token.TokenType = WHITESPACE
	case UK:
		if l.isEOF() {
			return Token{
				// Value:     string(l.input[startCursor : curSymbol-1]),
				TokenType: "EOF",
			}, nil
		}
		return token, fmt.Errorf("%s, prev state : %v", token.String(), prevState)
	default:
		// fmt.Println(prevState, curState)
		return token, fmt.Errorf("%s, prev state : %v", token.String(), prevState)
	}
	// token.Value = string(l.input[startCursor:l.cursor])
	return token, nil

}

func (l *AutomataLexer) isEOF() bool {
	return l.cursor >= len(l.input)
}

func (l *AutomataLexer) isNewLine() bool {
	return l.input[l.cursor] == '\n'
}

func getCol(b byte) int {
	if b == '\n' {
		return NewLine
	}
	if unicode.IsSpace(rune(b)) {
		return Spaces
	}
	if unicode.IsDigit(rune(b)) {
		return Digit
	}
	switch b {
	case 'u':
		return LetterU
	case 'p':
		return LetterP
	case 'd':
		return LetterD
	case 'a':
		return LetterA
	case 't':
		return LetterT
	case 'e':
		return LetterE
	case 'w':
		return LetterW
	case 'h':
		return LetterH
	case 'r':
		return LetterR
	case '/':
		return BackSlash
	case '\\':
		return Slash
	case '!':
		return NotEqual
	case '=':
		return Equals
	case 0:
		return Eof
	default:
		if unicode.IsLetter(rune(b)) {
			return Letter
		}
		return Others
	}
}
