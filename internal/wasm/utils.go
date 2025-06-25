package wasm

import (
	"encoding/binary"
	"io"
	"math"
	"strconv"
	"strings"
)

// FuncName returns the naming convention of "moduleName.funcName".
//
//   - moduleName is the possibly empty name the module was instantiated with.
//   - funcName is the name in the Custom Name section.
//   - funcIdx is the position in the function index namespace, prefixed with
//     imported functions.
//
// Note: "moduleName.$funcIdx" is used when the funcName is empty, as commonly
// the case in TinyGo.
func wasmdebug_FuncName(moduleName, funcName string, funcIdx uint32) string {
	var ret strings.Builder

	// Start module.function
	ret.WriteString(moduleName)
	ret.WriteByte('.')
	if funcName == "" {
		ret.WriteByte('$')
		ret.WriteString(strconv.Itoa(int(funcIdx)))
	} else {
		ret.WriteString(funcName)
	}

	return ret.String()
}

// DecodeFloat32 decodes a float32 in IEEE 754 binary representation.
// See https://www.w3.org/TR/2019/REC-wasm-core-1-20191205/#floating-point%E2%91%A2
func ieee754_DecodeFloat32(buf []byte) (float32, error) {
	if len(buf) < 4 {
		return 0, io.ErrUnexpectedEOF
	}

	raw := binary.LittleEndian.Uint32(buf[:4])
	return math.Float32frombits(raw), nil
}

// DecodeFloat64 decodes a float64 in IEEE 754 binary representation.
// See https://www.w3.org/TR/2019/REC-wasm-core-1-20191205/#floating-point%E2%91%A2
func ieee754_DecodeFloat64(buf []byte) (float64, error) {
	if len(buf) < 8 {
		return 0, io.ErrUnexpectedEOF
	}

	raw := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(raw), nil
}

// ValueTypeName returns the type name of the given ValueType as a string.
// These type names match the names used in the WebAssembly text format.
//
// Note: This returns "unknown", if an undefined ValueType value is passed.
func api_ValueTypeName(t ValueType) string {
	switch t {
	case ValueTypeI32:
		return "i32"
	case ValueTypeI64:
		return "i64"
	case ValueTypeF32:
		return "f32"
	case ValueTypeF64:
		return "f64"
	case ValueTypeExternref:
		return "externref"
	}
	return "unknown"
}
