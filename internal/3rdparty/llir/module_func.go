package llir

import types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"

// --- [ Functions ] -----------------------------------------------------------

// NewFunc appends a new function to the module based on the given function
// name, return type and function parameters.
//
// The Parent field of the function is set to m.
func (m *Module) NewFunc(name string, retType types.Type, params ...*Param) *Func {
	f := NewFunc(name, retType, params...)
	f.Parent = m
	m.Funcs = append(m.Funcs, f)
	return f
}
