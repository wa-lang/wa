// 版权 @2022 凹语言 作者。保留所有权利。

//go:build wasm
// +build wasm

// wa 命令 js/wasm 版本, 用于 playground 环境.
package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"testing/fstest"
	"time"

	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/wat/watutil"
)

var waError error

func main() {
	runMain()
}

//go:wasmexport runMain
func runMain() {
	window := js.Global().Get("window")

	// __WA_FILE_NAME__ 表示文件名, 用于区分中英文语法
	// __WA_CODE__ 代码内容
	waName := getJsValue(window, "__WA_FILE_NAME__", "hello.wa")
	waCode := getJsValue(window, "__WA_CODE__", "// no code")
	waVfsJson := getJsValue(window, "__WA_VFS_JSON__", "")

	waClearError()

	var (
		outWat  string
		outFmt  string
		outWasm []byte
	)

	if waVfsJson == "" {
		outWat = waGenerateWat(waName, waCode)
		outFmt = waFormatCode(waName, waCode)
		outWasm = waGenerateWasm(waName, outWat)
	} else {
		vfs := waLoadVFSFromJson(waVfsJson, waName)

		outWat = waGenerateWatVFS(vfs, waName)
		outFmt = waFormatCode(waName, waCode)
		outWasm = waGenerateWasm(waName, outWat)
	}

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

// 从 vfs json 字符串构建 vfs
func waLoadVFSFromJson(waVfsJson, filename string) *config.PkgVFS {
	if waGetError() != nil {
		return nil
	}

	// json解析为扁平的map
	var m map[string]string
	if err := json.Unmarshal([]byte(waVfsJson), &m); err != nil {
		waSetError(err)
	}

	// 构建 map fs
	mapfs := make(fstest.MapFS)
	now := time.Now()

	// 遍历 JSON 文件路径
	for filePath, content := range m {
		mapfs[filePath] = &fstest.MapFile{
			Data:    []byte(content),
			Mode:    0644,
			ModTime: now,
		}
	}

	// 设置 app 文件系统
	return &config.PkgVFS{
		App: mapfs,
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

func waGenerateWatVFS(vfs *config.PkgVFS, filename string) string {
	cfg := api.DefaultConfig()
	cfg.Target = api.WaOS_js

	wat, err := waBuildVFS(vfs, cfg, filename)
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

// 从文件构建 wat 目标
func waBuildFile(cfg *config.Config, filename string, src interface{}) (wat []byte, err error) {
	prog, err := loader.LoadProgramFile(cfg, filename, src)
	if err != nil || prog == nil {
		logger.Tracef(&config.EnableTrace_api, "loader.LoadProgramFile failed, err = %v", err)
		return nil, err
	}

	watOut, err := compiler_wat.New().Compile(prog)
	return []byte(watOut), err
}

// 从文件系统构建 wat 目标
func waBuildVFS(vfs *config.PkgVFS, cfg *config.Config, pkgpath string) (wat []byte, err error) {
	prog, err := loader.LoadProgramVFS(vfs, cfg, pkgpath)
	if err != nil || prog == nil {
		logger.Tracef(&config.EnableTrace_api, "loader.LoadProgramVFS failed, err = %v", err)
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
