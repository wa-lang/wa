// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import "wa-lang.org/wa/internal/wat/token"

// (table 3 funcref)
func (p *parser) parseModuleSection_table() {
	p.acceptToken(token.TABLE)

	p.consumeComments()
	p.parseIntLit()

	p.consumeComments()
	p.acceptToken(token.FUNCREF)
}
