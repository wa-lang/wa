package wazero

import (
	"wa-lang.org/wa/internal/3rdparty/wazero/sys"
)

// 单次执行 wasm
func RunWasm(wasmName string, wasmBytes []byte, mainFunc string, wasmArgs ...string) (stdout, stderr []byte, err error) {
	m, err := BuildModule(wasmName, wasmBytes, wasmArgs...)
	if err != nil {
		return
	}
	defer m.Close()

	return m.RunMain(mainFunc)
}

func AsExitError(err error) (exitCode int, ok bool) {
	errExit, ok := err.(*sys.ExitError)
	if ok {
		return int(errExit.ExitCode()), true
	}
	return 0, false
}
