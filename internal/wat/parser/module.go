// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import "wa-lang.org/wa/internal/wat/token"

// 解析 module
// (module $name ...)
func (p *parser) parseModule() {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)
	defer p.consumeComments()

	// module 关键字
	p.consumeComments()
	p.acceptToken(token.MODULE)

	// 模块名字
	p.consumeComments()
	if p.tok == token.IDENT {
		p.module.Name = p.lit
		p.next()
	}

	// 解析模块Section
	for {
		p.consumeComments()
		if p.tok == token.LPAREN {
			p.parseModuleSection()
		} else {
			break
		}
	}
}

// 解析模块段
func (p *parser) parseModuleSection() {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)
	defer p.consumeComments()

	p.consumeComments()

	switch p.tok {
	default:
		p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)

	case token.TYPE:
		p.parseModuleSection_type()

	case token.IMPORT:
		p.parseModuleSection_import()
	case token.EXPORT:
		p.parseModuleSection_export()

	case token.MEMORY:
		p.parseModuleSection_memory()
	case token.DATA:
		p.parseModuleSection_data()

	case token.TABLE:
		p.parseModuleSection_table()
	case token.ELEM:
		p.parseModuleSection_elem()

	case token.GLOBAL:
		p.parseModuleSection_global()

	case token.FUNC:
		p.parseModuleSection_func()

	case token.START:
		p.parseModuleSection_start()
	}
}
