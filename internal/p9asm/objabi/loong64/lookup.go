// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

import (
	"fmt"

	"wa-lang.org/wa/internal/p9asm/objabi"
)

// 寄存器名字列表
var Register = make([]string, RegMax-RBase)

func init() {
	// 初始化寄存器名字表格
	for i := RBase; i < RegMax; i++ {
		Register[i-RBase] = RegString(i)
	}
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
	switch {
	case r == 0:
		return "NONE"
	case REG_R0 <= r && r <= REG_R31:
		return fmt.Sprintf("R%d", r-REG_R0)
	case REG_F0 <= r && r <= REG_F31:
		return fmt.Sprintf("F%d", r-REG_F0)
	case REG_FCSR0 <= r && r <= REG_FCSR31:
		return fmt.Sprintf("FCSR%d", r-REG_FCSR0)
	case REG_FCC0 <= r && r <= REG_FCC31:
		return fmt.Sprintf("FCC%d", r-REG_FCC0)
	case REG_V0 <= r && r <= REG_V31:
		return fmt.Sprintf("V%d", r-REG_V0)
	case REG_X0 <= r && r <= REG_X31:
		return fmt.Sprintf("X%d", r-REG_X0)
	}

	// bits 0-4 indicates register: Vn or Xn
	// bits 5-9 indicates arrangement: <T>
	// bits 10 indicates SMID type: 0: LSX, 1: LASX
	simd_type := (int16(r) >> EXT_SIMDTYPE_SHIFT) & EXT_SIMDTYPE_MASK
	reg_num := (int16(r) >> EXT_REG_SHIFT) & EXT_REG_MASK
	arng_type := (int16(r) >> EXT_TYPE_SHIFT) & EXT_TYPE_MASK
	reg_prefix := "#"
	switch simd_type {
	case LSX:
		reg_prefix = "V"
	case LASX:
		reg_prefix = "X"
	}

	switch {
	case REG_ARNG <= r && r < REG_ELEM:
		return fmt.Sprintf("%s%d.%s", reg_prefix, reg_num, arrange(arng_type))

	case REG_ELEM <= r && r < REG_ELEM_END:
		return fmt.Sprintf("%s%d.%s", reg_prefix, reg_num, arrange(arng_type))
	}

	return fmt.Sprintf("loong64.badreg(%d)", r)
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
	panic("unreachable")
}

func arrange(a int16) string {
	switch a {
	case ARNG_32B:
		return "B32"
	case ARNG_16H:
		return "H16"
	case ARNG_8W:
		return "W8"
	case ARNG_4V:
		return "V4"
	case ARNG_2Q:
		return "Q2"
	case ARNG_16B:
		return "B16"
	case ARNG_8H:
		return "H8"
	case ARNG_4W:
		return "W4"
	case ARNG_2V:
		return "V2"
	case ARNG_B:
		return "B"
	case ARNG_H:
		return "H"
	case ARNG_W:
		return "W"
	case ARNG_V:
		return "V"
	case ARNG_BU:
		return "BU"
	case ARNG_HU:
		return "HU"
	case ARNG_WU:
		return "WU"
	case ARNG_VU:
		return "VU"
	default:
		return "ARNG_???"
	}
}
