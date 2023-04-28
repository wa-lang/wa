// 版权 @2022 凹语言 作者。保留所有权利。

package main

import (
	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/logger"
)

var waError error

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

func waDebugString(filename, code string) string {
	cfg := api.DefaultConfig()
	cfg.WaArch = api.WaArch_wasm
	cfg.WaOS = api.WaOS_chrome

	prog, err := loader.LoadProgramFile(cfg, filename, code)
	if err != nil || prog == nil {
		return err.Error()
	}

	return prog.DebugString()
}

func waGenerateWat(filename, code string) string {
	cfg := api.DefaultConfig()
	cfg.WaArch = api.WaArch_wasm
	cfg.WaOS = api.WaOS_chrome

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

	watOut, err := compiler_wat.New().Compile(prog, "main")
	return []byte(watOut), err
}
