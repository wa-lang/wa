package llir_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	ir "github.com/wa-lang/wa/internal/3rdparty/llir"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

func TestAssignGlobalIDs(t *testing.T) {
	// ref: https://github.com/llir/llvm/issues/148
	const want = `@0 = global [15 x i8] c"Hello, world!\0A\00"
@1 = global i32 0

define i32 @2() {
0:
	ret i32 1
}
`

	m := ir.NewModule()

	// should be @0
	s := "Hello, world!\n\x00"
	i := constant.NewCharArrayFromString(s)
	m.NewGlobalDef("", i)

	// should be @1
	i32 := types.I32
	zero := constant.NewInt(i32, 0)
	m.NewGlobalDef("", zero)

	// should be @2
	one := constant.NewInt(i32, 1)
	m.NewFunc("", i32).NewBlock("").NewRet(one)

	// Compare module output.
	got := m.String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("module mismatch (-want +got):\n%s", diff)
	}
}
