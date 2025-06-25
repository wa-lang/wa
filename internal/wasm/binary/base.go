package binary

type ValueType = byte

const (
	// ValueTypeI32 is a 32-bit integer.
	ValueTypeI32 ValueType = 0x7f
	// ValueTypeI64 is a 64-bit integer.
	ValueTypeI64 ValueType = 0x7e
	// ValueTypeF32 is a 32-bit floating point number.
	ValueTypeF32 ValueType = 0x7d
	// ValueTypeF64 is a 64-bit floating point number.
	ValueTypeF64 ValueType = 0x7c

	ValueTypeExternref ValueType = 0x6f
)

// Magic is the 4 byte preamble (literally "\0asm") of the binary format
// See https://www.w3.org/TR/2019/REC-wasm-core-1-20191205/#binary-magic
const Magic = "\x00\x61\x73\x6D" // []byte{0x00, 0x61, 0x73, 0x6D}

// version is format version and doesn't change between known specification versions
// See https://www.w3.org/TR/2019/REC-wasm-core-1-20191205/#binary-version
const version = "\x01\x00\x00\x00" // []byte{0x01, 0x00, 0x00, 0x00}
