package wazero

import (
	"wa-lang.org/wa/internal/config"
)

// 单次执行 wasm
func RunWasm(cfg *config.Config, wasmName string, wasmBytes []byte, wasmArgs ...string) (stdout, stderr []byte, err error) {
	m, err := BuildModule(cfg, wasmName, wasmBytes, wasmArgs...)
	if err != nil {
		return
	}
	defer m.Close()

	return m.RunMain()
}
