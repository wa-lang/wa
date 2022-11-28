// 版权 @2022 凹语言 作者。保留所有权利。

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	"wa-lang.org/wa/api"
)

func main() {
	output, err := api.RunCode("hello.wa", code)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(string(output))
}

const code = `
fn main() {
	println(add(40, 2))
}

fn add(a: i32, b: i32) => i32 {
	return a+b
}
`
