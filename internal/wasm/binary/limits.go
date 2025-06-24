package binary

import (
	"bytes"
	"fmt"

	"wa-lang.org/wa/internal/wasm/api"
	"wa-lang.org/wa/internal/wasm/leb128"
)

// decodeLimitsType returns the `limitsType` (min, max) decoded with the WebAssembly 1.0 (20191205) Binary Format.
//
// See https://www.w3.org/TR/2019/REC-wasm-core-1-20191205/#limits%E2%91%A6
func decodeLimitsType(r *bytes.Reader) (addrType api.ValueType, min uint32, max *uint32, err error) {
	var flag byte
	if flag, err = r.ReadByte(); err != nil {
		err = fmt.Errorf("read leading byte: %v", err)
		return
	}

	if flag < 0x04 {
		addrType = api.ValueTypeI32
	} else {
		addrType = api.ValueTypeI64
	}

	switch flag {
	case 0x00, 0x04:
		min, _, err = leb128.DecodeUint32(r)
		if err != nil {
			err = fmt.Errorf("read min of limit: %v", err)
		}
	case 0x01, 0x05:
		min, _, err = leb128.DecodeUint32(r)
		if err != nil {
			err = fmt.Errorf("read min of limit: %v", err)
			return
		}
		var m uint32
		if m, _, err = leb128.DecodeUint32(r); err != nil {
			err = fmt.Errorf("read max of limit: %v", err)
		} else {
			max = &m
		}
	default:
		err = fmt.Errorf("%v for limits: %#x != 0x00 or 0x01 or 0x04 or 0x05", ErrInvalidByte, flag)
	}
	return
}

// encodeLimitsType returns the `limitsType` (min, max) encoded in WebAssembly 1.0 (20191205) Binary Format.
//
// See https://www.w3.org/TR/2019/REC-wasm-core-1-20191205/#limits%E2%91%A6
func encodeLimitsType(min uint32, max *uint32) []byte {
	if max == nil {
		return append(leb128.EncodeUint32(0x00), leb128.EncodeUint32(min)...)
	}
	return append(leb128.EncodeUint32(0x01), append(leb128.EncodeUint32(min), leb128.EncodeUint32(*max)...)...)
}

// encodeLimitsType for memory64
func encodeLimitsType_i64(min uint32, max *uint32) []byte {
	// 即使在 memory64, 使用 uint32 表示 page 数目依然足够
	if max == nil {
		return append(leb128.EncodeUint32(0x04), leb128.EncodeUint32(min)...)
	}
	return append(leb128.EncodeUint32(0x05), append(leb128.EncodeUint32(min), leb128.EncodeUint32(*max)...)...)
}
