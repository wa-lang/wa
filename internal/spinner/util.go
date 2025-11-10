// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package spinner

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/types"
)

// 判断 e 是否为 "_"
func isBlankIdent(e ast.Expr) bool {
	id, ok := e.(*ast.Ident)
	return ok && id.Name == "_"
}

// 若 typ 为指针类型，返回其基类型，否则返回自身
func deref(typ types.Type) types.Type {
	if p, ok := typ.Underlying().(*types.Pointer); ok {
		return p.Elem()
	}
	return typ
}
