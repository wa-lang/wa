// 版权 @2022 凹语言 作者。保留所有权利。

//go:build !wasm
// +build !wasm

package main

import "fmt"

func main() {
	fmt.Println("=== WARN: only for build wa.wasm === ")

	waClearError()

	outWat := waGenerateWat(waName, waCode)
	outFmt := waFormatCode(waName, waCode)

	fmt.Println("ERR:", waGetErrorText())
	fmt.Println("FMT:", outFmt)
	fmt.Println("WAT:", outWat)
}

const waName = "hello.wa"
const waCode = `// 版权 @2019 凹语言 作者。保留所有权利。

import "fmt"
import "runtime"

func main {
	println("你好，凹语言！", runtime.WAOS)
	println(add(40, 2))

	fmt.Println(1+1)
}

func add(a: i32, b: i32) => i32 {
	return a+b
}
`
