package compiler

type Compiler struct {
	lexer Lexer
	parser Parser
}