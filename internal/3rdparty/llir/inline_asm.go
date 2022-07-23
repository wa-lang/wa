package llir

import (
	"fmt"
	"strings"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// InlineAsm is an inline assembler expression.
type InlineAsm struct {
	// Assembly instructions.
	Asm string
	// Constraints.
	Constraint string

	// extra.

	// Type of result produced by the inline assembler expression.
	Typ types.Type
	// (optional) Side effect.
	SideEffect bool
	// (optional) Stack alignment.
	AlignStack bool
	// (optional) Intel dialect.
	IntelDialect bool
}

// NewInlineAsm returns a new inline assembler expression based on the given
// type, assembly instructions and constraints.
func NewInlineAsm(typ types.Type, asm, constraint string) *InlineAsm {
	return &InlineAsm{Typ: typ, Asm: asm, Constraint: constraint}
}

// String returns the LLVM syntax representation of the inline assembler
// expression as a type-value pair.
func (asm *InlineAsm) String() string {
	return fmt.Sprintf("%s %s", asm.Type(), asm.Ident())
}

// Type returns the type of the inline assembler expression.
func (asm *InlineAsm) Type() types.Type {
	return asm.Typ
}

// Ident returns the identifier associated with the inline assembler expression.
func (asm *InlineAsm) Ident() string {
	// "asm" OptSideEffect OptAlignStack OptIntelDialect StringLit "," StringLit
	buf := &strings.Builder{}
	buf.WriteString("asm")
	if asm.SideEffect {
		buf.WriteString(" sideeffect")
	}
	if asm.AlignStack {
		buf.WriteString(" alignstack")
	}
	if asm.IntelDialect {
		buf.WriteString(" inteldialect")
	}
	fmt.Fprintf(buf, " %s, %s", quote(asm.Asm), quote(asm.Constraint))
	return buf.String()
}
