package llconstant

import (
	"fmt"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ Undefined values ] ----------------------------------------------------

// Undef is an LLVM IR undefined value.
type Undef struct {
	// Undefined value type.
	Typ types.Type
}

// NewUndef returns a new undefined value based on the given type.
func NewUndef(typ types.Type) *Undef {
	return &Undef{Typ: typ}
}

// String returns the LLVM syntax representation of the constant as a type-value
// pair.
func (c *Undef) String() string {
	return fmt.Sprintf("%s %s", c.Type(), c.Ident())
}

// Type returns the type of the constant.
func (c *Undef) Type() types.Type {
	return c.Typ
}

// Ident returns the identifier associated with the constant.
func (*Undef) Ident() string {
	// 'undef'
	return "undef"
}
