package internal

import "gonum.org/v1/gonum/graph/formats/dot/ast"

type BranchStmt struct {
	Stmt
	cond string
	positive Vertex
	negative Vertex
}
