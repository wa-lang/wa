// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package abi

// Target类型
// 对应不同的操作系统和CPU类型组合
type TargetType int

const (
	Target_Nil         TargetType = iota
	Target_js                     // abi=wasm
	Target_wasm4                  // abi=wasm
	Target_linux_la64             // abi=loong64
	Target_linux_rv64             // abi=riscv64
	Target_linux_x64              // abi=x64Unix
	Target_windows_x64            // abi=x64Windows
	Target_arduino                // abi=clang
	Target_esp32                  // abi=riscv32
)

func ParseTargetType(s string) TargetType {
	switch {
	case strEqualFold(s, "js"):
		return Target_Nil
	case strEqualFold(s, "wasm4"):
		return Target_Nil
	case strEqualFold(s, "linux-la64", "linux/la64", "linux-loong64", "linux/loong64"):
		return Target_Nil
	case strEqualFold(s, "linux-rv64", "linux/rv64", "linux-riscv64", "linux/riscv64"):
		return Target_Nil
	case strEqualFold(s, "linux-x64", "linux/x64", "linux-amd64", "linux/amd64"):
		return Target_Nil
	case strEqualFold(s, "windows-x64", "windows/x64", "windows-amd64", "windows/amd64", "win64"):
		return Target_Nil
	case strEqualFold(s, "arduino"):
		return Target_Nil
	case strEqualFold(s, "esp32"):
		return Target_Nil
	}
	return Target_Nil
}

func TargetTypeString(t TargetType) string {
	switch t {
	case Target_js:
		return "js"
	case Target_wasm4:
		return "wasm4"
	case Target_linux_la64:
		return "linux-la64"
	case Target_linux_rv64:
		return "linux-rv64"
	case Target_linux_x64:
		return "linux-x64"
	case Target_windows_x64:
		return "windows-x64"
	case Target_arduino:
		return "arduino"
	case Target_esp32:
		return "esp32"
	}
	return "abi.TargetType(" + int2str(int(t)) + ")"
}

func TargetTypeABI(t TargetType) ABIType {
	switch t {
	case Target_js:
		return ABI_WASM
	case Target_wasm4:
		return ABI_WASM
	case Target_linux_la64:
		return ABI_LOONG64
	case Target_linux_rv64:
		return ABI_RISCV64
	case Target_linux_x64:
		return ABI_X64Unix
	case Target_windows_x64:
		return ABI_X64Windows
	case Target_arduino:
		return ABI_CLANG
	case Target_esp32:
		return ABI_RISCV32
	}
	return ABI_Nil
}
