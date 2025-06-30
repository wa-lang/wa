// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/wat/watutil"
	"wa-lang.org/wa/internal/wazero"
)

// 执行凹代码
func RunCode(cfg *config.Config, filename, code string, args ...string) (stdoutStderr []byte, err error) {
	// 编译为 wat 格式
	mainFunc, watBytes, fsetBytes, err := BuildFile(cfg, filename, code)
	if err != nil {
		return
	}

	// wat 编译为 wasm
	wasmBytes, err := watutil.Wat2Wasm(filename, watBytes)
	if err != nil {
		return
	}

	// main 执行
	stdout, stderr, err := wazero.RunWasm(filename, wasmBytes, fsetBytes, mainFunc, args...)
	stdoutStderr = append(stdout, stderr...)
	return
}
