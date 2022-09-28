// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

/**************************************
varBlock:
**************************************/
type varBlock struct {
	aVar
}

func newVarBlock(name string, kind ValueKind, base_type ValueType) *varBlock {
	return &varBlock{aVar: aVar{name: name, kind: kind, typ: NewBlock(base_type)}}
}
func (v *varBlock) raw() []wat.Value { return []wat.Value{wat.NewVarI32(v.name)} }
func (v *varBlock) EmitInit() []wat.Inst {
	return []wat.Inst{wat.NewInstConst(wat.I32{}, "0"), v.pop(v.name)}
}
func (v *varBlock) EmitPush() (insts []wat.Inst) {
	insts = append(insts, v.push(v.name))
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Retain"))
	return
}
func (v *varBlock) EmitPop() (insts []wat.Inst) {
	insts = append(insts, v.EmitRelease()...)
	insts = append(insts, v.pop(v.name))
	return
}

func (v *varBlock) EmitRelease() (insts []wat.Inst) {
	insts = append(insts, v.push(v.name))
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Release"))
	return
}

func (v *varBlock) emitLoad(addr Value) (insts []wat.Inst) {
	insts = append(insts, addr.EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.I32{}, 0, 1))
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Retain"))
	return
}

func (v *varBlock) emitStore(addr Value) (insts []wat.Inst) {
	insts = append(insts, v.push(v.name))
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Retain"))
	insts = append(insts, wat.NewInstDrop())

	insts = append(insts, addr.EmitPush()...)
	insts = append(insts, wat.NewInstLoad(wat.I32{}, 0, 1))
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Release"))

	insts = append(insts, addr.EmitPush()...)
	insts = append(insts, v.push(v.name))
	insts = append(insts, wat.NewInstStore(toWatType(v.Type()), 0, 1))
	return
}

func (v *varBlock) emitHeapAlloc(item_count Value) (insts []wat.Inst) {
	switch item_count.Kind() {
	case ValueKindConst:
		c, err := strconv.Atoi(item_count.Name())
		if err != nil {
			logger.Fatalf("%v\n", err)
			return nil
		}
		insts = append(insts, NewConst(I32{}, strconv.Itoa(v.Type().(Block).Base.byteSize()*c+16)).EmitPush()...)
		insts = append(insts, wat.NewInstCall("$waHeapAlloc"))

		insts = append(insts, item_count.EmitPush()...)           //item_count
		insts = append(insts, NewConst(I32{}, "0").EmitPush()...) //release_method
		insts = append(insts, wat.NewInstCall("$wa.RT.Block.Init"))

	default:
		if !item_count.Type().Equal(I32{}) {
			logger.Fatal("item_count should be i32")
			return nil
		}

		insts = append(insts, item_count.EmitPush()...)
		insts = append(insts, NewConst(I32{}, strconv.Itoa(v.Type().(Block).Base.byteSize())).EmitPush()...)
		insts = append(insts, wat.NewInstMul(wat.I32{}))
		insts = append(insts, NewConst(I32{}, "16").EmitPush()...)
		insts = append(insts, wat.NewInstAdd(wat.I32{}))
		insts = append(insts, wat.NewInstCall("$waHeapAlloc"))

		insts = append(insts, item_count.EmitPush()...)
		insts = append(insts, NewConst(I32{}, "0").EmitPush()...) //release_method
		insts = append(insts, wat.NewInstCall("$wa.RT.Block.Init"))
	}

	return
}

/**************************************
VarStruct:
**************************************/
type VarStruct struct {
	aVar
}

func newVarStruct(name string, kind ValueKind, typ ValueType) *VarStruct {
	return &VarStruct{aVar: aVar{name: name, kind: kind, typ: typ}}
}
func (v *VarStruct) raw() []wat.Value {
	var r []wat.Value
	st := v.Type().(Struct)
	for _, m := range st.Members {
		t := NewVar(v.Name()+"."+m.Name(), v.kind, m.Type())
		r = append(r, t.raw()...)
	}
	return r
}
func (v *VarStruct) EmitInit() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for _, m := range st.Members {
		t := NewVar(v.Name()+"."+m.Name(), v.kind, m.Type())
		insts = append(insts, t.EmitInit()...)
	}
	return insts
}
func (v *VarStruct) EmitPush() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for _, m := range st.Members {
		t := NewVar(v.Name()+"."+m.Name(), v.kind, m.Type())
		insts = append(insts, t.EmitPush()...)
	}
	return insts
}
func (v *VarStruct) EmitPop() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := NewVar(v.Name()+"."+m.Name(), v.kind, m.Type())
		insts = append(insts, t.EmitPop()...)
	}
	return insts
}
func (v *VarStruct) EmitRelease() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := NewVar(v.Name()+"."+m.Name(), v.kind, m.Type())
		insts = append(insts, t.EmitRelease()...)
	}
	return insts
}
func (v *VarStruct) Extract(member_name string) Value {
	st := v.Type().(Struct)
	for _, m := range st.Members {
		if m.Name() == member_name {
			return NewVar(v.Name()+"."+m.Name(), v.kind, m.Type())
		}
	}
	return nil
}
func (v *VarStruct) emitLoad(addr Value) []wat.Inst {
	logger.Fatal("Todo")
	return nil
}
func (v *VarStruct) emitStore(addr Value) []wat.Inst {
	logger.Fatal("Todo")
	return nil
}

/**************************************
VarRef:
**************************************/
type VarRef struct {
	aVar
	underlying VarStruct
}

func NewVarRef(name string, kind ValueKind, base_type ValueType) *VarRef {
	var v VarRef
	ref_type := NewRef(base_type)
	v.aVar = aVar{name: name, kind: kind, typ: ref_type}
	v.underlying = *newVarStruct(name, kind, ref_type.underlying)
	return &v
}

func (v *VarRef) raw() []wat.Value                { return v.underlying.raw() }
func (v *VarRef) EmitInit() []wat.Inst            { return v.underlying.EmitInit() }
func (v *VarRef) EmitPush() []wat.Inst            { return v.underlying.EmitPush() }
func (v *VarRef) EmitPop() []wat.Inst             { return v.underlying.EmitPop() }
func (v *VarRef) EmitRelease() []wat.Inst         { return v.underlying.EmitRelease() }
func (v *VarRef) emitLoad(addr Value) []wat.Inst  { return v.underlying.emitLoad(addr) }
func (v *VarRef) emitStore(addr Value) []wat.Inst { return v.underlying.emitStore(addr) }

func (v *VarRef) EmitLoad() []wat.Inst {
	t := NewVar("", v.kind, v.Type().(Ref).Base)
	return t.emitLoad(v.underlying.Extract("data"))
}

func (v *VarRef) EmitStore(d Value) []wat.Inst {
	if !d.Type().Equal(v.Type().(Ref).Base) {
		logger.Fatal("Type not match")
		return nil
	}
	return d.emitStore(v.underlying.Extract("data"))
}

func (v *VarRef) emitHeapAlloc() (insts []wat.Inst) {
	insts = append(insts, wat.NewBlank())
	insts = append(insts, wat.NewComment(v.name+" Ref.emitHeapAlloc start"))

	insts = append(insts, newVarBlock("", v.Kind(), v.Type().(Ref).Base).emitHeapAlloc(NewConst(I32{}, "1"))...)
	insts = append(insts, wat.NewInstCall("$wa.RT.DupWatStack"))
	insts = append(insts, NewConst(I32{}, "16").EmitPush()...)
	insts = append(insts, wat.NewInstAdd(wat.I32{}))

	insts = append(insts, wat.NewComment(v.name+" Ref.emitHeapAlloc end"))
	insts = append(insts, wat.NewBlank())

	return
}

func (v *VarRef) emitStackAlloc() (insts []wat.Inst) {
	insts = append(insts, wat.NewBlank())
	insts = append(insts, wat.NewComment(v.name+" Ref.emitStackAlloc start"))

	insts = append(insts, NewConst(I32{}, "0").EmitPush()...)
	insts = append(insts, NewConst(I32{}, strconv.Itoa(v.Type().(Ref).Base.byteSize())).EmitPush()...)
	insts = append(insts, wat.NewInstCall("$waStackAlloc"))

	insts = append(insts, wat.NewComment(v.name+" Ref.emitStackAlloc end"))
	insts = append(insts, wat.NewBlank())
	return
}
