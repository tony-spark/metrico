package osexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check there is no os.Exit in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					ast.Inspect(x, func(node ast.Node) bool {
						if c, ok := node.(*ast.CallExpr); ok {
							if f, ok := c.Fun.(*ast.SelectorExpr); ok {
								x, okx := f.X.(*ast.Ident)
								if okx && x.Name == "os" && f.Sel.Name == "Exit" {
									pass.Reportf(f.Pos(), "os.Exit in main")
									return false
								}
							}
						}
						return true
					})
					return false
				}
			}
			return true
		})
	}
	return nil, nil
}
