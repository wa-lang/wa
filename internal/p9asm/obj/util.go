// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package obj

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const REG_NONE = 0

type Biobuf struct {
	f       *os.File
	r       *bufio.Reader
	w       *bufio.Writer
	linelen int
}

func Bopenw(name string) (*Biobuf, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &Biobuf{f: f, w: bufio.NewWriter(f)}, nil
}

func Bopenr(name string) (*Biobuf, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return &Biobuf{f: f, r: bufio.NewReader(f)}, nil
}

func Binitw(w io.Writer) *Biobuf {
	return &Biobuf{w: bufio.NewWriter(w)}
}

func (b *Biobuf) Write(p []byte) (int, error) {
	return b.w.Write(p)
}

func Bwritestring(b *Biobuf, p string) (int, error) {
	return b.w.WriteString(p)
}

func Bseek(b *Biobuf, offset int64, whence int) int64 {
	if b.w != nil {
		if err := b.w.Flush(); err != nil {
			log.Fatalf("writing output: %v", err)
		}
	} else if b.r != nil {
		if whence == 1 {
			offset -= int64(b.r.Buffered())
		}
	}
	off, err := b.f.Seek(offset, whence)
	if err != nil {
		log.Fatalf("seeking in output: %v", err)
	}
	if b.r != nil {
		b.r.Reset(b.f)
	}
	return off
}

func Boffset(b *Biobuf) int64 {
	if b.w != nil {
		if err := b.w.Flush(); err != nil {
			log.Fatalf("writing output: %v", err)
		}
	}
	off, err := b.f.Seek(0, 1)
	if err != nil {
		log.Fatalf("seeking in output [0, 1]: %v", err)
	}
	if b.r != nil {
		off -= int64(b.r.Buffered())
	}
	return off
}

func (b *Biobuf) Flush() error {
	return b.w.Flush()
}

func Bputc(b *Biobuf, c byte) {
	b.w.WriteByte(c)
}

const Beof = -1

func Bread(b *Biobuf, p []byte) int {
	n, err := io.ReadFull(b.r, p)
	if n == 0 {
		if err != nil && err != io.EOF {
			n = -1
		}
	}
	return n
}

func Bgetc(b *Biobuf) int {
	c, err := b.r.ReadByte()
	if err != nil {
		return -1
	}
	return int(c)
}

func Bgetrune(b *Biobuf) int {
	r, _, err := b.r.ReadRune()
	if err != nil {
		return -1
	}
	return int(r)
}

func Bungetrune(b *Biobuf) {
	b.r.UnreadRune()
}

func (b *Biobuf) Read(p []byte) (int, error) {
	return b.r.Read(p)
}

func (b *Biobuf) Peek(n int) ([]byte, error) {
	return b.r.Peek(n)
}

func Brdline(b *Biobuf, delim int) string {
	s, err := b.r.ReadBytes(byte(delim))
	if err != nil {
		log.Fatalf("reading input: %v", err)
	}
	b.linelen = len(s)
	return string(s)
}

func Brdstr(b *Biobuf, delim int, cut int) string {
	s, err := b.r.ReadString(byte(delim))
	if err != nil {
		log.Fatalf("reading input: %v", err)
	}
	if len(s) > 0 && cut > 0 {
		s = s[:len(s)-1]
	}
	return s
}

func Access(name string, mode int) int {
	if mode != 0 {
		panic("bad access")
	}
	_, err := os.Stat(name)
	if err != nil {
		return -1
	}
	return 0
}

func Blinelen(b *Biobuf) int {
	return b.linelen
}

func Bterm(b *Biobuf) error {
	var err error
	if b.w != nil {
		err = b.w.Flush()
	}
	err1 := b.f.Close()
	if err == nil {
		err = err1
	}
	return err
}

func (ctxt *Link) NewProg() *Prog {
	p := new(Prog) // should be the only call to this; all others should use ctxt.NewProg
	p.Ctxt = ctxt
	return p
}

const (
	ABase386 = (1 + iota) << 11
	ABaseARM
	ABaseAMD64
	ABaseARM64
	ABaseLoong64
	ABaseRISCV
	ABaseWasm

	AllowedOpCodes = 1 << 11            // The number of opcodes available for any given architecture.
	AMask          = AllowedOpCodes - 1 // AND with this to use the opcode as an array index.
)

type opSet struct {
	lo    int
	names []string
}

// Not even worth sorting
var aSpace []opSet

// RegisterOpcode binds a list of instruction names
// to a given instruction number range.
func RegisterOpcode(lo int, Anames []string) {
	aSpace = append(aSpace, opSet{lo, Anames})
}

func Aconv(a int) string {
	if a < int(A_ARCHSPECIFIC) {
		return Anames[a]
	}
	for i := range aSpace {
		as := &aSpace[i]
		if as.lo <= a && a < as.lo+len(as.names) {
			return as.names[a-as.lo]
		}
	}
	return fmt.Sprintf("A???%d", a)
}

var Anames = []string{
	"XXX",
	"CALL",
	"CHECKNIL",
	"DATA",
	"DUFFCOPY",
	"DUFFZERO",
	"END",
	"FUNCDATA",
	"GLOBL",
	"JMP",
	"NOP",
	"PCDATA",
	"RET",
	"TEXT",
	"TYPE",
	"UNDEF",
	"USEFIELD",
	"VARDEF",
	"VARKILL",
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
		// TODO(rsc): This special case is for x86 instructions like
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

type regSet struct {
	lo    int
	hi    int
	Rconv func(int) string
}

// Few enough architectures that a linear scan is fastest.
// Not even worth sorting.
var regSpace []regSet

func Rconv(reg int) string {
	if reg == REG_NONE {
		return "NONE"
	}
	for i := range regSpace {
		rs := &regSpace[i]
		if rs.lo <= reg && reg < rs.hi {
			return rs.Rconv(reg)
		}
	}
	return fmt.Sprintf("R???%d", reg)
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

/*
	Each architecture defines a register space as a unique
	integer range.
	Here is the list of architectures and the base of their register spaces.
*/

const (
	// Because of masking operations in the encodings, each register
	// space should start at 0 modulo some power of 2.
	RBase386     = 1 * 1024
	RBaseAMD64   = 2 * 1024
	RBaseARM     = 3 * 1024
	RBaseARM64   = 8 * 1024  // range [8k, 13k)
	RBaseRISCV   = 15 * 1024 // range [15k, 16k)
	RBaseWasm    = 16 * 1024
	RBaseLOONG64 = 19 * 1024 // range [19K, 22k)
)
