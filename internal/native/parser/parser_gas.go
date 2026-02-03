// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// 外部符号
func (p *parser) parseFile_gasExtern() {
	p.acceptToken(token.GAS_EXTERN)
	externName := p.parseIdent()
	p.gasExtern[externName] = true
	p.prog.Externs = append(p.prog.Externs, externName)
	p.consumeSemicolonList()
}

// 导出符号
func (p *parser) parseFile_gasGlobal() {
	p.acceptToken(token.GAS_GLOBL)
	ident := p.parseIdent()
	p.gasGlobal[ident] = true
	// 对象可能在后面定义, 最后再收集
	p.consumeSemicolonList()
}

// 段定义开始
func (p *parser) parseFile_gasSection() {
	p.acceptToken(token.GAS_SECTION)
	p.gasAlign = 0

	sectionName := p.parseIdent()
	switch sectionName {
	case ".data", ".radata", ".bss":
		p.consumeSemicolonList()
		p.parseFile_gasGlobalDefineList(sectionName)
	case ".text", ".init", ".fini":
		p.consumeSemicolonList()
		p.parseFile_gasFuncList(sectionName)
	default:
		p.errorf(p.pos, "invalid section name: %s", sectionName)
	}
}

// 全局变量定义列表
func (p *parser) parseFile_gasGlobalDefineList(sectionName string) {
	p.acceptToken(token.GAS_ALIGN)
	p.gasAlign = p.parseIntLit()
	p.consumeSemicolonList()

	for {
		if p.err != nil {
			return
		}
		if p.tok == token.EOF {
			return
		}
		if p.tok == token.GAS_SECTION {
			return
		}

		switch p.tok {
		case token.COMMENT:
			commentObj := p.parseCommentGroup(true)
			p.prog.Comments = append(p.prog.Comments, commentObj)
			p.prog.Objects = append(p.prog.Objects, commentObj)

		case token.GAS_GLOBL:
			p.parseFile_gasGlobal()
		case token.IDENT:
			p.prog.Globals = append(p.prog.Globals, p.parseFile_gasGlobalDefine(sectionName))

		default:
			p.errorf(p.pos, "unkonw token: %v", p.tok)
		}
	}
}

func (p *parser) parseFile_gasGlobalDefine(sectionName string) *ast.Global {
	g := &ast.Global{
		Section: sectionName,
		Name:    p.parseIdent(),
	}

	p.acceptToken(token.COLON)

	g.Tok = p.tok

	switch p.tok {
	case token.GAS_BYTE:
		g.Type = token.I8
		p.next()
	case token.GAS_SHORT:
		g.Type = token.I16
		p.next()
	case token.GAS_LONG:
		g.Type = token.I32
		p.next()
	case token.GAS_QUAD:
		g.Type = token.I64
		p.next()
	case token.GAS_ASSCII:
		g.Type = 0
		p.next()
	case token.GAS_SKIP:
		g.Type = 0
		p.next()
	case token.GAS_INCBIN:
		g.Type = 0
		p.next()
	default:
		p.errorf(p.pos, "unkonw token: %v", p.tok)
	}

	g.Init = &ast.InitValue{
		Lit: p.parseBasicLit(),
	}

	p.consumeSemicolonList()

	return g
}

func (p *parser) parseFile_gasFuncList(sectionName string) {
	for {
		if p.err != nil {
			return
		}
		if p.tok == token.EOF {
			return
		}

		switch p.tok {
		case token.COMMENT:
			commentObj := p.parseCommentGroup(true)
			p.prog.Comments = append(p.prog.Comments, commentObj)
			p.prog.Objects = append(p.prog.Objects, commentObj)

		case token.GAS_GLOBL:
			p.parseFile_gasGlobal()
		case token.IDENT:
			p.prog.Funcs = append(p.prog.Funcs, p.parseFile_gasFuncDefine(sectionName))

		default:
			return
		}
	}
}

func (p *parser) parseFile_gasFuncDefine(sectionName string) *ast.Func {
	fn := &ast.Func{
		Section: sectionName,
		Pos:     p.pos,
		Tok:     p.tok,
		Type:    new(ast.FuncType),
		Body:    new(ast.FuncBody),
	}

	p.parseFile_gasFuncBody(fn)
	p.consumeSemicolonList()

	return fn
}

func (p *parser) parseFile_gasFuncBody(fn *ast.Func) {
	assert(p.cpu == abi.RISCV64 || p.cpu == abi.RISCV32 || p.cpu == abi.LOONG64)

	fn.Body.Pos = p.pos

	for {
		if p.err != nil {
			return
		}
		if p.tok == token.EOF {
			return
		}

		// 注释
		if p.tok == token.COMMENT {
			commentObj := p.parseCommentGroup(false)
			fn.Body.Comments = append(fn.Body.Comments, commentObj)
			fn.Body.Objects = append(fn.Body.Objects, commentObj)
			continue
		}

		// 解析指令
		if p.tok == token.IDENT || p.tok.IsAs() {
			inst := p.parseInst(fn)
			fn.Body.Insts = append(fn.Body.Insts, inst)
			fn.Body.Objects = append(fn.Body.Objects, inst)
			continue
		}

		// 未知 token
		return
	}
}
