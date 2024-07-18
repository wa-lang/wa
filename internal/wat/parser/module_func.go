// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *parser) parseModuleSection_func() {
	p.acceptToken(token.FUNC)

	// (func $main (export "_start")

	if p.tok == token.IDENT {
		p.acceptToken(token.IDENT)
	}

	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)

	Loop0:
		for {
			p.consumeComments()
			switch p.tok {
			case token.EXPORT:
				p.acceptToken(token.EXPORT)
				p.parseStringLit()
				p.acceptToken(token.RPAREN)

			case token.PARAM:
				p.next()
				if p.tok == token.IDENT {
					p.parseIdent()
				}
				switch p.tok {
				case token.I32, token.I64, token.F32, token.F64:
					p.next()
					p.acceptToken(token.RPAREN)
				default:
					p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
				}
			default:
				break Loop0
			}
		}
	}

	for {
		p.consumeComments()
		if p.tok == token.LPAREN {
			p.acceptToken(token.LPAREN)
		} else {
			break
		}

		switch {
		case p.tok.IsIsntruction():
			p.parseInstruction()
			p.acceptToken(token.RPAREN)
		default:
			p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		}
	}
}

func (p *parser) parseImportFuncType(spec *ast.ImportSpec) {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)

	p.acceptToken(token.FUNC)

	spec.FuncName = p.parseIdent()
	spec.FuncType = &ast.FuncType{}

	if p.tok == token.LPAREN {
		p.acceptToken(token.LPAREN)

		switch p.tok {
		case token.PARAM:
			p.parseImportFuncType_param(spec)
			p.acceptToken(token.RPAREN)
		case token.RESULT:
			p.parseImportFuncType_result(spec)
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
			p.parseImportFuncType_result(spec)
			p.acceptToken(token.RPAREN)
		default:
			p.errorf(p.pos, "bad token: %v, lit: %q", p.tok, p.lit)
		}
	}
}

func (p *parser) parseImportFuncType_param(spec *ast.ImportSpec) {
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

func (p *parser) parseImportFuncType_result(spec *ast.ImportSpec) {
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
