// 版权 @2025 凹语言 作者。保留所有权利。

package leb128

// EncodeUint32 encodes the value into a buffer in LEB128 format
//
// See https://en.wikipedia.org/wiki/LEB128#Encode_unsigned_integer
func EncodeUint32(v uint32) []byte {
	dst := make([]byte, 10)
	n := encodeUint64(uint64(v), dst)
	return dst[:n]
}

// EncodeUint64 encodes the value into a buffer in LEB128 format
//
// See https://en.wikipedia.org/wiki/LEB128#Encode_unsigned_integer
func EncodeUint64(v uint64) (buf []byte) {
	dst := make([]byte, 10)
	n := encodeUint64(uint64(v), dst)
	return dst[:n]
}

// EncodeInt32 encodes the signed value into a buffer in LEB128 format
//
// See https://en.wikipedia.org/wiki/LEB128#Encode_signed_integer
func EncodeInt32(v int32) []byte {
	dst := make([]byte, 10)
	n := encodeInt64(int64(v), dst)
	return dst[:n]
}

// EncodeInt64 encodes the signed value into a buffer in LEB128 format
//
// See https://en.wikipedia.org/wiki/LEB128#Encode_signed_integer
func EncodeInt64(v int64) []byte {
	dst := make([]byte, 10)
	n := encodeInt64(int64(v), dst)
	return dst[:n]
}

func encodeUint64(v uint64, dst []byte) int {
	var c uint8

	length := uint8(0)
	for {
		c = uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80
		}
		if dst != nil {
			dst[0] = byte(c)
			dst = dst[1:]
		}
		length++
		if c&0x80 == 0 {
			break
		}
	}

	return int(length)
}

func encodeInt64(v int64, dst []byte) int {
	var c uint8
	var s uint8

	length := uint8(0)
	for {
		c = uint8(v & 0x7f)
		s = uint8(v & 0x40)
		v >>= 7
		if (v != -1 || s == 0) && (v != 0 || s != 0) {
			c |= 0x80
		}
		if dst != nil {
			dst[0] = byte(c)
			dst = dst[1:]
		}
		length++
		if c&0x80 == 0 {
			break
		}
	}

	return int(length)
}
