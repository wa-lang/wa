// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x86

import (
	"fmt"

	"wa-lang.org/wa/internal/p9asm/objabi"
)

var Register = []string{
	"AL", // [D_AL]
	"CL",
	"DL",
	"BL",
	"SPB",
	"BPB",
	"SIB",
	"DIB",
	"R8B",
	"R9B",
	"R10B",
	"R11B",
	"R12B",
	"R13B",
	"R14B",
	"R15B",
	"AX", // [D_AX]
	"CX",
	"DX",
	"BX",
	"SP",
	"BP",
	"SI",
	"DI",
	"R8",
	"R9",
	"R10",
	"R11",
	"R12",
	"R13",
	"R14",
	"R15",
	"AH",
	"CH",
	"DH",
	"BH",
	"F0", // [D_F0]
	"F1",
	"F2",
	"F3",
	"F4",
	"F5",
	"F6",
	"F7",
	"M0",
	"M1",
	"M2",
	"M3",
	"M4",
	"M5",
	"M6",
	"M7",
	"X0",
	"X1",
	"X2",
	"X3",
	"X4",
	"X5",
	"X6",
	"X7",
	"X8",
	"X9",
	"X10",
	"X11",
	"X12",
	"X13",
	"X14",
	"X15",
	"CS", // [D_CS]
	"SS",
	"DS",
	"ES",
	"FS",
	"GS",
	"GDTR", // [D_GDTR]
	"IDTR", // [D_IDTR]
	"LDTR", // [D_LDTR]
	"MSW",  // [D_MSW]
	"TASK", // [D_TASK]
	"CR0",  // [D_CR]
	"CR1",
	"CR2",
	"CR3",
	"CR4",
	"CR5",
	"CR6",
	"CR7",
	"CR8",
	"CR9",
	"CR10",
	"CR11",
	"CR12",
	"CR13",
	"CR14",
	"CR15",
	"DR0", // [D_DR]
	"DR1",
	"DR2",
	"DR3",
	"DR4",
	"DR5",
	"DR6",
	"DR7",
	"TR0", // [D_TR]
	"TR1",
	"TR2",
	"TR3",
	"TR4",
	"TR5",
	"TR6",
	"TR7",
	"TLS", // [D_TLS]
}

// 根据名字查找寄存器, 失败返回 objabi.REG_NONE
func LookupRegister(regName string) objabi.RBaseType {
	for i, s := range Register {
		if s == regName {
			return RBase + objabi.RBaseType(i)
		}
	}
	return objabi.REG_NONE
}

// 寄存器转字符串格式
func RegString(r objabi.RBaseType) string {
	if RBase <= r && r < RegMax {
		return Register[r-RBase]
	}
	return fmt.Sprintf("x86.Register(%d)", r-RBase)
}

// 根据名字查找汇编指令, 失败返回 objabi.AXXX
func LookupAs(asName string) objabi.As {
	for i, s := range Anames {
		if s == asName {
			return ABase + objabi.As(i)
		}
	}
	return objabi.AXXX
}

// 汇编指令转字符串格式
func AsString(as objabi.As) string {
	if ABase <= as && as < AsMax {
		return Anames[as-ABase]
	}
	return fmt.Sprintf("x86.As(%d)", as-ABase)
}
