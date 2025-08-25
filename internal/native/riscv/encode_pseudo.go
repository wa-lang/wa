// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "wa-lang.org/wa/internal/native/abi"

// 伪指令只产生1个机器指令
func (ctx *_OpContextType) encodePseudo(xlen int, as abi.As, arg *abi.AsArgument) (uint32, error) {
	if ctx.PseudoAs == 0 {
		panic("unreachable")
	}
	switch as {
	default:
		panic("unreachable")
	case ANOP:
		if arg.Rd != 0 {

		}
		if arg.Imm != 0 {
			panic("encodeR: imm was nonzero")
		}
		return ctx.encodeRaw(xlen, AADDI, &abi.AsArgument{
			Rd:  REG_X0,
			Rs1: REG_X0,
			Imm: 0,
		})
	case AMV:
	case ANOT:
	case ANEG:
	case ANEGW:
	case ASEXT_W:
	case ASEQZ:
	case ASNEZ:
	case ASLTZ:
	case ASGTZ:
	case AFMV_S:
	case AFABS_S:
	case AFNEG_S:
	case AFMV_D:
	case AFABS_D:
	case AFNEG_D:
	case ABEQZ:
	case ABNEZ:
	case ABLEZ:
	case ABGEZ:
	case ABLTZ:
	case ABGTZ:
	case ABGT:
	case ABLE:
	case ABGTU:
	case ABLEU:
	case AJ:
	case AJR:
	case ARET:
	case ARDINSTRET:
	case ARDCYCLE:
	case ARDTIME:
	case ACSRR:
	case ACSRW:
	case ACSRS:
	case ACSRC:
	case ACSRWI:
	case ACSRSI:
	case ACSRCI:
	case AFRCSR:
	case AFSCSR:
	case AFRRM:
	case AFSRM:
	case AFRFLAGS:
	case AFSFLAGS:
	}
	// TODO: 展开伪代码 后 继续调用 ctx.encode() 处理
	panic("TODO")
}

// 其他涉及地址的指令
// 超出范围的跳转需要拆分为2个指令
// if offset < -(1<<20) || (1<<20) <= offset {}
