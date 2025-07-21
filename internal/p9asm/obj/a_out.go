// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package obj

import "fmt"

// 各个平台通用的指令
// 平台特殊的指令从 A_ARCHSPECIFIC 开始定义
const (
	AXXX = 0 + iota
	ACALL
	ACHECKNIL
	ADATA
	ADUFFCOPY
	ADUFFZERO
	AEND
	AFUNCDATA
	AGLOBL
	AJMP
	ANOP
	APCDATA
	ARET
	ATEXT
	ATYPE
	AUNDEF
	AUSEFIELD
	AVARDEF
	AVARKILL
	A_ARCHSPECIFIC
)

// 指令的名字
var Anames = []string{
	AXXX:      "XXX",
	ACALL:     "CALL",
	ACHECKNIL: "CHECKNIL",
	ADATA:     "DATA",
	ADUFFCOPY: "DUFFCOPY",
	ADUFFZERO: "DUFFZERO",
	AEND:      "END",
	AFUNCDATA: "FUNCDATA",
	AGLOBL:    "GLOBL",
	AJMP:      "JMP",
	ANOP:      "NOP",
	APCDATA:   "PCDATA",
	ARET:      "RET",
	ATEXT:     "TEXT",
	ATYPE:     "TYPE",
	AUNDEF:    "UNDEF",
	AUSEFIELD: "USEFIELD",
	AVARDEF:   "VARDEF",
	AVARKILL:  "VARKILL",
}

// 每个平台有独立的指令空间
// 比如 ABaseAMD64 + A_ARCHSPECIFIC 开始的是 AMD64 特有的指令
const (
	ABase386 = (1 + iota) << 11
	ABaseARM
	ABaseAMD64
	ABaseARM64
	ABaseLoong64
	ABaseRISCV

	AllowedOpCodes = 1 << 11            // The number of opcodes available for any given architecture.
	AMask          = AllowedOpCodes - 1 // AND with this to use the opcode as an array index.
)

// 每个平台寄存也有独立的空间
// 比如 RBaseAMD64 开始的是 AMD64 平台的特有的寄存器
const (
	REG_NONE     = 0 // 寄存器编号为空
	RBase386     = 1 * 1024
	RBaseAMD64   = 2 * 1024
	RBaseARM     = 3 * 1024
	RBaseARM64   = 8 * 1024  // range [8k, 13k)
	RBaseRISCV   = 15 * 1024 // range [15k, 16k)
	RBaseLOONG64 = 19 * 1024 // range [19K, 22k)
)

// 机器码和对应的名字
type opSet struct {
	lo    int
	names []string
}

// 机器码区间集合
var aSpace []opSet

// 注册不同平台的指令区间
func RegisterOpcode(lo int, Anames []string) {
	aSpace = append(aSpace, opSet{lo, Anames})
}

// 指令机器码转为字符串名字
func Aconv(a int) string {
	if a < A_ARCHSPECIFIC {
		return Anames[a]
	}
	for i := range aSpace {
		as := &aSpace[i]
		if as.lo <= a && a < as.lo+len(as.names) {
			return as.names[a-as.lo]
		}
	}
	return fmt.Sprintf("A???%d", a)
}

// 寄存器区间
type regSet struct {
	lo    int              // 开始
	hi    int              // 结束(开区间)
	Rconv func(int) string // 用于打印
}

// 用于注册不同架构下的寄存器
// 不同架构的寄存器编号处于不同空间不会冲突
var regSpace []regSet

// 注册不同平台的寄存器区间
func RegisterRegister(lo, hi int, Rconv func(int) string) {
	regSpace = append(regSpace, regSet{lo, hi, Rconv})
}

// 将寄存器编号转为字符串名字
func Rconv(reg int) string {
	if reg == REG_NONE {
		return "NONE"
	}
	for i := range regSpace {
		rs := &regSpace[i]
		if rs.lo <= reg && reg < rs.hi {
			return rs.Rconv(reg)
		}
	}
	return fmt.Sprintf("R???%d", reg)
}
