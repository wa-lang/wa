package llconstant

import (
	"fmt"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ Vector expressions ] --------------------------------------------------

// ~~~ [ extractelement ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprExtractElement is an LLVM IR extractelement expression.
type ExprExtractElement struct {
	// Vector.
	X Constant
	// Element index.
	Index Constant

	// extra.

	// Type of result produced by the constant expression.
	Typ types.Type
}

// NewExtractElement returns a new extractelement expression based on the given
// vector and element index.
func NewExtractElement(x, index Constant) *ExprExtractElement {
	e := &ExprExtractElement{X: x, Index: index}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprExtractElement) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprExtractElement) Type() types.Type {
	// Cache type if not present.
	if e.Typ == nil {
		t, ok := e.X.Type().(*types.VectorType)
		if !ok {
			panic(fmt.Errorf("invalid vector type; expected *types.VectorType, got %T", e.X.Type()))
		}
		e.Typ = t.ElemType
	}
	return e.Typ
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprExtractElement) Ident() string {
	// 'extractelement' '(' X=TypeConst ',' Index=TypeConst ')'
	return fmt.Sprintf("extractelement (%s, %s)", e.X, e.Index)
}

// ~~~ [ insertelement ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprInsertElement is an LLVM IR insertelement expression.
type ExprInsertElement struct {
	// Vector.
	X Constant
	// Element to insert.
	Elem Constant
	// Element index.
	Index Constant

	// extra.

	// Type of result produced by the constant expression.
	Typ types.Type
}

// NewInsertElement returns a new insertelement expression based on the given
// vector, element and element index.
func NewInsertElement(x, elem, index Constant) *ExprInsertElement {
	e := &ExprInsertElement{X: x, Elem: elem, Index: index}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprInsertElement) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprInsertElement) Type() types.Type {
	// Cache type if not present.
	if e.Typ == nil {
		t, ok := e.X.Type().(*types.VectorType)
		if !ok {
			panic(fmt.Errorf("invalid vector type; expected *types.VectorType, got %T", e.X.Type()))
		}
		e.Typ = t
	}
	return e.Typ
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprInsertElement) Ident() string {
	// 'insertelement' '(' X=TypeConst ',' Elem=TypeConst ',' Index=TypeConst ')'
	return fmt.Sprintf("insertelement (%s, %s, %s)", e.X, e.Elem, e.Index)
}

// ~~~ [ shufflevector ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprShuffleVector is an LLVM IR shufflevector expression.
type ExprShuffleVector struct {
	// Vectors.
	X, Y Constant
	// Shuffle mask.
	Mask Constant

	// extra.

	// Type of result produced by the constant expression.
	Typ types.Type
}

// NewShuffleVector returns a new shufflevector expression based on the given
// vectors and shuffle mask.
func NewShuffleVector(x, y, mask Constant) *ExprShuffleVector {
	e := &ExprShuffleVector{X: x, Y: y, Mask: mask}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprShuffleVector) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprShuffleVector) Type() types.Type {
	// Cache type if not present.
	if e.Typ == nil {
		xType, ok := e.X.Type().(*types.VectorType)
		if !ok {
			panic(fmt.Errorf("invalid vector type; expected *types.VectorType, got %T", e.X.Type()))
		}
		maskType, ok := e.Mask.Type().(*types.VectorType)
		if !ok {
			panic(fmt.Errorf("invalid vector type; expected *types.VectorType, got %T", e.Mask.Type()))
		}
		e.Typ = types.NewVector(maskType.Len, xType.ElemType)
	}
	return e.Typ
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprShuffleVector) Ident() string {
	// 'shufflevector' '(' X=TypeConst ',' Y=TypeConst ',' Mask=TypeConst ')'
	return fmt.Sprintf("shufflevector (%s, %s, %s)", e.X, e.Y, e.Mask)
}
