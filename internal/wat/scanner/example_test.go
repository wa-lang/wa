// 版权 @2024 凹语言 作者。保留所有权利。

package scanner_test

import (
	"fmt"

	"wa-lang.org/wa/internal/wat/scanner"
	"wa-lang.org/wa/internal/wat/token"
)

func ExampleScanner_Scan() {
	var src = []byte("(module $__walang__)")
	var file = token.NewFile("", len(src))

	var s scanner.Scanner
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s %q\n", tok, lit)
	}

	// output:
	// ( ""
	// module "module"
	// IDENT "__walang__"
	// ) ""
}
