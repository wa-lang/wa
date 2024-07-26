// 版权 @2024 凹语言 作者。保留所有权利。

package parser

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// func ::= (func id? typeuse (t:local)* (in:instr)*)
//
// local := (local id? t:valtype)

// (func $runtime.i32_load (export "runtime.i32_load") (param $addr i32) (result i32)

func (p *parser) parseModuleSection_func() *ast.Func {
	p.acceptToken(token.FUNC)

	fn := &ast.Func{
		Type: &ast.FuncType{},
		Body: &ast.FuncBody{},
	}

	p.consumeComments()
	if p.tok == token.IDENT {
		fn.Name = p.parseIdent()
	}

Loop:
	// 解析 export/param/result/local
	for {
		p.consumeComments()
		if p.tok != token.LPAREN {
			break
		}

		p.acceptToken(token.LPAREN)
		switch p.tok {
		case token.EXPORT:
			p.acceptToken(token.EXPORT)
			fn.ExportName = p.parseStringLit()
			p.acceptToken(token.RPAREN)

		case token.PARAM:
			p.acceptToken(token.PARAM)

			var field ast.Field
			if p.tok == token.IDENT {
				field.Name = p.parseIdent()
				field.Type = p.parseNumberType()
				fn.Type.Params = append(fn.Type.Params, field)
			} else {
				for _, x := range p.parseNumberTypeList() {
					fn.Type.Params = append(fn.Type.Params, ast.Field{Type: x})
				}
			}

			p.acceptToken(token.RPAREN)

		case token.RESULT:
			p.acceptToken(token.RESULT)
			for _, x := range p.parseNumberTypeList() {
				fn.Type.Results = append(fn.Type.Results, x)
			}
			p.acceptToken(token.RPAREN)

		case token.LOCAL:
			if len(fn.Body.Insts) > 0 {
				p.errorf(p.pos, "local must befor instruction")
			}

			p.acceptToken(token.LOCAL)

			var field ast.Field
			if p.tok == token.IDENT {
				field.Name = p.parseIdent()
			}

			field.Type = p.parseNumberType()
			fn.Body.Locals = append(fn.Body.Locals, field)

			p.acceptToken(token.RPAREN)

		default:
			break Loop
		}
	}

	// 解析指令
	// 不支持指令折叠
	if fn.Body == nil {
		fn.Body = &ast.FuncBody{}
	}
	for {
		p.consumeComments()
		if !p.tok.IsIsntruction() {
			break
		}

		ins := p.parseInstruction()
		fn.Body.Insts = append(fn.Body.Insts, ins)
	}

	return fn
}
