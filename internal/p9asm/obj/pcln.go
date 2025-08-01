// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package obj

import (
	"fmt"
	"log"
)

type Pcdata struct {
	P []byte
}

func (d *Pcdata) addvarint(val uint32) {
	var v uint32
	for v = val; v >= 0x80; v >>= 7 {
		d.P = append(d.P, uint8(v|0x80))
	}
	d.P = append(d.P, uint8(v))
}

// funcpctab writes to dst a pc-value table mapping the code in func to the values
// returned by valfunc parameterized by arg. The invocation of valfunc to update the
// current value is, for each p,
//
//	val = valfunc(func, val, p, 0, arg);
//	record val as value at p->pc;
//	val = valfunc(func, val, p, 1, arg);
//
// where func is the function, val is the current value, p is the instruction being
// considered, and arg can be used to further parameterize valfunc.
func (dst *Pcdata) funcpctab(ctxt *Link, func_ *LSym, desc string, valfunc func(*Link, *LSym, int32, *Prog, int32, interface{}) int32, arg interface{}) {
	// To debug a specific function, uncomment second line and change name.
	dbg := 0

	//dbg = strcmp(func->name, "main.main") == 0;
	//dbg = strcmp(desc, "pctofile") == 0;

	ctxt.Debugpcln += int32(dbg)

	dst.P = dst.P[:0]

	if ctxt.Debugpcln != 0 {
		fmt.Fprintf(ctxt.Bso, "funcpctab %s [valfunc=%s]\n", func_.Name, desc)
	}

	val := int32(-1)
	oldval := val
	if func_.Text == nil {
		ctxt.Debugpcln -= int32(dbg)
		return
	}

	pc := func_.Text.Pc

	if ctxt.Debugpcln != 0 {
		fmt.Fprintf(ctxt.Bso, "%6x %6d %v\n", uint64(pc), val, func_.Text)
	}

	started := int32(0)
	var delta uint32
	for p := func_.Text; p != nil; p = p.Link {
		// Update val. If it's not changing, keep going.
		val = valfunc(ctxt, func_, val, p, 0, arg)

		if val == oldval && started != 0 {
			val = valfunc(ctxt, func_, val, p, 1, arg)
			if ctxt.Debugpcln != 0 {
				fmt.Fprintf(ctxt.Bso, "%6x %6s %v\n", uint64(int64(p.Pc)), "", p)
			}
			continue
		}

		// If the pc of the next instruction is the same as the
		// pc of this instruction, this instruction is not a real
		// instruction. Keep going, so that we only emit a delta
		// for a true instruction boundary in the program.
		if p.Link != nil && p.Link.Pc == p.Pc {
			val = valfunc(ctxt, func_, val, p, 1, arg)
			if ctxt.Debugpcln != 0 {
				fmt.Fprintf(ctxt.Bso, "%6x %6s %v\n", uint64(int64(p.Pc)), "", p)
			}
			continue
		}

		// The table is a sequence of (value, pc) pairs, where each
		// pair states that the given value is in effect from the current position
		// up to the given pc, which becomes the new current position.
		// To generate the table as we scan over the program instructions,
		// we emit a "(value" when pc == func->value, and then
		// each time we observe a change in value we emit ", pc) (value".
		// When the scan is over, we emit the closing ", pc)".
		//
		// The table is delta-encoded. The value deltas are signed and
		// transmitted in zig-zag form, where a complement bit is placed in bit 0,
		// and the pc deltas are unsigned. Both kinds of deltas are sent
		// as variable-length little-endian base-128 integers,
		// where the 0x80 bit indicates that the integer continues.

		if ctxt.Debugpcln != 0 {
			fmt.Fprintf(ctxt.Bso, "%6x %6d %v\n", uint64(int64(p.Pc)), val, p)
		}

		if started != 0 {
			dst.addvarint(uint32((p.Pc - pc) / int64(ctxt.Arch.Minlc)))
			pc = p.Pc
		}

		delta = uint32(val) - uint32(oldval)
		if delta>>31 != 0 {
			delta = 1 | ^(delta << 1)
		} else {
			delta <<= 1
		}
		dst.addvarint(delta)
		oldval = val
		started = 1
		val = valfunc(ctxt, func_, val, p, 1, arg)
	}

	if started != 0 {
		if ctxt.Debugpcln != 0 {
			fmt.Fprintf(ctxt.Bso, "%6x done\n", uint64(int64(func_.Text.Pc)+func_.Size))
		}
		dst.addvarint(uint32((func_.Value + func_.Size - pc) / int64(ctxt.Arch.Minlc)))
		dst.addvarint(0) // terminator
	}

	if ctxt.Debugpcln != 0 {
		fmt.Fprintf(ctxt.Bso, "wrote %d bytes to %p\n", len(dst.P), dst)
		for i := 0; i < len(dst.P); i++ {
			fmt.Fprintf(ctxt.Bso, " %02x", dst.P[i])
		}
		fmt.Fprintf(ctxt.Bso, "\n")
	}

	ctxt.Debugpcln -= int32(dbg)
}

// pctofileline computes either the file number (arg == 0)
// or the line number (arg == 1) to use at p.
// Because p->lineno applies to p, phase == 0 (before p)
// takes care of the update.
func (ctxt *Link) pctofileline(sym *LSym, oldval int32, p *Prog, phase int32, arg interface{}) int32 {
	if p.As == ATEXT || p.As == ANOP || p.As == AUSEFIELD || p.Lineno == 0 || phase == 1 {
		return oldval
	}
	var l int32
	var f *LSym
	ctxt.linkgetline(p.Lineno, &f, &l)
	if f == nil {
		//	print("getline failed for %s %v\n", ctxt->cursym->name, p);
		return oldval
	}

	if arg == nil {
		return l
	}
	pcln := arg.(*Pcln)

	if f == pcln.Lastfile {
		return int32(pcln.Lastindex)
	}

	var i int32
	for i = 0; i < int32(len(pcln.File)); i++ {
		file := pcln.File[i]
		if file == f {
			pcln.Lastfile = f
			pcln.Lastindex = int(i)
			return int32(i)
		}
	}
	pcln.File = append(pcln.File, f)
	pcln.Lastfile = f
	pcln.Lastindex = int(i)
	return i
}

// pctospadj computes the sp adjustment in effect.
// It is oldval plus any adjustment made by p itself.
// The adjustment by p takes effect only after p, so we
// apply the change during phase == 1.
func (ctxt *Link) pctospadj(sym *LSym, oldval int32, p *Prog, phase int32, arg interface{}) int32 {
	if oldval == -1 { // starting
		oldval = 0
	}
	if phase == 0 {
		return oldval
	}
	if oldval+p.Spadj < -10000 || oldval+p.Spadj > 1100000000 {
		ctxt.Diag("overflow in spadj: %d + %d = %d", oldval, p.Spadj, oldval+p.Spadj)
		log.Fatalf("bad code")
	}

	return oldval + p.Spadj
}

// pctopcdata computes the pcdata value in effect at p.
// A PCDATA instruction sets the value in effect at future
// non-PCDATA instructions.
// Since PCDATA instructions have no width in the final code,
// it does not matter which phase we use for the update.
func (*Link) pctopcdata(sym *LSym, oldval int32, p *Prog, phase int32, arg interface{}) int32 {
	if phase == 0 || p.As != APCDATA || p.From.Offset != int64(arg.(uint32)) {
		return oldval
	}
	if int64(int32(p.To.Offset)) != p.To.Offset {
		panic(fmt.Sprintf("overflow in PCDATA instruction: %v", p))
	}

	return int32(p.To.Offset)
}

func (cursym *LSym) linkpcln(ctxt *Link) {
	ctxt.Cursym = cursym

	pcln := new(Pcln)
	cursym.Pcln = pcln

	npcdata := 0
	nfuncdata := 0
	for p := cursym.Text; p != nil; p = p.Link {
		if p.As == APCDATA && p.From.Offset >= int64(npcdata) {
			npcdata = int(p.From.Offset + 1)
		}
		if p.As == AFUNCDATA && p.From.Offset >= int64(nfuncdata) {
			nfuncdata = int(p.From.Offset + 1)
		}
	}

	pcln.Pcdata = make([]Pcdata, npcdata)
	pcln.Pcdata = pcln.Pcdata[:npcdata]
	pcln.Funcdata = make([]*LSym, nfuncdata)
	pcln.Funcdataoff = make([]int64, nfuncdata)
	pcln.Funcdataoff = pcln.Funcdataoff[:nfuncdata]

	pcln.Pcsp.funcpctab(ctxt, cursym, "pctospadj", (*Link).pctospadj, nil)
	pcln.Pcfile.funcpctab(ctxt, cursym, "pctofile", (*Link).pctofileline, pcln)
	pcln.Pcline.funcpctab(ctxt, cursym, "pctoline", (*Link).pctofileline, nil)

	// tabulate which pc and func data we have.
	havepc := make([]uint32, (npcdata+31)/32)
	havefunc := make([]uint32, (nfuncdata+31)/32)
	for p := cursym.Text; p != nil; p = p.Link {
		if p.As == AFUNCDATA {
			if (havefunc[p.From.Offset/32]>>uint64(p.From.Offset%32))&1 != 0 {
				ctxt.Diag("multiple definitions for FUNCDATA $%d", p.From.Offset)
			}
			havefunc[p.From.Offset/32] |= 1 << uint64(p.From.Offset%32)
		}

		if p.As == APCDATA {
			havepc[p.From.Offset/32] |= 1 << uint64(p.From.Offset%32)
		}
	}

	// pcdata.
	for i := 0; i < npcdata; i++ {
		if (havepc[i/32]>>uint(i%32))&1 == 0 {
			continue
		}
		pcln.Pcdata[i].funcpctab(ctxt, cursym, "pctopcdata", (*Link).pctopcdata, interface{}(uint32(i)))
	}

	// funcdata
	if nfuncdata > 0 {
		var i int
		for p := cursym.Text; p != nil; p = p.Link {
			if p.As == AFUNCDATA {
				i = int(p.From.Offset)
				pcln.Funcdataoff[i] = p.To.Offset
				if p.To.Type != TYPE_CONST {
					// TODO: Dedup.
					//funcdata_bytes += p->to.sym->size;
					pcln.Funcdata[i] = p.To.Sym
				}
			}
		}
	}
}
