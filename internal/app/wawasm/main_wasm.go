// 版权 @2022 凹语言 作者。保留所有权利。

package main

import "syscall/js"

func getJsValue(window js.Value, key, defaultValue string) string {
	if x := window.Get(key); x.IsNull() {
		return defaultValue
	} else {
		return x.String()
	}
}

func main() {
	window := js.Global().Get("window")

	// __WA_FILE_NAME__ 表示文件名, 用于区分中英文语法
	// __WA_CODE__ 代码内容
	waName := getJsValue(window, "__WA_FILE_NAME__", "hello.wa")
	waCode := getJsValue(window, "__WA_CODE__", "// no code")

	waClearError()

	outWat := waGenerateWat(waName, waCode)
	outFmt := waFormatCode(waName, waCode)

	window.Set("__WA_WAT__", outWat)
	window.Set("__WA_FMT_CODE__", outFmt)
	window.Set("__WA_ERROR__", waGetErrorText())
}
