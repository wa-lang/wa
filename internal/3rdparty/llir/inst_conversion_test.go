package llir

import (
	"fmt"
	"testing"

	"github.com/wa-lang/wa/internal/3rdparty/errors"
	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

func TestTypeCheckTrunc(t *testing.T) {
	cases := []struct {
		fromTyp, toTyp types.Type
		panicMessage   string // "OK" if not panic'ing.
	}{
		{types.I64, types.I1,
			"OK"},
		{types.NewVector(2, types.I32), types.NewVector(2, types.I1),
			"OK"},

		{types.I32, types.I64,
			"invalid trunc operands: from.BitSize < to.BitSize (i32 is smaller than i64)"},
		{types.NewVector(2, types.I32), types.I1,
			"trunc operands are not compatible: from=<2 x i32>; to=i1"},
		{types.NewVector(1, types.I32), types.NewVector(2, types.I1),
			"trunc vector operand length mismatch: from=<1 x i32>; to=<2 x i1>"},
		{types.NewVector(2, types.I32), types.NewVector(2, types.I64),
			"invalid trunc operands: from.BitSize < to.BitSize (<2 x i32> is smaller than <2 x i64>)"},
	}

	errOK := errors.New("OK")

	for _, c := range cases {
		testName := fmt.Sprintf("%v to %v", c.fromTyp, c.toTyp)
		t.Run(testName, func(t *testing.T) {
			var panicErr error
			zeroVal := constant.NewZeroInitializer(c.fromTyp)
			func() {
				defer func() { panicErr = recover().(error) }()
				trunc := NewTrunc(zeroVal, c.toTyp)
				_ = trunc.String()
				panic(errOK)
			}()
			got := panicErr.Error()
			if got != c.panicMessage {
				t.Errorf("expected %q, got %q", c.panicMessage, got)
			}
		})
	}
}
