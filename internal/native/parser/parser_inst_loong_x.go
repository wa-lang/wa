// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/token"
)

func (p *parser) parseInst_loong(fn *ast.Func) (inst *ast.Instruction) {
	assert(p.cpu == abi.LOONG64)

	inst = new(ast.Instruction)
	inst.Arg = new(abi.AsArgument)

	inst.Doc = p.parseDocComment(&fn.Body.Comments, inst.Pos)
	if inst.Doc != nil {
		fn.Body.Objects = fn.Body.Objects[:len(fn.Body.Objects)-1]
	}

	defer func() {
		inst.Comment = p.parseTailComment(inst.Pos)
		p.consumeSemicolonList()
	}()

	if p.tok == token.IDENT {
		inst.Pos = p.pos
		inst.Label = p.parseIdent()
		p.acceptToken(token.COLON)

		// 后续如果不是指令则结束
		if !p.tok.IsAs() {
			return inst
		}
	}

	inst.Pos = p.pos
	inst.AsName = p.lit
	inst.As = p.parseAs()

	// 查询指令的参数格式
	if !loong64.AsValid(inst.As) {
		p.errorf(p.pos, "%v is not loong instruction", inst.As)
	}

	switch loong64.AsFormatType(inst.As) {
	case loong64.OpFormatType_2R:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case loong64.OpFormatType_3R:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case loong64.OpFormatType_4R:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs3 = p.parseRegister()
		return inst

	case loong64.OpFormatType_2RI8:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		return inst
	case loong64.OpFormatType_2RI12:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		return inst
	case loong64.OpFormatType_2RI14:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		return inst
	case loong64.OpFormatType_2RI16:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		return inst

	case loong64.OpFormatType_1RI20:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_immAddr(inst)
		return inst

	case loong64.OpFormatType_1RI21:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_immAddr(inst)
		return inst

	case loong64.OpFormatType_I26:
		p.parseInst_loong_immAddr(inst)
		return inst

	default:
		// TODO: 补齐不同类型的定义
		panic("TODO")
	}
}

func (p *parser) parseInst_loong_imm(inst *ast.Instruction) {
	switch p.tok {
	case token.IDENT:
		inst.Arg.Symbol = p.parseIdent()
		return
	case token.INT:
		inst.Arg.Imm = p.parseInt32Lit()
		return
	default:
		p.errorf(p.pos, "expect label or int, got %v", p.tok)
	}
	panic("unreachable")
}

// 解析地址立即数
// 不涉及寄存器解析
// 12
// _start
// %lo(_start)
func (p *parser) parseInst_loong_immAddr(inst *ast.Instruction) {
	if p.tok == token.INT {
		inst.Arg.Imm = p.parseInt32Lit()
		return
	}

	if p.tok != token.IDENT {
		p.errorf(p.pos, "export IDENT, got %v", p.tok)
	}

	pos := p.pos
	symbolOrDecor := p.parseIdent()

	// 没有重定位修饰函数
	if p.tok != token.LPAREN {
		inst.Arg.Symbol = symbolOrDecor
		return
	}

	// 判断重定位修饰函数
	switch symbolOrDecor {
	case "%hi":
		inst.Arg.SymbolDecor = abi.BuiltinFn_HI
	case "%lo":
		inst.Arg.SymbolDecor = abi.BuiltinFn_LO
	case "%pcrel_hi":
		inst.Arg.SymbolDecor = abi.BuiltinFn_PCREL_HI
	case "%pcrel_lo":
		inst.Arg.SymbolDecor = abi.BuiltinFn_PCREL_LO

	case "%高位":
		inst.Arg.SymbolDecor = abi.BuiltinFn_HI_zh
	case "%低位":
		inst.Arg.SymbolDecor = abi.BuiltinFn_LO_zh
	case "%相对高位":
		inst.Arg.SymbolDecor = abi.BuiltinFn_PCREL_HI_zh
	case "%相对低位":
		inst.Arg.SymbolDecor = abi.BuiltinFn_PCREL_LO_zh
	default:
		p.errorf(pos, "unknow symbol decorator %s", symbolOrDecor)
	}

	p.acceptToken(token.LPAREN)
	inst.Arg.Symbol = p.parseIdent()
	p.acceptToken(token.RPAREN)
}
