// 版权 @2022 凹语言 作者。保留所有权利。

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
