// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// func $add(%a:i32, %b:i32, %c:i32) => f64 {
//     local %d: i32 # 局部变量必须先声明, i32 大小的空间
//
//     # 指令
// Loop:
// }

func (p *parser) parseFunc() *ast.Func {
	fn := &ast.Func{
		Type: new(ast.FuncType),
		Body: new(ast.FuncBody),
	}

	p.acceptTokenAorB(token.FUNC, token.FUNC_zh)
	fn.Name = p.parseIdent()

	if p.tok == token.LPAREN {
		p.parseFunc_args(fn)
	}
	if p.tok == token.ARROW {
		p.parseFunc_return(fn)
	}

	p.parseFunc_body(fn)
	p.consumeTokenList(token.SEMICOLON)

	return fn
}

func (p *parser) parseFunc_args(fn *ast.Func) {
	p.acceptToken(token.LPAREN)
	defer p.acceptToken(token.RPAREN)

	for p.tok == token.IDENT {
		argPos := p.pos
		argName := p.parseIdent()
		argType := token.NONE

		p.acceptToken(token.COLON)

		switch p.tok {
		case token.I32, token.I64, token.F32, token.F64, token.PTR:
			argType = p.tok
			p.acceptToken(p.tok)
		default:
			p.errorf(p.pos, "expect argument type(i32/i64/f32/f64/ptr), got %v", p.tok)
		}
		fn.Type.Args = append(fn.Type.Args, ast.Argument{
			Pos:  argPos,
			Name: argName,
			Type: argType,
		})

		if p.tok == token.COMMA {
			p.acceptToken(token.COMMA)
			continue
		} else {
			break
		}
	}
}

func (p *parser) parseFunc_return(fn *ast.Func) {
	p.acceptToken(token.ARROW)

	switch p.tok {
	case token.I32, token.I64, token.F32, token.F64, token.PTR:
		fn.Type.Return = p.tok
		p.acceptToken(p.tok)
	default:
		p.errorf(p.pos, "expect return type(i32/i64/f32/f64/ptr), got %v", p.tok)
	}
}

func (p *parser) parseFunc_body(fn *ast.Func) {
	p.acceptToken(token.LBRACE)
	defer p.acceptToken(token.RBRACE)

	for p.tok == token.LOCAL || p.tok == token.LOCAL_zh {
		p.parseFunc_body_local(fn)
	}

	switch {
	case p.cpu == abi.RISCV64 || p.cpu == abi.RISCV32:
		for p.tok == token.IDENT || p.tok.IsAs() {
			fn.Body.Insts = append(fn.Body.Insts, p.parseInst_riscv())
		}
		return
	default:
		panic("unreachable")
	}
}

func (p *parser) parseFunc_body_local(fn *ast.Func) {
	for p.tok == token.LOCAL {
		localPos := p.pos
		p.acceptTokenAorB(token.LOCAL, token.LOCAL_zh)
		localName := p.parseIdent()
		p.acceptToken(token.COLON)

		switch p.tok {
		case token.I32, token.I64,
			token.U32, token.U64,
			token.F32, token.F64,
			token.PTR,
			token.I32_zh, token.I64_zh,
			token.U32_zh, token.U64_zh,
			token.F32_zh, token.F64_zh,
			token.PTR_zh:
			fn.Body.Locals = append(fn.Body.Locals, ast.Local{
				Pos:  localPos,
				Name: localName,
				Type: p.tok,
			})
			p.acceptToken(p.tok)
		default:
			p.errorf(p.pos, "expect local type(i32/i64/f32/f64/ptr), got %v", p.tok)
		}

		p.consumeTokenList(token.SEMICOLON)
	}
}
