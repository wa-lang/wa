// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// func $add {
//     # 指令
// Loop:
// }

// 函数 $add:
//     # 指令
// Loop:
// 完毕

func (p *parser) parseFunc(tok token.Token) *ast.Func {
	fn := &ast.Func{
		Pos:  p.pos,
		Tok:  tok,
		Body: new(ast.FuncBody),
	}

	fn.Doc = p.parseDocComment(&p.prog.Comments, fn.Pos)
	if fn.Doc != nil {
		p.prog.Objects = p.prog.Objects[:len(p.prog.Objects)-1]
	}

	fn.Tok = p.acceptTokenAorB(token.FUNC, token.FUNC_zh)
	fn.Name = p.parseIdent()

	p.parseFunc_body(fn)
	p.consumeSemicolonList()

	return fn
}

func (p *parser) parseFunc_body(fn *ast.Func) {
	assert(p.cpu == abi.RISCV64 || p.cpu == abi.RISCV32 || p.cpu == abi.LOONG64)

	fn.Body.Pos = p.pos

	if fn.Tok == token.FUNC_zh {
		p.acceptToken(token.COLON)
		defer p.acceptToken(token.END_zh)
	} else {
		p.acceptToken(token.LBRACE)
		defer p.acceptToken(token.RBRACE)
	}

Loop:
	for {
		switch p.tok {
		case token.EOF, token.ILLEGAL:
			break Loop
		case token.RBRACE, token.END_zh:
			break Loop
		case token.COMMENT:
			commentObj := p.parseCommentGroup(false)
			fn.Body.Comments = append(fn.Body.Comments, commentObj)
			fn.Body.Objects = append(fn.Body.Objects, commentObj)
		default:
			if p.tok == token.IDENT || p.tok.IsAs() {
				inst := p.parseInst(fn)
				fn.Body.Insts = append(fn.Body.Insts, inst)
				fn.Body.Objects = append(fn.Body.Objects, inst)
			} else {
				p.errorf(p.pos, "unknow as %v", p.tok)
			}
		}
	}
}
