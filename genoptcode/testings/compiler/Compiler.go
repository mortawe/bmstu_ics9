package compiler

type Compiler struct {
	parser Parser
}

func NewCompiler( input string) Compiler {
	return Compiler{
		parser: NewParser(input),
	}
}

func (c *Compiler) Exec() {

}