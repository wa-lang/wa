// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import "wa-lang.org/wa/internal/native/token"

func (p *parser) parseDebugInfo_file() {
	p.acceptToken(token.GAS_DEBUG_FILE)
	_ = p.parseIntLit()   // 文件索引
	_ = p.parseBasicLit() // 文件名
	p.consumeSemicolonList()
}

func (p *parser) parseDebugInfo_loc() {
	p.acceptToken(token.GAS_DEBUG_LOC)
	_ = p.parseIntLit() // 文件索引
	_ = p.parseIntLit() // 行号
	_ = p.parseIntLit() // 列号
	p.consumeSemicolonList()
}

func (p *parser) parseDebugInfo_size() {
	p.acceptToken(token.GAS_DEBUG_SIZE)

	// .size main, .-main
	_ = p.parseIdent()         // 函数名
	p.acceptToken(token.COMMA) // ,
	_ = p.parseIdent()         // .
	p.acceptToken(token.SUB)   // -
	_ = p.parseIdent()         // 函数名
	p.consumeSemicolonList()
}

func (p *parser) parseDebugInfo_type() {
	p.acceptToken(token.GAS_DEBUG_TYPE)
	_ = p.parseIdent()         // 函数或全局变量名
	p.acceptToken(token.COMMA) // ,
	_ = p.parseIdent()         // @function/@object
	p.consumeSemicolonList()
}
