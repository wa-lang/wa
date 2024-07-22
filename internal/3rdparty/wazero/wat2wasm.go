package wazero

import (
	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm/binary"
	"wa-lang.org/wa/internal/3rdparty/wazero/internalx/wasm/text"
)

func Wat2Wasm(source []byte) ([]byte, error) {
	if m, err := text.DecodeModule(source, api.CoreFeaturesV2); err != nil {
		return nil, err
	} else {
		return binary.EncodeModule(m), nil
	}
}
