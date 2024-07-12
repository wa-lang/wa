// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func ParseModule(file *token.File, src []byte) (f *ast.Module, err error) {
	var p parser
	p.init(file, src)
	f = p.parseFile()
	return
}
