package ctypes

import (
	"strconv"

	"wa-lang.org/wa/internal/logger"
)

/**************************************
Tuple:
**************************************/
type Tuple struct {
	Vars []Type
}

func NewTuple(v []Type) *Tuple {
	return &Tuple{Vars: v}
}

func (t *Tuple) String() string {
	return t.CIRString()
}

func (t *Tuple) CIRString() string {
	s := "$TupleStart_"
	for _, n := range t.Vars {
		s += n.Name()
		s += "_"
	}
	s += "$TupleEnd"
	return s
}

func (t *Tuple) Name() string {
	return t.CIRString()
}

func (t *Tuple) Equal(u Type) bool {
	if u, ok := u.(*Tuple); ok {
		return t.CIRString() == u.CIRString()
	}
	return false
}

func (t *Tuple) GenStruct() Struct {
	var m []Field
	for i, v := range t.Vars {
		m = append(m, *NewField("$m"+strconv.Itoa(i), v))
	}
	return *NewStruct(t.CIRString(), m)
}

/**************************************
Struct:
**************************************/
type Field struct {
	name string
	typ  Type
}

func (i *Field) CIRString() string {
	return i.name
}

func NewField(n string, t Type) *Field {
	return &Field{name: n, typ: t}
}

func (i *Field) Type() Type {
	return i.typ
}

type Struct struct {
	name    string
	Members []Field
}

func NewStruct(name string, m []Field) *Struct {
	return &Struct{name: name, Members: m}
}

func (t *Struct) String() string {
	return t.CIRString()
}

func (t *Struct) CIRString() string {
	return t.name
}

func (t *Struct) Name() string {
	return t.CIRString()
}

func (t *Struct) Equal(u Type) bool {
	if u, ok := u.(*Struct); ok {
		return t.CIRString() == u.CIRString()
	}
	return false
}

/**************************************
Array:
**************************************/
type Array struct {
	len  int64
	elem Type
}

func NewArray(len int64, elem Type) *Array {
	return &Array{len: len, elem: elem}
}

func (t *Array) GetElem() Type {
	return t.elem
}

func (t *Array) GetLen() int64 {
	return t.len
}

func (t *Array) String() string {
	return t.CIRString()
}

func (t *Array) CIRString() string {
	return "$wartc::Array<" + t.elem.CIRString() + ", " + strconv.Itoa(int(t.len)) + ">"
}

func (t *Array) Name() string {
	return t.CIRString()
}

func (t *Array) Equal(u Type) bool {
	if u, ok := u.(*Array); ok {
		return t.CIRString() == u.CIRString()
	}
	return false
}

/**************************************
Slice:
**************************************/
type Slice struct {
	elem Type
}

func NewSlice(elem Type) *Slice {
	return &Slice{elem: elem}
}

func (t *Slice) GetElem() Type {
	return t.elem
}

func (t *Slice) String() string {
	return t.CIRString()
}

func (t *Slice) CIRString() string {
	return "$wartc::Slice<" + t.elem.CIRString() + ">"
}

func (t *Slice) Name() string {
	logger.Fatal("Todo: Slice.Name()")
	return ""
}

func (t *Slice) Equal(u Type) bool {
	if u, ok := u.(*Slice); ok {
		return t.elem.Equal(u.elem)
	}
	return false
}
