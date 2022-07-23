package llconstant

import (
	"fmt"

	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ Other expressions ] ---------------------------------------------------

// ~~~ [ icmp ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprICmp is an LLVM IR icmp expression.
type ExprICmp struct {
	// Integer comparison predicate.
	Pred enum.IPred
	// Integer scalar or vector operands.
	X, Y Constant

	// extra.

	// Type of result produced by the constant expression.
	Typ types.Type
}

// NewICmp returns a new icmp expression based on the given integer comparison
// predicate and integer scalar or vector operands.
func NewICmp(pred enum.IPred, x, y Constant) *ExprICmp {
	e := &ExprICmp{Pred: pred, X: x, Y: y}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprICmp) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprICmp) Type() types.Type {
	// Cache type if not present.
	if e.Typ == nil {
		switch xType := e.X.Type().(type) {
		case *types.IntType, *types.PointerType:
			e.Typ = types.I1
		case *types.VectorType:
			e.Typ = types.NewVector(xType.Len, types.I1)
		default:
			panic(fmt.Errorf("invalid icmp operand type; expected *types.IntType, *types.PointerType or *types.VectorType, got %T", xType))
		}
	}
	return e.Typ
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprICmp) Ident() string {
	// 'icmp' Pred=IPred '(' X=TypeConst ',' Y=TypeConst ')'
	return fmt.Sprintf("icmp %s (%s, %s)", e.Pred, e.X, e.Y)
}

// ~~~ [ fcmp ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprFCmp is an LLVM IR fcmp expression.
type ExprFCmp struct {
	// Floating-point comparison predicate.
	Pred enum.FPred
	// Floating-point scalar or vector operands.
	X, Y Constant

	// extra.

	// Type of result produced by the constant expression.
	Typ types.Type
}

// NewFCmp returns a new fcmp expression based on the given floating-point
// comparison predicate and floating-point scalar or vector operands.
func NewFCmp(pred enum.FPred, x, y Constant) *ExprFCmp {
	e := &ExprFCmp{Pred: pred, X: x, Y: y}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprFCmp) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprFCmp) Type() types.Type {
	// Cache type if not present.
	if e.Typ == nil {
		switch xType := e.X.Type().(type) {
		case *types.FloatType:
			e.Typ = types.I1
		case *types.VectorType:
			e.Typ = types.NewVector(xType.Len, types.I1)
		default:
			panic(fmt.Errorf("invalid fcmp operand type; expected *types.FloatType or *types.VectorType, got %T", xType))
		}
	}
	return e.Typ
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprFCmp) Ident() string {
	// 'fcmp' Pred=FPred '(' X=TypeConst ',' Y=TypeConst ')'
	return fmt.Sprintf("fcmp %s (%s, %s)", e.Pred, e.X, e.Y)
}

// ~~~ [ select ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprSelect is an LLVM IR select expression.
type ExprSelect struct {
	// Selection condition.
	Cond Constant
	// Operands.
	X, Y Constant

	// extra.

	// Type of result produced by the constant expression.
	Typ types.Type
}

// NewSelect returns a new select expression based on the given selection
// condition and operands.
func NewSelect(cond, x, y Constant) *ExprSelect {
	e := &ExprSelect{Cond: cond, X: x, Y: y}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprSelect) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprSelect) Type() types.Type {
	// Cache type if not present.
	if e.Typ == nil {
		e.Typ = e.X.Type()
	}
	return e.Typ
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprSelect) Ident() string {
	// 'select' '(' Cond=TypeConst ',' X=TypeConst ',' Y=TypeConst ')'
	return fmt.Sprintf("select (%s, %s, %s)", e.Cond, e.X, e.Y)
}
