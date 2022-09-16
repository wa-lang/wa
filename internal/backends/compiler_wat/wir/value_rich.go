// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
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
	return []wat.Inst{wat.NewInstConst(wat.I32{}, "0"), v.set(v.name)}
}
func (v *varBlock) EmitGet() []wat.Inst { return []wat.Inst{v.get(v.name)} }
func (v *varBlock) EmitSet() []wat.Inst {
	var insts []wat.Inst
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Retain"))
	insts = append(insts, v.EmitRelease()...)
	insts = append(insts, v.set(v.name))
	return insts
}
func (v *varBlock) EmitRelease() []wat.Inst {
	insts := v.EmitGet()
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Release"))
	return insts
}
func (v *varBlock) emitLoad(addr Value) []wat.Inst {
	//if !addr.Type().(Pointer).Base.Equal(v.Type()) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}

	insts := addr.EmitGet()
	insts = append(insts, wat.NewInstLoad(wat.I32{}, 0, 1))
	return insts
}
func (v *varBlock) emitStore(addr Value) []wat.Inst {
	//if !addr.Type().(Pointer).Base.Equal(v.Type()) {
	//	logger.Fatal("Type not match")
	//	return nil
	//}

	insts := v.EmitGet()
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Retain"))
	insts = append(insts, wat.NewInstDrop())

	NewVar("", v.kind, v.Type()).emitLoad(addr)
	insts = append(insts, wat.NewInstCall("$wa.RT.Block.Release"))

	insts = append(insts, addr.EmitGet()...)
	insts = append(insts, v.EmitGet()...)
	insts = append(insts, wat.NewInstStore(toWatType(v.Type()), 0, 1))

	return insts
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
		t := NewVar(m.Name(), v.kind, m.Type())
		insts = append(insts, t.EmitInit()...)
	}
	return insts
}
func (v *VarStruct) EmitGet() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for _, m := range st.Members {
		t := NewVar(m.Name(), v.kind, m.Type())
		insts = append(insts, t.EmitGet()...)
	}
	return insts
}
func (v *VarStruct) EmitSet() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := NewVar(m.Name(), v.kind, m.Type())
		insts = append(insts, t.EmitSet()...)
	}
	return insts
}
func (v *VarStruct) EmitRelease() []wat.Inst {
	var insts []wat.Inst
	st := v.Type().(Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := NewVar(m.Name(), v.kind, m.Type())
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
	v.aVar = aVar{name: name, kind: kind, typ: NewRef(base_type)}
	v.underlying = *newVarStruct(name, kind, v.typ)
	return &v
}

func (v *VarRef) raw() []wat.Value        { return v.underlying.raw() }
func (v *VarRef) EmitInit() []wat.Inst    { return v.underlying.EmitInit() }
func (v *VarRef) EmitGet() []wat.Inst     { return v.underlying.EmitGet() }
func (v *VarRef) EmitSet() []wat.Inst     { return v.underlying.EmitSet() }
func (v *VarRef) EmitRelease() []wat.Inst { return v.underlying.EmitRelease() }
func (v *VarRef) emitLoad(addr Value) []wat.Inst {
	return v.underlying.emitLoad(addr)
}
func (v *VarRef) emitStore(addr Value) []wat.Inst {
	return v.underlying.emitStore(addr)
}
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
