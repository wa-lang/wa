// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package obj

import (
	"bytes"
	"fmt"
	"strings"
)

func (p *Prog) Line() string {
	return p.Ctxt.LineHist.LineString(int(p.Lineno))
}

func (p *Prog) String() string {
	if p.Ctxt == nil {
		return "<Prog without ctxt>"
	}

	sc := CConv(p.Scond)

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%.5d (%v)\t%v%s", p.Pc, p.Line(), p.As, sc)
	sep := "\t"
	if p.From.Type != TYPE_NONE {
		fmt.Fprintf(&buf, "%s%v", sep, Dconv(p, &p.From))
		sep = ", "
	}
	if p.Reg != REG_NONE {
		// Should not happen but might as well show it if it does.
		fmt.Fprintf(&buf, "%s%v", sep, p.Reg)
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
		fmt.Fprintf(&buf, "%s%v", sep, p.RegTo2)
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
			str = fmt.Sprintf("%v(%v)(NONE)", a, a.Reg)
		}

	case TYPE_REG:
		// TODO(chai2010): This special case is for x86 instructions like
		//	PINSRQ	CX,$1,X6
		// where the $1 is included in the p->to Addr.
		// Move into a new field.
		if a.Offset != 0 {
			str = fmt.Sprintf("$%d,%v", a.Offset, a.Reg)
			break
		}

		str = a.Reg.String()
		if a.Name != NAME_NONE || a.Sym != nil {
			str = fmt.Sprintf("%v(%v)(REG)", a, a.Reg)
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
		str = fmt.Sprintf("*%v", a)

	case TYPE_MEM:
		str = a.String()
		if a.Index != int16(REG_NONE) {
			str += fmt.Sprintf("(%v*%d)", RBaseType(a.Index), int(a.Scale))
		}

	case TYPE_CONST:
		if a.Reg != 0 {
			str = fmt.Sprintf("$%v(%v)", a, a.Reg)
		} else {
			str = fmt.Sprintf("$%v", a)
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
		str = fmt.Sprintf("$%v", a)

	case TYPE_SHIFT:
		v := int(a.Offset)
		op := string("<<>>->@>"[((v>>5)&3)<<1:])
		if v&(1<<4) != 0 {
			str = fmt.Sprintf("R%d%c%cR%d", v&15, op[0], op[1], (v>>8)&15)
		} else {
			str = fmt.Sprintf("R%d%c%c%d", v&15, op[0], op[1], (v>>7)&31)
		}
		if a.Reg != 0 {
			str += fmt.Sprintf("(%v)", a.Reg)
		}

	case TYPE_REGREG:
		str = fmt.Sprintf("(%v, %v)", a.Reg, RBaseType(a.Offset))

	case TYPE_REGREG2:
		str = fmt.Sprintf("%v, %v", a.Reg, RBaseType(a.Offset))

	case TYPE_REGLIST:
		str = regListConv(int(a.Offset))
	}

	return str
}

// 有bit位组成的寄存器列表转位字符串格式
func regListConv(list int) string {
	// 通常出现在ARM, 最多有16个寄存器列表
	var sb strings.Builder
	for i := 0; i < 16; i++ {
		if list&(1<<uint(i)) != 0 {
			if sb.Len() == 0 {
				// 需要区分是否为第一个出现的寄存器
				sb.WriteRune('[')
			} else {
				sb.WriteRune(',')
			}
			// 寄存器列表是 ARM 的用法, R10 对应 g 寄存器
			if i == 10 {
				sb.WriteRune('g')
			} else {
				sb.WriteString(fmt.Sprintf("R%d", i))
			}
		}
	}
	// 有可能没有任何寄存器
	if sb.Len() == 0 {
		sb.WriteRune('[')
	}
	sb.WriteRune(']')
	return sb.String()
}
