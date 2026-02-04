// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"os"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

var DebugMode = false

func ParseFile(cpu abi.CPUType, fset *token.FileSet, filename string, src []byte) (f *ast.File, err error) {
	if fset == nil {
		panic("parser.ParseFile: no token.FileSet provided (fset == nil)")
	}

	// get source
	text := src
	if text == nil {
		text, err = os.ReadFile(filename)
	}
	if err != nil {
		return nil, err
	}

	p := newParser(cpu, fset, filename, text)
	p.trace = DebugMode

	f, err = p.ParseFile()
	return
}
