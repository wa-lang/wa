// Package float implements floating-point representation utility functions.
package float

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/wa-lang/wa/internal/3rdparty/float/internal/strconv"
)

// IsExact16 reports whether x may be represented exactly as a 16-bit
// floating-point value.
func IsExact16(x *big.Float) bool {
	f, acc := x.Float64()
	if acc != big.Exact {
		return false
	}
	s1 := strconv.FormatFloat(f, 'e', -1, 16)
	s2 := trimZeros(x.Text('e', 100))
	return s1 == s2
}

// IsExact32 reports whether x may be represented exactly as a 32-bit
// floating-point value.
func IsExact32(x *big.Float) bool {
	f, acc := x.Float32()
	if acc != big.Exact {
		return false
	}
	s1 := strconv.FormatFloat(float64(f), 'e', -1, 32)
	s2 := trimZeros(x.Text('e', 100))
	return s1 == s2
}

// IsExact64 reports whether x may be represented exactly as a 64-bit
// floating-point value.
func IsExact64(x *big.Float) bool {
	f, acc := x.Float64()
	if acc != big.Exact {
		return false
	}
	s1 := strconv.FormatFloat(f, 'e', -1, 64)
	s2 := trimZeros(x.Text('e', 100))
	return s1 == s2
}

// trimZeros trims trailing zeroes after the decimal point in the given
// floating-point value (represented in scientific notation). If all digits
// after the decimal point are trimmed this way, the decimal point is also
// trimmed.
func trimZeros(s string) string {
	epos := strings.Index(s, "e")
	if epos == -1 {
		panic(fmt.Errorf("unable to locate position of exponent (e.g. e+02) in %q", s))
	}
	pos := epos - 1
	for ; pos > 0; pos-- {
		if s[pos] != '0' {
			break
		}
	}
	if s[pos] != '.' {
		pos++
	}
	return fmt.Sprintf("%s%s", s[:pos], s[epos:])
}
