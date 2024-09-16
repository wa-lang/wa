// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file sets up the universe scope and the unsafe package.

package types

import (
	"strings"

	"wa-lang.org/wa/internal/constant"
	"wa-lang.org/wa/internal/token"
)

// The Universe scope contains all predeclared objects of Go.
// It is the outermost scope of any chain of nested scopes.
var Universe *Scope

// The Unsafe package is the package returned by an importer
// for the import path "unsafe".
var Unsafe *Package

var (
	universeIota   *Const
	universeByte   *Basic // uint8 alias, but has name "byte"
	universeRune   *Basic // int32 alias, but has name "rune"
	universeString *Basic
	universeAny    Object

	universe__PACKAGE__ *Const
	universe__FILE__    *Const
	universe__LINE__    *Const
	universe__COLUMN__  *Const
	universe__FUNC__    *Const
)

// Typ contains the predeclared *Basic types indexed by their
// corresponding BasicKind.
//
// The *Basic type for Typ[Byte] will have the name "uint8".
// Use Universe.Lookup("byte").Type() to obtain the specific
// alias basic type named "byte" (and analogous for "rune").
var Typ = []*Basic{
	Invalid: {Invalid, 0, "invalid type"},

	Bool:          {Bool, IsBoolean, "bool"},
	Int:           {Int, IsInteger, "int"},
	Int8:          {Int8, IsInteger, "__wa_int8"},
	Int16:         {Int16, IsInteger, "__wa_int16"},
	Int32:         {Int32, IsInteger, "int32"},
	Int64:         {Int64, IsInteger, "int64"},
	Uint:          {Uint, IsInteger | IsUnsigned, "uint"},
	Uint8:         {Uint8, IsInteger | IsUnsigned, "uint8"},
	Uint16:        {Uint16, IsInteger | IsUnsigned, "uint16"},
	Uint32:        {Uint32, IsInteger | IsUnsigned, "uint32"},
	Uint64:        {Uint64, IsInteger | IsUnsigned, "uint64"},
	Uintptr:       {Uintptr, IsInteger | IsUnsigned, "uintptr"},
	Float32:       {Float32, IsFloat, "float32"},
	Float64:       {Float64, IsFloat, "float64"},
	Complex64:     {Complex64, IsComplex, "complex64"},
	Complex128:    {Complex128, IsComplex, "complex128"},
	String:        {String, IsString, "string"},
	UnsafePointer: {UnsafePointer, 0, "Pointer"},

	UntypedBool:    {UntypedBool, IsBoolean | IsUntyped, "untyped bool"},
	UntypedInt:     {UntypedInt, IsInteger | IsUntyped, "untyped int"},
	UntypedRune:    {UntypedRune, IsInteger | IsUntyped, "untyped rune"},
	UntypedFloat:   {UntypedFloat, IsFloat | IsUntyped, "untyped float"},
	UntypedComplex: {UntypedComplex, IsComplex | IsUntyped, "untyped complex"},
	UntypedString:  {UntypedString, IsString | IsUntyped, "untyped string"},
	UntypedNil:     {UntypedNil, IsUntyped, "untyped nil"},
}

var aliases = [...]*Basic{
	{Byte, IsInteger | IsUnsigned, "byte"},
	{Byte, IsInteger | IsUnsigned, "字"},
	{Rune, IsInteger, "rune"},

	{Int8, IsInteger, "__wa_i8"},
	{Int16, IsInteger, "__wa_i16"},
	{Int32, IsInteger, "i32"},
	{Int64, IsInteger, "i64"},
	{Int, IsInteger, "数"},

	{Uint8, IsInteger | IsUnsigned, "u8"},
	{Uint16, IsInteger | IsUnsigned, "u16"},
	{Uint32, IsInteger | IsUnsigned, "u32"},
	{Uint64, IsInteger | IsUnsigned, "u64"},
	{Uint64, IsInteger | IsUnsigned, "u128"},
	{Uint64, IsInteger | IsUnsigned, "u256"},

	{Float32, IsFloat, "f32"},
	{Float64, IsFloat, "f64"},

	{Float32, IsFloat, "float"},
	{Float64, IsFloat, "double"},

	{String, IsString, "文"},
}

func defPredeclaredTypes() {
	for _, t := range Typ {
		def(NewTypeName(token.NoPos, nil, t.name, t))
	}
	for _, t := range aliases {
		def(NewTypeName(token.NoPos, nil, t.name, t))
	}

	// type any = interface{}
	// Note: don't use &emptyInterface for the type of any. Using a unique
	// pointer allows us to detect any and format it as "any" rather than
	// interface{}, which clarifies user-facing error messages significantly.
	def(NewTypeName(token.NoPos, nil, "any", &Interface{}))

	// Error has a nil package in its qualified name since it is in no package
	res := NewVar(token.NoPos, nil, "", Typ[String])
	sig := &Signature{results: NewTuple(res)}
	err := NewFunc(token.NoPos, nil, "Error", sig)
	typ := &Named{underlying: NewInterfaceType([]*Func{err}, nil).Complete()}
	sig.recv = NewVar(token.NoPos, nil, "", typ)
	def(NewTypeName(token.NoPos, nil, "error", typ))
}

var predeclaredConsts = [...]struct {
	name string
	kind BasicKind
	val  constant.Value
}{
	{"true", UntypedBool, constant.MakeBool(true)},
	{"false", UntypedBool, constant.MakeBool(false)},
	{"iota", UntypedInt, constant.MakeInt64(0)},

	{"__PACKAGE__", UntypedString, constant.MakeString("<wa-lang:__PACKAGE__>")},
	{"__FILE__", UntypedString, constant.MakeString("<wa-lang:__FILE__>")},
	{"__LINE__", UntypedInt, constant.MakeInt64(0)},
	{"__COLUMN__", UntypedInt, constant.MakeInt64(0)},
	{"__FUNC__", UntypedString, constant.MakeString("")},
}

func defPredeclaredConsts() {
	for _, c := range predeclaredConsts {
		def(NewConst(token.NoPos, nil, c.name, Typ[c.kind], c.val))
	}
}

func defPredeclaredNil() {
	def(&Nil{object{name: "nil", typ: Typ[UntypedNil], color_: black}})
}

// A builtinId is the id of a builtin function.
type builtinId int

const (
	// universe scope
	_Append builtinId = iota
	_Cap
	_Complex
	_Copy
	_Delete
	_Imag
	_Len
	_Make
	_New
	_Panic
	_Print
	_Println
	_Real
	_Recover

	// wa
	_Raw
	_SetFinalizer
	_Printf

	// wz
	_长

	// package unsafe
	_Alignof
	_Offsetof
	_Sizeof

	// testing support
	_Assert
	_Trace
)

var predeclaredFuncs = [...]struct {
	name     string
	nargs    int
	variadic bool
	kind     exprKind
}{
	_Append:  {"append", 1, true, expression},
	_Cap:     {"cap", 1, false, expression},
	_Complex: {"complex", 2, false, expression},
	_Copy:    {"copy", 2, false, statement},
	_Delete:  {"delete", 2, false, statement},
	_Imag:    {"imag", 1, false, expression},
	_Len:     {"len", 1, false, expression},
	_Make:    {"make", 1, true, expression},
	_New:     {"new", 1, false, expression},
	_Panic:   {"panic", 1, false, statement},
	_Print:   {"print", 0, true, statement},
	_Println: {"println", 0, true, statement},
	_Real:    {"real", 1, false, expression},
	_Recover: {"recover", 0, false, statement},

	_Raw:          {"raw", 1, false, expression},
	_SetFinalizer: {"setFinalizer", 2, false, statement},
	_Printf:       {"printf", 1, true, statement},

	_长: {"长", 1, false, expression},

	_Alignof:  {"Alignof", 1, false, expression},
	_Offsetof: {"Offsetof", 1, false, expression},
	_Sizeof:   {"Sizeof", 1, false, expression},

	_Assert: {"assert", 1, true, statement},
	_Trace:  {"trace", 0, true, statement},
}

func defPredeclaredFuncs() {
	for i := range predeclaredFuncs {
		id := builtinId(i)
		if id == _Assert || id == _Trace {
			continue // only define these in testing environment
		}
		def(newBuiltin(id))
	}
}

// DefPredeclaredTestFuncs defines the assert and trace built-ins.
// These built-ins are intended for debugging and testing of this
// package only.
func DefPredeclaredTestFuncs() {
	if Universe.Lookup("assert") != nil {
		return // already defined
	}
	def(newBuiltin(_Assert))
	def(newBuiltin(_Trace))
}

func init() {
	Universe = NewScope(nil, token.NoPos, token.NoPos, "universe")
	Unsafe = NewPackage("unsafe", "unsafe")
	Unsafe.complete = true

	defPredeclaredTypes()
	defPredeclaredConsts()
	defPredeclaredNil()
	defPredeclaredFuncs()

	universeIota = Universe.Lookup("iota").(*Const)
	universeByte = Universe.Lookup("byte").(*TypeName).typ.(*Basic)
	universeRune = Universe.Lookup("rune").(*TypeName).typ.(*Basic)
	universeString = Universe.Lookup("string").(*TypeName).typ.(*Basic)
	universeAny = Universe.Lookup("any")

	universe__PACKAGE__ = Universe.Lookup("__PACKAGE__").(*Const)
	universe__FILE__ = Universe.Lookup("__FILE__").(*Const)
	universe__LINE__ = Universe.Lookup("__LINE__").(*Const)
	universe__COLUMN__ = Universe.Lookup("__COLUMN__").(*Const)
	universe__FUNC__ = Universe.Lookup("__FUNC__").(*Const)
}

// Objects with names containing blanks are internal and not entered into
// a scope. Objects with exported names are inserted in the unsafe package
// scope; other objects are inserted in the universe scope.
func isCNKeyword(name string) bool {
	return name == "数" || name == "文" || name == "字" || name == "长"
}

func def(obj Object) {
	assert(obj.color() == black)
	name := obj.Name()
	if strings.Contains(name, " ") {
		return // nothing to do
	}
	// fix Obj link for named types
	if typ, ok := obj.Type().(*Named); ok {
		typ.obj = obj.(*TypeName)
	}
	// exported identifiers go into package unsafe
	scope := Universe
	// 这里中文的关键字被obj.Exported给过滤了，所以被放进了Unsafe的scope
	if !isCNKeyword(obj.Name()) && obj.Exported() {
		scope = Unsafe.scope
		// set Pkg field
		switch obj := obj.(type) {
		case *TypeName:
			obj.pkg = Unsafe
		case *Builtin:
			obj.pkg = Unsafe
		default:
			unreachable()
		}
	}
	if scope.Insert(obj) != nil {
		panic("internal error: double declaration")
	}
}
