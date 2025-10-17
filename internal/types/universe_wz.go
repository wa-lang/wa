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
	{Byte, IsInteger | IsUnsigned, token.K_字节},
	{Rune, IsInteger, "rune"},
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

	{Float32, IsFloat, token.K_单精},
	{Float64, IsFloat, token.K_双精},

	{String, IsString, token.K_字串},
}

func wzDefPredeclaredTypes() {
	for _, t := range Typ {
		wzDef(NewTypeName(token.NoPos, nil, t.name, t))
	}
	for _, t := range aliases {
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
	err := NewFunc(token.NoPos, nil, token.K_报错, sig)
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
	{token.K_约塔, UntypedInt, constant.MakeInt64(0)},

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

	_unsafe_Alignof:  {"Alignof", 1, false, expression},
	_unsafe_Offsetof: {"Offsetof", 1, false, expression},
	_unsafe_Sizeof:   {"Sizeof", 1, false, expression},

	_断言: {token.K_断言, 1, true, statement},
}

func wzDefPredeclaredFuncs() {
	for i := range wzPredeclaredFuncs {
		id := builtinId(i)
		if id == _runtime_SetFinalizer {
			continue // 在加载 runtime 包时导入
		}
		if id == _Assert || id == _Trace || id == _断言 {
			continue // only define these in testing environment
		}
		wzDef(newBuiltin(id))
	}
}

// 注册运行时函数
func WzDefPredeclaredRuntimeFuncs(runtimePkg *Package) {
	if runtimePkg.scope.Lookup("SetFinalize") != nil {
		return // already defined
	}
	wzDefInPackage(runtimePkg, newBuiltin(_runtime_SetFinalizer))
}

// DefPredeclaredTestFuncs defines the assert and trace built-ins.
// These built-ins are intended for debugging and testing of this
// package only.
func WzDefPredeclaredTestFuncs() {
	if WzUniverse.Lookup(token.K_断言) != nil {
		return // already defined
	}
	def(newBuiltin(_断言))
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

	// TODO: 中文的名字都是导出的, 因此需要区分 builtin 和 unsafe

	scope := WzUniverse
	if !isCNKeyword(obj.Name()) && obj.Exported() {
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
	WzUniverse = NewScope(nil, token.NoPos, token.NoPos, token.K_太初)
	WzUniverse.isUniverse = true

	WzUnsafe = NewPackage(token.K_鸿蒙, token.K_鸿蒙)
	WzUnsafe.complete = true

	wzDefPredeclaredTypes()
	wzDefPredeclaredConsts()
	wzDefPredeclaredNil()
	wzDefPredeclaredFuncs()

	wzUniverseIota = WzUniverse.Lookup(token.K_约塔).(*Const)
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
