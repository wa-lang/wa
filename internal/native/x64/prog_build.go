// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/x64/p9x86"
)

func (prog *Prog) buildProg(as abi.As, arg *abi.X64Argument) (inst *Prog, err error) {
	dst := prog.operand2P9Addr(arg.Dst)
	src := prog.operand2P9Addr(arg.Src)

	switch as {
	case AADD: // add
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AADDSD
	case AADDSS: // addss
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AADDSS
	case AAND: // and
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ACALL
	case ACDQ: // cdq
		prog.As = p9x86.ACDQ
	case ACMP: // cmp
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 1)
		switch prog.xLen(arg) {
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
		switch prog.xLen(arg) {
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
		switch prog.xLen(arg) {
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
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 1)
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMOVQ
	case AMOVQ: // movq
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMOVQ
	case AMOVSD: // movsd
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMOVSD
	case AMOVSS: // movss
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMOVSS
	case AMOVSXD: // movsxd
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMULSD
	case AMULSS: // mulss
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMULSS
	case ANEG: // neg
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 0)
		prog.As = p9x86.ANOPW
	case AOR: // or
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 1)
		prog.To = dst
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 1)
		prog.From = src
		switch prog.xLen(arg) {
		case 1:
			// push
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
		assert(prog.nArg(arg) == 0)
		prog.As = p9x86.ARET
	case ASAR: // sar
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.ASUBSD
	case ASUBSS: // subss
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.ASUBSS
	case ASYSCALL: // syscall
		prog.As = p9x86.ASYSCALL
	case ATEST: // test
		assert(prog.nArg(arg) == 1)
		switch prog.xLen(arg) {
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
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
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
