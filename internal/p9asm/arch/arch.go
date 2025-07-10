// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package arch

import (
	"fmt"
	"strings"

	"wa-lang.org/wa/internal/p9asm/arm"
	"wa-lang.org/wa/internal/p9asm/arm64"
	"wa-lang.org/wa/internal/p9asm/obj"
	"wa-lang.org/wa/internal/p9asm/riscv"
	"wa-lang.org/wa/internal/p9asm/wasm"
	"wa-lang.org/wa/internal/p9asm/x86"
)

// 处理器类型
type CPUType byte

const (
	NoArch CPUType = iota
	AMD64
	ARM
	ARM64
	I386
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

// Plan9 汇编语言的伪指令
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
	case I386:
		return archX86(I386)
	case AMD64:
		return archX86(AMD64)
	case ARM:
		return archArm(ARM)
	case ARM64:
		return archArm64(ARM64)
	case RISCV64:
		return archRISCV64(false)
	case Wasm:
		return archWasm()
	}

	return nil
}

func archX86(CPU CPUType) *Arch {
	p := &Arch{CPU: CPU}

	// 构造寄存器名查询表
	for i, s := range x86.Register {
		p.Register[s] = int16(i + x86.REG_AL)
	}

	// Plan9 汇编语言伪寄存器
	// X86 已经有了 SP 名字的寄存器, 因此跳过 SP 伪寄存器注册
	p.Register["SB"] = RSB
	p.Register["FP"] = RFP
	p.Register["PC"] = RPC

	// X86 没有通用前缀名的寄存器
	p.RegisterPrefix = map[string]bool{}
	p.RegisterNumber = func(name string, n int16) (int16, bool) { return 0, false }

	// 初始化 Plan9 汇编语言的通用指令
	for i, s := range obj.Anames {
		p.Instructions[s] = i
	}

	// 初始化 X86 平台的指令
	for i, s := range x86.Anames {
		if i >= obj.A_ARCHSPECIFIC {
			p.Instructions[s] = i + obj.ABaseAMD64
		}
	}

	// 哪些指令是只有一个目的操作数的, 比如某些一元运算
	switch CPU {
	case I386:
		p.UnaryDst = x86.Link386.UnaryDst
	case AMD64:
		p.UnaryDst = x86.Linkamd64.UnaryDst
	}

	// 判断是否为 X86 平台的 Jump 指令
	p.IsJump = func(word string) bool {
		return word[0] == 'J' || word == "CALL" || strings.HasPrefix(word, "LOOP")
	}

	// 一些有别名的指令
	p.Instructions["JA"] = x86.AJHI   /* alternate */
	p.Instructions["JAE"] = x86.AJCC  /* alternate */
	p.Instructions["JB"] = x86.AJCS   /* alternate */
	p.Instructions["JBE"] = x86.AJLS  /* alternate */
	p.Instructions["JC"] = x86.AJCS   /* alternate */
	p.Instructions["JCC"] = x86.AJCC  /* carry clear (CF = 0) */
	p.Instructions["JCS"] = x86.AJCS  /* carry set (CF = 1) */
	p.Instructions["JE"] = x86.AJEQ   /* alternate */
	p.Instructions["JEQ"] = x86.AJEQ  /* equal (ZF = 1) */
	p.Instructions["JG"] = x86.AJGT   /* alternate */
	p.Instructions["JGE"] = x86.AJGE  /* greater than or equal (signed) (SF = OF) */
	p.Instructions["JGT"] = x86.AJGT  /* greater than (signed) (ZF = 0 && SF = OF) */
	p.Instructions["JHI"] = x86.AJHI  /* higher (unsigned) (CF = 0 && ZF = 0) */
	p.Instructions["JHS"] = x86.AJCC  /* alternate */
	p.Instructions["JL"] = x86.AJLT   /* alternate */
	p.Instructions["JLE"] = x86.AJLE  /* less than or equal (signed) (ZF = 1 || SF != OF) */
	p.Instructions["JLO"] = x86.AJCS  /* alternate */
	p.Instructions["JLS"] = x86.AJLS  /* lower or same (unsigned) (CF = 1 || ZF = 1) */
	p.Instructions["JLT"] = x86.AJLT  /* less than (signed) (SF != OF) */
	p.Instructions["JMI"] = x86.AJMI  /* negative (minus) (SF = 1) */
	p.Instructions["JNA"] = x86.AJLS  /* alternate */
	p.Instructions["JNAE"] = x86.AJCS /* alternate */
	p.Instructions["JNB"] = x86.AJCC  /* alternate */
	p.Instructions["JNBE"] = x86.AJHI /* alternate */
	p.Instructions["JNC"] = x86.AJCC  /* alternate */
	p.Instructions["JNE"] = x86.AJNE  /* not equal (ZF = 0) */
	p.Instructions["JNG"] = x86.AJLE  /* alternate */
	p.Instructions["JNGE"] = x86.AJLT /* alternate */
	p.Instructions["JNL"] = x86.AJGE  /* alternate */
	p.Instructions["JNLE"] = x86.AJGT /* alternate */
	p.Instructions["JNO"] = x86.AJOC  /* alternate */
	p.Instructions["JNP"] = x86.AJPC  /* alternate */
	p.Instructions["JNS"] = x86.AJPL  /* alternate */
	p.Instructions["JNZ"] = x86.AJNE  /* alternate */
	p.Instructions["JO"] = x86.AJOS   /* alternate */
	p.Instructions["JOC"] = x86.AJOC  /* overflow clear (OF = 0) */
	p.Instructions["JOS"] = x86.AJOS  /* overflow set (OF = 1) */
	p.Instructions["JP"] = x86.AJPS   /* alternate */
	p.Instructions["JPC"] = x86.AJPC  /* parity clear (PF = 0) */
	p.Instructions["JPE"] = x86.AJPS  /* alternate */
	p.Instructions["JPL"] = x86.AJPL  /* non-negative (plus) (SF = 0) */
	p.Instructions["JPO"] = x86.AJPC  /* alternate */
	p.Instructions["JPS"] = x86.AJPS  /* parity set (PF = 1) */
	p.Instructions["JS"] = x86.AJMI   /* alternate */
	p.Instructions["JZ"] = x86.AJEQ   /* alternate */
	p.Instructions["MASKMOVDQU"] = x86.AMASKMOVOU
	p.Instructions["MOVD"] = x86.AMOVQ
	p.Instructions["MOVDQ2Q"] = x86.AMOVQ
	p.Instructions["MOVNTDQ"] = x86.AMOVNTO
	p.Instructions["MOVOA"] = x86.AMOVO
	p.Instructions["MOVOA"] = x86.AMOVO
	p.Instructions["PF2ID"] = x86.APF2IL
	p.Instructions["PI2FD"] = x86.API2FL
	p.Instructions["PSLLDQ"] = x86.APSLLO
	p.Instructions["PSRLDQ"] = x86.APSRLO

	return p
}

func archArm64(CPU CPUType) *Arch {
	p := &Arch{CPU: CPU}

	// ARM64 的寄存器表
	// 注册方式和 AMD64 稍微有点区别
	p.Register[arm64.Rconv(arm64.REGSP)] = int16(arm64.REGSP)
	for i := arm64.REG_R0; i <= arm64.REG_R31; i++ {
		p.Register[arm64.Rconv(i)] = int16(i)
	}
	for i := arm64.REG_F0; i <= arm64.REG_F31; i++ {
		p.Register[arm64.Rconv(i)] = int16(i)
	}
	for i := arm64.REG_V0; i <= arm64.REG_V31; i++ {
		p.Register[arm64.Rconv(i)] = int16(i)
	}

	p.Register["LR"] = arm64.REGLINK
	p.Register["DAIF"] = arm64.REG_DAIF
	p.Register["NZCV"] = arm64.REG_NZCV
	p.Register["FPSR"] = arm64.REG_FPSR
	p.Register["FPCR"] = arm64.REG_FPCR
	p.Register["SPSR_EL1"] = arm64.REG_SPSR_EL1
	p.Register["ELR_EL1"] = arm64.REG_ELR_EL1
	p.Register["SPSR_EL2"] = arm64.REG_SPSR_EL2
	p.Register["ELR_EL2"] = arm64.REG_ELR_EL2
	p.Register["CurrentEL"] = arm64.REG_CurrentEL
	p.Register["SP_EL0"] = arm64.REG_SP_EL0
	p.Register["SPSel"] = arm64.REG_SPSel
	p.Register["DAIFSet"] = arm64.REG_DAIFSet
	p.Register["DAIFClr"] = arm64.REG_DAIFClr

	// Conditional operators, like EQ, NE, etc.
	p.Register["EQ"] = arm64.COND_EQ
	p.Register["NE"] = arm64.COND_NE
	p.Register["HS"] = arm64.COND_HS
	p.Register["LO"] = arm64.COND_LO
	p.Register["MI"] = arm64.COND_MI
	p.Register["PL"] = arm64.COND_PL
	p.Register["VS"] = arm64.COND_VS
	p.Register["VC"] = arm64.COND_VC
	p.Register["HI"] = arm64.COND_HI
	p.Register["LS"] = arm64.COND_LS
	p.Register["GE"] = arm64.COND_GE
	p.Register["LT"] = arm64.COND_LT
	p.Register["GT"] = arm64.COND_GT
	p.Register["LE"] = arm64.COND_LE
	p.Register["AL"] = arm64.COND_AL
	p.Register["NV"] = arm64.COND_NV

	// Pseudo-registers.
	p.Register["SB"] = RSB
	p.Register["FP"] = RFP
	p.Register["PC"] = RPC
	p.Register["SP"] = RSP

	// 保留 R28 给 g 使用(TODO: 可以删除)
	// Avoid unintentionally clobbering g using R28.
	delete(p.Register, "R28")
	p.Register["g"] = arm64.REG_R28

	// 有通用前缀名的寄存器
	p.RegisterNumber = arm64RegisterNumber
	p.RegisterPrefix = map[string]bool{
		"F": true,
		"R": true,
		"V": true,
	}

	// 注册 Plan9 汇编语言通用指令
	for i, s := range obj.Anames {
		p.Instructions[s] = i
	}

	// 注册 ARM64 指令
	for i, s := range arm64.Anames {
		if i >= obj.A_ARCHSPECIFIC {
			p.Instructions[s] = i + obj.ABaseARM64
		}
	}

	// 一些有别名的指令
	// Annoying aliases.
	p.Instructions["B"] = arm64.AB
	p.Instructions["BL"] = arm64.ABL

	// 哪些指令是只有一个目的操作数的, 比如某些一元运算
	p.UnaryDst = arm64.Linkarm64.UnaryDst

	// 判断是否为 ARM64 平台的 Jump 指令
	p.IsJump = jumpArm64

	return p
}

func archArm(CPU CPUType) *Arch {
	register := make(map[string]int16)
	// Create maps for easy lookup of instruction names etc.
	// Note that there is no list of names as there is for x86.
	for i := arm.REG_R0; i < arm.REG_SPSR; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	// Avoid unintentionally clobbering g using R10.
	delete(register, "R10")
	register["g"] = arm.REG_R10
	for i := 0; i < 16; i++ {
		register[fmt.Sprintf("C%d", i)] = int16(i)
	}

	// Pseudo-registers.
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC
	register["SP"] = RSP
	registerPrefix := map[string]bool{
		"F": true,
		"R": true,
	}

	// special operands for DMB/DSB instructions
	register["MB_SY"] = arm.REG_MB_SY
	register["MB_ST"] = arm.REG_MB_ST
	register["MB_ISH"] = arm.REG_MB_ISH
	register["MB_ISHST"] = arm.REG_MB_ISHST
	register["MB_NSH"] = arm.REG_MB_NSH
	register["MB_NSHST"] = arm.REG_MB_NSHST
	register["MB_OSH"] = arm.REG_MB_OSH
	register["MB_OSHST"] = arm.REG_MB_OSHST

	instructions := make(map[string]int)
	for i, s := range obj.Anames {
		instructions[s] = i
	}
	for i, s := range arm.Anames {
		if i >= obj.A_ARCHSPECIFIC {
			instructions[s] = i + obj.ABaseARM
		}
	}
	// Annoying aliases.
	instructions["B"] = obj.AJMP
	instructions["BL"] = obj.ACALL
	// MCR differs from MRC by the way fields of the word are encoded.
	// (Details in arm.go). Here we add the instruction so parse will find
	// it, but give it an opcode number known only to us.
	instructions["MCR"] = aMCR

	return &Arch{
		CPU:            CPU,
		Instructions:   instructions,
		Register:       register,
		RegisterPrefix: registerPrefix,
		RegisterNumber: armRegisterNumber,
		UnaryDst:       arm.Linkarm.UnaryDst,
		IsJump:         jumpArm,
	}
}

func archRISCV64(shared bool) *Arch {
	register := make(map[string]int16)

	// Standard register names.
	for i := riscv.REG_X0; i <= riscv.REG_X31; i++ {
		// Disallow X3 in shared mode, as this will likely be used as the
		// GP register, which could result in problems in non-Go code,
		// including signal handlers.
		if shared && i == riscv.REG_GP {
			continue
		}
		if i == riscv.REG_TP || i == riscv.REG_G {
			continue
		}
		name := fmt.Sprintf("X%d", i-riscv.REG_X0)
		register[name] = int16(i)
	}
	for i := riscv.REG_F0; i <= riscv.REG_F31; i++ {
		name := fmt.Sprintf("F%d", i-riscv.REG_F0)
		register[name] = int16(i)
	}
	for i := riscv.REG_V0; i <= riscv.REG_V31; i++ {
		name := fmt.Sprintf("V%d", i-riscv.REG_V0)
		register[name] = int16(i)
	}

	// General registers with ABI names.
	register["ZERO"] = riscv.REG_ZERO
	register["RA"] = riscv.REG_RA
	register["SP"] = riscv.REG_SP
	register["GP"] = riscv.REG_GP
	register["TP"] = riscv.REG_TP
	register["T0"] = riscv.REG_T0
	register["T1"] = riscv.REG_T1
	register["T2"] = riscv.REG_T2
	register["S0"] = riscv.REG_S0
	register["S1"] = riscv.REG_S1
	register["A0"] = riscv.REG_A0
	register["A1"] = riscv.REG_A1
	register["A2"] = riscv.REG_A2
	register["A3"] = riscv.REG_A3
	register["A4"] = riscv.REG_A4
	register["A5"] = riscv.REG_A5
	register["A6"] = riscv.REG_A6
	register["A7"] = riscv.REG_A7
	register["S2"] = riscv.REG_S2
	register["S3"] = riscv.REG_S3
	register["S4"] = riscv.REG_S4
	register["S5"] = riscv.REG_S5
	register["S6"] = riscv.REG_S6
	register["S7"] = riscv.REG_S7
	register["S8"] = riscv.REG_S8
	register["S9"] = riscv.REG_S9
	register["S10"] = riscv.REG_S10
	// Skip S11 as it is the g register.
	register["T3"] = riscv.REG_T3
	register["T4"] = riscv.REG_T4
	register["T5"] = riscv.REG_T5
	register["T6"] = riscv.REG_T6

	// Go runtime register names.
	register["g"] = riscv.REG_G
	register["CTXT"] = riscv.REG_CTXT
	register["TMP"] = riscv.REG_TMP

	// ABI names for floating point register.
	register["FT0"] = riscv.REG_FT0
	register["FT1"] = riscv.REG_FT1
	register["FT2"] = riscv.REG_FT2
	register["FT3"] = riscv.REG_FT3
	register["FT4"] = riscv.REG_FT4
	register["FT5"] = riscv.REG_FT5
	register["FT6"] = riscv.REG_FT6
	register["FT7"] = riscv.REG_FT7
	register["FS0"] = riscv.REG_FS0
	register["FS1"] = riscv.REG_FS1
	register["FA0"] = riscv.REG_FA0
	register["FA1"] = riscv.REG_FA1
	register["FA2"] = riscv.REG_FA2
	register["FA3"] = riscv.REG_FA3
	register["FA4"] = riscv.REG_FA4
	register["FA5"] = riscv.REG_FA5
	register["FA6"] = riscv.REG_FA6
	register["FA7"] = riscv.REG_FA7
	register["FS2"] = riscv.REG_FS2
	register["FS3"] = riscv.REG_FS3
	register["FS4"] = riscv.REG_FS4
	register["FS5"] = riscv.REG_FS5
	register["FS6"] = riscv.REG_FS6
	register["FS7"] = riscv.REG_FS7
	register["FS8"] = riscv.REG_FS8
	register["FS9"] = riscv.REG_FS9
	register["FS10"] = riscv.REG_FS10
	register["FS11"] = riscv.REG_FS11
	register["FT8"] = riscv.REG_FT8
	register["FT9"] = riscv.REG_FT9
	register["FT10"] = riscv.REG_FT10
	register["FT11"] = riscv.REG_FT11

	// Pseudo-registers.
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC

	instructions := make(map[string]int)
	for i, s := range obj.Anames {
		instructions[s] = i
	}
	for i, s := range riscv.Anames {
		if i >= obj.A_ARCHSPECIFIC {
			instructions[s] = i + obj.ABaseRISCV
		}
	}

	nilRegisterNumber := func(name string, n int16) (int16, bool) {
		return 0, false
	}
	jumpRISCV := func(word string) bool {
		switch word {
		case "BEQ", "BEQZ", "BGE", "BGEU", "BGEZ", "BGT", "BGTU", "BGTZ", "BLE", "BLEU", "BLEZ",
			"BLT", "BLTU", "BLTZ", "BNE", "BNEZ", "CALL", "JAL", "JALR", "JMP":
			return true
		}
		return false
	}

	return &Arch{
		//LinkArch:       &riscv.LinkRISCV64,
		Instructions:   instructions,
		Register:       register,
		RegisterPrefix: nil,
		RegisterNumber: nilRegisterNumber,
		UnaryDst:       riscv.LinkRISCV64.UnaryDst,
		IsJump:         jumpRISCV,
	}
}

func archWasm() *Arch {
	instructions := make(map[string]int)
	for i, s := range obj.Anames {
		instructions[s] = i
	}
	for i, s := range wasm.Anames {
		if i >= obj.A_ARCHSPECIFIC {
			instructions[s] = i + obj.ABaseWasm
		}
	}

	nilRegisterNumber := func(name string, n int16) (int16, bool) {
		return 0, false
	}
	jumpWasm := func(word string) bool {
		return word == "JMP" || word == "CALL" || word == "Call" || word == "Br" || word == "BrIf"
	}

	return &Arch{
		//LinkArch:       &wasm.Linkwasm,
		Instructions:   instructions,
		Register:       wasm.Register,
		RegisterPrefix: nil,
		RegisterNumber: nilRegisterNumber,
		UnaryDst:       wasm.Linkwasm.UnaryDst,
		IsJump:         jumpWasm,
	}
}
