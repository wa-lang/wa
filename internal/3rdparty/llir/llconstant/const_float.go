package llconstant

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/wa-lang/wa/internal/3rdparty/errors"
	"github.com/wa-lang/wa/internal/3rdparty/float"
	"github.com/wa-lang/wa/internal/3rdparty/float/bfloat"
	"github.com/wa-lang/wa/internal/3rdparty/float/binary128"
	"github.com/wa-lang/wa/internal/3rdparty/float/binary16"
	"github.com/wa-lang/wa/internal/3rdparty/float/float128ppc"
	"github.com/wa-lang/wa/internal/3rdparty/float/float80x86"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
)

// --- [ Floating-point constants ] --------------------------------------------

// Float is an LLVM IR floating-point constant.
type Float struct {
	// Floating-point type.
	Typ *types.FloatType
	// Floating-point constant.
	X *big.Float
	// NaN specifies whether the floating-point constant is Not-a-Number.
	NaN bool
}

// NewFloat returns a new floating-point constant based on the given
// floating-point type and double precision floating-point value.
func NewFloat(typ *types.FloatType, x float64) *Float {
	if math.IsNaN(x) {
		f := &Float{Typ: typ, X: &big.Float{}, NaN: true}
		// Store sign of NaN.
		if math.Signbit(x) {
			f.X.SetFloat64(-1)
		}
		return f
	}
	return &Float{Typ: typ, X: big.NewFloat(x)}
}

// NewFloatFromString returns a new floating-point constant based on the given
// floating-point type and floating-point string.
//
// The floating-point string may be expressed in one of the following forms.
//
//    * fraction floating-point literal
//         [+-]? [0-9]+ [.] [0-9]*
//    * scientific notation floating-point literal
//         [+-]? [0-9]+ [.] [0-9]* [eE] [+-]? [0-9]+
//    * hexadecimal floating-point literal
//         0x[0-9A-Fa-f]{16}  // HexFP
//         0xK[0-9A-Fa-f]{20} // HexFP80
//         0xL[0-9A-Fa-f]{32} // HexFP128
//         0xM[0-9A-Fa-f]{32} // HexPPC128
//         0xH[0-9A-Fa-f]{4}  // HexHalf
func NewFloatFromString(typ *types.FloatType, s string) (*Float, error) {
	// Hexadecimal floating-point literal.
	if strings.HasPrefix(s, "0x") {
		switch {
		// x86_fp80 (x86 extended precision)
		case strings.HasPrefix(s, "0xK"):
			// From https://llvm.org/docs/LangRef.html#simple-constants
			//
			// > The 80-bit format used by x86 is represented as 0xK followed by 20
			// > hexadecimal digits.
			hex := strings.TrimPrefix(s, "0xK")
			const hexLen = 8
			part1 := hex[:hexLen/2]
			part2 := hex[hexLen/2:]
			se, err := strconv.ParseUint(part1, 16, 16)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			m, err := strconv.ParseUint(part2, 16, 64)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			f := float80x86.NewFromBits(uint16(se), m)
			x, nan := f.Big()
			return &Float{Typ: typ, X: x, NaN: nan}, nil
		// fp128 (IEEE 754 quadruple precision)
		case strings.HasPrefix(s, "0xL"):
			// From https://llvm.org/docs/LangRef.html#simple-constants
			//
			// > The IEEE 128-bit format is represented by 0xL followed by 32
			// > hexadecimal digits.
			hex := strings.TrimPrefix(s, "0xL")
			const maxHexLen = 32
			if len(hex) < maxHexLen {
				// pad with leading zeroes (e.g. for case like `0xL01`)
				hex = strings.Repeat("0", maxHexLen-len(hex)) + hex
			}
			part1 := hex[:maxHexLen/2]
			part2 := hex[maxHexLen/2:]
			a, err := strconv.ParseUint(part1, 16, 64)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			b, err := strconv.ParseUint(part2, 16, 64)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			f := binary128.NewFromBits(a, b)
			x, nan := f.Big()
			return &Float{Typ: typ, X: x, NaN: nan}, nil
		// ppc_fp128 (PowerPC double-double arithmetic)
		case strings.HasPrefix(s, "0xM"):
			// From https://llvm.org/docs/LangRef.html#simple-constants
			//
			// > The 128-bit format used by PowerPC (two adjacent doubles) is
			// > represented by 0xM followed by 32 hexadecimal digits.
			hex := strings.TrimPrefix(s, "0xM")
			const maxHexLen = 32
			part1 := hex[:maxHexLen/2]
			part2 := hex[maxHexLen/2:]
			a, err := strconv.ParseUint(part1, 16, 64)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			b, err := strconv.ParseUint(part2, 16, 64)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			f := float128ppc.NewFromBits(a, b)
			x, nan := f.Big()
			return &Float{Typ: typ, X: x, NaN: nan}, nil
		// half (IEEE 754 half precision)
		case strings.HasPrefix(s, "0xH"):
			// From https://llvm.org/docs/LangRef.html#simple-constants
			//
			// > The IEEE 16-bit format (half precision) is represented by 0xH
			// > followed by 4 hexadecimal digits.
			hex := strings.TrimPrefix(s, "0xH")
			bits, err := strconv.ParseUint(hex, 16, 16)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			f := binary16.NewFromBits(uint16(bits))
			x, nan := f.Big()
			return &Float{Typ: typ, X: x, NaN: nan}, nil
		// Hexadecimal floating-point literal.
		case strings.HasPrefix(s, "0xR"):
			// From https://llvm.org/docs/LangRef.html#simple-constants
			//
			// > The bfloat 16-bit format is represented by 0xR followed by 4 hexadecimal digits.
			hex := strings.TrimPrefix(s, "0xR")
			bits, err := strconv.ParseUint(hex, 16, 16)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			f := bfloat.NewFromBits(uint16(bits))
			x, nan := f.Big()
			return &Float{Typ: typ, X: x, NaN: nan}, nil
			// HexBFloatConstant
		default:
			// From https://llvm.org/docs/LangRef.html#simple-constants
			//
			// > When using the hexadecimal form, constants of types half, float,
			// > and double are represented using the 16-digit form shown above
			// > (which matches the IEEE754 representation for double).
			hex := strings.TrimPrefix(s, "0x")
			bits, err := strconv.ParseUint(hex, 16, 64)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			switch typ.Kind {
			case types.FloatKindHalf:
				// TODO: verify if this is a correct implementation. We should
				// probably be using binary16.NewFromBits.
				f16 := math.Float64frombits(bits)
				if math.IsNaN(f16) {
					f := &Float{Typ: typ, X: &big.Float{}, NaN: true}
					// Store sign of NaN.
					if math.Signbit(f16) {
						f.X.SetFloat64(-1)
					}
					return f, nil
				}
				c := big.NewFloat(f16)
				const precision = 11
				c.SetPrec(precision)
				return &Float{
					Typ: typ,
					X:   c,
				}, nil
			case types.FloatKindFloat:
				// ref: https://groups.google.com/d/msg/llvm-dev/IlqV3TbSk6M/27dAggZOMb0J
				//
				// The exact bit representation of the float is laid out with the
				// corresponding bitwise representation of a double: the sign bit is
				// copied over, the exponent is encoded in the larger width, and the
				// 23 bits of significand fills in the top 23 bits of significand in
				// the double. A double has 52 bits of significand, so this means
				// that the last 29 bits of significand will always be ignored. As
				// an error-detection measure, the IR parser requires them to be
				// zero.
				f32 := math.Float64frombits(bits)
				if math.IsNaN(f32) {
					f := &Float{Typ: typ, X: &big.Float{}, NaN: true}
					// Store sign of NaN.
					if math.Signbit(f32) {
						f.X.SetFloat64(-1)
					}
					return f, nil
				}
				x := big.NewFloat(f32)
				const precision = 24
				x.SetPrec(precision)
				return &Float{Typ: typ, X: x}, nil
			case types.FloatKindDouble:
				f64 := math.Float64frombits(bits)
				if math.IsNaN(f64) {
					f := &Float{Typ: typ, X: &big.Float{}, NaN: true}
					// Store sign of NaN.
					if math.Signbit(f64) {
						f.X.SetFloat64(-1)
					}
					return f, nil
				}
				x := big.NewFloat(f64)
				const precision = 53
				x.SetPrec(precision)
				return &Float{Typ: typ, X: x}, nil
			default:
				panic(fmt.Errorf("support for hexadecimal floating-point literal %q of kind %v not yet implemented", s, typ.Kind))
			}
		}
	}
	const base = 10
	switch typ.Kind {
	case types.FloatKindHalf:
		const precision = 11
		x, _, err := big.ParseFloat(s, base, precision, big.ToNearestEven)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		c := &Float{
			Typ: typ,
			X:   x,
		}
		return c, nil
	case types.FloatKindFloat:
		const precision = 24
		x, _, err := big.ParseFloat(s, base, precision, big.ToNearestEven)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		c := &Float{
			Typ: typ,
			X:   x,
		}
		return c, nil
	case types.FloatKindDouble:
		const precision = 53
		x, _, err := big.ParseFloat(s, base, precision, big.ToNearestEven)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		c := &Float{
			Typ: typ,
			X:   x,
		}
		return c, nil
	default:
		panic(fmt.Errorf("support for floating-point kind %v not yet implemented", typ.Kind))
	}
}

// String returns the LLVM syntax representation of the constant as a type-value
// pair.
func (c *Float) String() string {
	return fmt.Sprintf("%s %s", c.Type(), c.Ident())
}

// Type returns the type of the constant.
func (c *Float) Type() types.Type {
	return c.Typ
}

// Ident returns the identifier associated with the constant.
func (c *Float) Ident() string {
	// FloatLit
	//
	// Print hexadecimal representation of floating-point literal if NaN, Inf,
	// inexact or extended precision (x86_fp80, fp128 or ppc_fp128).
	switch c.Typ.Kind {
	// half (IEEE 754 half precision)
	case types.FloatKindHalf:
		const hexPrefix = 'H'
		if c.NaN {
			bits := binary16.NaN.Bits()
			if c.X != nil && c.X.Signbit() {
				bits = binary16.NegNaN.Bits()
			}
			return fmt.Sprintf("0x%c%04X", hexPrefix, bits)
		}
		if c.X.IsInf() || !float.IsExact16(c.X) {
			f, acc := binary16.NewFromBig(c.X)
			if acc != big.Exact {
				log.Printf("unable to represent floating-point constant %v of type %v exactly; please submit a bug report to llir/llvm with this error message", c.X, c.Typ)
			}
			bits := f.Bits()
			return fmt.Sprintf("0x%c%04X", hexPrefix, bits)
		}
		// c is representable without loss as floating-point literal, this case is
		// handled for half, float and double below the switch statement.
	// float (IEEE 754 single precision)
	case types.FloatKindFloat:
		// ref: https://groups.google.com/d/msg/llvm-dev/IlqV3TbSk6M/27dAggZOMb0J
		//
		// The exact bit representation of the float is laid out with the
		// corresponding bitwise representation of a double: the sign bit is
		// copied over, the exponent is encoded in the larger width, and the 23
		// bits of significand fills in the top 23 bits of significand in the
		// double. A double has 52 bits of significand, so this means that the
		// last 29 bits of significand will always be ignored. As an error
		// detection measure, the IR parser requires them to be zero.
		if c.NaN {
			f := math.NaN()
			if c.X != nil && c.X.Signbit() {
				f = math.Copysign(f, -1)
			}
			bits := math.Float64bits(f)
			// zero out last 29 bits.
			bits &^= 0x1FFFFFFF
			return fmt.Sprintf("0x%X", bits)
		}
		if c.X.IsInf() || !float.IsExact32(c.X) {
			f, _ := c.X.Float64()
			bits := math.Float64bits(f)
			// Note, to match Clang output we do not zero-pad the hexadecimal
			// output.
			// zero out last 29 bits.
			bits &^= 0x1FFFFFFF
			return fmt.Sprintf("0x%X", bits)
		}
		// c is representable without loss as floating-point literal, this case is
		// handled for half, float and double below the switch statement.
	// double (IEEE 754 double precision)
	case types.FloatKindDouble:
		if c.NaN {
			// To use the same cannonical representation as LLVM IR for NaN values, we
			// explicitly a qNaN value (quiet NaN) with the leading bit in the mantissa
			// set, rather than the trailing bit as used for the cannonical
			// representation in Go (see math.NaN).
			//
			// For further background, see https://github.com/llir/llvm/issues/133
			//
			//      exponent    mantissa
			//    s 11111111111 1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx = quiet     (qNaN)
			//    s 11111111111 0xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx = signaling (sNaN) **
			//                  ^ quiet bit
			//
			// Where ** denote that at least one of the 'x' bits has to be set, since the
			// mantissa must be non-zero to denote NaN.
			//
			// quiet NaN:
			//    0x7FF8000000000000 = 0b0_11111111111_100000000000000000000000000000000000000000000000000
			f := math.Float64frombits(0x7FF8000000000000) // quiet NaN
			if c.X != nil && c.X.Signbit() {
				f = math.Copysign(f, -1)
			}
			bits := math.Float64bits(f)
			return fmt.Sprintf("0x%X", bits)
		}
		if c.X.IsInf() || !float.IsExact64(c.X) {
			f, _ := c.X.Float64()
			bits := math.Float64bits(f)
			// Note, to match Clang output we do not zero-pad the hexadecimal
			// output.
			return fmt.Sprintf("0x%X", bits)
		}
		// c is representable without loss as floating-point literal, this case is
		// handled for half, float and double below the switch statement.
	// x86_fp80 (x86 extended precision)
	case types.FloatKindX86_FP80:
		// always represent x86_fp80 in hexadecimal floating-point notation.
		const hexPrefix = 'K'
		if c.NaN {
			se, m := float80x86.NaN.Bits()
			if c.X != nil && c.X.Signbit() {
				se, m = float80x86.NegNaN.Bits()
			}
			return fmt.Sprintf("0x%c%04X%016X", hexPrefix, se, m)
		}
		f, acc := float80x86.NewFromBig(c.X)
		if acc != big.Exact {
			log.Printf("unable to represent floating-point constant %v of type %v exactly; please submit a bug report to llir/llvm with this error message", c.X, c.Typ)
		}
		se, m := f.Bits()
		return fmt.Sprintf("0x%c%04X%016X", hexPrefix, se, m)
	// fp128 (IEEE 754 quadruple precision)
	case types.FloatKindFP128:
		// always represent fp128 in hexadecimal floating-point notation.
		const hexPrefix = 'L'
		if c.NaN {
			a, b := binary128.NaN.Bits()
			if c.X != nil && c.X.Signbit() {
				a, b = binary128.NegNaN.Bits()
			}
			return fmt.Sprintf("0x%c%016X%016X", hexPrefix, a, b)
		}
		f, acc := binary128.NewFromBig(c.X)
		if acc != big.Exact {
			log.Printf("unable to represent floating-point constant %v of type %v exactly; please submit a bug report to llir/llvm with this error message", c.X, c.Typ)
		}
		a, b := f.Bits()
		return fmt.Sprintf("0x%c%016X%016X", hexPrefix, a, b)
	// ppc_fp128 (PowerPC double-double arithmetic)
	case types.FloatKindPPC_FP128:
		// always represent ppc_fp128 in hexadecimal floating-point notation.
		const hexPrefix = 'M'
		if c.NaN {
			a, b := float128ppc.NaN.Bits()
			if c.X != nil && c.X.Signbit() {
				a, b = float128ppc.NegNaN.Bits()
			}
			return fmt.Sprintf("0x%c%016X%016X", hexPrefix, a, b)
		}
		f, acc := float128ppc.NewFromBig(c.X)
		if acc != big.Exact {
			log.Printf("unable to represent floating-point constant %v of type %v exactly; please submit a bug report to llir/llvm with this error message", c.X, c.Typ)
		}
		a, b := f.Bits()
		return fmt.Sprintf("0x%c%016X%016X", hexPrefix, a, b)
	default:
		panic(fmt.Errorf("support for floating-point kind %v not yet implemented", c.Typ.Kind))
	}
	// Insert decimal point if not present.
	//    3e4 -> 3.0e4
	//    42  -> 42.0
	s := c.X.Text('g', -1)
	if !strings.ContainsRune(s, '.') {
		if pos := strings.IndexByte(s, 'e'); pos != -1 {
			s = s[:pos] + ".0" + s[pos:]
		} else {
			s += ".0"
		}
	}
	return s
}
