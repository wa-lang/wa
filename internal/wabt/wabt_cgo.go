// 版权 @2022 凹语言 作者。保留所有权利。

//go:build cgo
// +build cgo

package wabt

import (
	"fmt"
	"os"

	"github.com/wa-lang/wa/internal/3rdparty/wabt"
)

func Wat2WasmCmd(args ...string) {
	if err := wabt.Wat2WasmCmd(args...); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
