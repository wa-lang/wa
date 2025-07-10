// Inferno utils/6l/pass.c
// http://code.google.com/p/inferno-os/source/browse/utils/6l/pass.c
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

import (
	"encoding/binary"

	"wa-lang.org/wa/internal/p9asm/obj"
)

func preprocess(ctxt *obj.Link, cursym *obj.LSym) {
	panic("TODO")
}

func span6(ctxt *obj.Link, s *obj.LSym) {
	panic("TODO")
}

func follow(ctxt *obj.Link, s *obj.LSym) {
	panic("TODO")
}

func progedit(ctxt *obj.Link, p *obj.Prog) {
	panic("TODO")
}

var unaryDst = map[int]bool{
	ABSWAPL:    true,
	ABSWAPQ:    true,
	ACMPXCHG8B: true,
	ADECB:      true,
	ADECL:      true,
	ADECQ:      true,
	ADECW:      true,
	AINCB:      true,
	AINCL:      true,
	AINCQ:      true,
	AINCW:      true,
	ANEGB:      true,
	ANEGL:      true,
	ANEGQ:      true,
	ANEGW:      true,
	ANOTB:      true,
	ANOTL:      true,
	ANOTQ:      true,
	ANOTW:      true,
	APOPL:      true,
	APOPQ:      true,
	APOPW:      true,
	ASETCC:     true,
	ASETCS:     true,
	ASETEQ:     true,
	ASETGE:     true,
	ASETGT:     true,
	ASETHI:     true,
	ASETLE:     true,
	ASETLS:     true,
	ASETLT:     true,
	ASETMI:     true,
	ASETNE:     true,
	ASETOC:     true,
	ASETOS:     true,
	ASETPC:     true,
	ASETPL:     true,
	ASETPS:     true,
	AFFREE:     true,
	AFLDENV:    true,
	AFSAVE:     true,
	AFSTCW:     true,
	AFSTENV:    true,
	AFSTSW:     true,
	AFXSAVE:    true,
	AFXSAVE64:  true,
	ASTMXCSR:   true,
}

var Linkamd64 = obj.LinkArch{
	ByteOrder:  binary.LittleEndian,
	Name:       "amd64",
	Thechar:    '6',
	Preprocess: preprocess,
	Assemble:   span6,
	Follow:     follow,
	Progedit:   progedit,
	UnaryDst:   unaryDst,
	Minlc:      1,
	Ptrsize:    8,
	Regsize:    8,
}

var Link386 = obj.LinkArch{
	ByteOrder:  binary.LittleEndian,
	Name:       "386",
	Thechar:    '8',
	Preprocess: preprocess,
	Assemble:   span6,
	Follow:     follow,
	Progedit:   progedit,
	UnaryDst:   unaryDst,
	Minlc:      1,
	Ptrsize:    4,
	Regsize:    4,
}
