// Inferno utils/5c/5.out.h
// https://bitbucket.org/inferno-os/inferno-os/src/master/utils/5c/5.out.h
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors. All rights reserved.
//	Portions Copyright © 2025 武汉凹语言科技有限公司.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

//go:generate go run ../stringer.go -i $GOFILE -o anames.go -p arm

package arm

import (
	"wa-lang.org/wa/internal/p9asm/objabi"
)

// 编号范围
const (
	ABase = objabi.ABaseARM + objabi.A_ARCHSPECIFIC
	AsMax = ALAST

	RBase  = objabi.RBaseARM
	RegMax = MAXREG
)

const (
	REG_R0 = RBase + iota // must be 16-aligned
	REG_R1
	REG_R2
	REG_R3
	REG_R4
	REG_R5
	REG_R6
	REG_R7
	REG_R8
	REG_R9
	REG_R10
	REG_R11
	REG_R12
	REG_R13 // SP 寄存器
	REG_R14 // LINK 寄存器
	REG_R15 // PC 寄存器

	// 浮点数寄存器
	REG_F0 // must be 16-aligned
	REG_F1
	REG_F2
	REG_F3
	REG_F4
	REG_F5
	REG_F6
	REG_F7
	REG_F8
	REG_F9
	REG_F10
	REG_F11
	REG_F12
	REG_F13
	REG_F14
	REG_F15

	// 浮点数状态?
	REG_FPSR // must be 2-aligned
	REG_FPCR

	REG_CPSR // 当前状态寄存器, must be 2-aligned
	REG_SPSR // 异常时保存当前状态寄存器

	REGRET = REG_R0 // 入参/返回值

	// compiler allocates R1 up as temps
	// compiler allocates register variables R3 up
	// compiler allocates external registers R10 down
	REGEXT = REG_R10
	// these two registers are declared in runtime.h
	REGM = REGEXT - 1 // TODO(chai2010): 删除

	REGCTXT = REG_R7  // 闭包上下文?
	REGTMP  = REG_R11 // 临时数据?
	REGSP   = REG_R13 // SP 寄存器
	REGLINK = REG_R14 // LINK 寄存器
	REGPC   = REG_R15 // PC 寄存器

	NFREG = 16 // 浮点数寄存器的个数

	// compiler allocates register variables F0 up
	// compiler allocates external registers F7 down
	FREGRET = REG_F0 // 入参/返回值
	FREGEXT = REG_F7
	FREGTMP = REG_F15 // 临时数据?
)

// 特殊的寄存器, 必须控制在 1024 范围内
const (
	REG_SPECIAL = RBase + 1<<9 + iota
	REG_MB_SY
	REG_MB_ST
	REG_MB_ISH
	REG_MB_ISHST
	REG_MB_NSH
	REG_MB_NSHST
	REG_MB_OSH
	REG_MB_OSHST

	MAXREG
)

const (
	AAND = ABase + iota
	AEOR
	ASUB
	ARSB
	AADD
	AADC
	ASBC
	ARSC
	ATST
	ATEQ
	ACMP
	ACMN
	AORR
	ABIC

	AMVN

	// Do not reorder or fragment the conditional branch
	// opcodes, or the predication code will break

	ABEQ
	ABNE
	ABCS
	ABHS
	ABCC
	ABLO
	ABMI
	ABPL
	ABVS
	ABVC
	ABHI
	ABLS
	ABGE
	ABLT
	ABGT
	ABLE

	AMOVWD
	AMOVWF
	AMOVDW
	AMOVFW
	AMOVFD
	AMOVDF
	AMOVF
	AMOVD

	ACMPF
	ACMPD
	AADDF
	AADDD
	ASUBF
	ASUBD
	AMULF
	AMULD
	ANMULF
	ANMULD
	AMULAF
	AMULAD
	ANMULAF
	ANMULAD
	AMULSF
	AMULSD
	ANMULSF
	ANMULSD
	AFMULAF
	AFMULAD
	AFNMULAF
	AFNMULAD
	AFMULSF
	AFMULSD
	AFNMULSF
	AFNMULSD
	ADIVF
	ADIVD
	ASQRTF
	ASQRTD
	AABSF
	AABSD
	ANEGF
	ANEGD

	ASRL
	ASRA
	ASLL
	AMULU
	ADIVU
	AMUL
	AMMUL
	ADIV
	AMOD
	AMODU
	ADIVHW
	ADIVUHW

	AMOVB
	AMOVBS
	AMOVBU
	AMOVH
	AMOVHS
	AMOVHU
	AMOVW
	AMOVM
	ASWPBU
	ASWPW

	ARFE
	ASWI
	AMULA
	AMULS
	AMMULA
	AMMULS

	AWORD

	AMULL
	AMULAL
	AMULLU
	AMULALU

	ABX
	ABXRET
	ADWORD

	ALDREX  // Load i32
	ASTREX  // Store i32
	ALDREXD // Load i64
	ALDREXB // Load i8
	ASTREXD // Store i64
	ASTREXB // Store i8

	ADMB

	APLD

	ACLZ
	AREV
	AREV16
	AREVSH
	ARBIT

	AXTAB
	AXTAH
	AXTABU
	AXTAHU

	ABFX
	ABFXU
	ABFC
	ABFI

	AMULWT
	AMULWB
	AMULBB
	AMULAWT
	AMULAWB
	AMULABB

	AMRC // MRC/MCR

	ALAST

	// aliases
	AB  = objabi.AJMP
	ABL = objabi.ACALL
)
