// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// 解析值
func (p *parser) parseLitValue() *ast.BasicLit {
	pVal := &ast.BasicLit{Pos: p.pos}

	// 带类型的常量
	// const $D = i64('D') # i64
	// 常量 $甲 = 整64（‘A’） # i64

	switch p.tok {
	case token.I32, token.I32_zh,
		token.I64, token.I64_zh,
		token.U32, token.U32_zh,
		token.U64, token.U64_zh,
		token.F32, token.F32_zh,
		token.F64, token.F64_zh,
		token.PTR, token.PTR_zh:
		pVal.TypeCast = p.tok
		p.acceptToken(token.LPAREN)
		defer p.acceptToken(token.RPAREN)
	}

	// 默认类型
	switch p.tok {
	case token.CHAR, token.INT, token.FLOAT:
		pVal.Kind = p.tok
		if pVal.TypeCast == token.NONE {
			pVal.TypeCast = p.tok.DefaultNumberType()
		}
		pVal.Value = p.lit
		return pVal
	case token.STRING:
		pVal.Kind = p.tok
		pVal.Value = p.lit
		return pVal
	default:
		p.errorf(p.pos, "expect type or lit, got %v", p.tok)
	}

	return pVal
}
