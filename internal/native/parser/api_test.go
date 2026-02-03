// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser_test

import (
	"testing"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/parser"
	"wa-lang.org/wa/internal/native/token"
)

func TestParseFile(t *testing.T) {
	parser.DebugMode = false

	fset := token.NewFileSet()

	fileGas, err := parser.ParseFile(abi.LOONG64, fset, "./testdata/hello-01/app.wa.s", nil)
	if err != nil {
		t.Fatal(err)
	}
	if fileGas.CPU != abi.LOONG64 {
		t.Fatalf("CPU: invalid")
	}

	fileZh, err := parser.ParseFile(abi.LOONG64, fset, "./testdata/hello-01/app.wz.s", nil)
	if err != nil {
		t.Fatal(err)
	}
	if fileZh.CPU != abi.LOONG64 {
		t.Fatalf("CPU: invalid")
	}
}
