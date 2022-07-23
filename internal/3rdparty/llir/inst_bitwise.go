package llir

import (
	"fmt"
	"strings"

	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// --- [ Bitwise instructions ] ------------------------------------------------

// ~~~ [ shl ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstShl is an LLVM IR shl instruction.
type InstShl struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Operands.
	X, Y value.Value // integer scalar or integer vector

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type
	// (optional) Overflow flags.
	OverflowFlags []enum.OverflowFlag
	// (optional) Metadata.
	Metadata
}

// NewShl returns a new shl instruction based on the given operands.
func NewShl(x, y value.Value) *InstShl {
	inst := &InstShl{X: x, Y: y}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstShl) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstShl) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = inst.X.Type()
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'shl' OverflowFlags=OverflowFlag* X=TypeValue ',' Y=Value Metadata=(',' MetadataAttachment)+?
func (inst *InstShl) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("shl")
	for _, flag := range inst.OverflowFlags {
		fmt.Fprintf(buf, " %s", flag)
	}
	fmt.Fprintf(buf, " %s, %s", inst.X, inst.Y.Ident())
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ lshr ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstLShr is an LLVM IR lshr instruction.
type InstLShr struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Operands.
	X, Y value.Value // integer scalars or vectors

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type
	// (optional) Exact.
	Exact bool
	// (optional) Metadata.
	Metadata
}

// NewLShr returns a new lshr instruction based on the given operands.
func NewLShr(x, y value.Value) *InstLShr {
	inst := &InstLShr{X: x, Y: y}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstLShr) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstLShr) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = inst.X.Type()
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'lshr' Exactopt X=TypeValue ',' Y=Value Metadata=(',' MetadataAttachment)+?
func (inst *InstLShr) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("lshr")
	if inst.Exact {
		buf.WriteString(" exact")
	}
	fmt.Fprintf(buf, " %s, %s", inst.X, inst.Y.Ident())
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ ashr ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstAShr is an LLVM IR ashr instruction.
type InstAShr struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Operands.
	X, Y value.Value // integer scalars or vectors

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type
	// (optional) Exact.
	Exact bool
	// (optional) Metadata.
	Metadata
}

// NewAShr returns a new ashr instruction based on the given operands.
func NewAShr(x, y value.Value) *InstAShr {
	inst := &InstAShr{X: x, Y: y}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstAShr) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstAShr) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = inst.X.Type()
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'ashr' Exactopt X=TypeValue ',' Y=Value Metadata=(',' MetadataAttachment)+?
func (inst *InstAShr) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("ashr")
	if inst.Exact {
		buf.WriteString(" exact")
	}
	fmt.Fprintf(buf, " %s, %s", inst.X, inst.Y.Ident())
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ and ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstAnd is an LLVM IR and instruction.
type InstAnd struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Operands.
	X, Y value.Value // integer scalars or vectors

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type
	// (optional) Metadata.
	Metadata
}

// NewAnd returns a new and instruction based on the given operands.
func NewAnd(x, y value.Value) *InstAnd {
	inst := &InstAnd{X: x, Y: y}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstAnd) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstAnd) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = inst.X.Type()
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'and' X=TypeValue ',' Y=Value Metadata=(',' MetadataAttachment)+?
func (inst *InstAnd) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "and %s, %s", inst.X, inst.Y.Ident())
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ or ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstOr is an LLVM IR or instruction.
type InstOr struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Operands.
	X, Y value.Value // integer scalars or vectors

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type
	// (optional) Metadata.
	Metadata
}

// NewOr returns a new or instruction based on the given operands.
func NewOr(x, y value.Value) *InstOr {
	inst := &InstOr{X: x, Y: y}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstOr) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstOr) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = inst.X.Type()
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'or' X=TypeValue ',' Y=Value Metadata=(',' MetadataAttachment)+?
func (inst *InstOr) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "or %s, %s", inst.X, inst.Y.Ident())
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ xor ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstXor is an LLVM IR xor instruction.
type InstXor struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Operands.
	X, Y value.Value // integer scalars or vectors

	// extra.

	// Type of result produced by the instruction.
	Typ types.Type
	// (optional) Metadata.
	Metadata
}

// NewXor returns a new xor instruction based on the given operands.
func NewXor(x, y value.Value) *InstXor {
	inst := &InstXor{X: x, Y: y}
	// Compute type.
	inst.Type()
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstXor) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstXor) Type() types.Type {
	// Cache type if not present.
	if inst.Typ == nil {
		inst.Typ = inst.X.Type()
	}
	return inst.Typ
}

// LLString returns the LLVM syntax representation of the instruction.
//
// 'xor' X=TypeValue ',' Y=Value Metadata=(',' MetadataAttachment)+?
func (inst *InstXor) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "xor %s, %s", inst.X, inst.Y.Ident())
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}
