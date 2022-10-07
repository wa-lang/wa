package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/ssa"
)

/**************************************
Module:
**************************************/
type Module struct {
	Funcs []*Function

	Globals     []Value
	globals_map map[ssa.Value]Value

	BaseWat string
}

func NewModule() *Module {
	var m Module
	m.globals_map = make(map[ssa.Value]Value)
	return &m
}

func (m *Module) ToWatModule() *wat.Module {
	var wat_module wat.Module
	wat_module.BaseWat = m.BaseWat

	for _, f := range m.Funcs {
		wat_module.Funcs = append(wat_module.Funcs, f.ToWatFunc())
	}

	for _, g := range m.Globals {
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
