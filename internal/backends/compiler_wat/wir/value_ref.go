// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

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
func (t Ref) Name() string         { return "ref$" + t.Base.Name() }
func (t Ref) size() int            { return 8 }
func (t Ref) align() int           { return 4 }
func (t Ref) onFree(m *Module) int { return t.underlying.onFree(m) }
func (t Ref) Raw() []wat.ValueType { return t.underlying.Raw() }
func (t Ref) Equal(u ValueType) bool {
	if ut, ok := u.(Ref); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

func (t Ref) emitHeapAlloc(module *Module) (insts []wat.Inst) {
	insts = append(insts, wat.NewBlank())
	insts = append(insts, wat.NewComment("Ref.emitHeapAlloc start"))

	insts = append(insts, NewBlock(t.Base).emitHeapAlloc(NewConst("1", I32{}), module)...)
	insts = append(insts, wat.NewInstCall("$wa.RT.DupWatStack"))
	insts = append(insts, NewConst("16", I32{}).EmitPush()...)
	insts = append(insts, wat.NewInstAdd(wat.I32{}))

	insts = append(insts, wat.NewComment("Ref.emitHeapAlloc end"))
	insts = append(insts, wat.NewBlank())

	return
}

func (t Ref) emitStackAlloc(module *Module) (insts []wat.Inst) {
	insts = append(insts, wat.NewBlank())
	insts = append(insts, wat.NewComment("Ref.emitStackAlloc start"))

	logger.Fatal("Todo")

	insts = append(insts, NewConst("0", I32{}).EmitPush()...)
	insts = append(insts, NewConst(strconv.Itoa(t.Base.size()), I32{}).EmitPush()...)
	insts = append(insts, wat.NewInstCall("$waStackAlloc"))

	insts = append(insts, wat.NewComment("Ref.emitStackAlloc end"))
	insts = append(insts, wat.NewBlank())
	return
}

func (t Ref) emitLoadFromAddr(addr Value, offset int) []wat.Inst {
	return t.underlying.emitLoadFromAddr(addr, offset)
}

/**************************************
VarRef:
**************************************/
type VarRef struct {
	aVar
	underlying VarStruct
}

func newVarRef(name string, kind ValueKind, base_type ValueType) *VarRef {
	var v VarRef
	ref_type := NewRef(base_type)
	v.aVar = aVar{name: name, kind: kind, typ: ref_type}
	v.underlying = *newVarStruct(name, kind, ref_type.underlying)
	return &v
}

func (v *VarRef) raw() []wat.Value        { return v.underlying.raw() }
func (v *VarRef) EmitInit() []wat.Inst    { return v.underlying.EmitInit() }
func (v *VarRef) EmitPush() []wat.Inst    { return v.underlying.EmitPush() }
func (v *VarRef) EmitPop() []wat.Inst     { return v.underlying.EmitPop() }
func (v *VarRef) EmitRelease() []wat.Inst { return v.underlying.EmitRelease() }

func (v *VarRef) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	return v.underlying.emitStoreToAddr(addr, offset)
}

func (v *VarRef) emitGetValue() []wat.Inst {
	t := v.Type().(Ref).Base
	return t.emitLoadFromAddr(v.underlying.Extract("data"), 0)
}

func (v *VarRef) emitSetValue(d Value) []wat.Inst {
	if !d.Type().Equal(v.Type().(Ref).Base) {
		logger.Fatal("Type not match")
		return nil
	}
	return d.emitStoreToAddr(v.underlying.Extract("data"), 0)
}
