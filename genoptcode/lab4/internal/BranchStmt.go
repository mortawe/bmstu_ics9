package internal

type BranchStmt struct {
	Stmt
	cond string
	positive Vertex
	negative Vertex
}
