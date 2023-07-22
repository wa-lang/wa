package appast

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
)

func PrintAST(filename string) error {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(nil, fset, filename, nil, 0)
	if err != nil {
		return err
	}

	return ast.Print(fset, f)
}
