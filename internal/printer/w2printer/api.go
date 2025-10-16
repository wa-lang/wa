// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"io"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func Fprint(output io.Writer, fset *token.FileSet, f *ast.File) (err error) {
	var p printer
	p.init(fset, make(map[ast.Node]int))
	if err = p.printFile(f); err != nil {
		return
	}
	// print outstanding comments
	p.impliedSemi = false // EOF acts like a newline
	p.flush(token.Position{Offset: infinity, Line: infinity}, token.EOF)

	output = &trimmer{output: output}
	_, err = output.Write(p.output)
	return
}

// fprint implements Fprint and takes a nodesSizes map for setting up the printer state.
func fprintNode(output io.Writer, fset *token.FileSet, node interface{}, nodeSizes map[ast.Node]int) (err error) {
	var p printer
	p.init(fset, nodeSizes)
	if err = p.printNode(node); err != nil {
		return
	}
	// print outstanding comments
	p.impliedSemi = false // EOF acts like a newline
	p.flush(token.Position{Offset: infinity, Line: infinity}, token.EOF)

	// redirect output through a trimmer to eliminate trailing whitespace
	// (Input to a tabwriter must be untrimmed since trailing tabs provide
	// formatting information. The tabwriter could provide trimming
	// functionality but no tabwriter is used when RawFormat is set.)
	output = &trimmer{output: output}

	// write printer result via tabwriter/trimmer to output
	if _, err = output.Write(p.output); err != nil {
		return
	}

	return
}
