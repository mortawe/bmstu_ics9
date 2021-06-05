package lab4

type Expr struct {
	Vars []Var
}

func NewExpr() *Expr {
	return &Expr{
		Vars: make([]Var, 0),
	}
}

func (e *Expr) String() string {
	result := ""
	for _, e := range e.Vars {
		result += e.String()
	}
	return result
}

