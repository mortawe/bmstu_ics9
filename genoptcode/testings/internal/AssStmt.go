package internal

type AssStmt struct {
	Stmt
}

func (s *AssStmt) UpdateRhsVarVersion(version, indexInRhs int) bool {
	return false
}