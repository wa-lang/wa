// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "sort"

// A Type represents a type of Go.
// All types implement the Type interface.
type Type interface {
	// Underlying returns the underlying type of a type.
	Underlying() Type

	// String returns a string representation of a type.
	String() string
}

// BasicKind describes the kind of basic type.
type BasicKind int

const (
	Invalid BasicKind = iota // type is invalid

	// predeclared types
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String
	UnsafePointer

	// types for untyped values
	UntypedBool
	UntypedInt
	UntypedRune
	UntypedFloat
	UntypedComplex
	UntypedString
	UntypedNil

	// aliases
	Byte = Uint8
	Rune = Int32
)

// BasicInfo is a set of flags describing properties of a basic type.
type BasicInfo int

// Properties of basic types.
const (
	IsBoolean BasicInfo = 1 << iota
	IsInteger
	IsUnsigned
	IsFloat
	IsComplex
	IsString
	IsUntyped

	IsOrdered   = IsInteger | IsFloat | IsString
	IsNumeric   = IsInteger | IsFloat | IsComplex
	IsConstType = IsBoolean | IsNumeric | IsString
)

// A Basic represents a basic type.
type Basic struct {
	kind BasicKind
	info BasicInfo
	name string
}

// Kind returns the kind of basic type b.
func (b *Basic) Kind() BasicKind { return b.kind }

// Info returns information about properties of basic type b.
func (b *Basic) Info() BasicInfo { return b.info }

// Name returns the name of basic type b.
func (b *Basic) Name() string { return b.name }

// An Array represents an array type.
type Array struct {
	len  int64
	elem Type
}

// NewArray returns a new array type for the given element type and length.
// A negative length indicates an unknown length.
func NewArray(elem Type, len int64) *Array { return &Array{len, elem} }

// Len returns the length of array a.
// A negative result indicates an unknown length.
func (a *Array) Len() int64 { return a.len }

// Elem returns element type of array a.
func (a *Array) Elem() Type { return a.elem }

// A Slice represents a slice type.
type Slice struct {
	elem Type
}

// NewSlice returns a new slice type for the given element type.
func NewSlice(elem Type) *Slice { return &Slice{elem} }

// Elem returns the element type of slice s.
func (s *Slice) Elem() Type { return s.elem }

// A Struct represents a struct type.
type Struct struct {
	fields []*Var
	tags   []string // field tags; nil if there are no tags
}

// NewStruct returns a new struct with the given fields and corresponding field tags.
// If a field with index i has a tag, tags[i] must be that tag, but len(tags) may be
// only as long as required to hold the tag with the largest index i. Consequently,
// if no field has a tag, tags may be nil.
func NewStruct(fields []*Var, tags []string) *Struct {
	var fset objset
	for _, f := range fields {
		if f.name != "_" && fset.insert(f) != nil {
			panic("multiple fields with the same name")
		}
	}
	if len(tags) > len(fields) {
		panic("more tags than fields")
	}
	return &Struct{fields: fields, tags: tags}
}

// NumFields returns the number of fields in the struct (including blank and embedded fields).
func (s *Struct) NumFields() int { return len(s.fields) }

// Field returns the i'th field for 0 <= i < NumFields().
func (s *Struct) Field(i int) *Var { return s.fields[i] }

// Tag returns the i'th field tag for 0 <= i < NumFields().
func (s *Struct) Tag(i int) string {
	if i < len(s.tags) {
		return s.tags[i]
	}
	return ""
}

// A Pointer represents a pointer type.
type Pointer struct {
	base Type // element type
}

// NewPointer returns a new pointer type for the given element (base) type.
func NewPointer(elem Type) *Pointer { return &Pointer{base: elem} }

// Elem returns the element type for the given pointer p.
func (p *Pointer) Elem() Type { return p.base }

// A Tuple represents an ordered list of variables; a nil *Tuple is a valid (empty) tuple.
// Tuples are used as components of signatures and to represent the type of multiple
// assignments; they are not first class types of Go.
type Tuple struct {
	vars []*Var
}

// NewTuple returns a new tuple for the given variables.
func NewTuple(x ...*Var) *Tuple {
	if len(x) > 0 {
		return &Tuple{x}
	}
	return nil
}

// Len returns the number variables of tuple t.
func (t *Tuple) Len() int {
	if t != nil {
		return len(t.vars)
	}
	return 0
}

// At returns the i'th variable of tuple t.
func (t *Tuple) At(i int) *Var { return t.vars[i] }

// An Interface represents an interface type.
type Interface struct {
	methods   []*Func // ordered list of explicitly declared methods
	embeddeds []Type  // ordered list of explicitly embedded types

	allMethods []*Func // ordered list of methods declared with or embedded in this interface (TODO(gri): replace with mset)
}

// emptyInterface represents the empty (completed) interface
var emptyInterface = Interface{allMethods: markComplete}

// markComplete is used to mark an empty interface as completely
// set up by setting the allMethods field to a non-nil empty slice.
var markComplete = make([]*Func, 0)

// NewInterface returns a new (incomplete) interface for the given methods and embedded types.
// Each embedded type must have an underlying type of interface type.
// NewInterface takes ownership of the provided methods and may modify their types by setting
// missing receivers. To compute the method set of the interface, Complete must be called.
//
// Deprecated: Use NewInterfaceType instead which allows any (even non-defined) interface types
// to be embedded. This is necessary for interfaces that embed alias type names referring to
// non-defined (literal) interface types.
func NewInterface(methods []*Func, embeddeds []*Named) *Interface {
	tnames := make([]Type, len(embeddeds))
	for i, t := range embeddeds {
		tnames[i] = t
	}
	return NewInterfaceType(methods, tnames)
}

// NewInterfaceType returns a new (incomplete) interface for the given methods and embedded types.
// Each embedded type must have an underlying type of interface type (this property is not
// verified for defined types, which may be in the process of being set up and which don't
// have a valid underlying type yet).
// NewInterfaceType takes ownership of the provided methods and may modify their types by setting
// missing receivers. To compute the method set of the interface, Complete must be called.
func NewInterfaceType(methods []*Func, embeddeds []Type) *Interface {
	typ := new(Interface)

	if len(methods) == 0 && len(embeddeds) == 0 {
		return typ
	}

	var mset objset
	for _, m := range methods {
		if mset.insert(m) != nil {
			panic("multiple methods with the same name")
		}
		// set receiver if we don't have one
		if sig := m.typ.(*Signature); sig.recv == nil {
			sig.recv = NewVar(m.pos, m.pkg, "", typ)
		}
	}
	sort.Sort(byUniqueMethodName(methods))

	if len(embeddeds) > 0 {
		// All embedded types should be interfaces; however, defined types
		// may not yet be fully resolved. Only verify that non-defined types
		// are interfaces. This matches the behavior of the code before the
		// fix for #25301 (issue #25596).
		for _, t := range embeddeds {
			if _, ok := t.(*Named); !ok && !IsInterface(t) {
				panic("embedded type is not an interface")
			}
		}
		sort.Stable(byUniqueTypeName(embeddeds))
	}

	typ.methods = methods
	typ.embeddeds = embeddeds
	return typ
}

// NumExplicitMethods returns the number of explicitly declared methods of interface t.
func (t *Interface) NumExplicitMethods() int { return len(t.methods) }

// ExplicitMethod returns the i'th explicitly declared method of interface t for 0 <= i < t.NumExplicitMethods().
// The methods are ordered by their unique Id.
func (t *Interface) ExplicitMethod(i int) *Func { return t.methods[i] }

// NumEmbeddeds returns the number of embedded types in interface t.
func (t *Interface) NumEmbeddeds() int { return len(t.embeddeds) }

// Embedded returns the i'th embedded defined (*Named) type of interface t for 0 <= i < t.NumEmbeddeds().
// The result is nil if the i'th embedded type is not a defined type.
//
// Deprecated: Use EmbeddedType which is not restricted to defined (*Named) types.
func (t *Interface) Embedded(i int) *Named { tname, _ := t.embeddeds[i].(*Named); return tname }

// EmbeddedType returns the i'th embedded type of interface t for 0 <= i < t.NumEmbeddeds().
func (t *Interface) EmbeddedType(i int) Type { return t.embeddeds[i] }

// NumMethods returns the total number of methods of interface t.
func (t *Interface) NumMethods() int { return len(t.allMethods) }

// Method returns the i'th method of interface t for 0 <= i < t.NumMethods().
// The methods are ordered by their unique Id.
func (t *Interface) Method(i int) *Func { return t.allMethods[i] }

// Empty reports whether t is the empty interface.
func (t *Interface) Empty() bool { return len(t.allMethods) == 0 }

// Complete computes the interface's method set. It must be called by users of
// NewInterfaceType and NewInterface after the interface's embedded types are
// fully defined and before using the interface type in any way other than to
// form other types. Complete returns the receiver.
func (t *Interface) Complete() *Interface {
	if t.allMethods != nil {
		return t
	}

	// collect all methods
	var allMethods []*Func
	allMethods = append(allMethods, t.methods...)
	for _, et := range t.embeddeds {
		it := et.Underlying().(*Interface)
		it.Complete()
		// copy embedded methods unchanged (see issue #28282)
		allMethods = append(allMethods, it.allMethods...)
	}
	sort.Sort(byUniqueMethodName(allMethods))

	// t.methods and/or t.embeddeds may have been empty
	if allMethods == nil {
		allMethods = markComplete
	}
	t.allMethods = allMethods

	return t
}

// A Map represents a map type.
type Map struct {
	key, elem Type
}

// NewMap returns a new map for the given key and element types.
func NewMap(key, elem Type) *Map {
	return &Map{key, elem}
}

// Key returns the key type of map m.
func (m *Map) Key() Type { return m.key }

// Elem returns the element type of map m.
func (m *Map) Elem() Type { return m.elem }

// Implementations for Type methods.

func (b *Basic) Underlying() Type     { return b }
func (a *Array) Underlying() Type     { return a }
func (s *Slice) Underlying() Type     { return s }
func (s *Struct) Underlying() Type    { return s }
func (p *Pointer) Underlying() Type   { return p }
func (t *Tuple) Underlying() Type     { return t }
func (t *Interface) Underlying() Type { return t }
func (m *Map) Underlying() Type       { return m }

func (b *Basic) String() string     { return TypeString(b, nil) }
func (a *Array) String() string     { return TypeString(a, nil) }
func (s *Slice) String() string     { return TypeString(s, nil) }
func (s *Struct) String() string    { return TypeString(s, nil) }
func (p *Pointer) String() string   { return TypeString(p, nil) }
func (t *Tuple) String() string     { return TypeString(t, nil) }
func (t *Interface) String() string { return TypeString(t, nil) }
func (m *Map) String() string       { return TypeString(m, nil) }
