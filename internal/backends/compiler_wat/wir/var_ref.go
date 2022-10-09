// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

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

func (v *VarRef) raw() []wat.Value                       { return v.underlying.raw() }
func (v *VarRef) EmitInit() []wat.Inst                   { return v.underlying.EmitInit() }
func (v *VarRef) EmitPush() []wat.Inst                   { return v.underlying.EmitPush() }
func (v *VarRef) EmitPop() []wat.Inst                    { return v.underlying.EmitPop() }
func (v *VarRef) EmitRelease() []wat.Inst                { return v.underlying.EmitRelease() }
func (v *VarRef) emitLoadFromAddr(addr Value) []wat.Inst { return v.underlying.emitLoadFromAddr(addr) }
func (v *VarRef) emitStoreToAddr(addr Value) []wat.Inst  { return v.underlying.emitStoreToAddr(addr) }

func (v *VarRef) emitGetValue() []wat.Inst {
	t := NewVar("", v.kind, v.Type().(Ref).Base)
	return t.emitLoadFromAddr(v.underlying.Extract("data"))
}

func (v *VarRef) emitSetValue(d Value) []wat.Inst {
	if !d.Type().Equal(v.Type().(Ref).Base) {
		logger.Fatal("Type not match")
		return nil
	}
	return d.emitStoreToAddr(v.underlying.Extract("data"))
}

func (v *VarRef) emitHeapAlloc(module *Module) (insts []wat.Inst) {
	insts = append(insts, wat.NewBlank())
	insts = append(insts, wat.NewComment(v.name+" Ref.emitHeapAlloc start"))

	insts = append(insts, newVarBlock("", v.Kind(), v.Type().(Ref).Base).emitHeapAlloc(NewConst(I32{}, "1"), module)...)
	insts = append(insts, wat.NewInstCall("$wa.RT.DupWatStack"))
	insts = append(insts, NewConst(I32{}, "16").EmitPush()...)
	insts = append(insts, wat.NewInstAdd(wat.I32{}))

	insts = append(insts, wat.NewComment(v.name+" Ref.emitHeapAlloc end"))
	insts = append(insts, wat.NewBlank())

	return
}

func (v *VarRef) emitStackAlloc(module *Module) (insts []wat.Inst) {
	insts = append(insts, wat.NewBlank())
	insts = append(insts, wat.NewComment(v.name+" Ref.emitStackAlloc start"))

	insts = append(insts, NewConst(I32{}, "0").EmitPush()...)
	insts = append(insts, NewConst(I32{}, strconv.Itoa(v.Type().(Ref).Base.size())).EmitPush()...)
	insts = append(insts, wat.NewInstCall("$waStackAlloc"))

	insts = append(insts, wat.NewComment(v.name+" Ref.emitStackAlloc end"))
	insts = append(insts, wat.NewBlank())
	return
}
