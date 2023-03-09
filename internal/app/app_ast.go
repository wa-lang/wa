// 版权 @2023 凹语言 作者。保留所有权利。

package app

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
)

func (p *App) AST(filename string) error {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(nil, fset, filename, nil, 0)
	if err != nil {
		return err
	}

	return ast.Print(fset, f)
}
