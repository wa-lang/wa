package llmetadata

import (
	enum "github.com/wa-lang/wa/internal/3rdparty/llir/llenum"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// Assert that each metadata node implements the metadata.Node interface.
var (
	_ Node = (Definition)(nil)
	_ Node = (*DIExpression)(nil)
)

// Assert that each metadata node implements the metadata.MDNode interface.
var (
	_ MDNode = (*Tuple)(nil)
	_ MDNode = (Definition)(nil)
	_ MDNode = (SpecializedNode)(nil)
)

// Assert that each metadata field implements the metadata.Field interface.
var (
	_ Field = (*NullLit)(nil)
	_ Field = (Metadata)(nil)
)

// Assert that each specialized metadata node implements the
// metadata.SpecializedNode interface.
var (
	_ SpecializedNode = (*DIBasicType)(nil)
	_ SpecializedNode = (*DICommonBlock)(nil)
	_ SpecializedNode = (*DICompileUnit)(nil)
	_ SpecializedNode = (*DICompositeType)(nil)
	_ SpecializedNode = (*DIDerivedType)(nil)
	_ SpecializedNode = (*DIEnumerator)(nil)
	_ SpecializedNode = (*DIExpression)(nil)
	_ SpecializedNode = (*DIFile)(nil)
	_ SpecializedNode = (*DIGlobalVariable)(nil)
	_ SpecializedNode = (*DIGlobalVariableExpression)(nil)
	_ SpecializedNode = (*DIImportedEntity)(nil)
	_ SpecializedNode = (*DILabel)(nil)
	_ SpecializedNode = (*DILexicalBlock)(nil)
	_ SpecializedNode = (*DILexicalBlockFile)(nil)
	_ SpecializedNode = (*DILocalVariable)(nil)
	_ SpecializedNode = (*DILocation)(nil)
	_ SpecializedNode = (*DIMacro)(nil)
	_ SpecializedNode = (*DIMacroFile)(nil)
	_ SpecializedNode = (*DIModule)(nil)
	_ SpecializedNode = (*DINamespace)(nil)
	_ SpecializedNode = (*DIObjCProperty)(nil)
	_ SpecializedNode = (*DISubprogram)(nil)
	_ SpecializedNode = (*DISubrange)(nil)
	_ SpecializedNode = (*DISubroutineType)(nil)
	_ SpecializedNode = (*DITemplateTypeParameter)(nil)
	_ SpecializedNode = (*DITemplateValueParameter)(nil)
	_ SpecializedNode = (*GenericDINode)(nil)
)

// Assert that each metadata field or integer implements the metadata.FieldOrInt
// interface.
var (
	_ FieldOrInt = (Field)(nil)
	_ FieldOrInt = IntLit(0)
)

// Assert that each metadata DIExpression field implements the
// metadata.DIExpressionField interface.
var (
	_ DIExpressionField = UintLit(0)
	_ DIExpressionField = enum.DwarfAttEncoding(0)
	_ DIExpressionField = enum.DwarfOp(0)
)

// Assert that each metadata type implements the metadata.Metadata interface.
var (
	_ Metadata = (value.Value)(nil)
	_ Metadata = (*String)(nil)
	_ Metadata = (*Tuple)(nil)
	_ Metadata = (Definition)(nil)
	_ Metadata = (SpecializedNode)(nil)
)
