// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import "wa-lang.org/wa/internal/wat/token"

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

	p.consumeComments()
	p.parseIdent()
}
