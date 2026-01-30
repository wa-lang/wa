// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package gparser

import (
	"wa-lang.org/wa/internal/native/token"
)

func (p *parser) parseFile_gasExtern() {
	p.acceptToken(token.GAS_EXTERN)
	externName := p.parseIdent()
	p.gasExtern[externName] = true
	p.consumeSemicolonList()
}

func (p *parser) parseFile_gasSet() {
	p.acceptToken(token.GAS_SET)
	srcName := p.parseIdent()
	p.acceptToken(token.COMMA)
	dstName := p.parseIdent()
	p.gasAliases[srcName] = dstName
	p.consumeSemicolonList()
}

func (p *parser) parseFile_gasSection() {
	p.acceptToken(token.GAS_SECTION)
	p.gasAlign = 0

	sectionName := p.parseIdent()
	switch sectionName {
	case ".data", ".radata", ".bss":
		p.consumeSemicolonList()
		p.parseFile_gasGlobalList(sectionName)
	case ".text", ".init", ".fini":
		p.consumeSemicolonList()
		p.parseFile_gasFuncList(sectionName)
	default:
		p.errorf(p.pos, "invalid section name: %s", sectionName)
	}
}

func (p *parser) parseFile_gasGlobalList(sectionName string) {
	p.acceptToken(token.GAS_ALIGN)
	p.gasAlign = p.parseIntLit()
	p.consumeSemicolonList()

	// for

	panic("TODO: parseFile_gasGlobal")
}

func (p *parser) parseFile_gasFuncList(sectionName string) {
	// for
	panic("TODO: parseFile_gasFunc")
}
