package wtypes

import "fmt"

type ValueType interface {
	fmt.Stringer
	Name() string

	GetByteSize() int
	Raw() []ValueType
	Equal(ValueType) bool
}

/**************************************
Void:
**************************************/
type Void struct {
}

func (t Void) String() string   { return t.Name() }
func (t Void) Name() string     { return "void" }
func (t Void) GetByteSize() int { return 0 }
func (t Void) Raw() []ValueType { return append([]ValueType(nil), t) }
func (t Void) Equal(u ValueType) bool {
	if _, ok := u.(Void); ok {
		return true
	}
	return false
}

/**************************************
Int32:
**************************************/
type Int32 struct {
}

func (t Int32) String() string   { return t.Name() }
func (t Int32) Name() string     { return "i32" }
func (t Int32) GetByteSize() int { return 4 }
func (t Int32) Raw() []ValueType { return append([]ValueType(nil), t) }
func (t Int32) Equal(u ValueType) bool {
	if _, ok := u.(Int32); ok {
		return true
	}
	return false
}

/**************************************
Int64:
**************************************/
type Int64 struct {
}

func (t Int64) String() string   { return t.Name() }
func (t Int64) Name() string     { return "i64" }
func (t Int64) GetByteSize() int { return 8 }
func (t Int64) Raw() []ValueType { return append([]ValueType(nil), t) }
func (t Int64) Equal(u ValueType) bool {
	if _, ok := u.(Int64); ok {
		return true
	}
	return false
}

/**************************************
Pointer:
**************************************/
type Pointer struct {
	Base ValueType
}

func NewPointer(base ValueType) Pointer { return Pointer{Base: base} }
func (t Pointer) String() string        { return "*" + t.Base.Name() }
func (t Pointer) Name() string          { return "Todo" }
func (t Pointer) GetByteSize() int      { return 4 }
func (t Pointer) Raw() []ValueType      { return append([]ValueType(nil), Int32{}) }
func (t Pointer) Equal(u ValueType) bool {
	if ut, ok := u.(Pointer); ok {
		return t.Base.Equal(ut.Base)
	}
	return false
}
