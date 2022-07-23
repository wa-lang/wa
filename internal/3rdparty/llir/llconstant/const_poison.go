package llconstant

import (
	"fmt"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ Poison values ] -------------------------------------------------------

// Poison is an LLVM IR poison value.
type Poison struct {
	// Poison value type.
	Typ types.Type
}

// NewPoison returns a new poison value based on the given type.
func NewPoison(typ types.Type) *Poison {
	return &Poison{Typ: typ}
}

// String returns the LLVM syntax representation of the constant as a type-value
// pair.
func (c *Poison) String() string {
	return fmt.Sprintf("%s %s", c.Type(), c.Ident())
}

// Type returns the type of the constant.
func (c *Poison) Type() types.Type {
	return c.Typ
}

// Ident returns the identifier associated with the constant.
func (*Poison) Ident() string {
	// 'poison'
	return "poison"
}
