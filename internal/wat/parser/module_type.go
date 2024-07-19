// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import "wa-lang.org/wa/internal/wat/token"

// type ::= (type id? functype)

// (type $$onFree (func (param i32)))
// (type $$wa.runtime.comp (func (param i32) (param i32) (result i32)))
func (p *parser) parseModuleSection_type() {
	p.acceptToken(token.TYPE)

	p.consumeComments()
	if p.tok == token.IDENT {
		p.parseIdent()
	}

	p.consumeComments()
	p.parseModuleSection_type_funcType()
}

func (p *parser) parseModuleSection_type_funcType() {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)
	defer p.consumeComments()

	for {
		p.consumeComments()
		if p.tok != token.LPAREN {
			break
		}

		switch p.tok {
		case token.PARAM:
			p.parseModuleSection_type_funcType_param()
		case token.RESULT:
			p.parseModuleSection_type_funcType_result()
		}
	}
}

// (param i32)
// (param $release_func i32)
func (p *parser) parseModuleSection_type_funcType_param() {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)

	p.consumeComments()
	p.acceptToken(token.PARAM)

	p.consumeComments()
	if p.tok == token.IDENT {
		p.parseIdent()
	}

	p.consumeComments()
	p.parseNumberType()
}

// (result i32)
// (result i32 i32)
func (p *parser) parseModuleSection_type_funcType_result() {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)

	p.consumeComments()
	p.acceptToken(token.RESULT)

	for {
		p.consumeComments()
		if p.tok != token.RPAREN {
			p.parseNumberType()
		}
	}
}
