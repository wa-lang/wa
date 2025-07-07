// Inferno utils/6l/span.c
// http://code.google.com/p/inferno-os/source/browse/utils/6l/span.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
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

package x86

import "wa-lang.org/wa/internal/p9asm/obj"

// Instruction layout.

const (
	MaxAlign = 32 // max data alignment

	// Loop alignment constants:
	// want to align loop entry to LoopAlign-byte boundary,
	// and willing to insert at most MaxLoopPad bytes of NOP to do so.
	// We define a loop entry as the target of a backward jump.
	//
	// gcc uses MaxLoopPad = 10 for its 'generic x86-64' config,
	// and it aligns all jump targets, not just backward jump targets.
	//
	// As of 6/1/2012, the effect of setting MaxLoopPad = 10 here
	// is very slight but negative, so the alignment is disabled by
	// setting MaxLoopPad = 0. The code is here for reference and
	// for future experiments.
	//
	LoopAlign  = 16
	MaxLoopPad = 0
	FuncAlign  = 16
)

type Optab struct {
	as     int16
	ytab   []ytab
	prefix uint8
	op     [23]uint8
}

type ytab struct {
	from    uint8
	from3   uint8
	to      uint8
	zcase   uint8
	zoffset uint8
}

type Movtab struct {
	as   int16
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
	Yxr
	Yxm
	Ytls
	Ytextsize
	Yindir
	Ymax
)

const (
	Zxxx = iota
	Zlit
	Zlitm_r
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
	Ziqo_m
	Zjmp
	Zjmpcon
	Zloop
	Zo_iw
	Zm_o
	Zm_r
	Zm2_r
	Zm_r_xm
	Zm_r_i_xm
	Zm_r_3d
	Zm_r_xm_nr
	Zr_m_xm_nr
	Zibm_r /* mmx1,mmx2/mem64,imm8 */
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
	Zclr
	Zbyte
	Zmax
)

const (
	Px  = 0
	Px1 = 1    // symbolic; exact value doesn't matter
	P32 = 0x32 /* 32-bit only */
	Pe  = 0x66 /* operand escape */
	Pm  = 0x0f /* 2byte opcode escape */
	Pq  = 0xff /* both escapes: 66 0f */
	Pb  = 0xfe /* byte operands */
	Pf2 = 0xf2 /* xmm escape 1: f2 0f */
	Pf3 = 0xf3 /* xmm escape 2: f3 0f */
	Pq3 = 0x67 /* xmm escape 3: 66 48 0f */
	Pw  = 0x48 /* Rex.w */
	Pw8 = 0x90 // symbolic; exact value doesn't matter
	Py  = 0x80 /* defaults to 64-bit mode */
	Py1 = 0x81 // symbolic; exact value doesn't matter
	Py3 = 0x83 // symbolic; exact value doesn't matter

	Rxf = 1 << 9 /* internal flag for Rxr on from */
	Rxt = 1 << 8 /* internal flag for Rxr on to */
	Rxw = 1 << 3 /* =1, 64-bit operand size */
	Rxr = 1 << 2 /* extend modrm reg */
	Rxx = 1 << 1 /* extend sib index */
	Rxb = 1 << 0 /* extend modrm r/m, sib base, or opcode reg */

	Maxand = 10 /* in -a output width of the byte codes */
)

var ycover [Ymax * Ymax]uint8

var reg [MAXREG]int

var regrex [MAXREG + 1]int

var ynone = []ytab{
	{Ynone, Ynone, Ynone, Zlit, 1},
}

var ysahf = []ytab{
	{Ynone, Ynone, Ynone, Zlit, 2},
	{Ynone, Ynone, Ynone, Zlit, 1},
}

var ytext = []ytab{
	{Ymb, Ynone, Ytextsize, Zpseudo, 0},
	{Ymb, Yi32, Ytextsize, Zpseudo, 1},
}

var ynop = []ytab{
	{Ynone, Ynone, Ynone, Zpseudo, 0},
	{Ynone, Ynone, Yiauto, Zpseudo, 0},
	{Ynone, Ynone, Yml, Zpseudo, 0},
	{Ynone, Ynone, Yrf, Zpseudo, 0},
	{Ynone, Ynone, Yxr, Zpseudo, 0},
	{Yiauto, Ynone, Ynone, Zpseudo, 0},
	{Yml, Ynone, Ynone, Zpseudo, 0},
	{Yrf, Ynone, Ynone, Zpseudo, 0},
	{Yxr, Ynone, Ynone, Zpseudo, 1},
}

var yfuncdata = []ytab{
	{Yi32, Ynone, Ym, Zpseudo, 0},
}

var ypcdata = []ytab{
	{Yi32, Ynone, Yi32, Zpseudo, 0},
}

var yxorb = []ytab{
	{Yi32, Ynone, Yal, Zib_, 1},
	{Yi32, Ynone, Ymb, Zibo_m, 2},
	{Yrb, Ynone, Ymb, Zr_m, 1},
	{Ymb, Ynone, Yrb, Zm_r, 1},
}

var yxorl = []ytab{
	{Yi8, Ynone, Yml, Zibo_m, 2},
	{Yi32, Ynone, Yax, Zil_, 1},
	{Yi32, Ynone, Yml, Zilo_m, 2},
	{Yrl, Ynone, Yml, Zr_m, 1},
	{Yml, Ynone, Yrl, Zm_r, 1},
}

var yaddl = []ytab{
	{Yi8, Ynone, Yml, Zibo_m, 2},
	{Yi32, Ynone, Yax, Zil_, 1},
	{Yi32, Ynone, Yml, Zilo_m, 2},
	{Yrl, Ynone, Yml, Zr_m, 1},
	{Yml, Ynone, Yrl, Zm_r, 1},
}

var yincb = []ytab{
	{Ynone, Ynone, Ymb, Zo_m, 2},
}

var yincw = []ytab{
	{Ynone, Ynone, Yml, Zo_m, 2},
}

var yincl = []ytab{
	{Ynone, Ynone, Yrl, Z_rp, 1},
	{Ynone, Ynone, Yml, Zo_m, 2},
}

var yincq = []ytab{
	{Ynone, Ynone, Yml, Zo_m, 2},
}

var ycmpb = []ytab{
	{Yal, Ynone, Yi32, Z_ib, 1},
	{Ymb, Ynone, Yi32, Zm_ibo, 2},
	{Ymb, Ynone, Yrb, Zm_r, 1},
	{Yrb, Ynone, Ymb, Zr_m, 1},
}

var ycmpl = []ytab{
	{Yml, Ynone, Yi8, Zm_ibo, 2},
	{Yax, Ynone, Yi32, Z_il, 1},
	{Yml, Ynone, Yi32, Zm_ilo, 2},
	{Yml, Ynone, Yrl, Zm_r, 1},
	{Yrl, Ynone, Yml, Zr_m, 1},
}

var yshb = []ytab{
	{Yi1, Ynone, Ymb, Zo_m, 2},
	{Yi32, Ynone, Ymb, Zibo_m, 2},
	{Ycx, Ynone, Ymb, Zo_m, 2},
}

var yshl = []ytab{
	{Yi1, Ynone, Yml, Zo_m, 2},
	{Yi32, Ynone, Yml, Zibo_m, 2},
	{Ycl, Ynone, Yml, Zo_m, 2},
	{Ycx, Ynone, Yml, Zo_m, 2},
}

var ytestb = []ytab{
	{Yi32, Ynone, Yal, Zib_, 1},
	{Yi32, Ynone, Ymb, Zibo_m, 2},
	{Yrb, Ynone, Ymb, Zr_m, 1},
	{Ymb, Ynone, Yrb, Zm_r, 1},
}

var ytestl = []ytab{
	{Yi32, Ynone, Yax, Zil_, 1},
	{Yi32, Ynone, Yml, Zilo_m, 2},
	{Yrl, Ynone, Yml, Zr_m, 1},
	{Yml, Ynone, Yrl, Zm_r, 1},
}

var ymovb = []ytab{
	{Yrb, Ynone, Ymb, Zr_m, 1},
	{Ymb, Ynone, Yrb, Zm_r, 1},
	{Yi32, Ynone, Yrb, Zib_rp, 1},
	{Yi32, Ynone, Ymb, Zibo_m, 2},
}

var ymbs = []ytab{
	{Ymb, Ynone, Ynone, Zm_o, 2},
}

var ybtl = []ytab{
	{Yi8, Ynone, Yml, Zibo_m, 2},
	{Yrl, Ynone, Yml, Zr_m, 1},
}

var ymovw = []ytab{
	{Yrl, Ynone, Yml, Zr_m, 1},
	{Yml, Ynone, Yrl, Zm_r, 1},
	{Yi0, Ynone, Yrl, Zclr, 1},
	{Yi32, Ynone, Yrl, Zil_rp, 1},
	{Yi32, Ynone, Yml, Zilo_m, 2},
	{Yiauto, Ynone, Yrl, Zaut_r, 2},
}

var ymovl = []ytab{
	{Yrl, Ynone, Yml, Zr_m, 1},
	{Yml, Ynone, Yrl, Zm_r, 1},
	{Yi0, Ynone, Yrl, Zclr, 1},
	{Yi32, Ynone, Yrl, Zil_rp, 1},
	{Yi32, Ynone, Yml, Zilo_m, 2},
	{Yml, Ynone, Ymr, Zm_r_xm, 1}, // MMX MOVD
	{Ymr, Ynone, Yml, Zr_m_xm, 1}, // MMX MOVD
	{Yml, Ynone, Yxr, Zm_r_xm, 2}, // XMM MOVD (32 bit)
	{Yxr, Ynone, Yml, Zr_m_xm, 2}, // XMM MOVD (32 bit)
	{Yiauto, Ynone, Yrl, Zaut_r, 2},
}

var yret = []ytab{
	{Ynone, Ynone, Ynone, Zo_iw, 1},
	{Yi32, Ynone, Ynone, Zo_iw, 1},
}

var ymovq = []ytab{
	// valid in 32-bit mode
	{Ym, Ynone, Ymr, Zm_r_xm_nr, 1},  // 0x6f MMX MOVQ (shorter encoding)
	{Ymr, Ynone, Ym, Zr_m_xm_nr, 1},  // 0x7f MMX MOVQ
	{Yxr, Ynone, Ymr, Zm_r_xm_nr, 2}, // Pf2, 0xd6 MOVDQ2Q
	{Yxm, Ynone, Yxr, Zm_r_xm_nr, 2}, // Pf3, 0x7e MOVQ xmm1/m64 -> xmm2
	{Yxr, Ynone, Yxm, Zr_m_xm_nr, 2}, // Pe, 0xd6 MOVQ xmm1 -> xmm2/m64

	// valid only in 64-bit mode, usually with 64-bit prefix
	{Yrl, Ynone, Yml, Zr_m, 1},      // 0x89
	{Yml, Ynone, Yrl, Zm_r, 1},      // 0x8b
	{Yi0, Ynone, Yrl, Zclr, 1},      // 0x31
	{Ys32, Ynone, Yrl, Zilo_m, 2},   // 32 bit signed 0xc7,(0)
	{Yi64, Ynone, Yrl, Ziq_rp, 1},   // 0xb8 -- 32/64 bit immediate
	{Yi32, Ynone, Yml, Zilo_m, 2},   // 0xc7,(0)
	{Ymm, Ynone, Ymr, Zm_r_xm, 1},   // 0x6e MMX MOVD
	{Ymr, Ynone, Ymm, Zr_m_xm, 1},   // 0x7e MMX MOVD
	{Yml, Ynone, Yxr, Zm_r_xm, 2},   // Pe, 0x6e MOVD xmm load
	{Yxr, Ynone, Yml, Zr_m_xm, 2},   // Pe, 0x7e MOVD xmm store
	{Yiauto, Ynone, Yrl, Zaut_r, 1}, // 0 built-in LEAQ
}

var ym_rl = []ytab{
	{Ym, Ynone, Yrl, Zm_r, 1},
}

var yrl_m = []ytab{
	{Yrl, Ynone, Ym, Zr_m, 1},
}

var ymb_rl = []ytab{
	{Ymb, Ynone, Yrl, Zmb_r, 1},
}

var yml_rl = []ytab{
	{Yml, Ynone, Yrl, Zm_r, 1},
}

var yrl_ml = []ytab{
	{Yrl, Ynone, Yml, Zr_m, 1},
}

var yml_mb = []ytab{
	{Yrb, Ynone, Ymb, Zr_m, 1},
	{Ymb, Ynone, Yrb, Zm_r, 1},
}

var yrb_mb = []ytab{
	{Yrb, Ynone, Ymb, Zr_m, 1},
}

var yxchg = []ytab{
	{Yax, Ynone, Yrl, Z_rp, 1},
	{Yrl, Ynone, Yax, Zrp_, 1},
	{Yrl, Ynone, Yml, Zr_m, 1},
	{Yml, Ynone, Yrl, Zm_r, 1},
}

var ydivl = []ytab{
	{Yml, Ynone, Ynone, Zm_o, 2},
}

var ydivb = []ytab{
	{Ymb, Ynone, Ynone, Zm_o, 2},
}

var yimul = []ytab{
	{Yml, Ynone, Ynone, Zm_o, 2},
	{Yi8, Ynone, Yrl, Zib_rr, 1},
	{Yi32, Ynone, Yrl, Zil_rr, 1},
	{Yml, Ynone, Yrl, Zm_r, 2},
}

var yimul3 = []ytab{
	{Yi8, Yml, Yrl, Zibm_r, 2},
}

var ybyte = []ytab{
	{Yi64, Ynone, Ynone, Zbyte, 1},
}

var yin = []ytab{
	{Yi32, Ynone, Ynone, Zib_, 1},
	{Ynone, Ynone, Ynone, Zlit, 1},
}

var yint = []ytab{
	{Yi32, Ynone, Ynone, Zib_, 1},
}

var ypushl = []ytab{
	{Yrl, Ynone, Ynone, Zrp_, 1},
	{Ym, Ynone, Ynone, Zm_o, 2},
	{Yi8, Ynone, Ynone, Zib_, 1},
	{Yi32, Ynone, Ynone, Zil_, 1},
}

var ypopl = []ytab{
	{Ynone, Ynone, Yrl, Z_rp, 1},
	{Ynone, Ynone, Ym, Zo_m, 2},
}

var ybswap = []ytab{
	{Ynone, Ynone, Yrl, Z_rp, 2},
}

var yscond = []ytab{
	{Ynone, Ynone, Ymb, Zo_m, 2},
}

var yjcond = []ytab{
	{Ynone, Ynone, Ybr, Zbr, 0},
	{Yi0, Ynone, Ybr, Zbr, 0},
	{Yi1, Ynone, Ybr, Zbr, 1},
}

var yloop = []ytab{
	{Ynone, Ynone, Ybr, Zloop, 1},
}

var ycall = []ytab{
	{Ynone, Ynone, Yml, Zcallindreg, 0},
	{Yrx, Ynone, Yrx, Zcallindreg, 2},
	{Ynone, Ynone, Yindir, Zcallind, 2},
	{Ynone, Ynone, Ybr, Zcall, 0},
	{Ynone, Ynone, Yi32, Zcallcon, 1},
}

var yduff = []ytab{
	{Ynone, Ynone, Yi32, Zcallduff, 1},
}

var yjmp = []ytab{
	{Ynone, Ynone, Yml, Zo_m64, 2},
	{Ynone, Ynone, Ybr, Zjmp, 0},
	{Ynone, Ynone, Yi32, Zjmpcon, 1},
}

var yfmvd = []ytab{
	{Ym, Ynone, Yf0, Zm_o, 2},
	{Yf0, Ynone, Ym, Zo_m, 2},
	{Yrf, Ynone, Yf0, Zm_o, 2},
	{Yf0, Ynone, Yrf, Zo_m, 2},
}

var yfmvdp = []ytab{
	{Yf0, Ynone, Ym, Zo_m, 2},
	{Yf0, Ynone, Yrf, Zo_m, 2},
}

var yfmvf = []ytab{
	{Ym, Ynone, Yf0, Zm_o, 2},
	{Yf0, Ynone, Ym, Zo_m, 2},
}

var yfmvx = []ytab{
	{Ym, Ynone, Yf0, Zm_o, 2},
}

var yfmvp = []ytab{
	{Yf0, Ynone, Ym, Zo_m, 2},
}

var yfcmv = []ytab{
	{Yrf, Ynone, Yf0, Zm_o, 2},
}

var yfadd = []ytab{
	{Ym, Ynone, Yf0, Zm_o, 2},
	{Yrf, Ynone, Yf0, Zm_o, 2},
	{Yf0, Ynone, Yrf, Zo_m, 2},
}

var yfaddp = []ytab{
	{Yf0, Ynone, Yrf, Zo_m, 2},
}

var yfxch = []ytab{
	{Yf0, Ynone, Yrf, Zo_m, 2},
	{Yrf, Ynone, Yf0, Zm_o, 2},
}

var ycompp = []ytab{
	{Yf0, Ynone, Yrf, Zo_m, 2}, /* botch is really f0,f1 */
}

var ystsw = []ytab{
	{Ynone, Ynone, Ym, Zo_m, 2},
	{Ynone, Ynone, Yax, Zlit, 1},
}

var ystcw = []ytab{
	{Ynone, Ynone, Ym, Zo_m, 2},
	{Ym, Ynone, Ynone, Zm_o, 2},
}

var ysvrs = []ytab{
	{Ynone, Ynone, Ym, Zo_m, 2},
	{Ym, Ynone, Ynone, Zm_o, 2},
}

var ymm = []ytab{
	{Ymm, Ynone, Ymr, Zm_r_xm, 1},
	{Yxm, Ynone, Yxr, Zm_r_xm, 2},
}

var yxm = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 1},
}

var yxcvm1 = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 2},
	{Yxm, Ynone, Ymr, Zm_r_xm, 2},
}

var yxcvm2 = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 2},
	{Ymm, Ynone, Yxr, Zm_r_xm, 2},
}

/*
var yxmq = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 2},
}
*/

var yxr = []ytab{
	{Yxr, Ynone, Yxr, Zm_r_xm, 1},
}

var yxr_ml = []ytab{
	{Yxr, Ynone, Yml, Zr_m_xm, 1},
}

var ymr = []ytab{
	{Ymr, Ynone, Ymr, Zm_r, 1},
}

var ymr_ml = []ytab{
	{Ymr, Ynone, Yml, Zr_m_xm, 1},
}

var yxcmp = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 1},
}

var yxcmpi = []ytab{
	{Yxm, Yxr, Yi8, Zm_r_i_xm, 2},
}

var yxmov = []ytab{
	{Yxm, Ynone, Yxr, Zm_r_xm, 1},
	{Yxr, Ynone, Yxm, Zr_m_xm, 1},
}

var yxcvfl = []ytab{
	{Yxm, Ynone, Yrl, Zm_r_xm, 1},
}

var yxcvlf = []ytab{
	{Yml, Ynone, Yxr, Zm_r_xm, 1},
}

var yxcvfq = []ytab{
	{Yxm, Ynone, Yrl, Zm_r_xm, 2},
}

var yxcvqf = []ytab{
	{Yml, Ynone, Yxr, Zm_r_xm, 2},
}

var yps = []ytab{
	{Ymm, Ynone, Ymr, Zm_r_xm, 1},
	{Yi8, Ynone, Ymr, Zibo_m_xm, 2},
	{Yxm, Ynone, Yxr, Zm_r_xm, 2},
	{Yi8, Ynone, Yxr, Zibo_m_xm, 3},
}

var yxrrl = []ytab{
	{Yxr, Ynone, Yrl, Zm_r, 1},
}

var ymfp = []ytab{
	{Ymm, Ynone, Ymr, Zm_r_3d, 1},
}

var ymrxr = []ytab{
	{Ymr, Ynone, Yxr, Zm_r, 1},
	{Yxm, Ynone, Yxr, Zm_r_xm, 1},
}

var ymshuf = []ytab{
	{Yi8, Ymm, Ymr, Zibm_r, 2},
}

var ymshufb = []ytab{
	{Yxm, Ynone, Yxr, Zm2_r, 2},
}

var yxshuf = []ytab{
	{Yu8, Yxm, Yxr, Zibm_r, 2},
}

var yextrw = []ytab{
	{Yu8, Yxr, Yrl, Zibm_r, 2},
}

var yinsrw = []ytab{
	{Yu8, Yml, Yxr, Zibm_r, 2},
}

var yinsr = []ytab{
	{Yu8, Ymm, Yxr, Zibm_r, 3},
}

var ypsdq = []ytab{
	{Yi8, Ynone, Yxr, Zibo_m, 2},
}

var ymskb = []ytab{
	{Yxr, Ynone, Yrl, Zm_r_xm, 2},
	{Ymr, Ynone, Yrl, Zm_r_xm, 1},
}

var ycrc32l = []ytab{
	{Yml, Ynone, Yrl, Zlitm_r, 0},
}

var yprefetch = []ytab{
	{Ym, Ynone, Ynone, Zm_o, 2},
}

var yaes = []ytab{
	{Yxm, Ynone, Yxr, Zlitm_r, 2},
}

var yaes2 = []ytab{
	{Yu8, Yxm, Yxr, Zibm_r, 2},
}

/*
 * You are doasm, holding in your hand a Prog* with p->as set to, say, ACRC32,
 * and p->from and p->to as operands (Addr*).  The linker scans optab to find
 * the entry with the given p->as and then looks through the ytable for that
 * instruction (the second field in the optab struct) for a line whose first
 * two values match the Ytypes of the p->from and p->to operands.  The function
 * oclass in span.c computes the specific Ytype of an operand and then the set
 * of more general Ytypes that it satisfies is implied by the ycover table, set
 * up in instinit.  For example, oclass distinguishes the constants 0 and 1
 * from the more general 8-bit constants, but instinit says
 *
 *        ycover[Yi0*Ymax + Ys32] = 1;
 *        ycover[Yi1*Ymax + Ys32] = 1;
 *        ycover[Yi8*Ymax + Ys32] = 1;
 *
 * which means that Yi0, Yi1, and Yi8 all count as Ys32 (signed 32)
 * if that's what an instruction can handle.
 *
 * In parallel with the scan through the ytable for the appropriate line, there
 * is a z pointer that starts out pointing at the strange magic byte list in
 * the Optab struct.  With each step past a non-matching ytable line, z
 * advances by the 4th entry in the line.  When a matching line is found, that
 * z pointer has the extra data to use in laying down the instruction bytes.
 * The actual bytes laid down are a function of the 3rd entry in the line (that
 * is, the Ztype) and the z bytes.
 *
 * For example, let's look at AADDL.  The optab line says:
 *        { AADDL,        yaddl,  Px, 0x83,(00),0x05,0x81,(00),0x01,0x03 },
 *
 * and yaddl says
 *        uchar   yaddl[] =
 *        {
 *                Yi8,    Yml,    Zibo_m, 2,
 *                Yi32,   Yax,    Zil_,   1,
 *                Yi32,   Yml,    Zilo_m, 2,
 *                Yrl,    Yml,    Zr_m,   1,
 *                Yml,    Yrl,    Zm_r,   1,
 *                0
 *        };
 *
 * so there are 5 possible types of ADDL instruction that can be laid down, and
 * possible states used to lay them down (Ztype and z pointer, assuming z
 * points at {0x83,(00),0x05,0x81,(00),0x01,0x03}) are:
 *
 *        Yi8, Yml -> Zibo_m, z (0x83, 00)
 *        Yi32, Yax -> Zil_, z+2 (0x05)
 *        Yi32, Yml -> Zilo_m, z+2+1 (0x81, 0x00)
 *        Yrl, Yml -> Zr_m, z+2+1+2 (0x01)
 *        Yml, Yrl -> Zm_r, z+2+1+2+1 (0x03)
 *
 * The Pconstant in the optab line controls the prefix bytes to emit.  That's
 * relatively straightforward as this program goes.
 *
 * The switch on t[2] in doasm implements the various Z cases.  Zibo_m, for
 * example, is an opcode byte (z[0]) then an asmando (which is some kind of
 * encoded addressing mode for the Yml arg), and then a single immediate byte.
 * Zilo_m is the same but a long (32-bit) immediate.
 */
var optab = []Optab{
	// as, ytab, andproto, opcode
	{AAAA, ynone, P32, [23]uint8{0x37}},
	{AAAD, ynone, P32, [23]uint8{0xd5, 0x0a}},
	{AAAM, ynone, P32, [23]uint8{0xd4, 0x0a}},
	{AAAS, ynone, P32, [23]uint8{0x3f}},
	{AADCB, yxorb, Pb, [23]uint8{0x14, 0x80, 02, 0x10, 0x10}},
	{AADCL, yxorl, Px, [23]uint8{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADCQ, yxorl, Pw, [23]uint8{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADCW, yxorl, Pe, [23]uint8{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADDB, yxorb, Pb, [23]uint8{0x04, 0x80, 00, 0x00, 0x02}},
	{AADDL, yaddl, Px, [23]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADDPD, yxm, Pq, [23]uint8{0x58}},
	{AADDPS, yxm, Pm, [23]uint8{0x58}},
	{AADDQ, yaddl, Pw, [23]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADDSD, yxm, Pf2, [23]uint8{0x58}},
	{AADDSS, yxm, Pf3, [23]uint8{0x58}},
	{AADDW, yaddl, Pe, [23]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADJSP, nil, 0, [23]uint8{}},
	{AANDB, yxorb, Pb, [23]uint8{0x24, 0x80, 04, 0x20, 0x22}},
	{AANDL, yxorl, Px, [23]uint8{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AANDNPD, yxm, Pq, [23]uint8{0x55}},
	{AANDNPS, yxm, Pm, [23]uint8{0x55}},
	{AANDPD, yxm, Pq, [23]uint8{0x54}},
	{AANDPS, yxm, Pq, [23]uint8{0x54}},
	{AANDQ, yxorl, Pw, [23]uint8{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AANDW, yxorl, Pe, [23]uint8{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AARPL, yrl_ml, P32, [23]uint8{0x63}},
	{ABOUNDL, yrl_m, P32, [23]uint8{0x62}},
	{ABOUNDW, yrl_m, Pe, [23]uint8{0x62}},
	{ABSFL, yml_rl, Pm, [23]uint8{0xbc}},
	{ABSFQ, yml_rl, Pw, [23]uint8{0x0f, 0xbc}},
	{ABSFW, yml_rl, Pq, [23]uint8{0xbc}},
	{ABSRL, yml_rl, Pm, [23]uint8{0xbd}},
	{ABSRQ, yml_rl, Pw, [23]uint8{0x0f, 0xbd}},
	{ABSRW, yml_rl, Pq, [23]uint8{0xbd}},
	{ABSWAPL, ybswap, Px, [23]uint8{0x0f, 0xc8}},
	{ABSWAPQ, ybswap, Pw, [23]uint8{0x0f, 0xc8}},
	{ABTCL, ybtl, Pm, [23]uint8{0xba, 07, 0xbb}},
	{ABTCQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 07, 0x0f, 0xbb}},
	{ABTCW, ybtl, Pq, [23]uint8{0xba, 07, 0xbb}},
	{ABTL, ybtl, Pm, [23]uint8{0xba, 04, 0xa3}},
	{ABTQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 04, 0x0f, 0xa3}},
	{ABTRL, ybtl, Pm, [23]uint8{0xba, 06, 0xb3}},
	{ABTRQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 06, 0x0f, 0xb3}},
	{ABTRW, ybtl, Pq, [23]uint8{0xba, 06, 0xb3}},
	{ABTSL, ybtl, Pm, [23]uint8{0xba, 05, 0xab}},
	{ABTSQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 05, 0x0f, 0xab}},
	{ABTSW, ybtl, Pq, [23]uint8{0xba, 05, 0xab}},
	{ABTW, ybtl, Pq, [23]uint8{0xba, 04, 0xa3}},
	{ABYTE, ybyte, Px, [23]uint8{1}},
	{obj.ACALL, ycall, Px, [23]uint8{0xff, 02, 0xff, 0x15, 0xe8}},
	{ACDQ, ynone, Px, [23]uint8{0x99}},
	{ACLC, ynone, Px, [23]uint8{0xf8}},
	{ACLD, ynone, Px, [23]uint8{0xfc}},
	{ACLI, ynone, Px, [23]uint8{0xfa}},
	{ACLTS, ynone, Pm, [23]uint8{0x06}},
	{ACMC, ynone, Px, [23]uint8{0xf5}},
	{ACMOVLCC, yml_rl, Pm, [23]uint8{0x43}},
	{ACMOVLCS, yml_rl, Pm, [23]uint8{0x42}},
	{ACMOVLEQ, yml_rl, Pm, [23]uint8{0x44}},
	{ACMOVLGE, yml_rl, Pm, [23]uint8{0x4d}},
	{ACMOVLGT, yml_rl, Pm, [23]uint8{0x4f}},
	{ACMOVLHI, yml_rl, Pm, [23]uint8{0x47}},
	{ACMOVLLE, yml_rl, Pm, [23]uint8{0x4e}},
	{ACMOVLLS, yml_rl, Pm, [23]uint8{0x46}},
	{ACMOVLLT, yml_rl, Pm, [23]uint8{0x4c}},
	{ACMOVLMI, yml_rl, Pm, [23]uint8{0x48}},
	{ACMOVLNE, yml_rl, Pm, [23]uint8{0x45}},
	{ACMOVLOC, yml_rl, Pm, [23]uint8{0x41}},
	{ACMOVLOS, yml_rl, Pm, [23]uint8{0x40}},
	{ACMOVLPC, yml_rl, Pm, [23]uint8{0x4b}},
	{ACMOVLPL, yml_rl, Pm, [23]uint8{0x49}},
	{ACMOVLPS, yml_rl, Pm, [23]uint8{0x4a}},
	{ACMOVQCC, yml_rl, Pw, [23]uint8{0x0f, 0x43}},
	{ACMOVQCS, yml_rl, Pw, [23]uint8{0x0f, 0x42}},
	{ACMOVQEQ, yml_rl, Pw, [23]uint8{0x0f, 0x44}},
	{ACMOVQGE, yml_rl, Pw, [23]uint8{0x0f, 0x4d}},
	{ACMOVQGT, yml_rl, Pw, [23]uint8{0x0f, 0x4f}},
	{ACMOVQHI, yml_rl, Pw, [23]uint8{0x0f, 0x47}},
	{ACMOVQLE, yml_rl, Pw, [23]uint8{0x0f, 0x4e}},
	{ACMOVQLS, yml_rl, Pw, [23]uint8{0x0f, 0x46}},
	{ACMOVQLT, yml_rl, Pw, [23]uint8{0x0f, 0x4c}},
	{ACMOVQMI, yml_rl, Pw, [23]uint8{0x0f, 0x48}},
	{ACMOVQNE, yml_rl, Pw, [23]uint8{0x0f, 0x45}},
	{ACMOVQOC, yml_rl, Pw, [23]uint8{0x0f, 0x41}},
	{ACMOVQOS, yml_rl, Pw, [23]uint8{0x0f, 0x40}},
	{ACMOVQPC, yml_rl, Pw, [23]uint8{0x0f, 0x4b}},
	{ACMOVQPL, yml_rl, Pw, [23]uint8{0x0f, 0x49}},
	{ACMOVQPS, yml_rl, Pw, [23]uint8{0x0f, 0x4a}},
	{ACMOVWCC, yml_rl, Pq, [23]uint8{0x43}},
	{ACMOVWCS, yml_rl, Pq, [23]uint8{0x42}},
	{ACMOVWEQ, yml_rl, Pq, [23]uint8{0x44}},
	{ACMOVWGE, yml_rl, Pq, [23]uint8{0x4d}},
	{ACMOVWGT, yml_rl, Pq, [23]uint8{0x4f}},
	{ACMOVWHI, yml_rl, Pq, [23]uint8{0x47}},
	{ACMOVWLE, yml_rl, Pq, [23]uint8{0x4e}},
	{ACMOVWLS, yml_rl, Pq, [23]uint8{0x46}},
	{ACMOVWLT, yml_rl, Pq, [23]uint8{0x4c}},
	{ACMOVWMI, yml_rl, Pq, [23]uint8{0x48}},
	{ACMOVWNE, yml_rl, Pq, [23]uint8{0x45}},
	{ACMOVWOC, yml_rl, Pq, [23]uint8{0x41}},
	{ACMOVWOS, yml_rl, Pq, [23]uint8{0x40}},
	{ACMOVWPC, yml_rl, Pq, [23]uint8{0x4b}},
	{ACMOVWPL, yml_rl, Pq, [23]uint8{0x49}},
	{ACMOVWPS, yml_rl, Pq, [23]uint8{0x4a}},
	{ACMPB, ycmpb, Pb, [23]uint8{0x3c, 0x80, 07, 0x38, 0x3a}},
	{ACMPL, ycmpl, Px, [23]uint8{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACMPPD, yxcmpi, Px, [23]uint8{Pe, 0xc2}},
	{ACMPPS, yxcmpi, Pm, [23]uint8{0xc2, 0}},
	{ACMPQ, ycmpl, Pw, [23]uint8{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACMPSB, ynone, Pb, [23]uint8{0xa6}},
	{ACMPSD, yxcmpi, Px, [23]uint8{Pf2, 0xc2}},
	{ACMPSL, ynone, Px, [23]uint8{0xa7}},
	{ACMPSQ, ynone, Pw, [23]uint8{0xa7}},
	{ACMPSS, yxcmpi, Px, [23]uint8{Pf3, 0xc2}},
	{ACMPSW, ynone, Pe, [23]uint8{0xa7}},
	{ACMPW, ycmpl, Pe, [23]uint8{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACOMISD, yxcmp, Pe, [23]uint8{0x2f}},
	{ACOMISS, yxcmp, Pm, [23]uint8{0x2f}},
	{ACPUID, ynone, Pm, [23]uint8{0xa2}},
	{ACVTPL2PD, yxcvm2, Px, [23]uint8{Pf3, 0xe6, Pe, 0x2a}},
	{ACVTPL2PS, yxcvm2, Pm, [23]uint8{0x5b, 0, 0x2a, 0}},
	{ACVTPD2PL, yxcvm1, Px, [23]uint8{Pf2, 0xe6, Pe, 0x2d}},
	{ACVTPD2PS, yxm, Pe, [23]uint8{0x5a}},
	{ACVTPS2PL, yxcvm1, Px, [23]uint8{Pe, 0x5b, Pm, 0x2d}},
	{ACVTPS2PD, yxm, Pm, [23]uint8{0x5a}},
	{API2FW, ymfp, Px, [23]uint8{0x0c}},
	{ACVTSD2SL, yxcvfl, Pf2, [23]uint8{0x2d}},
	{ACVTSD2SQ, yxcvfq, Pw, [23]uint8{Pf2, 0x2d}},
	{ACVTSD2SS, yxm, Pf2, [23]uint8{0x5a}},
	{ACVTSL2SD, yxcvlf, Pf2, [23]uint8{0x2a}},
	{ACVTSQ2SD, yxcvqf, Pw, [23]uint8{Pf2, 0x2a}},
	{ACVTSL2SS, yxcvlf, Pf3, [23]uint8{0x2a}},
	{ACVTSQ2SS, yxcvqf, Pw, [23]uint8{Pf3, 0x2a}},
	{ACVTSS2SD, yxm, Pf3, [23]uint8{0x5a}},
	{ACVTSS2SL, yxcvfl, Pf3, [23]uint8{0x2d}},
	{ACVTSS2SQ, yxcvfq, Pw, [23]uint8{Pf3, 0x2d}},
	{ACVTTPD2PL, yxcvm1, Px, [23]uint8{Pe, 0xe6, Pe, 0x2c}},
	{ACVTTPS2PL, yxcvm1, Px, [23]uint8{Pf3, 0x5b, Pm, 0x2c}},
	{ACVTTSD2SL, yxcvfl, Pf2, [23]uint8{0x2c}},
	{ACVTTSD2SQ, yxcvfq, Pw, [23]uint8{Pf2, 0x2c}},
	{ACVTTSS2SL, yxcvfl, Pf3, [23]uint8{0x2c}},
	{ACVTTSS2SQ, yxcvfq, Pw, [23]uint8{Pf3, 0x2c}},
	{ACWD, ynone, Pe, [23]uint8{0x99}},
	{ACQO, ynone, Pw, [23]uint8{0x99}},
	{ADAA, ynone, P32, [23]uint8{0x27}},
	{ADAS, ynone, P32, [23]uint8{0x2f}},
	{obj.ADATA, nil, 0, [23]uint8{}},
	{ADECB, yincb, Pb, [23]uint8{0xfe, 01}},
	{ADECL, yincl, Px1, [23]uint8{0x48, 0xff, 01}},
	{ADECQ, yincq, Pw, [23]uint8{0xff, 01}},
	{ADECW, yincw, Pe, [23]uint8{0xff, 01}},
	{ADIVB, ydivb, Pb, [23]uint8{0xf6, 06}},
	{ADIVL, ydivl, Px, [23]uint8{0xf7, 06}},
	{ADIVPD, yxm, Pe, [23]uint8{0x5e}},
	{ADIVPS, yxm, Pm, [23]uint8{0x5e}},
	{ADIVQ, ydivl, Pw, [23]uint8{0xf7, 06}},
	{ADIVSD, yxm, Pf2, [23]uint8{0x5e}},
	{ADIVSS, yxm, Pf3, [23]uint8{0x5e}},
	{ADIVW, ydivl, Pe, [23]uint8{0xf7, 06}},
	{AEMMS, ynone, Pm, [23]uint8{0x77}},
	{AENTER, nil, 0, [23]uint8{}}, /* botch */
	{AFXRSTOR, ysvrs, Pm, [23]uint8{0xae, 01, 0xae, 01}},
	{AFXSAVE, ysvrs, Pm, [23]uint8{0xae, 00, 0xae, 00}},
	{AFXRSTOR64, ysvrs, Pw, [23]uint8{0x0f, 0xae, 01, 0x0f, 0xae, 01}},
	{AFXSAVE64, ysvrs, Pw, [23]uint8{0x0f, 0xae, 00, 0x0f, 0xae, 00}},
	{obj.AGLOBL, nil, 0, [23]uint8{}},
	{AHLT, ynone, Px, [23]uint8{0xf4}},
	{AIDIVB, ydivb, Pb, [23]uint8{0xf6, 07}},
	{AIDIVL, ydivl, Px, [23]uint8{0xf7, 07}},
	{AIDIVQ, ydivl, Pw, [23]uint8{0xf7, 07}},
	{AIDIVW, ydivl, Pe, [23]uint8{0xf7, 07}},
	{AIMULB, ydivb, Pb, [23]uint8{0xf6, 05}},
	{AIMULL, yimul, Px, [23]uint8{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMULQ, yimul, Pw, [23]uint8{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMULW, yimul, Pe, [23]uint8{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMUL3Q, yimul3, Pw, [23]uint8{0x6b, 00}},
	{AINB, yin, Pb, [23]uint8{0xe4, 0xec}},
	{AINCB, yincb, Pb, [23]uint8{0xfe, 00}},
	{AINCL, yincl, Px1, [23]uint8{0x40, 0xff, 00}},
	{AINCQ, yincq, Pw, [23]uint8{0xff, 00}},
	{AINCW, yincw, Pe, [23]uint8{0xff, 00}},
	{AINL, yin, Px, [23]uint8{0xe5, 0xed}},
	{AINSB, ynone, Pb, [23]uint8{0x6c}},
	{AINSL, ynone, Px, [23]uint8{0x6d}},
	{AINSW, ynone, Pe, [23]uint8{0x6d}},
	{AINT, yint, Px, [23]uint8{0xcd}},
	{AINTO, ynone, P32, [23]uint8{0xce}},
	{AINW, yin, Pe, [23]uint8{0xe5, 0xed}},
	{AIRETL, ynone, Px, [23]uint8{0xcf}},
	{AIRETQ, ynone, Pw, [23]uint8{0xcf}},
	{AIRETW, ynone, Pe, [23]uint8{0xcf}},
	{AJCC, yjcond, Px, [23]uint8{0x73, 0x83, 00}},
	{AJCS, yjcond, Px, [23]uint8{0x72, 0x82}},
	{AJCXZL, yloop, Px, [23]uint8{0xe3}},
	{AJCXZW, yloop, Px, [23]uint8{0xe3}},
	{AJCXZQ, yloop, Px, [23]uint8{0xe3}},
	{AJEQ, yjcond, Px, [23]uint8{0x74, 0x84}},
	{AJGE, yjcond, Px, [23]uint8{0x7d, 0x8d}},
	{AJGT, yjcond, Px, [23]uint8{0x7f, 0x8f}},
	{AJHI, yjcond, Px, [23]uint8{0x77, 0x87}},
	{AJLE, yjcond, Px, [23]uint8{0x7e, 0x8e}},
	{AJLS, yjcond, Px, [23]uint8{0x76, 0x86}},
	{AJLT, yjcond, Px, [23]uint8{0x7c, 0x8c}},
	{AJMI, yjcond, Px, [23]uint8{0x78, 0x88}},
	{obj.AJMP, yjmp, Px, [23]uint8{0xff, 04, 0xeb, 0xe9}},
	{AJNE, yjcond, Px, [23]uint8{0x75, 0x85}},
	{AJOC, yjcond, Px, [23]uint8{0x71, 0x81, 00}},
	{AJOS, yjcond, Px, [23]uint8{0x70, 0x80, 00}},
	{AJPC, yjcond, Px, [23]uint8{0x7b, 0x8b}},
	{AJPL, yjcond, Px, [23]uint8{0x79, 0x89}},
	{AJPS, yjcond, Px, [23]uint8{0x7a, 0x8a}},
	{ALAHF, ynone, Px, [23]uint8{0x9f}},
	{ALARL, yml_rl, Pm, [23]uint8{0x02}},
	{ALARW, yml_rl, Pq, [23]uint8{0x02}},
	{ALDMXCSR, ysvrs, Pm, [23]uint8{0xae, 02, 0xae, 02}},
	{ALEAL, ym_rl, Px, [23]uint8{0x8d}},
	{ALEAQ, ym_rl, Pw, [23]uint8{0x8d}},
	{ALEAVEL, ynone, P32, [23]uint8{0xc9}},
	{ALEAVEQ, ynone, Py, [23]uint8{0xc9}},
	{ALEAVEW, ynone, Pe, [23]uint8{0xc9}},
	{ALEAW, ym_rl, Pe, [23]uint8{0x8d}},
	{ALOCK, ynone, Px, [23]uint8{0xf0}},
	{ALODSB, ynone, Pb, [23]uint8{0xac}},
	{ALODSL, ynone, Px, [23]uint8{0xad}},
	{ALODSQ, ynone, Pw, [23]uint8{0xad}},
	{ALODSW, ynone, Pe, [23]uint8{0xad}},
	{ALONG, ybyte, Px, [23]uint8{4}},
	{ALOOP, yloop, Px, [23]uint8{0xe2}},
	{ALOOPEQ, yloop, Px, [23]uint8{0xe1}},
	{ALOOPNE, yloop, Px, [23]uint8{0xe0}},
	{ALSLL, yml_rl, Pm, [23]uint8{0x03}},
	{ALSLW, yml_rl, Pq, [23]uint8{0x03}},
	{AMASKMOVOU, yxr, Pe, [23]uint8{0xf7}},
	{AMASKMOVQ, ymr, Pm, [23]uint8{0xf7}},
	{AMAXPD, yxm, Pe, [23]uint8{0x5f}},
	{AMAXPS, yxm, Pm, [23]uint8{0x5f}},
	{AMAXSD, yxm, Pf2, [23]uint8{0x5f}},
	{AMAXSS, yxm, Pf3, [23]uint8{0x5f}},
	{AMINPD, yxm, Pe, [23]uint8{0x5d}},
	{AMINPS, yxm, Pm, [23]uint8{0x5d}},
	{AMINSD, yxm, Pf2, [23]uint8{0x5d}},
	{AMINSS, yxm, Pf3, [23]uint8{0x5d}},
	{AMOVAPD, yxmov, Pe, [23]uint8{0x28, 0x29}},
	{AMOVAPS, yxmov, Pm, [23]uint8{0x28, 0x29}},
	{AMOVB, ymovb, Pb, [23]uint8{0x88, 0x8a, 0xb0, 0xc6, 00}},
	{AMOVBLSX, ymb_rl, Pm, [23]uint8{0xbe}},
	{AMOVBLZX, ymb_rl, Pm, [23]uint8{0xb6}},
	{AMOVBQSX, ymb_rl, Pw, [23]uint8{0x0f, 0xbe}},
	{AMOVBQZX, ymb_rl, Pm, [23]uint8{0xb6}},
	{AMOVBWSX, ymb_rl, Pq, [23]uint8{0xbe}},
	{AMOVBWZX, ymb_rl, Pq, [23]uint8{0xb6}},
	{AMOVO, yxmov, Pe, [23]uint8{0x6f, 0x7f}},
	{AMOVOU, yxmov, Pf3, [23]uint8{0x6f, 0x7f}},
	{AMOVHLPS, yxr, Pm, [23]uint8{0x12}},
	{AMOVHPD, yxmov, Pe, [23]uint8{0x16, 0x17}},
	{AMOVHPS, yxmov, Pm, [23]uint8{0x16, 0x17}},
	{AMOVL, ymovl, Px, [23]uint8{0x89, 0x8b, 0x31, 0xb8, 0xc7, 00, 0x6e, 0x7e, Pe, 0x6e, Pe, 0x7e, 0}},
	{AMOVLHPS, yxr, Pm, [23]uint8{0x16}},
	{AMOVLPD, yxmov, Pe, [23]uint8{0x12, 0x13}},
	{AMOVLPS, yxmov, Pm, [23]uint8{0x12, 0x13}},
	{AMOVLQSX, yml_rl, Pw, [23]uint8{0x63}},
	{AMOVLQZX, yml_rl, Px, [23]uint8{0x8b}},
	{AMOVMSKPD, yxrrl, Pq, [23]uint8{0x50}},
	{AMOVMSKPS, yxrrl, Pm, [23]uint8{0x50}},
	{AMOVNTO, yxr_ml, Pe, [23]uint8{0xe7}},
	{AMOVNTPD, yxr_ml, Pe, [23]uint8{0x2b}},
	{AMOVNTPS, yxr_ml, Pm, [23]uint8{0x2b}},
	{AMOVNTQ, ymr_ml, Pm, [23]uint8{0xe7}},
	{AMOVQ, ymovq, Pw8, [23]uint8{0x6f, 0x7f, Pf2, 0xd6, Pf3, 0x7e, Pe, 0xd6, 0x89, 0x8b, 0x31, 0xc7, 00, 0xb8, 0xc7, 00, 0x6e, 0x7e, Pe, 0x6e, Pe, 0x7e, 0}},
	{AMOVQOZX, ymrxr, Pf3, [23]uint8{0xd6, 0x7e}},
	{AMOVSB, ynone, Pb, [23]uint8{0xa4}},
	{AMOVSD, yxmov, Pf2, [23]uint8{0x10, 0x11}},
	{AMOVSL, ynone, Px, [23]uint8{0xa5}},
	{AMOVSQ, ynone, Pw, [23]uint8{0xa5}},
	{AMOVSS, yxmov, Pf3, [23]uint8{0x10, 0x11}},
	{AMOVSW, ynone, Pe, [23]uint8{0xa5}},
	{AMOVUPD, yxmov, Pe, [23]uint8{0x10, 0x11}},
	{AMOVUPS, yxmov, Pm, [23]uint8{0x10, 0x11}},
	{AMOVW, ymovw, Pe, [23]uint8{0x89, 0x8b, 0x31, 0xb8, 0xc7, 00, 0}},
	{AMOVWLSX, yml_rl, Pm, [23]uint8{0xbf}},
	{AMOVWLZX, yml_rl, Pm, [23]uint8{0xb7}},
	{AMOVWQSX, yml_rl, Pw, [23]uint8{0x0f, 0xbf}},
	{AMOVWQZX, yml_rl, Pw, [23]uint8{0x0f, 0xb7}},
	{AMULB, ydivb, Pb, [23]uint8{0xf6, 04}},
	{AMULL, ydivl, Px, [23]uint8{0xf7, 04}},
	{AMULPD, yxm, Pe, [23]uint8{0x59}},
	{AMULPS, yxm, Ym, [23]uint8{0x59}},
	{AMULQ, ydivl, Pw, [23]uint8{0xf7, 04}},
	{AMULSD, yxm, Pf2, [23]uint8{0x59}},
	{AMULSS, yxm, Pf3, [23]uint8{0x59}},
	{AMULW, ydivl, Pe, [23]uint8{0xf7, 04}},
	{ANEGB, yscond, Pb, [23]uint8{0xf6, 03}},
	{ANEGL, yscond, Px, [23]uint8{0xf7, 03}},
	{ANEGQ, yscond, Pw, [23]uint8{0xf7, 03}},
	{ANEGW, yscond, Pe, [23]uint8{0xf7, 03}},
	{obj.ANOP, ynop, Px, [23]uint8{0, 0}},
	{ANOTB, yscond, Pb, [23]uint8{0xf6, 02}},
	{ANOTL, yscond, Px, [23]uint8{0xf7, 02}}, // TODO(rsc): yscond is wrong here.
	{ANOTQ, yscond, Pw, [23]uint8{0xf7, 02}},
	{ANOTW, yscond, Pe, [23]uint8{0xf7, 02}},
	{AORB, yxorb, Pb, [23]uint8{0x0c, 0x80, 01, 0x08, 0x0a}},
	{AORL, yxorl, Px, [23]uint8{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AORPD, yxm, Pq, [23]uint8{0x56}},
	{AORPS, yxm, Pm, [23]uint8{0x56}},
	{AORQ, yxorl, Pw, [23]uint8{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AORW, yxorl, Pe, [23]uint8{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AOUTB, yin, Pb, [23]uint8{0xe6, 0xee}},
	{AOUTL, yin, Px, [23]uint8{0xe7, 0xef}},
	{AOUTSB, ynone, Pb, [23]uint8{0x6e}},
	{AOUTSL, ynone, Px, [23]uint8{0x6f}},
	{AOUTSW, ynone, Pe, [23]uint8{0x6f}},
	{AOUTW, yin, Pe, [23]uint8{0xe7, 0xef}},
	{APACKSSLW, ymm, Py1, [23]uint8{0x6b, Pe, 0x6b}},
	{APACKSSWB, ymm, Py1, [23]uint8{0x63, Pe, 0x63}},
	{APACKUSWB, ymm, Py1, [23]uint8{0x67, Pe, 0x67}},
	{APADDB, ymm, Py1, [23]uint8{0xfc, Pe, 0xfc}},
	{APADDL, ymm, Py1, [23]uint8{0xfe, Pe, 0xfe}},
	{APADDQ, yxm, Pe, [23]uint8{0xd4}},
	{APADDSB, ymm, Py1, [23]uint8{0xec, Pe, 0xec}},
	{APADDSW, ymm, Py1, [23]uint8{0xed, Pe, 0xed}},
	{APADDUSB, ymm, Py1, [23]uint8{0xdc, Pe, 0xdc}},
	{APADDUSW, ymm, Py1, [23]uint8{0xdd, Pe, 0xdd}},
	{APADDW, ymm, Py1, [23]uint8{0xfd, Pe, 0xfd}},
	{APAND, ymm, Py1, [23]uint8{0xdb, Pe, 0xdb}},
	{APANDN, ymm, Py1, [23]uint8{0xdf, Pe, 0xdf}},
	{APAUSE, ynone, Px, [23]uint8{0xf3, 0x90}},
	{APAVGB, ymm, Py1, [23]uint8{0xe0, Pe, 0xe0}},
	{APAVGW, ymm, Py1, [23]uint8{0xe3, Pe, 0xe3}},
	{APCMPEQB, ymm, Py1, [23]uint8{0x74, Pe, 0x74}},
	{APCMPEQL, ymm, Py1, [23]uint8{0x76, Pe, 0x76}},
	{APCMPEQW, ymm, Py1, [23]uint8{0x75, Pe, 0x75}},
	{APCMPGTB, ymm, Py1, [23]uint8{0x64, Pe, 0x64}},
	{APCMPGTL, ymm, Py1, [23]uint8{0x66, Pe, 0x66}},
	{APCMPGTW, ymm, Py1, [23]uint8{0x65, Pe, 0x65}},
	{APEXTRW, yextrw, Pq, [23]uint8{0xc5, 00}},
	{APF2IL, ymfp, Px, [23]uint8{0x1d}},
	{APF2IW, ymfp, Px, [23]uint8{0x1c}},
	{API2FL, ymfp, Px, [23]uint8{0x0d}},
	{APFACC, ymfp, Px, [23]uint8{0xae}},
	{APFADD, ymfp, Px, [23]uint8{0x9e}},
	{APFCMPEQ, ymfp, Px, [23]uint8{0xb0}},
	{APFCMPGE, ymfp, Px, [23]uint8{0x90}},
	{APFCMPGT, ymfp, Px, [23]uint8{0xa0}},
	{APFMAX, ymfp, Px, [23]uint8{0xa4}},
	{APFMIN, ymfp, Px, [23]uint8{0x94}},
	{APFMUL, ymfp, Px, [23]uint8{0xb4}},
	{APFNACC, ymfp, Px, [23]uint8{0x8a}},
	{APFPNACC, ymfp, Px, [23]uint8{0x8e}},
	{APFRCP, ymfp, Px, [23]uint8{0x96}},
	{APFRCPIT1, ymfp, Px, [23]uint8{0xa6}},
	{APFRCPI2T, ymfp, Px, [23]uint8{0xb6}},
	{APFRSQIT1, ymfp, Px, [23]uint8{0xa7}},
	{APFRSQRT, ymfp, Px, [23]uint8{0x97}},
	{APFSUB, ymfp, Px, [23]uint8{0x9a}},
	{APFSUBR, ymfp, Px, [23]uint8{0xaa}},
	{APINSRW, yinsrw, Pq, [23]uint8{0xc4, 00}},
	{APINSRD, yinsr, Pq, [23]uint8{0x3a, 0x22, 00}},
	{APINSRQ, yinsr, Pq3, [23]uint8{0x3a, 0x22, 00}},
	{APMADDWL, ymm, Py1, [23]uint8{0xf5, Pe, 0xf5}},
	{APMAXSW, yxm, Pe, [23]uint8{0xee}},
	{APMAXUB, yxm, Pe, [23]uint8{0xde}},
	{APMINSW, yxm, Pe, [23]uint8{0xea}},
	{APMINUB, yxm, Pe, [23]uint8{0xda}},
	{APMOVMSKB, ymskb, Px, [23]uint8{Pe, 0xd7, 0xd7}},
	{APMULHRW, ymfp, Px, [23]uint8{0xb7}},
	{APMULHUW, ymm, Py1, [23]uint8{0xe4, Pe, 0xe4}},
	{APMULHW, ymm, Py1, [23]uint8{0xe5, Pe, 0xe5}},
	{APMULLW, ymm, Py1, [23]uint8{0xd5, Pe, 0xd5}},
	{APMULULQ, ymm, Py1, [23]uint8{0xf4, Pe, 0xf4}},
	{APOPAL, ynone, P32, [23]uint8{0x61}},
	{APOPAW, ynone, Pe, [23]uint8{0x61}},
	{APOPFL, ynone, P32, [23]uint8{0x9d}},
	{APOPFQ, ynone, Py, [23]uint8{0x9d}},
	{APOPFW, ynone, Pe, [23]uint8{0x9d}},
	{APOPL, ypopl, P32, [23]uint8{0x58, 0x8f, 00}},
	{APOPQ, ypopl, Py, [23]uint8{0x58, 0x8f, 00}},
	{APOPW, ypopl, Pe, [23]uint8{0x58, 0x8f, 00}},
	{APOR, ymm, Py1, [23]uint8{0xeb, Pe, 0xeb}},
	{APSADBW, yxm, Pq, [23]uint8{0xf6}},
	{APSHUFHW, yxshuf, Pf3, [23]uint8{0x70, 00}},
	{APSHUFL, yxshuf, Pq, [23]uint8{0x70, 00}},
	{APSHUFLW, yxshuf, Pf2, [23]uint8{0x70, 00}},
	{APSHUFW, ymshuf, Pm, [23]uint8{0x70, 00}},
	{APSHUFB, ymshufb, Pq, [23]uint8{0x38, 0x00}},
	{APSLLO, ypsdq, Pq, [23]uint8{0x73, 07}},
	{APSLLL, yps, Py3, [23]uint8{0xf2, 0x72, 06, Pe, 0xf2, Pe, 0x72, 06}},
	{APSLLQ, yps, Py3, [23]uint8{0xf3, 0x73, 06, Pe, 0xf3, Pe, 0x73, 06}},
	{APSLLW, yps, Py3, [23]uint8{0xf1, 0x71, 06, Pe, 0xf1, Pe, 0x71, 06}},
	{APSRAL, yps, Py3, [23]uint8{0xe2, 0x72, 04, Pe, 0xe2, Pe, 0x72, 04}},
	{APSRAW, yps, Py3, [23]uint8{0xe1, 0x71, 04, Pe, 0xe1, Pe, 0x71, 04}},
	{APSRLO, ypsdq, Pq, [23]uint8{0x73, 03}},
	{APSRLL, yps, Py3, [23]uint8{0xd2, 0x72, 02, Pe, 0xd2, Pe, 0x72, 02}},
	{APSRLQ, yps, Py3, [23]uint8{0xd3, 0x73, 02, Pe, 0xd3, Pe, 0x73, 02}},
	{APSRLW, yps, Py3, [23]uint8{0xd1, 0x71, 02, Pe, 0xe1, Pe, 0x71, 02}},
	{APSUBB, yxm, Pe, [23]uint8{0xf8}},
	{APSUBL, yxm, Pe, [23]uint8{0xfa}},
	{APSUBQ, yxm, Pe, [23]uint8{0xfb}},
	{APSUBSB, yxm, Pe, [23]uint8{0xe8}},
	{APSUBSW, yxm, Pe, [23]uint8{0xe9}},
	{APSUBUSB, yxm, Pe, [23]uint8{0xd8}},
	{APSUBUSW, yxm, Pe, [23]uint8{0xd9}},
	{APSUBW, yxm, Pe, [23]uint8{0xf9}},
	{APSWAPL, ymfp, Px, [23]uint8{0xbb}},
	{APUNPCKHBW, ymm, Py1, [23]uint8{0x68, Pe, 0x68}},
	{APUNPCKHLQ, ymm, Py1, [23]uint8{0x6a, Pe, 0x6a}},
	{APUNPCKHQDQ, yxm, Pe, [23]uint8{0x6d}},
	{APUNPCKHWL, ymm, Py1, [23]uint8{0x69, Pe, 0x69}},
	{APUNPCKLBW, ymm, Py1, [23]uint8{0x60, Pe, 0x60}},
	{APUNPCKLLQ, ymm, Py1, [23]uint8{0x62, Pe, 0x62}},
	{APUNPCKLQDQ, yxm, Pe, [23]uint8{0x6c}},
	{APUNPCKLWL, ymm, Py1, [23]uint8{0x61, Pe, 0x61}},
	{APUSHAL, ynone, P32, [23]uint8{0x60}},
	{APUSHAW, ynone, Pe, [23]uint8{0x60}},
	{APUSHFL, ynone, P32, [23]uint8{0x9c}},
	{APUSHFQ, ynone, Py, [23]uint8{0x9c}},
	{APUSHFW, ynone, Pe, [23]uint8{0x9c}},
	{APUSHL, ypushl, P32, [23]uint8{0x50, 0xff, 06, 0x6a, 0x68}},
	{APUSHQ, ypushl, Py, [23]uint8{0x50, 0xff, 06, 0x6a, 0x68}},
	{APUSHW, ypushl, Pe, [23]uint8{0x50, 0xff, 06, 0x6a, 0x68}},
	{APXOR, ymm, Py1, [23]uint8{0xef, Pe, 0xef}},
	{AQUAD, ybyte, Px, [23]uint8{8}},
	{ARCLB, yshb, Pb, [23]uint8{0xd0, 02, 0xc0, 02, 0xd2, 02}},
	{ARCLL, yshl, Px, [23]uint8{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCLQ, yshl, Pw, [23]uint8{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCLW, yshl, Pe, [23]uint8{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCPPS, yxm, Pm, [23]uint8{0x53}},
	{ARCPSS, yxm, Pf3, [23]uint8{0x53}},
	{ARCRB, yshb, Pb, [23]uint8{0xd0, 03, 0xc0, 03, 0xd2, 03}},
	{ARCRL, yshl, Px, [23]uint8{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{ARCRQ, yshl, Pw, [23]uint8{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{ARCRW, yshl, Pe, [23]uint8{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{AREP, ynone, Px, [23]uint8{0xf3}},
	{AREPN, ynone, Px, [23]uint8{0xf2}},
	{obj.ARET, ynone, Px, [23]uint8{0xc3}},
	{ARETFW, yret, Pe, [23]uint8{0xcb, 0xca}},
	{ARETFL, yret, Px, [23]uint8{0xcb, 0xca}},
	{ARETFQ, yret, Pw, [23]uint8{0xcb, 0xca}},
	{AROLB, yshb, Pb, [23]uint8{0xd0, 00, 0xc0, 00, 0xd2, 00}},
	{AROLL, yshl, Px, [23]uint8{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{AROLQ, yshl, Pw, [23]uint8{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{AROLW, yshl, Pe, [23]uint8{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{ARORB, yshb, Pb, [23]uint8{0xd0, 01, 0xc0, 01, 0xd2, 01}},
	{ARORL, yshl, Px, [23]uint8{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARORQ, yshl, Pw, [23]uint8{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARORW, yshl, Pe, [23]uint8{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARSQRTPS, yxm, Pm, [23]uint8{0x52}},
	{ARSQRTSS, yxm, Pf3, [23]uint8{0x52}},
	{ASAHF, ynone, Px1, [23]uint8{0x9e, 00, 0x86, 0xe0, 0x50, 0x9d}}, /* XCHGB AH,AL; PUSH AX; POPFL */
	{ASALB, yshb, Pb, [23]uint8{0xd0, 04, 0xc0, 04, 0xd2, 04}},
	{ASALL, yshl, Px, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASALQ, yshl, Pw, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASALW, yshl, Pe, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASARB, yshb, Pb, [23]uint8{0xd0, 07, 0xc0, 07, 0xd2, 07}},
	{ASARL, yshl, Px, [23]uint8{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASARQ, yshl, Pw, [23]uint8{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASARW, yshl, Pe, [23]uint8{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASBBB, yxorb, Pb, [23]uint8{0x1c, 0x80, 03, 0x18, 0x1a}},
	{ASBBL, yxorl, Px, [23]uint8{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASBBQ, yxorl, Pw, [23]uint8{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASBBW, yxorl, Pe, [23]uint8{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASCASB, ynone, Pb, [23]uint8{0xae}},
	{ASCASL, ynone, Px, [23]uint8{0xaf}},
	{ASCASQ, ynone, Pw, [23]uint8{0xaf}},
	{ASCASW, ynone, Pe, [23]uint8{0xaf}},
	{ASETCC, yscond, Pb, [23]uint8{0x0f, 0x93, 00}},
	{ASETCS, yscond, Pb, [23]uint8{0x0f, 0x92, 00}},
	{ASETEQ, yscond, Pb, [23]uint8{0x0f, 0x94, 00}},
	{ASETGE, yscond, Pb, [23]uint8{0x0f, 0x9d, 00}},
	{ASETGT, yscond, Pb, [23]uint8{0x0f, 0x9f, 00}},
	{ASETHI, yscond, Pb, [23]uint8{0x0f, 0x97, 00}},
	{ASETLE, yscond, Pb, [23]uint8{0x0f, 0x9e, 00}},
	{ASETLS, yscond, Pb, [23]uint8{0x0f, 0x96, 00}},
	{ASETLT, yscond, Pb, [23]uint8{0x0f, 0x9c, 00}},
	{ASETMI, yscond, Pb, [23]uint8{0x0f, 0x98, 00}},
	{ASETNE, yscond, Pb, [23]uint8{0x0f, 0x95, 00}},
	{ASETOC, yscond, Pb, [23]uint8{0x0f, 0x91, 00}},
	{ASETOS, yscond, Pb, [23]uint8{0x0f, 0x90, 00}},
	{ASETPC, yscond, Pb, [23]uint8{0x0f, 0x9b, 00}},
	{ASETPL, yscond, Pb, [23]uint8{0x0f, 0x99, 00}},
	{ASETPS, yscond, Pb, [23]uint8{0x0f, 0x9a, 00}},
	{ASHLB, yshb, Pb, [23]uint8{0xd0, 04, 0xc0, 04, 0xd2, 04}},
	{ASHLL, yshl, Px, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHLQ, yshl, Pw, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHLW, yshl, Pe, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHRB, yshb, Pb, [23]uint8{0xd0, 05, 0xc0, 05, 0xd2, 05}},
	{ASHRL, yshl, Px, [23]uint8{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHRQ, yshl, Pw, [23]uint8{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHRW, yshl, Pe, [23]uint8{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHUFPD, yxshuf, Pq, [23]uint8{0xc6, 00}},
	{ASHUFPS, yxshuf, Pm, [23]uint8{0xc6, 00}},
	{ASQRTPD, yxm, Pe, [23]uint8{0x51}},
	{ASQRTPS, yxm, Pm, [23]uint8{0x51}},
	{ASQRTSD, yxm, Pf2, [23]uint8{0x51}},
	{ASQRTSS, yxm, Pf3, [23]uint8{0x51}},
	{ASTC, ynone, Px, [23]uint8{0xf9}},
	{ASTD, ynone, Px, [23]uint8{0xfd}},
	{ASTI, ynone, Px, [23]uint8{0xfb}},
	{ASTMXCSR, ysvrs, Pm, [23]uint8{0xae, 03, 0xae, 03}},
	{ASTOSB, ynone, Pb, [23]uint8{0xaa}},
	{ASTOSL, ynone, Px, [23]uint8{0xab}},
	{ASTOSQ, ynone, Pw, [23]uint8{0xab}},
	{ASTOSW, ynone, Pe, [23]uint8{0xab}},
	{ASUBB, yxorb, Pb, [23]uint8{0x2c, 0x80, 05, 0x28, 0x2a}},
	{ASUBL, yaddl, Px, [23]uint8{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASUBPD, yxm, Pe, [23]uint8{0x5c}},
	{ASUBPS, yxm, Pm, [23]uint8{0x5c}},
	{ASUBQ, yaddl, Pw, [23]uint8{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASUBSD, yxm, Pf2, [23]uint8{0x5c}},
	{ASUBSS, yxm, Pf3, [23]uint8{0x5c}},
	{ASUBW, yaddl, Pe, [23]uint8{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASWAPGS, ynone, Pm, [23]uint8{0x01, 0xf8}},
	{ASYSCALL, ynone, Px, [23]uint8{0x0f, 0x05}}, /* fast syscall */
	{ATESTB, ytestb, Pb, [23]uint8{0xa8, 0xf6, 00, 0x84, 0x84}},
	{ATESTL, ytestl, Px, [23]uint8{0xa9, 0xf7, 00, 0x85, 0x85}},
	{ATESTQ, ytestl, Pw, [23]uint8{0xa9, 0xf7, 00, 0x85, 0x85}},
	{ATESTW, ytestl, Pe, [23]uint8{0xa9, 0xf7, 00, 0x85, 0x85}},
	{obj.ATEXT, ytext, Px, [23]uint8{}},
	{AUCOMISD, yxcmp, Pe, [23]uint8{0x2e}},
	{AUCOMISS, yxcmp, Pm, [23]uint8{0x2e}},
	{AUNPCKHPD, yxm, Pe, [23]uint8{0x15}},
	{AUNPCKHPS, yxm, Pm, [23]uint8{0x15}},
	{AUNPCKLPD, yxm, Pe, [23]uint8{0x14}},
	{AUNPCKLPS, yxm, Pm, [23]uint8{0x14}},
	{AVERR, ydivl, Pm, [23]uint8{0x00, 04}},
	{AVERW, ydivl, Pm, [23]uint8{0x00, 05}},
	{AWAIT, ynone, Px, [23]uint8{0x9b}},
	{AWORD, ybyte, Px, [23]uint8{2}},
	{AXCHGB, yml_mb, Pb, [23]uint8{0x86, 0x86}},
	{AXCHGL, yxchg, Px, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXCHGQ, yxchg, Pw, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXCHGW, yxchg, Pe, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXLAT, ynone, Px, [23]uint8{0xd7}},
	{AXORB, yxorb, Pb, [23]uint8{0x34, 0x80, 06, 0x30, 0x32}},
	{AXORL, yxorl, Px, [23]uint8{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AXORPD, yxm, Pe, [23]uint8{0x57}},
	{AXORPS, yxm, Pm, [23]uint8{0x57}},
	{AXORQ, yxorl, Pw, [23]uint8{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AXORW, yxorl, Pe, [23]uint8{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AFMOVB, yfmvx, Px, [23]uint8{0xdf, 04}},
	{AFMOVBP, yfmvp, Px, [23]uint8{0xdf, 06}},
	{AFMOVD, yfmvd, Px, [23]uint8{0xdd, 00, 0xdd, 02, 0xd9, 00, 0xdd, 02}},
	{AFMOVDP, yfmvdp, Px, [23]uint8{0xdd, 03, 0xdd, 03}},
	{AFMOVF, yfmvf, Px, [23]uint8{0xd9, 00, 0xd9, 02}},
	{AFMOVFP, yfmvp, Px, [23]uint8{0xd9, 03}},
	{AFMOVL, yfmvf, Px, [23]uint8{0xdb, 00, 0xdb, 02}},
	{AFMOVLP, yfmvp, Px, [23]uint8{0xdb, 03}},
	{AFMOVV, yfmvx, Px, [23]uint8{0xdf, 05}},
	{AFMOVVP, yfmvp, Px, [23]uint8{0xdf, 07}},
	{AFMOVW, yfmvf, Px, [23]uint8{0xdf, 00, 0xdf, 02}},
	{AFMOVWP, yfmvp, Px, [23]uint8{0xdf, 03}},
	{AFMOVX, yfmvx, Px, [23]uint8{0xdb, 05}},
	{AFMOVXP, yfmvp, Px, [23]uint8{0xdb, 07}},
	{AFCMOVCC, yfcmv, Px, [23]uint8{0xdb, 00}},
	{AFCMOVCS, yfcmv, Px, [23]uint8{0xda, 00}},
	{AFCMOVEQ, yfcmv, Px, [23]uint8{0xda, 01}},
	{AFCMOVHI, yfcmv, Px, [23]uint8{0xdb, 02}},
	{AFCMOVLS, yfcmv, Px, [23]uint8{0xda, 02}},
	{AFCMOVNE, yfcmv, Px, [23]uint8{0xdb, 01}},
	{AFCMOVNU, yfcmv, Px, [23]uint8{0xdb, 03}},
	{AFCMOVUN, yfcmv, Px, [23]uint8{0xda, 03}},
	{AFCOMB, nil, 0, [23]uint8{}},
	{AFCOMBP, nil, 0, [23]uint8{}},
	{AFCOMD, yfadd, Px, [23]uint8{0xdc, 02, 0xd8, 02, 0xdc, 02}},  /* botch */
	{AFCOMDP, yfadd, Px, [23]uint8{0xdc, 03, 0xd8, 03, 0xdc, 03}}, /* botch */
	{AFCOMDPP, ycompp, Px, [23]uint8{0xde, 03}},
	{AFCOMF, yfmvx, Px, [23]uint8{0xd8, 02}},
	{AFCOMFP, yfmvx, Px, [23]uint8{0xd8, 03}},
	{AFCOMI, yfmvx, Px, [23]uint8{0xdb, 06}},
	{AFCOMIP, yfmvx, Px, [23]uint8{0xdf, 06}},
	{AFCOML, yfmvx, Px, [23]uint8{0xda, 02}},
	{AFCOMLP, yfmvx, Px, [23]uint8{0xda, 03}},
	{AFCOMW, yfmvx, Px, [23]uint8{0xde, 02}},
	{AFCOMWP, yfmvx, Px, [23]uint8{0xde, 03}},
	{AFUCOM, ycompp, Px, [23]uint8{0xdd, 04}},
	{AFUCOMI, ycompp, Px, [23]uint8{0xdb, 05}},
	{AFUCOMIP, ycompp, Px, [23]uint8{0xdf, 05}},
	{AFUCOMP, ycompp, Px, [23]uint8{0xdd, 05}},
	{AFUCOMPP, ycompp, Px, [23]uint8{0xda, 13}},
	{AFADDDP, yfaddp, Px, [23]uint8{0xde, 00}},
	{AFADDW, yfmvx, Px, [23]uint8{0xde, 00}},
	{AFADDL, yfmvx, Px, [23]uint8{0xda, 00}},
	{AFADDF, yfmvx, Px, [23]uint8{0xd8, 00}},
	{AFADDD, yfadd, Px, [23]uint8{0xdc, 00, 0xd8, 00, 0xdc, 00}},
	{AFMULDP, yfaddp, Px, [23]uint8{0xde, 01}},
	{AFMULW, yfmvx, Px, [23]uint8{0xde, 01}},
	{AFMULL, yfmvx, Px, [23]uint8{0xda, 01}},
	{AFMULF, yfmvx, Px, [23]uint8{0xd8, 01}},
	{AFMULD, yfadd, Px, [23]uint8{0xdc, 01, 0xd8, 01, 0xdc, 01}},
	{AFSUBDP, yfaddp, Px, [23]uint8{0xde, 05}},
	{AFSUBW, yfmvx, Px, [23]uint8{0xde, 04}},
	{AFSUBL, yfmvx, Px, [23]uint8{0xda, 04}},
	{AFSUBF, yfmvx, Px, [23]uint8{0xd8, 04}},
	{AFSUBD, yfadd, Px, [23]uint8{0xdc, 04, 0xd8, 04, 0xdc, 05}},
	{AFSUBRDP, yfaddp, Px, [23]uint8{0xde, 04}},
	{AFSUBRW, yfmvx, Px, [23]uint8{0xde, 05}},
	{AFSUBRL, yfmvx, Px, [23]uint8{0xda, 05}},
	{AFSUBRF, yfmvx, Px, [23]uint8{0xd8, 05}},
	{AFSUBRD, yfadd, Px, [23]uint8{0xdc, 05, 0xd8, 05, 0xdc, 04}},
	{AFDIVDP, yfaddp, Px, [23]uint8{0xde, 07}},
	{AFDIVW, yfmvx, Px, [23]uint8{0xde, 06}},
	{AFDIVL, yfmvx, Px, [23]uint8{0xda, 06}},
	{AFDIVF, yfmvx, Px, [23]uint8{0xd8, 06}},
	{AFDIVD, yfadd, Px, [23]uint8{0xdc, 06, 0xd8, 06, 0xdc, 07}},
	{AFDIVRDP, yfaddp, Px, [23]uint8{0xde, 06}},
	{AFDIVRW, yfmvx, Px, [23]uint8{0xde, 07}},
	{AFDIVRL, yfmvx, Px, [23]uint8{0xda, 07}},
	{AFDIVRF, yfmvx, Px, [23]uint8{0xd8, 07}},
	{AFDIVRD, yfadd, Px, [23]uint8{0xdc, 07, 0xd8, 07, 0xdc, 06}},
	{AFXCHD, yfxch, Px, [23]uint8{0xd9, 01, 0xd9, 01}},
	{AFFREE, nil, 0, [23]uint8{}},
	{AFLDCW, ystcw, Px, [23]uint8{0xd9, 05, 0xd9, 05}},
	{AFLDENV, ystcw, Px, [23]uint8{0xd9, 04, 0xd9, 04}},
	{AFRSTOR, ysvrs, Px, [23]uint8{0xdd, 04, 0xdd, 04}},
	{AFSAVE, ysvrs, Px, [23]uint8{0xdd, 06, 0xdd, 06}},
	{AFSTCW, ystcw, Px, [23]uint8{0xd9, 07, 0xd9, 07}},
	{AFSTENV, ystcw, Px, [23]uint8{0xd9, 06, 0xd9, 06}},
	{AFSTSW, ystsw, Px, [23]uint8{0xdd, 07, 0xdf, 0xe0}},
	{AF2XM1, ynone, Px, [23]uint8{0xd9, 0xf0}},
	{AFABS, ynone, Px, [23]uint8{0xd9, 0xe1}},
	{AFCHS, ynone, Px, [23]uint8{0xd9, 0xe0}},
	{AFCLEX, ynone, Px, [23]uint8{0xdb, 0xe2}},
	{AFCOS, ynone, Px, [23]uint8{0xd9, 0xff}},
	{AFDECSTP, ynone, Px, [23]uint8{0xd9, 0xf6}},
	{AFINCSTP, ynone, Px, [23]uint8{0xd9, 0xf7}},
	{AFINIT, ynone, Px, [23]uint8{0xdb, 0xe3}},
	{AFLD1, ynone, Px, [23]uint8{0xd9, 0xe8}},
	{AFLDL2E, ynone, Px, [23]uint8{0xd9, 0xea}},
	{AFLDL2T, ynone, Px, [23]uint8{0xd9, 0xe9}},
	{AFLDLG2, ynone, Px, [23]uint8{0xd9, 0xec}},
	{AFLDLN2, ynone, Px, [23]uint8{0xd9, 0xed}},
	{AFLDPI, ynone, Px, [23]uint8{0xd9, 0xeb}},
	{AFLDZ, ynone, Px, [23]uint8{0xd9, 0xee}},
	{AFNOP, ynone, Px, [23]uint8{0xd9, 0xd0}},
	{AFPATAN, ynone, Px, [23]uint8{0xd9, 0xf3}},
	{AFPREM, ynone, Px, [23]uint8{0xd9, 0xf8}},
	{AFPREM1, ynone, Px, [23]uint8{0xd9, 0xf5}},
	{AFPTAN, ynone, Px, [23]uint8{0xd9, 0xf2}},
	{AFRNDINT, ynone, Px, [23]uint8{0xd9, 0xfc}},
	{AFSCALE, ynone, Px, [23]uint8{0xd9, 0xfd}},
	{AFSIN, ynone, Px, [23]uint8{0xd9, 0xfe}},
	{AFSINCOS, ynone, Px, [23]uint8{0xd9, 0xfb}},
	{AFSQRT, ynone, Px, [23]uint8{0xd9, 0xfa}},
	{AFTST, ynone, Px, [23]uint8{0xd9, 0xe4}},
	{AFXAM, ynone, Px, [23]uint8{0xd9, 0xe5}},
	{AFXTRACT, ynone, Px, [23]uint8{0xd9, 0xf4}},
	{AFYL2X, ynone, Px, [23]uint8{0xd9, 0xf1}},
	{AFYL2XP1, ynone, Px, [23]uint8{0xd9, 0xf9}},
	{ACMPXCHGB, yrb_mb, Pb, [23]uint8{0x0f, 0xb0}},
	{ACMPXCHGL, yrl_ml, Px, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHGW, yrl_ml, Pe, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHGQ, yrl_ml, Pw, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHG8B, yscond, Pm, [23]uint8{0xc7, 01}},
	{AINVD, ynone, Pm, [23]uint8{0x08}},
	{AINVLPG, ymbs, Pm, [23]uint8{0x01, 07}},
	{ALFENCE, ynone, Pm, [23]uint8{0xae, 0xe8}},
	{AMFENCE, ynone, Pm, [23]uint8{0xae, 0xf0}},
	{AMOVNTIL, yrl_ml, Pm, [23]uint8{0xc3}},
	{AMOVNTIQ, yrl_ml, Pw, [23]uint8{0x0f, 0xc3}},
	{ARDMSR, ynone, Pm, [23]uint8{0x32}},
	{ARDPMC, ynone, Pm, [23]uint8{0x33}},
	{ARDTSC, ynone, Pm, [23]uint8{0x31}},
	{ARSM, ynone, Pm, [23]uint8{0xaa}},
	{ASFENCE, ynone, Pm, [23]uint8{0xae, 0xf8}},
	{ASYSRET, ynone, Pm, [23]uint8{0x07}},
	{AWBINVD, ynone, Pm, [23]uint8{0x09}},
	{AWRMSR, ynone, Pm, [23]uint8{0x30}},
	{AXADDB, yrb_mb, Pb, [23]uint8{0x0f, 0xc0}},
	{AXADDL, yrl_ml, Px, [23]uint8{0x0f, 0xc1}},
	{AXADDQ, yrl_ml, Pw, [23]uint8{0x0f, 0xc1}},
	{AXADDW, yrl_ml, Pe, [23]uint8{0x0f, 0xc1}},
	{ACRC32B, ycrc32l, Px, [23]uint8{0xf2, 0x0f, 0x38, 0xf0, 0}},
	{ACRC32Q, ycrc32l, Pw, [23]uint8{0xf2, 0x0f, 0x38, 0xf1, 0}},
	{APREFETCHT0, yprefetch, Pm, [23]uint8{0x18, 01}},
	{APREFETCHT1, yprefetch, Pm, [23]uint8{0x18, 02}},
	{APREFETCHT2, yprefetch, Pm, [23]uint8{0x18, 03}},
	{APREFETCHNTA, yprefetch, Pm, [23]uint8{0x18, 00}},
	{AMOVQL, yrl_ml, Px, [23]uint8{0x89}},
	{obj.AUNDEF, ynone, Px, [23]uint8{0x0f, 0x0b}},
	{AAESENC, yaes, Pq, [23]uint8{0x38, 0xdc, 0}},
	{AAESENCLAST, yaes, Pq, [23]uint8{0x38, 0xdd, 0}},
	{AAESDEC, yaes, Pq, [23]uint8{0x38, 0xde, 0}},
	{AAESDECLAST, yaes, Pq, [23]uint8{0x38, 0xdf, 0}},
	{AAESIMC, yaes, Pq, [23]uint8{0x38, 0xdb, 0}},
	{AAESKEYGENASSIST, yaes2, Pq, [23]uint8{0x3a, 0xdf, 0}},
	{APSHUFD, yxshuf, Pq, [23]uint8{0x70, 0}},
	{APCLMULQDQ, yxshuf, Pq, [23]uint8{0x3a, 0x44, 0}},
	{obj.AUSEFIELD, ynop, Px, [23]uint8{0, 0}},
	{obj.ATYPE, nil, 0, [23]uint8{}},
	{obj.AFUNCDATA, yfuncdata, Px, [23]uint8{0, 0}},
	{obj.APCDATA, ypcdata, Px, [23]uint8{0, 0}},
	{obj.ACHECKNIL, nil, 0, [23]uint8{}},
	{obj.AVARDEF, nil, 0, [23]uint8{}},
	{obj.AVARKILL, nil, 0, [23]uint8{}},
	{obj.ADUFFCOPY, yduff, Px, [23]uint8{0xe8}},
	{obj.ADUFFZERO, yduff, Px, [23]uint8{0xe8}},
	{obj.AEND, nil, 0, [23]uint8{}},
	{0, nil, 0, [23]uint8{}},
}

var opindex [(ALAST + 1) & obj.AMask]*Optab

// single-instruction no-ops of various lengths.
// constructed by hand and disassembled with gdb to verify.
// see http://www.agner.org/optimize/optimizing_assembly.pdf for discussion.
var nop = [][16]uint8{
	{0x90},
	{0x66, 0x90},
	{0x0F, 0x1F, 0x00},
	{0x0F, 0x1F, 0x40, 0x00},
	{0x0F, 0x1F, 0x44, 0x00, 0x00},
	{0x66, 0x0F, 0x1F, 0x44, 0x00, 0x00},
	{0x0F, 0x1F, 0x80, 0x00, 0x00, 0x00, 0x00},
	{0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
	{0x66, 0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
}

// Native Client rejects the repeated 0x66 prefix.
// {0x66, 0x66, 0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
func fillnop(p []byte, n int) {
	var m int

	for n > 0 {
		m = n
		if m > len(nop) {
			m = len(nop)
		}
		copy(p[:m], nop[m-1][:m])
		p = p[m:]
		n -= m
	}
}
