// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements a simple printer performance benchmark:
// go test -bench=BenchmarkPrint

//go:build ignore
// +build ignore

package printer

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"testing"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/parser"
)

var testfile *ast.File

func testprint(out io.Writer, file *ast.File) {
	if err := (&Config{TabIndent | UseSpaces, 8, 0}).Fprint(out, fset, file); err != nil {
		log.Fatalf("print error: %s", err)
	}
}

// cannot initialize in init because (printer) Fprint launches goroutines.
func initialize() {
	const filename = "testdata/parser.go"

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("%s", err)
	}

	file, err := parser.ParseFile(nil, fset, filename, src, parser.ParseComments)
	if err != nil {
		log.Fatalf("%s", err)
	}

	var buf bytes.Buffer
	testprint(&buf, file)
	if !bytes.Equal(buf.Bytes(), src) {
		log.Fatalf("print error: %s not idempotent", filename)
	}

	testfile = file
}

func BenchmarkPrint(b *testing.B) {
	if testfile == nil {
		initialize()
	}
	for i := 0; i < b.N; i++ {
		testprint(ioutil.Discard, testfile)
	}
}
