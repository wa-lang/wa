// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"bytes"
	"fmt"

	"wa-lang.org/wa/internal/ast"
)

func (p *printer) printNode(node interface{}) error {
	if n, ok := node.(*ast.File); ok {
		return p.printFile(n)
	}

	// if there are no comments, use node comments
	p.useNodeComments = p.comments == nil

	// get comments ready for use
	p.nextComment()

	// format node
	switch n := node.(type) {
	case ast.Expr:
		p.expr(n)
	case ast.Stmt:
		// A labeled statement will un-indent to position the label.
		// Set p.indent to 1 so we don't get indent "underflow".
		if _, ok := n.(*ast.LabeledStmt); ok {
			p.indent = 1
		}
		p.stmt(n, false)
	case ast.Decl:
		p.decl(n)
	case ast.Spec:
		if s, ok := n.(*ast.ValueSpec); ok {
			p.spec_ValueSpec(s, 1, false)
		} else {
			panic("unreachable")
		}
	case []ast.Stmt:
		// A labeled statement will un-indent to position the label.
		// Set p.indent to 1 so we don't get indent "underflow".
		for _, s := range n {
			if _, ok := s.(*ast.LabeledStmt); ok {
				p.indent = 1
			}
		}
		p.stmtList(n, 0, false)
	default:
		goto unsupported
	}

	return nil

unsupported:
	return fmt.Errorf("wa-lang.org/wa/internal/printer/w2printer: unsupported node type %T", node)
}

// nodeSize determines the size of n in chars after formatting.
// The result is <= maxSize if the node fits on one line with at
// most maxSize chars and the formatted output doesn't contain
// any control chars. Otherwise, the result is > maxSize.
func (p *printer) nodeSize(n ast.Node, maxSize int) (size int) {
	// nodeSize invokes the printer, which may invoke nodeSize
	// recursively. For deep composite literal nests, this can
	// lead to an exponential algorithm. Remember previous
	// results to prune the recursion (was issue 1628).
	if size, found := p.nodeSizes[n]; found {
		return size
	}

	size = maxSize + 1 // assume n doesn't fit
	p.nodeSizes[n] = size

	// nodeSize computation must be independent of particular
	// style so that we always get the same decision; print
	// in RawFormat
	var buf bytes.Buffer
	if err := fprintNode(&buf, p.fset, n, p.nodeSizes); err != nil {
		return
	}
	if buf.Len() <= maxSize {
		for _, ch := range buf.Bytes() {
			if ch < ' ' {
				return
			}
		}
		size = buf.Len() // n fits
		p.nodeSizes[n] = size
	}
	return
}
