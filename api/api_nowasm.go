// 版权 @2022 凹语言 作者。保留所有权利。

//go:build !wasm
// +build !wasm

package api

import (
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/wabt"
	"wa-lang.org/wa/internal/wazero"
)

// 执行凹代码
func RunCode(cfg *config.Config, filename, code string, mainFunc string, args ...string) (stdoutStderr []byte, err error) {
	// 编译为 wat 格式
	watBytes, err := BuildFile(cfg, filename, code)
	if err != nil {
		return
	}

	// wat 编译为 wasm
	wasmBytes, err := wabt.Wat2Wasm(watBytes)
	if err != nil {
		return
	}

	// main 执行
	stdout, stderr, err := wazero.RunWasm(cfg, filename, wasmBytes, mainFunc, args...)
	stdoutStderr = append(stdout, stderr...)
	return
}
