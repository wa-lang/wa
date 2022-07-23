package llmetadata

import "fmt"

// TODO: constraint what types may be assigned to Node, MDNode, etc (i.e. make
// them sum types).

// Node is a metadata node.
//
// A Node has one of the following underlying types.
//
//    metadata.Definition      // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#Definition
//    *metadata.DIExpression   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIExpression
type Node interface {
	// Ident returns the identifier associated with the metadata node.
	Ident() string
}

// Definition is a metadata definition.
//
// A Definition has one of the following underlying types.
//
//    metadata.MDNode   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#MDNode
type Definition interface {
	// String returns the LLVM syntax representation of the metadata.
	fmt.Stringer
	// Ident returns the identifier associated with the metadata definition.
	Ident() string
	// ID returns the ID of the metadata definition.
	ID() int64
	// SetID sets the ID of the metadata definition.
	SetID(id int64)
	// LLString returns the LLVM syntax representation of the metadata
	// definition.
	LLString() string
	// SetDistinct specifies whether the metadata definition is dinstict.
	SetDistinct(distinct bool)
}

// MDNode is a metadata node.
//
// A MDNode has one of the following underlying types.
//
//    *metadata.Tuple            // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#Tuple
//    metadata.Definition        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#Definition
//    metadata.SpecializedNode   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#SpecializedNode
type MDNode interface {
	// Ident returns the identifier associated with the metadata node.
	Ident() string
	// LLString returns the LLVM syntax representation of the metadata node.
	LLString() string
}

// Field is a metadata field.
//
// A Field has one of the following underlying types.
//
//    *metadata.NullLit   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#NullLit
//    metadata.Metadata   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#Metadata
type Field interface {
	// String returns the LLVM syntax representation of the metadata field.
	fmt.Stringer
}

// SpecializedNode is a specialized metadata node.
//
// A SpecializedNode has one of the following underlying types.
//
//    *metadata.DIBasicType                  // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIBasicType
//    *metadata.DICommonBlock                // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DICommonBlock
//    *metadata.DICompileUnit                // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DICompileUnit
//    *metadata.DICompositeType              // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DICompositeType
//    *metadata.DIDerivedType                // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIDerivedType
//    *metadata.DIEnumerator                 // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIEnumerator
//    *metadata.DIExpression                 // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIExpression
//    *metadata.DIFile                       // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIFile
//    *metadata.DIGlobalVariable             // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIGlobalVariable
//    *metadata.DIGlobalVariableExpression   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIGlobalVariableExpression
//    *metadata.DIImportedEntity             // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIImportedEntity
//    *metadata.DILabel                      // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DILabel
//    *metadata.DILexicalBlock               // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DILexicalBlock
//    *metadata.DILexicalBlockFile           // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DILexicalBlockFile
//    *metadata.DILocalVariable              // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DILocalVariable
//    *metadata.DILocation                   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DILocation
//    *metadata.DIMacro                      // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIMacro
//    *metadata.DIMacroFile                  // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIMacroFile
//    *metadata.DIModule                     // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIModule
//    *metadata.DINamespace                  // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DINamespace
//    *metadata.DIObjCProperty               // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DIObjCProperty
//    *metadata.DISubprogram                 // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DISubprogram
//    *metadata.DISubrange                   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DISubrange
//    *metadata.DISubroutineType             // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DISubroutineType
//    *metadata.DITemplateTypeParameter      // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DITemplateTypeParameter
//    *metadata.DITemplateValueParameter     // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#DITemplateValueParameter
//    *metadata.GenericDINode                // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#GenericDINode
type SpecializedNode interface {
	Definition
}

// FieldOrInt is a metadata field or integer.
//
// A FieldOrInt has one of the following underlying types.
//
//    metadata.Field    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#Field
//    metadata.IntLit   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#IntLit
type FieldOrInt interface {
	Field
}

// DIExpressionField is a metadata DIExpression field.
//
// A DIExpressionField has one of the following underlying types.
//
//    metadata.UintLit        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#UintLit
//    enum.DwarfAttEncoding   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/enum#DwarfAttEncoding
//    enum.DwarfOp            // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/enum#DwarfOp
type DIExpressionField interface {
	fmt.Stringer
	// IsDIExpressionField ensures that only DIExpression fields can be assigned
	// to the metadata.DIExpressionField interface.
	IsDIExpressionField()
}

// IsDIExpressionField ensures that only DIExpression fields can be assigned to
// the metadata.DIExpressionField interface.
func (UintLit) IsDIExpressionField() {}

// Metadata is a sumtype of metadata.
//
// A Metadata has one of the following underlying types.
//
//    value.Value                // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/value#Value
//    *metadata.String           // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#String
//    *metadata.Tuple            // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#Tuple
//    metadata.Definition        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#Definition
//    metadata.SpecializedNode   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata#SpecializedNode
type Metadata interface {
	// String returns the LLVM syntax representation of the metadata.
	fmt.Stringer
}
