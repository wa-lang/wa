// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// 解析值
func (p *parser) parseValue() *ast.Value {
	pVal := &ast.Value{Pos: p.pos}

	// 无类型的面值
	switch p.tok {
	case token.INT:
		pVal.IntValue = p.parseInt64Lit()
		return pVal
	case token.FLOAT:
		pVal.FloatValue = p.parseFloat64Lit()
		return pVal
	case token.STRING:
		pVal.StrValue = p.parseStringLit()
		return pVal
	case token.IDENT:
		pVal.Symbal = p.parseIdent()
		return pVal
	}

	// 带类型的常量
	// const $D = i64('D') # i64
	// 常量 $甲 = 整64（‘A’） # i64

	switch p.tok {
	case token.I32, token.I32_zh:
		pVal.TypeDecor = p.tok
		p.acceptToken(p.tok)
		if p.tok == token.IDENT {
			pVal.Symbal = p.parseIdent()
		} else {
			pVal.IntValue = int64(p.parseInt32Lit())
		}
	case token.I64, token.I64_zh:
		pVal.TypeDecor = p.tok
		p.acceptToken(p.tok)
		if p.tok == token.IDENT {
			pVal.Symbal = p.parseIdent()
		} else {
			pVal.IntValue = int64(p.parseInt64Lit())
		}
	case token.U32, token.U32_zh:
		pVal.TypeDecor = p.tok
		p.acceptToken(p.tok)
		if p.tok == token.IDENT {
			pVal.Symbal = p.parseIdent()
		} else {
			pVal.UintValue = uint64(p.parseUint32Lit())
		}
	case token.U64, token.U64_zh:
		pVal.TypeDecor = p.tok
		p.acceptToken(p.tok)
		if p.tok == token.IDENT {
			pVal.Symbal = p.parseIdent()
		} else {
			pVal.UintValue = uint64(p.parseUint64Lit())
		}
	case token.F32, token.F32_zh:
		pVal.TypeDecor = p.tok
		p.acceptToken(p.tok)
		if p.tok == token.IDENT {
			pVal.Symbal = p.parseIdent()
		} else {
			pVal.FloatValue = float64(p.parseFloat32Lit())
		}
	case token.F64, token.F64_zh:
		pVal.TypeDecor = p.tok
		p.acceptToken(p.tok)
		if p.tok == token.IDENT {
			pVal.Symbal = p.parseIdent()
		} else {
			pVal.FloatValue = float64(p.parseFloat64Lit())
		}
	case token.PTR, token.PTR_zh:
		pVal.TypeDecor = p.tok
		p.acceptToken(p.tok)
		if p.tok == token.IDENT {
			pVal.Symbal = p.parseIdent()
		} else {
			pVal.UintValue = p.parseUint64Lit()
		}
	default:
		p.errorf(p.pos, "expect type or lit, got %v", p.tok)
	}

	return pVal
}
