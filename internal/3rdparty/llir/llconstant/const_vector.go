package llconstant

import (
	"fmt"
	"strings"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ Vector constants ] ----------------------------------------------------

// Vector is an LLVM IR vector constant.
type Vector struct {
	// Vector type.
	Typ *types.VectorType
	// Vector elements.
	Elems []Constant
}

// NewVector returns a new vector constant based on the given vector type and
// elements. The vector type is infered from the type of the elements if t is
// nil.
func NewVector(t *types.VectorType, elems ...Constant) *Vector {
	c := &Vector{
		Elems: elems,
		Typ:   t,
	}
	// Compute type.
	c.Type()
	return c
}

// String returns the LLVM syntax representation of the constant as a type-value
// pair.
func (c *Vector) String() string {
	return fmt.Sprintf("%s %s", c.Type(), c.Ident())
}

// Type returns the type of the constant.
func (c *Vector) Type() types.Type {
	// Cache type if not present.
	if c.Typ == nil {
		elemType := c.Elems[0].Type()
		c.Typ = types.NewVector(uint64(len(c.Elems)), elemType)
	}
	return c.Typ
}

// Ident returns the identifier associated with the constant.
func (c *Vector) Ident() string {
	// '<' Elems=(TypeConst separator ',')* '>'
	buf := &strings.Builder{}
	buf.WriteString("<")
	for i, elem := range c.Elems {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(elem.String())
	}
	buf.WriteString(">")
	return buf.String()
}
