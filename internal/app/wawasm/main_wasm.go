// 版权 @2022 凹语言 作者。保留所有权利。

package main

import (
	"syscall/js"

	"wa-lang.org/wa/api"
)

var waError string

func waGetError() string {
	return waError
}
func waSetError(err error) {
	if err != nil {
		waError = err.Error()
	} else {
		waError = ""
	}
}
func waClearError() {
	waError = ""
}

func waGenerateWat(code string) string {
	wat, err := api.BuildFile(api.DefaultConfig(), "hello.wa", code)
	if err != nil {
		waSetError(err)
		return ""
	}
	return string(wat)
}

func waFormatCode(code string) string {
	newCode, err := api.FormatCode("hello.wa", code)
	if err != nil {
		waSetError(err)
		return code
	}
	return newCode
}

func main() {
	window := js.Global().Get("window")
	waCode := window.Get("__WA_CODE__").String()

	waClearError()
	window.Set("__WA_WAT__", waGenerateWat(waCode))
	window.Set("__WA_FMT_CODE__", waFormatCode(waCode))
	window.Set("__WA_ERROR__", waGetError())
}
