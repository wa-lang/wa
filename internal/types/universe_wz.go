// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package types

import (
	"strings"

	"wa-lang.org/wa/internal/constant"
	"wa-lang.org/wa/internal/token"
)

var (
	WzUniverse *Scope
	WzUnsafe   *Package
)

var (
	wzUniverseIota   *Const
	wzUniverseByte   *Basic // uint8 alias, but has name "byte"
	wzUniverseRune   *Basic // int32 alias, but has name "rune"
	wzUniverseString *Basic
	wzUniverseAny    Object

	wzUniverse__PACKAGE__ *Const
	wzUniverse__FILE__    *Const
	wzUniverse__LINE__    *Const
	wzUniverse__COLUMN__  *Const
	wzUniverse__FUNC__    *Const
	wzUniverse__POS__     *Const
)

var wzAliases = [...]*Basic{
	{Bool, IsBoolean, token.K_布尔},
	{Byte, IsInteger | IsUnsigned, token.K_字节},
	{Rune, IsInteger, token.K_符文},

	{Int8, IsInteger, token.K_微整型},
	{Int16, IsInteger, token.K_短整型},
	{Int32, IsInteger, token.K_普整型},
	{Int64, IsInteger, token.K_长整型},
	{Int, IsInteger, token.K_整型},

	{Uint8, IsInteger | IsUnsigned, token.K_微正整},
	{Uint16, IsInteger | IsUnsigned, token.K_短正整},
	{Uint32, IsInteger | IsUnsigned, token.K_普正整},
	{Uint64, IsInteger | IsUnsigned, token.K_长正整},
	{Uint, IsInteger | IsUnsigned, token.K_正整},

	{Float32, IsFloat, token.K_单精},
	{Float64, IsFloat, token.K_双精},

	{String, IsString, token.K_字串},

	{Uintptr, IsInteger | IsUnsigned, token.K_地址型},
	{UnsafePointer, 0, token.K_unsafe_指针},
}

func wzDefPredeclaredTypes() {
	for _, t := range Typ {
		wzDef(NewTypeName(token.NoPos, nil, t.name, t))
	}
	for _, t := range wzAliases {
		wzDef(NewTypeName(token.NoPos, nil, t.name, t))
	}

	// type any = interface{}
	// Note: don't use &emptyInterface for the type of any. Using a unique
	// pointer allows us to detect any and format it as "any" rather than
	// interface{}, which clarifies user-facing error messages significantly.
	wzDef(NewTypeName(token.NoPos, nil, token.K_皮囊, &Interface{}))

	// Error has a nil package in its qualified name since it is in no package
	res := NewVar(token.NoPos, nil, "", Typ[String])
	sig := &Signature{results: NewTuple(res)}
	err := NewFunc(token.NoPos, nil, token.K_报错信息, sig)
	typ := &Named{underlying: NewInterfaceType([]*Func{err}, nil).Complete()}
	sig.recv = NewVar(token.NoPos, nil, "", typ)
	wzDef(NewTypeName(token.NoPos, nil, token.K_错误, typ))
}

var wzPredeclaredConsts = [...]struct {
	name string
	kind BasicKind
	val  constant.Value
}{
	{token.K_真, UntypedBool, constant.MakeBool(true)},
	{token.K_假, UntypedBool, constant.MakeBool(false)},
	{token.K_嘀嗒, UntypedInt, constant.MakeInt64(0)},

	{token.K__包__, UntypedString, constant.MakeString("<wa-lang:__PACKAGE__>")},
	{token.K__文件__, UntypedString, constant.MakeString("<wa-lang:__FILE__>")},
	{token.K__行号__, UntypedInt, constant.MakeInt64(0)},
	{token.K__列号__, UntypedInt, constant.MakeInt64(0)},
	{token.K__函数__, UntypedString, constant.MakeString("")},
	{token.K__位置__, UntypedInt, constant.MakeInt64(0)},
}

func wzDefPredeclaredConsts() {
	for _, c := range wzPredeclaredConsts {
		wzDef(NewConst(token.NoPos, nil, c.name, Typ[c.kind], c.val))
	}
}

func wzDefPredeclaredNil() {
	wzDef(&Nil{object{name: token.K_空, typ: Typ[UntypedNil], color_: black}})
}

var wzPredeclaredFuncs = [...]struct {
	name     string
	nargs    int
	variadic bool
	kind     exprKind
}{
	_Append:  {token.K_追加, 1, true, expression},
	_Cap:     {token.K_容量, 1, false, expression},
	_Complex: {token.K_复数, 2, false, expression},
	_Copy:    {token.K_拷贝, 2, false, statement},
	_Delete:  {token.K_删除, 2, false, statement},
	_Imag:    {token.K_虚部, 1, false, expression},
	_Len:     {token.K_长度, 1, false, expression},
	_Make:    {token.K_构建, 1, true, expression},
	_New:     {token.K_新建, 1, true, expression},
	_Panic:   {token.K_崩溃, 1, false, statement},
	_Print:   {token.K_输出, 0, true, statement},
	_Println: {token.K_打印, 0, true, statement},
	_Real:    {token.K_实部, 1, false, expression},

	_unsafe_Raw:      {token.K_unsafe_原生, 1, false, expression},
	_unsafe_Alignof:  {token.K_unsafe_对齐倍数, 1, false, expression},
	_unsafe_Offsetof: {token.K_unsafe_字节偏移量, 1, false, expression},
	_unsafe_Sizeof:   {token.K_unsafe_字节大小, 1, false, expression},

	_runtime_SetFinalizer: {token.K_runtime_设置终结函数, 2, false, statement},

	// 测试环境
	_Assert: {token.K_断言, 1, true, statement},
	_Trace:  {token.K_跟踪, 0, true, statement},
}

func wzNewBuiltin(id builtinId) *Builtin {
	return &Builtin{object{name: wzPredeclaredFuncs[id].name, typ: Typ[Invalid], color_: black}, id}
}

func wzDefPredeclaredFuncs() {
	for i := range wzPredeclaredFuncs {
		id := builtinId(i)
		if id == _runtime_SetFinalizer {
			continue // 在加载 runtime 包时导入
		}
		if id == _Assert || id == _Trace {
			continue // only define these in testing environment
		}
		wzDef(wzNewBuiltin(id))
	}
}

// 注册运行时函数
func WzDefPredeclaredRuntimeFuncs(runtimePkg *Package) {
	if runtimePkg.scope.Lookup(token.K_runtime_设置终结函数) != nil {
		return // already defined
	}
	wzDefInPackage(runtimePkg, wzNewBuiltin(_runtime_SetFinalizer))
}

// DefPredeclaredTestFuncs defines the assert and trace built-ins.
// These built-ins are intended for debugging and testing of this
// package only.
func WzDefPredeclaredTestFuncs() {
	if WzUniverse.Lookup(token.K_断言) != nil {
		return // already defined
	}
	wzDef(wzNewBuiltin(_Assert))
	wzDef(wzNewBuiltin(_Trace))
}

func wzDef(obj Object) {
	assert(obj.color() == black)
	name := obj.Name()
	if strings.Contains(name, " ") {
		return // nothing to do
	}
	// fix Obj link for named types
	if typ, ok := obj.Type().(*Named); ok {
		typ.obj = obj.(*TypeName)
	}

	// 中文无法区分 builtin 还是 unsafe (runtime 只有 1 个函数已经被过滤了)
	scope := WzUniverse
	switch obj.Name() {
	case token.K_unsafe_指针,
		token.K_unsafe_原生,
		token.K_unsafe_对齐倍数,
		token.K_unsafe_字节偏移量,
		token.K_unsafe_字节大小:

		scope = WzUnsafe.scope
		// set Pkg field
		switch obj := obj.(type) {
		case *TypeName:
			obj.pkg = WzUnsafe
		case *Builtin:
			obj.pkg = WzUnsafe
		default:
			unreachable()
		}
	}
	if scope.Insert(obj) != nil {
		panic("internal error: double declaration:" + obj.Name())
	}
}

func wzDefInPackage(pkg *Package, obj Object) {
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

func init() {
	WzUniverse = NewScope(nil, token.NoPos, token.NoPos, token.K_pkg_太初)
	WzUnsafe = NewPackage(token.K_pkg_洪荒, token.K_pkg_洪荒, true)
	WzUnsafe.complete = true

	wzDefPredeclaredTypes()
	wzDefPredeclaredConsts()
	wzDefPredeclaredNil()
	wzDefPredeclaredFuncs()

	wzUniverseIota = WzUniverse.Lookup(token.K_嘀嗒).(*Const)
	wzUniverseByte = WzUniverse.Lookup(token.K_字节).(*TypeName).typ.(*Basic)
	wzUniverseRune = WzUniverse.Lookup(token.K_符文).(*TypeName).typ.(*Basic)
	wzUniverseString = WzUniverse.Lookup(token.K_字串).(*TypeName).typ.(*Basic)
	wzUniverseAny = WzUniverse.Lookup(token.K_皮囊)

	wzUniverse__PACKAGE__ = WzUniverse.Lookup(token.K__包__).(*Const)
	wzUniverse__FILE__ = WzUniverse.Lookup(token.K__文件__).(*Const)
	wzUniverse__LINE__ = WzUniverse.Lookup(token.K__行号__).(*Const)
	wzUniverse__COLUMN__ = WzUniverse.Lookup(token.K__列号__).(*Const)
	wzUniverse__FUNC__ = WzUniverse.Lookup(token.K__函数__).(*Const)
	wzUniverse__POS__ = WzUniverse.Lookup(token.K__位置__).(*Const)
}
