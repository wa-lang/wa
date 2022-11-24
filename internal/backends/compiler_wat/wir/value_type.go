// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

func toWatType(t ValueType) wat.ValueType {
	switch t.(type) {
	case I32, RUNE, I8, I16:
		return wat.I32{}

	case U32, U8, U16:
		return wat.U32{}

	case I64:
		return wat.I64{}

	case U64:
		return wat.U64{}

	case F32:
		return wat.F32{}

	case F64:
		return wat.F64{}

	case Pointer:
		return wat.U32{}

	case Block:
		return wat.I32{}

	default:
		logger.Fatalf("Todo:%v\n", t)
	}

	return nil
}

/**************************************
VOID:
**************************************/
type VOID struct{}

func (t VOID) Name() string           { return "void" }
func (t VOID) size() int              { return 0 }
func (t VOID) align() int             { return 0 }
func (t VOID) onFree() int            { return 0 }
func (t VOID) Raw() []wat.ValueType   { return []wat.ValueType{} }
func (t VOID) Equal(u ValueType) bool { _, ok := u.(VOID); return ok }
func (t VOID) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	logger.Fatal("Unreachable")
	return nil
}

/**************************************
RUNE:
**************************************/
type RUNE struct{}

func (t RUNE) Name() string           { return "rune" }
func (t RUNE) size() int              { return 4 }
func (t RUNE) align() int             { return 4 }
func (t RUNE) onFree() int            { return 0 }
func (t RUNE) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t RUNE) Equal(u ValueType) bool { _, ok := u.(RUNE); return ok }
func (t RUNE) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
I8:
**************************************/
type I8 struct{}

func (t I8) Name() string           { return "i8" }
func (t I8) size() int              { return 1 }
func (t I8) align() int             { return 1 }
func (t I8) onFree() int            { return 0 }
func (t I8) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t I8) Equal(u ValueType) bool { _, ok := u.(I8); return ok }
func (t I8) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8s(offset, 1))
	return insts
}

/**************************************
U8:
**************************************/
type U8 struct{}

func (t U8) Name() string           { return "u8" }
func (t U8) size() int              { return 1 }
func (t U8) align() int             { return 1 }
func (t U8) onFree() int            { return 0 }
func (t U8) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t U8) Equal(u ValueType) bool { _, ok := u.(U8); return ok }
func (t U8) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad8u(offset, 1))
	return insts
}

/**************************************
I16:
**************************************/
type I16 struct{}

func (t I16) Name() string           { return "i16" }
func (t I16) size() int              { return 2 }
func (t I16) align() int             { return 2 }
func (t I16) onFree() int            { return 0 }
func (t I16) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t I16) Equal(u ValueType) bool { _, ok := u.(I16); return ok }
func (t I16) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16s(offset, 1))
	return insts
}

/**************************************
U16:
**************************************/
type U16 struct{}

func (t U16) Name() string           { return "u16" }
func (t U16) size() int              { return 2 }
func (t U16) align() int             { return 2 }
func (t U16) onFree() int            { return 0 }
func (t U16) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t U16) Equal(u ValueType) bool { _, ok := u.(U16); return ok }
func (t U16) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad16u(offset, 1))
	return insts
}

/**************************************
I32:
**************************************/
type I32 struct{}

func (t I32) Name() string           { return "i32" }
func (t I32) size() int              { return 4 }
func (t I32) align() int             { return 4 }
func (t I32) onFree() int            { return 0 }
func (t I32) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t I32) Equal(u ValueType) bool { _, ok := u.(I32); return ok }
func (t I32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
U32:
**************************************/
type U32 struct{}

func (t U32) Name() string           { return "u32" }
func (t U32) size() int              { return 4 }
func (t U32) align() int             { return 4 }
func (t U32) onFree() int            { return 0 }
func (t U32) Raw() []wat.ValueType   { return []wat.ValueType{wat.U32{}} }
func (t U32) Equal(u ValueType) bool { _, ok := u.(U32); return ok }
func (t U32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
I64:
**************************************/
type I64 struct{}

func (t I64) Name() string           { return "i64" }
func (t I64) size() int              { return 8 }
func (t I64) align() int             { return 8 }
func (t I64) onFree() int            { return 0 }
func (t I64) Raw() []wat.ValueType   { return []wat.ValueType{wat.I64{}} }
func (t I64) Equal(u ValueType) bool { _, ok := u.(I64); return ok }
func (t I64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
Uint64:
**************************************/
type U64 struct{}

func (t U64) Name() string           { return "u64" }
func (t U64) size() int              { return 8 }
func (t U64) align() int             { return 8 }
func (t U64) onFree() int            { return 0 }
func (t U64) Raw() []wat.ValueType   { return []wat.ValueType{wat.U64{}} }
func (t U64) Equal(u ValueType) bool { _, ok := u.(U64); return ok }
func (t U64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
F32:
**************************************/
type F32 struct{}

func (t F32) Name() string           { return "f32" }
func (t F32) size() int              { return 4 }
func (t F32) align() int             { return 4 }
func (t F32) onFree() int            { return 0 }
func (t F32) Raw() []wat.ValueType   { return []wat.ValueType{wat.F32{}} }
func (t F32) Equal(u ValueType) bool { _, ok := u.(F32); return ok }
func (t F32) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
F64:
**************************************/
type F64 struct{}

func (t F64) Name() string           { return "f64" }
func (t F64) size() int              { return 8 }
func (t F64) align() int             { return 8 }
func (t F64) onFree() int            { return 0 }
func (t F64) Raw() []wat.ValueType   { return []wat.ValueType{wat.F64{}} }
func (t F64) Equal(u ValueType) bool { _, ok := u.(F64); return ok }
func (t F64) EmitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}
