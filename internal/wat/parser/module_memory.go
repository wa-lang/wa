// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *parser) parseModuleSection_memory() {
	p.acceptToken(token.MEMORY)

	if p.module.Memory == nil {
		p.module.Memory = &ast.Memory{}
	}

	if p.tok == token.IDENT {
		p.module.Memory.Name = p.lit
		p.acceptToken(token.IDENT)
	}

	p.module.Memory.Pages = p.parseIntLit()
}
