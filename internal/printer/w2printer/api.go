// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"io"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

// A Mode value is a set of flags (or 0). They control printing.
type Mode uint

const (
	RawFormat Mode = 1 << iota // do not use a tabwriter; if set, UseSpaces is ignored
	TabIndent                  // use tabs for indentation independent of UseSpaces
	UseSpaces                  // use spaces instead of tabs for alignment
	SourcePos                  // emit //line directives to preserve original source positions
)

// A Config node controls the output of Fprint.
type Config struct {
	Mode     Mode // default: 0
	Tabwidth int  // default: 8
	Indent   int  // default: 0 (all code is indented at least by this much)
}

func Fprint(output io.Writer, fset *token.FileSet, node interface{}) error {
	return (&Config{Tabwidth: 8}).Fprint(output, fset, node)
}

func (cfg *Config) Fprint(output io.Writer, fset *token.FileSet, node interface{}) error {
	return cfg.fprint(output, fset, node, make(map[ast.Node]int))
}
