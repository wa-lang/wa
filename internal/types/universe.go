// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file sets up the universe scope and the unsafe package.

package types

import (
	"strings"
	"unicode/utf8"

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
	universe__POS__     *Const
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
	{Byte, IsInteger | IsUnsigned, token.K_字节},
	{Byte, IsInteger | IsUnsigned, "字"},
	{Rune, IsInteger, "rune"},
	{Rune, IsInteger, token.K_符文},

	{Int8, IsInteger, "__wa_i8"},
	{Int16, IsInteger, "__wa_i16"},
	{Int32, IsInteger, "i32"},
	{Int64, IsInteger, "i64"},
	{Int, IsInteger, "数"},

	{Int8, IsInteger, token.K_微整型},
	{Int16, IsInteger, token.K_短整型},
	{Int32, IsInteger, token.K_普整型},
	{Int64, IsInteger, token.K_长整型},
	{Int, IsInteger, token.K_整型},

	{Uint8, IsInteger | IsUnsigned, "u8"},
	{Uint16, IsInteger | IsUnsigned, "u16"},
	{Uint32, IsInteger | IsUnsigned, "u32"},
	{Uint64, IsInteger | IsUnsigned, "u64"},

	{Uint8, IsInteger | IsUnsigned, token.K_微正整},
	{Uint16, IsInteger | IsUnsigned, token.K_短正整},
	{Uint32, IsInteger | IsUnsigned, token.K_普正整},
	{Uint64, IsInteger | IsUnsigned, token.K_长正整},

	{Float32, IsFloat, "f32"},
	{Float64, IsFloat, "f64"},

	{Float32, IsFloat, token.K_单精},
	{Float64, IsFloat, token.K_双精},

	{String, IsString, "文"},
	{String, IsString, token.K_字串},
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
	{"__POS__", UntypedInt, constant.MakeInt64(0)},
}

func defPredeclaredConsts() {
	for _, c := range predeclaredConsts {
		def(NewConst(token.NoPos, nil, c.name, Typ[c.kind], c.val))
	}
}

func defPredeclaredNil() {
	def(&Nil{object{name: "nil", typ: Typ[UntypedNil], color_: black}})
	def(&Nil{object{name: token.K_空, typ: Typ[UntypedNil], color_: black}})
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

	// w2
	_追加 // append
	_容量 // cap
	_复数 // complex
	_拷贝 // copy
	_删除 // delete
	_虚部 // imag
	_长度 // len
	_构建 // make
	_新建 // new
	_崩溃 // panic
	_输出 // println
	_打印 // print
	_实部 // real

	// wa
	_Raw
	_SetFinalizer

	// wz
	_长

	// package unsafe
	_unsafe_Alignof
	_unsafe_Offsetof
	_unsafe_Sizeof
	_unsafe_Raw

	// package runtime
	_runtime_SetFinalizer

	// testing support
	_Assert
	_Trace

	_断言
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
	_New:     {"new", 1, true, expression},
	_Panic:   {"panic", 1, false, statement},
	_Print:   {"print", 0, true, statement},
	_Println: {"println", 0, true, statement},
	_Real:    {"real", 1, false, expression},
	_Recover: {"recover", 0, false, statement},

	_追加: {token.K_追加, 1, true, expression},
	_容量: {token.K_容量, 1, false, expression},
	_复数: {token.K_复数, 2, false, expression},
	_拷贝: {token.K_拷贝, 2, false, statement},
	_删除: {token.K_删除, 2, false, statement},
	_虚部: {token.K_虚部, 1, false, expression},
	_长度: {token.K_长度, 1, false, expression},
	_构建: {token.K_构建, 1, true, expression},
	_新建: {token.K_新建, 1, true, expression},
	_崩溃: {token.K_崩溃, 1, false, statement},
	_输出: {token.K_输出, 0, true, statement},
	_打印: {token.K_打印, 0, true, statement},
	_实部: {token.K_实部, 1, false, expression},

	_Raw:        {"raw", 1, false, expression},
	_unsafe_Raw: {"Raw", 1, false, expression},

	_SetFinalizer:         {"setFinalizer", 2, false, statement},
	_runtime_SetFinalizer: {"SetFinalizer", 2, false, statement},

	_长: {"长", 1, false, expression},

	_unsafe_Alignof:  {"Alignof", 1, false, expression},
	_unsafe_Offsetof: {"Offsetof", 1, false, expression},
	_unsafe_Sizeof:   {"Sizeof", 1, false, expression},

	_Assert: {"assert", 1, true, statement},
	_Trace:  {"trace", 0, true, statement},

	_断言: {token.K_断言, 1, true, statement},
}

func defPredeclaredFuncs() {
	for i := range predeclaredFuncs {
		id := builtinId(i)
		if id == _runtime_SetFinalizer {
			continue // 在加载 runtime 包时导入
		}
		if id == _Assert || id == _Trace || id == _断言 {
			continue // only define these in testing environment
		}
		def(newBuiltin(id))
	}
}

// 注册运行时函数
func DefPredeclaredRuntimeFuncs(runtimePkg *Package) {
	if runtimePkg.scope.Lookup("SetFinalize") != nil {
		return // already defined
	}
	defInPackage(runtimePkg, newBuiltin(_runtime_SetFinalizer))
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

	def(newBuiltin(_断言))
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
	universe__POS__ = Universe.Lookup("__POS__").(*Const)
}

// Objects with names containing blanks are internal and not entered into
// a scope. Objects with exported names are inserted in the unsafe package
// scope; other objects are inserted in the universe scope.
func isCNKeyword(name string) bool {
	_, size := utf8.DecodeRuneInString(name)
	return size > 1
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
		panic("internal error: double declaration:" + obj.Name())
	}
}

func defInPackage(pkg *Package, obj Object) {
	assert(obj.color() == black)
	name := obj.Name()
	if strings.Contains(name, " ") {
		return // nothing to do
	}
	// fix Obj link for named types
	if typ, ok := obj.Type().(*Named); ok {
		typ.obj = obj.(*TypeName)
	}
	scope := pkg.Scope()
	switch obj := obj.(type) {
	case *TypeName:
		obj.pkg = pkg
	case *Builtin:
		obj.pkg = pkg
	default:
		unreachable()
	}
	if scope.Insert(obj) != nil {
		panic("internal error: double declaration")
	}
}
