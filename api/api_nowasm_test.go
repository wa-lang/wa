// 版权 @2022 凹语言 作者。保留所有权利。

package api_test

import (
	"fmt"
	"log"

	"wa-lang.org/wa/api"
)

func ExampleRunCode() {
	const code = `
		var gBase: i32 = 1000

		fn main() {
			println(add(40, 2) + gBase)
		}

		fn add(a: i32, b: i32) => i32 {
			return a+b
		}
	`

	output, err := api.RunCode(api.DefaultConfig(), "hello.wa", code)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(output))

	// Output:
	// 1042
}
