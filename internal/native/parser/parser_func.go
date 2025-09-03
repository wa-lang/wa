// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
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
	fn := new(ast.Func)

	p.acceptToken(token.FUNC)
	fn.Name = p.parseIdent()

	p.parseFunc_args(fn)
	p.parseFunc_return(fn)
	p.parseFunc_body(fn)

	return fn
}

func (p *parser) parseFunc_args(fn *ast.Func) {
	_ = fn
	panic("TODO")
}

func (p *parser) parseFunc_return(fn *ast.Func) {
	_ = fn
	panic("TODO")
}

func (p *parser) parseFunc_body(fn *ast.Func) {
	for p.tok == token.LOCAL {
		p.parseFunc_body_local(fn)
	}
	for p.tok.IsAs() {
		p.parseFunc_body_inst(fn)
	}
}

func (p *parser) parseFunc_body_local(fn *ast.Func) {
	_ = fn
	p.next()
	panic("TODO")
}

func (p *parser) parseFunc_body_inst(fn *ast.Func) {
	_ = fn
	p.next()
	panic("TODO")
}
