package lab4

type Vertex struct {
	Name     string
	precs    []*Vertex
	succs    map[string]*Vertex
	children map[string]*Vertex

	immediateDom *Vertex
	stmts        []*Stmt
	phis         []*Stmt
}

func NewVertex(name string) *Vertex {
	return &Vertex{
		Name:         "",
		precs:        []*Vertex{},
		succs:        map[string]*Vertex{},
		children:     map[string]*Vertex{},
		immediateDom: nil,
		stmts:        []*Stmt{},
		phis:         []*Stmt{},
	}
}

func (v *Vertex) prependStmt(s Stmt) {
	t := []*Stmt{&s}
	t = append(t, v.stmts...)
	v.stmts = t
}

func (v *Vertex) String() string {
	stmts := "\n"
	for _, s := range  v.stmts {
		stmts += "\t" + s.String() + "\n"
	}
	return v.Name + "{" + stmts + "}\n"
}