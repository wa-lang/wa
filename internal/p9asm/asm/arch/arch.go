// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package arch

import (
	"fmt"
	"strings"

	"wa-lang.org/wa/internal/p9asm/obj"
	"wa-lang.org/wa/internal/p9asm/obj/arm"
	"wa-lang.org/wa/internal/p9asm/obj/arm64"
	"wa-lang.org/wa/internal/p9asm/obj/loong64"
	"wa-lang.org/wa/internal/p9asm/obj/riscv"
	"wa-lang.org/wa/internal/p9asm/obj/x86"
)

// 处理器类型
type CPUType byte

const (
	NoArch CPUType = iota
	AMD64
	ARM
	ARM64
	I386
	LOONG64
	RISCV64
)

// Plan9 汇编语言的伪寄存器
const (
	RFP = -(iota + 1)
	RSB
	RSP
	RPC
)

// Arch 封装了链接器的架构对象，并附加了更多与架构相关的具体信息。
type Arch struct {
	CPU CPUType

	// 连接器相关配置
	LinkArch *obj.LinkArch

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

	// 判断是否为 jump 跳转指令
	IsJump func(word string) bool
}

// 设置并返回对应的CPU对象
func Set(cpu CPUType) *Arch {
	switch cpu {
	case I386:
		return archX86(I386, &x86.Link386)
	case AMD64:
		return archX86(AMD64, &x86.Linkamd64)
	case ARM:
		return archArm(ARM, &arm.Linkarm)
	case ARM64:
		return archArm64(ARM64, &arm64.Linkarm64)
	case LOONG64:
		return archLoong64(LOONG64, &loong64.Linkloong64)
	case RISCV64:
		return archRISCV64(RISCV64, &riscv.LinkRISCV64, false)
	default:
		panic("unreachable")
	}
}

func archX86(CPU CPUType, linkArch *obj.LinkArch) *Arch {
	p := &Arch{
		CPU:      CPU,
		LinkArch: linkArch,

		Instructions:   map[string]int{},
		Register:       map[string]int16{},
		RegisterPrefix: map[string]bool{},
	}

	// 构造寄存器名查询表
	p.Register = make(map[string]int16)
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

func archArm64(CPU CPUType, linkArch *obj.LinkArch) *Arch {
	p := &Arch{
		CPU:      CPU,
		LinkArch: linkArch,

		Instructions:   map[string]int{},
		Register:       map[string]int16{},
		RegisterPrefix: map[string]bool{},
	}

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

	// 判断是否为 ARM64 平台的 Jump 指令
	p.IsJump = jumpArm64

	return p
}

func archArm(CPU CPUType, linkArch *obj.LinkArch) *Arch {
	p := &Arch{
		CPU:      CPU,
		LinkArch: linkArch,

		Instructions:   map[string]int{},
		Register:       map[string]int16{},
		RegisterPrefix: map[string]bool{},
	}

	// Create maps for easy lookup of instruction names etc.
	// Note that there is no list of names as there is for x86.
	for i := arm.REG_R0; i < arm.REG_SPSR; i++ {
		p.Register[obj.Rconv(i)] = int16(i)
	}
	// Avoid unintentionally clobbering g using R10.
	delete(p.Register, "R10")
	p.Register["g"] = arm.REG_R10
	for i := 0; i < 16; i++ {
		p.Register[fmt.Sprintf("C%d", i)] = int16(i)
	}

	// Pseudo-registers.
	p.Register["SB"] = RSB
	p.Register["FP"] = RFP
	p.Register["PC"] = RPC
	p.Register["SP"] = RSP

	p.RegisterPrefix = map[string]bool{"F": true, "R": true}
	p.RegisterNumber = armRegisterNumber

	// special operands for DMB/DSB instructions
	p.Register["MB_SY"] = arm.REG_MB_SY
	p.Register["MB_ST"] = arm.REG_MB_ST
	p.Register["MB_ISH"] = arm.REG_MB_ISH
	p.Register["MB_ISHST"] = arm.REG_MB_ISHST
	p.Register["MB_NSH"] = arm.REG_MB_NSH
	p.Register["MB_NSHST"] = arm.REG_MB_NSHST
	p.Register["MB_OSH"] = arm.REG_MB_OSH
	p.Register["MB_OSHST"] = arm.REG_MB_OSHST

	for i, s := range obj.Anames {
		p.Instructions[s] = i
	}
	for i, s := range arm.Anames {
		if i >= obj.A_ARCHSPECIFIC {
			p.Instructions[s] = i + obj.ABaseARM
		}
	}
	// Annoying aliases.
	p.Instructions["B"] = obj.AJMP
	p.Instructions["BL"] = obj.ACALL
	// MCR differs from MRC by the way fields of the word are encoded.
	// (Details in arm.go). Here we add the instruction so parse will find
	// it, but give it an opcode number known only to us.
	p.Instructions["MCR"] = aMCR

	p.IsJump = jumpArm

	return p
}

func archLoong64(CPU CPUType, linkArch *obj.LinkArch) *Arch {
	p := &Arch{
		CPU:      CPU,
		LinkArch: linkArch,

		Instructions:   map[string]int{},
		Register:       map[string]int16{},
		RegisterPrefix: map[string]bool{},
	}

	// Create maps for easy lookup of instruction names etc.
	// Note that there is no list of names as there is for x86.
	for i := loong64.REG_R0; i <= loong64.REG_R31; i++ {
		p.Register[obj.Rconv(i)] = int16(i)
	}

	for i := loong64.REG_F0; i <= loong64.REG_F31; i++ {
		p.Register[obj.Rconv(i)] = int16(i)
	}

	for i := loong64.REG_FCSR0; i <= loong64.REG_FCSR31; i++ {
		p.Register[obj.Rconv(i)] = int16(i)
	}

	for i := loong64.REG_FCC0; i <= loong64.REG_FCC31; i++ {
		p.Register[obj.Rconv(i)] = int16(i)
	}

	for i := loong64.REG_V0; i <= loong64.REG_V31; i++ {
		p.Register[obj.Rconv(i)] = int16(i)
	}

	for i := loong64.REG_X0; i <= loong64.REG_X31; i++ {
		p.Register[obj.Rconv(i)] = int16(i)
	}

	// Pseudo-registers.
	p.Register["SB"] = RSB
	p.Register["FP"] = RFP
	p.Register["PC"] = RPC

	// Avoid unintentionally clobbering g using R22.
	delete(p.Register, "R22")
	p.Register["g"] = loong64.REG_R22

	p.RegisterPrefix = map[string]bool{
		"F":    true,
		"FCSR": true,
		"FCC":  true,
		"R":    true,
		"V":    true,
		"X":    true,
	}
	p.RegisterNumber = loong64RegisterNumber

	for i, s := range obj.Anames {
		p.Instructions[s] = i
	}
	for i, s := range loong64.Anames {
		if i >= obj.A_ARCHSPECIFIC {
			p.Instructions[s] = i + obj.ABaseLoong64
		}
	}

	// Annoying alias.
	p.Instructions["JAL"] = loong64.AJAL

	p.IsJump = jumpLoong64

	return p
}

func archRISCV64(CPU CPUType, linkArch *obj.LinkArch, shared bool) *Arch {
	p := &Arch{
		CPU:      CPU,
		LinkArch: linkArch,

		Instructions:   map[string]int{},
		Register:       map[string]int16{},
		RegisterPrefix: map[string]bool{},
	}

	// Standard register names.
	for i := riscv.REG_X0; i <= riscv.REG_X31; i++ {
		// Disallow X3 in shared mode, as this will likely be used as the
		// GP register, which could result in problems in non-Wa code,
		// including signal handlers.
		if shared && i == riscv.REG_GP {
			continue
		}
		if i == riscv.REG_TP || i == riscv.REG_G {
			continue
		}
		name := fmt.Sprintf("X%d", i-riscv.REG_X0)
		p.Register[name] = int16(i)
	}
	for i := riscv.REG_F0; i <= riscv.REG_F31; i++ {
		name := fmt.Sprintf("F%d", i-riscv.REG_F0)
		p.Register[name] = int16(i)
	}
	for i := riscv.REG_V0; i <= riscv.REG_V31; i++ {
		name := fmt.Sprintf("V%d", i-riscv.REG_V0)
		p.Register[name] = int16(i)
	}

	// General registers with ABI names.
	p.Register["ZERO"] = riscv.REG_ZERO
	p.Register["RA"] = riscv.REG_RA
	p.Register["SP"] = riscv.REG_SP
	p.Register["GP"] = riscv.REG_GP
	p.Register["TP"] = riscv.REG_TP
	p.Register["T0"] = riscv.REG_T0
	p.Register["T1"] = riscv.REG_T1
	p.Register["T2"] = riscv.REG_T2
	p.Register["S0"] = riscv.REG_S0
	p.Register["S1"] = riscv.REG_S1
	p.Register["A0"] = riscv.REG_A0
	p.Register["A1"] = riscv.REG_A1
	p.Register["A2"] = riscv.REG_A2
	p.Register["A3"] = riscv.REG_A3
	p.Register["A4"] = riscv.REG_A4
	p.Register["A5"] = riscv.REG_A5
	p.Register["A6"] = riscv.REG_A6
	p.Register["A7"] = riscv.REG_A7
	p.Register["S2"] = riscv.REG_S2
	p.Register["S3"] = riscv.REG_S3
	p.Register["S4"] = riscv.REG_S4
	p.Register["S5"] = riscv.REG_S5
	p.Register["S6"] = riscv.REG_S6
	p.Register["S7"] = riscv.REG_S7
	p.Register["S8"] = riscv.REG_S8
	p.Register["S9"] = riscv.REG_S9
	p.Register["S10"] = riscv.REG_S10
	// Skip S11 as it is the g register.
	p.Register["T3"] = riscv.REG_T3
	p.Register["T4"] = riscv.REG_T4
	p.Register["T5"] = riscv.REG_T5
	p.Register["T6"] = riscv.REG_T6

	// Go runtime register names.
	p.Register["g"] = riscv.REG_G
	p.Register["CTXT"] = riscv.REG_CTXT
	p.Register["TMP"] = riscv.REG_TMP

	// ABI names for floating point register.
	p.Register["FT0"] = riscv.REG_FT0
	p.Register["FT1"] = riscv.REG_FT1
	p.Register["FT2"] = riscv.REG_FT2
	p.Register["FT3"] = riscv.REG_FT3
	p.Register["FT4"] = riscv.REG_FT4
	p.Register["FT5"] = riscv.REG_FT5
	p.Register["FT6"] = riscv.REG_FT6
	p.Register["FT7"] = riscv.REG_FT7
	p.Register["FS0"] = riscv.REG_FS0
	p.Register["FS1"] = riscv.REG_FS1
	p.Register["FA0"] = riscv.REG_FA0
	p.Register["FA1"] = riscv.REG_FA1
	p.Register["FA2"] = riscv.REG_FA2
	p.Register["FA3"] = riscv.REG_FA3
	p.Register["FA4"] = riscv.REG_FA4
	p.Register["FA5"] = riscv.REG_FA5
	p.Register["FA6"] = riscv.REG_FA6
	p.Register["FA7"] = riscv.REG_FA7
	p.Register["FS2"] = riscv.REG_FS2
	p.Register["FS3"] = riscv.REG_FS3
	p.Register["FS4"] = riscv.REG_FS4
	p.Register["FS5"] = riscv.REG_FS5
	p.Register["FS6"] = riscv.REG_FS6
	p.Register["FS7"] = riscv.REG_FS7
	p.Register["FS8"] = riscv.REG_FS8
	p.Register["FS9"] = riscv.REG_FS9
	p.Register["FS10"] = riscv.REG_FS10
	p.Register["FS11"] = riscv.REG_FS11
	p.Register["FT8"] = riscv.REG_FT8
	p.Register["FT9"] = riscv.REG_FT9
	p.Register["FT10"] = riscv.REG_FT10
	p.Register["FT11"] = riscv.REG_FT11

	// Pseudo-registers.
	p.Register["SB"] = RSB
	p.Register["FP"] = RFP
	p.Register["PC"] = RPC

	for i, s := range obj.Anames {
		p.Instructions[s] = i
	}
	for i, s := range riscv.Anames {
		if i >= obj.A_ARCHSPECIFIC {
			p.Instructions[s] = i + obj.ABaseRISCV
		}
	}

	p.RegisterNumber = func(name string, n int16) (int16, bool) {
		return 0, false
	}

	p.IsJump = func(word string) bool {
		switch word {
		case "BEQ", "BEQZ", "BGE", "BGEU", "BGEZ", "BGT", "BGTU", "BGTZ", "BLE", "BLEU", "BLEZ",
			"BLT", "BLTU", "BLTZ", "BNE", "BNEZ", "CALL", "JAL", "JALR", "JMP":
			return true
		}
		return false
	}

	return p
}
