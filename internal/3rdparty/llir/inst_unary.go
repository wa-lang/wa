package llir

import (
	"fmt"
	"strings"

	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// --- [ Unary instructions ] --------------------------------------------------

// ~~~ [ fneg ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstFNeg is an LLVM IR fneg instruction.
type InstFNeg struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Operand.
	X value.Value // floating-point scalar or floating-point vector

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type
	// (optional) Fast math flags.
	FastMathFlags []enum.FastMathFlag
	// (optional) Metadata.
	Metadata
}

// NewFNeg returns a new fneg instruction based on the given operand.
func NewFNeg(x value.Value) *InstFNeg {
	inst := &InstFNeg{X: x}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstFNeg) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstFNeg) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = inst.X.Type()
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'fneg' FastMathFlags=FastMathFlag* X=TypeValue Metadata=(',' MetadataAttachment)+?
func (inst *InstFNeg) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("fneg")
	for _, flag := range inst.FastMathFlags {
		fmt.Fprintf(buf, " %s", flag)
	}
	fmt.Fprintf(buf, " %s", inst.X)
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}
