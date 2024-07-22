// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import "wa-lang.org/wa/internal/wat/token"

// start ::= (start funcidx)

func (p *parser) parseModuleSection_start() string {
	p.acceptToken(token.START)

	p.consumeComments()
	return p.parseIdent()
}
