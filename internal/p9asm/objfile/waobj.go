// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Parsing of Go intermediate object files and archives.

package objfile

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/p9asm/obj/waobj"
)

type waobjFile struct {
	waobj *waobj.Package
}

func openWaobj(r *os.File) (rawFile, error) {
	f, err := waobj.Parse(r, `""`)
	if err != nil {
		return nil, err
	}
	return &waobjFile{f}, nil
}

func goobjName(id waobj.SymID) string {
	if id.Version == 0 {
		return id.Name
	}
	return fmt.Sprintf("%s<%d>", id.Name, id.Version)
}

func (f *waobjFile) symbols() ([]Sym, error) {
	seen := make(map[waobj.SymID]bool)

	var syms []Sym
	for _, s := range f.waobj.Syms {
		seen[s.SymID] = true
		sym := Sym{Addr: uint64(s.Data.Offset), Name: goobjName(s.SymID), Size: int64(s.Size), Type: s.Type.Name, Code: '?'}
		switch s.Kind {
		case waobj.STEXT, waobj.SELFRXSECT:
			sym.Code = 'T'
		case waobj.STYPE, waobj.SSTRING, waobj.SGOSTRING, waobj.SGOFUNC, waobj.SRODATA, waobj.SFUNCTAB, waobj.STYPELINK, waobj.SSYMTAB, waobj.SPCLNTAB, waobj.SELFROSECT:
			sym.Code = 'R'
		case waobj.SMACHOPLT, waobj.SELFSECT, waobj.SMACHO, waobj.SMACHOGOT, waobj.SNOPTRDATA, waobj.SINITARR, waobj.SDATA, waobj.SWINDOWS:
			sym.Code = 'D'
		case waobj.SBSS, waobj.SNOPTRBSS, waobj.STLSBSS:
			sym.Code = 'B'
		case waobj.SXREF, waobj.SMACHOSYMSTR, waobj.SMACHOSYMTAB, waobj.SMACHOINDIRECTPLT, waobj.SMACHOINDIRECTGOT, waobj.SFILE, waobj.SFILEPATH, waobj.SCONST, waobj.SDYNIMPORT, waobj.SHOSTOBJ:
			sym.Code = 'X' // should not see
		}
		if s.Version != 0 {
			sym.Code += 'a' - 'A'
		}
		syms = append(syms, sym)
	}

	for _, s := range f.waobj.Syms {
		for _, r := range s.Reloc {
			if !seen[r.Sym] {
				seen[r.Sym] = true
				sym := Sym{Name: goobjName(r.Sym), Code: 'U'}
				if s.Version != 0 {
					// should not happen but handle anyway
					sym.Code = 'u'
				}
				syms = append(syms, sym)
			}
		}
	}

	return syms, nil
}

// pcln does not make sense for Go object files, because each
// symbol has its own individual pcln table, so there is no global
// space of addresses to map.
func (f *waobjFile) pcln() (textStart uint64, symtab, pclntab []byte, err error) {
	return 0, nil, nil, fmt.Errorf("pcln not available in go object file")
}

// text does not make sense for Go object files, because
// each function has a separate section.
func (f *waobjFile) text() (textStart uint64, text []byte, err error) {
	return 0, nil, fmt.Errorf("text not available in go object file")
}

// waarch makes sense but is not exposed in debug/waobj's API,
// and we don't need it yet for any users of internal/objfile.
func (f *waobjFile) waarch() string {
	return "WAARCH unimplemented for debug/waobj files"
}
