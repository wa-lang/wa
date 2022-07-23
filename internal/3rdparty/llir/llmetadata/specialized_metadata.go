package llmetadata

import (
	"fmt"
	"strings"

	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
)

// ~~~ [ DIBasicType ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIBasicType is a specialized metadata node.
type DIBasicType struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Tag      enum.DwarfTag         // optional; zero value if not present.
	Name     string                // optional; empty if not present.
	Size     uint64                // optional; zero value if not present.
	Align    uint64                // optional; zero value if not present.
	Encoding enum.DwarfAttEncoding // optional; zero value if not present.
	Flags    enum.DIFlag           // optional.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIBasicType) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIBasicType) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIBasicType' '(' Fields=(DIBasicTypeField separator ',')* ')'
func (md *DIBasicType) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	if md.Tag != 0 {
		field := fmt.Sprintf("tag: %s", dwarfTagString(md.Tag))
		fields = append(fields, field)
	}
	if len(md.Name) > 0 {
		field := fmt.Sprintf("name: %s", quote(md.Name))
		fields = append(fields, field)
	}
	if md.Size != 0 {
		field := fmt.Sprintf("size: %d", md.Size)
		fields = append(fields, field)
	}
	if md.Align != 0 {
		field := fmt.Sprintf("align: %d", md.Align)
		fields = append(fields, field)
	}
	if md.Encoding != 0 {
		field := fmt.Sprintf("encoding: %s", md.Encoding)
		fields = append(fields, field)
	}
	if md.Flags != 0 {
		field := fmt.Sprintf("flags: %s", diFlagsString(md.Flags))
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DIBasicType(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIBasicType) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DICommonBlock ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DICommonBlock is a specialized metadata node.
type DICommonBlock struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Scope       Field   // required.
	Declaration Field   // optional; nil if not present.
	Name        string  // optional; empty if not present.
	File        *DIFile // required.
	Line        int64   // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DICommonBlock) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DICommonBlock) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DICommonBlock' '(' Fields=(DICommonBlockField separator ',')* ')'
func (md *DICommonBlock) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("scope: %s", md.Scope)
	fields = append(fields, field)
	if md.Declaration != nil {
		field := fmt.Sprintf("declaration: %s", md.Declaration)
		fields = append(fields, field)
	}
	if len(md.Name) > 0 {
		field := fmt.Sprintf("name: %s", quote(md.Name))
		fields = append(fields, field)
	}
	if md.File != nil {
		field := fmt.Sprintf("file: %s", md.File)
		fields = append(fields, field)
	}
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DICommonBlock(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DICommonBlock) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DICompileUnit ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DICompileUnit is a specialized metadata node.
type DICompileUnit struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Language              enum.DwarfLang     // required.
	File                  *DIFile            // required.
	Producer              string             // optional; empty if not present.
	IsOptimized           bool               // optional; zero value if not present.
	Flags                 string             // optional; empty if not present.
	RuntimeVersion        uint64             // optional; zero value if not present.
	SplitDebugFilename    string             // optional; empty if not present.
	EmissionKind          enum.EmissionKind  // optional; zero value if not present.
	Enums                 *Tuple             // optional; nil if not present.
	RetainedTypes         *Tuple             // optional; nil if not present.
	Globals               *Tuple             // optional; nil if not present.
	Imports               *Tuple             // optional; nil if not present.
	Macros                *Tuple             // optional; nil if not present.
	DwoID                 uint64             // optional; zero value if not present.
	SplitDebugInlining    bool               // optional; zero value if not present.
	DebugInfoForProfiling bool               // optional; zero value if not present.
	NameTableKind         enum.NameTableKind // optional; zero value if not present.
	RangesBaseAddress     bool               // optional; zero value if not present.
	Sysroot               string             // optional; zero value if not present.
	SDK                   string             // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DICompileUnit) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DICompileUnit) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DICompileUnit' '(' Fields=(DICompileUnitField separator ',')* ')'
func (md *DICompileUnit) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("language: %s", md.Language)
	fields = append(fields, field)
	field = fmt.Sprintf("file: %s", md.File)
	fields = append(fields, field)
	if len(md.Producer) > 0 {
		field = fmt.Sprintf("producer: %s", quote(md.Producer))
		fields = append(fields, field)
	}
	if md.IsOptimized {
		field = fmt.Sprintf("isOptimized: %t", md.IsOptimized)
		fields = append(fields, field)
	}
	if len(md.Flags) > 0 {
		field = fmt.Sprintf("flags: %s", quote(md.Flags))
		fields = append(fields, field)
	}
	if md.RuntimeVersion != 0 {
		field = fmt.Sprintf("runtimeVersion: %d", md.RuntimeVersion)
		fields = append(fields, field)
	}
	if len(md.SplitDebugFilename) > 0 {
		field = fmt.Sprintf("splitDebugFilename: %s", quote(md.SplitDebugFilename))
		fields = append(fields, field)
	}
	if md.EmissionKind != 0 {
		field = fmt.Sprintf("emissionKind: %s", md.EmissionKind)
		fields = append(fields, field)
	}
	if md.Enums != nil {
		field = fmt.Sprintf("enums: %s", md.Enums)
		fields = append(fields, field)
	}
	if md.RetainedTypes != nil {
		field = fmt.Sprintf("retainedTypes: %s", md.RetainedTypes)
		fields = append(fields, field)
	}
	if md.Globals != nil {
		field = fmt.Sprintf("globals: %s", md.Globals)
		fields = append(fields, field)
	}
	if md.Imports != nil {
		field = fmt.Sprintf("imports: %s", md.Imports)
		fields = append(fields, field)
	}
	if md.Macros != nil {
		field = fmt.Sprintf("macros: %s", md.Macros)
		fields = append(fields, field)
	}
	if md.DwoID != 0 {
		field = fmt.Sprintf("dwoId: %d", md.DwoID)
		fields = append(fields, field)
	}
	if md.SplitDebugInlining {
		field = fmt.Sprintf("splitDebugInlining: %t", md.SplitDebugInlining)
		fields = append(fields, field)
	}
	if md.DebugInfoForProfiling {
		field = fmt.Sprintf("debugInfoForProfiling: %t", md.DebugInfoForProfiling)
		fields = append(fields, field)
	}
	if md.NameTableKind != 0 {
		field = fmt.Sprintf("nameTableKind: %s", md.NameTableKind)
		fields = append(fields, field)
	}
	if md.RangesBaseAddress {
		field = fmt.Sprintf("rangesBaseAddress: %t", md.RangesBaseAddress)
		fields = append(fields, field)
	}
	if len(md.Sysroot) > 0 {
		field = fmt.Sprintf("sysroot: %s", quote(md.Sysroot))
		fields = append(fields, field)
	}
	if len(md.SDK) > 0 {
		field = fmt.Sprintf("sdk: %s", quote(md.SDK))
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DICompileUnit(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DICompileUnit) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DICompositeType ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DICompositeType is a specialized metadata node.
type DICompositeType struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Tag         enum.DwarfTag  // required.
	Name        string         // optional; empty if not present.
	Scope       Field          // optional; nil if not present.
	File        *DIFile        // optional; nil if not present.
	Line        int64          // optional; zero value if not present.
	BaseType    Field          // optional; nil if not present.
	Size        uint64         // optional; zero value if not present.
	Align       uint64         // optional; zero value if not present.
	Offset      uint64         // optional; zero value if not present.
	Flags       enum.DIFlag    // optional.
	Elements    *Tuple         // optional; nil if not present.
	RuntimeLang enum.DwarfLang // optional; zero value if not present.
	// *DIBasicType or *DICompositeType
	VtableHolder   Field      // optional; nil if not present.
	TemplateParams *Tuple     // optional; nil if not present.
	Identifier     string     // optional; empty if not present.
	Discriminator  Field      // optional; nil if not present.
	DataLocation   Field      // optional; nil if not present.
	Associated     Field      // optional; nil if not present.
	Allocated      Field      // optional; nil if not present.
	Rank           FieldOrInt // optional; nil if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DICompositeType) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DICompositeType) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DICompositeType' '(' Fields=(DICompositeTypeField separator ',')* ')'
func (md *DICompositeType) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("tag: %s", dwarfTagString(md.Tag))
	fields = append(fields, field)
	if len(md.Name) > 0 {
		field := fmt.Sprintf("name: %s", quote(md.Name))
		fields = append(fields, field)
	}
	if md.Scope != nil {
		field := fmt.Sprintf("scope: %s", md.Scope)
		fields = append(fields, field)
	}
	if md.File != nil {
		field := fmt.Sprintf("file: %s", md.File)
		fields = append(fields, field)
	}
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	if md.BaseType != nil {
		field := fmt.Sprintf("baseType: %s", md.BaseType)
		fields = append(fields, field)
	}
	if md.Size != 0 {
		field := fmt.Sprintf("size: %d", md.Size)
		fields = append(fields, field)
	}
	if md.Align != 0 {
		field := fmt.Sprintf("align: %d", md.Align)
		fields = append(fields, field)
	}
	if md.Offset != 0 {
		field := fmt.Sprintf("offset: %d", md.Offset)
		fields = append(fields, field)
	}
	if md.Flags != 0 {
		field = fmt.Sprintf("flags: %s", diFlagsString(md.Flags))
		fields = append(fields, field)
	}
	if md.Elements != nil {
		field := fmt.Sprintf("elements: %s", md.Elements)
		fields = append(fields, field)
	}
	if md.RuntimeLang != 0 {
		field := fmt.Sprintf("runtimeLang: %s", md.RuntimeLang)
		fields = append(fields, field)
	}
	if md.VtableHolder != nil {
		field := fmt.Sprintf("vtableHolder: %s", md.VtableHolder)
		fields = append(fields, field)
	}
	if md.TemplateParams != nil {
		field := fmt.Sprintf("templateParams: %s", md.TemplateParams)
		fields = append(fields, field)
	}
	if len(md.Identifier) > 0 {
		field := fmt.Sprintf("identifier: %s", quote(md.Identifier))
		fields = append(fields, field)
	}
	if md.Discriminator != nil {
		field := fmt.Sprintf("discriminator: %s", md.Discriminator)
		fields = append(fields, field)
	}
	if md.DataLocation != nil {
		field := fmt.Sprintf("dataLocation: %s", md.DataLocation)
		fields = append(fields, field)
	}
	if md.Associated != nil {
		field := fmt.Sprintf("associated: %s", md.Associated)
		fields = append(fields, field)
	}
	if md.Allocated != nil {
		field := fmt.Sprintf("allocated: %s", md.Allocated)
		fields = append(fields, field)
	}
	if md.Rank != nil {
		field := fmt.Sprintf("rank: %s", md.Rank)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DICompositeType(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DICompositeType) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DIDerivedType ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIDerivedType is a specialized metadata node.
type DIDerivedType struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Tag               enum.DwarfTag // required.
	Name              string        // optional; empty if not present.
	Scope             Field         // optional; nil if not present.
	File              *DIFile       // optional; nil if not present.
	Line              int64         // optional; zero value if not present.
	BaseType          Field         // required.
	Size              uint64        // optional; zero value if not present.
	Align             uint64        // optional; zero value if not present.
	Offset            uint64        // optional; zero value if not present.
	Flags             enum.DIFlag   // optional.
	ExtraData         Field         // optional; nil if not present.
	DwarfAddressSpace uint64        // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIDerivedType) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIDerivedType) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIDerivedType' '(' Fields=(DIDerivedTypeField separator ',')* ')'
func (md *DIDerivedType) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("tag: %s", dwarfTagString(md.Tag))
	fields = append(fields, field)
	if len(md.Name) > 0 {
		field := fmt.Sprintf("name: %s", quote(md.Name))
		fields = append(fields, field)
	}
	if md.Scope != nil {
		field := fmt.Sprintf("scope: %s", md.Scope)
		fields = append(fields, field)
	}
	if md.File != nil {
		field := fmt.Sprintf("file: %s", md.File)
		fields = append(fields, field)
	}
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	field = fmt.Sprintf("baseType: %s", md.BaseType)
	fields = append(fields, field)
	if md.Size != 0 {
		field := fmt.Sprintf("size: %d", md.Size)
		fields = append(fields, field)
	}
	if md.Align != 0 {
		field := fmt.Sprintf("align: %d", md.Align)
		fields = append(fields, field)
	}
	if md.Offset != 0 {
		field := fmt.Sprintf("offset: %d", md.Offset)
		fields = append(fields, field)
	}
	if md.Flags != 0 {
		field = fmt.Sprintf("flags: %s", diFlagsString(md.Flags))
		fields = append(fields, field)
	}
	if md.ExtraData != nil {
		field := fmt.Sprintf("extraData: %s", md.ExtraData)
		fields = append(fields, field)
	}
	if md.DwarfAddressSpace != 0 {
		field := fmt.Sprintf("dwarfAddressSpace: %d", md.DwarfAddressSpace)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DIDerivedType(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIDerivedType) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DIEnumerator ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIEnumerator is a specialized metadata node.
type DIEnumerator struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Name       string // required.
	Value      int64  // required.
	IsUnsigned bool   // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIEnumerator) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIEnumerator) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIEnumerator' '(' Fields=(DIEnumeratorField separator ',')* ')'
func (md *DIEnumerator) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("name: %s", quote(md.Name))
	fields = append(fields, field)
	if md.IsUnsigned {
		field = fmt.Sprintf("value: %d", uint64(md.Value))
	} else {
		field = fmt.Sprintf("value: %d", md.Value)
	}
	fields = append(fields, field)
	if md.IsUnsigned {
		field := fmt.Sprintf("isUnsigned: %t", md.IsUnsigned)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DIEnumerator(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIEnumerator) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DIExpression ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIExpression is a specialized metadata node.
type DIExpression struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Fields []DIExpressionField
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIExpression) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIExpression) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIExpression' '(' Fields=(DIExpressionField separator ',')* ')'
func (md *DIExpression) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	buf.WriteString("!DIExpression(")
	for i, field := range md.Fields {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(field.String())
	}
	buf.WriteString(")")
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIExpression) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DIFile ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIFile is a specialized metadata node.
type DIFile struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Filename     string            // required.
	Directory    string            // required.
	Checksumkind enum.ChecksumKind // optional; zero value if not present.
	Checksum     string            // optional; empty if not present.
	Source       string            // optional; empty if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIFile) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIFile) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIFile' '(' Fields=(DIFileField separator ',')* ')'
func (md *DIFile) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("filename: %s", quote(md.Filename))
	fields = append(fields, field)
	field = fmt.Sprintf("directory: %s", quote(md.Directory))
	fields = append(fields, field)
	if md.Checksumkind != 0 {
		field := fmt.Sprintf("checksumkind: %s", md.Checksumkind)
		fields = append(fields, field)
	}
	if len(md.Checksum) > 0 {
		field := fmt.Sprintf("checksum: %s", quote(md.Checksum))
		fields = append(fields, field)
	}
	if len(md.Source) > 0 {
		field := fmt.Sprintf("source: %s", quote(md.Source))
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DIFile(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIFile) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DIGlobalVariable ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIGlobalVariable is a specialized metadata node.
type DIGlobalVariable struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Name           string  // required.
	Scope          Field   // optional; nil if not present.
	LinkageName    string  // optional; empty if not present.
	File           *DIFile // optional; nil if not present.
	Line           int64   // optional; zero value if not present.
	Type           Field   // optional; nil if not present.
	IsLocal        bool    // optional; zero value if not present.
	IsDefinition   bool    // optional; zero value if not present.
	TemplateParams *Tuple  // optional; nil if not present.
	Declaration    Field   // optional; nil if not present.
	Align          uint64  // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIGlobalVariable) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIGlobalVariable) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIGlobalVariable' '(' Fields=(DIGlobalVariableField separator ',')* ')'
func (md *DIGlobalVariable) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("name: %s", quote(md.Name))
	fields = append(fields, field)
	if md.Scope != nil {
		field := fmt.Sprintf("scope: %s", md.Scope)
		fields = append(fields, field)
	}
	if len(md.LinkageName) > 0 {
		field := fmt.Sprintf("linkageName: %s", quote(md.LinkageName))
		fields = append(fields, field)
	}
	if md.File != nil {
		field := fmt.Sprintf("file: %s", md.File)
		fields = append(fields, field)
	}
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	if md.Type != nil {
		field := fmt.Sprintf("type: %s", md.Type)
		fields = append(fields, field)
	}
	if md.IsLocal {
		field = fmt.Sprintf("isLocal: %t", md.IsLocal)
		fields = append(fields, field)
	}
	if md.IsDefinition {
		field := fmt.Sprintf("isDefinition: %t", md.IsDefinition)
		fields = append(fields, field)
	}
	if md.TemplateParams != nil {
		field := fmt.Sprintf("templateParams: %s", md.TemplateParams)
		fields = append(fields, field)
	}
	if md.Declaration != nil {
		field := fmt.Sprintf("declaration: %s", md.Declaration)
		fields = append(fields, field)
	}
	if md.Align != 0 {
		field := fmt.Sprintf("align: %d", md.Align)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DIGlobalVariable(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIGlobalVariable) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DIGlobalVariableExpression ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIGlobalVariableExpression is a specialized metadata node.
type DIGlobalVariableExpression struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Var  *DIGlobalVariable // required.
	Expr *DIExpression     // required.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIGlobalVariableExpression) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIGlobalVariableExpression) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIGlobalVariableExpression' '(' Fields=(DIGlobalVariableExpressionField separator ',')* ')'
func (md *DIGlobalVariableExpression) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("var: %s", md.Var)
	fields = append(fields, field)
	// NOTE: Should be required. Thus nil check should not be needed. However,
	// Clang outputs `!0 = !DIGlobalVariableExpression(var: !1)` in cat.ll.
	if md.Expr != nil {
		field = fmt.Sprintf("expr: %s", md.Expr)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DIGlobalVariableExpression(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIGlobalVariableExpression) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DIImportedEntity ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIImportedEntity is a specialized metadata node.
type DIImportedEntity struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Tag    enum.DwarfTag // required.
	Scope  Field         // required.
	Entity Field         // optional; nil if not present.
	File   *DIFile       // optional; nil if not present.
	Line   int64         // optional; zero value if not present.
	Name   string        // optional; empty if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIImportedEntity) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIImportedEntity) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIImportedEntity' '(' Fields=(DIImportedEntityField separator ',')* ')'
func (md *DIImportedEntity) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("tag: %s", dwarfTagString(md.Tag))
	fields = append(fields, field)
	field = fmt.Sprintf("scope: %s", md.Scope)
	fields = append(fields, field)
	if md.Entity != nil {
		field := fmt.Sprintf("entity: %s", md.Entity)
		fields = append(fields, field)
	}
	if md.File != nil {
		field := fmt.Sprintf("file: %s", md.File)
		fields = append(fields, field)
	}
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	if len(md.Name) > 0 {
		field := fmt.Sprintf("name: %s", quote(md.Name))
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DIImportedEntity(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIImportedEntity) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DILabel ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DILabel is a specialized metadata node.
type DILabel struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Scope Field   // required.
	Name  string  // required.
	File  *DIFile // required.
	Line  int64   // required.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DILabel) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DILabel) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DILabel' '(' Fields=(DILabelField separator ',')* ')'
func (md *DILabel) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("scope: %s", md.Scope)
	fields = append(fields, field)
	field = fmt.Sprintf("name: %s", quote(md.Name))
	fields = append(fields, field)
	field = fmt.Sprintf("file: %s", md.File)
	fields = append(fields, field)
	field = fmt.Sprintf("line: %d", md.Line)
	fields = append(fields, field)
	fmt.Fprintf(buf, "!DILabel(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DILabel) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DILexicalBlock ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DILexicalBlock is a specialized metadata node.
type DILexicalBlock struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Scope  Field   // required.
	File   *DIFile // optional; nil if not present.
	Line   int64   // optional; zero value if not present.
	Column int64   // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DILexicalBlock) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DILexicalBlock) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DILexicalBlock' '(' Fields=(DILexicalBlockField separator ',')* ')'
func (md *DILexicalBlock) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("scope: %s", md.Scope)
	fields = append(fields, field)
	if md.File != nil {
		field := fmt.Sprintf("file: %s", md.File)
		fields = append(fields, field)
	}
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	if md.Column != 0 {
		field := fmt.Sprintf("column: %d", md.Column)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DILexicalBlock(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DILexicalBlock) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DILexicalBlockFile ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DILexicalBlockFile is a specialized metadata node.
type DILexicalBlockFile struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Scope         Field   // required.
	File          *DIFile // optional; nil if not present.
	Discriminator uint64  // required.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DILexicalBlockFile) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DILexicalBlockFile) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DILexicalBlockFile' '(' Fields=(DILexicalBlockFileField separator ',')* ')'
func (md *DILexicalBlockFile) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("scope: %s", md.Scope)
	fields = append(fields, field)
	if md.File != nil {
		field := fmt.Sprintf("file: %s", md.File)
		fields = append(fields, field)
	}
	field = fmt.Sprintf("discriminator: %d", md.Discriminator)
	fields = append(fields, field)
	fmt.Fprintf(buf, "!DILexicalBlockFile(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DILexicalBlockFile) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DILocalVariable ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DILocalVariable is a specialized metadata node.
type DILocalVariable struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Scope Field       // required.
	Name  string      // optional; empty if not present.
	Arg   uint64      // optional; zero value if not present.
	File  *DIFile     // optional; nil if not present.
	Line  int64       // optional; zero value if not present.
	Type  Field       // optional; nil if not present.
	Flags enum.DIFlag // optional.
	Align uint64      // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DILocalVariable) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DILocalVariable) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DILocalVariable' '(' Fields=(DILocalVariableField separator ',')* ')'
func (md *DILocalVariable) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	// Note, to match Clang output, the output order is changed to output name
	// before scope.
	if len(md.Name) > 0 {
		field := fmt.Sprintf("name: %s", quote(md.Name))
		fields = append(fields, field)
	}
	if md.Arg != 0 {
		field := fmt.Sprintf("arg: %d", md.Arg)
		fields = append(fields, field)
	}
	field := fmt.Sprintf("scope: %s", md.Scope)
	fields = append(fields, field)
	if md.File != nil {
		field := fmt.Sprintf("file: %s", md.File)
		fields = append(fields, field)
	}
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	if md.Type != nil {
		field := fmt.Sprintf("type: %s", md.Type)
		fields = append(fields, field)
	}
	if md.Flags != 0 {
		field = fmt.Sprintf("flags: %s", diFlagsString(md.Flags))
		fields = append(fields, field)
	}
	if md.Align != 0 {
		field := fmt.Sprintf("align: %d", md.Align)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DILocalVariable(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DILocalVariable) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DILocation ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DILocation is a specialized metadata node.
type DILocation struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Line           int64       // optional; zero value if not present.
	Column         int64       // optional; zero value if not present.
	Scope          Field       // required.
	InlinedAt      *DILocation // optional; nil if not present.
	IsImplicitCode bool        // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DILocation) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DILocation) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DILocation' '(' Fields=(DILocationField separator ',')* ')'
func (md *DILocation) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	if md.Column != 0 {
		field := fmt.Sprintf("column: %d", md.Column)
		fields = append(fields, field)
	}
	field := fmt.Sprintf("scope: %s", md.Scope)
	fields = append(fields, field)
	if md.InlinedAt != nil {
		field := fmt.Sprintf("inlinedAt: %s", md.InlinedAt)
		fields = append(fields, field)
	}
	if md.IsImplicitCode {
		field := fmt.Sprintf("isImplicitCode: %t", md.IsImplicitCode)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DILocation(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DILocation) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DIMacro ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIMacro is a specialized metadata node.
type DIMacro struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Type  enum.DwarfMacinfo // required.
	Line  int64             // optional; zero value if not present.
	Name  string            // required.
	Value string            // optional; empty if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIMacro) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIMacro) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIMacro' '(' Fields=(DIMacroField separator ',')* ')'
func (md *DIMacro) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("type: %s", md.Type)
	fields = append(fields, field)
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	field = fmt.Sprintf("name: %s", quote(md.Name))
	fields = append(fields, field)
	if len(md.Value) > 0 {
		field := fmt.Sprintf("value: %s", quote(md.Value))
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DIMacro(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIMacro) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DIMacroFile ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIMacroFile is a specialized metadata node.
type DIMacroFile struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Type  enum.DwarfMacinfo // optional; zero value if not present.
	Line  int64             // optional; zero value if not present.
	File  *DIFile           // required.
	Nodes *Tuple            // optional; nil if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIMacroFile) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIMacroFile) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIMacroFile' '(' Fields=(DIMacroFileField separator ',')* ')'
func (md *DIMacroFile) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	if md.Type != 0 {
		field := fmt.Sprintf("type: %s", md.Type)
		fields = append(fields, field)
	}
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	field := fmt.Sprintf("file: %s", md.File)
	fields = append(fields, field)
	if md.Nodes != nil {
		field := fmt.Sprintf("nodes: %s", md.Nodes)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DIMacroFile(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIMacroFile) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DIModule ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIModule is a specialized metadata node.
type DIModule struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Scope        Field  // required.
	Name         string // required.
	ConfigMacros string // optional; empty if not present.
	IncludePath  string // optional; empty if not present.
	APINotes     string // optional; empty if not present.
	File         Field  // optional; empty if not present.
	Line         int64  // optional; zero value if not present.
	IsDecl       bool   // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIModule) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIModule) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIModule' '(' Fields=(DIModuleField separator ',')* ')'
func (md *DIModule) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("scope: %s", md.Scope)
	fields = append(fields, field)
	field = fmt.Sprintf("name: %s", quote(md.Name))
	fields = append(fields, field)
	if len(md.ConfigMacros) > 0 {
		field := fmt.Sprintf("configMacros: %s", quote(md.ConfigMacros))
		fields = append(fields, field)
	}
	if len(md.IncludePath) > 0 {
		field := fmt.Sprintf("includePath: %s", quote(md.IncludePath))
		fields = append(fields, field)
	}
	if len(md.APINotes) > 0 {
		field := fmt.Sprintf("apinotes: %s", quote(md.APINotes))
		fields = append(fields, field)
	}
	if md.File != nil {
		field := fmt.Sprintf("file: %s", md.File)
		fields = append(fields, field)
	}
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	if md.IsDecl {
		field := fmt.Sprintf("isDecl: %v", md.IsDecl)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DIModule(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIModule) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DINamespace ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DINamespace is a specialized metadata node.
type DINamespace struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Scope         Field  // required.
	Name          string // optional; empty if not present.
	ExportSymbols bool   // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DINamespace) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DINamespace) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DINamespace' '(' Fields=(DINamespaceField separator ',')* ')'
func (md *DINamespace) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("scope: %s", md.Scope)
	fields = append(fields, field)
	if len(md.Name) > 0 {
		field := fmt.Sprintf("name: %s", quote(md.Name))
		fields = append(fields, field)
	}
	if md.ExportSymbols {
		field := fmt.Sprintf("exportSymbols: %t", md.ExportSymbols)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DINamespace(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DINamespace) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DIObjCProperty ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DIObjCProperty is a specialized metadata node.
type DIObjCProperty struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Name       string  // optional; empty if not present.
	File       *DIFile // optional; nil if not present.
	Line       int64   // optional; zero value if not present.
	Setter     string  // optional; empty if not present.
	Getter     string  // optional; empty if not present.
	Attributes uint64  // optional; zero value if not present.
	Type       Field   // optional; nil if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DIObjCProperty) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DIObjCProperty) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DIObjCProperty' '(' Fields=(DIObjCPropertyField separator ',')* ')'
func (md *DIObjCProperty) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	if len(md.Name) > 0 {
		field := fmt.Sprintf("name: %s", quote(md.Name))
		fields = append(fields, field)
	}
	if md.File != nil {
		field := fmt.Sprintf("file: %s", md.File)
		fields = append(fields, field)
	}
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	if len(md.Setter) > 0 {
		field := fmt.Sprintf("setter: %s", quote(md.Setter))
		fields = append(fields, field)
	}
	if len(md.Getter) > 0 {
		field := fmt.Sprintf("getter: %s", quote(md.Getter))
		fields = append(fields, field)
	}
	if md.Attributes != 0 {
		field := fmt.Sprintf("attributes: %d", md.Attributes)
		fields = append(fields, field)
	}
	if md.Type != nil {
		field := fmt.Sprintf("type: %s", md.Type)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DIObjCProperty(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DIObjCProperty) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DISubprogram ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DISubprogram is a specialized metadata node.
type DISubprogram struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Scope          Field                // optional; nil if not present.
	Name           string               // optional; empty if not present.
	LinkageName    string               // optional; empty if not present.
	File           *DIFile              // optional; nil if not present.
	Line           int64                // optional; zero value if not present.
	Type           Field                // optional; nil if not present.
	IsLocal        bool                 // optional; zero value if not present.
	IsDefinition   bool                 // optional; zero value if not present.
	ScopeLine      int64                // optional; zero value if not present.
	ContainingType Field                // optional; nil if not present.
	Virtuality     enum.DwarfVirtuality // optional; zero value if not present.
	VirtualIndex   uint64               // optional; zero value if not present.
	ThisAdjustment int64                // optional; zero value if not present.
	Flags          enum.DIFlag          // optional.
	SPFlags        enum.DISPFlag        // optional.
	IsOptimized    bool                 // optional; zero value if not present.
	Unit           *DICompileUnit       // optional; nil if not present.
	TemplateParams *Tuple               // optional; nil if not present.
	Declaration    Field                // optional; nil if not present.
	RetainedNodes  *Tuple               // optional; nil if not present.
	ThrownTypes    *Tuple               // optional; nil if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DISubprogram) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DISubprogram) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DISubprogram' '(' Fields=(DISubprogramField separator ',')* ')'
func (md *DISubprogram) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	// Note, to match Clang output, the output order is changed to output name
	// before scope.
	if len(md.Name) > 0 {
		field := fmt.Sprintf("name: %s", quote(md.Name))
		fields = append(fields, field)
	}
	// Note, to match Clang output, the output order is changed to output
	// linkageName before scope.
	if len(md.LinkageName) > 0 {
		field := fmt.Sprintf("linkageName: %s", quote(md.LinkageName))
		fields = append(fields, field)
	}
	if md.Scope != nil {
		field := fmt.Sprintf("scope: %s", md.Scope)
		fields = append(fields, field)
	}
	if md.File != nil {
		field := fmt.Sprintf("file: %s", md.File)
		fields = append(fields, field)
	}
	if md.Line != 0 {
		field := fmt.Sprintf("line: %d", md.Line)
		fields = append(fields, field)
	}
	if md.Type != nil {
		field := fmt.Sprintf("type: %s", md.Type)
		fields = append(fields, field)
	}
	if md.IsLocal {
		field := fmt.Sprintf("isLocal: %t", md.IsLocal)
		fields = append(fields, field)
	}
	// Note: IsDefinition should be optional. However, Clang 9.0 produces
	// !DISubprogram specialized metadata nodes which looks as follows.
	//
	//    !80 = !DISubprogram(name: "abs", scope: !81, file: !81, line: 840, type: !82, flags: DIFlagPrototyped, spFlags: 0)
	//
	// Since `spFlags` is optional and has the zero value, we should be able to
	// remove it.
	//
	//    !80 = !DISubprogram(name: "abs", scope: !81, file: !81, line: 840, type: !82, flags: DIFlagPrototyped)
	//
	// However, doing so results in an error when run through `lli`, namely
	// `missing 'distinct', required for !DISubprogram that is a Definition`.
	//
	// Rather than including a "dummy" zero value for the spFlags enum, we
	// specify intent more clearly by stating whether the subprogram is a
	// definition or not.
	//
	// For this reason, we output isDefinition if it has a non-zero value or if
	// !DISubProgram is not distinct.
	if md.IsDefinition || !md.Distinct {
		field := fmt.Sprintf("isDefinition: %t", md.IsDefinition)
		fields = append(fields, field)
	}
	if md.ScopeLine != 0 {
		field := fmt.Sprintf("scopeLine: %d", md.ScopeLine)
		fields = append(fields, field)
	}
	if md.ContainingType != nil {
		field := fmt.Sprintf("containingType: %s", md.ContainingType)
		fields = append(fields, field)
	}
	if md.Virtuality != 0 {
		field := fmt.Sprintf("virtuality: %s", md.Virtuality)
		fields = append(fields, field)
	}
	if md.VirtualIndex != 0 {
		field := fmt.Sprintf("virtualIndex: %d", md.VirtualIndex)
		fields = append(fields, field)
	}
	if md.ThisAdjustment != 0 {
		field := fmt.Sprintf("thisAdjustment: %d", md.ThisAdjustment)
		fields = append(fields, field)
	}
	if md.Flags != 0 {
		field := fmt.Sprintf("flags: %s", diFlagsString(md.Flags))
		fields = append(fields, field)
	}
	if md.SPFlags != 0 {
		field := fmt.Sprintf("spFlags: %s", dispFlagsString(md.SPFlags))
		fields = append(fields, field)
	}
	if md.IsOptimized {
		field := fmt.Sprintf("isOptimized: %t", md.IsOptimized)
		fields = append(fields, field)
	}
	if md.Unit != nil {
		field := fmt.Sprintf("unit: %s", md.Unit)
		fields = append(fields, field)
	}
	if md.TemplateParams != nil {
		field := fmt.Sprintf("templateParams: %s", md.TemplateParams)
		fields = append(fields, field)
	}
	if md.Declaration != nil {
		field := fmt.Sprintf("declaration: %s", md.Declaration)
		fields = append(fields, field)
	}
	if md.RetainedNodes != nil {
		field := fmt.Sprintf("retainedNodes: %s", md.RetainedNodes)
		fields = append(fields, field)
	}
	if md.ThrownTypes != nil {
		field := fmt.Sprintf("thrownTypes: %s", md.ThrownTypes)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DISubprogram(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DISubprogram) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DISubrange ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DISubrange is a specialized metadata node.
type DISubrange struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Count      FieldOrInt // optional
	LowerBound FieldOrInt // optional
	UpperBound FieldOrInt // optional
	Stride     FieldOrInt // optional
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DISubrange) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DISubrange) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DISubrange' '(' Fields=(DISubrangeField separator ',')* ')'
func (md *DISubrange) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	if md.Count != nil {
		field := fmt.Sprintf("count: %s", md.Count)
		fields = append(fields, field)
	}
	if md.LowerBound != nil {
		field := fmt.Sprintf("lowerBound: %s", md.LowerBound)
		fields = append(fields, field)
	}
	if md.UpperBound != nil {
		field := fmt.Sprintf("upperBound: %s", md.UpperBound)
		fields = append(fields, field)
	}
	if md.Stride != nil {
		field := fmt.Sprintf("stride: %s", md.Stride)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!DISubrange(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DISubrange) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DISubroutineType ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DISubroutineType is a specialized metadata node.
type DISubroutineType struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Flags enum.DIFlag  // optional.
	CC    enum.DwarfCC // optional; zero value if not present.
	Types *Tuple       // required.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DISubroutineType) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DISubroutineType) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DISubroutineType' '(' Fields=(DISubroutineTypeField separator ',')* ')'
func (md *DISubroutineType) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	if md.Flags != 0 {
		field := fmt.Sprintf("flags: %s", diFlagsString(md.Flags))
		fields = append(fields, field)
	}
	if md.CC != 0 {
		field := fmt.Sprintf("cc: %s", md.CC)
		fields = append(fields, field)
	}
	field := fmt.Sprintf("types: %s", md.Types)
	fields = append(fields, field)
	fmt.Fprintf(buf, "!DISubroutineType(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DISubroutineType) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DITemplateTypeParameter ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DITemplateTypeParameter is a specialized metadata node.
type DITemplateTypeParameter struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Name      string // optional; empty if not present.
	Type      Field  // required.
	Defaulted bool   // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DITemplateTypeParameter) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DITemplateTypeParameter) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DITemplateTypeParameter' '(' Fields=(DITemplateTypeParameterField separator ',')* ')'
func (md *DITemplateTypeParameter) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	if len(md.Name) > 0 {
		field := fmt.Sprintf("name: %s", quote(md.Name))
		fields = append(fields, field)
	}
	field := fmt.Sprintf("type: %s", md.Type)
	if md.Defaulted {
		field := fmt.Sprintf("defaulted: %v", md.Defaulted)
		fields = append(fields, field)
	}
	fields = append(fields, field)
	fmt.Fprintf(buf, "!DITemplateTypeParameter(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DITemplateTypeParameter) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ DITemplateValueParameter ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DITemplateValueParameter is a specialized metadata node.
type DITemplateValueParameter struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Tag       enum.DwarfTag // optional; zero value if not present.
	Name      string        // optional; empty if not present.
	Type      Field         // optional; nil if not present.
	Value     Field         // required.
	Defaulted bool          // optional; zero value if not present.
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *DITemplateValueParameter) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *DITemplateValueParameter) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!DITemplateValueParameter' '(' Fields=(DITemplateValueParameterField separator ',')* ')'
func (md *DITemplateValueParameter) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	if md.Tag != 0 {
		field := fmt.Sprintf("tag: %s", dwarfTagString(md.Tag))
		fields = append(fields, field)
	}
	if len(md.Name) > 0 {
		field := fmt.Sprintf("name: %s", quote(md.Name))
		fields = append(fields, field)
	}
	if md.Type != nil {
		field := fmt.Sprintf("type: %s", md.Type)
		fields = append(fields, field)
	}
	field := fmt.Sprintf("value: %s", md.Value)
	if md.Defaulted {
		field := fmt.Sprintf("defaulted: %v", md.Defaulted)
		fields = append(fields, field)
	}
	fields = append(fields, field)
	fmt.Fprintf(buf, "!DITemplateValueParameter(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *DITemplateValueParameter) SetDistinct(distinct bool) {
	md.Distinct = distinct
}

// ~~~ [ GenericDINode ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GenericDINode is a specialized GenericDINode metadata node.
type GenericDINode struct {
	// Metadata ID associated with the specialized metadata node; -1 if not
	// present.
	MetadataID
	// (optional) Distinct.
	Distinct bool

	Tag      enum.DwarfTag // required
	Header   string        // optional; empty if not present
	Operands []Field       // optional
}

// String returns the LLVM syntax representation of the specialized metadata
// node.
func (md *GenericDINode) String() string {
	return md.Ident()
}

// Ident returns the identifier associated with the specialized metadata node.
func (md *GenericDINode) Ident() string {
	if md == nil {
		return "null"
	}
	if md.MetadataID != -1 {
		return md.MetadataID.Ident()
	}
	return md.LLString()
}

// LLString returns the LLVM syntax representation of the specialized metadata
// node.
//
// '!GenericDINode' '(' Fields=(GenericDINodeField separator ',')* ')'
func (md *GenericDINode) LLString() string {
	buf := &strings.Builder{}
	if md.Distinct {
		buf.WriteString("distinct ")
	}
	var fields []string
	field := fmt.Sprintf("tag: %s", dwarfTagString(md.Tag))
	fields = append(fields, field)
	if len(md.Header) > 0 {
		field := fmt.Sprintf("header: %s", quote(md.Header))
		fields = append(fields, field)
	}
	if len(md.Operands) > 0 {
		// TODO: figure out what operands output should look like.
		buf := &strings.Builder{}
		buf.WriteString("{")
		for i, o := range md.Operands {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(o.String())
		}
		buf.WriteString("}")
		field = fmt.Sprintf("operands: %s", buf)
		fields = append(fields, field)
	}
	fmt.Fprintf(buf, "!GenericDINode(%s)", strings.Join(fields, ", "))
	return buf.String()
}

// SetDistinct specifies whether the metadata definition is dinstict.
func (md *GenericDINode) SetDistinct(distinct bool) {
	md.Distinct = distinct
}
