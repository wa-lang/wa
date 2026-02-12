// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

// allKeys records the list of all possible child keys for actions that support "any".
var allKeys = map[string][]string{
	"is64":     {"0", "1"},
	"ismem":    {"0", "1"},
	"addrsize": {"16", "32", "64"},
	"datasize": {"16", "32", "64"},
}

// isPrefix records the x86 opcode prefix bytes.
var isPrefix = map[string]bool{
	"26": true,
	"2E": true,
	"36": true,
	"3E": true,
	"64": true,
	"65": true,
	"66": true,
	"67": true,
	"F0": true,
	"F2": true,
	"F3": true,
}

// usesReg records the argument codes that use the modrm reg field.
var usesReg = map[string]bool{
	"r8":  true,
	"r16": true,
	"r32": true,
	"r64": true,
}

// usesRM records the argument codes that use the modrm r/m field.
var usesRM = map[string]bool{
	"r/m8":  true,
	"r/m16": true,
	"r/m32": true,
	"r/m64": true,
}

var isVexEncodablePrefix = map[string]bool{
	"0F":   true,
	"0F38": true,
	"0F3A": true,
	"66":   true,
	"F3":   true,
	"F2":   true,
}
