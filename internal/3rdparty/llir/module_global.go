package llir

import (
	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ Global variables ] ----------------------------------------------------

// NewGlobal appends a new global variable declaration to the module based on
// the given global variable name and content type.
func (m *Module) NewGlobal(name string, contentType types.Type) *Global {
	g := NewGlobal(name, contentType)
	m.Globals = append(m.Globals, g)
	return g
}

// NewGlobalDef appends a new global variable definition to the module based on
// the given global variable name and initial value.
func (m *Module) NewGlobalDef(name string, init constant.Constant) *Global {
	g := NewGlobalDef(name, init)
	m.Globals = append(m.Globals, g)
	return g
}
