// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/logger"
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
tVoid:
**************************************/
type tVoid struct{}

func (t *tVoid) Name() string           { return "void" }
func (t *tVoid) size() int              { return 0 }
func (t *tVoid) align() int             { return 0 }
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
type tRune struct{}

func (t *tRune) Name() string           { return "rune" }
func (t *tRune) size() int              { return 4 }
func (t *tRune) align() int             { return 4 }
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
type tI8 struct{}

func (t *tI8) Name() string           { return "i8" }
func (t *tI8) size() int              { return 1 }
func (t *tI8) align() int             { return 1 }
func (t *tI8) onFree() int            { return 0 }
func (t *tI8) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tI8) Equal(u ValueType) bool { _, ok := u.(*tI8); return ok }
func (t *tI8) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8s(offset, 1))
	return insts
}

/**************************************
tU8:
**************************************/
type tU8 struct{}

func (t *tU8) Name() string           { return "u8" }
func (t *tU8) size() int              { return 1 }
func (t *tU8) align() int             { return 1 }
func (t *tU8) onFree() int            { return 0 }
func (t *tU8) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tU8) Equal(u ValueType) bool { _, ok := u.(*tU8); return ok }
func (t *tU8) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8u(offset, 1))
	return insts
}

/**************************************
tI16:
**************************************/
type tI16 struct{}

func (t *tI16) Name() string           { return "i16" }
func (t *tI16) size() int              { return 2 }
func (t *tI16) align() int             { return 2 }
func (t *tI16) onFree() int            { return 0 }
func (t *tI16) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tI16) Equal(u ValueType) bool { _, ok := u.(*tI16); return ok }
func (t *tI16) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16s(offset, 1))
	return insts
}

/**************************************
tU16:
**************************************/
type tU16 struct{}

func (t *tU16) Name() string           { return "u16" }
func (t *tU16) size() int              { return 2 }
func (t *tU16) align() int             { return 2 }
func (t *tU16) onFree() int            { return 0 }
func (t *tU16) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tU16) Equal(u ValueType) bool { _, ok := u.(*tU16); return ok }
func (t *tU16) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16u(offset, 1))
	return insts
}

/**************************************
tI32:
**************************************/
type tI32 struct{}

func (t *tI32) Name() string           { return "i32" }
func (t *tI32) size() int              { return 4 }
func (t *tI32) align() int             { return 4 }
func (t *tI32) onFree() int            { return 0 }
func (t *tI32) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t *tI32) Equal(u ValueType) bool { _, ok := u.(*tI32); return ok }
func (t *tI32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tU32:
**************************************/
type tU32 struct{}

func (t *tU32) Name() string           { return "u32" }
func (t *tU32) size() int              { return 4 }
func (t *tU32) align() int             { return 4 }
func (t *tU32) onFree() int            { return 0 }
func (t *tU32) Raw() []wat.ValueType   { return []wat.ValueType{wat.U32{}} }
func (t *tU32) Equal(u ValueType) bool { _, ok := u.(*tU32); return ok }
func (t *tU32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tI64:
**************************************/
type tI64 struct{}

func (t *tI64) Name() string           { return "i64" }
func (t *tI64) size() int              { return 8 }
func (t *tI64) align() int             { return 8 }
func (t *tI64) onFree() int            { return 0 }
func (t *tI64) Raw() []wat.ValueType   { return []wat.ValueType{wat.I64{}} }
func (t *tI64) Equal(u ValueType) bool { _, ok := u.(*tI64); return ok }
func (t *tI64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tUint64:
**************************************/
type tU64 struct{}

func (t *tU64) Name() string           { return "u64" }
func (t *tU64) size() int              { return 8 }
func (t *tU64) align() int             { return 8 }
func (t *tU64) onFree() int            { return 0 }
func (t *tU64) Raw() []wat.ValueType   { return []wat.ValueType{wat.U64{}} }
func (t *tU64) Equal(u ValueType) bool { _, ok := u.(*tU64); return ok }
func (t *tU64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tF32:
**************************************/
type tF32 struct{}

func (t *tF32) Name() string           { return "f32" }
func (t *tF32) size() int              { return 4 }
func (t *tF32) align() int             { return 4 }
func (t *tF32) onFree() int            { return 0 }
func (t *tF32) Raw() []wat.ValueType   { return []wat.ValueType{wat.F32{}} }
func (t *tF32) Equal(u ValueType) bool { _, ok := u.(*tF32); return ok }
func (t *tF32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
tF64:
**************************************/
type tF64 struct{}

func (t *tF64) Name() string           { return "f64" }
func (t *tF64) size() int              { return 8 }
func (t *tF64) align() int             { return 8 }
func (t *tF64) onFree() int            { return 0 }
func (t *tF64) Raw() []wat.ValueType   { return []wat.ValueType{wat.F64{}} }
func (t *tF64) Equal(u ValueType) bool { _, ok := u.(*tF64); return ok }
func (t *tF64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(*Ptr).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}
