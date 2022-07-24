package types

import "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"

var (
	I8      = &Int{Type: lltypes.I8, TypeName: "int8", TypeSize: 8 / 8, Signed: true}
	U8      = &Int{Type: lltypes.I8, TypeName: "uint8", TypeSize: 8 / 8}
	I16     = &Int{Type: lltypes.I16, TypeName: "int16", TypeSize: 18 / 8, Signed: true}
	U16     = &Int{Type: lltypes.I16, TypeName: "uint16", TypeSize: 18 / 8}
	I32     = &Int{Type: lltypes.I32, TypeName: "int32", TypeSize: 32 / 8, Signed: true}
	U32     = &Int{Type: lltypes.I32, TypeName: "uint32", TypeSize: 32 / 8}
	I64     = &Int{Type: lltypes.I64, TypeName: "int64", TypeSize: 64 / 8, Signed: true}
	U64     = &Int{Type: lltypes.I64, TypeName: "uint64", TypeSize: 64 / 8}
	Uintptr = &Int{Type: lltypes.I64, TypeName: "uintptr", TypeSize: 64 / 8}

	Void   = &VoidType{}
	Bool   = &BoolType{}
	String = &StringType{}
)
