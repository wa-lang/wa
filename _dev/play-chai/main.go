// 版权 @2022 凹语言 作者。保留所有权利。

package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flagDbgOut = flag.String("dbg", "a.out.dbg", "debug output file")
	flagWatOut = flag.String("wat", "a.out.wat", "wat output file")
)

func main() {
	flag.Parse()

	os.Remove(*flagDbgOut)
	os.Remove(*flagWatOut)

	fmt.Println("=== WARN: only for build wa.wasm === ")

	waClearError()

	outDbg := waDebugString(waName, waCode)
	outWat := waGenerateWat(waName, waCode)

	if serr := waGetErrorText(); serr != "" {
		fmt.Println("ERR:", serr)
		os.Exit(1)
	}

	os.WriteFile(*flagDbgOut, []byte(outDbg), 0666)
	os.WriteFile(*flagWatOut, []byte(outWat), 0666)

	fmt.Println("Done")
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
