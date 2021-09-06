package compiler

import (
	"errors"

	"github.com/dlclark/regexp2"
)

var skipWSRegex = regexp2.MustCompile("[ \\r\\t\\f\\v\\n]*", 0)
var idRegex = regexp2.MustCompile("[a-zA-Z][a-zA-Z0-9]*", 0)
var numberRegex = regexp2.MustCompile("[0-9]+", 0)

type Lexer struct {
	input   string
	current int
}

func (l *Lexer) suffix() string {
	return l.input[l.current:]
}

func (l *Lexer) nextToken() (Token, error) {
	match, err := skipWSRegex.FindStringMatchStartingAt(l.suffix(), 0)
	if err != nil {
		return Token{}, err
	}
	l.current += match.Length
	match, err = idRegex.FindStringMatchStartingAt(l.suffix(), 0)
	if err != nil {
		return Token{}, err
	}
	if match != nil && match.Index == 0 && match.String() != "" {
		l.current += match.Length
		switch match.String() {
		case "var":
			return Token{Tag: TagVar}, nil
		case "if":
			return Token{Tag: TagIf}, nil
		case "then":
			return Token{Tag: TagThen}, nil
		case "else":
			return Token{Tag: TagElse}, nil
		default:
			return Token{Tag: TagID, Value: match.String()}, nil
		}
	}
	match, err = numberRegex.FindStringMatchStartingAt(l.suffix(), 0)
	if err != nil {
		return Token{}, err
	}
	if match != nil && match.Index == 0 && match.String() != "" {
		l.current += match.Length
		return Token{
			Tag:   TagNumber,
			Value: match.String(),
		}, nil
	}
	suffix := l.suffix()
	if len(l.suffix()) == 0 {
		return Token{}, errors.New("eof")
	}
	l.current += 1
	if suffix[0] == '=' {
		return Token{Tag: TagAssign}, nil
	}
	if suffix[0] == '+' {
		return Token{Tag: TagPlus}, nil
	}
	if suffix[0] == ';' {
		return Token{Tag: TagSC}, nil
	}
	if suffix[0] == '>' {
		return Token{Tag: TagMore}, nil
	}
	if suffix[0] == '<' {
		return Token{Tag: TagLess}, nil
	}
	return Token{}, errors.New("unknown token")
}
