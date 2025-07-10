// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import "wa-lang.org/wa/internal/p9asm/obj"

var Register = map[string]int16{
	"SP":    REG_SP,
	"CTXT":  REG_CTXT,
	"g":     REG_g,
	"RET0":  REG_RET0,
	"RET1":  REG_RET1,
	"RET2":  REG_RET2,
	"RET3":  REG_RET3,
	"PAUSE": REG_PAUSE,

	"R0":  REG_R0,
	"R1":  REG_R1,
	"R2":  REG_R2,
	"R3":  REG_R3,
	"R4":  REG_R4,
	"R5":  REG_R5,
	"R6":  REG_R6,
	"R7":  REG_R7,
	"R8":  REG_R8,
	"R9":  REG_R9,
	"R10": REG_R10,
	"R11": REG_R11,
	"R12": REG_R12,
	"R13": REG_R13,
	"R14": REG_R14,
	"R15": REG_R15,

	"F0":  REG_F0,
	"F1":  REG_F1,
	"F2":  REG_F2,
	"F3":  REG_F3,
	"F4":  REG_F4,
	"F5":  REG_F5,
	"F6":  REG_F6,
	"F7":  REG_F7,
	"F8":  REG_F8,
	"F9":  REG_F9,
	"F10": REG_F10,
	"F11": REG_F11,
	"F12": REG_F12,
	"F13": REG_F13,
	"F14": REG_F14,
	"F15": REG_F15,

	"F16": REG_F16,
	"F17": REG_F17,
	"F18": REG_F18,
	"F19": REG_F19,
	"F20": REG_F20,
	"F21": REG_F21,
	"F22": REG_F22,
	"F23": REG_F23,
	"F24": REG_F24,
	"F25": REG_F25,
	"F26": REG_F26,
	"F27": REG_F27,
	"F28": REG_F28,
	"F29": REG_F29,
	"F30": REG_F30,
	"F31": REG_F31,

	"PC_B": REG_PC_B,
}

var unaryDst = map[int]bool{
	ASet:          true,
	ATee:          true,
	ACall:         true,
	ACallIndirect: true,
	ABr:           true,
	ABrIf:         true,
	ABrTable:      true,
	AI32Store:     true,
	AI64Store:     true,
	AF32Store:     true,
	AF64Store:     true,
	AI32Store8:    true,
	AI32Store16:   true,
	AI64Store8:    true,
	AI64Store16:   true,
	AI64Store32:   true,
	ACALLNORESUME: true,
}

var Linkwasm = obj.LinkArch{
	//Arch:       sys.ArchWasm,
	//Init:       instinit,
	//Preprocess: preprocess,
	//Assemble:   assemble,
	UnaryDst: unaryDst,
}
