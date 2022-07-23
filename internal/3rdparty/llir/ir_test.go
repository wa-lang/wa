package llir

import (
	"strings"
	"testing"

	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	metadata "github.com/wa-lang/wa/internal/3rdparty/llir/llmetadata"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

func TestModuleString(t *testing.T) {
	golden := []struct {
		in   *Module
		want string
	}{
		// Empty module.
		{
			in:   &Module{},
			want: "",
		},
		// Type definition.
		{
			in: &Module{
				TypeDefs: []types.Type{&types.StructType{
					TypeName: "foo",
					Fields:   []types.Type{types.I32},
				}},
			},
			want: "%foo = type { i32 }",
		},
	}
	for _, g := range golden {
		got := strings.TrimSpace(g.in.String())
		if g.want != got {
			t.Errorf("module mismatch; expected `%v`, got `%v`", g.want, got)
		}
	}
}

// Assert that each constant implements the constant.Constant interface.
var (
	// Constants.
	_ constant.Constant = (*Global)(nil)
	_ constant.Constant = (*Func)(nil)
	_ constant.Constant = (*Alias)(nil)
	_ constant.Constant = (*IFunc)(nil)
)

// Assert that each instruction implements the ir.Instruction interface.
var (
	// Binary instructions.
	_ Instruction = (*InstAdd)(nil)
	_ Instruction = (*InstFAdd)(nil)
	_ Instruction = (*InstSub)(nil)
	_ Instruction = (*InstFSub)(nil)
	_ Instruction = (*InstMul)(nil)
	_ Instruction = (*InstFMul)(nil)
	_ Instruction = (*InstUDiv)(nil)
	_ Instruction = (*InstSDiv)(nil)
	_ Instruction = (*InstFDiv)(nil)
	_ Instruction = (*InstURem)(nil)
	_ Instruction = (*InstSRem)(nil)
	_ Instruction = (*InstFRem)(nil)
	// Bitwise instructions.
	_ Instruction = (*InstShl)(nil)
	_ Instruction = (*InstLShr)(nil)
	_ Instruction = (*InstAShr)(nil)
	_ Instruction = (*InstAnd)(nil)
	_ Instruction = (*InstOr)(nil)
	_ Instruction = (*InstXor)(nil)
	// Vector instructions.
	_ Instruction = (*InstExtractElement)(nil)
	_ Instruction = (*InstInsertElement)(nil)
	_ Instruction = (*InstShuffleVector)(nil)
	// Aggregate instructions.
	_ Instruction = (*InstExtractValue)(nil)
	_ Instruction = (*InstInsertValue)(nil)
	// Memory instructions.
	_ Instruction = (*InstAlloca)(nil)
	_ Instruction = (*InstLoad)(nil)
	_ Instruction = (*InstStore)(nil)
	_ Instruction = (*InstFence)(nil)
	_ Instruction = (*InstCmpXchg)(nil)
	_ Instruction = (*InstAtomicRMW)(nil)
	_ Instruction = (*InstGetElementPtr)(nil)
	// Conversion instructions.
	_ Instruction = (*InstTrunc)(nil)
	_ Instruction = (*InstZExt)(nil)
	_ Instruction = (*InstSExt)(nil)
	_ Instruction = (*InstFPTrunc)(nil)
	_ Instruction = (*InstFPExt)(nil)
	_ Instruction = (*InstFPToUI)(nil)
	_ Instruction = (*InstFPToSI)(nil)
	_ Instruction = (*InstUIToFP)(nil)
	_ Instruction = (*InstSIToFP)(nil)
	_ Instruction = (*InstPtrToInt)(nil)
	_ Instruction = (*InstIntToPtr)(nil)
	_ Instruction = (*InstBitCast)(nil)
	_ Instruction = (*InstAddrSpaceCast)(nil)
	// Other instructions.
	_ Instruction = (*InstICmp)(nil)
	_ Instruction = (*InstFCmp)(nil)
	_ Instruction = (*InstPhi)(nil)
	_ Instruction = (*InstSelect)(nil)
	_ Instruction = (*InstCall)(nil)
	_ Instruction = (*InstVAArg)(nil)
	_ Instruction = (*InstLandingPad)(nil)
	_ Instruction = (*InstCatchPad)(nil)
	_ Instruction = (*InstCleanupPad)(nil)
)

// Assert that each terminator implements the ir.Terminator interface.
var (
	// Terminators.
	_ Terminator = (*TermRet)(nil)
	_ Terminator = (*TermBr)(nil)
	_ Terminator = (*TermCondBr)(nil)
	_ Terminator = (*TermSwitch)(nil)
	_ Terminator = (*TermIndirectBr)(nil)
	_ Terminator = (*TermInvoke)(nil)
	_ Terminator = (*TermResume)(nil)
	_ Terminator = (*TermCatchSwitch)(nil)
	_ Terminator = (*TermCatchRet)(nil)
	_ Terminator = (*TermCleanupRet)(nil)
	_ Terminator = (*TermUnreachable)(nil)
)

// Assert that each value implements the value.Value interface.
var (
	// Constants.
	// Checked in constant_test.go as ir.Constant embeds value.Value.
	_ value.Value = constant.Constant(nil)
	// Named values.
	// Checked in value_test.go as value.Named embeds value.Value.
	_ value.Value = value.Named(nil)
	// Inline assembler expressions.
	_ value.Value = (*InlineAsm)(nil)
	// Metadata values.
	_ value.Value = (*metadata.Value)(nil)
)

// Assert that each named value implements the value.Named interface.
var (
	// Other values.
	_ value.Named = (*Global)(nil)
	_ value.Named = (*Func)(nil)
	_ value.Named = (*Param)(nil)
	_ value.Named = (*Block)(nil)

	// Instructions.
	// Binary instructions.
	_ value.Named = (*InstAdd)(nil)
	_ value.Named = (*InstFAdd)(nil)
	_ value.Named = (*InstSub)(nil)
	_ value.Named = (*InstFSub)(nil)
	_ value.Named = (*InstMul)(nil)
	_ value.Named = (*InstFMul)(nil)
	_ value.Named = (*InstUDiv)(nil)
	_ value.Named = (*InstSDiv)(nil)
	_ value.Named = (*InstFDiv)(nil)
	_ value.Named = (*InstURem)(nil)
	_ value.Named = (*InstSRem)(nil)
	_ value.Named = (*InstFRem)(nil)
	// Bitwise instructions.
	_ value.Named = (*InstShl)(nil)
	_ value.Named = (*InstLShr)(nil)
	_ value.Named = (*InstAShr)(nil)
	_ value.Named = (*InstAnd)(nil)
	_ value.Named = (*InstOr)(nil)
	_ value.Named = (*InstXor)(nil)
	// Vector instructions.
	_ value.Named = (*InstExtractElement)(nil)
	_ value.Named = (*InstInsertElement)(nil)
	_ value.Named = (*InstShuffleVector)(nil)
	// Aggregate instructions.
	_ value.Named = (*InstExtractValue)(nil)
	_ value.Named = (*InstInsertValue)(nil)
	// Memory instructions.
	_ value.Named = (*InstAlloca)(nil)
	_ value.Named = (*InstLoad)(nil)
	_ value.Named = (*InstCmpXchg)(nil)
	_ value.Named = (*InstAtomicRMW)(nil)
	_ value.Named = (*InstGetElementPtr)(nil)
	// Conversion instructions.
	_ value.Named = (*InstTrunc)(nil)
	_ value.Named = (*InstZExt)(nil)
	_ value.Named = (*InstSExt)(nil)
	_ value.Named = (*InstFPTrunc)(nil)
	_ value.Named = (*InstFPExt)(nil)
	_ value.Named = (*InstFPToUI)(nil)
	_ value.Named = (*InstFPToSI)(nil)
	_ value.Named = (*InstUIToFP)(nil)
	_ value.Named = (*InstSIToFP)(nil)
	_ value.Named = (*InstPtrToInt)(nil)
	_ value.Named = (*InstIntToPtr)(nil)
	_ value.Named = (*InstBitCast)(nil)
	_ value.Named = (*InstAddrSpaceCast)(nil)
	// Other instructions.
	_ value.Named = (*InstICmp)(nil)
	_ value.Named = (*InstFCmp)(nil)
	_ value.Named = (*InstPhi)(nil)
	_ value.Named = (*InstSelect)(nil)
	_ value.Named = (*InstCall)(nil)
	_ value.Named = (*InstVAArg)(nil)
	_ value.Named = (*InstLandingPad)(nil)
	_ value.Named = (*InstCatchPad)(nil)
	_ value.Named = (*InstCleanupPad)(nil)

	// Terminators.
	_ value.Named = (*TermInvoke)(nil)
	_ value.Named = (*TermCatchSwitch)(nil) // token result used by catchpad
)
