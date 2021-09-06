package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func insertIntVar(file *ast.File, name string, value int) {
	var before, after []ast.Decl

	if len(file.Decls) > 0 {
		hasImport := false
		if genDecl, ok := file.Decls[0].(*ast.GenDecl); ok {
			hasImport = genDecl.Tok == token.IMPORT
		}

		if hasImport {
			before, after = []ast.Decl{file.Decls[0]}, file.Decls[1:]
		} else {
			after = file.Decls
		}
	}

	file.Decls = append(before,
		&ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{ast.NewIdent(name)},
					Type:  ast.NewIdent("int"),
					Values: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.INT,
							Value: fmt.Sprintf("%d", value),
						},
					},
				},
			},
		},
	)
	file.Decls = append(file.Decls, after...)
}

func countIncOnGo(file *ast.File, counterName string) {
	ast.Inspect(file, func(node ast.Node) bool {

		if blockStmt, ok := node.(*ast.BlockStmt); ok {
			for i := range blockStmt.List {
				if _, ok := blockStmt.List[i].(*ast.GoStmt); ok {
					newList := append(blockStmt.List[:i+1], &ast.IncDecStmt{
						X: &ast.Ident{
							Name: counterName,
							Obj:  nil,
						},
						Tok: token.INC,
					})
					blockStmt.List = append(newList, blockStmt.List[i:]...)
				}
			}
		}

		return true
	})
}

func insertCounterPrint(file *ast.File, counterName string) {
	isCurrentNodeMain := false
	ast.Inspect(file, func(node ast.Node) bool {
		if isCurrentNodeMain {
			if block, ok := node.(*ast.BlockStmt); ok {
				block.List = append(

					block.List, []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("fmt"),
									Sel: ast.NewIdent("Printf"),
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.INT,
										Value: counterName,
									},
								},
							},
						},
					}...,
				)
				isCurrentNodeMain = false
			}
		}
		if ident, ok := node.(*ast.Ident); ok {
			if ident.Name == "main" {
				isCurrentNodeMain = true
			}
		}
		return true
	})
}

func main() {
	if len(os.Args) != 2 {
		return
	}

	fset := token.NewFileSet()
	if file, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments); err == nil {
		insertIntVar(file, "routinesCounter", 0)
		countIncOnGo(file, "routinesCounter")
		insertCounterPrint(file, "routinesCounter")
		if format.Node(os.Stdout, fset, file) != nil {
			fmt.Printf("Formatter error: %v\n", err)
		}
		//ast.Fprint(os.Stdout, fset, file, nil)
	} else {
		fmt.Printf("Errors in %s\n", os.Args[1])
	}
}
