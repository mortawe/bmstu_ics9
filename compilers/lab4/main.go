package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"unicode"
)

type TokenType string

var IsLetterOrDigit = regexp.MustCompile(`^([0-9]|[a-zA-Z])$`).MatchString
var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

var (
	Varname TokenType = "varname"
	Ident   TokenType = "ident"
	Unknown TokenType = "unknown"
	EOF     TokenType = "eof"
)

func splitTokens(input []byte, i, line, pos int) ([]Token, int, int, int, error) {
	tokens := []Token{}
	var currentTokenType TokenType
	j := 0

	for ; i < len(input); i++ {
		if input[i] == '\n' {
			line++
			pos = 0
			continue
		}
		if unicode.IsSpace(rune(input[i])) {
			pos++
			continue
		}
		j = i
		value := ""
		for {
			pos++
			if input[j] != 's' && input[j] != 'e' && input[j] != 't' {
				currentTokenType = Ident
				break
			}
			value += string(input[j])
			j++
			if j >= len(input) {
				return tokens, j, line, pos + j - i, fmt.Errorf("unexpected eof")
			}
			if IsLetterOrDigit(string(input[j])) {
				currentTokenType = Varname
				value += string(input[j])
				tokens = append(tokens, Token{
					StartPos:  pos,
					EndPos:    pos + 2,
					Line:      line,
					Value:     value,
					TokenType: Varname,
				})
				break
			}
			if input[j] != '.' {
				currentTokenType = Ident
				if j >= len(input) {
					return tokens, j, line, pos + j - i, fmt.Errorf("unexpected eof")
				}
			}
			value += string(input[j])
			j++
			if j >= len(input) {
				return tokens, j, line, pos + j - i, fmt.Errorf("unexpected eof")
			}
			for {
				if IsLetterOrDigit(string(input[j])) {
					value += string(input[j])
					j++
					if j >= len(input) {
						break
					}
				} else {
					j--
					break
				}
			}
			if len(value) < 3 {
				return tokens, j, line, pos + j - i, fmt.Errorf("short varname")
			}
			currentTokenType = Varname
			tokens = append(tokens, Token{
				StartPos:  pos,
				EndPos:    pos + j - i,
				Line:      line,
				Value:     value,
				TokenType: Varname,
			})
			break
		}
		if currentTokenType == Ident {
			if IsLetter(string(input[j])) {
				value += string(input[j])
				j++
				if j >= len(input) {
					tokens = append(tokens, Token{
						StartPos:  pos,
						EndPos:    pos + j - i,
						Line:      line,
						Value:     value,
						TokenType: Ident,
					})
					break
				}
			} else {
				currentTokenType = Unknown
			}
			if currentTokenType == Unknown {
				return tokens, j, line, pos + j - i, fmt.Errorf("unknown token")
			}
			for {
				if IsLetterOrDigit(string(input[j])) {
					value += string(input[j])
					j++
					if j >= len(input) {
						break
					}
				} else {
					j--
					break
				}
			}
			tokens = append(tokens, Token{
				StartPos:  pos,
				EndPos:    pos + j - i,
				Line:      line,
				Value:     value,
				TokenType: Ident,
			})
		}

		if currentTokenType == EOF {
			return tokens, j, line, pos + j - i, nil
		}
		pos += j - i
		i = j
	}

	tokens = append(tokens, Token{
		StartPos:  pos,
		EndPos:    pos,
		Line:      line,
		Value:     "",
		TokenType: EOF,
	})
	return tokens, 0, line, pos, nil
}

func main() {
	content, err := ioutil.ReadFile("./input.txt")

	if err != nil {
		log.Fatal(err)
	}
	line := 1
	pos := -1
	i := -1
	tokens := []Token{}
	errors := []error{}
	for {
		tokens, i, line, pos, err = splitTokens(content, i+1, line, pos+1)
		for _, t := range tokens {
			fmt.Println(t.String())
		}
		if err != nil {
			errors = append(errors, fmt.Errorf("err : %v - (%d, %d) \n", err, line, pos))
		}
		if err == nil {
			break
		}
	}

	for _, err := range errors {
		fmt.Println(err)
	}
}
