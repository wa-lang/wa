// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package obj

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var start time.Time

func Cputime() float64 {
	if start.IsZero() {
		start = time.Now()
	}
	return time.Since(start).Seconds()
}

func envOr(key, value string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return value
}

func Getwaroot() string {
	return envOr("WAROOT", "")
}

func Getwaarch() string {
	return envOr("WAARCH", runtime.GOARCH)
}

func Getwaos() string {
	return envOr("WAOS", runtime.GOOS)
}

func (p *Prog) Line() string {
	return p.Ctxt.LineHist.LineString(int(p.Lineno))
}

func (p *Prog) String() string {
	if p.Ctxt == nil {
		return "<Prog without ctxt>"
	}

	sc := CConv(p.Scond)

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%.5d (%v)\t%v%s", p.Pc, p.Line(), Aconv(int(p.As)), sc)
	sep := "\t"
	if p.From.Type != TYPE_NONE {
		fmt.Fprintf(&buf, "%s%v", sep, Dconv(p, &p.From))
		sep = ", "
	}
	if p.Reg != REG_NONE {
		// Should not happen but might as well show it if it does.
		fmt.Fprintf(&buf, "%s%v", sep, Rconv(int(p.Reg)))
		sep = ", "
	}
	if p.From3Type() != TYPE_NONE {
		if p.From3.Type == TYPE_CONST && (p.As == ADATA || p.As == ATEXT || p.As == AGLOBL) {
			// Special case - omit $.
			fmt.Fprintf(&buf, "%s%d", sep, p.From3.Offset)
		} else {
			fmt.Fprintf(&buf, "%s%v", sep, Dconv(p, p.From3))
		}
		sep = ", "
	}
	if p.To.Type != TYPE_NONE {
		fmt.Fprintf(&buf, "%s%v", sep, Dconv(p, &p.To))
	}
	if p.RegTo2 != REG_NONE {
		fmt.Fprintf(&buf, "%s%v", sep, Rconv(int(p.RegTo2)))
	}
	return buf.String()
}

func (ctxt *Link) NewProg() *Prog {
	p := new(Prog) // should be the only call to this; all others should use ctxt.NewProg
	p.Ctxt = ctxt
	return p
}

func (ctxt *Link) Line(n int) string {
	return ctxt.LineHist.LineString(n)
}

func Getcallerpc(interface{}) uintptr {
	return 1
}

func (ctxt *Link) Dconv(a *Addr) string {
	return Dconv(nil, a)
}

func Dconv(p *Prog, a *Addr) string {
	var str string

	switch a.Type {
	default:
		str = fmt.Sprintf("type=%d", a.Type)

	case TYPE_NONE:
		str = ""
		if a.Name != NAME_NONE || a.Reg != 0 || a.Sym != nil {
			str = fmt.Sprintf("%v(%v)(NONE)", Mconv(a), Rconv(int(a.Reg)))
		}

	case TYPE_REG:
		// TODO(chai2010): This special case is for x86 instructions like
		//	PINSRQ	CX,$1,X6
		// where the $1 is included in the p->to Addr.
		// Move into a new field.
		if a.Offset != 0 {
			str = fmt.Sprintf("$%d,%v", a.Offset, Rconv(int(a.Reg)))
			break
		}

		str = Rconv(int(a.Reg))
		if a.Name != TYPE_NONE || a.Sym != nil {
			str = fmt.Sprintf("%v(%v)(REG)", Mconv(a), Rconv(int(a.Reg)))
		}

	case TYPE_BRANCH:
		if a.Sym != nil {
			str = fmt.Sprintf("%s(SB)", a.Sym.Name)
		} else if p != nil && p.Pcond != nil {
			str = fmt.Sprint(p.Pcond.Pc)
		} else if a.Val != nil {
			str = fmt.Sprint(a.Val.(*Prog).Pc)
		} else {
			str = fmt.Sprintf("%d(PC)", a.Offset)
		}

	case TYPE_INDIR:
		str = fmt.Sprintf("*%s", Mconv(a))

	case TYPE_MEM:
		str = Mconv(a)
		if a.Index != REG_NONE {
			str += fmt.Sprintf("(%v*%d)", Rconv(int(a.Index)), int(a.Scale))
		}

	case TYPE_CONST:
		if a.Reg != 0 {
			str = fmt.Sprintf("$%v(%v)", Mconv(a), Rconv(int(a.Reg)))
		} else {
			str = fmt.Sprintf("$%v", Mconv(a))
		}

	case TYPE_TEXTSIZE:
		if a.Val.(int32) == ArgsSizeUnknown {
			str = fmt.Sprintf("$%d", a.Offset)
		} else {
			str = fmt.Sprintf("$%d-%d", a.Offset, a.Val.(int32))
		}

	case TYPE_FCONST:
		str = fmt.Sprintf("%.17g", a.Val.(float64))
		// Make sure 1 prints as 1.0
		if !strings.ContainsAny(str, ".e") {
			str += ".0"
		}
		str = fmt.Sprintf("$(%s)", str)

	case TYPE_SCONST:
		str = fmt.Sprintf("$%q", a.Val.(string))

	case TYPE_ADDR:
		str = fmt.Sprintf("$%s", Mconv(a))

	case TYPE_SHIFT:
		v := int(a.Offset)
		op := string("<<>>->@>"[((v>>5)&3)<<1:])
		if v&(1<<4) != 0 {
			str = fmt.Sprintf("R%d%c%cR%d", v&15, op[0], op[1], (v>>8)&15)
		} else {
			str = fmt.Sprintf("R%d%c%c%d", v&15, op[0], op[1], (v>>7)&31)
		}
		if a.Reg != 0 {
			str += fmt.Sprintf("(%v)", Rconv(int(a.Reg)))
		}

	case TYPE_REGREG:
		str = fmt.Sprintf("(%v, %v)", Rconv(int(a.Reg)), Rconv(int(a.Offset)))

	case TYPE_REGREG2:
		str = fmt.Sprintf("%v, %v", Rconv(int(a.Reg)), Rconv(int(a.Offset)))

	case TYPE_REGLIST:
		str = regListConv(int(a.Offset))
	}

	return str
}

func Mconv(a *Addr) string {
	var str string

	switch a.Name {
	default:
		str = fmt.Sprintf("name=%d", a.Name)

	case NAME_NONE:
		switch {
		case a.Reg == REG_NONE:
			str = fmt.Sprint(a.Offset)
		case a.Offset == 0:
			str = fmt.Sprintf("(%v)", Rconv(int(a.Reg)))
		case a.Offset != 0:
			str = fmt.Sprintf("%d(%v)", a.Offset, Rconv(int(a.Reg)))
		}

	case NAME_EXTERN:
		str = fmt.Sprintf("%s%s(SB)", a.Sym.Name, offConv(a.Offset))

	case NAME_GOTREF:
		str = fmt.Sprintf("%s%s@GOT(SB)", a.Sym.Name, offConv(a.Offset))

	case NAME_STATIC:
		str = fmt.Sprintf("%s<>%s(SB)", a.Sym.Name, offConv(a.Offset))

	case NAME_AUTO:
		if a.Sym != nil {
			str = fmt.Sprintf("%s%s(SP)", a.Sym.Name, offConv(a.Offset))
		} else {
			str = fmt.Sprintf("%s(SP)", offConv(a.Offset))
		}

	case NAME_PARAM:
		if a.Sym != nil {
			str = fmt.Sprintf("%s%s(FP)", a.Sym.Name, offConv(a.Offset))
		} else {
			str = fmt.Sprintf("%s(FP)", offConv(a.Offset))
		}
	}
	return str
}

func offConv(off int64) string {
	if off == 0 {
		return ""
	}
	return fmt.Sprintf("%+d", off)
}

func regListConv(list int) string {
	str := ""

	for i := 0; i < 16; i++ { // TODO: 16 is ARM-specific.
		if list&(1<<uint(i)) != 0 {
			if str == "" {
				str += "["
			} else {
				str += ","
			}
			// This is ARM-specific; R10 is g.
			if i == 10 {
				str += "g"
			} else {
				str += fmt.Sprintf("R%d", i)
			}
		}
	}

	str += "]"
	return str
}
