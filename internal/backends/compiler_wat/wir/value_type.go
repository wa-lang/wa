// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

func toWatType(t ValueType) wat.ValueType {
	switch t.(type) {
	case RUNE:
		return wat.I32{}

	case I32:
		return wat.I32{}

	case U32:
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
		return wat.I32{}

	default:
		logger.Fatalf("Todo:%v\n", t)
	}

	return nil
}

/**************************************
Void:
**************************************/
type Void struct{}

func (t Void) byteSize() int          { return 0 }
func (t Void) Raw() []wat.ValueType   { return []wat.ValueType{} }
func (t Void) Equal(u ValueType) bool { _, ok := u.(Void); return ok }

/**************************************
RUNE:
**************************************/
type RUNE struct{}

func (t RUNE) byteSize() int          { return 4 }
func (t RUNE) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t RUNE) Equal(u ValueType) bool { _, ok := u.(RUNE); return ok }

/**************************************
I32:
**************************************/
type I32 struct{}

func (t I32) byteSize() int          { return 4 }
func (t I32) Raw() []wat.ValueType   { return []wat.ValueType{wat.I32{}} }
func (t I32) Equal(u ValueType) bool { _, ok := u.(I32); return ok }

/**************************************
U32:
**************************************/
type U32 struct{}

func (t U32) byteSize() int          { return 4 }
func (t U32) Raw() []wat.ValueType   { return []wat.ValueType{wat.U32{}} }
func (t U32) Equal(u ValueType) bool { _, ok := u.(U32); return ok }

/**************************************
I64:
**************************************/
type I64 struct{}

func (t I64) byteSize() int          { return 8 }
func (t I64) Raw() []wat.ValueType   { return []wat.ValueType{wat.I64{}} }
func (t I64) Equal(u ValueType) bool { _, ok := u.(I64); return ok }

/**************************************
Uint64:
**************************************/
type U64 struct{}

func (t U64) byteSize() int          { return 8 }
func (t U64) Raw() []wat.ValueType   { return []wat.ValueType{wat.U64{}} }
func (t U64) Equal(u ValueType) bool { _, ok := u.(U64); return ok }

/**************************************
F32:
**************************************/
type F32 struct{}

func (t F32) byteSize() int          { return 4 }
func (t F32) Raw() []wat.ValueType   { return []wat.ValueType{wat.F32{}} }
func (t F32) Equal(u ValueType) bool { _, ok := u.(F32); return ok }

/**************************************
F64:
**************************************/
type F64 struct{}

func (t F64) byteSize() int          { return 8 }
func (t F64) Raw() []wat.ValueType   { return []wat.ValueType{wat.F64{}} }
func (t F64) Equal(u ValueType) bool { _, ok := u.(F64); return ok }

/**************************************
Pointer:
**************************************/
type Pointer struct {
	Base ValueType
}

func NewPointer(base ValueType) Pointer { return Pointer{Base: base} }
func (t Pointer) byteSize() int         { return 4 }
func (t Pointer) Raw() []wat.ValueType  { return []wat.ValueType{wat.I32{}} }
func (t Pointer) Equal(u ValueType) bool {
	if ut, ok := u.(Pointer); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

/**************************************
Block:
**************************************/
type Block struct {
	Base ValueType
}

func NewBlock(base ValueType) Block  { return Block{Base: base} }
func (t Block) byteSize() int        { return 4 }
func (t Block) Raw() []wat.ValueType { return []wat.ValueType{wat.I32{}} }
func (t Block) Equal(u ValueType) bool {
	if ut, ok := u.(Block); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

/**************************************
Struct:
**************************************/
type Struct struct {
	name    string
	Members []Field
}

type Field struct {
	name string
	typ  ValueType
}

func NewField(n string, t ValueType) Field { return Field{name: n, typ: t} }
func (i Field) Name() string               { return i.name }
func (i Field) Type() ValueType            { return i.typ }
func (i Field) Equal(u Field) bool         { return i.name == u.name && i.typ.Equal(u.typ) }

func NewStruct(name string, m []Field) Struct {
	return Struct{name: name, Members: m}
}
func (t Struct) byteSize() int { logger.Fatal("Todo"); return 0 }
func (t Struct) Raw() []wat.ValueType {
	var r []wat.ValueType
	for _, f := range t.Members {
		r = append(r, f.Type().Raw()...)
	}
	return r
}
func (t Struct) Equal(u ValueType) bool {
	if u, ok := u.(Struct); ok {
		if len(t.Members) != len(u.Members) {
			return false
		}

		for i := range t.Members {
			if !t.Members[i].Equal(u.Members[i]) {
				return false
			}
		}

		return true
	}
	return false
}

/**************************************
Ref:
**************************************/
type Ref struct {
	Base       ValueType
	underlying Struct
}

func NewRef(base ValueType) Ref {
	var v Ref
	v.Base = base
	var m []Field
	m = append(m, NewField("block", NewBlock(base)))
	m = append(m, NewField("data", NewPointer(base)))
	v.underlying = NewStruct("", m)
	return v
}
func (t Ref) byteSize() int        { return 8 }
func (t Ref) Raw() []wat.ValueType { return t.underlying.Raw() }
func (t Ref) Equal(u ValueType) bool {
	if ut, ok := u.(Ref); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}
