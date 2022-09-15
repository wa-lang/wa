// 版权 @2022 凹语言 作者。保留所有权利。

//go:build !wasm
// +build !wasm

package api

import (
	"os"

	"github.com/wa-lang/wa/internal/app/apputil"
)

// 执行凹代码
func RunCode(filename, code string) (stdoutStderr []byte, err error) {
	// 编译为 wat 格式
	wat, err := BuildFile(filename, code, Machine_Wasm32_wa)

	// wat 写到临时文件
	outfile := "a.out.wat"
	if !*FlagDebugMode {
		defer os.Remove(outfile)
	}
	if err = os.WriteFile(outfile, []byte(wat), 0666); err != nil {
		return nil, err
	}

	// 执行 wat 文件
	stdoutStderr, err = apputil.RunWasm(outfile)
	return
}
