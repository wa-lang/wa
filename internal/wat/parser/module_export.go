// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import "wa-lang.org/wa/internal/wat/token"

// export ::= (export nm:name d:exportdesc)
//
// exportdesc ::= (func funcidx)
//             |  (table tableidx)
//             |  (memory memidx)
//             |  (global globalidx)

func (p *parser) parseModuleSection_export() {
	p.acceptToken(token.EXPORT)

	// 解析导出的名字, 字符串类型
	p.parseStringLit()

	// 导出的对象
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)

	// (memory 0)
	// (func $name)

	switch p.tok {
	case token.MEMORY:
		p.acceptToken(token.MEMORY)
		p.parseIntLit() // todo

	case token.FUNC:
		p.acceptToken(token.FUNC)
		p.parseIdent() // todo

	default:
		p.errorf(p.pos, "expect int, got %q", p.lit)
	}
}
