// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// data ::= (data id? b:datastring)
//       |  (data id? x:memuse (offset e:expr) b:datastring)
//
// datastring ::= (b:string)
//
// memuse ::= (memory x:memidx)

// (data (i32.const 8) "hello world\n")
func (p *parser) parseModuleSection_data() *ast.DataSection {
	p.acceptToken(token.DATA)

	dataSection := &ast.DataSection{}

	if p.tok == token.IDENT {
		dataSection.Name = p.parseIdent()
	}

	p.acceptToken(token.LPAREN)
	p.acceptToken(token.INS_I32_CONST)
	dataSection.Offset = uint32(p.parseIntLit())
	p.acceptToken(token.RPAREN)

	dataSection.Value = []byte(p.parseStringLit())

	return dataSection
}
