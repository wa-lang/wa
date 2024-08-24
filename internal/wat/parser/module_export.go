// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"strconv"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// export ::= (export nm:name d:exportdesc)
//
// exportdesc ::= (func funcidx)
//             |  (table tableidx)
//             |  (memory memidx)
//             |  (global globalidx)

func (p *parser) parseModuleSection_export() *ast.ExportSpec {
	p.acceptToken(token.EXPORT)

	spec := &ast.ExportSpec{}

	// 解析导出的名字, 字符串类型
	p.consumeComments()
	spec.Name = p.parseStringLit()

	// 导出的对象
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)
	defer p.consumeComments()

	// (memory 0)
	// (func $name)

	switch p.tok {
	case token.MEMORY:
		spec.Kind = p.tok
		p.acceptToken(token.MEMORY)
		p.consumeComments()

		if p.tok == token.IDENT {
			spec.MemoryIdx = p.parseIdent()
		} else {
			idx := p.parseIntLit()
			spec.MemoryIdx = strconv.Itoa(idx)
		}

	case token.TABLE:
		spec.Kind = p.tok
		p.acceptToken(token.TABLE)
		p.consumeComments()

		if p.tok == token.IDENT {
			spec.TableIdx = p.parseIdent()
		} else {
			idx := p.parseIntLit()
			spec.TableIdx = strconv.Itoa(idx)
		}

	case token.FUNC:
		spec.Kind = p.tok
		p.acceptToken(token.FUNC)
		p.consumeComments()

		if p.tok == token.IDENT {
			spec.FuncIdx = p.parseIdent()
		} else {
			idx := p.parseIntLit()
			spec.FuncIdx = strconv.Itoa(idx)
		}

	case token.GLOBAL:
		spec.Kind = p.tok
		p.acceptToken(token.GLOBAL)
		p.consumeComments()

		if p.tok == token.IDENT {
			spec.GlobalIdx = p.parseIdent()
		} else {
			idx := p.parseIntLit()
			spec.GlobalIdx = strconv.Itoa(idx)
		}

	default:
		p.errorf(p.pos, "unexpect %q", p.lit)
	}

	return spec
}
