// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package api_test

import (
	"fmt"

	"wa-lang.org/wa/api"
)

func ExampleFormatCode() {
	s, err := api.FormatCode("hello.wa", "func add(a:i32, b:i32)=>i32 {return a+b}")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)

	// Output:
	// func add(a: i32, b: i32) => i32 { return a + b }
}

func ExampleFormatCode_structEmbeddingField() {
	s, err := api.FormatCode("hello.wa", `
type A :struct{  }
type B :struct{  A  }
`)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)

	// Output:
	// type A :struct{}
	// type B :struct{ A }
}
