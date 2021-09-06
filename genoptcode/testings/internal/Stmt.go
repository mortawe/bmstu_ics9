package internal

type Stmt struct {
	Lhs Var
	Rhs Expr

	isPhi bool
	isAss bool
}

type StmtI interface {
	UpdateRhsVarVersion(version, indexInRhs int) bool
	String() string
}

func (s *Stmt) RenameRhsVar(name string, version int) {
	for i, v := range s.Rhs.Vars {
		if v.name == name {
			t := NewVar(name, v.sign, version)
			s.Rhs.Vars[i] = *t
		}
	}
}

func (s *Stmt) RenameLhsVar(name string, version int) bool {
	if name == s.Lhs.name {
		t := NewVar(name, s.Lhs.sign, version)
		s.Lhs = *t
		return true
	}
	return false
}

