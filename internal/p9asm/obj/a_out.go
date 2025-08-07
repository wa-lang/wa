// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package obj

// ARM scond byte
const (
	C_SCOND     = (1 << 4) - 1
	C_SBIT      = 1 << 4
	C_PBIT      = 1 << 5
	C_WBIT      = 1 << 6
	C_FBIT      = 1 << 7
	C_UBIT      = 1 << 7
	C_SCOND_XOR = 14
)

// ARM条件指令
// 比如 ADD 指令通过添加 .NE 后缀变成条件支持 ADD.NE
var armCondCode = []string{
	".EQ",
	".NE",
	".CS",
	".CC",
	".MI",
	".PL",
	".VS",
	".VC",
	".HI",
	".LS",
	".GE",
	".LT",
	".GT",
	".LE",
	"",
	".NV",
}

// CConv formats ARM condition codes.
func CConv(s uint8) string {
	if s == 0 {
		return ""
	}
	sc := armCondCode[(s&C_SCOND)^C_SCOND_XOR]
	if s&C_SBIT != 0 {
		sc += ".S"
	}
	if s&C_PBIT != 0 {
		sc += ".P"
	}
	if s&C_WBIT != 0 {
		sc += ".W"
	}
	if s&C_UBIT != 0 { // ambiguous with FBIT
		sc += ".U"
	}
	return sc
}

const (
	StackBig   = 4096 // 大帧判断阈值
	StackSmall = 128  // 最小预留空间
	StackGuard = 640  // 栈溢出保护区

	StackLimit = StackGuard - StackSmall // 剩余触发点, 可触发自动报告或任务失败处理
)

// Wa obj 文件魔数
// Writeobjdirect 函数使用
const (
	MagicHeader = "\x00\x00wa01ld"
	MagicFooter = "\xff\xffwa01ld"

	MagicSymbolStart = 0xfe
	MagicFooterStart = 0xff
)

// symbol version, incremented each time a file is loaded.
// version==1 is reserved for savehist.
const (
	HistVersion = 1
)
