// Inferno utils/6l/span.c
// https://bitbucket.org/inferno-os/inferno-os/src/master/utils/6l/span.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors. All rights reserved.
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

package p9x86

// Instruction layout.

// Bit flags that are used to express jump target properties.
const (
	// branchBackwards marks targets that are located behind.
	// Used to express jumps to loop headers.
	branchBackwards = (1 << iota)
	// branchShort marks branches those target is close,
	// with offset is in -128..127 range.
	branchShort
	// branchLoopHead marks loop entry.
	// Used to insert padding for misaligned loops.
	branchLoopHead
)

// opBytes holds optab encoding bytes.
// Each ytab reserves fixed amount of bytes in this array.
//
// The size should be the minimal number of bytes that
// are enough to hold biggest optab op lines.
type opBytes [31]uint8

type Optab struct {
	as     As
	ytab   []ytab
	prefix uint8
	op     opBytes
}

type movtab struct {
	as   As
	ft   uint8
	f3t  uint8
	tt   uint8
	code uint8
	op   [4]uint8
}

const (
	Yxxx = iota
	Ynone
	Yi0 // $0
	Yi1 // $1
	Yu2 // $x, x fits in uint2
	Yi8 // $x, x fits in int8
	Yu8 // $x, x fits in uint8
	Yu7 // $x, x in 0..127 (fits in both int8 and uint8)
	Ys32
	Yi32
	Yi64
	Yiauto
	Yal
	Ycl
	Yax
	Ycx
	Yrb
	Yrl
	Yrl32 // Yrl on 32-bit system
	Yrf
	Yf0
	Yrx
	Ymb
	Yml
	Ym
	Ybr
	Ycs
	Yss
	Yds
	Yes
	Yfs
	Ygs
	Ygdtr
	Yidtr
	Yldtr
	Ymsw
	Ytask
	Ycr0
	Ycr1
	Ycr2
	Ycr3
	Ycr4
	Ycr5
	Ycr6
	Ycr7
	Ycr8
	Ydr0
	Ydr1
	Ydr2
	Ydr3
	Ydr4
	Ydr5
	Ydr6
	Ydr7
	Ytr0
	Ytr1
	Ytr2
	Ytr3
	Ytr4
	Ytr5
	Ytr6
	Ytr7
	Ymr
	Ymm
	Yxr0          // X0 only. "<XMM0>" notation in Intel manual.
	YxrEvexMulti4 // [ X<n> - X<n+3> ]; multisource YxrEvex
	Yxr           // X0..X15
	YxrEvex       // X0..X31
	Yxm
	YxmEvex       // YxrEvex+Ym
	Yxvm          // VSIB vector array; vm32x/vm64x
	YxvmEvex      // Yxvm which permits High-16 X register as index.
	YyrEvexMulti4 // [ Y<n> - Y<n+3> ]; multisource YyrEvex
	Yyr           // Y0..Y15
	YyrEvex       // Y0..Y31
	Yym
	YymEvex   // YyrEvex+Ym
	Yyvm      // VSIB vector array; vm32y/vm64y
	YyvmEvex  // Yyvm which permits High-16 Y register as index.
	YzrMulti4 // [ Z<n> - Z<n+3> ]; multisource YzrEvex
	Yzr       // Z0..Z31
	Yzm       // Yzr+Ym
	Yzvm      // VSIB vector array; vm32z/vm64z
	Yk0       // K0
	Yknot0    // K1..K7; write mask
	Yk        // K0..K7; used for KOP
	Ykm       // Yk+Ym; used for KOP
	Ytls
	Ytextsize
	Yindir
	Ymax
)

const (
	Zxxx = iota
	Zlit
	Zlitm_r
	Zlitr_m
	Zlit_m_r
	Z_rp
	Zbr
	Zcall
	Zcallcon
	Zcallduff
	Zcallind
	Zcallindreg
	Zib_
	Zib_rp
	Zibo_m
	Zibo_m_xm
	Zil_
	Zil_rp
	Ziq_rp
	Zilo_m
	Zjmp
	Zjmpcon
	Zloop
	Zo_iw
	Zm_o
	Zm_r
	Z_m_r
	Zm2_r
	Zm_r_xm
	Zm_r_i_xm
	Zm_r_xm_nr
	Zr_m_xm_nr
	Zibm_r // mmx1,mmx2/mem64,imm8
	Zibr_m
	Zmb_r
	Zaut_r
	Zo_m
	Zo_m64
	Zpseudo
	Zr_m
	Zr_m_xm
	Zrp_
	Z_ib
	Z_il
	Zm_ibo
	Zm_ilo
	Zib_rr
	Zil_rr
	Zbyte

	Zmax
)

const (
	Px   = 0
	Px1  = 1    // symbolic; exact value doesn't matter
	P32  = 0x32 // 32-bit only
	Pe   = 0x66 // operand escape
	Pm   = 0x0f // 2byte opcode escape
	Pq   = 0xff // both escapes: 66 0f
	Pb   = 0xfe // byte operands
	Pf2  = 0xf2 // xmm escape 1: f2 0f
	Pf3  = 0xf3 // xmm escape 2: f3 0f
	Pef3 = 0xf5 // xmm escape 2 with 16-bit prefix: 66 f3 0f
	Pq3  = 0x67 // xmm escape 3: 66 48 0f
	Pq4  = 0x68 // xmm escape 4: 66 0F 38
	Pq4w = 0x69 // Pq4 with Rex.w 66 0F 38
	Pq5  = 0x6a // xmm escape 5: F3 0F 38
	Pq5w = 0x6b // Pq5 with Rex.w F3 0F 38
	Pfw  = 0xf4 // Pf3 with Rex.w: f3 48 0f
	Pw   = 0x48 // Rex.w
	Pw8  = 0x90 // symbolic; exact value doesn't matter
	Py   = 0x80 // defaults to 64-bit mode
	Py1  = 0x81 // symbolic; exact value doesn't matter
	Py3  = 0x83 // symbolic; exact value doesn't matter
	Pavx = 0x84 // symbolic: exact value doesn't matter

	RxrEvex = 1 << 4 // AVX512 extension to REX.R/VEX.R
	Rxw     = 1 << 3 // =1, 64-bit operand size
	Rxr     = 1 << 2 // extend modrm reg
	Rxx     = 1 << 1 // extend sib index
	Rxb     = 1 << 0 // extend modrm r/m, sib base, or opcode reg
)

var ycover [Ymax * Ymax]uint8

var reg [MAXREG]int

var regrex [MAXREG + 1]int

var ynone = []ytab{
	{Zlit, 1, argList{}},
}

var ynop = []ytab{
	{Zpseudo, 0, argList{}},
	{Zpseudo, 0, argList{Yiauto}},
	{Zpseudo, 0, argList{Yml}},
	{Zpseudo, 0, argList{Yrf}},
	{Zpseudo, 0, argList{Yxr}},
	{Zpseudo, 0, argList{Yiauto}},
	{Zpseudo, 0, argList{Yml}},
	{Zpseudo, 0, argList{Yrf}},
	{Zpseudo, 1, argList{Yxr}},
}

var yxorb = []ytab{
	{Zib_, 1, argList{Yi32, Yal}},
	{Zibo_m, 2, argList{Yi32, Ymb}},
	{Zr_m, 1, argList{Yrb, Ymb}},
	{Zm_r, 1, argList{Ymb, Yrb}},
}

var yaddl = []ytab{
	{Zibo_m, 2, argList{Yi8, Yml}},
	{Zil_, 1, argList{Yi32, Yax}},
	{Zilo_m, 2, argList{Yi32, Yml}},
	{Zr_m, 1, argList{Yrl, Yml}},
	{Zm_r, 1, argList{Yml, Yrl}},
}

var yincl = []ytab{
	{Z_rp, 1, argList{Yrl}},
	{Zo_m, 2, argList{Yml}},
}

var yincq = []ytab{
	{Zo_m, 2, argList{Yml}},
}

var ycmpb = []ytab{
	{Z_ib, 1, argList{Yal, Yi32}},
	{Zm_ibo, 2, argList{Ymb, Yi32}},
	{Zm_r, 1, argList{Ymb, Yrb}},
	{Zr_m, 1, argList{Yrb, Ymb}},
}

var ycmpl = []ytab{
	{Zm_ibo, 2, argList{Yml, Yi8}},
	{Z_il, 1, argList{Yax, Yi32}},
	{Zm_ilo, 2, argList{Yml, Yi32}},
	{Zm_r, 1, argList{Yml, Yrl}},
	{Zr_m, 1, argList{Yrl, Yml}},
}

var yshb = []ytab{
	{Zo_m, 2, argList{Yi1, Ymb}},
	{Zibo_m, 2, argList{Yu8, Ymb}},
	{Zo_m, 2, argList{Ycx, Ymb}},
}

var yshl = []ytab{
	{Zo_m, 2, argList{Yi1, Yml}},
	{Zibo_m, 2, argList{Yu8, Yml}},
	{Zo_m, 2, argList{Ycl, Yml}},
	{Zo_m, 2, argList{Ycx, Yml}},
}

var ytestl = []ytab{
	{Zil_, 1, argList{Yi32, Yax}},
	{Zilo_m, 2, argList{Yi32, Yml}},
	{Zr_m, 1, argList{Yrl, Yml}},
	{Zm_r, 1, argList{Yml, Yrl}},
}

var ymovb = []ytab{
	{Zr_m, 1, argList{Yrb, Ymb}},
	{Zm_r, 1, argList{Ymb, Yrb}},
	{Zib_rp, 1, argList{Yi32, Yrb}},
	{Zibo_m, 2, argList{Yi32, Ymb}},
}

var ybtl = []ytab{
	{Zibo_m, 2, argList{Yi8, Yml}},
	{Zr_m, 1, argList{Yrl, Yml}},
}

var ymovw = []ytab{
	{Zr_m, 1, argList{Yrl, Yml}},
	{Zm_r, 1, argList{Yml, Yrl}},
	{Zil_rp, 1, argList{Yi32, Yrl}},
	{Zilo_m, 2, argList{Yi32, Yml}},
	{Zaut_r, 2, argList{Yiauto, Yrl}},
}

var ymovl = []ytab{
	{Zr_m, 1, argList{Yrl, Yml}},
	{Zm_r, 1, argList{Yml, Yrl}},
	{Zil_rp, 1, argList{Yi32, Yrl}},
	{Zilo_m, 2, argList{Yi32, Yml}},
	{Zm_r_xm, 1, argList{Yml, Ymr}}, // MMX MOVD
	{Zr_m_xm, 1, argList{Ymr, Yml}}, // MMX MOVD
	{Zm_r_xm, 2, argList{Yml, Yxr}}, // XMM MOVD (32 bit)
	{Zr_m_xm, 2, argList{Yxr, Yml}}, // XMM MOVD (32 bit)
	{Zaut_r, 2, argList{Yiauto, Yrl}},
}

var yret = []ytab{
	{Zo_iw, 1, argList{}},
	{Zo_iw, 1, argList{Yi32}},
}

var ymovq = []ytab{
	// valid in 32-bit mode
	{Zm_r_xm_nr, 1, argList{Ym, Ymr}},  // 0x6f MMX MOVQ (shorter encoding)
	{Zr_m_xm_nr, 1, argList{Ymr, Ym}},  // 0x7f MMX MOVQ
	{Zm_r_xm_nr, 2, argList{Yxr, Ymr}}, // Pf2, 0xd6 MOVDQ2Q
	{Zm_r_xm_nr, 2, argList{Yxm, Yxr}}, // Pf3, 0x7e MOVQ xmm1/m64 -> xmm2
	{Zr_m_xm_nr, 2, argList{Yxr, Yxm}}, // Pe, 0xd6 MOVQ xmm1 -> xmm2/m64

	// valid only in 64-bit mode, usually with 64-bit prefix
	{Zr_m, 1, argList{Yrl, Yml}},      // 0x89
	{Zm_r, 1, argList{Yml, Yrl}},      // 0x8b
	{Zilo_m, 2, argList{Ys32, Yrl}},   // 32 bit signed 0xc7,(0)
	{Ziq_rp, 1, argList{Yi64, Yrl}},   // 0xb8 -- 32/64 bit immediate
	{Zilo_m, 2, argList{Yi32, Yml}},   // 0xc7,(0)
	{Zm_r_xm, 1, argList{Ymm, Ymr}},   // 0x6e MMX MOVD
	{Zr_m_xm, 1, argList{Ymr, Ymm}},   // 0x7e MMX MOVD
	{Zm_r_xm, 2, argList{Yml, Yxr}},   // Pe, 0x6e MOVD xmm load
	{Zr_m_xm, 2, argList{Yxr, Yml}},   // Pe, 0x7e MOVD xmm store
	{Zaut_r, 1, argList{Yiauto, Yrl}}, // 0 built-in LEAQ
}

var ymovbe = []ytab{
	{Zlitm_r, 3, argList{Ym, Yrl}},
	{Zlitr_m, 3, argList{Yrl, Ym}},
}

var ym_rl = []ytab{
	{Zm_r, 1, argList{Ym, Yrl}},
}

var yrl_m = []ytab{
	{Zr_m, 1, argList{Yrl, Ym}},
}

var ymb_rl = []ytab{
	{Zmb_r, 1, argList{Ymb, Yrl}},
}

var yml_rl = []ytab{
	{Zm_r, 1, argList{Yml, Yrl}},
}

var yrl_ml = []ytab{
	{Zr_m, 1, argList{Yrl, Yml}},
}

var yml_mb = []ytab{
	{Zr_m, 1, argList{Yrb, Ymb}},
	{Zm_r, 1, argList{Ymb, Yrb}},
}

var yrb_mb = []ytab{
	{Zr_m, 1, argList{Yrb, Ymb}},
}

var yxchg = []ytab{
	{Z_rp, 1, argList{Yax, Yrl}},
	{Zrp_, 1, argList{Yrl, Yax}},
	{Zr_m, 1, argList{Yrl, Yml}},
	{Zm_r, 1, argList{Yml, Yrl}},
}

var ydivl = []ytab{
	{Zm_o, 2, argList{Yml}},
}

var ydivb = []ytab{
	{Zm_o, 2, argList{Ymb}},
}

var yimul = []ytab{
	{Zm_o, 2, argList{Yml}},
	{Zib_rr, 1, argList{Yi8, Yrl}},
	{Zil_rr, 1, argList{Yi32, Yrl}},
	{Zm_r, 2, argList{Yml, Yrl}},
}

var yimul3 = []ytab{
	{Zibm_r, 2, argList{Yi8, Yml, Yrl}},
	{Zibm_r, 2, argList{Yi32, Yml, Yrl}},
}

var ybyte = []ytab{
	{Zbyte, 1, argList{Yi64}},
}

var yin = []ytab{
	{Zib_, 1, argList{Yi32}},
	{Zlit, 1, argList{}},
}

var yint = []ytab{
	{Zib_, 1, argList{Yi32}},
}

var ypushl = []ytab{
	{Zrp_, 1, argList{Yrl}},
	{Zm_o, 2, argList{Ym}},
	{Zib_, 1, argList{Yi8}},
	{Zil_, 1, argList{Yi32}},
}

var ypopl = []ytab{
	{Z_rp, 1, argList{Yrl}},
	{Zo_m, 2, argList{Ym}},
}

var ywrfsbase = []ytab{
	{Zm_o, 2, argList{Yrl}},
}

var yrdrand = []ytab{
	{Zo_m, 2, argList{Yrl}},
}

var yclflush = []ytab{
	{Zo_m, 2, argList{Ym}},
}

var ybswap = []ytab{
	{Z_rp, 2, argList{Yrl}},
}

var yscond = []ytab{
	{Zo_m, 2, argList{Ymb}},
}

var yjcond = []ytab{
	{Zbr, 0, argList{Ybr}},
	{Zbr, 0, argList{Yi0, Ybr}},
	{Zbr, 1, argList{Yi1, Ybr}},
}

var yloop = []ytab{
	{Zloop, 1, argList{Ybr}},
}

var ycall = []ytab{
	{Zcallindreg, 0, argList{Yml}},
	{Zcallindreg, 2, argList{Yrx, Yrx}},
	{Zcallind, 2, argList{Yindir}},
	{Zcall, 0, argList{Ybr}},
	{Zcallcon, 1, argList{Yi32}},
}

var yjmp = []ytab{
	{Zo_m64, 2, argList{Yml}},
	{Zjmp, 0, argList{Ybr}},
	{Zjmpcon, 1, argList{Yi32}},
}

var yfmvd = []ytab{
	{Zm_o, 2, argList{Ym, Yf0}},
	{Zo_m, 2, argList{Yf0, Ym}},
	{Zm_o, 2, argList{Yrf, Yf0}},
	{Zo_m, 2, argList{Yf0, Yrf}},
}

var yfmvdp = []ytab{
	{Zo_m, 2, argList{Yf0, Ym}},
	{Zo_m, 2, argList{Yf0, Yrf}},
}

var yfmvf = []ytab{
	{Zm_o, 2, argList{Ym, Yf0}},
	{Zo_m, 2, argList{Yf0, Ym}},
}

var yfmvx = []ytab{
	{Zm_o, 2, argList{Ym, Yf0}},
}

var yfmvp = []ytab{
	{Zo_m, 2, argList{Yf0, Ym}},
}

var yfcmv = []ytab{
	{Zm_o, 2, argList{Yrf, Yf0}},
}

var yfadd = []ytab{
	{Zm_o, 2, argList{Ym, Yf0}},
	{Zm_o, 2, argList{Yrf, Yf0}},
	{Zo_m, 2, argList{Yf0, Yrf}},
}

var yfxch = []ytab{
	{Zo_m, 2, argList{Yf0, Yrf}},
	{Zm_o, 2, argList{Yrf, Yf0}},
}

var ycompp = []ytab{
	{Zo_m, 2, argList{Yf0, Yrf}}, // botch is really f0,f1
}

var ystsw = []ytab{
	{Zo_m, 2, argList{Ym}},
	{Zlit, 1, argList{Yax}},
}

var ysvrs_mo = []ytab{
	{Zm_o, 2, argList{Ym}},
}

// unaryDst version of "ysvrs_mo".
var ysvrs_om = []ytab{
	{Zo_m, 2, argList{Ym}},
}

var ymm = []ytab{
	{Zm_r_xm, 1, argList{Ymm, Ymr}},
	{Zm_r_xm, 2, argList{Yxm, Yxr}},
}

var yxm = []ytab{
	{Zm_r_xm, 1, argList{Yxm, Yxr}},
}

var yxm_q4 = []ytab{
	{Zm_r, 1, argList{Yxm, Yxr}},
}

var yxcvm1 = []ytab{
	{Zm_r_xm, 2, argList{Yxm, Yxr}},
	{Zm_r_xm, 2, argList{Yxm, Ymr}},
}

var yxcvm2 = []ytab{
	{Zm_r_xm, 2, argList{Yxm, Yxr}},
	{Zm_r_xm, 2, argList{Ymm, Yxr}},
}

var yxr = []ytab{
	{Zm_r_xm, 1, argList{Yxr, Yxr}},
}

var yxr_ml = []ytab{
	{Zr_m_xm, 1, argList{Yxr, Yml}},
}

var ymr = []ytab{
	{Zm_r, 1, argList{Ymr, Ymr}},
}

var ymr_ml = []ytab{
	{Zr_m_xm, 1, argList{Ymr, Yml}},
}

var yxcmpi = []ytab{
	{Zm_r_i_xm, 2, argList{Yxm, Yxr, Yi8}},
}

var yxmov = []ytab{
	{Zm_r_xm, 1, argList{Yxm, Yxr}},
	{Zr_m_xm, 1, argList{Yxr, Yxm}},
}

var yxcvfl = []ytab{
	{Zm_r_xm, 1, argList{Yxm, Yrl}},
}

var yxcvlf = []ytab{
	{Zm_r_xm, 1, argList{Yml, Yxr}},
}

var yxcvfq = []ytab{
	{Zm_r_xm, 2, argList{Yxm, Yrl}},
}

var yxcvqf = []ytab{
	{Zm_r_xm, 2, argList{Yml, Yxr}},
}

var yps = []ytab{
	{Zm_r_xm, 1, argList{Ymm, Ymr}},
	{Zibo_m_xm, 2, argList{Yi8, Ymr}},
	{Zm_r_xm, 2, argList{Yxm, Yxr}},
	{Zibo_m_xm, 3, argList{Yi8, Yxr}},
}

var yxrrl = []ytab{
	{Zm_r, 1, argList{Yxr, Yrl}},
}

var ymrxr = []ytab{
	{Zm_r, 1, argList{Ymr, Yxr}},
	{Zm_r_xm, 1, argList{Yxm, Yxr}},
}

var ymshuf = []ytab{
	{Zibm_r, 2, argList{Yi8, Ymm, Ymr}},
}

var ymshufb = []ytab{
	{Zm2_r, 2, argList{Yxm, Yxr}},
}

// It should never have more than 1 entry,
// because some optab entries you opcode secuences that
// are longer than 2 bytes (zoffset=2 here),
// ROUNDPD and ROUNDPS and recently added BLENDPD,
// to name a few.
var yxshuf = []ytab{
	{Zibm_r, 2, argList{Yu8, Yxm, Yxr}},
}

var yextrw = []ytab{
	{Zibm_r, 2, argList{Yu8, Yxr, Yrl}},
	{Zibr_m, 2, argList{Yu8, Yxr, Yml}},
}

var yextr = []ytab{
	{Zibr_m, 3, argList{Yu8, Yxr, Ymm}},
}

var yinsrw = []ytab{
	{Zibm_r, 2, argList{Yu8, Yml, Yxr}},
}

var yinsr = []ytab{
	{Zibm_r, 3, argList{Yu8, Ymm, Yxr}},
}

var ypsdq = []ytab{
	{Zibo_m, 2, argList{Yi8, Yxr}},
}

var ymskb = []ytab{
	{Zm_r_xm, 2, argList{Yxr, Yrl}},
	{Zm_r_xm, 1, argList{Ymr, Yrl}},
}

var ycrc32l = []ytab{
	{Zlitm_r, 0, argList{Yml, Yrl}},
}

var ycrc32b = []ytab{
	{Zlitm_r, 0, argList{Ymb, Yrl}},
}

var yprefetch = []ytab{
	{Zm_o, 2, argList{Ym}},
}

var yaes = []ytab{
	{Zlitm_r, 2, argList{Yxm, Yxr}},
}

var yxbegin = []ytab{
	{Zjmp, 1, argList{Ybr}},
}

var yxabort = []ytab{
	{Zib_, 1, argList{Yu8}},
}

var ylddqu = []ytab{
	{Zm_r, 1, argList{Ym, Yxr}},
}

var ypalignr = []ytab{
	{Zibm_r, 2, argList{Yu8, Yxm, Yxr}},
}

var ysha256rnds2 = []ytab{
	{Zlit_m_r, 0, argList{Yxr0, Yxm, Yxr}},
}

var yblendvpd = []ytab{
	{Z_m_r, 1, argList{Yxr0, Yxm, Yxr}},
}

var ymmxmm0f38 = []ytab{
	{Zlitm_r, 3, argList{Ymm, Ymr}},
	{Zlitm_r, 5, argList{Yxm, Yxr}},
}

var yextractps = []ytab{
	{Zibr_m, 2, argList{Yu2, Yxr, Yml}},
}

var ysha1rnds4 = []ytab{
	{Zibm_r, 2, argList{Yu2, Yxm, Yxr}},
}

// You are doasm, holding in your hand a *Prog with p.As set to, say,
// ACRC32, and p.From and p.To as operands (Addr).  The linker scans optab
// to find the entry with the given p.As and then looks through the ytable for
// that instruction (the second field in the optab struct) for a line whose
// first two values match the Ytypes of the p.From and p.To operands.  The
// function oclass computes the specific Ytype of an operand and then the set
// of more general Ytypes that it satisfies is implied by the ycover table, set
// up in instinit.  For example, oclass distinguishes the constants 0 and 1
// from the more general 8-bit constants, but instinit says
//
//	ycover[Yi0*Ymax+Ys32] = 1
//	ycover[Yi1*Ymax+Ys32] = 1
//	ycover[Yi8*Ymax+Ys32] = 1
//
// which means that Yi0, Yi1, and Yi8 all count as Ys32 (signed 32)
// if that's what an instruction can handle.
//
// In parallel with the scan through the ytable for the appropriate line, there
// is a z pointer that starts out pointing at the strange magic byte list in
// the Optab struct.  With each step past a non-matching ytable line, z
// advances by the 4th entry in the line.  When a matching line is found, that
// z pointer has the extra data to use in laying down the instruction bytes.
// The actual bytes laid down are a function of the 3rd entry in the line (that
// is, the Ztype) and the z bytes.
//
// For example, let's look at AADDL.  The optab line says:
//
//	{AADDL, yaddl, Px, opBytes{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
//
// and yaddl says
//
//	var yaddl = []ytab{
//	        {Yi8, Ynone, Yml, Zibo_m, 2},
//	        {Yi32, Ynone, Yax, Zil_, 1},
//	        {Yi32, Ynone, Yml, Zilo_m, 2},
//	        {Yrl, Ynone, Yml, Zr_m, 1},
//	        {Yml, Ynone, Yrl, Zm_r, 1},
//	}
//
// so there are 5 possible types of ADDL instruction that can be laid down, and
// possible states used to lay them down (Ztype and z pointer, assuming z
// points at opBytes{0x83, 00, 0x05,0x81, 00, 0x01, 0x03}) are:
//
//	Yi8, Yml -> Zibo_m, z (0x83, 00)
//	Yi32, Yax -> Zil_, z+2 (0x05)
//	Yi32, Yml -> Zilo_m, z+2+1 (0x81, 0x00)
//	Yrl, Yml -> Zr_m, z+2+1+2 (0x01)
//	Yml, Yrl -> Zm_r, z+2+1+2+1 (0x03)
//
// The Pconstant in the optab line controls the prefix bytes to emit.  That's
// relatively straightforward as this program goes.
//
// The switch on yt.zcase in doasm implements the various Z cases.  Zibo_m, for
// example, is an opcode byte (z[0]) then an asmando (which is some kind of
// encoded addressing mode for the Yml arg), and then a single immediate byte.
// Zilo_m is the same but a long (32-bit) immediate.
var optab =
// as, ytab, andproto, opcode
[...]Optab{
	{AXXX, nil, 0, opBytes{}},
	{AAAA, ynone, P32, opBytes{0x37}},
	{AAAD, ynone, P32, opBytes{0xd5, 0x0a}},
	{AAAM, ynone, P32, opBytes{0xd4, 0x0a}},
	{AAAS, ynone, P32, opBytes{0x3f}},
	{AADCB, yxorb, Pb, opBytes{0x14, 0x80, 02, 0x10, 0x12}},
	{AADCL, yaddl, Px, opBytes{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADCQ, yaddl, Pw, opBytes{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADCW, yaddl, Pe, opBytes{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADCXL, yml_rl, Pq4, opBytes{0xf6}},
	{AADCXQ, yml_rl, Pq4w, opBytes{0xf6}},
	{AADDB, yxorb, Pb, opBytes{0x04, 0x80, 00, 0x00, 0x02}},
	{AADDL, yaddl, Px, opBytes{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADDPD, yxm, Pq, opBytes{0x58}},
	{AADDPS, yxm, Pm, opBytes{0x58}},
	{AADDQ, yaddl, Pw, opBytes{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADDSD, yxm, Pf2, opBytes{0x58}},
	{AADDSS, yxm, Pf3, opBytes{0x58}},
	{AADDSUBPD, yxm, Pq, opBytes{0xd0}},
	{AADDSUBPS, yxm, Pf2, opBytes{0xd0}},
	{AADDW, yaddl, Pe, opBytes{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADOXL, yml_rl, Pq5, opBytes{0xf6}},
	{AADOXQ, yml_rl, Pq5w, opBytes{0xf6}},
	{AADJSP, nil, 0, opBytes{}},
	{AANDB, yxorb, Pb, opBytes{0x24, 0x80, 04, 0x20, 0x22}},
	{AANDL, yaddl, Px, opBytes{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AANDNPD, yxm, Pq, opBytes{0x55}},
	{AANDNPS, yxm, Pm, opBytes{0x55}},
	{AANDPD, yxm, Pq, opBytes{0x54}},
	{AANDPS, yxm, Pm, opBytes{0x54}},
	{AANDQ, yaddl, Pw, opBytes{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AANDW, yaddl, Pe, opBytes{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AARPL, yrl_ml, P32, opBytes{0x63}},
	{ABOUNDL, yrl_m, P32, opBytes{0x62}},
	{ABOUNDW, yrl_m, Pe, opBytes{0x62}},
	{ABSFL, yml_rl, Pm, opBytes{0xbc}},
	{ABSFQ, yml_rl, Pw, opBytes{0x0f, 0xbc}},
	{ABSFW, yml_rl, Pq, opBytes{0xbc}},
	{ABSRL, yml_rl, Pm, opBytes{0xbd}},
	{ABSRQ, yml_rl, Pw, opBytes{0x0f, 0xbd}},
	{ABSRW, yml_rl, Pq, opBytes{0xbd}},
	{ABSWAPL, ybswap, Px, opBytes{0x0f, 0xc8}},
	{ABSWAPQ, ybswap, Pw, opBytes{0x0f, 0xc8}},
	{ABTCL, ybtl, Pm, opBytes{0xba, 07, 0xbb}},
	{ABTCQ, ybtl, Pw, opBytes{0x0f, 0xba, 07, 0x0f, 0xbb}},
	{ABTCW, ybtl, Pq, opBytes{0xba, 07, 0xbb}},
	{ABTL, ybtl, Pm, opBytes{0xba, 04, 0xa3}},
	{ABTQ, ybtl, Pw, opBytes{0x0f, 0xba, 04, 0x0f, 0xa3}},
	{ABTRL, ybtl, Pm, opBytes{0xba, 06, 0xb3}},
	{ABTRQ, ybtl, Pw, opBytes{0x0f, 0xba, 06, 0x0f, 0xb3}},
	{ABTRW, ybtl, Pq, opBytes{0xba, 06, 0xb3}},
	{ABTSL, ybtl, Pm, opBytes{0xba, 05, 0xab}},
	{ABTSQ, ybtl, Pw, opBytes{0x0f, 0xba, 05, 0x0f, 0xab}},
	{ABTSW, ybtl, Pq, opBytes{0xba, 05, 0xab}},
	{ABTW, ybtl, Pq, opBytes{0xba, 04, 0xa3}},
	{ABYTE, ybyte, Px, opBytes{1}},
	{ACALL, ycall, Px, opBytes{0xff, 02, 0xff, 0x15, 0xe8}},
	{ACBW, ynone, Pe, opBytes{0x98}},
	{ACDQ, ynone, Px, opBytes{0x99}},
	{ACDQE, ynone, Pw, opBytes{0x98}},
	{ACLAC, ynone, Pm, opBytes{01, 0xca}},
	{ACLC, ynone, Px, opBytes{0xf8}},
	{ACLD, ynone, Px, opBytes{0xfc}},
	{ACLDEMOTE, yclflush, Pm, opBytes{0x1c, 00}},
	{ACLFLUSH, yclflush, Pm, opBytes{0xae, 07}},
	{ACLFLUSHOPT, yclflush, Pq, opBytes{0xae, 07}},
	{ACLI, ynone, Px, opBytes{0xfa}},
	{ACLTS, ynone, Pm, opBytes{0x06}},
	{ACLWB, yclflush, Pq, opBytes{0xae, 06}},
	{ACMC, ynone, Px, opBytes{0xf5}},
	{ACMOVLCC, yml_rl, Pm, opBytes{0x43}},
	{ACMOVLCS, yml_rl, Pm, opBytes{0x42}},
	{ACMOVLEQ, yml_rl, Pm, opBytes{0x44}},
	{ACMOVLGE, yml_rl, Pm, opBytes{0x4d}},
	{ACMOVLGT, yml_rl, Pm, opBytes{0x4f}},
	{ACMOVLHI, yml_rl, Pm, opBytes{0x47}},
	{ACMOVLLE, yml_rl, Pm, opBytes{0x4e}},
	{ACMOVLLS, yml_rl, Pm, opBytes{0x46}},
	{ACMOVLLT, yml_rl, Pm, opBytes{0x4c}},
	{ACMOVLMI, yml_rl, Pm, opBytes{0x48}},
	{ACMOVLNE, yml_rl, Pm, opBytes{0x45}},
	{ACMOVLOC, yml_rl, Pm, opBytes{0x41}},
	{ACMOVLOS, yml_rl, Pm, opBytes{0x40}},
	{ACMOVLPC, yml_rl, Pm, opBytes{0x4b}},
	{ACMOVLPL, yml_rl, Pm, opBytes{0x49}},
	{ACMOVLPS, yml_rl, Pm, opBytes{0x4a}},
	{ACMOVQCC, yml_rl, Pw, opBytes{0x0f, 0x43}},
	{ACMOVQCS, yml_rl, Pw, opBytes{0x0f, 0x42}},
	{ACMOVQEQ, yml_rl, Pw, opBytes{0x0f, 0x44}},
	{ACMOVQGE, yml_rl, Pw, opBytes{0x0f, 0x4d}},
	{ACMOVQGT, yml_rl, Pw, opBytes{0x0f, 0x4f}},
	{ACMOVQHI, yml_rl, Pw, opBytes{0x0f, 0x47}},
	{ACMOVQLE, yml_rl, Pw, opBytes{0x0f, 0x4e}},
	{ACMOVQLS, yml_rl, Pw, opBytes{0x0f, 0x46}},
	{ACMOVQLT, yml_rl, Pw, opBytes{0x0f, 0x4c}},
	{ACMOVQMI, yml_rl, Pw, opBytes{0x0f, 0x48}},
	{ACMOVQNE, yml_rl, Pw, opBytes{0x0f, 0x45}},
	{ACMOVQOC, yml_rl, Pw, opBytes{0x0f, 0x41}},
	{ACMOVQOS, yml_rl, Pw, opBytes{0x0f, 0x40}},
	{ACMOVQPC, yml_rl, Pw, opBytes{0x0f, 0x4b}},
	{ACMOVQPL, yml_rl, Pw, opBytes{0x0f, 0x49}},
	{ACMOVQPS, yml_rl, Pw, opBytes{0x0f, 0x4a}},
	{ACMOVWCC, yml_rl, Pq, opBytes{0x43}},
	{ACMOVWCS, yml_rl, Pq, opBytes{0x42}},
	{ACMOVWEQ, yml_rl, Pq, opBytes{0x44}},
	{ACMOVWGE, yml_rl, Pq, opBytes{0x4d}},
	{ACMOVWGT, yml_rl, Pq, opBytes{0x4f}},
	{ACMOVWHI, yml_rl, Pq, opBytes{0x47}},
	{ACMOVWLE, yml_rl, Pq, opBytes{0x4e}},
	{ACMOVWLS, yml_rl, Pq, opBytes{0x46}},
	{ACMOVWLT, yml_rl, Pq, opBytes{0x4c}},
	{ACMOVWMI, yml_rl, Pq, opBytes{0x48}},
	{ACMOVWNE, yml_rl, Pq, opBytes{0x45}},
	{ACMOVWOC, yml_rl, Pq, opBytes{0x41}},
	{ACMOVWOS, yml_rl, Pq, opBytes{0x40}},
	{ACMOVWPC, yml_rl, Pq, opBytes{0x4b}},
	{ACMOVWPL, yml_rl, Pq, opBytes{0x49}},
	{ACMOVWPS, yml_rl, Pq, opBytes{0x4a}},
	{ACMPB, ycmpb, Pb, opBytes{0x3c, 0x80, 07, 0x38, 0x3a}},
	{ACMPL, ycmpl, Px, opBytes{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACMPPD, yxcmpi, Px, opBytes{Pe, 0xc2}},
	{ACMPPS, yxcmpi, Pm, opBytes{0xc2, 0}},
	{ACMPQ, ycmpl, Pw, opBytes{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACMPSB, ynone, Pb, opBytes{0xa6}},
	{ACMPSD, yxcmpi, Px, opBytes{Pf2, 0xc2}},
	{ACMPSL, ynone, Px, opBytes{0xa7}},
	{ACMPSQ, ynone, Pw, opBytes{0xa7}},
	{ACMPSS, yxcmpi, Px, opBytes{Pf3, 0xc2}},
	{ACMPSW, ynone, Pe, opBytes{0xa7}},
	{ACMPW, ycmpl, Pe, opBytes{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACOMISD, yxm, Pe, opBytes{0x2f}},
	{ACOMISS, yxm, Pm, opBytes{0x2f}},
	{ACPUID, ynone, Pm, opBytes{0xa2}},
	{ACVTPL2PD, yxcvm2, Px, opBytes{Pf3, 0xe6, Pe, 0x2a}},
	{ACVTPL2PS, yxcvm2, Pm, opBytes{0x5b, 0, 0x2a, 0}},
	{ACVTPD2PL, yxcvm1, Px, opBytes{Pf2, 0xe6, Pe, 0x2d}},
	{ACVTPD2PS, yxm, Pe, opBytes{0x5a}},
	{ACVTPS2PL, yxcvm1, Px, opBytes{Pe, 0x5b, Pm, 0x2d}},
	{ACVTPS2PD, yxm, Pm, opBytes{0x5a}},
	{ACVTSD2SL, yxcvfl, Pf2, opBytes{0x2d}},
	{ACVTSD2SQ, yxcvfq, Pw, opBytes{Pf2, 0x2d}},
	{ACVTSD2SS, yxm, Pf2, opBytes{0x5a}},
	{ACVTSL2SD, yxcvlf, Pf2, opBytes{0x2a}},
	{ACVTSQ2SD, yxcvqf, Pw, opBytes{Pf2, 0x2a}},
	{ACVTSL2SS, yxcvlf, Pf3, opBytes{0x2a}},
	{ACVTSQ2SS, yxcvqf, Pw, opBytes{Pf3, 0x2a}},
	{ACVTSS2SD, yxm, Pf3, opBytes{0x5a}},
	{ACVTSS2SL, yxcvfl, Pf3, opBytes{0x2d}},
	{ACVTSS2SQ, yxcvfq, Pw, opBytes{Pf3, 0x2d}},
	{ACVTTPD2PL, yxcvm1, Px, opBytes{Pe, 0xe6, Pe, 0x2c}},
	{ACVTTPS2PL, yxcvm1, Px, opBytes{Pf3, 0x5b, Pm, 0x2c}},
	{ACVTTSD2SL, yxcvfl, Pf2, opBytes{0x2c}},
	{ACVTTSD2SQ, yxcvfq, Pw, opBytes{Pf2, 0x2c}},
	{ACVTTSS2SL, yxcvfl, Pf3, opBytes{0x2c}},
	{ACVTTSS2SQ, yxcvfq, Pw, opBytes{Pf3, 0x2c}},
	{ACWD, ynone, Pe, opBytes{0x99}},
	{ACWDE, ynone, Px, opBytes{0x98}},
	{ACQO, ynone, Pw, opBytes{0x99}},
	{ADAA, ynone, P32, opBytes{0x27}},
	{ADAS, ynone, P32, opBytes{0x2f}},
	{ADECB, yscond, Pb, opBytes{0xfe, 01}},
	{ADECL, yincl, Px1, opBytes{0x48, 0xff, 01}},
	{ADECQ, yincq, Pw, opBytes{0xff, 01}},
	{ADECW, yincq, Pe, opBytes{0xff, 01}},
	{ADIVB, ydivb, Pb, opBytes{0xf6, 06}},
	{ADIVL, ydivl, Px, opBytes{0xf7, 06}},
	{ADIVPD, yxm, Pe, opBytes{0x5e}},
	{ADIVPS, yxm, Pm, opBytes{0x5e}},
	{ADIVQ, ydivl, Pw, opBytes{0xf7, 06}},
	{ADIVSD, yxm, Pf2, opBytes{0x5e}},
	{ADIVSS, yxm, Pf3, opBytes{0x5e}},
	{ADIVW, ydivl, Pe, opBytes{0xf7, 06}},
	{ADPPD, yxshuf, Pq, opBytes{0x3a, 0x41, 0}},
	{ADPPS, yxshuf, Pq, opBytes{0x3a, 0x40, 0}},
	{AEMMS, ynone, Pm, opBytes{0x77}},
	{AEXTRACTPS, yextractps, Pq, opBytes{0x3a, 0x17, 0}},
	{AENTER, nil, 0, opBytes{}}, // botch
	{AFXRSTOR, ysvrs_mo, Pm, opBytes{0xae, 01, 0xae, 01}},
	{AFXSAVE, ysvrs_om, Pm, opBytes{0xae, 00, 0xae, 00}},
	{AFXRSTOR64, ysvrs_mo, Pw, opBytes{0x0f, 0xae, 01, 0x0f, 0xae, 01}},
	{AFXSAVE64, ysvrs_om, Pw, opBytes{0x0f, 0xae, 00, 0x0f, 0xae, 00}},
	{AHLT, ynone, Px, opBytes{0xf4}},
	{AIDIVB, ydivb, Pb, opBytes{0xf6, 07}},
	{AIDIVL, ydivl, Px, opBytes{0xf7, 07}},
	{AIDIVQ, ydivl, Pw, opBytes{0xf7, 07}},
	{AIDIVW, ydivl, Pe, opBytes{0xf7, 07}},
	{AIMULB, ydivb, Pb, opBytes{0xf6, 05}},
	{AIMULL, yimul, Px, opBytes{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMULQ, yimul, Pw, opBytes{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMULW, yimul, Pe, opBytes{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMUL3W, yimul3, Pe, opBytes{0x6b, 00, 0x69, 00}},
	{AIMUL3L, yimul3, Px, opBytes{0x6b, 00, 0x69, 00}},
	{AIMUL3Q, yimul3, Pw, opBytes{0x6b, 00, 0x69, 00}},
	{AINB, yin, Pb, opBytes{0xe4, 0xec}},
	{AINW, yin, Pe, opBytes{0xe5, 0xed}},
	{AINL, yin, Px, opBytes{0xe5, 0xed}},
	{AINCB, yscond, Pb, opBytes{0xfe, 00}},
	{AINCL, yincl, Px1, opBytes{0x40, 0xff, 00}},
	{AINCQ, yincq, Pw, opBytes{0xff, 00}},
	{AINCW, yincq, Pe, opBytes{0xff, 00}},
	{AINSB, ynone, Pb, opBytes{0x6c}},
	{AINSL, ynone, Px, opBytes{0x6d}},
	{AINSERTPS, yxshuf, Pq, opBytes{0x3a, 0x21, 0}},
	{AINSW, ynone, Pe, opBytes{0x6d}},
	{AICEBP, ynone, Px, opBytes{0xf1}},
	{AINT, yint, Px, opBytes{0xcd}},
	{AINTO, ynone, P32, opBytes{0xce}},
	{AIRETL, ynone, Px, opBytes{0xcf}},
	{AIRETQ, ynone, Pw, opBytes{0xcf}},
	{AIRETW, ynone, Pe, opBytes{0xcf}},
	{AJCC, yjcond, Px, opBytes{0x73, 0x83, 00}},
	{AJCS, yjcond, Px, opBytes{0x72, 0x82}},
	{AJCXZL, yloop, Px, opBytes{0xe3}},
	{AJCXZW, yloop, Px, opBytes{0xe3}},
	{AJCXZQ, yloop, Px, opBytes{0xe3}},
	{AJEQ, yjcond, Px, opBytes{0x74, 0x84}},
	{AJGE, yjcond, Px, opBytes{0x7d, 0x8d}},
	{AJGT, yjcond, Px, opBytes{0x7f, 0x8f}},
	{AJHI, yjcond, Px, opBytes{0x77, 0x87}},
	{AJLE, yjcond, Px, opBytes{0x7e, 0x8e}},
	{AJLS, yjcond, Px, opBytes{0x76, 0x86}},
	{AJLT, yjcond, Px, opBytes{0x7c, 0x8c}},
	{AJMI, yjcond, Px, opBytes{0x78, 0x88}},
	{AJMP, yjmp, Px, opBytes{0xff, 04, 0xeb, 0xe9}},
	{AJNE, yjcond, Px, opBytes{0x75, 0x85}},
	{AJOC, yjcond, Px, opBytes{0x71, 0x81, 00}},
	{AJOS, yjcond, Px, opBytes{0x70, 0x80, 00}},
	{AJPC, yjcond, Px, opBytes{0x7b, 0x8b}},
	{AJPL, yjcond, Px, opBytes{0x79, 0x89}},
	{AJPS, yjcond, Px, opBytes{0x7a, 0x8a}},
	{AHADDPD, yxm, Pq, opBytes{0x7c}},
	{AHADDPS, yxm, Pf2, opBytes{0x7c}},
	{AHSUBPD, yxm, Pq, opBytes{0x7d}},
	{AHSUBPS, yxm, Pf2, opBytes{0x7d}},
	{ALAHF, ynone, Px, opBytes{0x9f}},
	{ALARL, yml_rl, Pm, opBytes{0x02}},
	{ALARQ, yml_rl, Pw, opBytes{0x0f, 0x02}},
	{ALARW, yml_rl, Pq, opBytes{0x02}},
	{ALDDQU, ylddqu, Pf2, opBytes{0xf0}},
	{ALDMXCSR, ysvrs_mo, Pm, opBytes{0xae, 02, 0xae, 02}},
	{ALEAL, ym_rl, Px, opBytes{0x8d}},
	{ALEAQ, ym_rl, Pw, opBytes{0x8d}},
	{ALEAVEL, ynone, P32, opBytes{0xc9}},
	{ALEAVEQ, ynone, Py, opBytes{0xc9}},
	{ALEAVEW, ynone, Pe, opBytes{0xc9}},
	{ALEAW, ym_rl, Pe, opBytes{0x8d}},
	{ALOCK, ynone, Px, opBytes{0xf0}},
	{ALODSB, ynone, Pb, opBytes{0xac}},
	{ALODSL, ynone, Px, opBytes{0xad}},
	{ALODSQ, ynone, Pw, opBytes{0xad}},
	{ALODSW, ynone, Pe, opBytes{0xad}},
	{ALONG, ybyte, Px, opBytes{4}},
	{ALOOP, yloop, Px, opBytes{0xe2}},
	{ALOOPEQ, yloop, Px, opBytes{0xe1}},
	{ALOOPNE, yloop, Px, opBytes{0xe0}},
	{ALTR, ydivl, Pm, opBytes{0x00, 03}},
	{ALZCNTL, yml_rl, Pf3, opBytes{0xbd}},
	{ALZCNTQ, yml_rl, Pfw, opBytes{0xbd}},
	{ALZCNTW, yml_rl, Pef3, opBytes{0xbd}},
	{ALSLL, yml_rl, Pm, opBytes{0x03}},
	{ALSLW, yml_rl, Pq, opBytes{0x03}},
	{ALSLQ, yml_rl, Pw, opBytes{0x0f, 0x03}},
	{AMASKMOVOU, yxr, Pe, opBytes{0xf7}},
	{AMASKMOVQ, ymr, Pm, opBytes{0xf7}},
	{AMAXPD, yxm, Pe, opBytes{0x5f}},
	{AMAXPS, yxm, Pm, opBytes{0x5f}},
	{AMAXSD, yxm, Pf2, opBytes{0x5f}},
	{AMAXSS, yxm, Pf3, opBytes{0x5f}},
	{AMINPD, yxm, Pe, opBytes{0x5d}},
	{AMINPS, yxm, Pm, opBytes{0x5d}},
	{AMINSD, yxm, Pf2, opBytes{0x5d}},
	{AMINSS, yxm, Pf3, opBytes{0x5d}},
	{AMONITOR, ynone, Px, opBytes{0x0f, 0x01, 0xc8, 0}},
	{AMWAIT, ynone, Px, opBytes{0x0f, 0x01, 0xc9, 0}},
	{AMOVAPD, yxmov, Pe, opBytes{0x28, 0x29}},
	{AMOVAPS, yxmov, Pm, opBytes{0x28, 0x29}},
	{AMOVB, ymovb, Pb, opBytes{0x88, 0x8a, 0xb0, 0xc6, 00}},
	{AMOVBLSX, ymb_rl, Pm, opBytes{0xbe}},
	{AMOVBLZX, ymb_rl, Pm, opBytes{0xb6}},
	{AMOVBQSX, ymb_rl, Pw, opBytes{0x0f, 0xbe}},
	{AMOVBQZX, ymb_rl, Pw, opBytes{0x0f, 0xb6}},
	{AMOVBWSX, ymb_rl, Pq, opBytes{0xbe}},
	{AMOVSWW, ymb_rl, Pe, opBytes{0x0f, 0xbf}},
	{AMOVBWZX, ymb_rl, Pq, opBytes{0xb6}},
	{AMOVZWW, ymb_rl, Pe, opBytes{0x0f, 0xb7}},
	{AMOVO, yxmov, Pe, opBytes{0x6f, 0x7f}},
	{AMOVOU, yxmov, Pf3, opBytes{0x6f, 0x7f}},
	{AMOVHLPS, yxr, Pm, opBytes{0x12}},
	{AMOVHPD, yxmov, Pe, opBytes{0x16, 0x17}},
	{AMOVHPS, yxmov, Pm, opBytes{0x16, 0x17}},
	{AMOVL, ymovl, Px, opBytes{0x89, 0x8b, 0xb8, 0xc7, 00, 0x6e, 0x7e, Pe, 0x6e, Pe, 0x7e, 0}},
	{AMOVLHPS, yxr, Pm, opBytes{0x16}},
	{AMOVLPD, yxmov, Pe, opBytes{0x12, 0x13}},
	{AMOVLPS, yxmov, Pm, opBytes{0x12, 0x13}},
	{AMOVLQSX, yml_rl, Pw, opBytes{0x63}},
	{AMOVLQZX, yml_rl, Px, opBytes{0x8b}},
	{AMOVMSKPD, yxrrl, Pq, opBytes{0x50}},
	{AMOVMSKPS, yxrrl, Pm, opBytes{0x50}},
	{AMOVNTO, yxr_ml, Pe, opBytes{0xe7}},
	{AMOVNTDQA, ylddqu, Pq4, opBytes{0x2a}},
	{AMOVNTPD, yxr_ml, Pe, opBytes{0x2b}},
	{AMOVNTPS, yxr_ml, Pm, opBytes{0x2b}},
	{AMOVNTQ, ymr_ml, Pm, opBytes{0xe7}},
	{AMOVQ, ymovq, Pw8, opBytes{0x6f, 0x7f, Pf2, 0xd6, Pf3, 0x7e, Pe, 0xd6, 0x89, 0x8b, 0xc7, 00, 0xb8, 0xc7, 00, 0x6e, 0x7e, Pe, 0x6e, Pe, 0x7e, 0}},
	{AMOVQOZX, ymrxr, Pf3, opBytes{0xd6, 0x7e}},
	{AMOVSB, ynone, Pb, opBytes{0xa4}},
	{AMOVSD, yxmov, Pf2, opBytes{0x10, 0x11}},
	{AMOVSL, ynone, Px, opBytes{0xa5}},
	{AMOVSQ, ynone, Pw, opBytes{0xa5}},
	{AMOVSS, yxmov, Pf3, opBytes{0x10, 0x11}},
	{AMOVSW, ynone, Pe, opBytes{0xa5}},
	{AMOVUPD, yxmov, Pe, opBytes{0x10, 0x11}},
	{AMOVUPS, yxmov, Pm, opBytes{0x10, 0x11}},
	{AMOVW, ymovw, Pe, opBytes{0x89, 0x8b, 0xb8, 0xc7, 00, 0}},
	{AMOVWLSX, yml_rl, Pm, opBytes{0xbf}},
	{AMOVWLZX, yml_rl, Pm, opBytes{0xb7}},
	{AMOVWQSX, yml_rl, Pw, opBytes{0x0f, 0xbf}},
	{AMOVWQZX, yml_rl, Pw, opBytes{0x0f, 0xb7}},
	{AMPSADBW, yxshuf, Pq, opBytes{0x3a, 0x42, 0}},
	{AMULB, ydivb, Pb, opBytes{0xf6, 04}},
	{AMULL, ydivl, Px, opBytes{0xf7, 04}},
	{AMULPD, yxm, Pe, opBytes{0x59}},
	{AMULPS, yxm, Ym, opBytes{0x59}},
	{AMULQ, ydivl, Pw, opBytes{0xf7, 04}},
	{AMULSD, yxm, Pf2, opBytes{0x59}},
	{AMULSS, yxm, Pf3, opBytes{0x59}},
	{AMULW, ydivl, Pe, opBytes{0xf7, 04}},
	{ANEGB, yscond, Pb, opBytes{0xf6, 03}},
	{ANEGL, yscond, Px, opBytes{0xf7, 03}},
	{ANEGQ, yscond, Pw, opBytes{0xf7, 03}},
	{ANEGW, yscond, Pe, opBytes{0xf7, 03}},
	{ANOP, ynop, Px, opBytes{0, 0}},
	{ANOTB, yscond, Pb, opBytes{0xf6, 02}},
	{ANOTL, yscond, Px, opBytes{0xf7, 02}}, // TODO(rsc): yscond is wrong here.
	{ANOTQ, yscond, Pw, opBytes{0xf7, 02}},
	{ANOTW, yscond, Pe, opBytes{0xf7, 02}},
	{AORB, yxorb, Pb, opBytes{0x0c, 0x80, 01, 0x08, 0x0a}},
	{AORL, yaddl, Px, opBytes{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AORPD, yxm, Pq, opBytes{0x56}},
	{AORPS, yxm, Pm, opBytes{0x56}},
	{AORQ, yaddl, Pw, opBytes{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AORW, yaddl, Pe, opBytes{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AOUTB, yin, Pb, opBytes{0xe6, 0xee}},
	{AOUTL, yin, Px, opBytes{0xe7, 0xef}},
	{AOUTW, yin, Pe, opBytes{0xe7, 0xef}},
	{AOUTSB, ynone, Pb, opBytes{0x6e}},
	{AOUTSL, ynone, Px, opBytes{0x6f}},
	{AOUTSW, ynone, Pe, opBytes{0x6f}},
	{APABSB, yxm_q4, Pq4, opBytes{0x1c}},
	{APABSD, yxm_q4, Pq4, opBytes{0x1e}},
	{APABSW, yxm_q4, Pq4, opBytes{0x1d}},
	{APACKSSLW, ymm, Py1, opBytes{0x6b, Pe, 0x6b}},
	{APACKSSWB, ymm, Py1, opBytes{0x63, Pe, 0x63}},
	{APACKUSDW, yxm_q4, Pq4, opBytes{0x2b}},
	{APACKUSWB, ymm, Py1, opBytes{0x67, Pe, 0x67}},
	{APADDB, ymm, Py1, opBytes{0xfc, Pe, 0xfc}},
	{APADDL, ymm, Py1, opBytes{0xfe, Pe, 0xfe}},
	{APADDQ, yxm, Pe, opBytes{0xd4}},
	{APADDSB, ymm, Py1, opBytes{0xec, Pe, 0xec}},
	{APADDSW, ymm, Py1, opBytes{0xed, Pe, 0xed}},
	{APADDUSB, ymm, Py1, opBytes{0xdc, Pe, 0xdc}},
	{APADDUSW, ymm, Py1, opBytes{0xdd, Pe, 0xdd}},
	{APADDW, ymm, Py1, opBytes{0xfd, Pe, 0xfd}},
	{APALIGNR, ypalignr, Pq, opBytes{0x3a, 0x0f}},
	{APAND, ymm, Py1, opBytes{0xdb, Pe, 0xdb}},
	{APANDN, ymm, Py1, opBytes{0xdf, Pe, 0xdf}},
	{APAUSE, ynone, Px, opBytes{0xf3, 0x90}},
	{APAVGB, ymm, Py1, opBytes{0xe0, Pe, 0xe0}},
	{APAVGW, ymm, Py1, opBytes{0xe3, Pe, 0xe3}},
	{APBLENDW, yxshuf, Pq, opBytes{0x3a, 0x0e, 0}},
	{APCMPEQB, ymm, Py1, opBytes{0x74, Pe, 0x74}},
	{APCMPEQL, ymm, Py1, opBytes{0x76, Pe, 0x76}},
	{APCMPEQQ, yxm_q4, Pq4, opBytes{0x29}},
	{APCMPEQW, ymm, Py1, opBytes{0x75, Pe, 0x75}},
	{APCMPGTB, ymm, Py1, opBytes{0x64, Pe, 0x64}},
	{APCMPGTL, ymm, Py1, opBytes{0x66, Pe, 0x66}},
	{APCMPGTQ, yxm_q4, Pq4, opBytes{0x37}},
	{APCMPGTW, ymm, Py1, opBytes{0x65, Pe, 0x65}},
	{APCMPISTRI, yxshuf, Pq, opBytes{0x3a, 0x63, 0}},
	{APCMPISTRM, yxshuf, Pq, opBytes{0x3a, 0x62, 0}},
	{APEXTRW, yextrw, Pq, opBytes{0xc5, 0, 0x3a, 0x15, 0}},
	{APEXTRB, yextr, Pq, opBytes{0x3a, 0x14, 00}},
	{APEXTRD, yextr, Pq, opBytes{0x3a, 0x16, 00}},
	{APEXTRQ, yextr, Pq3, opBytes{0x3a, 0x16, 00}},
	{APHADDD, ymmxmm0f38, Px, opBytes{0x0F, 0x38, 0x02, 0, 0x66, 0x0F, 0x38, 0x02, 0}},
	{APHADDSW, yxm_q4, Pq4, opBytes{0x03}},
	{APHADDW, yxm_q4, Pq4, opBytes{0x01}},
	{APHMINPOSUW, yxm_q4, Pq4, opBytes{0x41}},
	{APHSUBD, yxm_q4, Pq4, opBytes{0x06}},
	{APHSUBSW, yxm_q4, Pq4, opBytes{0x07}},
	{APHSUBW, yxm_q4, Pq4, opBytes{0x05}},
	{APINSRW, yinsrw, Pq, opBytes{0xc4, 00}},
	{APINSRB, yinsr, Pq, opBytes{0x3a, 0x20, 00}},
	{APINSRD, yinsr, Pq, opBytes{0x3a, 0x22, 00}},
	{APINSRQ, yinsr, Pq3, opBytes{0x3a, 0x22, 00}},
	{APMADDUBSW, yxm_q4, Pq4, opBytes{0x04}},
	{APMADDWL, ymm, Py1, opBytes{0xf5, Pe, 0xf5}},
	{APMAXSB, yxm_q4, Pq4, opBytes{0x3c}},
	{APMAXSD, yxm_q4, Pq4, opBytes{0x3d}},
	{APMAXSW, yxm, Pe, opBytes{0xee}},
	{APMAXUB, yxm, Pe, opBytes{0xde}},
	{APMAXUD, yxm_q4, Pq4, opBytes{0x3f}},
	{APMAXUW, yxm_q4, Pq4, opBytes{0x3e}},
	{APMINSB, yxm_q4, Pq4, opBytes{0x38}},
	{APMINSD, yxm_q4, Pq4, opBytes{0x39}},
	{APMINSW, yxm, Pe, opBytes{0xea}},
	{APMINUB, yxm, Pe, opBytes{0xda}},
	{APMINUD, yxm_q4, Pq4, opBytes{0x3b}},
	{APMINUW, yxm_q4, Pq4, opBytes{0x3a}},
	{APMOVMSKB, ymskb, Px, opBytes{Pe, 0xd7, 0xd7}},
	{APMOVSXBD, yxm_q4, Pq4, opBytes{0x21}},
	{APMOVSXBQ, yxm_q4, Pq4, opBytes{0x22}},
	{APMOVSXBW, yxm_q4, Pq4, opBytes{0x20}},
	{APMOVSXDQ, yxm_q4, Pq4, opBytes{0x25}},
	{APMOVSXWD, yxm_q4, Pq4, opBytes{0x23}},
	{APMOVSXWQ, yxm_q4, Pq4, opBytes{0x24}},
	{APMOVZXBD, yxm_q4, Pq4, opBytes{0x31}},
	{APMOVZXBQ, yxm_q4, Pq4, opBytes{0x32}},
	{APMOVZXBW, yxm_q4, Pq4, opBytes{0x30}},
	{APMOVZXDQ, yxm_q4, Pq4, opBytes{0x35}},
	{APMOVZXWD, yxm_q4, Pq4, opBytes{0x33}},
	{APMOVZXWQ, yxm_q4, Pq4, opBytes{0x34}},
	{APMULDQ, yxm_q4, Pq4, opBytes{0x28}},
	{APMULHRSW, yxm_q4, Pq4, opBytes{0x0b}},
	{APMULHUW, ymm, Py1, opBytes{0xe4, Pe, 0xe4}},
	{APMULHW, ymm, Py1, opBytes{0xe5, Pe, 0xe5}},
	{APMULLD, yxm_q4, Pq4, opBytes{0x40}},
	{APMULLW, ymm, Py1, opBytes{0xd5, Pe, 0xd5}},
	{APMULULQ, ymm, Py1, opBytes{0xf4, Pe, 0xf4}},
	{APOPAL, ynone, P32, opBytes{0x61}},
	{APOPAW, ynone, Pe, opBytes{0x61}},
	{APOPCNTW, yml_rl, Pef3, opBytes{0xb8}},
	{APOPCNTL, yml_rl, Pf3, opBytes{0xb8}},
	{APOPCNTQ, yml_rl, Pfw, opBytes{0xb8}},
	{APOPFL, ynone, P32, opBytes{0x9d}},
	{APOPFQ, ynone, Py, opBytes{0x9d}},
	{APOPFW, ynone, Pe, opBytes{0x9d}},
	{APOPL, ypopl, P32, opBytes{0x58, 0x8f, 00}},
	{APOPQ, ypopl, Py, opBytes{0x58, 0x8f, 00}},
	{APOPW, ypopl, Pe, opBytes{0x58, 0x8f, 00}},
	{APOR, ymm, Py1, opBytes{0xeb, Pe, 0xeb}},
	{APSADBW, yxm, Pq, opBytes{0xf6}},
	{APSHUFHW, yxshuf, Pf3, opBytes{0x70, 00}},
	{APSHUFL, yxshuf, Pq, opBytes{0x70, 00}},
	{APSHUFLW, yxshuf, Pf2, opBytes{0x70, 00}},
	{APSHUFW, ymshuf, Pm, opBytes{0x70, 00}},
	{APSHUFB, ymshufb, Pq, opBytes{0x38, 0x00}},
	{APSIGNB, yxm_q4, Pq4, opBytes{0x08}},
	{APSIGND, yxm_q4, Pq4, opBytes{0x0a}},
	{APSIGNW, yxm_q4, Pq4, opBytes{0x09}},
	{APSLLO, ypsdq, Pq, opBytes{0x73, 07}},
	{APSLLL, yps, Py3, opBytes{0xf2, 0x72, 06, Pe, 0xf2, Pe, 0x72, 06}},
	{APSLLQ, yps, Py3, opBytes{0xf3, 0x73, 06, Pe, 0xf3, Pe, 0x73, 06}},
	{APSLLW, yps, Py3, opBytes{0xf1, 0x71, 06, Pe, 0xf1, Pe, 0x71, 06}},
	{APSRAL, yps, Py3, opBytes{0xe2, 0x72, 04, Pe, 0xe2, Pe, 0x72, 04}},
	{APSRAW, yps, Py3, opBytes{0xe1, 0x71, 04, Pe, 0xe1, Pe, 0x71, 04}},
	{APSRLO, ypsdq, Pq, opBytes{0x73, 03}},
	{APSRLL, yps, Py3, opBytes{0xd2, 0x72, 02, Pe, 0xd2, Pe, 0x72, 02}},
	{APSRLQ, yps, Py3, opBytes{0xd3, 0x73, 02, Pe, 0xd3, Pe, 0x73, 02}},
	{APSRLW, yps, Py3, opBytes{0xd1, 0x71, 02, Pe, 0xd1, Pe, 0x71, 02}},
	{APSUBB, yxm, Pe, opBytes{0xf8}},
	{APSUBL, yxm, Pe, opBytes{0xfa}},
	{APSUBQ, yxm, Pe, opBytes{0xfb}},
	{APSUBSB, yxm, Pe, opBytes{0xe8}},
	{APSUBSW, yxm, Pe, opBytes{0xe9}},
	{APSUBUSB, yxm, Pe, opBytes{0xd8}},
	{APSUBUSW, yxm, Pe, opBytes{0xd9}},
	{APSUBW, yxm, Pe, opBytes{0xf9}},
	{APTEST, yxm_q4, Pq4, opBytes{0x17}},
	{APUNPCKHBW, ymm, Py1, opBytes{0x68, Pe, 0x68}},
	{APUNPCKHLQ, ymm, Py1, opBytes{0x6a, Pe, 0x6a}},
	{APUNPCKHQDQ, yxm, Pe, opBytes{0x6d}},
	{APUNPCKHWL, ymm, Py1, opBytes{0x69, Pe, 0x69}},
	{APUNPCKLBW, ymm, Py1, opBytes{0x60, Pe, 0x60}},
	{APUNPCKLLQ, ymm, Py1, opBytes{0x62, Pe, 0x62}},
	{APUNPCKLQDQ, yxm, Pe, opBytes{0x6c}},
	{APUNPCKLWL, ymm, Py1, opBytes{0x61, Pe, 0x61}},
	{APUSHAL, ynone, P32, opBytes{0x60}},
	{APUSHAW, ynone, Pe, opBytes{0x60}},
	{APUSHFL, ynone, P32, opBytes{0x9c}},
	{APUSHFQ, ynone, Py, opBytes{0x9c}},
	{APUSHFW, ynone, Pe, opBytes{0x9c}},
	{APUSHL, ypushl, P32, opBytes{0x50, 0xff, 06, 0x6a, 0x68}},
	{APUSHQ, ypushl, Py, opBytes{0x50, 0xff, 06, 0x6a, 0x68}},
	{APUSHW, ypushl, Pe, opBytes{0x50, 0xff, 06, 0x6a, 0x68}},
	{APXOR, ymm, Py1, opBytes{0xef, Pe, 0xef}},
	{AQUAD, ybyte, Px, opBytes{8}},
	{ARCLB, yshb, Pb, opBytes{0xd0, 02, 0xc0, 02, 0xd2, 02}},
	{ARCLL, yshl, Px, opBytes{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCLQ, yshl, Pw, opBytes{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCLW, yshl, Pe, opBytes{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCPPS, yxm, Pm, opBytes{0x53}},
	{ARCPSS, yxm, Pf3, opBytes{0x53}},
	{ARCRB, yshb, Pb, opBytes{0xd0, 03, 0xc0, 03, 0xd2, 03}},
	{ARCRL, yshl, Px, opBytes{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{ARCRQ, yshl, Pw, opBytes{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{ARCRW, yshl, Pe, opBytes{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{AREP, ynone, Px, opBytes{0xf3}},
	{AREPN, ynone, Px, opBytes{0xf2}},
	{ARET, ynone, Px, opBytes{0xc3}},
	{ARETFW, yret, Pe, opBytes{0xcb, 0xca}},
	{ARETFL, yret, Px, opBytes{0xcb, 0xca}},
	{ARETFQ, yret, Pw, opBytes{0xcb, 0xca}},
	{AROLB, yshb, Pb, opBytes{0xd0, 00, 0xc0, 00, 0xd2, 00}},
	{AROLL, yshl, Px, opBytes{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{AROLQ, yshl, Pw, opBytes{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{AROLW, yshl, Pe, opBytes{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{ARORB, yshb, Pb, opBytes{0xd0, 01, 0xc0, 01, 0xd2, 01}},
	{ARORL, yshl, Px, opBytes{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARORQ, yshl, Pw, opBytes{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARORW, yshl, Pe, opBytes{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARSQRTPS, yxm, Pm, opBytes{0x52}},
	{ARSQRTSS, yxm, Pf3, opBytes{0x52}},
	{ASAHF, ynone, Px, opBytes{0x9e, 00, 0x86, 0xe0, 0x50, 0x9d}}, // XCHGB AH,AL; PUSH AX; POPFL
	{ASALB, yshb, Pb, opBytes{0xd0, 04, 0xc0, 04, 0xd2, 04}},
	{ASALL, yshl, Px, opBytes{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASALQ, yshl, Pw, opBytes{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASALW, yshl, Pe, opBytes{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASARB, yshb, Pb, opBytes{0xd0, 07, 0xc0, 07, 0xd2, 07}},
	{ASARL, yshl, Px, opBytes{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASARQ, yshl, Pw, opBytes{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASARW, yshl, Pe, opBytes{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASBBB, yxorb, Pb, opBytes{0x1c, 0x80, 03, 0x18, 0x1a}},
	{ASBBL, yaddl, Px, opBytes{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASBBQ, yaddl, Pw, opBytes{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASBBW, yaddl, Pe, opBytes{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASCASB, ynone, Pb, opBytes{0xae}},
	{ASCASL, ynone, Px, opBytes{0xaf}},
	{ASCASQ, ynone, Pw, opBytes{0xaf}},
	{ASCASW, ynone, Pe, opBytes{0xaf}},
	{ASETCC, yscond, Pb, opBytes{0x0f, 0x93, 00}},
	{ASETCS, yscond, Pb, opBytes{0x0f, 0x92, 00}},
	{ASETEQ, yscond, Pb, opBytes{0x0f, 0x94, 00}},
	{ASETGE, yscond, Pb, opBytes{0x0f, 0x9d, 00}},
	{ASETGT, yscond, Pb, opBytes{0x0f, 0x9f, 00}},
	{ASETHI, yscond, Pb, opBytes{0x0f, 0x97, 00}},
	{ASETLE, yscond, Pb, opBytes{0x0f, 0x9e, 00}},
	{ASETLS, yscond, Pb, opBytes{0x0f, 0x96, 00}},
	{ASETLT, yscond, Pb, opBytes{0x0f, 0x9c, 00}},
	{ASETMI, yscond, Pb, opBytes{0x0f, 0x98, 00}},
	{ASETNE, yscond, Pb, opBytes{0x0f, 0x95, 00}},
	{ASETOC, yscond, Pb, opBytes{0x0f, 0x91, 00}},
	{ASETOS, yscond, Pb, opBytes{0x0f, 0x90, 00}},
	{ASETPC, yscond, Pb, opBytes{0x0f, 0x9b, 00}},
	{ASETPL, yscond, Pb, opBytes{0x0f, 0x99, 00}},
	{ASETPS, yscond, Pb, opBytes{0x0f, 0x9a, 00}},
	{ASHLB, yshb, Pb, opBytes{0xd0, 04, 0xc0, 04, 0xd2, 04}},
	{ASHLL, yshl, Px, opBytes{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHLQ, yshl, Pw, opBytes{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHLW, yshl, Pe, opBytes{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHRB, yshb, Pb, opBytes{0xd0, 05, 0xc0, 05, 0xd2, 05}},
	{ASHRL, yshl, Px, opBytes{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHRQ, yshl, Pw, opBytes{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHRW, yshl, Pe, opBytes{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHUFPD, yxshuf, Pq, opBytes{0xc6, 00}},
	{ASHUFPS, yxshuf, Pm, opBytes{0xc6, 00}},
	{ASQRTPD, yxm, Pe, opBytes{0x51}},
	{ASQRTPS, yxm, Pm, opBytes{0x51}},
	{ASQRTSD, yxm, Pf2, opBytes{0x51}},
	{ASQRTSS, yxm, Pf3, opBytes{0x51}},
	{ASTC, ynone, Px, opBytes{0xf9}},
	{ASTD, ynone, Px, opBytes{0xfd}},
	{ASTI, ynone, Px, opBytes{0xfb}},
	{ASTMXCSR, ysvrs_om, Pm, opBytes{0xae, 03, 0xae, 03}},
	{ASTOSB, ynone, Pb, opBytes{0xaa}},
	{ASTOSL, ynone, Px, opBytes{0xab}},
	{ASTOSQ, ynone, Pw, opBytes{0xab}},
	{ASTOSW, ynone, Pe, opBytes{0xab}},
	{ASUBB, yxorb, Pb, opBytes{0x2c, 0x80, 05, 0x28, 0x2a}},
	{ASUBL, yaddl, Px, opBytes{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASUBPD, yxm, Pe, opBytes{0x5c}},
	{ASUBPS, yxm, Pm, opBytes{0x5c}},
	{ASUBQ, yaddl, Pw, opBytes{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASUBSD, yxm, Pf2, opBytes{0x5c}},
	{ASUBSS, yxm, Pf3, opBytes{0x5c}},
	{ASUBW, yaddl, Pe, opBytes{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASWAPGS, ynone, Pm, opBytes{0x01, 0xf8}},
	{ASYSCALL, ynone, Px, opBytes{0x0f, 0x05}}, // fast syscall
	{ATESTB, yxorb, Pb, opBytes{0xa8, 0xf6, 00, 0x84, 0x84}},
	{ATESTL, ytestl, Px, opBytes{0xa9, 0xf7, 00, 0x85, 0x85}},
	{ATESTQ, ytestl, Pw, opBytes{0xa9, 0xf7, 00, 0x85, 0x85}},
	{ATESTW, ytestl, Pe, opBytes{0xa9, 0xf7, 00, 0x85, 0x85}},
	{ATPAUSE, ywrfsbase, Pq, opBytes{0xae, 06}},
	{AUCOMISD, yxm, Pe, opBytes{0x2e}},
	{AUCOMISS, yxm, Pm, opBytes{0x2e}},
	{AUNPCKHPD, yxm, Pe, opBytes{0x15}},
	{AUNPCKHPS, yxm, Pm, opBytes{0x15}},
	{AUNPCKLPD, yxm, Pe, opBytes{0x14}},
	{AUNPCKLPS, yxm, Pm, opBytes{0x14}},
	{AUMONITOR, ywrfsbase, Pf3, opBytes{0xae, 06}},
	{AVERR, ydivl, Pm, opBytes{0x00, 04}},
	{AVERW, ydivl, Pm, opBytes{0x00, 05}},
	{AWAIT, ynone, Px, opBytes{0x9b}},
	{AWORD, ybyte, Px, opBytes{2}},
	{AXCHGB, yml_mb, Pb, opBytes{0x86, 0x86}},
	{AXCHGL, yxchg, Px, opBytes{0x90, 0x90, 0x87, 0x87}},
	{AXCHGQ, yxchg, Pw, opBytes{0x90, 0x90, 0x87, 0x87}},
	{AXCHGW, yxchg, Pe, opBytes{0x90, 0x90, 0x87, 0x87}},
	{AXLAT, ynone, Px, opBytes{0xd7}},
	{AXORB, yxorb, Pb, opBytes{0x34, 0x80, 06, 0x30, 0x32}},
	{AXORL, yaddl, Px, opBytes{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AXORPD, yxm, Pe, opBytes{0x57}},
	{AXORPS, yxm, Pm, opBytes{0x57}},
	{AXORQ, yaddl, Pw, opBytes{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AXORW, yaddl, Pe, opBytes{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AFMOVB, yfmvx, Px, opBytes{0xdf, 04}},
	{AFMOVBP, yfmvp, Px, opBytes{0xdf, 06}},
	{AFMOVD, yfmvd, Px, opBytes{0xdd, 00, 0xdd, 02, 0xd9, 00, 0xdd, 02}},
	{AFMOVDP, yfmvdp, Px, opBytes{0xdd, 03, 0xdd, 03}},
	{AFMOVF, yfmvf, Px, opBytes{0xd9, 00, 0xd9, 02}},
	{AFMOVFP, yfmvp, Px, opBytes{0xd9, 03}},
	{AFMOVL, yfmvf, Px, opBytes{0xdb, 00, 0xdb, 02}},
	{AFMOVLP, yfmvp, Px, opBytes{0xdb, 03}},
	{AFMOVV, yfmvx, Px, opBytes{0xdf, 05}},
	{AFMOVVP, yfmvp, Px, opBytes{0xdf, 07}},
	{AFMOVW, yfmvf, Px, opBytes{0xdf, 00, 0xdf, 02}},
	{AFMOVWP, yfmvp, Px, opBytes{0xdf, 03}},
	{AFMOVX, yfmvx, Px, opBytes{0xdb, 05}},
	{AFMOVXP, yfmvp, Px, opBytes{0xdb, 07}},
	{AFCMOVCC, yfcmv, Px, opBytes{0xdb, 00}},
	{AFCMOVCS, yfcmv, Px, opBytes{0xda, 00}},
	{AFCMOVEQ, yfcmv, Px, opBytes{0xda, 01}},
	{AFCMOVHI, yfcmv, Px, opBytes{0xdb, 02}},
	{AFCMOVLS, yfcmv, Px, opBytes{0xda, 02}},
	{AFCMOVB, yfcmv, Px, opBytes{0xda, 00}},
	{AFCMOVBE, yfcmv, Px, opBytes{0xda, 02}},
	{AFCMOVNB, yfcmv, Px, opBytes{0xdb, 00}},
	{AFCMOVNBE, yfcmv, Px, opBytes{0xdb, 02}},
	{AFCMOVE, yfcmv, Px, opBytes{0xda, 01}},
	{AFCMOVNE, yfcmv, Px, opBytes{0xdb, 01}},
	{AFCMOVNU, yfcmv, Px, opBytes{0xdb, 03}},
	{AFCMOVU, yfcmv, Px, opBytes{0xda, 03}},
	{AFCMOVUN, yfcmv, Px, opBytes{0xda, 03}},
	{AFCOMD, yfadd, Px, opBytes{0xdc, 02, 0xd8, 02, 0xdc, 02}},  // botch
	{AFCOMDP, yfadd, Px, opBytes{0xdc, 03, 0xd8, 03, 0xdc, 03}}, // botch
	{AFCOMDPP, ycompp, Px, opBytes{0xde, 03}},
	{AFCOMF, yfmvx, Px, opBytes{0xd8, 02}},
	{AFCOMFP, yfmvx, Px, opBytes{0xd8, 03}},
	{AFCOMI, yfcmv, Px, opBytes{0xdb, 06}},
	{AFCOMIP, yfcmv, Px, opBytes{0xdf, 06}},
	{AFCOML, yfmvx, Px, opBytes{0xda, 02}},
	{AFCOMLP, yfmvx, Px, opBytes{0xda, 03}},
	{AFCOMW, yfmvx, Px, opBytes{0xde, 02}},
	{AFCOMWP, yfmvx, Px, opBytes{0xde, 03}},
	{AFUCOM, ycompp, Px, opBytes{0xdd, 04}},
	{AFUCOMI, ycompp, Px, opBytes{0xdb, 05}},
	{AFUCOMIP, ycompp, Px, opBytes{0xdf, 05}},
	{AFUCOMP, ycompp, Px, opBytes{0xdd, 05}},
	{AFUCOMPP, ycompp, Px, opBytes{0xda, 13}},
	{AFADDDP, ycompp, Px, opBytes{0xde, 00}},
	{AFADDW, yfmvx, Px, opBytes{0xde, 00}},
	{AFADDL, yfmvx, Px, opBytes{0xda, 00}},
	{AFADDF, yfmvx, Px, opBytes{0xd8, 00}},
	{AFADDD, yfadd, Px, opBytes{0xdc, 00, 0xd8, 00, 0xdc, 00}},
	{AFMULDP, ycompp, Px, opBytes{0xde, 01}},
	{AFMULW, yfmvx, Px, opBytes{0xde, 01}},
	{AFMULL, yfmvx, Px, opBytes{0xda, 01}},
	{AFMULF, yfmvx, Px, opBytes{0xd8, 01}},
	{AFMULD, yfadd, Px, opBytes{0xdc, 01, 0xd8, 01, 0xdc, 01}},
	{AFSUBDP, ycompp, Px, opBytes{0xde, 05}},
	{AFSUBW, yfmvx, Px, opBytes{0xde, 04}},
	{AFSUBL, yfmvx, Px, opBytes{0xda, 04}},
	{AFSUBF, yfmvx, Px, opBytes{0xd8, 04}},
	{AFSUBD, yfadd, Px, opBytes{0xdc, 04, 0xd8, 04, 0xdc, 05}},
	{AFSUBRDP, ycompp, Px, opBytes{0xde, 04}},
	{AFSUBRW, yfmvx, Px, opBytes{0xde, 05}},
	{AFSUBRL, yfmvx, Px, opBytes{0xda, 05}},
	{AFSUBRF, yfmvx, Px, opBytes{0xd8, 05}},
	{AFSUBRD, yfadd, Px, opBytes{0xdc, 05, 0xd8, 05, 0xdc, 04}},
	{AFDIVDP, ycompp, Px, opBytes{0xde, 07}},
	{AFDIVW, yfmvx, Px, opBytes{0xde, 06}},
	{AFDIVL, yfmvx, Px, opBytes{0xda, 06}},
	{AFDIVF, yfmvx, Px, opBytes{0xd8, 06}},
	{AFDIVD, yfadd, Px, opBytes{0xdc, 06, 0xd8, 06, 0xdc, 07}},
	{AFDIVRDP, ycompp, Px, opBytes{0xde, 06}},
	{AFDIVRW, yfmvx, Px, opBytes{0xde, 07}},
	{AFDIVRL, yfmvx, Px, opBytes{0xda, 07}},
	{AFDIVRF, yfmvx, Px, opBytes{0xd8, 07}},
	{AFDIVRD, yfadd, Px, opBytes{0xdc, 07, 0xd8, 07, 0xdc, 06}},
	{AFXCHD, yfxch, Px, opBytes{0xd9, 01, 0xd9, 01}},
	{AFFREE, nil, 0, opBytes{}},
	{AFLDCW, ysvrs_mo, Px, opBytes{0xd9, 05, 0xd9, 05}},
	{AFLDENV, ysvrs_mo, Px, opBytes{0xd9, 04, 0xd9, 04}},
	{AFRSTOR, ysvrs_mo, Px, opBytes{0xdd, 04, 0xdd, 04}},
	{AFSAVE, ysvrs_om, Px, opBytes{0xdd, 06, 0xdd, 06}},
	{AFSTCW, ysvrs_om, Px, opBytes{0xd9, 07, 0xd9, 07}},
	{AFSTENV, ysvrs_om, Px, opBytes{0xd9, 06, 0xd9, 06}},
	{AFSTSW, ystsw, Px, opBytes{0xdd, 07, 0xdf, 0xe0}},
	{AF2XM1, ynone, Px, opBytes{0xd9, 0xf0}},
	{AFABS, ynone, Px, opBytes{0xd9, 0xe1}},
	{AFBLD, ysvrs_mo, Px, opBytes{0xdf, 04}},
	{AFBSTP, yclflush, Px, opBytes{0xdf, 06}},
	{AFCHS, ynone, Px, opBytes{0xd9, 0xe0}},
	{AFCLEX, ynone, Px, opBytes{0xdb, 0xe2}},
	{AFCOS, ynone, Px, opBytes{0xd9, 0xff}},
	{AFDECSTP, ynone, Px, opBytes{0xd9, 0xf6}},
	{AFINCSTP, ynone, Px, opBytes{0xd9, 0xf7}},
	{AFINIT, ynone, Px, opBytes{0xdb, 0xe3}},
	{AFLD1, ynone, Px, opBytes{0xd9, 0xe8}},
	{AFLDL2E, ynone, Px, opBytes{0xd9, 0xea}},
	{AFLDL2T, ynone, Px, opBytes{0xd9, 0xe9}},
	{AFLDLG2, ynone, Px, opBytes{0xd9, 0xec}},
	{AFLDLN2, ynone, Px, opBytes{0xd9, 0xed}},
	{AFLDPI, ynone, Px, opBytes{0xd9, 0xeb}},
	{AFLDZ, ynone, Px, opBytes{0xd9, 0xee}},
	{AFNOP, ynone, Px, opBytes{0xd9, 0xd0}},
	{AFPATAN, ynone, Px, opBytes{0xd9, 0xf3}},
	{AFPREM, ynone, Px, opBytes{0xd9, 0xf8}},
	{AFPREM1, ynone, Px, opBytes{0xd9, 0xf5}},
	{AFPTAN, ynone, Px, opBytes{0xd9, 0xf2}},
	{AFRNDINT, ynone, Px, opBytes{0xd9, 0xfc}},
	{AFSCALE, ynone, Px, opBytes{0xd9, 0xfd}},
	{AFSIN, ynone, Px, opBytes{0xd9, 0xfe}},
	{AFSINCOS, ynone, Px, opBytes{0xd9, 0xfb}},
	{AFSQRT, ynone, Px, opBytes{0xd9, 0xfa}},
	{AFTST, ynone, Px, opBytes{0xd9, 0xe4}},
	{AFXAM, ynone, Px, opBytes{0xd9, 0xe5}},
	{AFXTRACT, ynone, Px, opBytes{0xd9, 0xf4}},
	{AFYL2X, ynone, Px, opBytes{0xd9, 0xf1}},
	{AFYL2XP1, ynone, Px, opBytes{0xd9, 0xf9}},
	{ACMPXCHGB, yrb_mb, Pb, opBytes{0x0f, 0xb0}},
	{ACMPXCHGL, yrl_ml, Px, opBytes{0x0f, 0xb1}},
	{ACMPXCHGW, yrl_ml, Pe, opBytes{0x0f, 0xb1}},
	{ACMPXCHGQ, yrl_ml, Pw, opBytes{0x0f, 0xb1}},
	{ACMPXCHG8B, yscond, Pm, opBytes{0xc7, 01}},
	{ACMPXCHG16B, yscond, Pw, opBytes{0x0f, 0xc7, 01}},
	{AINVD, ynone, Pm, opBytes{0x08}},
	{AINVLPG, ydivb, Pm, opBytes{0x01, 07}},
	{AINVPCID, ycrc32l, Pe, opBytes{0x0f, 0x38, 0x82, 0}},
	{ALFENCE, ynone, Pm, opBytes{0xae, 0xe8}},
	{AMFENCE, ynone, Pm, opBytes{0xae, 0xf0}},
	{AMOVNTIL, yrl_ml, Pm, opBytes{0xc3}},
	{AMOVNTIQ, yrl_ml, Pw, opBytes{0x0f, 0xc3}},
	{ARDPKRU, ynone, Pm, opBytes{0x01, 0xee, 0}},
	{ARDMSR, ynone, Pm, opBytes{0x32}},
	{ARDPMC, ynone, Pm, opBytes{0x33}},
	{ARDTSC, ynone, Pm, opBytes{0x31}},
	{ARSM, ynone, Pm, opBytes{0xaa}},
	{ASFENCE, ynone, Pm, opBytes{0xae, 0xf8}},
	{ASYSRET, ynone, Pm, opBytes{0x07}},
	{AWBINVD, ynone, Pm, opBytes{0x09}},
	{AWRMSR, ynone, Pm, opBytes{0x30}},
	{AWRPKRU, ynone, Pm, opBytes{0x01, 0xef, 0}},
	{AXADDB, yrb_mb, Pb, opBytes{0x0f, 0xc0}},
	{AXADDL, yrl_ml, Px, opBytes{0x0f, 0xc1}},
	{AXADDQ, yrl_ml, Pw, opBytes{0x0f, 0xc1}},
	{AXADDW, yrl_ml, Pe, opBytes{0x0f, 0xc1}},
	{ACRC32B, ycrc32b, Px, opBytes{0xf2, 0x0f, 0x38, 0xf0, 0}},
	{ACRC32L, ycrc32l, Px, opBytes{0xf2, 0x0f, 0x38, 0xf1, 0}},
	{ACRC32Q, ycrc32l, Pw, opBytes{0xf2, 0x0f, 0x38, 0xf1, 0}},
	{ACRC32W, ycrc32l, Pe, opBytes{0xf2, 0x0f, 0x38, 0xf1, 0}},
	{APREFETCHT0, yprefetch, Pm, opBytes{0x18, 01}},
	{APREFETCHT1, yprefetch, Pm, opBytes{0x18, 02}},
	{APREFETCHT2, yprefetch, Pm, opBytes{0x18, 03}},
	{APREFETCHNTA, yprefetch, Pm, opBytes{0x18, 00}},
	{AMOVQL, yrl_ml, Px, opBytes{0x89}},
	{AUNDEF, ynone, Px, opBytes{0x0f, 0x0b}},
	{AAESENC, yaes, Pq, opBytes{0x38, 0xdc, 0}},
	{AAESENCLAST, yaes, Pq, opBytes{0x38, 0xdd, 0}},
	{AAESDEC, yaes, Pq, opBytes{0x38, 0xde, 0}},
	{AAESDECLAST, yaes, Pq, opBytes{0x38, 0xdf, 0}},
	{AAESIMC, yaes, Pq, opBytes{0x38, 0xdb, 0}},
	{AAESKEYGENASSIST, yxshuf, Pq, opBytes{0x3a, 0xdf, 0}},
	{AROUNDPD, yxshuf, Pq, opBytes{0x3a, 0x09, 0}},
	{AROUNDPS, yxshuf, Pq, opBytes{0x3a, 0x08, 0}},
	{AROUNDSD, yxshuf, Pq, opBytes{0x3a, 0x0b, 0}},
	{AROUNDSS, yxshuf, Pq, opBytes{0x3a, 0x0a, 0}},
	{APSHUFD, yxshuf, Pq, opBytes{0x70, 0}},
	{APCLMULQDQ, yxshuf, Pq, opBytes{0x3a, 0x44, 0}},
	{APCMPESTRI, yxshuf, Pq, opBytes{0x3a, 0x61, 0}},
	{APCMPESTRM, yxshuf, Pq, opBytes{0x3a, 0x60, 0}},
	{AMOVDDUP, yxm, Pf2, opBytes{0x12}},
	{AMOVSHDUP, yxm, Pf3, opBytes{0x16}},
	{AMOVSLDUP, yxm, Pf3, opBytes{0x12}},
	{ARDTSCP, ynone, Pm, opBytes{0x01, 0xf9, 0}},
	{ASTAC, ynone, Pm, opBytes{0x01, 0xcb, 0}},
	{AUD1, ynone, Pm, opBytes{0xb9, 0}},
	{AUD2, ynone, Pm, opBytes{0x0b, 0}},
	{AUMWAIT, ywrfsbase, Pf2, opBytes{0xae, 06}},
	{ASYSENTER, ynone, Px, opBytes{0x0f, 0x34, 0}},
	{ASYSENTER64, ynone, Pw, opBytes{0x0f, 0x34, 0}},
	{ASYSEXIT, ynone, Px, opBytes{0x0f, 0x35, 0}},
	{ASYSEXIT64, ynone, Pw, opBytes{0x0f, 0x35, 0}},
	{ALMSW, ydivl, Pm, opBytes{0x01, 06}},
	{ALLDT, ydivl, Pm, opBytes{0x00, 02}},
	{ALIDT, ysvrs_mo, Pm, opBytes{0x01, 03}},
	{ALGDT, ysvrs_mo, Pm, opBytes{0x01, 02}},
	{ATZCNTW, ycrc32l, Pe, opBytes{0xf3, 0x0f, 0xbc, 0}},
	{ATZCNTL, ycrc32l, Px, opBytes{0xf3, 0x0f, 0xbc, 0}},
	{ATZCNTQ, ycrc32l, Pw, opBytes{0xf3, 0x0f, 0xbc, 0}},
	{AXRSTOR, ydivl, Px, opBytes{0x0f, 0xae, 05}},
	{AXRSTOR64, ydivl, Pw, opBytes{0x0f, 0xae, 05}},
	{AXRSTORS, ydivl, Px, opBytes{0x0f, 0xc7, 03}},
	{AXRSTORS64, ydivl, Pw, opBytes{0x0f, 0xc7, 03}},
	{AXSAVE, yclflush, Px, opBytes{0x0f, 0xae, 04}},
	{AXSAVE64, yclflush, Pw, opBytes{0x0f, 0xae, 04}},
	{AXSAVEOPT, yclflush, Px, opBytes{0x0f, 0xae, 06}},
	{AXSAVEOPT64, yclflush, Pw, opBytes{0x0f, 0xae, 06}},
	{AXSAVEC, yclflush, Px, opBytes{0x0f, 0xc7, 04}},
	{AXSAVEC64, yclflush, Pw, opBytes{0x0f, 0xc7, 04}},
	{AXSAVES, yclflush, Px, opBytes{0x0f, 0xc7, 05}},
	{AXSAVES64, yclflush, Pw, opBytes{0x0f, 0xc7, 05}},
	{ASGDT, yclflush, Pm, opBytes{0x01, 00}},
	{ASIDT, yclflush, Pm, opBytes{0x01, 01}},
	{ARDRANDW, yrdrand, Pe, opBytes{0x0f, 0xc7, 06}},
	{ARDRANDL, yrdrand, Px, opBytes{0x0f, 0xc7, 06}},
	{ARDRANDQ, yrdrand, Pw, opBytes{0x0f, 0xc7, 06}},
	{ARDSEEDW, yrdrand, Pe, opBytes{0x0f, 0xc7, 07}},
	{ARDSEEDL, yrdrand, Px, opBytes{0x0f, 0xc7, 07}},
	{ARDSEEDQ, yrdrand, Pw, opBytes{0x0f, 0xc7, 07}},
	{ASTRW, yincq, Pe, opBytes{0x0f, 0x00, 01}},
	{ASTRL, yincq, Px, opBytes{0x0f, 0x00, 01}},
	{ASTRQ, yincq, Pw, opBytes{0x0f, 0x00, 01}},
	{AXSETBV, ynone, Pm, opBytes{0x01, 0xd1, 0}},
	{AMOVBEWW, ymovbe, Pq, opBytes{0x38, 0xf0, 0, 0x38, 0xf1, 0}},
	{AMOVBELL, ymovbe, Pm, opBytes{0x38, 0xf0, 0, 0x38, 0xf1, 0}},
	{AMOVBEQQ, ymovbe, Pw, opBytes{0x0f, 0x38, 0xf0, 0, 0x0f, 0x38, 0xf1, 0}},
	{ANOPW, ydivl, Pe, opBytes{0x0f, 0x1f, 00}},
	{ANOPL, ydivl, Px, opBytes{0x0f, 0x1f, 00}},
	{ASLDTW, yincq, Pe, opBytes{0x0f, 0x00, 00}},
	{ASLDTL, yincq, Px, opBytes{0x0f, 0x00, 00}},
	{ASLDTQ, yincq, Pw, opBytes{0x0f, 0x00, 00}},
	{ASMSWW, yincq, Pe, opBytes{0x0f, 0x01, 04}},
	{ASMSWL, yincq, Px, opBytes{0x0f, 0x01, 04}},
	{ASMSWQ, yincq, Pw, opBytes{0x0f, 0x01, 04}},
	{ABLENDVPS, yblendvpd, Pq4, opBytes{0x14}},
	{ABLENDVPD, yblendvpd, Pq4, opBytes{0x15}},
	{APBLENDVB, yblendvpd, Pq4, opBytes{0x10}},
	{ASHA1MSG1, yaes, Px, opBytes{0x0f, 0x38, 0xc9, 0}},
	{ASHA1MSG2, yaes, Px, opBytes{0x0f, 0x38, 0xca, 0}},
	{ASHA1NEXTE, yaes, Px, opBytes{0x0f, 0x38, 0xc8, 0}},
	{ASHA256MSG1, yaes, Px, opBytes{0x0f, 0x38, 0xcc, 0}},
	{ASHA256MSG2, yaes, Px, opBytes{0x0f, 0x38, 0xcd, 0}},
	{ASHA1RNDS4, ysha1rnds4, Pm, opBytes{0x3a, 0xcc, 0}},
	{ASHA256RNDS2, ysha256rnds2, Px, opBytes{0x0f, 0x38, 0xcb, 0}},
	{ARDFSBASEL, yrdrand, Pf3, opBytes{0xae, 00}},
	{ARDFSBASEQ, yrdrand, Pfw, opBytes{0xae, 00}},
	{ARDGSBASEL, yrdrand, Pf3, opBytes{0xae, 01}},
	{ARDGSBASEQ, yrdrand, Pfw, opBytes{0xae, 01}},
	{AWRFSBASEL, ywrfsbase, Pf3, opBytes{0xae, 02}},
	{AWRFSBASEQ, ywrfsbase, Pfw, opBytes{0xae, 02}},
	{AWRGSBASEL, ywrfsbase, Pf3, opBytes{0xae, 03}},
	{AWRGSBASEQ, ywrfsbase, Pfw, opBytes{0xae, 03}},
	{ALFSW, ym_rl, Pe, opBytes{0x0f, 0xb4}},
	{ALFSL, ym_rl, Px, opBytes{0x0f, 0xb4}},
	{ALFSQ, ym_rl, Pw, opBytes{0x0f, 0xb4}},
	{ALGSW, ym_rl, Pe, opBytes{0x0f, 0xb5}},
	{ALGSL, ym_rl, Px, opBytes{0x0f, 0xb5}},
	{ALGSQ, ym_rl, Pw, opBytes{0x0f, 0xb5}},
	{ALSSW, ym_rl, Pe, opBytes{0x0f, 0xb2}},
	{ALSSL, ym_rl, Px, opBytes{0x0f, 0xb2}},
	{ALSSQ, ym_rl, Pw, opBytes{0x0f, 0xb2}},

	{ABLENDPD, yxshuf, Pq, opBytes{0x3a, 0x0d, 0}},
	{ABLENDPS, yxshuf, Pq, opBytes{0x3a, 0x0c, 0}},
	{AXACQUIRE, ynone, Px, opBytes{0xf2}},
	{AXRELEASE, ynone, Px, opBytes{0xf3}},
	{AXBEGIN, yxbegin, Px, opBytes{0xc7, 0xf8}},
	{AXABORT, yxabort, Px, opBytes{0xc6, 0xf8}},
	{AXEND, ynone, Px, opBytes{0x0f, 01, 0xd5}},
	{AXTEST, ynone, Px, opBytes{0x0f, 01, 0xd6}},
	{AXGETBV, ynone, Pm, opBytes{01, 0xd0}},

	{0, nil, 0, opBytes{}},
}

var opindex [ALAST + 1]*Optab

const (
	movLit uint8 = iota // Like Zlit
	movRegMem
	movMemReg
	movRegMem2op
	movMemReg2op
	movFullPtr // Load full pointer, trash heap (unsupported)
	movDoubleShift
	movTLSReg
)

var ymovtab = []movtab{
	// push
	{APUSHL, Ycs, Ynone, Ynone, movLit, [4]uint8{0x0e, 0}},
	{APUSHL, Yss, Ynone, Ynone, movLit, [4]uint8{0x16, 0}},
	{APUSHL, Yds, Ynone, Ynone, movLit, [4]uint8{0x1e, 0}},
	{APUSHL, Yes, Ynone, Ynone, movLit, [4]uint8{0x06, 0}},
	{APUSHL, Yfs, Ynone, Ynone, movLit, [4]uint8{0x0f, 0xa0, 0}},
	{APUSHL, Ygs, Ynone, Ynone, movLit, [4]uint8{0x0f, 0xa8, 0}},
	{APUSHQ, Yfs, Ynone, Ynone, movLit, [4]uint8{0x0f, 0xa0, 0}},
	{APUSHQ, Ygs, Ynone, Ynone, movLit, [4]uint8{0x0f, 0xa8, 0}},
	{APUSHW, Ycs, Ynone, Ynone, movLit, [4]uint8{Pe, 0x0e, 0}},
	{APUSHW, Yss, Ynone, Ynone, movLit, [4]uint8{Pe, 0x16, 0}},
	{APUSHW, Yds, Ynone, Ynone, movLit, [4]uint8{Pe, 0x1e, 0}},
	{APUSHW, Yes, Ynone, Ynone, movLit, [4]uint8{Pe, 0x06, 0}},
	{APUSHW, Yfs, Ynone, Ynone, movLit, [4]uint8{Pe, 0x0f, 0xa0, 0}},
	{APUSHW, Ygs, Ynone, Ynone, movLit, [4]uint8{Pe, 0x0f, 0xa8, 0}},

	// pop
	{APOPL, Ynone, Ynone, Yds, movLit, [4]uint8{0x1f, 0}},
	{APOPL, Ynone, Ynone, Yes, movLit, [4]uint8{0x07, 0}},
	{APOPL, Ynone, Ynone, Yss, movLit, [4]uint8{0x17, 0}},
	{APOPL, Ynone, Ynone, Yfs, movLit, [4]uint8{0x0f, 0xa1, 0}},
	{APOPL, Ynone, Ynone, Ygs, movLit, [4]uint8{0x0f, 0xa9, 0}},
	{APOPQ, Ynone, Ynone, Yfs, movLit, [4]uint8{0x0f, 0xa1, 0}},
	{APOPQ, Ynone, Ynone, Ygs, movLit, [4]uint8{0x0f, 0xa9, 0}},
	{APOPW, Ynone, Ynone, Yds, movLit, [4]uint8{Pe, 0x1f, 0}},
	{APOPW, Ynone, Ynone, Yes, movLit, [4]uint8{Pe, 0x07, 0}},
	{APOPW, Ynone, Ynone, Yss, movLit, [4]uint8{Pe, 0x17, 0}},
	{APOPW, Ynone, Ynone, Yfs, movLit, [4]uint8{Pe, 0x0f, 0xa1, 0}},
	{APOPW, Ynone, Ynone, Ygs, movLit, [4]uint8{Pe, 0x0f, 0xa9, 0}},

	// mov seg
	{AMOVW, Yes, Ynone, Yml, movRegMem, [4]uint8{0x8c, 0, 0, 0}},
	{AMOVW, Ycs, Ynone, Yml, movRegMem, [4]uint8{0x8c, 1, 0, 0}},
	{AMOVW, Yss, Ynone, Yml, movRegMem, [4]uint8{0x8c, 2, 0, 0}},
	{AMOVW, Yds, Ynone, Yml, movRegMem, [4]uint8{0x8c, 3, 0, 0}},
	{AMOVW, Yfs, Ynone, Yml, movRegMem, [4]uint8{0x8c, 4, 0, 0}},
	{AMOVW, Ygs, Ynone, Yml, movRegMem, [4]uint8{0x8c, 5, 0, 0}},
	{AMOVW, Yml, Ynone, Yes, movMemReg, [4]uint8{0x8e, 0, 0, 0}},
	{AMOVW, Yml, Ynone, Ycs, movMemReg, [4]uint8{0x8e, 1, 0, 0}},
	{AMOVW, Yml, Ynone, Yss, movMemReg, [4]uint8{0x8e, 2, 0, 0}},
	{AMOVW, Yml, Ynone, Yds, movMemReg, [4]uint8{0x8e, 3, 0, 0}},
	{AMOVW, Yml, Ynone, Yfs, movMemReg, [4]uint8{0x8e, 4, 0, 0}},
	{AMOVW, Yml, Ynone, Ygs, movMemReg, [4]uint8{0x8e, 5, 0, 0}},

	// mov cr
	{AMOVL, Ycr0, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x20, 0, 0}},
	{AMOVL, Ycr2, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x20, 2, 0}},
	{AMOVL, Ycr3, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x20, 3, 0}},
	{AMOVL, Ycr4, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x20, 4, 0}},
	{AMOVL, Ycr8, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x20, 8, 0}},
	{AMOVQ, Ycr0, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x20, 0, 0}},
	{AMOVQ, Ycr2, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x20, 2, 0}},
	{AMOVQ, Ycr3, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x20, 3, 0}},
	{AMOVQ, Ycr4, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x20, 4, 0}},
	{AMOVQ, Ycr8, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x20, 8, 0}},
	{AMOVL, Yrl, Ynone, Ycr0, movMemReg2op, [4]uint8{0x0f, 0x22, 0, 0}},
	{AMOVL, Yrl, Ynone, Ycr2, movMemReg2op, [4]uint8{0x0f, 0x22, 2, 0}},
	{AMOVL, Yrl, Ynone, Ycr3, movMemReg2op, [4]uint8{0x0f, 0x22, 3, 0}},
	{AMOVL, Yrl, Ynone, Ycr4, movMemReg2op, [4]uint8{0x0f, 0x22, 4, 0}},
	{AMOVL, Yrl, Ynone, Ycr8, movMemReg2op, [4]uint8{0x0f, 0x22, 8, 0}},
	{AMOVQ, Yrl, Ynone, Ycr0, movMemReg2op, [4]uint8{0x0f, 0x22, 0, 0}},
	{AMOVQ, Yrl, Ynone, Ycr2, movMemReg2op, [4]uint8{0x0f, 0x22, 2, 0}},
	{AMOVQ, Yrl, Ynone, Ycr3, movMemReg2op, [4]uint8{0x0f, 0x22, 3, 0}},
	{AMOVQ, Yrl, Ynone, Ycr4, movMemReg2op, [4]uint8{0x0f, 0x22, 4, 0}},
	{AMOVQ, Yrl, Ynone, Ycr8, movMemReg2op, [4]uint8{0x0f, 0x22, 8, 0}},

	// mov dr
	{AMOVL, Ydr0, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x21, 0, 0}},
	{AMOVL, Ydr6, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x21, 6, 0}},
	{AMOVL, Ydr7, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x21, 7, 0}},
	{AMOVQ, Ydr0, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x21, 0, 0}},
	{AMOVQ, Ydr2, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x21, 2, 0}},
	{AMOVQ, Ydr3, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x21, 3, 0}},
	{AMOVQ, Ydr6, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x21, 6, 0}},
	{AMOVQ, Ydr7, Ynone, Yrl, movRegMem2op, [4]uint8{0x0f, 0x21, 7, 0}},
	{AMOVL, Yrl, Ynone, Ydr0, movMemReg2op, [4]uint8{0x0f, 0x23, 0, 0}},
	{AMOVL, Yrl, Ynone, Ydr6, movMemReg2op, [4]uint8{0x0f, 0x23, 6, 0}},
	{AMOVL, Yrl, Ynone, Ydr7, movMemReg2op, [4]uint8{0x0f, 0x23, 7, 0}},
	{AMOVQ, Yrl, Ynone, Ydr0, movMemReg2op, [4]uint8{0x0f, 0x23, 0, 0}},
	{AMOVQ, Yrl, Ynone, Ydr2, movMemReg2op, [4]uint8{0x0f, 0x23, 2, 0}},
	{AMOVQ, Yrl, Ynone, Ydr3, movMemReg2op, [4]uint8{0x0f, 0x23, 3, 0}},
	{AMOVQ, Yrl, Ynone, Ydr6, movMemReg2op, [4]uint8{0x0f, 0x23, 6, 0}},
	{AMOVQ, Yrl, Ynone, Ydr7, movMemReg2op, [4]uint8{0x0f, 0x23, 7, 0}},

	// mov tr
	{AMOVL, Ytr6, Ynone, Yml, movRegMem2op, [4]uint8{0x0f, 0x24, 6, 0}},
	{AMOVL, Ytr7, Ynone, Yml, movRegMem2op, [4]uint8{0x0f, 0x24, 7, 0}},
	{AMOVL, Yml, Ynone, Ytr6, movMemReg2op, [4]uint8{0x0f, 0x26, 6, 0xff}},
	{AMOVL, Yml, Ynone, Ytr7, movMemReg2op, [4]uint8{0x0f, 0x26, 7, 0xff}},

	// lgdt, sgdt, lidt, sidt
	{AMOVL, Ym, Ynone, Ygdtr, movMemReg2op, [4]uint8{0x0f, 0x01, 2, 0}},
	{AMOVL, Ygdtr, Ynone, Ym, movRegMem2op, [4]uint8{0x0f, 0x01, 0, 0}},
	{AMOVL, Ym, Ynone, Yidtr, movMemReg2op, [4]uint8{0x0f, 0x01, 3, 0}},
	{AMOVL, Yidtr, Ynone, Ym, movRegMem2op, [4]uint8{0x0f, 0x01, 1, 0}},
	{AMOVQ, Ym, Ynone, Ygdtr, movMemReg2op, [4]uint8{0x0f, 0x01, 2, 0}},
	{AMOVQ, Ygdtr, Ynone, Ym, movRegMem2op, [4]uint8{0x0f, 0x01, 0, 0}},
	{AMOVQ, Ym, Ynone, Yidtr, movMemReg2op, [4]uint8{0x0f, 0x01, 3, 0}},
	{AMOVQ, Yidtr, Ynone, Ym, movRegMem2op, [4]uint8{0x0f, 0x01, 1, 0}},

	// lldt, sldt
	{AMOVW, Yml, Ynone, Yldtr, movMemReg2op, [4]uint8{0x0f, 0x00, 2, 0}},
	{AMOVW, Yldtr, Ynone, Yml, movRegMem2op, [4]uint8{0x0f, 0x00, 0, 0}},

	// lmsw, smsw
	{AMOVW, Yml, Ynone, Ymsw, movMemReg2op, [4]uint8{0x0f, 0x01, 6, 0}},
	{AMOVW, Ymsw, Ynone, Yml, movRegMem2op, [4]uint8{0x0f, 0x01, 4, 0}},

	// ltr, str
	{AMOVW, Yml, Ynone, Ytask, movMemReg2op, [4]uint8{0x0f, 0x00, 3, 0}},
	{AMOVW, Ytask, Ynone, Yml, movRegMem2op, [4]uint8{0x0f, 0x00, 1, 0}},

	/* load full pointer - unsupported
	{AMOVL, Yml, Ycol, movFullPtr, [4]uint8{0, 0, 0, 0}},
	{AMOVW, Yml, Ycol, movFullPtr, [4]uint8{Pe, 0, 0, 0}},
	*/

	// double shift
	{ASHLL, Yi8, Yrl, Yml, movDoubleShift, [4]uint8{0xa4, 0xa5, 0, 0}},
	{ASHLL, Ycl, Yrl, Yml, movDoubleShift, [4]uint8{0xa4, 0xa5, 0, 0}},
	{ASHLL, Ycx, Yrl, Yml, movDoubleShift, [4]uint8{0xa4, 0xa5, 0, 0}},
	{ASHRL, Yi8, Yrl, Yml, movDoubleShift, [4]uint8{0xac, 0xad, 0, 0}},
	{ASHRL, Ycl, Yrl, Yml, movDoubleShift, [4]uint8{0xac, 0xad, 0, 0}},
	{ASHRL, Ycx, Yrl, Yml, movDoubleShift, [4]uint8{0xac, 0xad, 0, 0}},
	{ASHLQ, Yi8, Yrl, Yml, movDoubleShift, [4]uint8{Pw, 0xa4, 0xa5, 0}},
	{ASHLQ, Ycl, Yrl, Yml, movDoubleShift, [4]uint8{Pw, 0xa4, 0xa5, 0}},
	{ASHLQ, Ycx, Yrl, Yml, movDoubleShift, [4]uint8{Pw, 0xa4, 0xa5, 0}},
	{ASHRQ, Yi8, Yrl, Yml, movDoubleShift, [4]uint8{Pw, 0xac, 0xad, 0}},
	{ASHRQ, Ycl, Yrl, Yml, movDoubleShift, [4]uint8{Pw, 0xac, 0xad, 0}},
	{ASHRQ, Ycx, Yrl, Yml, movDoubleShift, [4]uint8{Pw, 0xac, 0xad, 0}},
	{ASHLW, Yi8, Yrl, Yml, movDoubleShift, [4]uint8{Pe, 0xa4, 0xa5, 0}},
	{ASHLW, Ycl, Yrl, Yml, movDoubleShift, [4]uint8{Pe, 0xa4, 0xa5, 0}},
	{ASHLW, Ycx, Yrl, Yml, movDoubleShift, [4]uint8{Pe, 0xa4, 0xa5, 0}},
	{ASHRW, Yi8, Yrl, Yml, movDoubleShift, [4]uint8{Pe, 0xac, 0xad, 0}},
	{ASHRW, Ycl, Yrl, Yml, movDoubleShift, [4]uint8{Pe, 0xac, 0xad, 0}},
	{ASHRW, Ycx, Yrl, Yml, movDoubleShift, [4]uint8{Pe, 0xac, 0xad, 0}},

	// load TLS base
	{AMOVL, Ytls, Ynone, Yrl, movTLSReg, [4]uint8{0, 0, 0, 0}},
	{AMOVQ, Ytls, Ynone, Yrl, movTLSReg, [4]uint8{0, 0, 0, 0}},
	{0, 0, 0, 0, 0, [4]uint8{}},
}
