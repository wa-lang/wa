// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import "wa-lang.org/wa/internal/wat/token"

// (data (i32.const 8) "hello world\n")
func (p *parser) parseModuleSection_data() {
	p.acceptToken(token.DATA)

	p.consumeComments()
	p.acceptToken(token.LPAREN)

	p.consumeComments()
	p.acceptToken(token.INS_I32_CONST)

	p.consumeComments()
	p.parseIntLit()

	p.consumeComments()
	p.acceptToken(token.RPAREN)

	p.consumeComments()
	p.parseStringLit()
}
