// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package types

import (
	"strings"

	"wa-lang.org/wa/internal/constant"
	"wa-lang.org/wa/internal/token"
)

var (
	WaUniverse *Scope
	WaUnsafe   *Package
)

var (
	waUniverseIota   *Const
	waUniverseByte   *Basic // uint8 alias, but has name "byte"
	waUniverseRune   *Basic // int32 alias, but has name "rune"
	waUniverseString *Basic
	waUniverseAny    Object
	waUniverseError  Object

	waUniverse__PACKAGE__ *Const
	waUniverse__FILE__    *Const
	waUniverse__LINE__    *Const
	waUniverse__COLUMN__  *Const
	waUniverse__FUNC__    *Const
	waUniverse__POS__     *Const
)

var waAliases = [...]*Basic{
	{Byte, IsInteger | IsUnsigned, token.K_byte},
	{Rune, IsInteger, token.K_rune},

	{Int8, IsInteger, token.K_i8},
	{Int16, IsInteger, token.K_i16},
	{Int32, IsInteger, token.K_i32},
	{Int64, IsInteger, token.K_i64},

	{Uint8, IsInteger | IsUnsigned, token.K_u8},
	{Uint16, IsInteger | IsUnsigned, token.K_u16},
	{Uint32, IsInteger | IsUnsigned, token.K_u32},
	{Uint64, IsInteger | IsUnsigned, token.K_u64},

	{Float32, IsFloat, token.K_f32},
	{Float64, IsFloat, token.K_f64},
}

func waDefPredeclaredTypes() {
	for _, t := range Typ {
		waDef(NewTypeName(token.NoPos, nil, t.name, t))
	}
	for _, t := range waAliases {
		waDef(NewTypeName(token.NoPos, nil, t.name, t))
	}

	// type any = interface{}
	// Note: don't use &emptyInterface for the type of any. Using a unique
	// pointer allows us to detect any and format it as "any" rather than
	// interface{}, which clarifies user-facing error messages significantly.
	waDef(NewTypeName(token.NoPos, nil, token.K_any, &Interface{}))

	// Error has a nil package in its qualified name since it is in no package
	res := NewVar(token.NoPos, nil, "", Typ[String])
	sig := &Signature{results: NewTuple(res)}
	err := NewFunc(token.NoPos, nil, token.K_Error, sig)
	typ := &Named{underlying: NewInterfaceType([]*Func{err}, nil).Complete()}
	sig.recv = NewVar(token.NoPos, nil, "", typ)
	waDef(NewTypeName(token.NoPos, nil, token.K_error, typ))
}

var waPredeclaredConsts = [...]struct {
	name string
	kind BasicKind
	val  constant.Value
}{
	{token.K_true, UntypedBool, constant.MakeBool(true)},
	{token.K_false, UntypedBool, constant.MakeBool(false)},
	{token.K_iota, UntypedInt, constant.MakeInt64(0)},

	{token.K__PACKAGE__, UntypedString, constant.MakeString("<wa-lang:__PACKAGE__>")},
	{token.K__FILE__, UntypedString, constant.MakeString("<wa-lang:__FILE__>")},
	{token.K__LINE__, UntypedInt, constant.MakeInt64(0)},
	{token.K__COLUMN__, UntypedInt, constant.MakeInt64(0)},
	{token.K__FUNC__, UntypedString, constant.MakeString("")},
	{token.K__POS__, UntypedInt, constant.MakeInt64(0)},
}

func waDefPredeclaredConsts() {
	for _, c := range waPredeclaredConsts {
		waDef(NewConst(token.NoPos, nil, c.name, Typ[c.kind], c.val))
	}
}

func waDefPredeclaredNil() {
	waDef(&Nil{object{name: token.K_nil, typ: Typ[UntypedNil], color_: black}})
}

var waPredeclaredFuncs = [...]struct {
	name     string
	nargs    int
	variadic bool
	kind     exprKind
}{
	_Append:  {token.K_append, 1, true, expression},
	_Cap:     {token.K_cap, 1, false, expression},
	_Complex: {token.K_complex, 2, false, expression},
	_Copy:    {token.K_copy, 2, false, statement},
	_Delete:  {token.K_delete, 2, false, statement},
	_Imag:    {token.K_imag, 1, false, expression},
	_Len:     {token.K_len, 1, false, expression},
	_Make:    {token.K_make, 1, true, expression},
	_New:     {token.K_new, 1, true, expression},
	_Panic:   {token.K_panic, 1, false, statement},
	_Print:   {token.K_print, 0, true, statement},
	_Println: {token.K_println, 0, true, statement},
	_Real:    {token.K_real, 1, false, expression},

	// test
	_Assert: {token.K_assert, 1, true, statement},
	_Trace:  {token.K_trace, 0, true, statement},

	_unsafe_Raw:        {token.K_unsafe_Raw, 1, false, expression},
	_unsafe_Alignof:    {token.K_unsafe_Alignof, 1, false, expression},
	_unsafe_Offsetof:   {token.K_unsafe_Offsetof, 1, false, expression},
	_unsafe_Sizeof:     {token.K_unsafe_Sizeof, 1, false, expression},
	_unsafe_SliceData:  {token.K_unsafe_SliceData, 1, false, expression},
	_unsafe_StringData: {token.K_unsafe_StringData, 1, false, expression},
	_unsafe_MakeSlice:  {token.K_unsafe_MakeSlice, 2, false, expression},
	_unsafe_MakeString: {token.K_unsafe_MakeString, 2, false, expression},

	_runtime_SetFinalizer: {token.K_runtime_SetFinalizer, 2, false, statement},
}

func waNewBuiltin(id builtinId) *Builtin {
	return &Builtin{object{name: waPredeclaredFuncs[id].name, typ: Typ[Invalid], color_: black}, id}
}

func waDefPredeclaredFuncs() {
	for i := range waPredeclaredFuncs {
		id := builtinId(i)
		if id == _runtime_SetFinalizer {
			continue // 在加载 runtime 包时导入
		}
		if id == _Assert || id == _Trace {
			continue // only define these in testing environment
		}
		waDef(waNewBuiltin(id))
	}
}

// 注册运行时函数
func WaDefPredeclaredRuntimeFuncs(runtimePkg *Package) {
	if runtimePkg.scope.Lookup(token.K_runtime_SetFinalizer) != nil {
		return // already defined
	}
	waDefInPackage(runtimePkg, waNewBuiltin(_runtime_SetFinalizer))
}

// DefPredeclaredTestFuncs defines the assert and trace built-ins.
// These built-ins are intended for debugging and testing of this
// package only.
func WaDefPredeclaredTestFuncs() {
	if WaUniverse.Lookup(token.K_assert) != nil {
		return // already defined
	}
	waDef(waNewBuiltin(_Assert))
	waDef(waNewBuiltin(_Trace))
}

func waDef(obj Object) {
	assert(obj.color() == black)
	name := obj.Name()
	if strings.Contains(name, " ") {
		return // nothing to do
	}
	// fix Obj link for named types
	if typ, ok := obj.Type().(*Named); ok {
		typ.obj = obj.(*TypeName)
	}

	// 区分 builtin 还是 unsafe
	// runtime 只有 1 个函数已经被过滤了
	// 因此如果是导出的一定是 unsafe 包的
	scope := WaUniverse
	if obj.Exported() && obj.Name() != token.K_错误 {
		scope = WaUnsafe.scope
		// set Pkg field
		switch obj := obj.(type) {
		case *TypeName:
			obj.pkg = WaUnsafe
		case *Builtin:
			obj.pkg = WaUnsafe
		default:
			unreachable()
		}
	}
	if scope.Insert(obj) != nil {
		panic("internal error: double declaration:" + obj.Name())
	}
}

func waDefInPackage(pkg *Package, obj Object) {
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

func initWa() {
	WaUniverse = NewScope(nil, token.NoPos, token.NoPos, token.K_pkg_universe)
	WaUnsafe = NewPackage(token.K_pkg_unsafe, token.K_pkg_unsafe, false)
	WaUnsafe.complete = true

	waDefPredeclaredTypes()
	waDefPredeclaredConsts()
	waDefPredeclaredNil()
	waDefPredeclaredFuncs()

	waUniverseIota = WaUniverse.Lookup(token.K_iota).(*Const)
	waUniverseByte = WaUniverse.Lookup(token.K_byte).(*TypeName).typ.(*Basic)
	waUniverseRune = WaUniverse.Lookup(token.K_rune).(*TypeName).typ.(*Basic)
	waUniverseString = WaUniverse.Lookup(token.K_string).(*TypeName).typ.(*Basic)
	waUniverseAny = WaUniverse.Lookup(token.K_any)
	waUniverseError = WaUniverse.Lookup(token.K_error)

	waUniverse__PACKAGE__ = WaUniverse.Lookup(token.K__PACKAGE__).(*Const)
	waUniverse__FILE__ = WaUniverse.Lookup(token.K__FILE__).(*Const)
	waUniverse__LINE__ = WaUniverse.Lookup(token.K__LINE__).(*Const)
	waUniverse__COLUMN__ = WaUniverse.Lookup(token.K__COLUMN__).(*Const)
	waUniverse__FUNC__ = WaUniverse.Lookup(token.K__FUNC__).(*Const)
	waUniverse__POS__ = WaUniverse.Lookup(token.K__POS__).(*Const)
}
