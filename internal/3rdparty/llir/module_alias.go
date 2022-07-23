package llir

import constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"

// NewAlias appends a new alias to the module based on the given alias name and
// aliasee.
func (m *Module) NewAlias(name string, aliasee constant.Constant) *Alias {
	alias := NewAlias(name, aliasee)
	m.Aliases = append(m.Aliases, alias)
	return alias
}
