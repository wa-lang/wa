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
		prog.From = src
		prog.To = dst

	case AADDSD: // addsd
		// addsd xmm4, qword ptr [rbp-64]
		assert(prog.nArg(arg) == 2)
		assert(RegValid_Float(arg.Dst.Reg))
		assert(RegValid_Float(arg.Src.Reg) || arg.Src.Kind == abi.X64Operand_Mem)
		prog.As = p9x86.AADDSD
		prog.From = src
		prog.To = dst

	case AADDSS: // addss
		// addss xmm4, dword ptr [rbp-64]
		assert(prog.nArg(arg) == 2)
		assert(RegValid_Float(arg.Dst.Reg))
		assert(RegValid_Float(arg.Src.Reg) || arg.Src.Kind == abi.X64Operand_Mem)
		prog.As = p9x86.AADDSS
		prog.From = src
		prog.To = dst

	case AAND: // and
		// and eax, dword ptr [rbp-1304]
		// and al, cl
		assert(prog.nArg(arg) == 2)
		assert(RegValid_Int(arg.Dst.Reg))
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
		prog.From = src
		prog.To = dst

	case ACALL: // call
		// call .Wa.F.runtime.printI64
		// call r11
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ACALL
		prog.To = dst

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
		prog.From = src
		prog.To = dst

	case ACMOVNE: // cmovne
		// cmovne r10d, r11d
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
		case 1:
			panic("unreachable")
		case 2:
			prog.As = p9x86.ACMOVWNE
		case 4:
			prog.As = p9x86.ACMOVLNE
		case 8:
			prog.As = p9x86.ACMOVQNE
		default:
			panic("unreachable")
		}
		prog.From = src
		prog.To = dst

	case ACQO: // cqo
		//  cqo # rdx = copysign(rax)
		assert(prog.nArg(arg) == 0)
		prog.As = p9x86.ACQO

	case ACVTSI2SD: // cvtsi2sd
		// cvtsi2sd xmm4, rax
		assert(prog.nArg(arg) == 2)
		assert(arg.Dst.Kind == abi.X64Operand_Reg)
		assert(arg.Src.Kind == abi.X64Operand_Reg || arg.Src.Kind == abi.X64Operand_Mem)
		switch prog.xLen(arg) {
		case 4:
			assert(RegValid_I32(arg.Src.Reg))
			assert(RegValid_Float(arg.Dst.Reg))
			prog.As = p9x86.ACVTSL2SD
		case 8:
			assert(RegValid_I64(arg.Src.Reg))
			assert(RegValid_Float(arg.Dst.Reg))
			prog.As = p9x86.ACVTSQ2SD
		default:
			panic("unreachable")
		}
		prog.From = src
		prog.To = dst

	case ACVTSS2SD: // cvtss2sd
		// cvtss2sd xmm4, xmm4
		assert(prog.nArg(arg) == 2)
		assert(arg.Dst.Kind == abi.X64Operand_Reg)
		assert(arg.Src.Kind == abi.X64Operand_Reg)
		prog.As = p9x86.ACVTSS2SD
		prog.From = src
		prog.To = dst

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
			prog.As = p9x86.ACVTSD2SL
		case 8:
			assert(RegValid_Float(arg.Src.Reg))
			assert(RegValid_I64(arg.Dst.Reg))
			prog.As = p9x86.ACVTSD2SQ
		default:
			panic("unreachable")
		}
		prog.From = src
		prog.To = dst

	case ACVTTSS2SI: // cvttss2si
		// cvttss2si eax, xmm4
		// cvttss2si rax, xmm4
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
		case 4:
			// float32 -> int32
			assert(RegValid_Float(arg.Src.Reg))
			assert(RegValid_I32(arg.Dst.Reg))
			prog.As = p9x86.ACVTSS2SL
		case 8:
			// float32 -> int64
			assert(RegValid_Float(arg.Src.Reg))
			assert(RegValid_I64(arg.Dst.Reg))
			prog.As = p9x86.ACVTSS2SQ
		default:
			panic("unreachable")
		}
		prog.From = src
		prog.To = dst

	case ACVTSI2SS: // cvtsi2ss
		// cvtsi2ss xmm4, eax
		// cvtsi2ss xmm4, rax
		assert(prog.nArg(arg) == 2)

		switch prog.xLen(arg) {
		case 4:
			// int32 -> float32
			assert(RegValid_I32(arg.Src.Reg))
			assert(RegValid_Float(arg.Dst.Reg))
			prog.As = p9x86.ACVTSL2SS
		case 8:
			// int64 -> float32
			assert(RegValid_I64(arg.Src.Reg))
			assert(RegValid_Float(arg.Dst.Reg))
			prog.As = p9x86.ACVTSQ2SS
		default:
			panic("unreachable")
		}
		prog.From = src
		prog.To = dst

	case ADEC: // dec
		// dec rcx
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
		prog.To = dst

	case ADIV: // div
		// div dword ptr [rbp-768]
		assert(prog.nArg(arg) == 1)
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
		prog.To = dst

	case ADIVSD: // divsd
		// divsd xmm4, qword ptr [rbp-160]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.ADIVSD
		prog.From = src
		prog.To = dst

	case ADIVSS: // divss
		// divss xmm4, dword ptr [rbp-160]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.ADIVSS
		prog.From = src
		prog.To = dst

	case AIDIV: // idiv
		// idiv dword ptr [rbp-768]
		assert(prog.nArg(arg) == 1)
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
		prog.To = dst

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
		prog.From = src
		prog.To = dst

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
		prog.To = dst

	case AJA: // ja
		// # memory.grow
		// ja .Wa.L.else.00000051
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.AJHI
		prog.To = dst

	case AJB: // jb
		// jb .Wa.L.memmove.forward
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.AJCS
		prog.To = dst

	case AJE: // je
		// je .Wa.L.brNext.00000052
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.AJEQ
		prog.To = dst

	case AJGE: // jge
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.AJEQ
		prog.To = dst

	case AJMP: // jmp
		// jmp .Wa.L.brNext.00000052
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.AJMP
		prog.To = dst

	case AJNS: // jns
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.AJNE
		prog.To = dst

	case AJNZ: // jnz
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.AJNE
		prog.To = dst

	case AJZ: // jz
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.AJEQ
		prog.To = dst

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
		prog.From = src
		prog.To = dst

	case ALZCNT: // lzcnt
		// lzcnt eax, eax
		// lzcnt rax, rax
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
		case 1:
			panic("unreachable")
		case 2:
			prog.As = p9x86.ALZCNTW
		case 4:
			prog.As = p9x86.ALZCNTL
		case 8:
			prog.As = p9x86.ALZCNTQ
		default:
			panic("unreachable")
		}
		prog.From = src
		prog.To = dst

	case AMAXSD: // maxsd
		//  maxsd xmm4, qword ptr [rbp+8]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMAXSD
		prog.From = src
		prog.To = dst

	case AMAXSS: // maxss
		// maxss xmm4, dword ptr [rbp+8]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMAXSS
		prog.From = src
		prog.To = dst

	case AMINSD: // minsd
		// minsd xmm4, qword ptr [rbp+8]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMINSD
		prog.From = src
		prog.To = dst

	case AMINSS: // minss
		// minss xmm4, dword ptr [rbp+8]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMINSS
		prog.From = src
		prog.To = dst

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
		prog.From = src
		prog.To = dst

	case AMOVABS: // movabs
		// movabs rax, 0
		// movabs rax, 0x3FF0000000000000
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMOVQ
		prog.From = src
		prog.To = dst

	case AMOVSD: // movsd
		// movsd xmm4, qword ptr [rbp-320]
		// movsd qword ptr [rbp-760], xmm4
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMOVSD
		prog.From = src
		prog.To = dst

	case AMOVSS: // movss
		// movss xmm4, dword ptr [rbp+16]
		// movss dword ptr [rbp-24], xmm4
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMOVSS
		prog.From = src
		prog.To = dst

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
		prog.From = src
		prog.To = dst

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
		prog.From = src
		prog.To = dst

	case AMULSD: // mulsd
		// mulsd xmm4, qword ptr [rbp-160]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMULSD
		prog.From = src
		prog.To = dst

	case AMULSS: // mulss
		// mulss xmm4, dword ptr [rbp-160]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AMULSS
		prog.From = src
		prog.To = dst

	case ANEG: // neg
		// neg rax
		assert(prog.nArg(arg) == 1)
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
		prog.To = dst

	case ANOP: // nop
		assert(prog.nArg(arg) == 0)
		prog.As = p9x86.ANOPW

	case AOR: // or
		// or eax, dword ptr [rbp-1304]
		assert(prog.nArg(arg) == 2)
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
		prog.From = src
		prog.To = dst

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
		prog.To = dst

	case APOPCNT: // popcnt
		// popcnt eax, eax
		// popcnt rax, rax
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
		case 1:
			panic("unreachable")
		case 2:
			prog.As = p9x86.APOPCNTW
		case 4:
			prog.As = p9x86.APOPCNTL
		case 8:
			prog.As = p9x86.APOPCNTQ
		default:
			panic("unreachable")
		}
		prog.From = src
		prog.To = dst

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
		prog.To = dst

	case ARET: // ret
		assert(prog.nArg(arg) == 0)
		prog.As = p9x86.ARET

	case AROL: // rol
		// rol  eax, cl # cl 是 ecx 低8位
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
		case 1:
			prog.As = p9x86.AROLB
		case 2:
			prog.As = p9x86.AROLW
		case 4:
			prog.As = p9x86.AROLL
		case 8:
			prog.As = p9x86.AROLQ
		default:
			panic("unreachable")
		}
		prog.From = src
		prog.To = dst

	case AROR: // ror
		// ror eax, cl # cl 是 ecx 低8位
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
		case 1:
			prog.As = p9x86.ARORB
		case 2:
			prog.As = p9x86.ARORW
		case 4:
			prog.As = p9x86.ARORL
		case 8:
			prog.As = p9x86.ARORQ
		default:
			panic("unreachable")
		}
		prog.From = src
		prog.To = dst

	case AROUNDSD: // roundsd
		// roundsd xmm4, xmm4, 2
		assert(prog.nArg(arg) == 3)
		prog.As = p9x86.AROUNDSD
		prog.From = src
		prog.To = dst
		prog.RestArgs = append(prog.RestArgs, prog.operand2P9Addr(arg.Rest[0]))

	case AROUNDSS: // roundss
		// roundss xmm4, xmm4, 2
		assert(prog.nArg(arg) == 3)
		prog.As = p9x86.AROUNDSS
		prog.From = src
		prog.To = dst
		prog.RestArgs = append(prog.RestArgs, prog.operand2P9Addr(arg.Rest[0]))

	case ASAR: // sar
		// sar eax, cl # cl 是 ecx 低8位
		assert(prog.nArg(arg) == 2)
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
		prog.From = src
		prog.To = dst

	case ASETA: // seta
		// seta al
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ASETHI
		prog.To = dst

	case ASETAE: // setae
		// setae al
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ASETCC
		prog.To = dst

	case ASETB: // setb
		// setb al
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ASETCS
		prog.To = dst

	case ASETBE: // setbe
		// setbe al
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ASETLS
		prog.To = dst

	case ASETE: // sete
		// sete al # al = (eax==0)? 1: 0
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ASETEQ
		prog.To = dst

	case ASETG: // setg
		// setg al
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ASETGT
		prog.To = dst

	case ASETGE: // setge
		// setge al
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ASETGE
		prog.To = dst

	case ASETL: // setl
		// setl al
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ASETLT
		prog.To = dst

	case ASETLE: // setle
		// setle al
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ASETLE
		prog.To = dst

	case ASETNE: // setne
		// setne al # al = (r10d==r11d)? 1: 0
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ASETNE
		prog.To = dst

	case ASETNP: // setnp
		// setnp cl # set if not NaN
		assert(prog.nArg(arg) == 1)
		prog.As = p9x86.ASETPC
		prog.To = dst

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
		prog.From = src
		prog.To = dst

	case ASHR: // shr
		// shr eax, cl # cl 是 ecx 低8位
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
		case 1:
			prog.As = p9x86.ASHRB
		case 2:
			prog.As = p9x86.ASHRW
		case 4:
			prog.As = p9x86.ASHRL
		case 8:
			prog.As = p9x86.ASHRQ
		default:
			panic("unreachable")
		}
		prog.From = src
		prog.To = dst

	case ASQRTSD: // sqrtsd
		// sqrtsd xmm4, xmm4
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.ASQRTSD
		prog.From = src
		prog.To = dst

	case ASQRTSS: // sqrtss
		// sqrtss xmm4, xmm4
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.ASQRTSS
		prog.From = src
		prog.To = dst

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
		prog.From = src
		prog.To = dst

	case ASUBSD: // subsd
		// subsd xmm4, qword ptr [rbp-768]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.ASUBSD
		prog.From = src
		prog.To = dst

	case ASUBSS: // subss
		// subss xmm4, dword ptr [rbp-160]
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.ASUBSS
		prog.From = src
		prog.To = dst

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
		prog.From = src
		prog.To = dst

	case ATZCNT: // tzcnt
		// tzcnt eax, eax
		assert(prog.nArg(arg) == 2)
		switch prog.xLen(arg) {
		case 1:
			panic("unreachable")
		case 2:
			prog.As = p9x86.ATZCNTW
		case 4:
			prog.As = p9x86.ATZCNTL
		case 8:
			prog.As = p9x86.ATZCNTQ
		default:
			panic("unreachable")
		}
		prog.From = src
		prog.To = dst

	case AUCOMISD: // ucomisd
		// ucomisd xmm4, xmm5
		assert(prog.nArg(arg) == 2)
		prog.As = p9x86.AUCOMISD
		prog.From = src
		prog.To = dst

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
		prog.From = src
		prog.To = dst

	default:
		panic(fmt.Sprintf("x64: unsupport as: %v", as))
	}

	return prog, nil
}
