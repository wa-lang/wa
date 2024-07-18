// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import "wa-lang.org/wa/internal/wat/token"

// (global $__stack_ptr (mut i32) (i32.const 1024))     ;; index=0
// (global $__heap_max  i32       (i32.const 67108864)) ;; 64MB, 1024 page
func (p *parser) parseModuleSection_global() {
	p.acceptToken(token.GLOBAL)

	p.consumeComments()
	p.parseIdent()

	panic("TODO")
}
