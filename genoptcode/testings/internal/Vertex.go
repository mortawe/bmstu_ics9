package internal

type Vertex struct {
	Name     string
	precs    []*Vertex
	succs    map[string]*Vertex
	children map[string]*Vertex

	immediateDom *Vertex
	stmts        []*StmtI
	phis         []*StmtI
}

func NewVertex(name string) *Vertex {
	return &Vertex{
		Name:         "",
		precs:        []*Vertex{},
		succs:        map[string]*Vertex{},
		children:     map[string]*Vertex{},
		immediateDom: nil,
		stmts:        []*StmtI{},
		phis:         []*StmtI{},
	}
}

func (v *Vertex) prependStmt(s StmtI) {
	t := []*StmtI{&s}
	t = append(t, v.stmts...)
	v.stmts = t
}

func (v *Vertex) String() string {
	stmts := "\n"
	for _, s := range  v.stmts {
		stmts += "\t" + (*s).String()+ "\n"
	}
	return v.Name + "{" + stmts + "}\n"
}