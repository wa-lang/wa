// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"fmt"

	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

func assert(ok bool, message ...interface{}) {
	if !ok {
		if len(message) != 0 {
			panic(fmt.Sprint(append([]interface{}{"assert failed:"}, message...)...))
		} else {
			panic("assert failed")
		}
	}
}

func localSize(x *ast.Local) int {
	switch x.Type {
	case token.I32, token.I32_zh:
		return x.Cap * 4
	case token.I64, token.I64_zh:
		return x.Cap * 8
	case token.F32, token.F32_zh:
		return x.Cap * 4
	case token.F64, token.F64_zh:
		return x.Cap * 8
	default:
		panic("unreachable")
	}
}
