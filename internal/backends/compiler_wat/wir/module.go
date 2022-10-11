package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/ssa"
)

/**************************************
Module:
**************************************/
type Module struct {
	Funcs []*Function

	globals     []Value
	Globals_map map[ssa.Value]Value

	BaseWat string
}

func NewModule() *Module {
	var m Module
	m.Globals_map = make(map[ssa.Value]Value)
	return &m
}

func (m *Module) AddGlobal(name string, typ ValueType, is_pointer bool, ssa_value ssa.Value) Value {
	v := NewGlobal(name, typ, is_pointer)
	m.globals = append(m.globals, v)
	m.Globals_map[ssa_value] = v
	return v
}

func (m *Module) genGlobalAlloc() *Function {
	var f Function
	f.Name = "$waGlobalAlloc"
	f.Result = VOID{}

	for _, g := range m.globals {
		if g.Kind() != ValueKindGlobal_Pointer {
			continue
		}

		t := g.Type().(Pointer).Base
		f.Insts = append(f.Insts, wat.NewInstConst(wat.I32{}, strconv.Itoa(t.size())))
		f.Insts = append(f.Insts, wat.NewInstCall("$waHeapAlloc"))
		f.Insts = append(f.Insts, g.EmitPop()...)
	}

	return &f
}

func (m *Module) ToWatModule() *wat.Module {
	var wat_module wat.Module
	wat_module.BaseWat = m.BaseWat

	wat_module.Funcs = append(wat_module.Funcs, m.genGlobalAlloc().ToWatFunc())

	for _, f := range m.Funcs {
		wat_module.Funcs = append(wat_module.Funcs, f.ToWatFunc())
	}

	for _, g := range m.globals {
		raw := g.raw()
		for _, r := range raw {
			var wat_global wat.Global
			wat_global.V = r
			wat_global.IsMut = true
			wat_module.Globals = append(wat_module.Globals, wat_global)
		}
	}

	return &wat_module
}
