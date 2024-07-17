// 版权 @2024 凹语言 作者。保留所有权利。

package parser_test

import (
	"fmt"

	"wa-lang.org/wa/internal/wat/parser"
)

func ExampleParseModule() {
	src := `(module $hello)`

	// Parse src but stop after processing the imports.
	m, err := parser.ParseModule("hello.wat", []byte(src))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(m)
	fmt.Println(m.Name)

	// output:
	// (module $hello)
	// $hello
}
