package llconstant

import (
	"fmt"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ zeroinitializer constants ] -------------------------------------------

// ZeroInitializer is an LLVM IR zeroinitializer constant.
type ZeroInitializer struct {
	// zeroinitializer type.
	Typ types.Type
}

// NewZeroInitializer returns a new zeroinitializer constant based on the given
// type.
func NewZeroInitializer(typ types.Type) *ZeroInitializer {
	return &ZeroInitializer{Typ: typ}
}

// String returns the LLVM syntax representation of the constant as a type-value
// pair.
func (c *ZeroInitializer) String() string {
	return fmt.Sprintf("%s %s", c.Type(), c.Ident())
}

// Type returns the type of the constant.
func (c *ZeroInitializer) Type() types.Type {
	return c.Typ
}

// Ident returns the identifier associated with the constant.
func (c *ZeroInitializer) Ident() string {
	// 'zeroinitializer'
	return "zeroinitializer"
}
