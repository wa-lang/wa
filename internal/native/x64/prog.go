// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/x64/p9x86"
)

// 底层的指令别名
type Prog = p9x86.Prog

// 初始化b编码表格
func init() {
	p9x86.Init()
}

// 从解析器获得的参数构建底层的指令
// 如果有标识符, 默认均以转化为了具体的数值
func BuildProg(as abi.As, arg *abi.X64Argument) (inst *Prog, err error) {
	prog := &p9x86.Prog{}
	prog.To = operand2P9Addr(arg.Dst)
	prog.From = operand2P9Addr(arg.Src)

	dstXlen := operandXlen(arg.Dst)
	srcXlen := operandXlen(arg.Src)

	switch as {
	case AADD: // add
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.AADDB
		case 2:
			prog.As = p9x86.AADDW
		case 4:
			prog.As = p9x86.AADDL
		case 8:
			prog.As = p9x86.AADDQ
		default:
			panic("unreachable")
		}
	case AADDSD: // addsd
		prog.As = p9x86.AADDSD
	case AADDSS: // addss
		prog.As = p9x86.AADDSS
	case AAND: // and
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.AANDB
		case 2:
			prog.As = p9x86.AANDW
		case 4:
			prog.As = p9x86.AANDL
		case 8:
			prog.As = p9x86.AANDQ
		default:
			panic("unreachable")
		}
	case ACALL: // call
		prog.As = p9x86.ACALL
	case ACDQ: // cdq
		prog.As = p9x86.ACDQ
	case ACMP: // cmp
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.ACMPB
		case 2:
			prog.As = p9x86.ACMPW
		case 4:
			prog.As = p9x86.ACMPL
		case 8:
			prog.As = p9x86.ACMPQ
		default:
			panic("unreachable")
		}
	case ACVTSI2SD: // cvtsi2sd
		panic("TODO")
	case ACVTSS2SD: // cvtss2sd
		panic("TODO")
	case ACVTTSD2SI: // cvttsd2si
		panic("TODO")
	case ADEC: // dec
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.ADECB
		case 2:
			prog.As = p9x86.ADECW
		case 4:
			prog.As = p9x86.ADECL
		case 8:
			prog.As = p9x86.ADECQ
		default:
			panic("unreachable")
		}
	case ADIV: // div
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.ADIVB
		case 2:
			prog.As = p9x86.ADIVW
		case 4:
			prog.As = p9x86.ADIVL
		case 8:
			prog.As = p9x86.ADIVQ
		default:
			panic("unreachable")
		}
	case ADIVSD: // divsd
		prog.As = p9x86.ADIVSD
	case ADIVSS: // divss
		prog.As = p9x86.ADIVSS
	case AIDIV: // idiv
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.AIDIVB
		case 2:
			prog.As = p9x86.AIDIVW
		case 4:
			prog.As = p9x86.AIDIVL
		case 8:
			prog.As = p9x86.AIDIVQ
		default:
			panic("unreachable")
		}
	case AIMUL: // imul
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.AIMULB
		case 2:
			prog.As = p9x86.AIMULW
		case 4:
			prog.As = p9x86.AIMULL
		case 8:
			prog.As = p9x86.AIMULQ
		default:
			panic("unreachable")
		}
	case AINC: // inc
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.AINCB
		case 2:
			prog.As = p9x86.AINCW
		case 4:
			prog.As = p9x86.AINCL
		case 8:
			prog.As = p9x86.AINCQ
		default:
			panic("unreachable")
		}
	case AJA: // ja
		prog.As = p9x86.AJHI
	case AJB: // jb
		prog.As = p9x86.AJCS
	case AJE: // je
		prog.As = p9x86.AJEQ
	case AJGE: // jge
		prog.As = p9x86.AJEQ
	case AJMP: // jmp
		prog.As = p9x86.AJMP
	case AJNS: // jns
		prog.As = p9x86.AJNE
	case AJNZ: // jnz
		prog.As = p9x86.AJNE
	case AJZ: // jz
		prog.As = p9x86.AJEQ
	case ALEA: // lea
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			panic("unreachable")
		case 2:
			prog.As = p9x86.ALEAW
		case 4:
			prog.As = p9x86.ALEAL
		case 8:
			prog.As = p9x86.ALEAQ
		default:
			panic("unreachable")
		}
	case AMOV: // mov
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.AMOVB
		case 2:
			prog.As = p9x86.AMOVW
		case 4:
			prog.As = p9x86.AMOVL
		case 8:
			prog.As = p9x86.AMOVQ
		default:
			panic("unreachable")
		}
	case AMOVABS: // movabs
		prog.As = p9x86.AMOVQ
	case AMOVQ: // movq
		prog.As = p9x86.AMOVQ
	case AMOVSD: // movsd
		prog.As = p9x86.AMOVSD
	case AMOVSS: // movss
		prog.As = p9x86.AMOVSS
	case AMOVSXD: // movsxd
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.AMOVBLSX
		case 2:
			prog.As = p9x86.AMOVBWSX
		case 4:
			prog.As = p9x86.AMOVBLSX
		case 8:
			prog.As = p9x86.AMOVQ
		default:
			panic("unreachable")
		}
	case AMOVZX: // movzx
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.AMOVBLZX
		case 2:
			prog.As = p9x86.AMOVBWZX
		case 4:
			prog.As = p9x86.AMOVBLZX
		case 8:
			prog.As = p9x86.AMOVQ
		default:
			panic("unreachable")
		}
	case AMULSD: // mulsd
		prog.As = p9x86.AMULSD
	case AMULSS: // mulss
		prog.As = p9x86.AMULSS
	case ANEG: // neg
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.ANEGB
		case 2:
			prog.As = p9x86.ANEGW
		case 4:
			prog.As = p9x86.ANEGL
		case 8:
			prog.As = p9x86.ANEGQ
		default:
			panic("unreachable")
		}
	case ANOP: // nop
		prog.As = p9x86.ANOPW
	case AOR: // or
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.AORB
		case 2:
			prog.As = p9x86.AORW
		case 4:
			prog.As = p9x86.AORL
		case 8:
			prog.As = p9x86.AORQ
		default:
			panic("unreachable")
		}
	case APOP: // pop
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			panic("TODO")
		case 2:
			prog.As = p9x86.APOPW
		case 4:
			prog.As = p9x86.APOPL
		case 8:
			prog.As = p9x86.APOPQ
		default:
			panic("unreachable")
		}
	case APUSH: // push
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			panic("TODO")
		case 2:
			prog.As = p9x86.APUSHW
		case 4:
			prog.As = p9x86.APUSHL
		case 8:
			prog.As = p9x86.APUSHQ
		default:
			panic("unreachable")
		}
	case ARET: // ret
		prog.As = p9x86.ARET
	case ASAR: // sar
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.ASARB
		case 2:
			prog.As = p9x86.ASARW
		case 4:
			prog.As = p9x86.ASARL
		case 8:
			prog.As = p9x86.ASARQ
		default:
			panic("unreachable")
		}
	case ASETA: // seta
		panic("TODO")
	case ASETAE: // setae
		panic("TODO")
	case ASETB: // setb
		panic("TODO")
	case ASETBE: // setbe
		panic("TODO")
	case ASETE: // sete
		panic("TODO")
	case ASETG: // setg
		panic("TODO")
	case ASETGE: // setge
		panic("TODO")
	case ASETL: // setl
		panic("TODO")
	case ASETLE: // setle
		panic("TODO")
	case ASETNE: // setne
		panic("TODO")
	case ASETNP: // setnp
		panic("TODO")
	case ASHL: // shl
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.ASHLB
		case 2:
			prog.As = p9x86.ASHLW
		case 4:
			prog.As = p9x86.ASHLL
		case 8:
			prog.As = p9x86.ASHLQ
		default:
			panic("unreachable")
		}
	case ASTD: // std
		prog.As = p9x86.ASTD
	case ASUB: // sub
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.ASUBB
		case 2:
			prog.As = p9x86.ASUBW
		case 4:
			prog.As = p9x86.ASUBL
		case 8:
			prog.As = p9x86.ASUBQ
		default:
			panic("unreachable")
		}
	case ASUBSD: // subsd
		prog.As = p9x86.ASUBSD
	case ASUBSS: // subss
		prog.As = p9x86.ASUBSS
	case ASYSCALL: // syscall
		prog.As = p9x86.ASYSCALL
	case ATEST: // test
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.ATESTB
		case 2:
			prog.As = p9x86.ATESTW
		case 4:
			prog.As = p9x86.ATESTL
		case 8:
			prog.As = p9x86.ATESTQ
		default:
			panic("unreachable")
		}
	case AUCOMISD: // ucomisd
		prog.As = p9x86.AUCOMISD
	case AXOR: // xor
		switch maxInt(dstXlen, srcXlen) {
		case 1:
			prog.As = p9x86.AXORB
		case 2:
			prog.As = p9x86.AXORW
		case 4:
			prog.As = p9x86.AXORL
		case 8:
			prog.As = p9x86.AXORQ
		default:
			panic("unreachable")
		}
	default:
		panic(fmt.Sprintf("TODO: %v", as))
	}

	return prog, nil
}

func reg2p9Reg(r abi.RegType) int16 {
	switch {
	case r >= REG_EAX && r <= REG_R15D:
		return p9x86.REG_AX + int16(r-REG_EAX)
	case r >= REG_RAX && r <= REG_R15:
		return p9x86.REG_AX + int16(r-REG_RAX)
	default:
		panic("unreachable")
	}
}

func operandXlen(op *abi.X64Operand) int {
	if op == nil {
		return 1
	}
	switch op.Kind {
	case abi.X64Operand_Reg:
		return RegXLen(op.Reg)
	case abi.X64Operand_Mem:
		switch op.PtrTyp {
		case abi.X64BytePtr:
			return 1
		case abi.X64WordPtr:
			return 2
		case abi.X64DWordPtr:
			return 4
		case abi.X64QWordPtr:
			return 8
		default:
			panic("unreachable")
		}
	case abi.X64Operand_Imm:
		return 1 // 立即数按照单字节计算, 后续再修复
	default:
		panic("unreachable")
	}
}

func operand2P9Addr(op *abi.X64Operand) p9x86.Addr {
	if op == nil {
		return p9x86.Addr{Type: p9x86.TYPE_NONE}
	}

	addr := p9x86.Addr{}

	switch op.Kind {
	case abi.X64Operand_Reg:
		addr.Type = p9x86.TYPE_REG
		addr.Reg = reg2p9Reg(op.Reg)

	case abi.X64Operand_Imm:
		addr.Type = p9x86.TYPE_CONST
		addr.Offset = op.Imm

	case abi.X64Operand_Mem:
		addr.Type = p9x86.TYPE_MEM // lea/jmp 等需要再次修复
		addr.Reg = reg2p9Reg(op.Reg)
		addr.Offset = op.Offset
	}

	return addr
}
