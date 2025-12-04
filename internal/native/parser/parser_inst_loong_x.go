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
		p.acceptToken(token.COMMA)
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
		p.acceptToken(token.COMMA)
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
		p.acceptToken(token.COMMA)
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
		p.acceptToken(token.COMMA)
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
		p.acceptToken(token.COMMA)
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
		p.acceptToken(token.COMMA)
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
		p.acceptToken(token.COMMA)
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
		p.acceptToken(token.COMMA)
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
		p.acceptToken(token.COMMA)
		sa3, sa3Symbol := p.parseInst_loong_imm_sa3()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Rs2 = rk
		inst.Arg.Imm = sa3
		inst.Arg.Symbol = sa3Symbol
		return inst
	case loong64.OpFormatType_code:
		code, codeSymbol := p.parseInst_loong_imm_code_15bit()
		inst.Arg.Imm = code
		inst.Arg.Symbol = codeSymbol
		return inst

	case loong64.OpFormatType_code_1R_si12:
		code, codeSymbol := p.parseInst_loong_imm_code_5bit()
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
		p.checkRegI_loong(rd, rj)
		p.acceptToken(token.COMMA)
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
		p.checkRegI_loong(rd, rj)
		p.acceptToken(token.COMMA)
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
		fcsr, fcsrSymbol := p.parseInst_loong_imm_fcsr_5bit()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rj)
		inst.Arg.Rd = abi.RegType(fcsr) // Rd 寄存器参数位置用于记录 fscr
		inst.Arg.RdName = fcsrSymbol
		inst.Arg.Rs1 = rj
		return inst
	case loong64.OpFormatType_1R_fcsr:
		rd := p.parseRegister()
		p.checkRegI_loong(rd)
		p.acceptToken(token.COMMA)
		fcsr, fcsrSymbol := p.parseInst_loong_imm_fcsr_5bit()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = abi.RegType(fcsr) // Rs1 寄存器参数位置用于记录 fscr
		inst.Arg.Rs1Name = fcsrSymbol
		return inst

	case loong64.OpFormatType_cd_1R:
		cd, cdSymbol := p.parseInst_loong_imm_cd_2bit()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rj)
		p.acceptToken(token.COMMA)
		inst.Arg.Rd = abi.RegType(cd) // Rd 寄存器参数位置用于记录 cd
		inst.Arg.RdName = cdSymbol
		inst.Arg.Rs1 = rj
		return inst
	case loong64.OpFormatType_cd_1F:
		cd, cdSymbol := p.parseInst_loong_imm_cd_2bit()
		p.acceptToken(token.COMMA)
		fj := p.parseRegister()
		p.checkRegF_loong(fj)
		p.acceptToken(token.COMMA)
		inst.Arg.Rd = abi.RegType(cd) // Rd 寄存器参数位置用于记录 cd
		inst.Arg.RdName = cdSymbol
		inst.Arg.Rs1 = fj
		return inst
	case loong64.OpFormatType_cd_2F:
		cd, cdSymbol := p.parseInst_loong_imm_cd_2bit()
		p.acceptToken(token.COMMA)
		fj := p.parseRegister()
		p.acceptToken(token.COMMA)
		fk := p.parseRegister()
		p.checkRegF_loong(fj, fk)
		p.acceptToken(token.COMMA)
		inst.Arg.Rd = abi.RegType(cd) // Rd 寄存器参数位置用于记录 cd
		inst.Arg.RdName = cdSymbol
		inst.Arg.Rs1 = fj
		inst.Arg.Rs2 = fk
		return inst
	case loong64.OpFormatType_1R_cj:
		rd := p.parseRegister()
		p.checkRegI_loong(rd)
		p.acceptToken(token.COMMA)
		cj, cjSymbol := p.parseInst_loong_imm_cj_3bit()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = abi.RegType(cj) // Rs1 寄存器参数位置用于记录 cj
		inst.Arg.Rs1Name = cjSymbol
		return inst
	case loong64.OpFormatType_1F_cj:
		fd := p.parseRegister()
		p.checkRegF_loong(fd)
		p.acceptToken(token.COMMA)
		cj, cjSymbol := p.parseInst_loong_imm_cj_3bit()
		inst.Arg.Rd = fd
		inst.Arg.Rs1 = abi.RegType(cj) // Rs1 寄存器参数位置用于记录 cj
		inst.Arg.Rs1Name = cjSymbol
		return inst
	case loong64.OpFormatType_1R_csr:
		rd := p.parseRegister()
		p.checkRegI_loong(rd)
		p.acceptToken(token.COMMA)
		csr, csrSymbol := p.parseInst_loong_imm_csr_14bit()
		inst.Arg.Rd = rd
		inst.Arg.Imm = csr // CSR 作为 imm 记录, 没有宏修饰
		inst.Arg.Symbol = csrSymbol
		return inst
	case loong64.OpFormatType_2R_csr:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rd, rj)
		p.acceptToken(token.COMMA)
		csr, csrSymbol := p.parseInst_loong_imm_csr_14bit()
		if rj == loong64.REG_R0 || rj != loong64.REG_R1 {
			p.errorf(p.pos, "%v: rj(%v) can notequal 0 or 1", inst.As, rj)
		}
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = csr // CSR 作为 imm 记录, 没有宏修饰
		inst.Arg.Symbol = csrSymbol
		return inst
	case loong64.OpFormatType_2R_level:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rd, rj)
		p.acceptToken(token.COMMA)
		level, levelSymbol := p.parseInst_loong_imm_level_8bit()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = level
		inst.Arg.Symbol = levelSymbol
		return inst
	case loong64.OpFormatType_level:
		level, levelSymbol := p.parseInst_loong_imm_level_15bit()
		inst.Arg.Imm = level
		inst.Arg.Symbol = levelSymbol
		return inst
	case loong64.OpFormatType_0_1R_seq:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		p.checkRegI_loong(rd)
		seq, seqSymbol := p.parseInst_loong_imm_seq_8bit()
		inst.Arg.Imm = seq
		inst.Arg.Symbol = seqSymbol
		return inst
	case loong64.OpFormatType_op_2R:
		op, opSymbol := p.parseInst_loong_imm_op_5bit()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.acceptToken(token.COMMA)
		rk := p.parseRegister()
		p.checkRegI_loong(rj, rk)
		inst.Arg.Rd = abi.RegType(op) // Rd 寄存器存放 op
		inst.Arg.RdName = opSymbol
		inst.Arg.Rs1 = rj
		inst.Arg.Rs2 = rk
		return inst
	case loong64.OpFormatType_3F_ca:
		fd := p.parseRegister()
		p.acceptToken(token.COMMA)
		fj := p.parseRegister()
		p.acceptToken(token.COMMA)
		fk := p.parseRegister()
		p.acceptToken(token.COMMA)
		p.checkRegF_loong(fd, fj, fk)
		ca, caSymbol := p.parseInst_loong_imm_ca_3bit()
		inst.Arg.Rd = fd
		inst.Arg.Rs1 = fj
		inst.Arg.Rs2 = fk
		inst.Arg.Imm = ca
		inst.Arg.Symbol = caSymbol
		return inst
	case loong64.OpFormatType_hint_1R_si12:
		hint, hintSymbol := p.parseInst_loong_imm_hint_5bit()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.checkRegI_loong(rj)
		p.acceptToken(token.COMMA)
		si12, si12Symbol, si12SymbolDecor := p.parseInst_loong_imm_si12()
		inst.Arg.Rd = abi.RegType(hint) // Rd 寄存器保存 hint
		inst.Arg.RdName = hintSymbol
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = si12
		inst.Arg.Symbol = si12Symbol
		inst.Arg.SymbolDecor = si12SymbolDecor
		return inst
	case loong64.OpFormatType_hint_2R:
		hint, hintSymbol := p.parseInst_loong_imm_hint_5bit()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.acceptToken(token.COMMA)
		rk := p.parseRegister()
		p.checkRegI_loong(rj, rk)
		p.acceptToken(token.COMMA)
		inst.Arg.Rd = abi.RegType(hint) // Rd 寄存器保存 hint
		inst.Arg.RdName = hintSymbol
		inst.Arg.Rs1 = rj
		inst.Arg.Rs2 = rk
		return inst
	case loong64.OpFormatType_hint:
		hint, hintSymbol := p.parseInst_loong_imm_hint_15bit()
		inst.Arg.Imm = hint
		inst.Arg.Symbol = hintSymbol
		return inst

	case loong64.OpFormatType_cj_offset:
		cj, cjSymbol := p.parseInst_loong_imm_cj_3bit()
		p.acceptToken(token.COMMA)
		off, offSymbol, offSymbolDecor := p.parseInst_loong_imm_offset_21bit()
		inst.Arg.Rs1 = abi.RegType(cj) // Rs1 寄存器保存 cj
		inst.Arg.Rs1Name = cjSymbol
		inst.Arg.Imm = off
		inst.Arg.Symbol = offSymbol
		inst.Arg.SymbolDecor = offSymbolDecor
		return inst
	case loong64.OpFormatType_rj_offset:
		rj := p.parseRegister()
		p.checkRegI_loong(rj)
		p.acceptToken(token.COMMA)
		off, offSymbol, offSymbolDecor := p.parseInst_loong_imm_offset_21bit()
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = off
		inst.Arg.Symbol = offSymbol
		inst.Arg.SymbolDecor = offSymbolDecor
		return inst
	case loong64.OpFormatType_rj_rd_offset:
		rj := p.parseRegister()
		p.acceptToken(token.COMMA)
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		p.checkRegI_loong(rj, rd)
		off, offSymbol, offSymbolDecor := p.parseInst_loong_imm_offset_21bit()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = off
		inst.Arg.Symbol = offSymbol
		inst.Arg.SymbolDecor = offSymbolDecor
		return inst
	case loong64.OpFormatType_rd_rj_offset:
		rd := p.parseRegister()
		p.acceptToken(token.COMMA)
		rj := p.parseRegister()
		p.acceptToken(token.COMMA)
		p.checkRegI_loong(rd, rj)
		off, offSymbol, offSymbolDecor := p.parseInst_loong_imm_offset_21bit()
		inst.Arg.Rd = rd
		inst.Arg.Rs1 = rj
		inst.Arg.Imm = off
		inst.Arg.Symbol = offSymbol
		inst.Arg.SymbolDecor = offSymbolDecor
		return inst
	case loong64.OpFormatType_offset:
		off, offSymbol, offSymbolDecor := p.parseInst_loong_imm_offset_21bit()
		inst.Arg.Imm = off
		inst.Arg.Symbol = offSymbol
		inst.Arg.SymbolDecor = offSymbolDecor
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
func (p *parser) parseInst_loong_imm_offset_21bit() (x int32, symbol string, symbolDecor abi.BuiltinFn) {
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

func (p *parser) parseInst_loong_imm_fcsr_5bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_cd_2bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_cj_3bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_csr_14bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_level_8bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_level_15bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_seq_8bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_op_5bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_ca_3bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_hint_5bit() (x int32, symbol string) {
	return p.parseInst_loong_imm_sa()
}
func (p *parser) parseInst_loong_imm_hint_15bit() (x int32, symbol string) {
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
