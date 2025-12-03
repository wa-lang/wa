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
		p.errorf(p.pos, "%v is not loongarch instruction", inst.As)
	}

	switch loong64.AsFormatType(inst.As) {
	case loong64.OpFormatType_NULL:
		return inst
	case loong64.OpFormatType_2R:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rd, rj)
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		return inst
	case loong64.OpFormatType_2F:
		fd := p.parseRegister()
		p.acceptToken(token.COMMA)
		fj := p.parseRegister()
		p.checkRegF_loong(fd, fj)
		inst.Arg.Rd = fd
		inst.Arg.Rs1 = fj
		return inst
	case loong64.OpFormatType_1F_1R:
		fd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegF_loong(fd)
		p.checkRegI_loong(rj)
		inst.Arg.Rd = fd
		inst.Arg.Rs1 = rj
		return inst
	case loong64.OpFormatType_1R_1F:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		fj := p.parseRegister()
		p.checkRegI_loong(rd)
		p.checkRegF_loong(fj)
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = fj
		return inst
	case loong64.OpFormatType_3R:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.acceptToken(token.COMMA)
		rk := p.parseRegister()
		p.checkRegI_loong(rd, rj, rk)
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Rs2 = rk
		return inst
	case loong64.OpFormatType_3F:
		fd := p.parseRegister()
		p.acceptToken(token.COMMA)
		fj := p.parseRegister()
		p.acceptToken(token.COMMA)
		fk := p.parseRegister()
		p.checkRegF_loong(fd, fj, fk)
		inst.Arg.Rd = fd
		inst.Arg.Rs1 = fj
		inst.Arg.Rs2 = fk
		return inst
	case loong64.OpFormatType_1F_2R:
		fd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.acceptToken(token.COMMA)
		rk := p.parseRegister()
		p.checkRegF_loong(fd)
		p.checkRegI_loong(rj, rk)
		inst.Arg.Rd = fd
		inst.Arg.Rs1 = rj
		inst.Arg.Rs2 = rk
		return inst
	case loong64.OpFormatType_4F:
		fd := p.parseRegister()
		p.acceptToken(token.COMMA)
		fj := p.parseRegister()
		p.acceptToken(token.COMMA)
		fk := p.parseRegister()
		p.acceptToken(token.COMMA)
		fa := p.parseRegister()
		p.checkRegF_loong(fd, fj, fk, fa)
		inst.Arg.Rd = fd
		inst.Arg.Rs1 = fj
		inst.Arg.Rs2 = fk
		inst.Arg.Rs3 = fa
		return inst

	case loong64.OpFormatType_2R_ui5:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rd, rj)
		ui5, ui5Symbol, ui5SymbolDecor := p.parseInst_loong_imm_ui5()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = ui5
		inst.Arg.Symbol = ui5Symbol
		inst.Arg.SymbolDecor = ui5SymbolDecor
		return inst
	case loong64.OpFormatType_2R_ui6:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rd, rj)
		ui6, ui6Symbol, ui6SymbolDecor := p.parseInst_loong_imm_ui6()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = ui6
		inst.Arg.Symbol = ui6Symbol
		inst.Arg.SymbolDecor = ui6SymbolDecor
		return inst
	case loong64.OpFormatType_2R_si12:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rd, rj)
		si12, si12Symbol, si12SymbolDecor := p.parseInst_loong_imm_si12()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = si12
		inst.Arg.Symbol = si12Symbol
		inst.Arg.SymbolDecor = si12SymbolDecor
		return inst

	case loong64.OpFormatType_2R_ui12:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rd, rj)
		ui12, ui12Symbol, ui12SymbolDecor := p.parseInst_loong_imm_ui12()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = ui12
		inst.Arg.Symbol = ui12Symbol
		inst.Arg.SymbolDecor = ui12SymbolDecor
		return inst
	case loong64.OpFormatType_2R_si14:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rd, rj)
		si14, si14Symbol, si14SymbolDecor := p.parseInst_loong_imm_si14()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = si14
		inst.Arg.Symbol = si14Symbol
		inst.Arg.SymbolDecor = si14SymbolDecor
		return inst
	case loong64.OpFormatType_2R_si16:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rd, rj)
		si16, si16Symbol, si16SymbolDecor := p.parseInst_loong_imm_si16()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = si16
		inst.Arg.Symbol = si16Symbol
		inst.Arg.SymbolDecor = si16SymbolDecor
		return inst
	case loong64.OpFormatType_1R_si20:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rd, rj)
		si20, si20Symbol, si20SymbolDecor := p.parseInst_loong_imm_si20()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = si20
		inst.Arg.Symbol = si20Symbol
		inst.Arg.SymbolDecor = si20SymbolDecor
		return inst
	case loong64.OpFormatType_0_2R:
		rj := p.parseRegister()
		p.acceptToken(token.COMMA)
		rk := p.parseRegister()
		p.checkRegI_loong(rj, rk)
		inst.Arg.Rs1 = rj
		inst.Arg.Rs2 = rk
		return inst
	case loong64.OpFormatType_3R_sa2:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.acceptToken(token.COMMA)
		rk := p.parseRegister()
		p.checkRegI_loong(rd, rj, rk)
		sa2, sa2Symbol := p.parseInst_loong_imm_sa2()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Rs2 = rk
		inst.Arg.Imm = sa2
		inst.Arg.Symbol = sa2Symbol
		return inst
	case loong64.OpFormatType_3R_sa3:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.acceptToken(token.COMMA)
		rk := p.parseRegister()
		p.checkRegI_loong(rd, rj, rk)
		sa3, sa3Symbol := p.parseInst_loong_imm_sa3()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Rs2 = rk
		inst.Arg.Imm = sa3
		inst.Arg.Symbol = sa3Symbol
		return inst
	case loong64.OpFormatType_code:
		code, codeSymbol := p.parseInst_loong_imm_code_5bit()
		inst.Arg.Imm = code
		inst.Arg.Symbol = codeSymbol
		return inst

	case loong64.OpFormatType_code_1R_si12:
		code, codeSymbol := p.parseInst_loong_imm_code_15bit()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rj)
		p.acceptToken(token.COMMA)
		si12, si12Symbol, si12SymbolDecor := p.parseInst_loong_imm_si12()
		inst.Arg.Rd = abi.RegType(code) // Rd 寄存器参数位置用于记录 code
		inst.Arg.RdName = codeSymbol
		inst.Arg.Imm = si12
		inst.Arg.Symbol = si12Symbol
		inst.Arg.SymbolDecor = si12SymbolDecor
		return inst

	case loong64.OpFormatType_2R_msbw_lsbw:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.acceptToken(token.COMMA)
		p.checkRegI_loong(rd, rj)
		msbw, msbwSymbol := p.parseInst_loong_imm_msbw_5bit()
		p.acceptToken(token.COMMA)
		lsbw, lsbwSymbol := p.parseInst_loong_imm_msbw_5bit()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Rs2 = abi.RegType(msbw) // Rs1 寄存器参数位置用于记录 msbw
		inst.Arg.Rs2Name = msbwSymbol
		inst.Arg.Rs2 = abi.RegType(lsbw) // Rs2 寄存器参数位置用于记录 lsbw
		inst.Arg.Rs2Name = lsbwSymbol
		return inst
	case loong64.OpFormatType_2R_msbd_lsbd:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.acceptToken(token.COMMA)
		p.checkRegI_loong(rd, rj)
		msbd, msbdSymbol := p.parseInst_loong_imm_msbd_6bit()
		p.acceptToken(token.COMMA)
		lsbd, lsbdSymbol := p.parseInst_loong_imm_msbd_6bit()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Rs2 = abi.RegType(msbd) // Rs1 寄存器参数位置用于记录 msbd
		inst.Arg.Rs2Name = msbdSymbol
		inst.Arg.Rs2 = abi.RegType(lsbd) // Rs2 寄存器参数位置用于记录 lsbd
		inst.Arg.Rs2Name = lsbdSymbol
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

// 检查普通寄存器类型
func (p *parser) checkRegI_loong(regs ...abi.RegType) {
	for _, x := range regs {
		if x < loong64.REG_R0 || x > loong64.REG_R31 {
			p.errorf(p.pos, "%v is not loongarch int register", x)
		}
	}
}

// 检查浮点数寄存器类型
func (p *parser) checkRegF_loong(regs ...abi.RegType) {
	for _, x := range regs {
		if x < loong64.REG_F0 || x > loong64.REG_F31 {
			p.errorf(p.pos, "%v is not loongarch float register", x)
		}
	}
}

func (p *parser) parseInst_loong_imm_ui5() (ui5 int32, symbol string, symbolDecor abi.BuiltinFn) {
	return p.parseInst_loong_imm_v2()
}
func (p *parser) parseInst_loong_imm_ui6() (ui6 int32, symbol string, symbolDecor abi.BuiltinFn) {
	return p.parseInst_loong_imm_v2()
}

func (p *parser) parseInst_loong_imm_si12() (s12 int32, symbol string, symbolDecor abi.BuiltinFn) {
	return p.parseInst_loong_imm_v2()
}
func (p *parser) parseInst_loong_imm_ui12() (ui12 int32, symbol string, symbolDecor abi.BuiltinFn) {
	return p.parseInst_loong_imm_v2()
}

func (p *parser) parseInst_loong_imm_si14() (si14 int32, symbol string, symbolDecor abi.BuiltinFn) {
	return p.parseInst_loong_imm_v2()
}
func (p *parser) parseInst_loong_imm_si16() (si16 int32, symbol string, symbolDecor abi.BuiltinFn) {
	return p.parseInst_loong_imm_v2()
}
func (p *parser) parseInst_loong_imm_si20() (si20 int32, symbol string, symbolDecor abi.BuiltinFn) {
	return p.parseInst_loong_imm_v2()
}

func (p *parser) parseInst_loong_imm_sa2() (sa2 int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}

func (p *parser) parseInst_loong_imm_sa3() (sa3 int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}

func (p *parser) parseInst_loong_imm_code_5bit() (code int32, symbol string) {
	return p.parseInst_loong_imm_sa() // todo: 检查整数的范围
}
func (p *parser) parseInst_loong_imm_code_15bit() (code int32, symbol string) {
	return p.parseInst_loong_imm_sa() // todo: 检查整数的范围
}

func (p *parser) parseInst_loong_imm_msbw_5bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_lsbw_5bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}

func (p *parser) parseInst_loong_imm_msbd_6bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_lsbd_6bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}

func (p *parser) parseInst_loong_imm_sa() (imm int32, symbol string) {
	switch p.tok {
	case token.INT:
		imm = p.parseInt32Lit()
		return
	case token.IDENT:
		symbol = p.parseIdent()
		return
	default:
		p.errorf(p.pos, "expect label or int, got %v", p.tok)
	}
	panic("unreachable")
}

func (p *parser) parseInst_loong_imm_v2() (imm int32, symbol string, symbolDecor abi.BuiltinFn) {
	if p.tok == token.INT {
		imm = p.parseInt32Lit()
		return
	}

	if p.tok != token.IDENT {
		p.errorf(p.pos, "export IDENT, got %v", p.tok)
	}

	pos := p.pos
	symbolOrDecor := p.parseIdent()

	// 没有重定位修饰函数
	if p.tok != token.LPAREN {
		symbol = symbolOrDecor
		return
	}

	// 判断重定位修饰函数
	switch symbolOrDecor {
	case "%hi":
		symbolDecor = abi.BuiltinFn_HI
	case "%lo":
		symbolDecor = abi.BuiltinFn_LO
	case "%pcrel_hi":
		symbolDecor = abi.BuiltinFn_PCREL_HI
	case "%pcrel_lo":
		symbolDecor = abi.BuiltinFn_PCREL_LO

	case "%高位":
		symbolDecor = abi.BuiltinFn_HI_zh
	case "%低位":
		symbolDecor = abi.BuiltinFn_LO_zh
	case "%相对高位":
		symbolDecor = abi.BuiltinFn_PCREL_HI_zh
	case "%相对低位":
		symbolDecor = abi.BuiltinFn_PCREL_LO_zh
	default:
		p.errorf(pos, "unknow symbol decorator %s", symbolOrDecor)
	}

	p.acceptToken(token.LPAREN)
	symbol = p.parseIdent()
	p.acceptToken(token.RPAREN)

	switch p.tok {
	case token.INT:
		imm = p.parseInt32Lit()
		return
	case token.IDENT:
		symbol = p.parseIdent()
		return
	default:
		p.errorf(p.pos, "expect label or int, got %v", p.tok)
	}
	panic("unreachable")
}

func (p *parser) parseInst_loong_imm(inst *ast.Instruction) {
	switch p.tok {
	case token.INT:
		inst.Arg.Imm = p.parseInt32Lit()
		return
	case token.IDENT:
		inst.Arg.Symbol = p.parseIdent()
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
