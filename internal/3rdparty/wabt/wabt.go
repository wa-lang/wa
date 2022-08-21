// 版权 @2022 凹语言 作者。保留所有权利。

//go:build cgo
// +build cgo

package wabt

/*
#cgo LDFLAGS:
#cgo CPPFLAGS: -I./internal/wabt-1.0.29
#cgo CXXFLAGS: -std=c++17

#include <stdlib.h>

extern int wat2wasmMain(int argc, char** argv);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const Version = "1.0.29"

func Wat2WasmCmd(args ...string) error {
	args = append([]string{"wat2wasm"}, args...)

	argc := len(args)
	argv := make([]*C.char, argc, argc+1)

	for i, s := range args {
		argv[i] = C.CString(s)
		defer C.free(unsafe.Pointer(argv[i]))
	}

	rv := C.wat2wasmMain(C.int(argc), (**C.char)(&argv[0]))
	if rv != 0 {
		return fmt.Errorf("wat2wasm failed: err-code=%d", rv)
	}

	return nil
}
