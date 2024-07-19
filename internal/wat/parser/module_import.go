// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// import ::= (import mod:name nm:nname d:importdesc)
//
// importdesc ::= (func id typeuse)
//             |  (table id tabletype)
//             |  (memory id memtype)
//             |  (global id globaltype)

func (p *parser) parseModuleSection_import() *ast.ImportSpec {
	p.acceptToken(token.IMPORT)

	spec := &ast.ImportSpec{}

	// 宿主模块和名字
	p.consumeComments()
	spec.ModulePath = p.parseStringLit()

	p.consumeComments()
	spec.FuncPath = p.parseStringLit()

Loop:
	for {
		p.consumeComments()
		p.acceptToken(token.LPAREN)

		p.consumeComments()
		switch p.tok {
		case token.FUNC:
			p.parseModuleSection_import_func(spec)
			p.consumeComments()
			p.acceptToken(token.RPAREN)
		case token.TABLE:
			p.parseModuleSection_import_table(spec)
			p.consumeComments()
			p.acceptToken(token.RPAREN)
		case token.MEMORY:
			p.parseModuleSection_import_memory(spec)
			p.consumeComments()
			p.acceptToken(token.RPAREN)
		case token.GLOBAL:
			p.parseModuleSection_import_global(spec)
			p.consumeComments()
			p.acceptToken(token.RPAREN)
		}

		break Loop
	}

	return spec
}

func (p *parser) parseModuleSection_import_table(spec *ast.ImportSpec) {
	p.acceptToken(token.TABLE)

	panic("TODO")
}

func (p *parser) parseModuleSection_import_memory(spec *ast.ImportSpec) {
	p.acceptToken(token.MEMORY)

	panic("TODO")
}

func (p *parser) parseModuleSection_import_global(spec *ast.ImportSpec) {
	p.acceptToken(token.GLOBAL)

	panic("TODO")
}

func (p *parser) parseModuleSection_import_func(spec *ast.ImportSpec) {
	p.acceptToken(token.FUNC)

	spec.FuncName = p.parseIdent()
	spec.FuncType = &ast.FuncType{}

	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)

		switch p.tok {
		case token.PARAM:
			p.parseModuleSection_import_func_param(spec)
			p.acceptToken(token.RPAREN)
		case token.RESULT:
			p.parseModuleSection_import_func_result(spec)
			p.acceptToken(token.RPAREN)
			return
		default:
			p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		}
	}

	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)

		switch p.tok {
		case token.RESULT:
			p.parseModuleSection_import_func_result(spec)
			p.acceptToken(token.RPAREN)
		default:
			p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		}
	}
}

func (p *parser) parseModuleSection_import_func_param(spec *ast.ImportSpec) {
	p.acceptToken(token.PARAM)

	for {
		var field ast.Field
		if p.tok == token.IDENT {
			field.Name = p.lit
			p.next()
		}

		switch p.tok {
		case token.I32, token.I64, token.F32, token.F64:
			field.Type = p.lit
			spec.FuncType.Params = append(spec.FuncType.Params, field)
			p.next()
		case token.RPAREN:
			return
		default:
			p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		}
	}
}

func (p *parser) parseModuleSection_import_func_result(spec *ast.ImportSpec) {
	p.acceptToken(token.RESULT)

	for {
		switch p.tok {
		case token.I32, token.I64, token.F32, token.F64:
			spec.FuncType.ResultsType = append(spec.FuncType.ResultsType, p.lit)
			p.next()
		case token.RPAREN:
			return
		default:
			p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		}
	}
}
