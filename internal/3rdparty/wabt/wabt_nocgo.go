// 版权 @2022 凹语言 作者。保留所有权利。

//go:build !cgo
// +build !cgo

package wabt

import (
	"fmt"
	"os/exec"
	"runtime"
)

const Version = "1.0.29"

func Wat2WasmCmd(args ...string) error {
	exe := "wat2wasm"
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	cmd := exec.Command(exe, args...)
	stdoutStderr, err := cmd.CombinedOutput()
	if len(stdoutStderr) != 0 {
		fmt.Printf("%s\n", stdoutStderr)
	}
	if err != nil {
		return err
	}

	return nil
}
