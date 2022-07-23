package llir

import types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"

// --- [ Type definitions ] ----------------------------------------------------

// NewTypeDef appends a new type definition to the module based on the given
// type name and underlying type.
func (m *Module) NewTypeDef(name string, typ types.Type) types.Type {
	typ.SetName(name)
	m.TypeDefs = append(m.TypeDefs, typ)
	return typ
}
