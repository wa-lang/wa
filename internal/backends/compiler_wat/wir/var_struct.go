package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
)

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

func (v *VarStruct) emitLoadFromAddr(addr Value, offset int) (insts []wat.Inst) {
	st := v.Type().(Struct)
	for _, m := range st.Members {
		t := NewVar(v.Name()+"."+m.Name(), v.kind, m.Type())
		insts = append(insts, t.emitLoadFromAddr(addr, m._start+offset)...)
	}
	return
}

func (v *VarStruct) emitStoreToAddr(addr Value, offset int) (insts []wat.Inst) {
	st := v.Type().(Struct)
	for i := range st.Members {
		m := st.Members[len(st.Members)-i-1]
		t := NewVar(v.Name()+"."+m.Name(), v.kind, m.Type())
		insts = append(insts, t.emitStoreToAddr(addr, m._start+offset)...)
	}
	return
}
