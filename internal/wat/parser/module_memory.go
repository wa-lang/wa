// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// mem ::= (memory id? memtype)
//
// memtype ::= lim:limits
//
// limits ::= n:u32 | n:u32 m:u32

// (memory $memory 1024)
// (memory $memory 1024 2028)
func (p *parser) parseModuleSection_memory() *ast.Memory {
	p.acceptToken(token.MEMORY)

	mem := &ast.Memory{}

	p.consumeComments()
	if p.tok == token.IDENT {
		mem.Name = p.lit
		p.acceptToken(token.IDENT)
	}

	p.consumeComments()
	mem.Pages = p.parseIntLit()

	p.consumeComments()
	if p.tok == token.INT {
		mem.MaxPages = p.parseIntLit()
	}

	return mem
}
