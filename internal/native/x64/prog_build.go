// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/x64/p9x86"
)

// https://www.felixcloutier.com/x86/

func (prog *Prog) buildProg(as abi.As, arg *abi.X64Argument) (inst *Prog, err error) {
	dst := prog.operand2P9Addr(arg.Dst)
	src := prog.operand2P9Addr(arg.Src)

	switch as {
	case AADD: // add
		// add rdi, rdx
		// add rsp, 16
		// add dl, '0'
		// add eax, dword ptr [rbp-16]
		assert(prog.nArg(arg) == 2)
		assert(RegValid_Int(arg.Dst.Reg))
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
		// addsd xmm4, qword ptr [rbp-64]
		assert(prog.nArg(arg) == 2)
		assert(RegValid_Float(arg.Dst.Reg))
		prog.As = p9x86.AADDSD
	case AADDSS: // addss
		// addss xmm4, dword ptr [rbp-64]
		assert(prog.nArg(arg) == 2)
		assert(RegValid_Float(arg.Dst.Reg))
		prog.As = p9x86.AADDSS
	case AAND: // and
		// and eax, dword ptr [rbp-1304]
		// and al, cl
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
		// call .Wa.F.runtime.printI64
		// call r11
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ACALL
	case ACDQ: // cdq
		// cdq # edx = copysign(eax)
		assert(prog.nArg(arg) == 0)
		prog.As = p9x86.ACDQ
	case ACMP: // cmp
		// cmp eax, 0 # (eax==0)?
		// cmp r10d, r11d
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
		// cvtsi2sd xmm4, rax
		assert(prog.nArg(arg) == 2)
		assert(arg.Dst.Kind == abi.X64Operand_Reg)
		assert(arg.Src.Kind == abi.X64Operand_Reg)
		switch prog.xLen(arg) {
		case 4:
			assert(RegValid_I32(arg.Src.Reg))
			assert(RegValid_Float(arg.Dst.Reg))
			prog.As = p9x86.AMOVL
		case 8:
			assert(RegValid_I64(arg.Src.Reg))
			assert(RegValid_Float(arg.Dst.Reg))
			prog.As = p9x86.AMOVQ
		default:
			panic("unreachable")
		}
	case ACVTSS2SD: // cvtss2sd
		assert(prog.nArg(arg) == 2)
		assert(arg.Dst.Kind == abi.X64Operand_Reg)
		assert(arg.Src.Kind == abi.X64Operand_Reg)
		prog.As = p9x86.AMOVQ
	case ACVTTSD2SI: // cvttsd2si
		// cvttsd2si rax, xmm4
		// cvttsd2si eax, xmm4
		assert(prog.nArg(arg) == 2)
		assert(arg.Dst.Kind == abi.X64Operand_Reg)
		assert(arg.Src.Kind == abi.X64Operand_Reg)
		switch prog.xLen(arg) {
		case 4:
			assert(RegValid_Float(arg.Src.Reg))
			assert(RegValid_I32(arg.Dst.Reg))
			prog.As = p9x86.AMOVL
		case 8:
			assert(RegValid_Float(arg.Src.Reg))
			assert(RegValid_I64(arg.Dst.Reg))
			prog.As = p9x86.AMOVQ
		default:
			panic("unreachable")
		}
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
		// div dword ptr [rbp-768]
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
		// divsd xmm4, qword ptr [rbp-160]
		prog.As = p9x86.ADIVSD
	case ADIVSS: // divss
		// divss xmm4, dword ptr [rbp-160]
		prog.As = p9x86.ADIVSS
	case AIDIV: // idiv
		// idiv dword ptr [rbp-768]
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
		// imul eax, dword ptr [rbp-784]
		assert(prog.nArg(arg) == 2)
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
		// # memory.grow
		// ja .Wa.L.else.00000051
		prog.As = p9x86.AJHI
	case AJB: // jb
		prog.As = p9x86.AJCS
	case AJE: // je
		// je .Wa.L.brNext.00000052
		prog.As = p9x86.AJEQ
	case AJGE: // jge
		prog.As = p9x86.AJEQ
	case AJMP: // jmp
		// jmp .Wa.L.brNext.00000052
		prog.As = p9x86.AJMP
	case AJNS: // jns
		prog.As = p9x86.AJNE
	case AJNZ: // jnz
		prog.As = p9x86.AJNE
	case AJZ: // jz
		prog.As = p9x86.AJEQ
	case ALEA: // lea
		// lea rcx, [rsp+40] # return address
		// lea rax, [rip+.Wa.Table.funcIndexList]
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
		// mov rax, -1
		// mov rsp, rbp
		// mov ecx, dword ptr [rbp-16] # arg 0
		// mov edx, dword ptr [rbp-24] # arg 1
		// mov r8, qword ptr [rip + .Wa.Memory.maxPages]
		// mov qword ptr [rax+1], 61
		// mov dword ptr [rbp-16], eax
		// mov qword ptr [rbp-16], rax
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
		// movabs rax, 0
		// movabs rax, 0x3FF0000000000000
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMOVQ
	case AMOVQ: // movq
		// movq [rbp-760], rax
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMOVQ
	case AMOVSD: // movsd
		// movsd xmm4, qword ptr [rbp-320]
		// movsd qword ptr [rbp-760], xmm4
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMOVSD
	case AMOVSS: // movss
		// movss xmm4, dword ptr [rbp+16]
		// movss dword ptr [rbp-24], xmm4
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
		// movzx eax, al # eax = al
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
		// mulsd xmm4, qword ptr [rbp-160]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMULSD
	case AMULSS: // mulss
		// mulss xmm4, dword ptr [rbp-160]
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
		// or eax, dword ptr [rbp-1304]
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
		// pop rcx
		assert(prog.nArg(arg) == 1)
		prog.To = dst
		switch prog.xLen(arg) {
		case 1:
			assert(arg.Dst.Reg == REG_AL)
			prog.As = p9x86.APOPAL
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
		// push rbp
		assert(prog.nArg(arg) == 1)
		prog.From = src
		switch prog.xLen(arg) {
		case 1:
			assert(arg.Src.Reg == REG_AL)
			prog.As = p9x86.APUSHAL
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
		// sar eax, cl # cl 是 ecx 低8位
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
		// seta al
		panic("TODO")
	case ASETAE: // setae
		// setae al
		panic("TODO")
	case ASETB: // setb
		// setb al
		panic("TODO")
	case ASETBE: // setbe
		// setbe al
		panic("TODO")
	case ASETE: // sete
		// sete al # al = (eax==0)? 1: 0
		panic("TODO")
	case ASETG: // setg
		// setg al
		panic("TODO")
	case ASETGE: // setge
		// setge al
		panic("TODO")
	case ASETL: // setl
		// setl al
		panic("TODO")
	case ASETLE: // setle
		// setle al
		panic("TODO")
	case ASETNE: // setne
		// setne al # al = (r10d==r11d)? 1: 0
		panic("TODO")
	case ASETNP: // setnp
		// setnp cl # set if not NaN
		panic("TODO")
	case ASHL: // shl
		// shl eax, cl # cl 是 ecx 低8位
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
		// sub rsp, 80
		// sub eax, dword ptr [rbp-464]
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
		// subsd xmm4, qword ptr [rbp-768]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.ASUBSD
	case ASUBSS: // subss
		// subss xmm4, dword ptr [rbp-160]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.ASUBSS
	case ASYSCALL: // syscall
		assert(prog.nArg(arg) == 0)
		prog.As = p9x86.ASYSCALL
	case ATEST: // test
		// test rdx, rdx
		assert(prog.nArg(arg) == 2)
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
		// xor rdx, rdx
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
