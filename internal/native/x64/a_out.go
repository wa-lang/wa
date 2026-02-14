// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import "wa-lang.org/wa/internal/native/abi"

const (
	// 通用寄存器
	_ abi.RegType = iota // 0 是无效的编号

	// 指针计数器
	REG_RIP

	// 低8位寄存器
	REG_AL
	REG_CL
	REG_DL
	REG_BL
	REG_SIL
	REG_DIL
	REG_SPL
	REG_BPL
	REG_R8B
	REG_R9B
	REG_R10B
	REG_R11B
	REG_R12B
	REG_R13B
	REG_R14B
	REG_R15B

	// 高8位寄存器
	REG_AH
	REG_CH
	REG_DH
	REG_BH

	// 16位寄存器
	REG_AX
	REG_CX
	REG_DX
	REG_BX
	REG_SI
	REG_DI
	REG_SP
	REG_BP
	REG_R8W
	REG_R9W
	REG_R10W
	REG_R11W
	REG_R12W
	REG_R13W
	REG_R14W
	REG_R15W

	// 32位寄存器
	REG_EAX
	REG_ECX
	REG_EDX
	REG_EBX
	REG_ESP
	REG_EBP
	REG_ESI
	REG_EDI
	REG_R8D
	REG_R9D
	REG_R10D
	REG_R11D
	REG_R12D
	REG_R13D
	REG_R14D
	REG_R15D

	// 64位寄存器
	REG_RAX
	REG_RCX
	REG_RDX
	REG_RBX
	REG_RSP
	REG_RBP
	REG_RSI
	REG_RDI
	REG_R8
	REG_R9
	REG_R10
	REG_R11
	REG_R12
	REG_R13
	REG_R14
	REG_R15

	// 浮点数寄存器
	REG_XMM0
	REG_XMM1
	REG_XMM2
	REG_XMM3
	REG_XMM4
	REG_XMM5
	REG_XMM6
	REG_XMM7

	// 寄存器编号结束
	REG_END
)

// 凹语言用到的部分指令
const (
	_          abi.As = iota
	AADD              // add
	AADDSD            // addsd
	AADDSS            // addss
	AAND              // and
	ACALL             // call
	ACDQ              // cdq
	ACMP              // cmp
	ACVTSI2SD         // cvtsi2sd
	ACVTSS2SD         // cvtss2sd
	ACVTTSD2SI        // cvttsd2si
	ADEC              // dec
	ADIV              // div
	ADIVSD            // divsd
	ADIVSS            // divss
	AIDIV             // idiv
	AIMUL             // imul
	AINC              // inc
	AJA               // ja
	AJB               // jb
	AJE               // je
	AJGE              // jge
	AJMP              // jmp
	AJNS              // jns
	AJNZ              // jnz
	AJZ               // jz
	ALEA              // lea
	AMOV              // mov
	AMOVABS           // movabs
	AMOVQ             // movq
	AMOVSD            // movsd
	AMOVSS            // movss
	AMOVSXD           // movsxd
	AMOVZX            // movzx
	AMULSD            // mulsd
	AMULSS            // mulss
	ANEG              // neg
	ANOP              // nop
	AOR               // or
	APOP              // pop
	APUSH             // push
	ARET              // ret
	ASAR              // sar
	ASETA             // seta
	ASETAE            // setae
	ASETB             // setb
	ASETBE            // setbe
	ASETE             // sete
	ASETG             // setg
	ASETGE            // setge
	ASETL             // setl
	ASETLE            // setle
	ASETNE            // setne
	ASETNP            // setnp
	ASHL              // shl
	ASTD              // std
	ASUB              // sub
	ASUBSD            // subsd
	ASUBSS            // subss
	ASYSCALL          // syscall
	ATEST             // test
	AUCOMISD          // ucomisd
	AXOR              // xor

	// End marker
	ALAST
)
