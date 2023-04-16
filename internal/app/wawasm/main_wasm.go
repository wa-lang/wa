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

func waGenerateWat(filename, code string) string {
	cfg := api.DefaultConfig()
	cfg.WaArch = api.WaArch_wasm
	cfg.WaOS = api.WaOS_chrome

	wat, err := api.BuildFile(cfg, filename, code)
	if err != nil {
		waSetError(err)
		return ""
	}
	return string(wat)
}

func waFormatCode(filename, code string) string {
	newCode, err := api.FormatCode(filename, code)
	if err != nil {
		waSetError(err)
		return code
	}
	return newCode
}

func getJsValue(x js.Value, key, defaultValue string) string {
	window := js.Global().Get("window")
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
	waCode := getJsValue(window, "__WA_CODE__", "")

	waClearError()
	window.Set("__WA_WAT__", waGenerateWat(waName, waCode))
	window.Set("__WA_FMT_CODE__", waFormatCode(waName, waCode))
	window.Set("__WA_ERROR__", waGetError())
}
