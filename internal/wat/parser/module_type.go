// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// type ::= (type id? functype)

// (type $$OnFree (func (param i32)))
// (type $$wa.runtime.comp (func (param i32) (param i32) (result i32)))
func (p *parser) parseModuleSection_type() *ast.TypeSection {
	p.acceptToken(token.TYPE)

	typ := &ast.TypeSection{
		Type: &ast.FuncType{},
	}

	if p.tok == token.IDENT {
		typ.Name = p.parseIdent()
	}

	p.acceptToken(token.LPAREN)
	p.parseModuleSection_type_funcType(typ.Type)
	p.acceptToken(token.RPAREN)

	return typ
}

func (p *parser) parseModuleSection_type_funcType(typ *ast.FuncType) {
	p.acceptToken(token.FUNC)

	for {
		if p.tok != token.LPAREN {
			break
		}

		p.acceptToken(token.LPAREN)
		switch p.tok {
		case token.PARAM:
			p.parseModuleSection_type_funcType_param(typ)
			p.acceptToken(token.RPAREN)
		case token.RESULT:
			p.parseModuleSection_type_funcType_result(typ)
			p.acceptToken(token.RPAREN)
		default:
			p.errorf(p.pos, "bad token %v, %q", p.tok, p.lit)
		}
	}
}

// (param i32)
// (param $release_func i32)
func (p *parser) parseModuleSection_type_funcType_param(typ *ast.FuncType) {
	p.acceptToken(token.PARAM)

	if p.tok == token.IDENT {
		var field ast.Field
		field.Name = p.parseIdent()
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
	p.acceptToken(token.RESULT)

	for _, x := range p.parseNumberTypeList() {
		typ.Results = append(typ.Results, x)
	}
}
