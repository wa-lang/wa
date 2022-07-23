package llutil

import (
	"encoding/binary"
	"strings"

	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// NewZero returns a new zero value of the given type.
func NewZero(typ types.Type) value.Value {
	switch typ := typ.(type) {
	case *types.IntType:
		return constant.NewInt(typ, 0)
	case *types.FloatType:
		return constant.NewFloat(typ, 0)
	default:
		return constant.NewZeroInitializer(typ)
	}
}

// NewCString returns a new NULL-terminated character array constant based on
// the given UTF-8 string contents.
func NewCString(s string) *constant.CharArray {
	return constant.NewCharArrayFromString(s + "\x00")
}

// NewPascalString returns a length-prefixed string, also known as a Pascal
// string. The length prefix is stored as a 32-bit integer with big-endian
// encoding.
func NewPascalString(s string) *constant.CharArray {
	var buf strings.Builder
	var length [4]byte
	binary.BigEndian.PutUint32(length[:], uint32(len(s)))
	for _, b := range length {
		buf.WriteByte(b)
	}
	buf.WriteString(s)
	return constant.NewCharArrayFromString(buf.String())
}
