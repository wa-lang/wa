// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p9x86

import (
	"fmt"
	"log"
)

// 指令变长编码缓存
type AsmBuf struct {
	buf     [100]byte
	off     int
	rexflag int
}

// Put1 appends one byte to the end of the buffer.
func (ab *AsmBuf) Put1(x byte) {
	ab.buf[ab.off] = x
	ab.off++
}

// Put2 appends two bytes to the end of the buffer.
func (ab *AsmBuf) Put2(x, y byte) {
	ab.buf[ab.off+0] = x
	ab.buf[ab.off+1] = y
	ab.off += 2
}

// Put3 appends three bytes to the end of the buffer.
func (ab *AsmBuf) Put3(x, y, z byte) {
	ab.buf[ab.off+0] = x
	ab.buf[ab.off+1] = y
	ab.buf[ab.off+2] = z
	ab.off += 3
}

// PutInt16 writes v into the buffer using little-endian encoding.
func (ab *AsmBuf) PutInt16(v int16) {
	ab.buf[ab.off+0] = byte(v)
	ab.buf[ab.off+1] = byte(v >> 8)
	ab.off += 2
}

// PutInt32 writes v into the buffer using little-endian encoding.
func (ab *AsmBuf) PutInt32(v int32) {
	ab.buf[ab.off+0] = byte(v)
	ab.buf[ab.off+1] = byte(v >> 8)
	ab.buf[ab.off+2] = byte(v >> 16)
	ab.buf[ab.off+3] = byte(v >> 24)
	ab.off += 4
}

// PutInt64 writes v into the buffer using little-endian encoding.
func (ab *AsmBuf) PutInt64(v int64) {
	ab.buf[ab.off+0] = byte(v)
	ab.buf[ab.off+1] = byte(v >> 8)
	ab.buf[ab.off+2] = byte(v >> 16)
	ab.buf[ab.off+3] = byte(v >> 24)
	ab.buf[ab.off+4] = byte(v >> 32)
	ab.buf[ab.off+5] = byte(v >> 40)
	ab.buf[ab.off+6] = byte(v >> 48)
	ab.buf[ab.off+7] = byte(v >> 56)
	ab.off += 8
}

// PutOpBytesLit writes zero terminated sequence of bytes from op,
// starting at specified offset (e.g. z counter value).
// Trailing 0 is not written.
//
// Intended to be used for literal Z cases.
// Literal Z cases usually have "Zlit" in their name (Zlit, Zlitr_m, Zlitm_r).
func (ab *AsmBuf) PutOpBytesLit(offset int, op *opBytes) {
	for int(op[offset]) != 0 {
		ab.Put1(byte(op[offset]))
		offset++
	}
}

func (ab *AsmBuf) Asmins(p *Prog) {
	ab.asmins(p)
}

// Insert inserts b at offset i.
func (ab *AsmBuf) Insert(i int, b byte) {
	ab.off++
	copy(ab.buf[i+1:ab.off], ab.buf[i:ab.off-1])
	ab.buf[i] = b
}

// Last returns the byte at the end of the buffer.
func (ab *AsmBuf) Last() byte { return ab.buf[ab.off-1] }

// Len returns the length of the buffer.
func (ab *AsmBuf) Len() int { return ab.off }

// Bytes returns the contents of the buffer.
func (ab *AsmBuf) Bytes() []byte { return ab.buf[:ab.off] }

// Reset empties the buffer.
func (ab *AsmBuf) Reset() { ab.off = 0 }

// At returns the byte at offset i.
func (ab *AsmBuf) At(i int) byte { return ab.buf[i] }

// asmidx emits SIB byte.
func (ab *AsmBuf) asmidx(scale int, index int, base int) {
	var i int

	// X/Y index register is used in VSIB.
	switch index {
	default:
		goto bad

	case REG_NONE:
		i = 4 << 3
		goto bas

	case REG_R8,
		REG_R9,
		REG_R10,
		REG_R11,
		REG_R12,
		REG_R13,
		REG_R14,
		REG_R15,
		REG_X8,
		REG_X9,
		REG_X10,
		REG_X11,
		REG_X12,
		REG_X13,
		REG_X14,
		REG_X15,
		REG_X16,
		REG_X17,
		REG_X18,
		REG_X19,
		REG_X20,
		REG_X21,
		REG_X22,
		REG_X23,
		REG_X24,
		REG_X25,
		REG_X26,
		REG_X27,
		REG_X28,
		REG_X29,
		REG_X30,
		REG_X31,
		REG_Y8,
		REG_Y9,
		REG_Y10,
		REG_Y11,
		REG_Y12,
		REG_Y13,
		REG_Y14,
		REG_Y15,
		REG_Y16,
		REG_Y17,
		REG_Y18,
		REG_Y19,
		REG_Y20,
		REG_Y21,
		REG_Y22,
		REG_Y23,
		REG_Y24,
		REG_Y25,
		REG_Y26,
		REG_Y27,
		REG_Y28,
		REG_Y29,
		REG_Y30,
		REG_Y31,
		REG_Z8,
		REG_Z9,
		REG_Z10,
		REG_Z11,
		REG_Z12,
		REG_Z13,
		REG_Z14,
		REG_Z15,
		REG_Z16,
		REG_Z17,
		REG_Z18,
		REG_Z19,
		REG_Z20,
		REG_Z21,
		REG_Z22,
		REG_Z23,
		REG_Z24,
		REG_Z25,
		REG_Z26,
		REG_Z27,
		REG_Z28,
		REG_Z29,
		REG_Z30,
		REG_Z31:

		fallthrough

	case REG_AX,
		REG_CX,
		REG_DX,
		REG_BX,
		REG_BP,
		REG_SI,
		REG_DI,
		REG_X0,
		REG_X1,
		REG_X2,
		REG_X3,
		REG_X4,
		REG_X5,
		REG_X6,
		REG_X7,
		REG_Y0,
		REG_Y1,
		REG_Y2,
		REG_Y3,
		REG_Y4,
		REG_Y5,
		REG_Y6,
		REG_Y7,
		REG_Z0,
		REG_Z1,
		REG_Z2,
		REG_Z3,
		REG_Z4,
		REG_Z5,
		REG_Z6,
		REG_Z7:
		i = reg[index] << 3
	}

	switch scale {
	default:
		goto bad

	case 1:
		break

	case 2:
		i |= 1 << 6

	case 4:
		i |= 2 << 6

	case 8:
		i |= 3 << 6
	}

bas:
	switch base {
	default:
		goto bad

	case REG_NONE: // must be mod=00
		i |= 5

	case REG_R8,
		REG_R9,
		REG_R10,
		REG_R11,
		REG_R12,
		REG_R13,
		REG_R14,
		REG_R15:

		fallthrough

	case REG_AX,
		REG_CX,
		REG_DX,
		REG_BX,
		REG_SP,
		REG_BP,
		REG_SI,
		REG_DI:
		i |= reg[base]
	}

	ab.Put1(byte(i))
	return

bad:
	panic(fmt.Errorf("asmidx: bad address %d/%d/%d", scale, index, base))
}

func (ab *AsmBuf) asmins(p *Prog) {
	ab.Reset()

	ab.rexflag = 0
	mark := ab.Len()
	ab.doasm(p)
	if ab.rexflag != 0 {
		// as befits the whole approach of the architecture,
		// the rex prefix must appear before the first opcode byte
		// (and thus after any 66/67/f2/f3/26/2e/3e prefix bytes, but
		// before the 0f opcode escape!), or it might be ignored.
		// note that the handbook often misleadingly shows 66/f2/f3 in `opcode'.
		n := ab.Len()
		var np int
		for np = mark; np < n; np++ {
			c := ab.At(np)
			if c != 0xf2 && c != 0xf3 && (c < 0x64 || c > 0x67) && c != 0x2e && c != 0x3e && c != 0x26 {
				break
			}
		}
		ab.Insert(np, byte(0x40|ab.rexflag))
	}
}

func (ab *AsmBuf) relput4(a *Addr) {
	ab.PutInt32(int32(a.Offset))
}

func (ab *AsmBuf) asmandsz(p *Prog, a *Addr, r int, rex int) {
	var base int

	rex &= 0x40 | Rxr
	if a.Offset != int64(int32(a.Offset)) {
		// The rules are slightly different for 386 and AMD64,
		// mostly for historical reasons. We may unify them later,
		// but it must be discussed beforehand.
		//
		// For 64bit mode only LEAL is allowed to overflow.
		// It's how https://golang.org/cl/59630 made it.
		// crypto/sha1/sha1block_amd64.s depends on this feature.
		//
		// For 32bit mode rules are more permissive.
		// If offset fits uint32, it's permitted.
		// This is allowed for assembly that wants to use 32-bit hex
		// constants, e.g. LEAL 0x99999999(AX), AX.
		overflowOK := (p.As == ALEAL)
		if !overflowOK {
			panic(fmt.Errorf("offset too large in %v", p))
		}
	}
	v := int32(a.Offset)

	switch a.Type {
	case TYPE_ADDR:
		panic(fmt.Errorf("unexpected TYPE_ADDR with NAME_NONE"))

	case TYPE_REG:
		const regFirst = REG_AL
		const regLast = REG_Z31
		if a.Reg < regFirst || regLast < a.Reg {
			goto bad
		}
		if v != 0 {
			goto bad
		}
		ab.Put1(byte(3<<6 | reg[a.Reg]<<0 | r<<3))
		ab.rexflag |= regrex[a.Reg]&(0x40|Rxb) | rex
		return
	}

	if a.Type != TYPE_MEM {
		goto bad
	}

	if a.Index != REG_NONE {
		base := int(a.Reg)

		ab.rexflag |= regrex[int(a.Index)]&Rxx | regrex[base]&Rxb | rex
		if base == REG_NONE {
			ab.Put1(byte(0<<6 | 4<<0 | r<<3))
			ab.asmidx(int(a.Scale), int(a.Index), base)
			goto putrelv
		}

		if v == 0 && base != REG_BP && base != REG_R13 {
			ab.Put1(byte(0<<6 | 4<<0 | r<<3))
			ab.asmidx(int(a.Scale), int(a.Index), base)
			return
		}

		if disp8, ok := toDisp8(v); ok {
			ab.Put1(byte(1<<6 | 4<<0 | r<<3))
			ab.asmidx(int(a.Scale), int(a.Index), base)
			ab.Put1(disp8)
			return
		}

		ab.Put1(byte(2<<6 | 4<<0 | r<<3))
		ab.asmidx(int(a.Scale), int(a.Index), base)
		goto putrelv
	}

	base = int(a.Reg)

	ab.rexflag |= regrex[base]&Rxb | rex
	if base == REG_NONE || (REG_CS <= base && base <= REG_GS) {

		// temporary
		ab.Put2(
			byte(0<<6|4<<0|r<<3), // sib present
			0<<6|4<<3|5<<0,       // DS:d32
		)
		goto putrelv
	}

	if base == REG_SP || base == REG_R12 {
		if v == 0 {
			ab.Put1(byte(0<<6 | reg[base]<<0 | r<<3))
			ab.asmidx(int(a.Scale), REG_NONE, base)
			return
		}

		if disp8, ok := toDisp8(v); ok {
			ab.Put1(byte(1<<6 | reg[base]<<0 | r<<3))
			ab.asmidx(int(a.Scale), REG_NONE, base)
			ab.Put1(disp8)
			return
		}

		ab.Put1(byte(2<<6 | reg[base]<<0 | r<<3))
		ab.asmidx(int(a.Scale), REG_NONE, base)
		goto putrelv
	}

	if REG_AX <= base && base <= REG_R15 {

		if v == 0 && base != REG_BP && base != REG_R13 {
			ab.Put1(byte(0<<6 | reg[base]<<0 | r<<3))
			return
		}

		if disp8, ok := toDisp8(v); ok {
			ab.Put2(byte(1<<6|reg[base]<<0|r<<3), disp8)
			return
		}

		ab.Put1(byte(2<<6 | reg[base]<<0 | r<<3))
		goto putrelv
	}

	goto bad

putrelv:

	ab.PutInt32(v)
	return

bad:
	panic(fmt.Errorf("asmand: bad address %+v", a))
}

func (ab *AsmBuf) asmand(p *Prog, a *Addr, ra *Addr) {
	ab.asmandsz(p, a, reg[ra.Reg], regrex[ra.Reg])
}

func (ab *AsmBuf) asmando(p *Prog, a *Addr, o int) {
	ab.asmandsz(p, a, o, 0)
}

func (ab *AsmBuf) mediaop(o *Optab, op int, osize int, z int) int {
	switch op {
	case Pm, Pe, Pf2, Pf3:
		if osize != 1 {
			if op != Pm {
				ab.Put1(byte(op))
			}
			ab.Put1(Pm)
			z++
			op = int(o.op[z])
			break
		}
		fallthrough

	default:
		if ab.Len() == 0 || ab.Last() != Pm {
			ab.Put1(Pm)
		}
	}

	ab.Put1(byte(op))
	return z
}

func (ab *AsmBuf) doasm(p *Prog) {
	o := opindex[p.As]

	if o == nil {
		panic(fmt.Errorf("asmins: missing op %v", p))
	}

	if pre := prefixof(&p.From); pre != 0 {
		ab.Put1(byte(pre))
	}
	if pre := prefixof(&p.To); pre != 0 {
		ab.Put1(byte(pre))
	}

	if p.Ft == 0 {
		p.Ft = uint8(oclass(&p.From))
	}
	if p.Tt == 0 {
		p.Tt = uint8(oclass(&p.To))
	}

	ft := int(p.Ft) * Ymax
	var f3t int
	tt := int(p.Tt) * Ymax

	xo := 0
	if o.op[0] == 0x0f {
		xo = 1
	}

	z := 0
	var a *Addr
	var l int
	var op int
	var v int64

	args := make([]int, 0, argListMax)
	if ft != Ynone*Ymax {
		args = append(args, ft)
	}
	for i := range p.RestArgs {
		args = append(args, oclass(&p.RestArgs[i])*Ymax)
	}
	if tt != Ynone*Ymax {
		args = append(args, tt)
	}

	for _, yt := range o.ytab {
		// ytab matching is purely args-based,
		// but AVX512 suffixes like "Z" or "RU_SAE" will
		// add EVEX-only filter that will reject non-EVEX matches.
		//
		// Consider "VADDPD.BCST 2032(DX), X0, X0".
		// Without this rule, operands will lead to VEX-encoded form
		// and produce "c5b15813" encoding.
		if !yt.match(args) {
			// "xo" is always zero for VEX/EVEX encoded insts.
			z += int(yt.zoffset) + xo
		} else {

			switch o.prefix {
			case Px1: // first option valid only in 32-bit mode
				if z == 0 {
					z += int(yt.zoffset) + xo
					continue
				}
			case Pq: // 16 bit escape and opcode escape
				ab.Put2(Pe, Pm)

			case Pq3: // 16 bit escape and opcode escape + REX.W
				ab.rexflag |= Pw
				ab.Put2(Pe, Pm)

			case Pq4: // 66 0F 38
				ab.Put3(0x66, 0x0F, 0x38)

			case Pq4w: // 66 0F 38 + REX.W
				ab.rexflag |= Pw
				ab.Put3(0x66, 0x0F, 0x38)

			case Pq5: // F3 0F 38
				ab.Put3(0xF3, 0x0F, 0x38)

			case Pq5w: //  F3 0F 38 + REX.W
				ab.rexflag |= Pw
				ab.Put3(0xF3, 0x0F, 0x38)

			case Pf2, // xmm opcode escape
				Pf3:
				ab.Put2(o.prefix, Pm)

			case Pef3:
				ab.Put3(Pe, Pf3, Pm)

			case Pfw: // xmm opcode escape + REX.W
				ab.rexflag |= Pw
				ab.Put2(Pf3, Pm)

			case Pm: // opcode escape
				ab.Put1(Pm)

			case Pe: // 16 bit escape
				ab.Put1(Pe)

			case Pw: // 64-bit escape
				ab.rexflag |= Pw

			case Pw8: // 64-bit escape if z >= 8
				if z >= 8 {
					ab.rexflag |= Pw
				}

			case Pb: // botch
				// NOTE(rsc): This is probably safe to do always,
				// but when enabled it chooses different encodings
				// than the old cmd/internal/obj/i386 code did,
				// which breaks our "same bits out" checks.
				// In particular, CMPB AX, $0 encodes as 80 f8 00
				// in the original obj/i386, and it would encode
				// (using a valid, shorter form) as 3c 00 if we enabled
				// the call to bytereg here.
				bytereg(&p.From, &p.Ft)
				bytereg(&p.To, &p.Tt)

			case P32: // 32 bit but illegal if 64-bit mode
				panic(fmt.Errorf("asmins: illegal in 64-bit mode: %v", p))

			case Py: // 64-bit only, no prefix

			case Py1: // 64-bit only if z < 1, no prefix

			case Py3: // 64-bit only if z < 3, no prefix
			}

			if z >= len(o.op) {
				log.Fatalf("asmins bad table %v", p)
			}
			op = int(o.op[z])
			if op == 0x0f {
				ab.Put1(byte(op))
				z++
				op = int(o.op[z])
			}

			switch yt.zcase {
			default:
				panic(fmt.Errorf("asmins: unknown z %d %v", yt.zcase, p))

			case Zpseudo:
				break

			case Zlit:
				ab.PutOpBytesLit(z, &o.op)

			case Zlitr_m:
				ab.PutOpBytesLit(z, &o.op)
				ab.asmand(p, &p.To, &p.From)

			case Zlitm_r:
				ab.PutOpBytesLit(z, &o.op)
				ab.asmand(p, &p.From, &p.To)

			case Zlit_m_r:
				ab.PutOpBytesLit(z, &o.op)
				ab.asmand(p, p.GetFrom3(), &p.To)

			case Zmb_r:
				bytereg(&p.From, &p.Ft)
				fallthrough

			case Zm_r:
				ab.Put1(byte(op))
				ab.asmand(p, &p.From, &p.To)

			case Z_m_r:
				ab.Put1(byte(op))
				ab.asmand(p, p.GetFrom3(), &p.To)

			case Zm2_r:
				ab.Put2(byte(op), o.op[z+1])
				ab.asmand(p, &p.From, &p.To)

			case Zm_r_xm:
				ab.mediaop(o, op, int(yt.zoffset), z)
				ab.asmand(p, &p.From, &p.To)

			case Zm_r_xm_nr:
				ab.rexflag = 0
				ab.mediaop(o, op, int(yt.zoffset), z)
				ab.asmand(p, &p.From, &p.To)

			case Zm_r_i_xm:
				ab.mediaop(o, op, int(yt.zoffset), z)
				ab.asmand(p, &p.From, p.GetFrom3())
				ab.Put1(byte(p.To.Offset))

			case Zibm_r, Zibr_m:
				ab.PutOpBytesLit(z, &o.op)
				if yt.zcase == Zibr_m {
					ab.asmand(p, &p.To, p.GetFrom3())
				} else {
					ab.asmand(p, p.GetFrom3(), &p.To)
				}
				switch {
				default:
					ab.Put1(byte(p.From.Offset))
				case yt.args[0] == Yi32 && o.prefix == Pe:
					ab.PutInt16(int16(p.From.Offset))
				case yt.args[0] == Yi32:
					ab.PutInt32(int32(p.From.Offset))
				}

			case Zaut_r:
				ab.Put1(0x8d) // leal
				if p.From.Type != TYPE_ADDR {
					panic(fmt.Errorf("asmins: Zaut sb type ADDR"))
				}
				p.From.Type = TYPE_MEM
				ab.asmand(p, &p.From, &p.To)
				p.From.Type = TYPE_ADDR

			case Zm_o:
				ab.Put1(byte(op))
				ab.asmando(p, &p.From, int(o.op[z+1]))

			case Zr_m:
				ab.Put1(byte(op))
				ab.asmand(p, &p.To, &p.From)

			case Zr_m_xm:
				ab.mediaop(o, op, int(yt.zoffset), z)
				ab.asmand(p, &p.To, &p.From)

			case Zr_m_xm_nr:
				ab.rexflag = 0
				ab.mediaop(o, op, int(yt.zoffset), z)
				ab.asmand(p, &p.To, &p.From)

			case Zo_m:
				ab.Put1(byte(op))
				ab.asmando(p, &p.To, int(o.op[z+1]))

			case Zcallindreg:
				fallthrough

			case Zo_m64:
				ab.Put1(byte(op))
				ab.asmandsz(p, &p.To, int(o.op[z+1]), 0)

			case Zm_ibo:
				ab.Put1(byte(op))
				ab.asmando(p, &p.From, int(o.op[z+1]))
				ab.Put1(byte(p.To.Offset))

			case Zibo_m:
				ab.Put1(byte(op))
				ab.asmando(p, &p.To, int(o.op[z+1]))
				ab.Put1(byte(p.From.Offset))

			case Zibo_m_xm:
				z = ab.mediaop(o, op, int(yt.zoffset), z)
				ab.asmando(p, &p.To, int(o.op[z+1]))
				ab.Put1(byte(p.From.Offset))

			case Z_ib, Zib_:
				if yt.zcase == Zib_ {
					a = &p.From
				} else {
					a = &p.To
				}
				ab.Put1(byte(op))
				if p.As == AXABORT {
					ab.Put1(o.op[z+1])
				}
				ab.Put1(byte(a.Offset))

			case Zib_rp:
				ab.rexflag |= regrex[p.To.Reg] & (Rxb | 0x40)
				ab.Put2(byte(op+reg[p.To.Reg]), byte(p.From.Offset))

			case Zil_rp:
				ab.rexflag |= regrex[p.To.Reg] & Rxb
				ab.Put1(byte(op + reg[p.To.Reg]))
				if o.prefix == Pe {
					v = p.From.Offset
					ab.PutInt16(int16(v))
				} else {
					ab.relput4(&p.From)
				}

			case Zo_iw:
				ab.Put1(byte(op))
				if p.From.Type != TYPE_NONE {
					v = p.From.Offset
					ab.PutInt16(int16(v))
				}

			case Ziq_rp:
				v = p.From.Offset
				l = int(v >> 32)
				if l == 0 {
					ab.rexflag &^= (0x40 | Rxw)

					ab.rexflag |= regrex[p.To.Reg] & Rxb
					ab.Put1(byte(0xb8 + reg[p.To.Reg]))

					ab.PutInt32(int32(v))
				} else if l == -1 && uint64(v)&(uint64(1)<<31) != 0 { // sign extend
					ab.Put1(0xc7)
					ab.asmando(p, &p.To, 0)

					ab.PutInt32(int32(v)) // need all 8
				} else {
					ab.rexflag |= regrex[p.To.Reg] & Rxb
					ab.Put1(byte(op + reg[p.To.Reg]))

					ab.PutInt64(v)
				}

			case Zib_rr:
				ab.Put1(byte(op))
				ab.asmand(p, &p.To, &p.To)
				ab.Put1(byte(p.From.Offset))

			case Z_il, Zil_:
				if yt.zcase == Zil_ {
					a = &p.From
				} else {
					a = &p.To
				}
				ab.Put1(byte(op))
				if o.prefix == Pe {
					v = a.Offset
					ab.PutInt16(int16(v))
				} else {
					ab.relput4(a)
				}

			case Zm_ilo, Zilo_m:
				ab.Put1(byte(op))
				if yt.zcase == Zilo_m {
					a = &p.From
					ab.asmando(p, &p.To, int(o.op[z+1]))
				} else {
					a = &p.To
					ab.asmando(p, &p.From, int(o.op[z+1]))
				}

				if o.prefix == Pe {
					v = a.Offset
					ab.PutInt16(int16(v))
				} else {
					ab.relput4(a)
				}

			case Zil_rr:
				ab.Put1(byte(op))
				ab.asmand(p, &p.To, &p.To)
				if o.prefix == Pe {
					v = p.From.Offset
					ab.PutInt16(int16(v))
				} else {
					ab.relput4(&p.From)
				}

			case Z_rp:
				ab.rexflag |= regrex[p.To.Reg] & (Rxb | 0x40)
				ab.Put1(byte(op + reg[p.To.Reg]))

			case Zrp_:
				ab.rexflag |= regrex[p.From.Reg] & (Rxb | 0x40)
				ab.Put1(byte(op + reg[p.From.Reg]))

			case Zcallcon, Zjmpcon:
				if yt.zcase == Zcallcon {
					ab.Put1(byte(op))
				} else {
					ab.Put1(o.op[z+1])
				}
				ab.PutInt32(0)

			case Zcallind:
				ab.Put2(byte(op), o.op[z+1])
				ab.PutInt32(0)

			case Zcall, Zcallduff:
				panic(fmt.Errorf("call without target"))

			// TODO: jump across functions needs reloc
			case Zbr, Zjmp, Zloop:
				if p.As == AXBEGIN {
					ab.Put1(byte(op))
				}

				// Assumes q is in this function.
				// TODO: Check in input, preserve in brchain.

				// Fill in backward jump now.
				panic(fmt.Errorf("jmp/branch/loop without target"))

			case Zbyte:
				v = p.From.Offset

				ab.Put1(byte(v))
				if op > 1 {
					ab.Put1(byte(v >> 8))
					if op > 2 {
						ab.PutInt16(int16(v >> 16))
						if op > 4 {
							ab.PutInt32(int32(v >> 32))
						}
					}
				}
			}

			return
		}
	}
	f3t = Ynone * Ymax
	if p.GetFrom3() != nil {
		f3t = oclass(p.GetFrom3()) * Ymax
	}
	for mo := ymovtab; mo[0].as != 0; mo = mo[1:] {
		var t []byte
		if p.As == mo[0].as {
			if ycover[ft+int(mo[0].ft)] != 0 && ycover[f3t+int(mo[0].f3t)] != 0 && ycover[tt+int(mo[0].tt)] != 0 {
				t = mo[0].op[:]
				switch mo[0].code {
				default:
					panic(fmt.Errorf("asmins: unknown mov %d %v", mo[0].code, p))

				case movLit:
					for z = 0; t[z] != 0; z++ {
						ab.Put1(t[z])
					}

				case movRegMem:
					ab.Put1(t[0])
					ab.asmando(p, &p.To, int(t[1]))

				case movMemReg:
					ab.Put1(t[0])
					ab.asmando(p, &p.From, int(t[1]))

				case movRegMem2op: // r,m - 2op
					ab.Put2(t[0], t[1])
					ab.asmando(p, &p.To, int(t[2]))
					ab.rexflag |= regrex[p.From.Reg] & (Rxr | 0x40)

				case movMemReg2op:
					ab.Put2(t[0], t[1])
					ab.asmando(p, &p.From, int(t[2]))
					ab.rexflag |= regrex[p.To.Reg] & (Rxr | 0x40)

				case movFullPtr:
					if t[0] != 0 {
						ab.Put1(t[0])
					}
					switch p.To.Index {
					default:
						goto bad

					case REG_DS:
						ab.Put1(0xc5)

					case REG_SS:
						ab.Put2(0x0f, 0xb2)

					case REG_ES:
						ab.Put1(0xc4)

					case REG_FS:
						ab.Put2(0x0f, 0xb4)

					case REG_GS:
						ab.Put2(0x0f, 0xb5)
					}

					ab.asmand(p, &p.From, &p.To)

				case movDoubleShift:
					switch t[0] {
					case Pw:
						ab.rexflag |= Pw
						t = t[1:]
					case Pe:
						ab.Put1(Pe)
						t = t[1:]
					}

					switch p.From.Type {
					default:
						goto bad

					case TYPE_CONST:
						ab.Put2(0x0f, t[0])
						ab.asmandsz(p, &p.To, reg[p.GetFrom3().Reg], regrex[p.GetFrom3().Reg])
						ab.Put1(byte(p.From.Offset))

					case TYPE_REG:
						switch p.From.Reg {
						default:
							goto bad

						case REG_CL, REG_CX:
							ab.Put2(0x0f, t[1])
							ab.asmandsz(p, &p.To, reg[p.GetFrom3().Reg], regrex[p.GetFrom3().Reg])
						}
					}

				// NOTE: The systems listed here are the ones that use the "TLS initial exec" model,
				// where you load the TLS base register into a register and then index off that
				// register to access the actual TLS variables. Systems that allow direct TLS access
				// are handled in prefixof above and should not be listed here.
				case movTLSReg:
					if p.As != AMOVQ {
						panic(fmt.Errorf("invalid load of TLS: %v", p))
					}

					log.Fatalf("unknown TLS base location for linux/freebsd without -shared")

					// Note that this is not generating the same insn as the other cases.
					//     MOV TLS, R_to
					// becomes
					//     movq g@gottpoff(%rip), R_to
					// which is encoded as
					//     movq 0(%rip), R_to
					// and a R_TLS_IE reloc. This all assumes the only tls variable we access
					// is g, which we can't check here, but will when we assemble the second
					// instruction.
					ab.rexflag = Pw | (regrex[p.To.Reg] & Rxr)

					ab.Put2(0x8B, byte(0x05|(reg[p.To.Reg]<<3)))
					ab.PutInt32(0)

				}
				return
			}
		}
	}
	goto bad

bad:

	panic(fmt.Errorf("invalid instruction: %v", p))
}

// toDisp8 tries to convert disp to proper 8-bit displacement value.
func toDisp8(disp int32) (disp8 byte, ok bool) {
	return byte(disp), disp >= -128 && disp < 128
}

func instinit() {
	if ycover[0] != 0 {
		// Already initialized; stop now.
		// This happens in the cmd/asm tests,
		// each of which re-initializes the arch.
		return
	}

	for i := 1; optab[i].as != 0; i++ {
		c := optab[i].as
		if opindex[c] != nil {
			panic(fmt.Errorf("phase error in optab: %d (%v)", i, c))
		}
		opindex[c] = &optab[i]
	}

	for i := 0; i < Ymax; i++ {
		ycover[i*Ymax+i] = 1
	}

	ycover[Yi0*Ymax+Yu2] = 1
	ycover[Yi1*Ymax+Yu2] = 1

	ycover[Yi0*Ymax+Yi8] = 1
	ycover[Yi1*Ymax+Yi8] = 1
	ycover[Yu2*Ymax+Yi8] = 1
	ycover[Yu7*Ymax+Yi8] = 1

	ycover[Yi0*Ymax+Yu7] = 1
	ycover[Yi1*Ymax+Yu7] = 1
	ycover[Yu2*Ymax+Yu7] = 1

	ycover[Yi0*Ymax+Yu8] = 1
	ycover[Yi1*Ymax+Yu8] = 1
	ycover[Yu2*Ymax+Yu8] = 1
	ycover[Yu7*Ymax+Yu8] = 1

	ycover[Yi0*Ymax+Ys32] = 1
	ycover[Yi1*Ymax+Ys32] = 1
	ycover[Yu2*Ymax+Ys32] = 1
	ycover[Yu7*Ymax+Ys32] = 1
	ycover[Yu8*Ymax+Ys32] = 1
	ycover[Yi8*Ymax+Ys32] = 1

	ycover[Yi0*Ymax+Yi32] = 1
	ycover[Yi1*Ymax+Yi32] = 1
	ycover[Yu2*Ymax+Yi32] = 1
	ycover[Yu7*Ymax+Yi32] = 1
	ycover[Yu8*Ymax+Yi32] = 1
	ycover[Yi8*Ymax+Yi32] = 1
	ycover[Ys32*Ymax+Yi32] = 1

	ycover[Yi0*Ymax+Yi64] = 1
	ycover[Yi1*Ymax+Yi64] = 1
	ycover[Yu7*Ymax+Yi64] = 1
	ycover[Yu2*Ymax+Yi64] = 1
	ycover[Yu8*Ymax+Yi64] = 1
	ycover[Yi8*Ymax+Yi64] = 1
	ycover[Ys32*Ymax+Yi64] = 1
	ycover[Yi32*Ymax+Yi64] = 1

	ycover[Yal*Ymax+Yrb] = 1
	ycover[Ycl*Ymax+Yrb] = 1
	ycover[Yax*Ymax+Yrb] = 1
	ycover[Ycx*Ymax+Yrb] = 1
	ycover[Yrx*Ymax+Yrb] = 1
	ycover[Yrl*Ymax+Yrb] = 1 // but not Yrl32

	ycover[Ycl*Ymax+Ycx] = 1

	ycover[Yax*Ymax+Yrx] = 1
	ycover[Ycx*Ymax+Yrx] = 1

	ycover[Yax*Ymax+Yrl] = 1
	ycover[Ycx*Ymax+Yrl] = 1
	ycover[Yrx*Ymax+Yrl] = 1
	ycover[Yrl32*Ymax+Yrl] = 1

	ycover[Yf0*Ymax+Yrf] = 1

	ycover[Yal*Ymax+Ymb] = 1
	ycover[Ycl*Ymax+Ymb] = 1
	ycover[Yax*Ymax+Ymb] = 1
	ycover[Ycx*Ymax+Ymb] = 1
	ycover[Yrx*Ymax+Ymb] = 1
	ycover[Yrb*Ymax+Ymb] = 1
	ycover[Yrl*Ymax+Ymb] = 1 // but not Yrl32
	ycover[Ym*Ymax+Ymb] = 1

	ycover[Yax*Ymax+Yml] = 1
	ycover[Ycx*Ymax+Yml] = 1
	ycover[Yrx*Ymax+Yml] = 1
	ycover[Yrl*Ymax+Yml] = 1
	ycover[Yrl32*Ymax+Yml] = 1
	ycover[Ym*Ymax+Yml] = 1

	ycover[Yax*Ymax+Ymm] = 1
	ycover[Ycx*Ymax+Ymm] = 1
	ycover[Yrx*Ymax+Ymm] = 1
	ycover[Yrl*Ymax+Ymm] = 1
	ycover[Yrl32*Ymax+Ymm] = 1
	ycover[Ym*Ymax+Ymm] = 1
	ycover[Ymr*Ymax+Ymm] = 1

	ycover[Yxr0*Ymax+Yxr] = 1

	ycover[Ym*Ymax+Yxm] = 1
	ycover[Yxr0*Ymax+Yxm] = 1
	ycover[Yxr*Ymax+Yxm] = 1

	ycover[Ym*Ymax+Yym] = 1
	ycover[Yyr*Ymax+Yym] = 1

	ycover[Yxr0*Ymax+YxrEvex] = 1
	ycover[Yxr*Ymax+YxrEvex] = 1

	ycover[Ym*Ymax+YxmEvex] = 1
	ycover[Yxr0*Ymax+YxmEvex] = 1
	ycover[Yxr*Ymax+YxmEvex] = 1
	ycover[YxrEvex*Ymax+YxmEvex] = 1

	ycover[Yyr*Ymax+YyrEvex] = 1

	ycover[Ym*Ymax+YymEvex] = 1
	ycover[Yyr*Ymax+YymEvex] = 1
	ycover[YyrEvex*Ymax+YymEvex] = 1

	ycover[Ym*Ymax+Yzm] = 1
	ycover[Yzr*Ymax+Yzm] = 1

	ycover[Yk0*Ymax+Yk] = 1
	ycover[Yknot0*Ymax+Yk] = 1

	ycover[Yk0*Ymax+Ykm] = 1
	ycover[Yknot0*Ymax+Ykm] = 1
	ycover[Yk*Ymax+Ykm] = 1
	ycover[Ym*Ymax+Ykm] = 1

	ycover[Yxvm*Ymax+YxvmEvex] = 1

	ycover[Yyvm*Ymax+YyvmEvex] = 1

	for i := 0; i < MAXREG; i++ {
		reg[i] = -1
		if i >= REG_AL && i <= REG_R15B {
			reg[i] = (i - REG_AL) & 7
			if i >= REG_SPB && i <= REG_DIB {
				regrex[i] = 0x40
			}
			if i >= REG_R8B && i <= REG_R15B {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}

		if i >= REG_AH && i <= REG_BH {
			reg[i] = 4 + ((i - REG_AH) & 7)
		}
		if i >= REG_AX && i <= REG_R15 {
			reg[i] = (i - REG_AX) & 7
			if i >= REG_R8 {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}

		if i >= REG_F0 && i <= REG_F0+7 {
			reg[i] = (i - REG_F0) & 7
		}
		if i >= REG_M0 && i <= REG_M0+7 {
			reg[i] = (i - REG_M0) & 7
		}
		if i >= REG_K0 && i <= REG_K0+7 {
			reg[i] = (i - REG_K0) & 7
		}
		if i >= REG_X0 && i <= REG_X0+15 {
			reg[i] = (i - REG_X0) & 7
			if i >= REG_X0+8 {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}
		if i >= REG_X16 && i <= REG_X16+15 {
			reg[i] = (i - REG_X16) & 7
			if i >= REG_X16+8 {
				regrex[i] = Rxr | Rxx | Rxb | RxrEvex
			} else {
				regrex[i] = RxrEvex
			}
		}
		if i >= REG_Y0 && i <= REG_Y0+15 {
			reg[i] = (i - REG_Y0) & 7
			if i >= REG_Y0+8 {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}
		if i >= REG_Y16 && i <= REG_Y16+15 {
			reg[i] = (i - REG_Y16) & 7
			if i >= REG_Y16+8 {
				regrex[i] = Rxr | Rxx | Rxb | RxrEvex
			} else {
				regrex[i] = RxrEvex
			}
		}
		if i >= REG_Z0 && i <= REG_Z0+15 {
			reg[i] = (i - REG_Z0) & 7
			if i > REG_Z0+7 {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}
		if i >= REG_Z16 && i <= REG_Z16+15 {
			reg[i] = (i - REG_Z16) & 7
			if i >= REG_Z16+8 {
				regrex[i] = Rxr | Rxx | Rxb | RxrEvex
			} else {
				regrex[i] = RxrEvex
			}
		}

		if i >= REG_CR+8 && i <= REG_CR+15 {
			regrex[i] = Rxr
		}
	}
}

func prefixof(a *Addr) int {
	if a.Reg < REG_CS && a.Index < REG_CS { // fast path
		return 0
	}
	if a.Type == TYPE_MEM {
		switch a.Reg {
		case REG_CS:
			return 0x2e

		case REG_DS:
			return 0x3e

		case REG_ES:
			return 0x26

		case REG_FS:
			return 0x64

		case REG_GS:
			return 0x65
		}
	}

	switch a.Index {
	case REG_CS:
		return 0x2e

	case REG_DS:
		return 0x3e

	case REG_ES:
		return 0x26

	case REG_FS:
		return 0x64

	case REG_GS:
		return 0x65
	}

	return 0
}

// oclassVMem returns V-mem (vector memory with VSIB) operand class.
// For addr that is not V-mem returns (Yxxx, false).
func oclassVMem(addr *Addr) (int, bool) {
	switch addr.Index {
	case REG_X0 + 0,
		REG_X0 + 1,
		REG_X0 + 2,
		REG_X0 + 3,
		REG_X0 + 4,
		REG_X0 + 5,
		REG_X0 + 6,
		REG_X0 + 7:
		return Yxvm, true
	case REG_X8 + 0,
		REG_X8 + 1,
		REG_X8 + 2,
		REG_X8 + 3,
		REG_X8 + 4,
		REG_X8 + 5,
		REG_X8 + 6,
		REG_X8 + 7:

		return Yxvm, true
	case REG_X16 + 0,
		REG_X16 + 1,
		REG_X16 + 2,
		REG_X16 + 3,
		REG_X16 + 4,
		REG_X16 + 5,
		REG_X16 + 6,
		REG_X16 + 7,
		REG_X16 + 8,
		REG_X16 + 9,
		REG_X16 + 10,
		REG_X16 + 11,
		REG_X16 + 12,
		REG_X16 + 13,
		REG_X16 + 14,
		REG_X16 + 15:

		return YxvmEvex, true

	case REG_Y0 + 0,
		REG_Y0 + 1,
		REG_Y0 + 2,
		REG_Y0 + 3,
		REG_Y0 + 4,
		REG_Y0 + 5,
		REG_Y0 + 6,
		REG_Y0 + 7:
		return Yyvm, true
	case REG_Y8 + 0,
		REG_Y8 + 1,
		REG_Y8 + 2,
		REG_Y8 + 3,
		REG_Y8 + 4,
		REG_Y8 + 5,
		REG_Y8 + 6,
		REG_Y8 + 7:

		return Yyvm, true
	case REG_Y16 + 0,
		REG_Y16 + 1,
		REG_Y16 + 2,
		REG_Y16 + 3,
		REG_Y16 + 4,
		REG_Y16 + 5,
		REG_Y16 + 6,
		REG_Y16 + 7,
		REG_Y16 + 8,
		REG_Y16 + 9,
		REG_Y16 + 10,
		REG_Y16 + 11,
		REG_Y16 + 12,
		REG_Y16 + 13,
		REG_Y16 + 14,
		REG_Y16 + 15:

		return YyvmEvex, true

	case REG_Z0 + 0,
		REG_Z0 + 1,
		REG_Z0 + 2,
		REG_Z0 + 3,
		REG_Z0 + 4,
		REG_Z0 + 5,
		REG_Z0 + 6,
		REG_Z0 + 7:
		return Yzvm, true
	case REG_Z8 + 0,
		REG_Z8 + 1,
		REG_Z8 + 2,
		REG_Z8 + 3,
		REG_Z8 + 4,
		REG_Z8 + 5,
		REG_Z8 + 6,
		REG_Z8 + 7,
		REG_Z8 + 8,
		REG_Z8 + 9,
		REG_Z8 + 10,
		REG_Z8 + 11,
		REG_Z8 + 12,
		REG_Z8 + 13,
		REG_Z8 + 14,
		REG_Z8 + 15,
		REG_Z8 + 16,
		REG_Z8 + 17,
		REG_Z8 + 18,
		REG_Z8 + 19,
		REG_Z8 + 20,
		REG_Z8 + 21,
		REG_Z8 + 22,
		REG_Z8 + 23:

		return Yzvm, true
	}

	return Yxxx, false
}

func oclass(a *Addr) int {
	switch a.Type {
	case TYPE_NONE:
		return Ynone

	case TYPE_BRANCH:
		return Ybr

	case TYPE_MEM:
		// Pseudo registers have negative index, but SP is
		// not pseudo on x86, hence REG_SP check is not redundant.
		if a.Index == REG_SP || a.Index < 0 {
			// Can't use FP/SB/PC/SP as the index register.
			return Yxxx
		}

		if vmem, ok := oclassVMem(a); ok {
			return vmem
		}

		return Ym

	case TYPE_ADDR:

		fallthrough

	case TYPE_CONST:

		v := a.Offset

		switch {
		case v == 0:
			return Yi0
		case v == 1:
			return Yi1
		case v >= 0 && v <= 3:
			return Yu2
		case v >= 0 && v <= 127:
			return Yu7
		case v >= 0 && v <= 255:
			return Yu8
		case v >= -128 && v <= 127:
			return Yi8
		}

		l := int32(v)
		if int64(l) == v {
			return Ys32 // can sign extend
		}
		if v>>32 == 0 {
			return Yi32 // unsigned
		}
		return Yi64
	}

	if a.Type != TYPE_REG {
		panic(fmt.Errorf("unexpected addr1: type=%d %+v", a.Type, a))
	}

	switch a.Reg {
	case REG_AL:
		return Yal

	case REG_AX:
		return Yax

		/*
			case REG_SPB:
		*/
	case REG_BPB,
		REG_SIB,
		REG_DIB,
		REG_R8B,
		REG_R9B,
		REG_R10B,
		REG_R11B,
		REG_R12B,
		REG_R13B,
		REG_R14B,
		REG_R15B:

		fallthrough

	case REG_DL,
		REG_BL,
		REG_AH,
		REG_CH,
		REG_DH,
		REG_BH:
		return Yrb

	case REG_CL:
		return Ycl

	case REG_CX:
		return Ycx

	case REG_DX, REG_BX:
		return Yrx

	case REG_R8, // not really Yrl
		REG_R9,
		REG_R10,
		REG_R11,
		REG_R12,
		REG_R13,
		REG_R14,
		REG_R15:

		fallthrough

	case REG_SP, REG_BP, REG_SI, REG_DI:
		return Yrl

	case REG_F0 + 0:
		return Yf0

	case REG_F0 + 1,
		REG_F0 + 2,
		REG_F0 + 3,
		REG_F0 + 4,
		REG_F0 + 5,
		REG_F0 + 6,
		REG_F0 + 7:
		return Yrf

	case REG_M0 + 0,
		REG_M0 + 1,
		REG_M0 + 2,
		REG_M0 + 3,
		REG_M0 + 4,
		REG_M0 + 5,
		REG_M0 + 6,
		REG_M0 + 7:
		return Ymr

	case REG_X0:
		return Yxr0

	case REG_X0 + 1,
		REG_X0 + 2,
		REG_X0 + 3,
		REG_X0 + 4,
		REG_X0 + 5,
		REG_X0 + 6,
		REG_X0 + 7,
		REG_X0 + 8,
		REG_X0 + 9,
		REG_X0 + 10,
		REG_X0 + 11,
		REG_X0 + 12,
		REG_X0 + 13,
		REG_X0 + 14,
		REG_X0 + 15:
		return Yxr

	case REG_X0 + 16,
		REG_X0 + 17,
		REG_X0 + 18,
		REG_X0 + 19,
		REG_X0 + 20,
		REG_X0 + 21,
		REG_X0 + 22,
		REG_X0 + 23,
		REG_X0 + 24,
		REG_X0 + 25,
		REG_X0 + 26,
		REG_X0 + 27,
		REG_X0 + 28,
		REG_X0 + 29,
		REG_X0 + 30,
		REG_X0 + 31:
		return YxrEvex

	case REG_Y0 + 0,
		REG_Y0 + 1,
		REG_Y0 + 2,
		REG_Y0 + 3,
		REG_Y0 + 4,
		REG_Y0 + 5,
		REG_Y0 + 6,
		REG_Y0 + 7,
		REG_Y0 + 8,
		REG_Y0 + 9,
		REG_Y0 + 10,
		REG_Y0 + 11,
		REG_Y0 + 12,
		REG_Y0 + 13,
		REG_Y0 + 14,
		REG_Y0 + 15:
		return Yyr

	case REG_Y0 + 16,
		REG_Y0 + 17,
		REG_Y0 + 18,
		REG_Y0 + 19,
		REG_Y0 + 20,
		REG_Y0 + 21,
		REG_Y0 + 22,
		REG_Y0 + 23,
		REG_Y0 + 24,
		REG_Y0 + 25,
		REG_Y0 + 26,
		REG_Y0 + 27,
		REG_Y0 + 28,
		REG_Y0 + 29,
		REG_Y0 + 30,
		REG_Y0 + 31:
		return YyrEvex

	case REG_Z0 + 0,
		REG_Z0 + 1,
		REG_Z0 + 2,
		REG_Z0 + 3,
		REG_Z0 + 4,
		REG_Z0 + 5,
		REG_Z0 + 6,
		REG_Z0 + 7:
		return Yzr

	case REG_Z0 + 8,
		REG_Z0 + 9,
		REG_Z0 + 10,
		REG_Z0 + 11,
		REG_Z0 + 12,
		REG_Z0 + 13,
		REG_Z0 + 14,
		REG_Z0 + 15,
		REG_Z0 + 16,
		REG_Z0 + 17,
		REG_Z0 + 18,
		REG_Z0 + 19,
		REG_Z0 + 20,
		REG_Z0 + 21,
		REG_Z0 + 22,
		REG_Z0 + 23,
		REG_Z0 + 24,
		REG_Z0 + 25,
		REG_Z0 + 26,
		REG_Z0 + 27,
		REG_Z0 + 28,
		REG_Z0 + 29,
		REG_Z0 + 30,
		REG_Z0 + 31:
		return Yzr

	case REG_K0:
		return Yk0

	case REG_K0 + 1,
		REG_K0 + 2,
		REG_K0 + 3,
		REG_K0 + 4,
		REG_K0 + 5,
		REG_K0 + 6,
		REG_K0 + 7:
		return Yknot0

	case REG_CS:
		return Ycs
	case REG_SS:
		return Yss
	case REG_DS:
		return Yds
	case REG_ES:
		return Yes
	case REG_FS:
		return Yfs
	case REG_GS:
		return Ygs

	case REG_GDTR:
		return Ygdtr
	case REG_IDTR:
		return Yidtr
	case REG_LDTR:
		return Yldtr
	case REG_MSW:
		return Ymsw
	case REG_TASK:
		return Ytask

	case REG_CR + 0:
		return Ycr0
	case REG_CR + 1:
		return Ycr1
	case REG_CR + 2:
		return Ycr2
	case REG_CR + 3:
		return Ycr3
	case REG_CR + 4:
		return Ycr4
	case REG_CR + 5:
		return Ycr5
	case REG_CR + 6:
		return Ycr6
	case REG_CR + 7:
		return Ycr7
	case REG_CR + 8:
		return Ycr8

	case REG_DR + 0:
		return Ydr0
	case REG_DR + 1:
		return Ydr1
	case REG_DR + 2:
		return Ydr2
	case REG_DR + 3:
		return Ydr3
	case REG_DR + 4:
		return Ydr4
	case REG_DR + 5:
		return Ydr5
	case REG_DR + 6:
		return Ydr6
	case REG_DR + 7:
		return Ydr7

	case REG_TR + 0:
		return Ytr0
	case REG_TR + 1:
		return Ytr1
	case REG_TR + 2:
		return Ytr2
	case REG_TR + 3:
		return Ytr3
	case REG_TR + 4:
		return Ytr4
	case REG_TR + 5:
		return Ytr5
	case REG_TR + 6:
		return Ytr6
	case REG_TR + 7:
		return Ytr7
	}

	return Yxxx
}

func bytereg(a *Addr, t *uint8) {
	if a.Type == TYPE_REG && a.Index == REG_NONE && (REG_AX <= a.Reg && a.Reg <= REG_R15) {
		a.Reg += REG_AL - REG_AX
		*t = 0
	}
}
