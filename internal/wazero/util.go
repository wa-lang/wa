package wazero

import (
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wazero"
	"wa-lang.org/wazero/sys"
)

func Wat2Wasm(source []byte) ([]byte, error) {
	return wazero.Wat2Wasm(source)
}

// 单次执行 wasm
func RunWasm(cfg *config.Config, wasmName string, wasmBytes []byte, mainFunc string, wasmArgs ...string) (stdout, stderr []byte, err error) {
	m, err := BuildModule(cfg, wasmName, wasmBytes, wasmArgs...)
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
