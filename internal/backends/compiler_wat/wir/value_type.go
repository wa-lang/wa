// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
)

type TypeKind uint8

const (
	kUnknown TypeKind = iota
	kVoid
	kU8
	kU16
	kU32
	kU64
	kI8
	kI16
	kI32
	kI64
	kF32
	kF64
	kRune
	kPtr
	kBlock
	kStruct
	kTuple
	kSPtr
	kString
	kSlice
	kArray
	kMap
	kInterface
	kRef
)

func toWatType(t ValueType) wat.ValueType {
	switch t.(type) {
	case *tI32, *tRune, *tI8, *tI16:
		return wat.I32{}

	case *tU32, *tU8, *tU16:
		return wat.U32{}

	case *tI64:
		return wat.I64{}

	case *tU64:
		return wat.U64{}

	case *tF32:
		return wat.F32{}

	case *tF64:
		return wat.F64{}

	case *Ptr:
		return wat.U32{}

	case *Block:
		return wat.I32{}

	default:
		logger.Fatalf("Todo:%v\n", t)
	}

	return nil
}

/**************************************
tCommon:
**************************************/
type tCommon struct {
	addr    int
	hash    int
	methods []Method
}

func (t *tCommon) Hash() int           { return t.hash }
func (t *tCommon) SetHash(h int)       { t.hash = h }
func (t *tCommon) AddMethod(m Method)  { t.methods = append(t.methods, m) }
func (t *tCommon) NumMethods() int     { return len(t.methods) }
func (t *tCommon) Method(i int) Method { return t.methods[i] }
func (t *tCommon) typeInfoAddr() int   { return t.addr }

//func (t *tCommon) AddMethodEntry(m FnType) { logger.Fatal("Can't add method for common type.") }
//
///**************************************
//tUncommon:
//**************************************/
//type tUncommon struct {
//	tCommon
//	methodTab []FnType
//}
//
//func (t *tUncommon) AddMethodEntry(m FnType) { t.methodTab = append(t.methodTab, m) }

/**************************************
tVoid:
**************************************/
type tVoid struct {
	tCommon
}

func (t *tVoid) Name() string           { return "void" }
func (t *tVoid) Size() int              { return 0 }
func (t *tVoid) align() int             { return 0 }
func (t *tVoid) Kind() TypeKind         { return kVoid }
func (t *tVoid) onFree() int            { return 0 }
func (t *tVoid) Raw() []wat.ValueType   { return []wat.ValueType{} }
func (t *tVoid) Equal(u ValueType) bool { _, ok := u.(*tVoid); return ok }
func (t *tVoid) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	logger.Fatal("Unreachable")
	return nil
}

/**************************************
tRune:
**************************************/
type tRune struct {
	tCommon
}

func (t *tRune) Name() string           { return "rune" }
func (t *tRune) Size() int              { return 4 }
func (t *tRune) align() int             { return 4 }
func (t *tRune) Kind() TypeKind         { return kRune }
func (t *tRune) onFree() int            { return 0 }
func (t *tRune) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tRune) Equal(u ValueType) bool { _, ok := u.(*tRune); return ok }
func (t *tRune) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tI8:
**************************************/
type tI8 struct {
	tCommon
}

func (t *tI8) Name() string           { return "i8" }
func (t *tI8) Size() int              { return 1 }
func (t *tI8) align() int             { return 1 }
func (t *tI8) Kind() TypeKind         { return kI8 }
func (t *tI8) onFree() int            { return 0 }
func (t *tI8) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tI8) Equal(u ValueType) bool { _, ok := u.(*tI8); return ok }
func (t *tI8) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}

	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8s(offset, 1))
	return insts
}

/**************************************
tU8:
**************************************/
type tU8 struct {
	tCommon
}

func (t *tU8) Name() string           { return "u8" }
func (t *tU8) Size() int              { return 1 }
func (t *tU8) align() int             { return 1 }
func (t *tU8) Kind() TypeKind         { return kU8 }
func (t *tU8) onFree() int            { return 0 }
func (t *tU8) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tU8) Equal(u ValueType) bool { _, ok := u.(*tU8); return ok }
func (t *tU8) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8u(offset, 1))
	return insts
}

/**************************************
tI16:
**************************************/
type tI16 struct {
	tCommon
}

func (t *tI16) Name() string           { return "i16" }
func (t *tI16) Size() int              { return 2 }
func (t *tI16) align() int             { return 2 }
func (t *tI16) Kind() TypeKind         { return kI16 }
func (t *tI16) onFree() int            { return 0 }
func (t *tI16) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tI16) Equal(u ValueType) bool { _, ok := u.(*tI16); return ok }
func (t *tI16) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16s(offset, 1))
	return insts
}

/**************************************
tU16:
**************************************/
type tU16 struct {
	tCommon
}

func (t *tU16) Name() string           { return "u16" }
func (t *tU16) Size() int              { return 2 }
func (t *tU16) align() int             { return 2 }
func (t *tU16) Kind() TypeKind         { return kU16 }
func (t *tU16) onFree() int            { return 0 }
func (t *tU16) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tU16) Equal(u ValueType) bool { _, ok := u.(*tU16); return ok }
func (t *tU16) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16u(offset, 1))
	return insts
}

/**************************************
tI32:
**************************************/
type tI32 struct {
	tCommon
}

func (t *tI32) Name() string           { return "i32" }
func (t *tI32) Size() int              { return 4 }
func (t *tI32) align() int             { return 4 }
func (t *tI32) Kind() TypeKind         { return kI32 }
func (t *tI32) onFree() int            { return 0 }
func (t *tI32) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tI32) Equal(u ValueType) bool { _, ok := u.(*tI32); return ok }
func (t *tI32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tU32:
**************************************/
type tU32 struct {
	tCommon
}

func (t *tU32) Name() string           { return "u32" }
func (t *tU32) Size() int              { return 4 }
func (t *tU32) align() int             { return 4 }
func (t *tU32) Kind() TypeKind         { return kU32 }
func (t *tU32) onFree() int            { return 0 }
func (t *tU32) Raw() []wat.ValueType   { return []wat.ValueType{wat.U32{}} }
func (t *tU32) Equal(u ValueType) bool { _, ok := u.(*tU32); return ok }
func (t *tU32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tI64:
**************************************/
type tI64 struct {
	tCommon
}

func (t *tI64) Name() string           { return "i64" }
func (t *tI64) Size() int              { return 8 }
func (t *tI64) align() int             { return 8 }
func (t *tI64) Kind() TypeKind         { return kI64 }
func (t *tI64) onFree() int            { return 0 }
func (t *tI64) Raw() []wat.ValueType   { return []wat.ValueType{wat.I64{}} }
func (t *tI64) Equal(u ValueType) bool { _, ok := u.(*tI64); return ok }
func (t *tI64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tUint64:
**************************************/
type tU64 struct {
	tCommon
}

func (t *tU64) Name() string           { return "u64" }
func (t *tU64) Size() int              { return 8 }
func (t *tU64) align() int             { return 8 }
func (t *tU64) Kind() TypeKind         { return kU64 }
func (t *tU64) onFree() int            { return 0 }
func (t *tU64) Raw() []wat.ValueType   { return []wat.ValueType{wat.U64{}} }
func (t *tU64) Equal(u ValueType) bool { _, ok := u.(*tU64); return ok }
func (t *tU64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tF32:
**************************************/
type tF32 struct {
	tCommon
}

func (t *tF32) Name() string           { return "f32" }
func (t *tF32) Size() int              { return 4 }
func (t *tF32) align() int             { return 4 }
func (t *tF32) Kind() TypeKind         { return kF32 }
func (t *tF32) onFree() int            { return 0 }
func (t *tF32) Raw() []wat.ValueType   { return []wat.ValueType{wat.F32{}} }
func (t *tF32) Equal(u ValueType) bool { _, ok := u.(*tF32); return ok }
func (t *tF32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tF64:
**************************************/
type tF64 struct {
	tCommon
}

func (t *tF64) Name() string           { return "f64" }
func (t *tF64) Size() int              { return 8 }
func (t *tF64) align() int             { return 8 }
func (t *tF64) Kind() TypeKind         { return kF64 }
func (t *tF64) onFree() int            { return 0 }
func (t *tF64) Raw() []wat.ValueType   { return []wat.ValueType{wat.F64{}} }
func (t *tF64) Equal(u ValueType) bool { _, ok := u.(*tF64); return ok }
func (t *tF64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	//if !addr.Type().(*Ptr).Base.Equal(t) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}
	if _, ok := addr.(*aPtr); !ok {
		logger.Fatal("addr should be `*aPtr`")
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}
