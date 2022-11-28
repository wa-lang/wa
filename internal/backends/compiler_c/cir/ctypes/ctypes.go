//ctypes包定义C-IR的类型
package ctypes

import (
	"fmt"
	"strings"

	"wa-lang.org/wa/internal/logger"
)

/**************************************
Type: C-IR类型接口
以下类型实现了 Type 接口：
   *types.VoidType
   *types.FuncType
   *types.IntType
   *types.FloatType
**************************************/
type Type interface {
	fmt.Stringer
	// CIRString 返回该类型的CIR表达
	CIRString() string
	// Name 返回该类型的名字，例如字符串类型的Name()=="string"，而其CIRString()=="$wartc::string"
	Name() string
	// Equal 判断u是否与当前类型相同
	Equal(u Type) bool
}

var (
	Void   = &VoidType{}
	Bool   = &BoolType{}
	Int8   = &IntType{Kind: IntKindI8}
	Uint8  = &IntType{Kind: IntKindU8}
	Int16  = &IntType{Kind: IntKindI16}
	Uint16 = &IntType{Kind: IntKindU16}
	Int32  = &IntType{Kind: IntKindI32}
	Uint32 = &IntType{Kind: IntKindU32}
	Int64  = &IntType{Kind: IntKindI64}
	Uint64 = &IntType{Kind: IntKindU64}
	Float  = &FloatType{Kind: FloatKindFloat}
	Double = &FloatType{Kind: FloatKindDouble}
)

/**************************************
VoidType:
**************************************/
type VoidType struct {
}

// String 返回该类型的字符串表达
func (t *VoidType) String() string {
	return t.Name()
}

// CIRString 返回该类型的CIR语法表达
func (t *VoidType) CIRString() string {
	return "void"
}

func (t *VoidType) Name() string {
	return "void"
}

// Equal 判断u是否与当前类型相同
func (t *VoidType) Equal(u Type) bool {
	if _, ok := u.(*VoidType); ok {
		return true
	}
	return false
}

/**************************************
BoolType:
**************************************/
type BoolType struct {
}

// String 返回该类型的字符串表达
func (t *BoolType) String() string {
	return t.Name()
}

// CIRString 返回该类型的CIR语法表达
func (t *BoolType) CIRString() string {
	return "bool"
}

func (t *BoolType) Name() string {
	return "bool"
}

// Equal 判断u是否与当前类型相同
func (t *BoolType) Equal(u Type) bool {
	if _, ok := u.(*BoolType); ok {
		return true
	}
	return false
}

/**************************************
FuncType:
**************************************/
type FuncType struct {
	// 返回值类型
	Ret Type
	// 参数
	Params []Type
}

// NewFunc 根据指定的返回值和参数类型创建函数类型
func NewFuncType(ret Type, params []Type) *FuncType {
	return &FuncType{
		Ret:    ret,
		Params: params,
	}
}

// String 返回该类型的字符串表达
func (t *FuncType) String() string {
	return t.CIRString()
}

// CIRString 返回该类型的CIR语法表达
func (t *FuncType) CIRString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s (", t.Ret.CIRString())
	for i, param := range t.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.String())
	}

	buf.WriteString(")")
	return buf.String()
}

func (t *FuncType) Name() string {
	logger.Fatal("Todo: FuncType.Name()")
	return ""
}

// Equal 判断u是否与当前类型相同
func (t *FuncType) Equal(u Type) bool {
	if u, ok := u.(*FuncType); ok {
		if !t.Ret.Equal(u.Ret) {
			return false
		}
		if len(t.Params) != len(u.Params) {
			return false
		}
		for i := range t.Params {
			if !t.Params[i].Equal(u.Params[i]) {
				return false
			}
		}
		return true
	}
	return false
}

/**************************************
IntType:
**************************************/
type IntType struct {
	// 整数类型
	Kind IntKind
}

// IntKind: 整数类型
type IntKind uint8

const (
	// 8位有符号整数
	IntKindI8 IntKind = iota
	// 8位无符号整数
	IntKindU8
	// 16位有符号整数
	IntKindI16
	// 16位无符号整数
	IntKindU16
	// 32位有符号整数
	IntKindI32
	// 32位无符号整数
	IntKindU32
	// 64位有符号整数
	IntKindI64
	// 64位无符号整数
	IntKindU64
)

var intKindName = [...][2]string{
	IntKindI8:  {"int8_t", "int8"},
	IntKindU8:  {"uint8_t", "uint8"},
	IntKindI16: {"int16_t", "int16"},
	IntKindU16: {"uint16_t", "uint16"},
	IntKindI32: {"int32_t", "int32"},
	IntKindU32: {"uint32_t", "uint32"},
	IntKindI64: {"int64_t", "int64"},
	IntKindU64: {"uint64_t", "uint64"},
}

func (i IntKind) cirString() string {
	if int(i) >= len(intKindName) {
		panic("Invalid IntKind.")
	}
	return intKindName[i][0]
}

func (i IntKind) name() string {
	if int(i) >= len(intKindName) {
		panic("Invalid IntKind.")
	}
	return intKindName[i][1]
}

// GetBitSize() 返回整数的位宽
func (i IntKind) GetBitSize() int {
	_IntKindSize := [...]int{
		IntKindI8:  8,
		IntKindU8:  8,
		IntKindI16: 16,
		IntKindU16: 16,
		IntKindI32: 32,
		IntKindU32: 32,
		IntKindI64: 64,
		IntKindU64: 64,
	}
	if int(i) >= len(_IntKindSize) {
		panic("Invalid IntKind.")
	}
	return _IntKindSize[i]
}

// NewInt 根据给定参数创建整数类型
func NewInt(kind IntKind) *IntType {
	return &IntType{
		Kind: kind,
	}
}

// String 返回该类型的字符串表达
func (t *IntType) String() string {
	return t.Name()
}

// CIRString 返回该类型的CIR语法表达
func (t *IntType) CIRString() string {
	return t.Kind.cirString()
}

func (t *IntType) Name() string {
	return t.Kind.name()
}

// Equal 判断u是否与当前类型相同
func (t *IntType) Equal(u Type) bool {
	if u, ok := u.(*IntType); ok {
		return u.Kind == t.Kind
	}
	return false
}

/*-----------------------------------*/
// FloatType: C-IR的浮点数类型
type FloatType struct {
	// 浮点数类型
	Kind FloatKind
}

// FloatKind: 浮点数类型
type FloatKind uint8

const (
	// 32位浮点数，IEEE 754 单精度
	FloatKindFloat FloatKind = iota
	// 64位浮点数，IEEE 754 双精度
	FloatKindDouble
)

func (i FloatKind) String() string {
	_FloatKindName := [...]string{
		FloatKindFloat:  "float",
		FloatKindDouble: "double",
	}
	if i >= FloatKind(len(_FloatKindName)) {
		panic("Invalid FloatKind.")
	}
	return _FloatKindName[i]
}

// NewFloat 根据给定参数创建浮点数类型
func NewFloat(kind FloatKind) *FloatType {
	return &FloatType{
		Kind: kind,
	}
}

// String 返回该类型的字符串表达
func (t *FloatType) String() string {
	return t.Name()
}

// CIRString 返回该类型的CIR语法表达
func (t *FloatType) CIRString() string {
	return t.Kind.String()
}

func (t *FloatType) Name() string {
	return t.Kind.String()
}

// Equal 判断u是否与当前类型相同
func (t *FloatType) Equal(u Type) bool {
	if u, ok := u.(*FloatType); ok {
		return u.Kind == t.Kind
	}
	return false
}

/**************************************
PointerType:
**************************************/
type PointerType struct {
	Base Type
}

func NewPointerType(base Type) *PointerType {
	return &PointerType{Base: base}
}

// String 返回该类型的字符串表达
func (t *PointerType) String() string {
	return t.CIRString()
}

// CIRString 返回该类型的CIR语法表达
func (t *PointerType) CIRString() string {
	return t.Base.CIRString() + "*"
}

func (t *PointerType) Name() string {
	logger.Fatal("Todo: PointerType.Name()")
	return ""
}

// Equal 判断u是否与当前类型相同
func (t *PointerType) Equal(u Type) bool {
	if ut, ok := u.(*PointerType); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

/**************************************
RefType:
**************************************/
type RefType struct {
	Base Type
}

func NewRefType(base Type) *RefType {
	return &RefType{Base: base}
}

// String 返回该类型的字符串表达
func (t *RefType) String() string {
	return t.CIRString()
}

// CIRString 返回该类型的CIR语法表达
func (t *RefType) CIRString() string {
	return "$wartc::Ref<" + t.Base.CIRString() + ">"
}

func (t *RefType) Name() string {
	logger.Fatal("Todo: RefType.Name()")
	return ""
}

// Equal 判断u是否与当前类型相同
func (t *RefType) Equal(u Type) bool {
	if ut, ok := u.(*RefType); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}

/**************************************
StringLit:
**************************************/
type StringLit struct {
}

// String 返回该类型的字符串表达
func (t *StringLit) String() string {
	return t.CIRString()
}

// CIRString 返回该类型的CIR语法表达
func (t *StringLit) CIRString() string {
	//字符串字面值类型没有CIR表达
	logger.Fatal("Unreachable")
	return ""
}

func (t *StringLit) Name() string {
	return "StringLit"
}

// Equal 判断u是否与当前类型相同
func (t *StringLit) Equal(u Type) bool {
	if _, ok := u.(*StringLit); ok {
		return true
	}
	return false
}

/**************************************
String:
**************************************/
type StringVar struct {
}

// String 返回该类型的字符串表达
func (t *StringVar) String() string {
	return t.CIRString()
}

// CIRString 返回该类型的CIR语法表达
func (t *StringVar) CIRString() string {
	return "$wartc::String"
}

func (t *StringVar) Name() string {
	return "string"
}

// Equal 判断u是否与当前类型相同
func (t *StringVar) Equal(u Type) bool {
	if _, ok := u.(*StringVar); ok {
		return true
	}
	return false
}
