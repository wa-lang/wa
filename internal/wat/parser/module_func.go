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
		Type: nil, // 作为标志
		Body: nil, // 作为标志
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
			if fn.Type != nil || fn.Body != nil {
				p.errorf(p.pos, "export must befor param/result/body")
			}

			p.acceptToken(token.EXPORT)

			p.consumeComments()
			fn.ExportName = p.parseStringLit()

			p.consumeComments()
			p.acceptToken(token.RPAREN)

		case token.PARAM:
			if fn.Body != nil {
				p.errorf(p.pos, "export must befor result/body")
			}
			if fn.Type != nil && len(fn.Type.Results) > 0 {
				p.errorf(p.pos, "export must befor result")
			}

			p.acceptToken(token.PARAM)

			if fn.Type == nil {
				fn.Type = &ast.FuncType{}
			}

			var field ast.Field
			p.consumeComments()
			if p.tok == token.IDENT {
				field.Name = p.parseIdent()
				p.consumeComments()
				field.Type = p.parseNumberType()
				fn.Type.Params = append(fn.Type.Params, field)
			} else {
				for _, x := range p.parseNumberTypeList() {
					fn.Type.Params = append(fn.Type.Params, ast.Field{Type: x})
				}
			}

			p.consumeComments()
			p.acceptToken(token.RPAREN)

		case token.RESULT:
			if fn.Body != nil {
				p.errorf(p.pos, "result must befor func body")
			}

			p.acceptToken(token.RESULT)

			if fn.Type == nil {
				fn.Type = &ast.FuncType{}
			}

			p.consumeComments()
			for _, x := range p.parseNumberTypeList() {
				fn.Type.Results = append(fn.Type.Results, x)
			}

			p.consumeComments()
			p.acceptToken(token.RPAREN)

		case token.LOCAL:
			if fn.Body == nil {
				fn.Body = &ast.FuncBody{}
			}
			if len(fn.Body.Insts) > 0 {
				p.errorf(p.pos, "local must befor instruction")
			}

			p.acceptToken(token.LOCAL)

			var field ast.Field
			p.consumeComments()
			if p.tok == token.IDENT {
				field.Name = p.parseIdent()
			}

			p.consumeComments()
			field.Type = p.parseNumberType()
			fn.Body.Locals = append(fn.Body.Locals, field)

			p.consumeComments()
			p.acceptToken(token.RPAREN)

		default:
			if p.tok.IsIsntruction() {
				if fn.Body == nil {
					fn.Body = &ast.FuncBody{}
				}

				ins := p.parseInstruction()
				fn.Body.Insts = append(fn.Body.Insts, ins)

				p.consumeComments()
				p.acceptToken(token.RPAREN)
			} else {
				break Loop
			}
		}
	}

	return fn
}
