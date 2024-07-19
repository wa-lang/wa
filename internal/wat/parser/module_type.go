// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// type ::= (type id? functype)

// (type $$onFree (func (param i32)))
// (type $$wa.runtime.comp (func (param i32) (param i32) (result i32)))
func (p *parser) parseModuleSection_type() *ast.TypeSection {
	p.acceptToken(token.TYPE)

	typ := &ast.TypeSection{
		Type: &ast.FuncType{},
	}

	p.consumeComments()
	if p.tok == token.IDENT {
		typ.Name = p.parseIdent()
	}

	p.consumeComments()
	p.parseModuleSection_type_funcType(typ.Type)

	return typ
}

func (p *parser) parseModuleSection_type_funcType(typ *ast.FuncType) {
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
			p.parseModuleSection_type_funcType_param(typ)
		case token.RESULT:
			p.parseModuleSection_type_funcType_result(typ)
		}
	}
}

// (param i32)
// (param $release_func i32)
func (p *parser) parseModuleSection_type_funcType_param(typ *ast.FuncType) {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)

	p.consumeComments()
	p.acceptToken(token.PARAM)

	p.consumeComments()
	if p.tok == token.IDENT {
		var field ast.Field
		field.Name = p.parseIdent()
		p.consumeComments()
		field.Type = p.parseNumberType()
		typ.Params = append(typ.Params, field)
	} else {
		for _, x := range p.parseNumberTypeList() {
			typ.Params = append(typ.Params, ast.Field{Type: x})
		}
	}
}

// (result i32)
// (result i32 i32)
func (p *parser) parseModuleSection_type_funcType_result(typ *ast.FuncType) {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)

	p.consumeComments()
	p.acceptToken(token.RESULT)

	p.consumeComments()
	for _, x := range p.parseNumberTypeList() {
		typ.Results = append(typ.Results, x)
	}
}
