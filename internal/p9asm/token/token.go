// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"fmt"
	"strconv"

	"wa-lang.org/wa/internal/p9asm/objabi"
	"wa-lang.org/wa/internal/p9asm/objabi/x86"
)

// 记号类型
type Token int

const (
	ILLEGAL Token = iota // 非法
	EOF                  // 结尾
	COMMENT              // 注释

	IDENT  // 标识符
	REG    // 寄存器
	INT    // 12345
	FLOAT  // 123.45
	CHAR   // 'a'
	STRING // "abc"

	LSH // << 左移
	RSH // >> 逻辑右移
	ARR // -> 数学右移, 用于 ARM 的第 3 类移动指令
	ROT // @> 循环右移, 用于 ARM 的第 4 类移动指令

	// 特殊指令和寄存器开始
	objabi_base

	// 汇编指令 开始
	instruction_beg = Token(objabi.ABase)
	instruction_end = Token(objabi.ABaseMax)

	// 寄存器开始
	register_beg = Token(objabi.RBase)
	register_end = Token(objabi.RBaseMax)
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	REG:    "REG",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	STRING: "STRING",

	LSH: "<<",
	RSH: ">>",
	ARR: "->",
	ROT: "@>",
}

func (tok Token) String() string {
	// 普通 token
	if 0 <= tok && tok < Token(len(tokens)) {
		return tokens[tok]
	}

	// 基础汇编指令
	if objabi.ABase <= tok && tok < Token(objabi.A_ARCHSPECIFIC) {
		return objabi.Anames[tok]
	}

	// 不同平台的汇编指令
	if tok.IsInstruction() {
		switch {
		case tok >= objabi.ABaseAMD64 && tok < Token(objabi.ABaseAMD64+len(x86.Anames)):
			return x86.Anames[tok-Token(objabi.ABaseAMD64+len(x86.Anames))]
		default:
			panic("TODO")
		}
	}

	// 不同平台的寄存器
	if tok.IsReginster() {
		switch {
		case tok >= Token(objabi.RBaseAMD64) && tok < Token(x86.RegMax):
			return x86.RegString(objabi.RBaseType(tok))
		default:
			panic("TODO")
		}
	}

	return "token(" + strconv.Itoa(int(tok)) + ")"
}

func Lookup(arch objabi.CPUType, ident string) Token {
	switch arch {
	case objabi.X86:
	case objabi.AMD64:
		if x := x86.LoopupRegister(ident); x != objabi.REG_NONE {
			return Token(x)
		}
		// TODO: 指令
		// TODO: 寄存器
	case objabi.ARM:
	case objabi.ARM64:
	case objabi.Loong64:
	case objabi.RISCV:
	default:
		panic(fmt.Sprintf("invalid arch: %v", arch))
	}

	return ILLEGAL
}

func (tok Token) IsInstruction() bool {
	return instruction_beg < tok && tok < instruction_end
}

func (tok Token) IsReginster() bool {
	return register_beg < tok && tok < register_end
}
