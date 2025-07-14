// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file encapsulates some of the odd characteristics of the
// Loong64 (LoongArch64) instruction set, to minimize its interaction
// with the core of the assembler.

package arch

import (
	"wa-lang.org/wa/internal/p9asm/obj/loong64"
)

func jumpLoong64(word string) bool {
	switch word {
	case "BEQ", "BFPF", "BFPT", "BLTZ", "BGEZ", "BLEZ", "BGTZ", "BLT", "BLTU", "JIRL", "BNE", "BGE", "BGEU", "JMP", "JAL", "CALL":
		return true
	}
	return false
}

// IsLoong64MUL reports whether the op (as defined by an loong64.A* constant) is
// one of the MUL/DIV/REM instructions that require special handling.
func IsLoong64MUL(op int) bool {
	switch op {
	case loong64.AMUL, loong64.AMULU, loong64.AMULV, loong64.AMULVU,
		loong64.ADIV, loong64.ADIVU, loong64.ADIVV, loong64.ADIVVU,
		loong64.AREM, loong64.AREMU, loong64.AREMV, loong64.AREMVU:
		return true
	}
	return false
}

// IsLoong64RDTIME reports whether the op (as defined by an loong64.A*
// constant) is one of the RDTIMELW/RDTIMEHW/RDTIMED instructions that
// require special handling.
func IsLoong64RDTIME(op int) bool {
	switch op {
	case loong64.ARDTIMELW, loong64.ARDTIMEHW, loong64.ARDTIMED:
		return true
	}
	return false
}

func IsLoong64AMO(op int) bool {
	return loong64.IsAtomicInst(op)
}

var loong64ElemExtMap = map[string]int16{
	"B":  loong64.ARNG_B,
	"H":  loong64.ARNG_H,
	"W":  loong64.ARNG_W,
	"V":  loong64.ARNG_V,
	"BU": loong64.ARNG_BU,
	"HU": loong64.ARNG_HU,
	"WU": loong64.ARNG_WU,
	"VU": loong64.ARNG_VU,
}

var loong64LsxArngExtMap = map[string]int16{
	"B16": loong64.ARNG_16B,
	"H8":  loong64.ARNG_8H,
	"W4":  loong64.ARNG_4W,
	"V2":  loong64.ARNG_2V,
}

var loong64LasxArngExtMap = map[string]int16{
	"B32": loong64.ARNG_32B,
	"H16": loong64.ARNG_16H,
	"W8":  loong64.ARNG_8W,
	"V4":  loong64.ARNG_4V,
	"Q2":  loong64.ARNG_2Q,
}

func loong64RegisterNumber(name string, n int16) (int16, bool) {
	switch name {
	case "F":
		if 0 <= n && n <= 31 {
			return loong64.REG_F0 + n, true
		}
	case "FCSR":
		if 0 <= n && n <= 31 {
			return loong64.REG_FCSR0 + n, true
		}
	case "FCC":
		if 0 <= n && n <= 31 {
			return loong64.REG_FCC0 + n, true
		}
	case "R":
		if 0 <= n && n <= 31 {
			return loong64.REG_R0 + n, true
		}
	case "V":
		if 0 <= n && n <= 31 {
			return loong64.REG_V0 + n, true
		}
	case "X":
		if 0 <= n && n <= 31 {
			return loong64.REG_X0 + n, true
		}
	}
	return 0, false
}
