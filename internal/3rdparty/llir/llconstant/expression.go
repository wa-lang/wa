package llconstant

// === [ Expressions ] =========================================================

// Expression is an LLVM IR constant expression.
//
// An Expression has one of the following underlying types.
//
// Unary expressions
//
// https://llvm.org/docs/LangRef.html#constant-expressions
//
//    *constant.ExprFNeg   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprFNeg
//
// Binary expressions
//
// https://llvm.org/docs/LangRef.html#constant-expressions
//
//    *constant.ExprAdd    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprAdd
//    *constant.ExprFAdd   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprFAdd
//    *constant.ExprSub    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprSub
//    *constant.ExprFSub   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprFSub
//    *constant.ExprMul    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprMul
//    *constant.ExprFMul   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprFMul
//    *constant.ExprUDiv   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprUDiv
//    *constant.ExprSDiv   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprSDiv
//    *constant.ExprFDiv   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprFDiv
//    *constant.ExprURem   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprURem
//    *constant.ExprSRem   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprSRem
//    *constant.ExprFRem   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprFRem
//
// Bitwise expressions
//
// https://llvm.org/docs/LangRef.html#constant-expressions
//
//    *constant.ExprShl    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprShl
//    *constant.ExprLShr   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprLShr
//    *constant.ExprAShr   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprAShr
//    *constant.ExprAnd    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprAnd
//    *constant.ExprOr     // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprOr
//    *constant.ExprXor    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprXor
//
// Vector expressions
//
// https://llvm.org/docs/LangRef.html#constant-expressions
//
//    *constant.ExprExtractElement   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprExtractElement
//    *constant.ExprInsertElement    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprInsertElement
//    *constant.ExprShuffleVector    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprShuffleVector
//
// Aggregate expressions
//
// https://llvm.org/docs/LangRef.html#constant-expressions
//
//    *constant.ExprExtractValue   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprExtractValue
//    *constant.ExprInsertValue    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprInsertValue
//
// Memory expressions
//
// https://llvm.org/docs/LangRef.html#constant-expressions
//
//    *constant.ExprGetElementPtr   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprGetElementPtr
//
// Conversion expressions
//
// https://llvm.org/docs/LangRef.html#constant-expressions
//
//    *constant.ExprTrunc           // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprTrunc
//    *constant.ExprZExt            // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprZExt
//    *constant.ExprSExt            // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprSExt
//    *constant.ExprFPTrunc         // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprFPTrunc
//    *constant.ExprFPExt           // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprFPExt
//    *constant.ExprFPToUI          // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprFPToUI
//    *constant.ExprFPToSI          // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprFPToSI
//    *constant.ExprUIToFP          // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprUIToFP
//    *constant.ExprSIToFP          // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprSIToFP
//    *constant.ExprPtrToInt        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprPtrToInt
//    *constant.ExprIntToPtr        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprIntToPtr
//    *constant.ExprBitCast         // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprBitCast
//    *constant.ExprAddrSpaceCast   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprAddrSpaceCast
//
// Other expressions
//
// https://llvm.org/docs/LangRef.html#constant-expressions
//
//    *constant.ExprICmp     // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprICmp
//    *constant.ExprFCmp     // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprFCmp
//    *constant.ExprSelect   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir/llconstant#ExprSelect
type Expression interface {
	Constant
	// IsExpression ensures that only constants expressions can be assigned to
	// the constant.Expression interface.
	IsExpression()
}
