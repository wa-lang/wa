// Derived from Inferno utils/6l/obj.c and utils/6l/span.c
// http://code.google.com/p/inferno-os/source/browse/utils/6l/obj.c
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

package obj

import (
	"encoding/binary"
	"io"
	"log"

	"wa-lang.org/wa/internal/p9asm/objabi"
)

// LinkArch is the definition of a single architecture.
type LinkArch struct {
	ByteOrder  binary.ByteOrder
	Name       string
	Thechar    int
	Preprocess func(*Link, *LSym)
	Assemble   func(*Link, *LSym)
	Follow     func(*Link, *LSym)
	Progedit   func(*Link, *Prog)
	UnaryDst   map[objabi.As]bool // Instruction takes one operand, a destination.
	Minlc      int
	Ptrsize    int
	Regsize    int
}

// Link holds the context for writing object code from a compiler
// to be linker input or for reading that input into the linker.
type Link struct {
	Headtype           objabi.HeadType
	Arch               *LinkArch
	Fset               *objabi.FileSet
	Debugasm           int32
	Debugvlog          int32
	Debugzerostack     int32
	Debugdivmod        int32
	Debugpcln          int32
	Flag_shared        int32
	Flag_dynlink       bool
	Bso                io.Writer
	Windows            int32
	Enforce_data_order int32
	Hash               map[SymVer]*LSym
	Imports            []string
	Plist              *Plist
	Plast              *Plist
	Sym_div            *LSym
	Sym_divu           *LSym
	Sym_mod            *LSym
	Sym_modu           *LSym
	Tlsg               *LSym
	Curp               *Prog
	Printp             *Prog
	Blitrl             *Prog
	Elitrl             *Prog
	Rexflag            int
	Rep                int
	Repn               int
	Lock               int
	Asmode             int
	Andptr             []byte
	And                [100]uint8
	Instoffset         int64
	Autosize           int32
	Armsize            int32
	Pc                 int64
	Tlsoffset          int
	Diag               func(string, ...interface{})
	Mode               int
	Cursym             *LSym
	Version            int
	Textp              *LSym
	Etextp             *LSym
}

func Linknew(arch *LinkArch, targetOS, _workDir string) *Link {
	ctxt := new(Link)
	ctxt.Fset = objabi.NewFileSet()
	ctxt.Hash = make(map[SymVer]*LSym)
	ctxt.Arch = arch
	ctxt.Version = HistVersion

	if err := ctxt.Headtype.Set(targetOS); err != nil {
		log.Fatal(err)
	}

	// Record thread-local storage offset.
	// TODO(chai2010): Move tlsoffset back into the linker.
	switch ctxt.Headtype {
	default:
		log.Fatalf("unknown thread-local storage offset for %v", ctxt.Headtype)

	case objabi.Hwindows:
		break

		/*
		 * ELF uses TLS offset negative from FS.
		 * Translate 0(FS) and 8(FS) into -16(FS) and -8(FS).
		 * Known to low-level assembly in package runtime and runtime/cgo.
		 */
	case objabi.Hlinux:
		ctxt.Tlsoffset = -1 * ctxt.Arch.Ptrsize

		/*
		 * OS X system constants - offset from 0(GS) to our TLS.
		 * Explained in ../../runtime/cgo/gcc_darwin_*.c.
		 */
	case objabi.Hdarwin:
		switch ctxt.Arch.Thechar {
		default:
			log.Fatalf("unknown thread-local storage offset for darwin/%s", ctxt.Arch.Name)

		case '5':
			ctxt.Tlsoffset = 0 // dummy value, not needed

		case '6':
			ctxt.Tlsoffset = 0x8a0

		case '7':
			ctxt.Tlsoffset = 0 // dummy value, not needed

		case '8':
			ctxt.Tlsoffset = 0x468
		}
	}

	return ctxt
}

func (ctxt *Link) Lookup(symb string, v int16) *LSym {
	s := ctxt.Hash[SymVer{symb, v}]
	if s != nil {
		return s
	}

	s = &LSym{Name: symb, Version: v}
	ctxt.Hash[SymVer{symb, v}] = s

	return s
}
