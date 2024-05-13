package wazero

import (
	"wa-lang.org/wazero/api"
	"wa-lang.org/wazero/internal/wasm/binary"
	"wa-lang.org/wazero/internal/wasm/text"
)

func Wat2Wasm(source []byte) ([]byte, error) {
	if m, err := text.DecodeModule(source, api.CoreFeaturesV2); err != nil {
		return nil, err
	} else {
		return binary.EncodeModule(m), nil
	}
}
