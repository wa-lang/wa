// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	"wa-lang.org/wa/api"
)

func main() {
	output, err := api.RunCode(api.DefaultConfig(), "hello.wa", code)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(string(output))

	output, err = api.RunCode(api.DefaultConfig(), "hello-zh.wz", code_zh)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(string(output))
}

const code = `
func main {
	println(add(40, 2))
}

func add(a: i32, b: i32) => i32 {
	return a+b
}
`

const code_zh = `
引入 "书"

函数·主控：
	书·说("你好，凹语言中文版！")
完毕
`
