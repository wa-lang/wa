// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package api_test

import (
	"fmt"
	"log"

	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/token"
)

func ExampleRunCode() {
	const code = `
		global gBase: i32 = 1000

		func main() {
			println(add(40, 2) + gBase)
		}

		func add(a: i32, b: i32) => i32 {
			return a+b
		}
	`

	output, err := api.RunCode(api.DefaultConfig(), "hello.wa", code, token.K_pkg_main+"."+token.K_main)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(output))

	// Output:
	// 1042
}

func ExampleRunCode_args() {
	const code = `
		import "os"

		func main {
			for i, s := range os.Args[1:] {
				println(i, ":", s)
			}
		}
	`

	cfg := api.DefaultConfig()
	cfg.Target = api.WaOS_js

	args := []string{"hello.test", "aa", "bb"}
	output, err := api.RunCode(cfg, "hello.wa", code, args...)
	if err != nil {
		if len(output) != 0 {
			log.Println(string(output))
		}
		log.Fatal(err)
	}

	fmt.Print(string(output))

	// Output:
	// 0 : aa
	// 1 : bb
}

func ExampleRunCode_genericChainCalls() {
	const code = `
		type A: struct { }

		#wa:generic AddStr
		func A.Add(i: int) => *A {
			println(i)
			return this
		}

		func A.AddStr(s: string) => *A {
			println(s)
			return this
		}

		func main() {
			a: A
			a.Add("abc").Add("def")
		}
	`

	output, err := api.RunCode(api.DefaultConfig(), "hello.wa", code, token.K_pkg_main+"."+token.K_main)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(output))

	// Output:
	// abc
	// def
}

func ExampleRunCode_wz() {
	const code = `
		注: 版权 @2022 _examples/hello-zh 作者。保留所有权利。

		引入 "书"

		函数 主控:
		    书·说("你好，凹语言中文版！")
		完毕
	`

	output, err := api.RunCode(api.DefaultConfig(), "hello.wz", code, token.K_pkg_main+"."+token.K_主控)
	if err != nil {
		if len(output) != 0 {
			log.Println(string(output))
		}
		log.Fatal(err)
	}

	fmt.Print(string(output))

	// Output:
	// 你好，凹语言中文版！
}
