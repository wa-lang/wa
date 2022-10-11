package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

/**************************************
Pointer:
**************************************/
type Pointer struct {
	Base ValueType
}

func NewPointer(base ValueType) Pointer { return Pointer{Base: base} }
func (t Pointer) Name() string          { return "pointer$" + t.Base.Name() }
func (t Pointer) size() int             { return 4 }
func (t Pointer) align() int            { return 4 }
func (t Pointer) onFree(m *Module) int  { return 0 }
func (t Pointer) Raw() []wat.ValueType  { return []wat.ValueType{wat.I32{}} }
func (t Pointer) Equal(u ValueType) bool {
	if ut, ok := u.(Pointer); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}
func (t Pointer) emitLoadFromAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(t) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, wat.NewInstLoad(toWatType(t), offset, 1))
	return insts
}

/**************************************
aPointer:
**************************************/
type aPointer struct {
	aBasic
}

func newValuePointer(name string, kind ValueKind, base_type ValueType) *aPointer {
	var v aPointer
	pointer_type := NewPointer(base_type)
	v.aValue = aValue{name: name, kind: kind, typ: pointer_type}
	return &v
}

func (v *aPointer) emitGetValue() []wat.Inst {
	t := v.Type().(Pointer).Base
	return t.emitLoadFromAddr(v, 0)
}

func (v *aPointer) emitSetValue(d Value) []wat.Inst {
	if !d.Type().Equal(v.Type().(Pointer).Base) {
		logger.Fatal("Type not match")
		return nil
	}
	return d.emitStoreToAddr(v, 0)
}
