// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner_test

import (
	"fmt"

	"wa-lang.org/wa/internal/scanner"
	"wa-lang.org/wa/internal/token"
)

func ExampleScanner_Scan() {
	// src is the input that we want to tokenize.
	src := []byte("cos(x) + 1i*sin(x) # Euler")

	// Initialize the scanner.
	var s scanner.Scanner
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

	// Repeated calls to Scan yield the token sequence found in the input.
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}

	// output:
	// 1:1	IDENT	"cos"
	// 1:4	(	""
	// 1:5	IDENT	"x"
	// 1:6	)	""
	// 1:8	+	""
	// 1:10	IMAG	"1i"
	// 1:12	*	""
	// 1:13	IDENT	"sin"
	// 1:16	(	""
	// 1:17	IDENT	"x"
	// 1:18	)	""
	// 1:20	;	"\n"
	// 1:20	COMMENT	"# Euler"
}

func ExampleScanner_Scan_w2() {
	const code = `
注: 这是注释

函数 主控:
	输出 ("Hello, world!")
完毕
`
	fset := token.NewFileSet()
	file := fset.AddFile("hello.wz", fset.Base(), len(code))

	var s scanner.Scanner
	s.W2Mode = true // 中文模式
	s.Init(file, []byte(code), nil, scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}

	// output:
	// hello.wz:2:1	COMMENT	"注: 这是注释"
	// hello.wz:4:1	函数	"函数"
	// hello.wz:4:8	IDENT	"主控"
	// hello.wz:4:14	:	""
	// hello.wz:5:2	IDENT	"输出"
	// hello.wz:5:9	(	""
	// hello.wz:5:10	STRING	"\"Hello, world!\""
	// hello.wz:5:25	)	""
	// hello.wz:5:26	;	"\n"
	// hello.wz:6:1	完毕	"完毕"
	// hello.wz:6:7	;	"\n"
}
