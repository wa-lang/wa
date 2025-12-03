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
	case loong64.OpFormatType_NULL:
		return inst
	case loong64.OpFormatType_2R:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case loong64.OpFormatType_2F:
		panic("TODO")
	case loong64.OpFormatType_1F_1R:
		panic("TODO")
	case loong64.OpFormatType_1R_1F:
		panic("TODO")
	case loong64.OpFormatType_3R:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case loong64.OpFormatType_3F:
		panic("TODO")
	case loong64.OpFormatType_1F_2R:
		panic("TODO")
	case loong64.OpFormatType_4F:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs3 = p.parseRegister()
		return inst

	case loong64.OpFormatType_2R_ui5:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		return inst
	case loong64.OpFormatType_2R_ui6:
		panic("TODO")
	case loong64.OpFormatType_2R_si12:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		return inst
	case loong64.OpFormatType_2R_ui12:
		panic("TODO")
	case loong64.OpFormatType_2R_si14:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		return inst
	case loong64.OpFormatType_2R_si16:
		panic("TODO")
	case loong64.OpFormatType_1R_si20:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_immAddr(inst)
		return inst

	case loong64.OpFormatType_0_2R:
		// 没有Rd
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst

	case loong64.OpFormatType_3R_s2:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		return inst
	case loong64.OpFormatType_3R_s3:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		return inst

	case loong64.OpFormatType_code:
		p.parseInst_loong_imm(inst)
		return inst

	case loong64.OpFormatType_code_1R_si12:
		// TODO: 可能需要通过名字解析code
		code := p.parseInt32Lit()
		inst.Arg.Rd = abi.RegType(code) + loong64.REG_R0
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		return inst

	case loong64.OpFormatType_msbw_lsbw:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		p.acceptToken(token.COMMA)
		inst.Arg.Imm2 = p.parseInt32Lit()
		return inst
	case loong64.OpFormatType_msbd_lsbd:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		p.parseInst_loong_imm(inst)
		p.acceptToken(token.COMMA)
		inst.Arg.Imm2 = p.parseInt32Lit()
		return inst

	case loong64.OpFormatType_fcsr_1R:
		// TODO: 解析 fcsr
		inst.Arg.Imm = p.parseInt32Lit()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case loong64.OpFormatType_1R_fcsr:
		// TODO: 解析 fcsr
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst

	case loong64.OpFormatType_cd_1R:
		inst.Arg.Imm = p.parseInt32Lit()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case loong64.OpFormatType_cd_1F:
		panic("TODO")
	case loong64.OpFormatType_cd_2R:
		inst.Arg.Imm = p.parseInt32Lit()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case loong64.OpFormatType_cd_2F:
		panic("TODO")
	case loong64.OpFormatType_1R_cj:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case loong64.OpFormatType_1F_cj:
		panic("TODO")
	case loong64.OpFormatType_1R_csr:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case loong64.OpFormatType_2R_csr:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case loong64.OpFormatType_2R_level:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case loong64.OpFormatType_level:
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case loong64.OpFormatType_0_1R_seq:
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case loong64.OpFormatType_op_2R:
		inst.Arg.Imm = p.parseInt32Lit()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case loong64.OpFormatType_3R_ca:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst

	case loong64.OpFormatType_hint_1R_si12:
		inst.Arg.Imm = p.parseInt32Lit()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		return inst
	case loong64.OpFormatType_hint_2R:
		inst.Arg.Imm = p.parseInt32Lit()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs2 = p.parseRegister()
		return inst
	case loong64.OpFormatType_hint:
		inst.Arg.Imm = p.parseInt32Lit()
		return inst

	case loong64.OpFormatType_cj_offset:
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case loong64.OpFormatType_rj_offset:
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case loong64.OpFormatType_rj_rd_offset:
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case loong64.OpFormatType_rd_rj_offset:
		inst.Arg.Rd = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Rs1 = p.parseRegister()
		p.acceptToken(token.COMMA)
		inst.Arg.Imm = p.parseInt32Lit()
		return inst
	case loong64.OpFormatType_offset:
		inst.Arg.Imm = p.parseInt32Lit()
		return inst

	default:
		panic("unreachable")
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
