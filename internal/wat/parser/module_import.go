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
	spec.ObjName = p.parseStringLit()

	p.acceptToken(token.LPAREN)
	{
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
	}
	p.acceptToken(token.RPAREN)

	return spec
}

// (import "env" "memory" (memory 1))
func (p *parser) parseModuleSection_import_memory(spec *ast.ImportSpec) {
	p.acceptToken(token.MEMORY)
	spec.ObjKind = token.MEMORY

	spec.Memory = &ast.Memory{}
	if p.tok == token.IDENT {
		spec.Memory.Name = p.parseIdent()
	}
	spec.Memory.Pages = p.parseIntLit()
	if p.tok == token.INT {
		spec.Memory.MaxPages = p.parseIntLit()
	}
}

func (p *parser) parseModuleSection_import_table(spec *ast.ImportSpec) {
	p.acceptToken(token.TABLE)
	p.consumeComments()

	spec.ObjKind = token.TABLE
	spec.TableName = p.parseIdentOrIndex()
}

func (p *parser) parseModuleSection_import_global(spec *ast.ImportSpec) {
	p.acceptToken(token.GLOBAL)

	spec.ObjKind = token.GLOBAL
	spec.GlobalName = p.parseIdent()
	spec.GlobalType = p.parseNumberType()
}

func (p *parser) parseModuleSection_import_func(spec *ast.ImportSpec) {
	p.acceptToken(token.FUNC)

	spec.ObjKind = token.FUNC
	spec.FuncType = &ast.FuncType{}

	if p.tok == token.IDENT {
		spec.FuncName = p.parseIdent()
	}

	for {
		if p.tok != token.LPAREN {
			break
		}

		p.acceptToken(token.LPAREN)
		switch p.tok {
		case token.PARAM:
			p.acceptToken(token.PARAM)

			if p.tok == token.IDENT {
				// (param $name i32)
				name := p.parseIdent()
				typ := p.parseNumberType()

				spec.FuncType.Params = append(spec.FuncType.Params, ast.Field{
					Name: name,
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
			for _, typ := range p.parseNumberTypeList() {
				spec.FuncType.Results = append(spec.FuncType.Results, typ)
			}
			p.acceptToken(token.RPAREN)

		default:
			p.errorf(p.pos, "bad token: %v %s", p.tok, p.lit)
		}
	}
}
