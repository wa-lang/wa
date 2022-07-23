// Inspiration for the LLVM IR walker was taken from Go fix.

// Note: we only walk IR values during walk, not type not metadata fields.

package llutil

import (
	"fmt"

	"github.com/wa-lang/wa/internal/3rdparty/llir"
	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	metadata "github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// Walk walks the LLVM IR AST in depth-first order; invoking visit recursively
// for each non-nil child of root. If visit returns false, the walk is
// terminated.
func Walk(root interface{}, visit func(n interface{}) bool) {
	visited := make(map[interface{}]bool)
	walk(root, visit, visited)
}

// walk walks the LLVM IR AST in depth-first order; invoking visit recursively
// for each non-nil child of root. If visit returns false, the walk is
// terminated. Visited tracks visited nodes.
func walk(root interface{}, visit func(n interface{}) bool, visited map[interface{}]bool) {
	if visited[root] {
		return
	}
	visited[root] = true
	if !visit(root) {
		return
	}
	switch root := root.(type) {
	// pointer to pointer to struct.
	case **llir.Arg:
		walk(*root, visit, visited)
	case **llir.Block:
		walk(*root, visit, visited)
	case **llir.Case:
		walk(*root, visit, visited)
	case **llir.Clause:
		walk(*root, visit, visited)
	case **llir.Incoming:
		walk(*root, visit, visited)
	case **llir.Module:
		walk(*root, visit, visited)
	case **llir.OperandBundle:
		walk(*root, visit, visited)
	case **llir.Param:
		walk(*root, visit, visited)
	case **llir.UseListOrder:
		walk(*root, visit, visited)
	case **llir.UseListOrderBB:
		walk(*root, visit, visited)
	// Constants
	// Simple constants
	case **constant.Int:
		walk(*root, visit, visited)
	case **constant.Float:
		walk(*root, visit, visited)
	case **constant.Null:
		walk(*root, visit, visited)
	case **constant.NoneToken:
		walk(*root, visit, visited)
	// Complex constants
	case **constant.Struct:
		walk(*root, visit, visited)
	case **constant.Array:
		walk(*root, visit, visited)
	case **constant.CharArray:
		walk(*root, visit, visited)
	case **constant.Vector:
		walk(*root, visit, visited)
	case **constant.ZeroInitializer:
		walk(*root, visit, visited)
	// Global variable and function addresses
	case **llir.Global:
		walk(*root, visit, visited)
	case **llir.Func:
		walk(*root, visit, visited)
	case **llir.Alias:
		walk(*root, visit, visited)
	case **llir.IFunc:
		walk(*root, visit, visited)
	// Undefined values
	case **constant.Undef:
		walk(*root, visit, visited)
	// Addresses of basic blocks
	case **constant.BlockAddress:
		walk(*root, visit, visited)
	// Constant expressions
	// Unary expressions
	case **constant.ExprFNeg:
		walk(*root, visit, visited)
	// Binary expressions
	case **constant.ExprAdd:
		walk(*root, visit, visited)
	case **constant.ExprFAdd:
		walk(*root, visit, visited)
	case **constant.ExprSub:
		walk(*root, visit, visited)
	case **constant.ExprFSub:
		walk(*root, visit, visited)
	case **constant.ExprMul:
		walk(*root, visit, visited)
	case **constant.ExprFMul:
		walk(*root, visit, visited)
	case **constant.ExprUDiv:
		walk(*root, visit, visited)
	case **constant.ExprSDiv:
		walk(*root, visit, visited)
	case **constant.ExprFDiv:
		walk(*root, visit, visited)
	case **constant.ExprURem:
		walk(*root, visit, visited)
	case **constant.ExprSRem:
		walk(*root, visit, visited)
	case **constant.ExprFRem:
		walk(*root, visit, visited)
	// Bitwise expressions
	case **constant.ExprShl:
		walk(*root, visit, visited)
	case **constant.ExprLShr:
		walk(*root, visit, visited)
	case **constant.ExprAShr:
		walk(*root, visit, visited)
	case **constant.ExprAnd:
		walk(*root, visit, visited)
	case **constant.ExprOr:
		walk(*root, visit, visited)
	case **constant.ExprXor:
		walk(*root, visit, visited)
	// Vector expressions
	case **constant.ExprExtractElement:
		walk(*root, visit, visited)
	case **constant.ExprInsertElement:
		walk(*root, visit, visited)
	case **constant.ExprShuffleVector:
		walk(*root, visit, visited)
	// Aggregate expressions
	case **constant.ExprExtractValue:
		walk(*root, visit, visited)
	case **constant.ExprInsertValue:
		walk(*root, visit, visited)
	// Memory expressions
	case **constant.ExprGetElementPtr:
		walk(*root, visit, visited)
	// Conversion expressions
	case **constant.ExprTrunc:
		walk(*root, visit, visited)
	case **constant.ExprZExt:
		walk(*root, visit, visited)
	case **constant.ExprSExt:
		walk(*root, visit, visited)
	case **constant.ExprFPTrunc:
		walk(*root, visit, visited)
	case **constant.ExprFPExt:
		walk(*root, visit, visited)
	case **constant.ExprFPToUI:
		walk(*root, visit, visited)
	case **constant.ExprFPToSI:
		walk(*root, visit, visited)
	case **constant.ExprUIToFP:
		walk(*root, visit, visited)
	case **constant.ExprSIToFP:
		walk(*root, visit, visited)
	case **constant.ExprPtrToInt:
		walk(*root, visit, visited)
	case **constant.ExprIntToPtr:
		walk(*root, visit, visited)
	case **constant.ExprBitCast:
		walk(*root, visit, visited)
	case **constant.ExprAddrSpaceCast:
		walk(*root, visit, visited)
	// Other expressions
	case **constant.ExprICmp:
		walk(*root, visit, visited)
	case **constant.ExprFCmp:
		walk(*root, visit, visited)
	case **constant.ExprSelect:
		walk(*root, visit, visited)
		// Instructions
	// Unary instructions
	case **llir.InstFNeg:
		walk(*root, visit, visited)
	// Binary instructions
	case **llir.InstAdd:
		walk(*root, visit, visited)
	case **llir.InstFAdd:
		walk(*root, visit, visited)
	case **llir.InstSub:
		walk(*root, visit, visited)
	case **llir.InstFSub:
		walk(*root, visit, visited)
	case **llir.InstMul:
		walk(*root, visit, visited)
	case **llir.InstFMul:
		walk(*root, visit, visited)
	case **llir.InstUDiv:
		walk(*root, visit, visited)
	case **llir.InstSDiv:
		walk(*root, visit, visited)
	case **llir.InstFDiv:
		walk(*root, visit, visited)
	case **llir.InstURem:
		walk(*root, visit, visited)
	case **llir.InstSRem:
		walk(*root, visit, visited)
	case **llir.InstFRem:
		walk(*root, visit, visited)
	// Bitwise instructions
	case **llir.InstShl:
		walk(*root, visit, visited)
	case **llir.InstLShr:
		walk(*root, visit, visited)
	case **llir.InstAShr:
		walk(*root, visit, visited)
	case **llir.InstAnd:
		walk(*root, visit, visited)
	case **llir.InstOr:
		walk(*root, visit, visited)
	case **llir.InstXor:
		walk(*root, visit, visited)
	// Vector instructions
	case **llir.InstExtractElement:
		walk(*root, visit, visited)
	case **llir.InstInsertElement:
		walk(*root, visit, visited)
	case **llir.InstShuffleVector:
		walk(*root, visit, visited)
	// Aggregate instructions
	case **llir.InstExtractValue:
		walk(*root, visit, visited)
	case **llir.InstInsertValue:
		walk(*root, visit, visited)
	// Memory instructions
	case **llir.InstAlloca:
		walk(*root, visit, visited)
	case **llir.InstLoad:
		walk(*root, visit, visited)
	case **llir.InstStore:
		walk(*root, visit, visited)
	case **llir.InstFence:
		walk(*root, visit, visited)
		// nothing to do
	case **llir.InstCmpXchg:
		walk(*root, visit, visited)
	case **llir.InstAtomicRMW:
		walk(*root, visit, visited)
	case **llir.InstGetElementPtr:
		walk(*root, visit, visited)
	// Conversion instructions
	case **llir.InstTrunc:
		walk(*root, visit, visited)
	case **llir.InstZExt:
		walk(*root, visit, visited)
	case **llir.InstSExt:
		walk(*root, visit, visited)
	case **llir.InstFPTrunc:
		walk(*root, visit, visited)
	case **llir.InstFPExt:
		walk(*root, visit, visited)
	case **llir.InstFPToUI:
		walk(*root, visit, visited)
	case **llir.InstFPToSI:
		walk(*root, visit, visited)
	case **llir.InstUIToFP:
		walk(*root, visit, visited)
	case **llir.InstSIToFP:
		walk(*root, visit, visited)
	case **llir.InstPtrToInt:
		walk(*root, visit, visited)
	case **llir.InstIntToPtr:
		walk(*root, visit, visited)
	case **llir.InstBitCast:
		walk(*root, visit, visited)
	case **llir.InstAddrSpaceCast:
		walk(*root, visit, visited)
	// Other instructions
	case **llir.InstICmp:
		walk(*root, visit, visited)
	case **llir.InstFCmp:
		walk(*root, visit, visited)
	case **llir.InstPhi:
		walk(*root, visit, visited)
	case **llir.InstSelect:
		walk(*root, visit, visited)
	case **llir.InstCall:
		walk(*root, visit, visited)
	case **llir.InstVAArg:
		walk(*root, visit, visited)
	case **llir.InstLandingPad:
		walk(*root, visit, visited)
	case **llir.InstCatchPad:
		walk(*root, visit, visited)
	case **llir.InstCleanupPad:
		walk(*root, visit, visited)
	// Terminators
	case **llir.TermRet:
		walk(*root, visit, visited)
	case **llir.TermBr:
		walk(*root, visit, visited)
	case **llir.TermCondBr:
		walk(*root, visit, visited)
	case **llir.TermSwitch:
		walk(*root, visit, visited)
	case **llir.TermIndirectBr:
		walk(*root, visit, visited)
	case **llir.TermInvoke:
		walk(*root, visit, visited)
	case **llir.TermResume:
		walk(*root, visit, visited)
	case **llir.TermCatchSwitch:
		walk(*root, visit, visited)
	case **llir.TermCatchRet:
		walk(*root, visit, visited)
	case **llir.TermCleanupRet:
		walk(*root, visit, visited)
	case **llir.TermUnreachable:
		walk(*root, visit, visited)
	// Metadata.
	case **metadata.NamedDef:
		walk(*root, visit, visited)
	case **metadata.Tuple:
		walk(*root, visit, visited)
	case **metadata.Value:
		walk(*root, visit, visited)
	case **metadata.String:
		walk(*root, visit, visited)
	case **metadata.Attachment:
		walk(*root, visit, visited)
	case **metadata.NullLit:
		walk(*root, visit, visited)
	// Specialized metadata node.
	case **metadata.DIBasicType:
		walk(*root, visit, visited)
	case **metadata.DICommonBlock:
		walk(*root, visit, visited)
	case **metadata.DICompileUnit:
		walk(*root, visit, visited)
	case **metadata.DICompositeType:
		walk(*root, visit, visited)
	case **metadata.DIDerivedType:
		walk(*root, visit, visited)
	case **metadata.DIEnumerator:
		walk(*root, visit, visited)
	case **metadata.DIExpression:
		walk(*root, visit, visited)
	case **metadata.DIFile:
		walk(*root, visit, visited)
	case **metadata.DIGlobalVariable:
		walk(*root, visit, visited)
	case **metadata.DIGlobalVariableExpression:
		walk(*root, visit, visited)
	case **metadata.DIImportedEntity:
		walk(*root, visit, visited)
	case **metadata.DILabel:
		walk(*root, visit, visited)
	case **metadata.DILexicalBlock:
		walk(*root, visit, visited)
	case **metadata.DILexicalBlockFile:
		walk(*root, visit, visited)
	case **metadata.DILocalVariable:
		walk(*root, visit, visited)
	case **metadata.DILocation:
		walk(*root, visit, visited)
	case **metadata.DIMacro:
		walk(*root, visit, visited)
	case **metadata.DIMacroFile:
		walk(*root, visit, visited)
	case **metadata.DIModule:
		walk(*root, visit, visited)
	case **metadata.DINamespace:
		walk(*root, visit, visited)
	case **metadata.DIObjCProperty:
		walk(*root, visit, visited)
	case **metadata.DISubprogram:
		walk(*root, visit, visited)
	case **metadata.DISubrange:
		walk(*root, visit, visited)
	case **metadata.DISubroutineType:
		walk(*root, visit, visited)
	case **metadata.DITemplateTypeParameter:
		walk(*root, visit, visited)
	case **metadata.DITemplateValueParameter:
		walk(*root, visit, visited)
	case **metadata.GenericDINode:
		walk(*root, visit, visited)

	// pointer to struct (with value receiver).
	case *metadata.IntLit:
		walk(*root, visit, visited)
	case *metadata.UintLit:
		walk(*root, visit, visited)

	// pointer to interface.
	case *constant.Constant:
		walk(*root, visit, visited)
	case *constant.Expression:
		walk(*root, visit, visited)
	case *llir.Instruction:
		walk(*root, visit, visited)
	case *llir.Terminator:
		walk(*root, visit, visited)
	case *value.Value:
		walk(*root, visit, visited)
	case *value.Named:
		walk(*root, visit, visited)
	// Metadata.
	case *metadata.Node:
		walk(*root, visit, visited)
	case *metadata.Definition:
		walk(*root, visit, visited)
	case *metadata.MDNode:
		walk(*root, visit, visited)
	case *metadata.Field:
		walk(*root, visit, visited)
	case *metadata.SpecializedNode:
		walk(*root, visit, visited)
	case *metadata.FieldOrInt:
		walk(*root, visit, visited)
	case *metadata.DIExpressionField:
		walk(*root, visit, visited)
	case *metadata.Metadata:
		walk(*root, visit, visited)

	// pointer to struct.
	case *llir.Arg:
		walk(&root.Value, visit, visited)
	case *llir.Block:
		for i := range root.Insts {
			walk(&root.Insts[i], visit, visited)
		}
		// allow walk on partial AST (terminator may not yet be set).
		if root.Term != nil {
			walk(&root.Term, visit, visited)
		}
	case *llir.Case:
		walk(&root.X, visit, visited)
		walk(&root.Target, visit, visited)
	case *llir.Clause:
		walk(&root.X, visit, visited)
	case *llir.Incoming:
		walk(&root.X, visit, visited)
		walk(&root.Pred, visit, visited)
	case *llir.Module:
		for i := range root.Globals {
			walk(&root.Globals[i], visit, visited)
		}
		for i := range root.Funcs {
			walk(&root.Funcs[i], visit, visited)
		}
		for i := range root.Aliases {
			walk(&root.Aliases[i], visit, visited)
		}
		for i := range root.IFuncs {
			walk(&root.IFuncs[i], visit, visited)
		}
		for i := range root.UseListOrders {
			walk(&root.UseListOrders[i], visit, visited)
		}
		for i := range root.UseListOrderBBs {
			walk(&root.UseListOrderBBs[i], visit, visited)
		}
	case *llir.OperandBundle:
		for i := range root.Inputs {
			walk(&root.Inputs[i], visit, visited)
		}
	case *llir.Param:
		// nothing to do
	case *llir.UseListOrder:
		walk(&root.Value, visit, visited)
	case *llir.UseListOrderBB:
		walk(&root.Func, visit, visited)
		walk(&root.Block, visit, visited)
	// Metadata.
	case *metadata.NamedDef:
		for i := range root.Nodes {
			walk(&root.Nodes[i], visit, visited)
		}
	case *metadata.Tuple:
		for i := range root.Fields {
			walk(&root.Fields[i], visit, visited)
		}
	case *metadata.Value:
		walk(&root.Value, visit, visited)
	case *metadata.String:
		// nothing to do.
	case *metadata.Attachment:
		walk(&root.Node, visit, visited)
	case *metadata.NullLit:
		// nothing to do.

	// struct (with value receiver).
	case metadata.IntLit:
		// nothing to do.
	case metadata.UintLit:
		// nothing to do.

	// interface.
	case constant.Constant:
		walkConst(root, visit, visited)
	case constant.Expression:
		walkConstExpr(root, visit, visited)
	case llir.Instruction:
		walkInst(root, visit, visited)
	case llir.Terminator:
		walkTerm(root, visit, visited)
	case value.Value:
		walkValue(root, visit, visited)
	case value.Named:
		walkValueNamed(root, visit, visited)
	case metadata.SpecializedNode:
		walkSpecializedMetadataNode(root, visit, visited)
	default:
		panic(fmt.Errorf("support for LLVM IR AST node type %T not yet implemented", root))
	}
}

// walkConst walks the LLVM IR AST in depth-first order; invoking visit
// recursively for each non-nil child of root. If visit returns false, the walk
// is terminated. Visited tracks visited nodes.
func walkConst(root constant.Constant, visit func(n interface{}) bool, visited map[interface{}]bool) {
	switch root := root.(type) {
	// Simple constants
	case *constant.Int:
		// nothing to do
	case *constant.Float:
		// nothing to do
	case *constant.Null:
		// nothing to do
	case *constant.NoneToken:
		// nothing to do
	// Complex constants
	case *constant.Struct:
		for i := range root.Fields {
			walk(&root.Fields[i], visit, visited)
		}
	case *constant.Array:
		for i := range root.Elems {
			walk(&root.Elems[i], visit, visited)
		}
	case *constant.CharArray:
		// nothing to do
	case *constant.Vector:
		for i := range root.Elems {
			walk(&root.Elems[i], visit, visited)
		}
	case *constant.ZeroInitializer:
		// nothing to do
	// Global variable and function addresses
	case *llir.Global:
		if root.Init != nil {
			walk(&root.Init, visit, visited)
		}
	case *llir.Func:
		for i := range root.Params {
			walk(&root.Params[i], visit, visited)
		}
		for i := range root.Blocks {
			walk(&root.Blocks[i], visit, visited)
		}
		if root.Prefix != nil {
			walk(&root.Prefix, visit, visited)
		}
		if root.Prologue != nil {
			walk(&root.Prologue, visit, visited)
		}
		if root.Personality != nil {
			walk(&root.Personality, visit, visited)
		}
		for i := range root.UseListOrders {
			walk(&root.UseListOrders[i], visit, visited)
		}
	case *llir.Alias:
		walk(&root.Aliasee, visit, visited)
	case *llir.IFunc:
		walk(&root.Resolver, visit, visited)
	// Undefined values
	case *constant.Undef:
		// nothing to do
	// Addresses of basic blocks
	case *constant.BlockAddress:
		walk(&root.Func, visit, visited)
		walk(&root.Block, visit, visited)
	// Constant expressions
	case constant.Expression:
		walk(root, visit, visited)
	default:
		panic(fmt.Errorf("support for LLVM IR AST node type %T not yet implemented", root))
	}
}

// walkConstExpr walks the LLVM IR AST in depth-first order; invoking visit
// recursively for each non-nil child of root. If visit returns false, the walk
// is terminated. Visited tracks visited nodes.
func walkConstExpr(root constant.Expression, visit func(n interface{}) bool, visited map[interface{}]bool) {
	switch root := root.(type) {
	// Unary expressions
	case *constant.ExprFNeg:
		walk(&root.X, visit, visited)
	// Binary expressions
	case *constant.ExprAdd:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprFAdd:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprSub:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprFSub:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprMul:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprFMul:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprUDiv:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprSDiv:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprFDiv:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprURem:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprSRem:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprFRem:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	// Bitwise expressions
	case *constant.ExprShl:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprLShr:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprAShr:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprAnd:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprOr:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprXor:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	// Vector expressions
	case *constant.ExprExtractElement:
		walk(&root.X, visit, visited)
		walk(&root.Index, visit, visited)
	case *constant.ExprInsertElement:
		walk(&root.X, visit, visited)
		walk(&root.Elem, visit, visited)
		walk(&root.Index, visit, visited)
	case *constant.ExprShuffleVector:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
		walk(&root.Mask, visit, visited)
	// Aggregate expressions
	case *constant.ExprExtractValue:
		walk(&root.X, visit, visited)
	case *constant.ExprInsertValue:
		walk(&root.X, visit, visited)
		walk(&root.Elem, visit, visited)
	// Memory expressions
	case *constant.ExprGetElementPtr:
		walk(&root.Src, visit, visited)
		for i := range root.Indices {
			walk(&root.Indices[i], visit, visited)
		}
	// Conversion expressions
	case *constant.ExprTrunc:
		walk(&root.From, visit, visited)
	case *constant.ExprZExt:
		walk(&root.From, visit, visited)
	case *constant.ExprSExt:
		walk(&root.From, visit, visited)
	case *constant.ExprFPTrunc:
		walk(&root.From, visit, visited)
	case *constant.ExprFPExt:
		walk(&root.From, visit, visited)
	case *constant.ExprFPToUI:
		walk(&root.From, visit, visited)
	case *constant.ExprFPToSI:
		walk(&root.From, visit, visited)
	case *constant.ExprUIToFP:
		walk(&root.From, visit, visited)
	case *constant.ExprSIToFP:
		walk(&root.From, visit, visited)
	case *constant.ExprPtrToInt:
		walk(&root.From, visit, visited)
	case *constant.ExprIntToPtr:
		walk(&root.From, visit, visited)
	case *constant.ExprBitCast:
		walk(&root.From, visit, visited)
	case *constant.ExprAddrSpaceCast:
		walk(&root.From, visit, visited)
	// Other expressions
	case *constant.ExprICmp:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprFCmp:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *constant.ExprSelect:
		walk(&root.Cond, visit, visited)
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	default:
		panic(fmt.Errorf("support for LLVM IR AST node type %T not yet implemented", root))
	}
}

// walkInst walks the LLVM IR AST in depth-first order; invoking visit
// recursively for each non-nil child of root. If visit returns false, the walk
// is terminated. Visited tracks visited nodes.
func walkInst(root llir.Instruction, visit func(n interface{}) bool, visited map[interface{}]bool) {
	switch root := root.(type) {
	// Unary instructions
	case *llir.InstFNeg:
		walk(&root.X, visit, visited)
	// Binary instructions
	case *llir.InstAdd:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstFAdd:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstSub:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstFSub:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstMul:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstFMul:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstUDiv:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstSDiv:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstFDiv:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstURem:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstSRem:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstFRem:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	// Bitwise instructions
	case *llir.InstShl:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstLShr:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstAShr:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstAnd:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstOr:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstXor:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	// Vector instructions
	case *llir.InstExtractElement:
		walk(&root.X, visit, visited)
		walk(&root.Index, visit, visited)
	case *llir.InstInsertElement:
		walk(&root.X, visit, visited)
		walk(&root.Elem, visit, visited)
		walk(&root.Index, visit, visited)
	case *llir.InstShuffleVector:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
		walk(&root.Mask, visit, visited)
	// Aggregate instructions
	case *llir.InstExtractValue:
		walk(&root.X, visit, visited)
	case *llir.InstInsertValue:
		walk(&root.X, visit, visited)
		walk(&root.Elem, visit, visited)
	// Memory instructions
	case *llir.InstAlloca:
		if root.NElems != nil {
			walk(&root.NElems, visit, visited)
		}
	case *llir.InstLoad:
		walk(&root.Src, visit, visited)
	case *llir.InstStore:
		walk(&root.Src, visit, visited)
		walk(&root.Dst, visit, visited)
	case *llir.InstFence:
		// nothing to do
	case *llir.InstCmpXchg:
		walk(&root.Ptr, visit, visited)
		walk(&root.Cmp, visit, visited)
		walk(&root.New, visit, visited)
	case *llir.InstAtomicRMW:
		walk(&root.Dst, visit, visited)
		walk(&root.X, visit, visited)
	case *llir.InstGetElementPtr:
		walk(&root.Src, visit, visited)
		for i := range root.Indices {
			walk(&root.Indices[i], visit, visited)
		}
	// Conversion instructions
	case *llir.InstTrunc:
		walk(&root.From, visit, visited)
	case *llir.InstZExt:
		walk(&root.From, visit, visited)
	case *llir.InstSExt:
		walk(&root.From, visit, visited)
	case *llir.InstFPTrunc:
		walk(&root.From, visit, visited)
	case *llir.InstFPExt:
		walk(&root.From, visit, visited)
	case *llir.InstFPToUI:
		walk(&root.From, visit, visited)
	case *llir.InstFPToSI:
		walk(&root.From, visit, visited)
	case *llir.InstUIToFP:
		walk(&root.From, visit, visited)
	case *llir.InstSIToFP:
		walk(&root.From, visit, visited)
	case *llir.InstPtrToInt:
		walk(&root.From, visit, visited)
	case *llir.InstIntToPtr:
		walk(&root.From, visit, visited)
	case *llir.InstBitCast:
		walk(&root.From, visit, visited)
	case *llir.InstAddrSpaceCast:
		walk(&root.From, visit, visited)
	// Other instructions
	case *llir.InstICmp:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstFCmp:
		walk(&root.X, visit, visited)
		walk(&root.Y, visit, visited)
	case *llir.InstPhi:
		for i := range root.Incs {
			walk(&root.Incs[i], visit, visited)
		}
	case *llir.InstSelect:
		walk(&root.Cond, visit, visited)
		walk(&root.ValueTrue, visit, visited)
		walk(&root.ValueFalse, visit, visited)
	case *llir.InstCall:
		walk(&root.Callee, visit, visited)
		for i := range root.Args {
			walk(&root.Args[i], visit, visited)
		}
	case *llir.InstVAArg:
		walk(&root.ArgList, visit, visited)
	case *llir.InstLandingPad:
		for i := range root.Clauses {
			walk(&root.Clauses[i], visit, visited)
		}
	case *llir.InstCatchPad:
		walk(&root.CatchSwitch, visit, visited)
		for i := range root.Args {
			walk(&root.Args[i], visit, visited)
		}
	case *llir.InstCleanupPad:
		walk(&root.ParentPad, visit, visited)
		for i := range root.Args {
			walk(&root.Args[i], visit, visited)
		}
	default:
		panic(fmt.Errorf("support for LLVM IR AST node type %T not yet implemented", root))
	}
}

// walkTerm walks the LLVM IR AST in depth-first order; invoking visit
// recursively for each non-nil child of root. If visit returns false, the walk
// is terminated. Visited tracks visited nodes.
func walkTerm(root llir.Terminator, visit func(n interface{}) bool, visited map[interface{}]bool) {
	switch root := root.(type) {
	// Terminators
	case *llir.TermRet:
		if root.X != nil {
			walk(&root.X, visit, visited)
		}
	case *llir.TermBr:
		walk(&root.Target, visit, visited)
	case *llir.TermCondBr:
		walk(&root.Cond, visit, visited)
		walk(&root.TargetTrue, visit, visited)
		walk(&root.TargetFalse, visit, visited)
	case *llir.TermSwitch:
		walk(&root.X, visit, visited)
		walk(&root.TargetDefault, visit, visited)
		for i := range root.Cases {
			walk(&root.Cases[i], visit, visited)
		}
	case *llir.TermIndirectBr:
		walk(&root.Addr, visit, visited)
		for i := range root.ValidTargets {
			walk(&root.ValidTargets[i], visit, visited)
		}
	case *llir.TermInvoke:
		walk(&root.Invokee, visit, visited)
		for i := range root.Args {
			walk(&root.Args[i], visit, visited)
		}
		walk(&root.NormalRetTarget, visit, visited)
		walk(&root.ExceptionRetTarget, visit, visited)
		for i := range root.OperandBundles {
			walk(&root.OperandBundles[i], visit, visited)
		}
	case *llir.TermResume:
		walk(&root.X, visit, visited)
	case *llir.TermCatchSwitch:
		walk(&root.ParentPad, visit, visited)
		for i := range root.Handlers {
			walk(&root.Handlers[i], visit, visited)
		}
		if root.DefaultUnwindTarget != nil {
			walk(&root.DefaultUnwindTarget, visit, visited)
		}
	case *llir.TermCatchRet:
		walk(&root.CatchPad, visit, visited)
		walk(&root.Target, visit, visited)
	case *llir.TermCleanupRet:
		walk(&root.CleanupPad, visit, visited)
		if root.UnwindTarget != nil {
			walk(&root.UnwindTarget, visit, visited)
		}
	case *llir.TermUnreachable:
		// nothing to do
	default:
		panic(fmt.Errorf("support for LLVM IR AST node type %T not yet implemented", root))
	}
}

// walkValue walks the LLVM IR AST in depth-first order; invoking visit
// recursively for each non-nil child of root. If visit returns false, the walk
// is terminated. Visited tracks visited nodes.
func walkValue(root value.Value, visit func(n interface{}) bool, visited map[interface{}]bool) {
	switch root := root.(type) {
	case constant.Constant:
		walk(root, visit, visited)
	case value.Named:
		walk(root, visit, visited)
	default:
		panic(fmt.Errorf("support for LLVM IR AST node type %T not yet implemented", root))
	}
}

// walkValueNamed walks the LLVM IR AST in depth-first order; invoking visit
// recursively for each non-nil child of root. If visit returns false, the walk
// is terminated. Visited tracks visited nodes.
func walkValueNamed(root value.Named, visit func(n interface{}) bool, visited map[interface{}]bool) {
	switch root := root.(type) {
	case *llir.Global:
		walk(root, visit, visited)
	case *llir.Func:
		walk(root, visit, visited)
	case *llir.Param:
		walk(root, visit, visited)
	case *llir.Block:
		walk(root, visit, visited)
	case llir.Instruction:
		walk(root, visit, visited)
	case *llir.TermInvoke:
		walk(root, visit, visited)
	case *llir.TermCatchSwitch:
		walk(root, visit, visited)
	default:
		panic(fmt.Errorf("support for LLVM IR AST node type %T not yet implemented", root))
	}
}

// walkSpecializedMetadataNode walks the LLVM IR AST in depth-first order;
// invoking visit recursively for each non-nil child of root. If visit returns
// false, the walk is terminated. Visited tracks visited nodes.
func walkSpecializedMetadataNode(root metadata.SpecializedNode, visit func(n interface{}) bool, visited map[interface{}]bool) {
	switch root := root.(type) {
	// Specialized metadata node.
	case *metadata.DIBasicType:
		// nothing to do.
	case *metadata.DICommonBlock:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
		if root.Declaration != nil {
			walk(&root.Declaration, visit, visited)
		}
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
	case *metadata.DICompileUnit:
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
		if root.Enums != nil {
			walk(&root.Enums, visit, visited)
		}
		if root.RetainedTypes != nil {
			walk(&root.RetainedTypes, visit, visited)
		}
		if root.Globals != nil {
			walk(&root.Globals, visit, visited)
		}
		if root.Imports != nil {
			walk(&root.Imports, visit, visited)
		}
		if root.Macros != nil {
			walk(&root.Macros, visit, visited)
		}
	case *metadata.DICompositeType:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
		if root.BaseType != nil {
			walk(&root.BaseType, visit, visited)
		}
		if root.Elements != nil {
			walk(&root.Elements, visit, visited)
		}
		if root.VtableHolder != nil {
			walk(&root.VtableHolder, visit, visited)
		}
		if root.TemplateParams != nil {
			walk(&root.TemplateParams, visit, visited)
		}
		if root.Discriminator != nil {
			walk(&root.Discriminator, visit, visited)
		}
	case *metadata.DIDerivedType:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
		if root.BaseType != nil {
			walk(&root.BaseType, visit, visited)
		}
		if root.ExtraData != nil {
			walk(&root.ExtraData, visit, visited)
		}
	case *metadata.DIEnumerator:
		// nothing to do.
	case *metadata.DIExpression:
		// nothing to do.
	case *metadata.DIFile:
		// nothing to do.
	case *metadata.DIGlobalVariable:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
		if root.Type != nil {
			walk(&root.Type, visit, visited)
		}
		if root.TemplateParams != nil {
			walk(&root.TemplateParams, visit, visited)
		}
		if root.Declaration != nil {
			walk(&root.Declaration, visit, visited)
		}
	case *metadata.DIGlobalVariableExpression:
		if root.Var != nil {
			walk(&root.Var, visit, visited)
		}
		if root.Expr != nil {
			walk(&root.Expr, visit, visited)
		}
	case *metadata.DIImportedEntity:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
		if root.Entity != nil {
			walk(&root.Entity, visit, visited)
		}
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
	case *metadata.DILabel:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
	case *metadata.DILexicalBlock:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
	case *metadata.DILexicalBlockFile:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
	case *metadata.DILocalVariable:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
		if root.Type != nil {
			walk(&root.Type, visit, visited)
		}
	case *metadata.DILocation:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
		if root.InlinedAt != nil {
			walk(&root.InlinedAt, visit, visited)
		}
	case *metadata.DIMacro:
		// nothing to do.
	case *metadata.DIMacroFile:
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
		if root.Nodes != nil {
			walk(&root.Nodes, visit, visited)
		}
	case *metadata.DIModule:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
	case *metadata.DINamespace:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
	case *metadata.DIObjCProperty:
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
		if root.Type != nil {
			walk(&root.Type, visit, visited)
		}
	case *metadata.DISubprogram:
		if root.Scope != nil {
			walk(&root.Scope, visit, visited)
		}
		if root.File != nil {
			walk(&root.File, visit, visited)
		}
		if root.Type != nil {
			walk(&root.Type, visit, visited)
		}
		if root.ContainingType != nil {
			walk(&root.ContainingType, visit, visited)
		}
		if root.Unit != nil {
			walk(&root.Unit, visit, visited)
		}
		if root.TemplateParams != nil {
			walk(&root.TemplateParams, visit, visited)
		}
		if root.Declaration != nil {
			walk(&root.Declaration, visit, visited)
		}
		if root.RetainedNodes != nil {
			walk(&root.RetainedNodes, visit, visited)
		}
		if root.ThrownTypes != nil {
			walk(&root.ThrownTypes, visit, visited)
		}
	case *metadata.DISubrange:
		if root.Count != nil {
			walk(&root.Count, visit, visited)
		}
	case *metadata.DISubroutineType:
		if root.Types != nil {
			walk(&root.Types, visit, visited)
		}
	case *metadata.DITemplateTypeParameter:
		if root.Type != nil {
			walk(&root.Type, visit, visited)
		}
	case *metadata.DITemplateValueParameter:
		if root.Type != nil {
			walk(&root.Type, visit, visited)
		}
		if root.Value != nil {
			walk(&root.Value, visit, visited)
		}
	case *metadata.GenericDINode:
		for i := range root.Operands {
			walk(&root.Operands[i], visit, visited)
		}
	}
}
