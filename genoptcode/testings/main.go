package main

import "lab4/compiler"

func main(){
	comp :=compiler.NewCompiler("var x = 5; if x > 3 then x = x + 1 else x = x + 2;")
	comp.Exec()
}