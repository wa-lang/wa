// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import "wa-lang.org/wa/internal/wat/token"

// elem ::= (elem id? elemlist)
//       |  (elem id? x:tableuse (offset e:expr) elemlist)
//       |  (elem id? declare elemlist)
//
// elemlist ::= reftype vec(elemexpr)
// elemexpr ::= (item e:expr)
// tableuse ::= (table x:tableidx)

// (elem (i32.const 1) $$u8.$$block.$$onFree)
func (p *parser) parseModuleSection_elem() {
	p.acceptToken(token.ELEM)

	p.consumeComments()
	p.acceptToken(token.LPAREN)

	p.consumeComments()
	p.acceptToken(token.INS_I32_CONST)

	p.consumeComments()
	p.parseIntLit()

	p.consumeComments()
	p.acceptToken(token.RPAREN)

	for {
		p.consumeComments()
		if p.tok != token.RPAREN {
			p.parseIdent()
		}
	}
}
