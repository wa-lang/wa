// 版权 @2022 凹语言 作者。保留所有权利。

//go:build wasm
// +build wasm

// wa 命令 js/wasm 版本, 用于 playground 环境.
package main

import (
	"fmt"
	"syscall/js"

	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/wat/watutil"
)

var waError error

func main() {
	window := js.Global().Get("window")

	// __WA_FILE_NAME__ 表示文件名, 用于区分中英文语法
	// __WA_CODE__ 代码内容
	waName := getJsValue(window, "__WA_FILE_NAME__", "hello.wa")
	waCode := getJsValue(window, "__WA_CODE__", "// no code")

	waClearError()

	outWat := waGenerateWat(waName, waCode)
	outFmt := waFormatCode(waName, waCode)
	outWasm := waGenerateWasm(waName, outWat)

	if !window.IsNull() && !window.IsUndefined() {
		window.Set("__WA_WAT__", outWat)
		window.Set("__WA_FMT_CODE__", outFmt)

		// 复制数组到 js
		jsArray := js.Global().Get("Uint8Array").New(len(outWasm))
		js.CopyBytesToJS(jsArray, outWasm)
		window.Set("__WA_WASM__", jsArray)

		window.Set("__WA_ERROR__", waGetErrorText())
	} else {
		fmt.Println(outWat)
	}
}

func waGenerateWasm(filename, code string) []byte {
	if waGetError() != nil {
		return nil
	}
	wasmBytes, err := watutil.Wat2Wasm(filename, []byte(code))
	if err != nil {
		if waGetError() == nil {
			waSetError(err)
		}
		return nil
	}
	return wasmBytes
}

func waGenerateWat(filename, code string) string {
	cfg := api.DefaultConfig()
	cfg.Target = api.WaOS_js

	wat, err := waBuildFile(cfg, filename, code)
	if err != nil {
		if waGetError() == nil {
			waSetError(err)
		}
		return ""
	}
	return string(wat)
}

func waFormatCode(filename, code string) string {
	newCode, err := api.FormatCode(filename, code)
	if err != nil {
		if waGetError() == nil {
			waSetError(err)
		}
		return code
	}
	return newCode
}

// 构建 wat 目标
func waBuildFile(cfg *config.Config, filename string, src interface{}) (wat []byte, err error) {
	prog, err := loader.LoadProgramFile(cfg, filename, src)
	if err != nil || prog == nil {
		logger.Tracef(&config.EnableTrace_api, "LoadProgramVFS failed, err = %v", err)
		return nil, err
	}

	watOut, err := compiler_wat.New().Compile(prog)
	return []byte(watOut), err
}

func getJsValue(window js.Value, key, defaultValue string) string {
	if window.IsNull() || window.IsUndefined() {
		return defaultValue
	}
	if x := window.Get(key); x.IsNull() || x.IsUndefined() {
		return defaultValue
	} else {
		return x.String()
	}
}

func waGetError() error {
	return waError
}

func waGetErrorText() string {
	if waError != nil {
		return waError.Error()
	} else {
		return ""
	}
}

func waSetError(err error) {
	waError = err
}

func waClearError() {
	waError = nil
}
