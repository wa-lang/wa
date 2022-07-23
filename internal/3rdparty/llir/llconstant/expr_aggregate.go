package llconstant

import (
	"fmt"
	"strings"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ Aggregate expressions ] -----------------------------------------------

// ~~~ [ extractvalue ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprExtractValue is an LLVM IR extractvalue expression.
type ExprExtractValue struct {
	// Aggregate value.
	X Constant
	// Element indices.
	Indices []uint64

	// extra.

	// Type of result produced by the constant expression.
	Typ types.Type
}

// NewExtractValue returns a new extractvalue expression based on the given
// aggregate value and indicies.
func NewExtractValue(x Constant, indices ...uint64) *ExprExtractValue {
	e := &ExprExtractValue{X: x, Indices: indices}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprExtractValue) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprExtractValue) Type() types.Type {
	// Cache type if not present.
	if e.Typ == nil {
		e.Typ = aggregateElemType(e.X.Type(), e.Indices)
	}
	return e.Typ
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprExtractValue) Ident() string {
	// 'extractvalue' '(' X=TypeConst Indices=(',' UintLit)* ')'
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "extractvalue (%s", e.X)
	for _, index := range e.Indices {
		fmt.Fprintf(buf, ", %d", index)
	}
	buf.WriteString(")")
	return buf.String()
}

// ~~~ [ insertvalue ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExprInsertValue is an LLVM IR insertvalue expression.
type ExprInsertValue struct {
	// Aggregate value.
	X Constant
	// Element to insert.
	Elem Constant
	// Element indices.
	Indices []uint64

	// extra.

	// Type of result produced by the constant expression.
	Typ types.Type
}

// NewInsertValue returns a new insertvalue expression based on the given
// aggregate value, element and indicies.
func NewInsertValue(x, elem Constant, indices ...uint64) *ExprInsertValue {
	e := &ExprInsertValue{X: x, Elem: elem, Indices: indices}
	// Compute type.
	e.Type()
	return e
}

// String returns the LLVM syntax representation of the constant expression as a
// type-value pair.
func (e *ExprInsertValue) String() string {
	return fmt.Sprintf("%s %s", e.Type(), e.Ident())
}

// Type returns the type of the constant expression.
func (e *ExprInsertValue) Type() types.Type {
	// Cache type if not present.
	if e.Typ == nil {
		e.Typ = e.X.Type()
	}
	return e.Typ
}

// Ident returns the identifier associated with the constant expression.
func (e *ExprInsertValue) Ident() string {
	// 'insertvalue' '(' X=TypeConst ',' Elem=TypeConst Indices=(',' UintLit)*
	// ')'
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "insertvalue (%s, %s", e.X, e.Elem)
	for _, index := range e.Indices {
		fmt.Fprintf(buf, ", %d", index)
	}
	buf.WriteString(")")
	return buf.String()
}

// ### [ Helper functions ] ####################################################

// aggregateElemType returns the element type at the position in the aggregate
// type specified by the given indices.
func aggregateElemType(t types.Type, indices []uint64) types.Type {
	// Base case.
	if len(indices) == 0 {
		return t
	}
	switch t := t.(type) {
	case *types.ArrayType:
		return aggregateElemType(t.ElemType, indices[1:])
	case *types.StructType:
		return aggregateElemType(t.Fields[indices[0]], indices[1:])
	default:
		panic(fmt.Errorf("support for aggregate type %T not yet implemented", t))
	}
}
