// Package metadata provides access to LLVM IR metadata.
package llmetadata

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wa-lang/wa/internal/3rdparty/llir/internal/enc"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// TODO: remove Null if possible.

// Convenience literals.
var (
	// Null metadata literal.
	Null = &NullLit{}
)

// ~~~ [ Named metadata definition ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NamedDef is a named metadata definition.
type NamedDef struct {
	// Metadata definition name (without '!' prefix).
	Name string
	// Metadata definition nodes.
	Nodes []Node
}

// Ident returns the identifier associated with the named metadata definition.
func (md *NamedDef) Ident() string {
	return enc.MetadataName(md.Name)
}

// LLString returns the LLVM syntax representation of the named metadata
// definition.
//
// Name=MetadataName '=' '!' '{' MDNodes=(MetadataNode separator ',')* '}'
func (md *NamedDef) LLString() string {
	buf := &strings.Builder{}
	buf.WriteString("!{")
	for i, node := range md.Nodes {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(node.Ident())
	}
	buf.WriteString("}")
	return buf.String()
}

// === [ Metadata nodes and metadata strings ] =================================

// --- [ Metadata tuple ] ------------------------------------------------------

// Tuple is a metadata node tuple.
type Tuple struct {
	// Metadata ID associated with the metadata tuple; -1 if not present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	// Metadata tuple fields.
	Fields []Field
}

// String returns the LLVM syntax representation of the metadata tuple.
func (md *Tuple) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the metadata tuple.
func (md *Tuple) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the metadata tuple.
//
// '!' MDFields
func (md *Tuple) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	buf.WriteString("!{")
	for i, field := range md.Fields {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(field.String())
	}
	buf.WriteString("}")
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *Tuple) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// --- [ Metadata value ] ------------------------------------------------------

// A Value is a metadata value.
type Value struct {
	// Metadata value.
	Value Metadata
}

// String returns the LLVM syntax representation of the metadata value as a
// type-value pair.
func (md *Value) String() string {
	return fmt.Sprintf("%s %s", md.Type(), md.Ident())
}

// Type returns the type of the metadata value.
func (md *Value) Type() types.Type {
	return types.Metadata
}

// Ident returns the identifier associated with the metadata value.
func (md *Value) Ident() string {
	return md.Value.String()
}

// --- [ Metadata string ] -----------------------------------------------------

// String is a metadata string.
type String struct {
	// Metadata string value.
	Value string
}

// String returns the LLVM syntax representation of the metadata string.
func (md *String) String() string {
	// '!' Val=StringLit
	return fmt.Sprintf("!%s", quote(md.Value))
}

// --- [ Metadata attachment ] -------------------------------------------------

// Attachment is a metadata attachment.
type Attachment struct {
	// Metadata attachment name (without '!' prefix); e.g. !dbg.
	Name string
	// Metadata attachment node.
	Node MDNode
}

// String returns the string representation of the metadata attachment.
func (m *Attachment) String() string {
	// Name=MetadataName MDNode
	return fmt.Sprintf("%s %s", enc.MetadataName(m.Name), m.Node.Ident())
}

// --- [ Integer literals ] -----------------------------------------------------

// IntLit is an integer literal.
type IntLit int64

// String returns the LLVM syntax representation of the integer literal.
func (i IntLit) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// UintLit is an unsigned integer literal.
type UintLit uint64

// String returns the LLVM syntax representation of the unsigned integer literal.
func (i UintLit) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

// --- [ Null literal ] --------------------------------------------------------

// TODO: remove NullLit if possible.

// NullLit is a metadata null literal.
type NullLit struct{}

// String returns the LLVM syntax representation of the null literal.
func (i *NullLit) String() string {
	return "null"
}

// --- [ Metadata identifiers ] ------------------------------------------------

// ~~~ [ Metadata ID ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// MetadataID is a metadata ID, as used by metadata definitions.
type MetadataID int64

// Ident returns the identifier associated with the metadata ID.
func (i MetadataID) Ident() string {
	return enc.MetadataID(int64(i))
}

// ID returns the ID of the metadata ID.
func (i MetadataID) ID() int64 {
	return int64(i)
}

// SetID sets the ID of the metadata ID.
func (i *MetadataID) SetID(id int64) {
	*i = MetadataID(id)
}
