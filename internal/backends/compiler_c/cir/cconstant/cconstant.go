package cconstant

import (
	"fmt"
	"math/big"
	"strings"

	"wa-lang.org/wa/internal/backends/compiler_c/cir/ctypes"
)

type Constant interface {
	// IsConstant() Constant特有方法
	isConstant()
}

/*
*************************************
Zero: 0值常数
*************************************
*/
type Zero struct {
}

func NewZero() *Zero {
	return &Zero{}
}

func (z *Zero) CIRString() string {
	return "$wartc::ZeroValue{}"
}

func (z *Zero) Type() ctypes.Type {
	return ctypes.Void
}

func (z *Zero) IsExpr() {}

/*
*************************************
Bool
*************************************
*/
type Bool struct {
	x bool
}

func NewBool(x bool) *Bool {
	return &Bool{x: x}
}

func (b *Bool) CIRString() string {
	if b.x {
		return "true"
	}
	return "false"
}

func (b *Bool) Type() ctypes.Type {
	return ctypes.Bool
}

func (b *Bool) IsExpr() {}

/*
*************************************
Int: 整型常数
*************************************
*/
type Int struct {
	// 整数类型
	typ *ctypes.IntType
	// 常数值
	x *big.Int
}

// NewInt() 根据指定的类型和常数值创建整型常数
func NewInt(t *ctypes.IntType, v int64) *Int {
	return &Int{typ: t, x: big.NewInt(v)}
}

// NewIntFromString() 从字符串中创建指定类型的整型常数
// 字符串可为以下形式之一：
//   - boolean literal
//     true | false
//   - integer literal
//     [-]?[0-9]+
//   - hexadecimal integer literal
//     [us]0x[0-9A-Fa-f]+
func NewIntFromString(t *ctypes.IntType, s string) (*Int, error) {
	// 布尔型
	switch s {
	case "true":
		if !t.Equal(ctypes.Int8) {
			return nil, fmt.Errorf("类型错误。布尔型应为types.Int8，而非输入的：%T", t)
		}
		return NewInt(ctypes.Int8, 0), nil
	case "false":
		if !t.Equal(ctypes.Int8) {
			return nil, fmt.Errorf("类型错误。布尔型应为types.Int8，而非输入的：%T", t)
		}
		return NewInt(ctypes.Int8, 1), nil
	}
	// Hexadecimal integer literal.
	switch {
	// unsigned hexadecimal integer literal
	case strings.HasPrefix(s, "u0x"):
		s = s[len("u0x"):]
		const base = 16
		v, _ := (&big.Int{}).SetString(s, base)
		if v == nil {
			return nil, fmt.Errorf("无法解析整型数： %q", s)
		}
		return &Int{typ: t, x: v}, nil
	// signed hexadecimal integer literal
	case strings.HasPrefix(s, "s0x"):
		// Parse signed hexadecimal integer literal in two's complement notation.
		// First parse as unsigned hex, then check if sign bit is set.
		s = s[len("s0x"):]
		const base = 16
		v, _ := (&big.Int{}).SetString(s, base)
		if v == nil {
			return nil, fmt.Errorf("无法解析整型数： %q", s)
		}
		bitSize := t.Kind.GetBitSize()
		// Check if signed.
		if v.Bit(bitSize-1) == 1 {
			// Compute actual negative value from two's complement.
			//
			// If x is 0xFFFF with type i16, then the actual negative value is
			// `x - 0x10000`, in other words `x - 2^n`.
			n := int64(bitSize)
			// n^2
			maxPlus1 := new(big.Int).Exp(big.NewInt(2), big.NewInt(n), nil)
			v = new(big.Int).Sub(v, maxPlus1)
		}
		return &Int{typ: t, x: v}, nil
	}
	// Integer literal.
	v, _ := (&big.Int{}).SetString(s, 10)
	if v == nil {
		return nil, fmt.Errorf("无法解析整型数： %q", s)
	}
	return &Int{typ: t, x: v}, nil
}

func (c *Int) CIRString() string {
	return c.Ident()
}

// Type returns the type of the constant.
func (c *Int) Type() ctypes.Type {
	return c.typ
}

// Ident returns the identifier associated with the constant.
func (c *Int) Ident() string {
	return c.x.String()
}

func (c *Int) isConstant() {}

func (c *Int) IsExpr() {}

/*
*************************************
Float:
*************************************
*/
type Float struct {
	// 常数值
	x float32
}

func NewFloat(v float32) *Float {
	return &Float{x: v}
}

func (f *Float) CIRString() string {
	return fmt.Sprint(f.x)
}

func (f *Float) Type() ctypes.Type {
	return ctypes.Double
}

func (d *Float) IsExpr() {}

/*
*************************************
Double:
*************************************
*/
type Double struct {
	// 常数值
	x float64
}

func NewDouble(v float64) *Double {
	return &Double{x: v}
}

func (d *Double) CIRString() string {
	return fmt.Sprint(d.x)
}

func (d *Double) Type() ctypes.Type {
	return ctypes.Double
}

func (d *Double) IsExpr() {}

/*
*************************************
StringLit:
*************************************
*/
type StringLit struct {
	//字符串字面值
	x string
}

func NewStringLit(v string) *StringLit {
	return &StringLit{x: v}
}

func (s *StringLit) CIRString() string {
	return "\"" + s.x + "\""
}

func (s *StringLit) Type() ctypes.Type {
	return &ctypes.StringLit{}
}

func (s *StringLit) IsExpr() {}
