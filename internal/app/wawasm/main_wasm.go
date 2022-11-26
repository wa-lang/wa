// 版权 @2022 凹语言 作者。保留所有权利。

package main

import (
	"syscall/js"

	"github.com/wa-lang/wa/api"
)

func main() {
	window := js.Global().Get("window")
	waCode := window.Get("waCode").String()

	wat, err := api.BuildFile(api.DefaultConfig(), "hello.wa", waCode)
	if err != nil {
		window.Set("waWat", err.Error())
	} else {
		window.Set("waWat", string(wat))
	}
}
