// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// import ::= (import mod:name nm:nname d:importdesc)
//
// importdesc ::= (memory id memtype)
//             |  (table id tabletype)
//             |  (global id globaltype)
//             |  (func id typeuse)

func (p *parser) parseModuleSection_import() *ast.ImportSpec {
	p.acceptToken(token.IMPORT)
	p.consumeComments()

	spec := &ast.ImportSpec{}

	// 宿主模块和名字
	spec.ObjModule = p.parseStringLit()
	p.consumeComments()

	spec.ObjName = p.parseStringLit()
	p.consumeComments()

	p.acceptToken(token.LPAREN)
	{
		p.consumeComments()
		switch p.tok {
		case token.MEMORY:
			p.parseModuleSection_import_memory(spec)
		case token.TABLE:
			p.parseModuleSection_import_table(spec)
		case token.FUNC:
			p.parseModuleSection_import_func(spec)
		case token.GLOBAL:
			p.parseModuleSection_import_global(spec)
		default:
			p.errorf(p.pos, "bad token, %v %s", p.tok, p.lit)
		}
		p.consumeComments()
	}
	p.acceptToken(token.RPAREN)

	return spec
}

func (p *parser) parseModuleSection_import_memory(spec *ast.ImportSpec) {
	p.acceptToken(token.MEMORY)
	p.consumeComments()

	spec.MemoryIdx = p.parseIntLit()
}

func (p *parser) parseModuleSection_import_table(spec *ast.ImportSpec) {
	p.acceptToken(token.TABLE)
	p.consumeComments()

	spec.TableIdx = p.parseIntLit()
}

func (p *parser) parseModuleSection_import_global(spec *ast.ImportSpec) {
	p.acceptToken(token.GLOBAL)

	spec.GlobalName = p.parseIdent()
	spec.GlobalType = p.parseNumberType()
}

func (p *parser) parseModuleSection_import_func(spec *ast.ImportSpec) {
	p.acceptToken(token.FUNC)
	p.consumeComments()

	spec.FuncType = &ast.FuncType{}
	spec.FuncName = p.parseIdent()

	for {
		p.consumeComments()
		if p.tok != token.LPAREN {
			break
		}

		p.acceptToken(token.LPAREN)
		p.consumeComments()

		switch p.tok {
		case token.PARAM:
			p.acceptToken(token.PARAM)
			p.consumeComments()

			if p.tok == token.IDENT {
				// (param $name i32)
				p.parseIdent()
				p.consumeComments()
				typ := p.parseNumberType()

				spec.FuncType.Params = append(spec.FuncType.Params, ast.Field{
					Type: typ,
				})
			} else {
				// (param i32)
				// (param i32 i64)
				for _, typ := range p.parseNumberTypeList() {
					spec.FuncType.Params = append(spec.FuncType.Params, ast.Field{
						Type: typ,
					})
				}
			}

			p.acceptToken(token.RPAREN)

		case token.RESULT:
			// (result i32)
			// (result i32 i64)
			p.acceptToken(token.RESULT)
			p.consumeComments()
			for _, typ := range p.parseNumberTypeList() {
				spec.FuncType.Results = append(spec.FuncType.Results, typ)
			}
			p.acceptToken(token.RPAREN)

		default:
			p.errorf(p.pos, "bad token: %v %s", p.tok, p.lit)
		}
	}
}
