// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package gparser

import (
	"fmt"
	"math/big"
	"strconv"

	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// 解析值
func (p *parser) parseBasicLit() *ast.BasicLit {
	pVal := &ast.BasicLit{Pos: p.pos}

	// 带类型的常量
	// const $D = i64('D') # i64
	// 常量 $甲 = 整64（‘A’） # i64

	switch p.tok {
	case token.GAS_BYTE,
		token.GAS_SHORT,
		token.GAS_LONG, token.GAS_QUAD:
		p.acceptToken(token.LPAREN)
		defer p.acceptToken(token.RPAREN)
	}

	// 默认类型
	switch p.tok {
	case token.CHAR, token.INT, token.FLOAT:
		pVal.LitKind = p.tok
		pVal.LitString = p.lit
		p.acceptToken(p.tok)
	case token.STRING:
		pVal.LitKind = p.tok
		pVal.LitString = p.lit
		p.acceptToken(p.tok)
	default:
		p.errorf(p.pos, "expect type or lit, got %v", p.tok)
	}

	// 解析常量的值
	switch pVal.LitKind {
	case token.INT:
		if x, err := strconv.ParseInt(pVal.LitString, 0, 64); err == nil {
			pVal.ConstV = int64(x)
		} else {
			p.errorf(pVal.Pos, "int %v %v", pVal.LitString, err)
		}
	case token.FLOAT:
		if f, ok := new(big.Float).SetString(pVal.LitString); ok {
			if f64V, acc := f.Float64(); acc != big.Exact {
				pVal.ConstV = float64(f64V)
			} else {
				p.errorf(pVal.Pos, "float %v %v", pVal.LitString, acc)
			}
		} else {
			p.errorf(pVal.Pos, "expect float, got %v", pVal.LitString)
		}

	case token.CHAR:
		if n := len(pVal.LitString); n >= 2 {
			if code, _, _, err := strconv.UnquoteChar(pVal.LitString[1:n-1], '\''); err == nil {
				pVal.ConstV = int64(code)
			} else {
				p.errorf(pVal.Pos, "char %v %v", pVal.LitString, err)
			}
		}

	case token.STRING:
		if s, err := strconv.Unquote(pVal.LitString); err == nil {
			pVal.ConstV = s
		} else {
			p.errorf(pVal.Pos, "string %q %v", pVal.LitString, err)
		}

	default:
		panic(fmt.Sprintf("%v is not a valid token", pVal.LitKind))
	}

	return pVal
}
