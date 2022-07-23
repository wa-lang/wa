package llir

// === [ Instructions ] ========================================================

// Instruction is an LLVM IR instruction. All instructions (except store and
// fence) implement the value.Named interface and may thus be used directly as
// values.
//
// An Instruction has one of the following underlying types.
//
// Unary instructions
//
// https://llvm.org/docs/LangRef.html#unary-operations
//
//    *ir.InstFNeg   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFNeg
//
// Binary instructions
//
// https://llvm.org/docs/LangRef.html#binary-operations
//
//    *ir.InstAdd    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstAdd
//    *ir.InstFAdd   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFAdd
//    *ir.InstSub    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstSub
//    *ir.InstFSub   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFSub
//    *ir.InstMul    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstMul
//    *ir.InstFMul   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFMul
//    *ir.InstUDiv   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstUDiv
//    *ir.InstSDiv   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstSDiv
//    *ir.InstFDiv   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFDiv
//    *ir.InstURem   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstURem
//    *ir.InstSRem   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstSRem
//    *ir.InstFRem   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFRem
//
// Bitwise instructions
//
// https://llvm.org/docs/LangRef.html#bitwise-binary-operations
//
//    *ir.InstShl    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstShl
//    *ir.InstLShr   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstLShr
//    *ir.InstAShr   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstAShr
//    *ir.InstAnd    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstAnd
//    *ir.InstOr     // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstOr
//    *ir.InstXor    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstXor
//
// Vector instructions
//
// https://llvm.org/docs/LangRef.html#vector-operations
//
//    *ir.InstExtractElement   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstExtractElement
//    *ir.InstInsertElement    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstInsertElement
//    *ir.InstShuffleVector    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstShuffleVector
//
// Aggregate instructions
//
// https://llvm.org/docs/LangRef.html#aggregate-operations
//
//    *ir.InstExtractValue   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstExtractValue
//    *ir.InstInsertValue    // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstInsertValue
//
// Memory instructions
//
// https://llvm.org/docs/LangRef.html#memory-access-and-addressing-operations
//
//    *ir.InstAlloca          // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstAlloca
//    *ir.InstLoad            // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstLoad
//    *ir.InstStore           // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstStore
//    *ir.InstFence           // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFence
//    *ir.InstCmpXchg         // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstCmpXchg
//    *ir.InstAtomicRMW       // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstAtomicRMW
//    *ir.InstGetElementPtr   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstGetElementPtr
//
// Conversion instructions
//
// https://llvm.org/docs/LangRef.html#conversion-operations
//
//    *ir.InstTrunc           // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstTrunc
//    *ir.InstZExt            // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstZExt
//    *ir.InstSExt            // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstSExt
//    *ir.InstFPTrunc         // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFPTrunc
//    *ir.InstFPExt           // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFPExt
//    *ir.InstFPToUI          // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFPToUI
//    *ir.InstFPToSI          // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFPToSI
//    *ir.InstUIToFP          // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstUIToFP
//    *ir.InstSIToFP          // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstSIToFP
//    *ir.InstPtrToInt        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstPtrToInt
//    *ir.InstIntToPtr        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstIntToPtr
//    *ir.InstBitCast         // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstBitCast
//    *ir.InstAddrSpaceCast   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstAddrSpaceCast
//
// Other instructions
//
// https://llvm.org/docs/LangRef.html#other-operations
//
//    *ir.InstICmp         // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstICmp
//    *ir.InstFCmp         // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFCmp
//    *ir.InstPhi          // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstPhi
//    *ir.InstSelect       // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstSelect
//    *ir.InstFreeze       // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstFreeze
//    *ir.InstCall         // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstCall
//    *ir.InstVAArg        // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstVAArg
//    *ir.InstLandingPad   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstLandingPad
//    *ir.InstCatchPad     // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstCatchPad
//    *ir.InstCleanupPad   // https://godoc.org/github.com/wa-lang/wa/internal/3rdparty/llir#InstCleanupPad
type Instruction interface {
	LLStringer
	// isInstruction ensures that only instructions can be assigned to the
	// instruction.Instruction interface.
	isInstruction()
}
