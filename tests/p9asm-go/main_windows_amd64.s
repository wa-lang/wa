// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

#include "textflag.h"

GLOBL ·Name(SB),NOPTR,$24

DATA ·Name+0(SB)/8,$·Name+16(SB)
DATA ·Name+8(SB)/8,$6
DATA ·Name+16(SB)/8,$"walang"

TEXT ·main(SB), $16-0
	MOVQ ·Name+0(SB), AX; MOVQ AX, 0(SP)
	MOVQ ·Name+8(SB), BX; MOVQ BX, 8(SP)
	CALL ·println(SB)
	RET
