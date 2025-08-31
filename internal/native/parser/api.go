// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/ast"
)

var DebugMode = false

func ParseFile(path string, src []byte) (f *ast.File, err error) {
	p := newParser(path, src)
	p.trace = DebugMode

	f, err = p.ParseFile()
	return
}
