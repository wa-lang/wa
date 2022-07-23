package llconstant

import (
	"fmt"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ Token constants ] -----------------------------------------------------

// NoneToken is an LLVM IR none token constant.
type NoneToken struct {
}

// String returns the LLVM syntax representation of the constant as a type-value
// pair.
func (c *NoneToken) String() string {
	return fmt.Sprintf("%v %v", c.Type(), c.Ident())
}

// Type returns the type of the constant.
func (*NoneToken) Type() types.Type {
	return types.Token
}

// Ident returns the identifier associated with the constant.
func (*NoneToken) Ident() string {
	// 'none'
	return "none"
}
