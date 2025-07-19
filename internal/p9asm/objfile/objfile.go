// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package objfile implements portable access to OS-specific executable files.
package objfile

import (
	"fmt"
	"os"
	"sort"

	"wa-lang.org/wa/internal/p9asm/debug/wasym"
)

type rawFile interface {
	symbols() (syms []Sym, err error)
	pcln() (textStart uint64, symtab, pclntab []byte, err error)
	text() (textStart uint64, text []byte, err error)
	waarch() string
}

// A File is an opened executable file.
type File struct {
	r   *os.File
	raw rawFile
}

// A Sym is a symbol defined in an executable file.
type Sym struct {
	Name string // symbol name
	Addr uint64 // virtual address of symbol
	Size int64  // size in bytes
	Code rune   // nm code (T for text, D for data, and so on)
	Type string // XXX?
}

// Open opens the named file.
// The caller must call f.Close when the file is no longer needed.
func Open(name string) (*File, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	if raw, err := openWaobj(r); err == nil {
		return &File{r, raw}, nil
	}
	if raw, err := openElf(r); err == nil {
		return &File{r, raw}, nil
	}
	if raw, err := openPE(r); err == nil {
		return &File{r, raw}, nil
	}
	if raw, err := openMacho(r); err == nil {
		return &File{r, raw}, nil
	}

	r.Close()
	return nil, fmt.Errorf("open %s: unrecognized object file", name)
}

func (f *File) Close() error {
	return f.r.Close()
}

func (f *File) Symbols() ([]Sym, error) {
	syms, err := f.raw.symbols()
	if err != nil {
		return nil, err
	}
	sort.Slice(syms, func(i, j int) bool {
		return syms[i].Addr < syms[j].Addr
	})
	return syms, nil
}

func (f *File) PCLineTable() (*wasym.Table, error) {
	textStart, symtab, pclntab, err := f.raw.pcln()
	if err != nil {
		return nil, err
	}
	return wasym.NewTable(symtab, wasym.NewLineTable(pclntab, textStart))
}

func (f *File) Text() (uint64, []byte, error) {
	return f.raw.text()
}

func (f *File) WAARCH() string {
	return f.raw.waarch()
}
