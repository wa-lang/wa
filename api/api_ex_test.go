// 版权 @2022 凹语言 作者。保留所有权利。

package api_test

import (
	"fmt"
	"log"

	"wa-lang.org/wa/api"
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

	output, err := api.RunCode(api.DefaultConfig(), "hello.wa", code, "__main__.main")
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
	cfg.Target = api.WaOS_wasi

	args := []string{"aa", "bb"}
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

	output, err := api.RunCode(api.DefaultConfig(), "hello.wa", code, "__main__.main")
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
		#wa:syntax=wz

		引于 "书"

		【启】：
		  书·说："你好，凹语言中文版！"
		。
	`

	output, err := api.RunCode(api.DefaultConfig(), "hello.wa", code, "__main__.main")
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
