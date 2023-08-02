// 版权 @2023 凹语言 作者。保留所有权利。

package wasm

import "io"

const (
	I32 = 0x7F
	I64 = 0x7E
	F32 = 0x7D
	F64 = 0x7C
)

const (
	sectionCustom   = 0
	sectionType     = 1
	sectionImport   = 2
	sectionFunction = 3
	sectionTable    = 4
	sectionMemory   = 5
	sectionGlobal   = 6
	sectionExport   = 7
	sectionStart    = 8
	sectionElement  = 9
	sectionCode     = 10
	sectionData     = 11
)

type wasmFunc struct {
	Name string
	Type uint32
	Code []byte
}

type wasmFuncType struct {
	Params  []byte
	Results []byte
}

type nameWriter interface {
	io.ByteWriter
	io.Writer
}

func writeI32Const(w io.ByteWriter, v int32) {
	w.WriteByte(0x41) // i32.const
	writeSleb128(w, int64(v))
}

func writeI64Const(w io.ByteWriter, v int64) {
	w.WriteByte(0x42) // i64.const
	writeSleb128(w, v)
}

func writeName(w nameWriter, name string) {
	writeUleb128(w, uint64(len(name)))
	w.Write([]byte(name))
}

func writeUleb128(w io.ByteWriter, v uint64) {
	if v < 128 {
		w.WriteByte(uint8(v))
		return
	}
	more := true
	for more {
		c := uint8(v & 0x7f)
		v >>= 7
		more = v != 0
		if more {
			c |= 0x80
		}
		w.WriteByte(c)
	}
}

func writeUleb128FixedLength(w io.ByteWriter, v uint64, length int) {
	for i := 0; i < length; i++ {
		c := uint8(v & 0x7f)
		v >>= 7
		if i < length-1 {
			c |= 0x80
		}
		w.WriteByte(c)
	}
	if v != 0 {
		panic("writeUleb128FixedLength: length too small")
	}
}

func writeSleb128(w io.ByteWriter, v int64) {
	more := true
	for more {
		c := uint8(v & 0x7f)
		s := uint8(v & 0x40)
		v >>= 7
		more = !((v == 0 && s == 0) || (v == -1 && s != 0))
		if more {
			c |= 0x80
		}
		w.WriteByte(c)
	}
}
