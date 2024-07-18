// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import "wa-lang.org/wa/internal/wat/token"

// (type $$onFree (func (param i32)))
func (p *parser) parseModuleSection_type() {
	p.acceptToken(token.TYPE)

	panic("TODO")
}
