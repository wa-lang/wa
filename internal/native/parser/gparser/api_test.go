// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package gparser_test

import (
	"testing"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/parser/gparser"
	"wa-lang.org/wa/internal/native/token"
)

func TestParseFile(t *testing.T) {
	gparser.DebugMode = true

	t.Skip("TODO")

	fset := token.NewFileSet()
	file, err := gparser.ParseFile(abi.LOONG64, fset, "./testdata/hello-01/app.s", nil)
	if err != nil {
		t.Fatal(err)
	}

	if file.CPU != abi.LOONG64 {
		t.Fatalf("CPU: invalid")
	}
}
