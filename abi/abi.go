// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 凹语言ABI和Target定义
package abi

// ABI类型
// 对应不同的汇编指令和调用约定
type ABIType int

const (
	ABI_Nil ABIType = iota
	ABI_WASM
	ABI_LOONG64
	ABI_RISCV64
	ABI_RISCV32
	ABI_X64Unix    // UNIX System V ABI
	ABI_X64Windows // Windows ABI
	ABI_CLANG
	ABI_Max
)

func ParseABIType(s string) ABIType {
	switch {
	case strEqualFold(s, "wasm"):
		return ABI_WASM
	case strEqualFold(s, "loong64", "la64"):
		return ABI_LOONG64
	case strEqualFold(s, "riscv64", "rv64"):
		return ABI_RISCV64
	case strEqualFold(s, "riscv32", "rv32"):
		return ABI_RISCV32
	case strEqualFold(s, "x64unix"):
		return ABI_X64Unix
	case strEqualFold(s, "x64windows", "win64"):
		return ABI_X64Windows
	case strEqualFold(s, "clang", "c"):
		return ABI_CLANG
	}
	return ABI_Nil
}

func ABITypeString(t ABIType) string {
	switch t {
	case ABI_WASM:
		return "wasm"
	case ABI_LOONG64:
		return "loong64"
	case ABI_RISCV64:
		return "riscv64"
	case ABI_RISCV32:
		return "riscv32"
	case ABI_X64Unix: // UNIX System V ABI
		return "x64unix"
	case ABI_X64Windows: // Windows ABI
		return "x64windows"
	case ABI_CLANG:
		return "clang"
	}
	return "abi.ABIType(" + int2str(int(t)) + ")"
}
