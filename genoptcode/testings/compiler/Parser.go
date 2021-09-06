package compiler

import "errors"

var lastID = 0
type Vertex struct {
	id int
	parents []*Vertex
	children []*Vertex
}
var typeID = 0
type Var struct {
	id int
	name string
}

type Parser struct {
	Lexer Lexer
	current int

	Tree []Vertex

	Vars []Var
}

func NewParser(input string) Parser {
	return Parser{
		Lexer: Lexer{
			input:   input,
			current: 0,
		},
		Tree:   []Vertex{{
			id:       0,
			parents:  nil,
			children: nil,
		}},
		Vars:   []Var{},
	}
}

func (p *Parser) Parse() error {
	currentNode := &p.Tree[0]
	var err error
	for currentNode != nil {
		currentNode, err = p.parseAxiom(currentNode)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *Parser) newVar(token Token) {
	v := Var{id: typeID, name: token.Value}
	p.Vars = append(p.Vars, v)
	typeID++
}
func (p *Parser) parseAxiom(currentNode *Vertex) (*Vertex, error) {
	token, err := p.Lexer.nextToken()
	if err != nil {
		return nil, err
	}
	switch token.Tag {
	case TagVar:
		token, err = p.Lexer.nextToken()
		if err != nil {
			return nil, err
		}
		if token.Tag != TagID {
			return nil, errors.New("expected id after var")
		}
		p.newVar(token)
		token, err =
	}

}

func (p *Parser) parseStmt(currentNode *Vertex)