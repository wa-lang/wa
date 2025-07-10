// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package arch

import (
	"strings"

	"wa-lang.org/wa/internal/p9asm/arm64"
	"wa-lang.org/wa/internal/p9asm/obj"
	"wa-lang.org/wa/internal/p9asm/x86"
)

// 处理器类型
type CPUType byte

const (
	NoArch CPUType = iota
	AMD64
	ARM
	ARM64
	Loong64
	RISCV64
	Wasm
)

// Plan9 汇编语言的伪寄存器
const (
	RFP = -(iota + 1)
	RSB
	RSP
	RPC
)

// Arch 封装了链接器的架构对象，并附加了更多与架构相关的具体信息。
//
// 注意: 尽量避免在该结构体类型中直接依赖 p9asm/obj 包的类型
type Arch struct {
	CPU CPUType

	// 指令集表
	// 比如 X86 平台的 JCC 指令对应 x86.AJCC
	Instructions map[string]int

	// 通用寄存器表
	// 比如 X86 平台的 AX 寄存器对应 x86.REG_AX
	Register map[string]int16

	// 通用寄存器名字前缀表格, 比如 R(0) 的前缀为 R, SPR(268) 的前缀为 SPR
	// 该表格需要和 this.RegisterNumber() 辅助函数配合使用
	RegisterPrefix map[string]bool

	// 根据通用寄存器的名字前缀和编号, 转化为底层唯一的机器码对应的编号.
	// 比如 R(10) 对应 arm64.REG_R10.
	RegisterNumber func(string, int16) (int16, bool)

	// 哪些指令是只有一个目的操作数的, 比如某些一元运算
	UnaryDst map[int]bool

	// 判断是否为 jump 跳转指令
	IsJump func(word string) bool
}

var Pseudos = map[string]int{
	"DATA":     obj.ADATA,
	"FUNCDATA": obj.AFUNCDATA,
	"GLOBL":    obj.AGLOBL,
	"PCDATA":   obj.APCDATA,
	"TEXT":     obj.ATEXT,
}

// 设置并返回对应的CPU对象
func Set(cpu CPUType) *Arch {
	switch cpu {
	case AMD64:
		return archX86(AMD64)
	case ARM64:
		return archArm64(ARM64)
	}

	return nil
}

func archX86(CPU CPUType) *Arch {
	register := make(map[string]int16)
	// Create maps for easy lookup of instruction names etc.
	for i, s := range x86.Register {
		register[s] = int16(i + x86.REG_AL)
	}
	// Pseudo-registers.
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC
	// Register prefix not used on this architecture.

	instructions := make(map[string]int)
	for i, s := range obj.Anames {
		instructions[s] = i
	}
	for i, s := range x86.Anames {
		if i >= obj.A_ARCHSPECIFIC {
			instructions[s] = i + obj.ABaseAMD64
		}
	}
	// Annoying aliases.
	instructions["JA"] = x86.AJHI   /* alternate */
	instructions["JAE"] = x86.AJCC  /* alternate */
	instructions["JB"] = x86.AJCS   /* alternate */
	instructions["JBE"] = x86.AJLS  /* alternate */
	instructions["JC"] = x86.AJCS   /* alternate */
	instructions["JCC"] = x86.AJCC  /* carry clear (CF = 0) */
	instructions["JCS"] = x86.AJCS  /* carry set (CF = 1) */
	instructions["JE"] = x86.AJEQ   /* alternate */
	instructions["JEQ"] = x86.AJEQ  /* equal (ZF = 1) */
	instructions["JG"] = x86.AJGT   /* alternate */
	instructions["JGE"] = x86.AJGE  /* greater than or equal (signed) (SF = OF) */
	instructions["JGT"] = x86.AJGT  /* greater than (signed) (ZF = 0 && SF = OF) */
	instructions["JHI"] = x86.AJHI  /* higher (unsigned) (CF = 0 && ZF = 0) */
	instructions["JHS"] = x86.AJCC  /* alternate */
	instructions["JL"] = x86.AJLT   /* alternate */
	instructions["JLE"] = x86.AJLE  /* less than or equal (signed) (ZF = 1 || SF != OF) */
	instructions["JLO"] = x86.AJCS  /* alternate */
	instructions["JLS"] = x86.AJLS  /* lower or same (unsigned) (CF = 1 || ZF = 1) */
	instructions["JLT"] = x86.AJLT  /* less than (signed) (SF != OF) */
	instructions["JMI"] = x86.AJMI  /* negative (minus) (SF = 1) */
	instructions["JNA"] = x86.AJLS  /* alternate */
	instructions["JNAE"] = x86.AJCS /* alternate */
	instructions["JNB"] = x86.AJCC  /* alternate */
	instructions["JNBE"] = x86.AJHI /* alternate */
	instructions["JNC"] = x86.AJCC  /* alternate */
	instructions["JNE"] = x86.AJNE  /* not equal (ZF = 0) */
	instructions["JNG"] = x86.AJLE  /* alternate */
	instructions["JNGE"] = x86.AJLT /* alternate */
	instructions["JNL"] = x86.AJGE  /* alternate */
	instructions["JNLE"] = x86.AJGT /* alternate */
	instructions["JNO"] = x86.AJOC  /* alternate */
	instructions["JNP"] = x86.AJPC  /* alternate */
	instructions["JNS"] = x86.AJPL  /* alternate */
	instructions["JNZ"] = x86.AJNE  /* alternate */
	instructions["JO"] = x86.AJOS   /* alternate */
	instructions["JOC"] = x86.AJOC  /* overflow clear (OF = 0) */
	instructions["JOS"] = x86.AJOS  /* overflow set (OF = 1) */
	instructions["JP"] = x86.AJPS   /* alternate */
	instructions["JPC"] = x86.AJPC  /* parity clear (PF = 0) */
	instructions["JPE"] = x86.AJPS  /* alternate */
	instructions["JPL"] = x86.AJPL  /* non-negative (plus) (SF = 0) */
	instructions["JPO"] = x86.AJPC  /* alternate */
	instructions["JPS"] = x86.AJPS  /* parity set (PF = 1) */
	instructions["JS"] = x86.AJMI   /* alternate */
	instructions["JZ"] = x86.AJEQ   /* alternate */
	instructions["MASKMOVDQU"] = x86.AMASKMOVOU
	instructions["MOVD"] = x86.AMOVQ
	instructions["MOVDQ2Q"] = x86.AMOVQ
	instructions["MOVNTDQ"] = x86.AMOVNTO
	instructions["MOVOA"] = x86.AMOVO
	instructions["MOVOA"] = x86.AMOVO
	instructions["PF2ID"] = x86.APF2IL
	instructions["PI2FD"] = x86.API2FL
	instructions["PSLLDQ"] = x86.APSLLO
	instructions["PSRLDQ"] = x86.APSRLO

	return &Arch{
		CPU:          CPU,
		Instructions: instructions,
		Register:     register,

		// nil的map依然可以正常返回 x[key] 为零值
		RegisterPrefix: nil,
		// nil的func被调用会导致panic, 因此需要一个空函数
		RegisterNumber: func(name string, n int16) (int16, bool) { return 0, false },

		UnaryDst: x86.Linkamd64.UnaryDst,

		// 判断是否为 X86 平台的 Jump 指令
		IsJump: func(word string) bool {
			return word[0] == 'J' || word == "CALL" || strings.HasPrefix(word, "LOOP")
		},
	}
}

func archArm64(CPU CPUType) *Arch {
	register := make(map[string]int16)
	// Create maps for easy lookup of instruction names etc.
	// Note that there is no list of names as there is for 386 and amd64.
	register[arm64.Rconv(arm64.REGSP)] = int16(arm64.REGSP)
	for i := arm64.REG_R0; i <= arm64.REG_R31; i++ {
		register[arm64.Rconv(i)] = int16(i)
	}
	for i := arm64.REG_F0; i <= arm64.REG_F31; i++ {
		register[arm64.Rconv(i)] = int16(i)
	}
	for i := arm64.REG_V0; i <= arm64.REG_V31; i++ {
		register[arm64.Rconv(i)] = int16(i)
	}
	register["LR"] = arm64.REGLINK
	register["DAIF"] = arm64.REG_DAIF
	register["NZCV"] = arm64.REG_NZCV
	register["FPSR"] = arm64.REG_FPSR
	register["FPCR"] = arm64.REG_FPCR
	register["SPSR_EL1"] = arm64.REG_SPSR_EL1
	register["ELR_EL1"] = arm64.REG_ELR_EL1
	register["SPSR_EL2"] = arm64.REG_SPSR_EL2
	register["ELR_EL2"] = arm64.REG_ELR_EL2
	register["CurrentEL"] = arm64.REG_CurrentEL
	register["SP_EL0"] = arm64.REG_SP_EL0
	register["SPSel"] = arm64.REG_SPSel
	register["DAIFSet"] = arm64.REG_DAIFSet
	register["DAIFClr"] = arm64.REG_DAIFClr
	// Conditional operators, like EQ, NE, etc.
	register["EQ"] = arm64.COND_EQ
	register["NE"] = arm64.COND_NE
	register["HS"] = arm64.COND_HS
	register["LO"] = arm64.COND_LO
	register["MI"] = arm64.COND_MI
	register["PL"] = arm64.COND_PL
	register["VS"] = arm64.COND_VS
	register["VC"] = arm64.COND_VC
	register["HI"] = arm64.COND_HI
	register["LS"] = arm64.COND_LS
	register["GE"] = arm64.COND_GE
	register["LT"] = arm64.COND_LT
	register["GT"] = arm64.COND_GT
	register["LE"] = arm64.COND_LE
	register["AL"] = arm64.COND_AL
	register["NV"] = arm64.COND_NV
	// Pseudo-registers.
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC
	register["SP"] = RSP
	// Avoid unintentionally clobbering g using R28.
	delete(register, "R28")
	register["g"] = arm64.REG_R28
	registerPrefix := map[string]bool{
		"F": true,
		"R": true,
		"V": true,
	}

	instructions := make(map[string]int)
	for i, s := range obj.Anames {
		instructions[s] = i
	}
	for i, s := range arm64.Anames {
		if i >= obj.A_ARCHSPECIFIC {
			instructions[s] = i + obj.ABaseARM64
		}
	}
	// Annoying aliases.
	instructions["B"] = arm64.AB
	instructions["BL"] = arm64.ABL

	return &Arch{
		CPU:            CPU,
		Instructions:   instructions,
		Register:       register,
		RegisterPrefix: registerPrefix,
		RegisterNumber: arm64RegisterNumber,
		UnaryDst:       arm64.Linkarm64.UnaryDst,
		IsJump:         jumpArm64,
	}

}
