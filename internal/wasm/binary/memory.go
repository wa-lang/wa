package binary

import (
	"bytes"

	"wa-lang.org/wa/internal/wasm"
	"wa-lang.org/wa/internal/wasm/api"
)

// decodeMemory returns the api.Memory decoded with the WebAssembly 1.0 (20191205) Binary Format.
//
// See https://www.w3.org/TR/2019/REC-wasm-core-1-20191205/#binary-memory
func decodeMemory(
	r *bytes.Reader,
	memorySizer func(minPages uint32, maxPages *uint32) (min, capacity, max uint32),
	memoryLimitPages uint32,
) (*wasm.Memory, error) {
	addrType, min, maxP, err := decodeLimitsType(r)
	if err != nil {
		return nil, err
	}

	min, capacity, max := memorySizer(min, maxP)
	mem := &wasm.Memory{AddrType: addrType, Min: min, Cap: capacity, Max: max, IsMaxEncoded: maxP != nil}

	return mem, mem.Validate(memoryLimitPages)
}

// encodeMemory returns the wasm.Memory encoded in WebAssembly 1.0 (20191205) Binary Format.
//
// See https://www.w3.org/TR/2019/REC-wasm-core-1-20191205/#binary-memory
func encodeMemory(i *wasm.Memory) []byte {
	maxPtr := &i.Max
	if !i.IsMaxEncoded {
		maxPtr = nil
	}
	if i.AddrType == api.ValueTypeI64 {
		return encodeLimitsType_i64(i.Min, maxPtr)
	}
	return encodeLimitsType(i.Min, maxPtr)
}
