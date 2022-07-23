package llconstant

import (
	"testing"

	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

func TestIntIdent(t *testing.T) {
	golden := []struct {
		in   string
		want string
	}{
		// integers < 0x1000 are always represented in decimal notation.
		{in: "100", want: "100"},
		{in: "256", want: "256"},
		{in: "1000", want: "1000"},
		{in: "4095", want: "4095"},
		// integers >= 0x1000 are represented in hexadecimal notation if lower
		// entropy than decimal notation.
		{in: "4096", want: "u0x1000"},
		{in: "2147483648", want: "u0x80000000"},
		{in: "9218868437227405312", want: "u0x7FF0000000000000"},
		{in: "1000000000000000000", want: "1000000000000000000"}, // hex would be u0xDE0B6B3A7640000
		// negative integers are always represented in decimal notation.
		{in: "-100", want: "-100"},
		{in: "-256", want: "-256"},
		{in: "-9218868437227405312", want: "-9218868437227405312"},
		{in: "-1000000000000000000", want: "-1000000000000000000"},
	}
	for _, g := range golden {
		c, err := NewIntFromString(types.I64, g.in)
		if err != nil {
			t.Errorf("unable to parse integer literal %q; %v", g.in, err)
			continue
		}
		got := c.Ident()
		if g.want != got {
			t.Errorf("integer constant string mismatch; expected %q, got %q", g.want, got)
			continue
		}
	}
}
