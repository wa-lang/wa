// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import "wa-lang.org/wa/internal/ast"

// getNode returns the ast.CommentGroup associated with n, if any.
func getDoc(n ast.Node) *ast.CommentGroup {
	switch n := n.(type) {
	case *ast.Field:
		return n.Doc
	case *ast.ImportSpec:
		return n.Doc
	case *ast.ValueSpec:
		return n.Doc
	case *ast.TypeSpec:
		return n.Doc
	case *ast.GenDecl:
		return n.Doc
	case *ast.FuncDecl:
		return n.Doc
	case *ast.File:
		return n.Doc
	}
	return nil
}
