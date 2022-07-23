package llutil

import (
	"fmt"

	"github.com/wa-lang/wa/internal/3rdparty/llir"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

// ResetTypes resets the (cached) types of instructions in the given function.
func ResetTypes(f *llir.Func) {
	for _, b := range f.Blocks {
		for _, inst := range b.Insts {
			valueInst, ok := inst.(value.Value)
			if !ok {
				continue
			}
			resetType(valueInst)
		}
	}
}

// resetType resets the (cached) type of the given instruction.
func resetType(inst value.Value) {
	switch inst := inst.(type) {
	// Unary instructions
	case *llir.InstFNeg:
		inst.Typ = nil
	// Binary instructions
	case *llir.InstAdd:
		inst.Typ = nil
	case *llir.InstFAdd:
		inst.Typ = nil
	case *llir.InstSub:
		inst.Typ = nil
	case *llir.InstFSub:
		inst.Typ = nil
	case *llir.InstMul:
		inst.Typ = nil
	case *llir.InstFMul:
		inst.Typ = nil
	case *llir.InstUDiv:
		inst.Typ = nil
	case *llir.InstSDiv:
		inst.Typ = nil
	case *llir.InstFDiv:
		inst.Typ = nil
	case *llir.InstURem:
		inst.Typ = nil
	case *llir.InstSRem:
		inst.Typ = nil
	case *llir.InstFRem:
		inst.Typ = nil
	// Bitwise instructions
	case *llir.InstShl:
		inst.Typ = nil
	case *llir.InstLShr:
		inst.Typ = nil
	case *llir.InstAShr:
		inst.Typ = nil
	case *llir.InstAnd:
		inst.Typ = nil
	case *llir.InstOr:
		inst.Typ = nil
	case *llir.InstXor:
		inst.Typ = nil
	// Vector instructions
	case *llir.InstExtractElement:
		inst.Typ = nil
	case *llir.InstInsertElement:
		inst.Typ = nil
	case *llir.InstShuffleVector:
		inst.Typ = nil
	// Aggregate instructions
	case *llir.InstExtractValue:
		inst.Typ = nil
	case *llir.InstInsertValue:
		inst.Typ = nil
	// Memory instructions
	case *llir.InstAlloca:
		inst.Typ = nil
	case *llir.InstLoad:
		// type not cached.
	case *llir.InstCmpXchg:
		inst.Typ = nil
	case *llir.InstAtomicRMW:
		inst.Typ = nil
	case *llir.InstGetElementPtr:
		inst.Typ = nil
	// Conversion instructions
	case *llir.InstTrunc:
		// type not cached.
	case *llir.InstZExt:
		// type not cached.
	case *llir.InstSExt:
		// type not cached.
	case *llir.InstFPTrunc:
		// type not cached.
	case *llir.InstFPExt:
		// type not cached.
	case *llir.InstFPToUI:
		// type not cached.
	case *llir.InstFPToSI:
		// type not cached.
	case *llir.InstUIToFP:
		// type not cached.
	case *llir.InstSIToFP:
		// type not cached.
	case *llir.InstPtrToInt:
		// type not cached.
	case *llir.InstIntToPtr:
		// type not cached.
	case *llir.InstBitCast:
		// type not cached.
	case *llir.InstAddrSpaceCast:
		// type not cached.
	// Other instructions
	case *llir.InstICmp:
		inst.Typ = nil
	case *llir.InstFCmp:
		inst.Typ = nil
	case *llir.InstPhi:
		inst.Typ = nil
	case *llir.InstSelect:
		inst.Typ = nil
	case *llir.InstCall:
		inst.Typ = nil
	case *llir.InstVAArg:
		// type not cached.
	case *llir.InstLandingPad:
		// type not cached.
	case *llir.InstCatchPad:
		// type not cached.
	case *llir.InstCleanupPad:
		// type not cached.
	default:
		panic(fmt.Errorf("support for instruction type %T not yet implemented", inst))
	}
}
