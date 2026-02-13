// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x86asm

type SymLookup func(uint64) (string, uint64)

func memArgToSymbol(a Mem, pc uint64, instrLen int, symname SymLookup) (string, int64) {
	if a.Segment != 0 || a.Disp == 0 || a.Index != 0 || a.Scale != 0 {
		return "", 0
	}

	var disp uint64
	switch a.Base {
	case IP, EIP, RIP:
		disp = uint64(a.Disp + int64(pc) + int64(instrLen))
	case 0:
		disp = uint64(a.Disp)
	default:
		return "", 0
	}

	s, base := symname(disp)
	return s, int64(disp) - int64(base)
}
