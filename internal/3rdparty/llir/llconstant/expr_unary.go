package llconstant

import (
	"fmt"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ Unary expressions ] ---------------------------------------------------

// ~~~ [ fneg ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprFNeg is an LLVM IR fneg expression.
type ExprFNeg struct {
	// Operand.
	X Constant // floating-point scalar or vector constant

	// extra.

	// Type of result produced by the constant expression.
	Typ types.Type
}

// NewFNeg returns a new fneg expression based on the given operand.
func NewFNeg(x Constant) *ExprFNeg {
	e := &ExprFNeg{X: x}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprFNeg) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprFNeg) Type() types.Type {
	// Cache type if not present.
	if e.Typ == nil {
		e.Typ = e.X.Type()
	}
	return e.Typ
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprFNeg) Ident() string {
	// 'fneg' '(' X=TypeConst ')'
	return fmt.Sprintf("fneg (%s)", e.X)
}
