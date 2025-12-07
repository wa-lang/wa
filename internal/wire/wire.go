// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import "strings"

/**************************************
本包定义了凹语言的高级中间表达 wire
**************************************/

//-------------------------------------

/**************************************
ScopeKind: 域的类别，取值见后续常量
**************************************/
type ScopeKind int

const (
	ScopeKindModule ScopeKind = iota
	ScopeKindFunc
	ScopeKindBlock
)

type Scope interface {
	// 该域的类型
	ScopeKind() ScopeKind

	// 父域（域可嵌套）
	ParentScope() Scope

	// 递归向上查找关联对象为 obj 的 变量位置，对凹语言前端，obj 应为 types.Object
	Lookup(obj interface{}, escaping bool) Location

	// 格式化输出
	Format(tab string, sb *strings.Builder)
}

//-------------------------------------

/**************************************
ValueKind: 值的类别，取值见后续常量
**************************************/
type ValueKind int

const (
	ValueKindLocal ValueKind = iota
	ValueKindGlobal
	ValueKindConst
)

/**************************************
TypeKind: 值类型的类别，取值见后续常量
**************************************/
type TypeKind int

const (
	TypeKindUnknown TypeKind = iota
	TypeKindVoid
	TypeKindBool
	TypeKindU8
	TypeKindU16
	TypeKindU32
	TypeKindU64
	TypeKindI8
	TypeKindI16
	TypeKindI32
	TypeKindI64
	TypeKindInt
	TypeKindUint
	TypeKindF32
	TypeKindF64
	TypeKindComplex64
	TypeKindComplex128
	TypeKindRune
	TypeKindString

	TypeKindPtr
	TypeKindStruct
	TypeKindTuple
	TypeKindSignature
	TypeKindSlice
	TypeKindArray
	TypeKindMap
	TypeKindInterface
)

/**************************************
Type: 值类型
**************************************/
type Type interface {
	Name() string    //类型名，自定义类型应包含包路径，需要进行名字修饰
	Kind() TypeKind  //该类型的类别
	Equal(Type) bool //判断该类型与输入类型是否相等，注意该比较仅关心类别和结构，不关心类型名
}

/**************************************
NamedType: 具名值类型
**************************************/
type NamedType interface {
	Name() string    //类型名，自定义类型应包含包路径，需要进行名字修饰
	Kind() TypeKind  //该类型的类别
	Equal(Type) bool //判断该类型与输入类型是否相等，注意该比较仅关心类别和结构，不关心类型名

	AddMethod(m Method) int //为该类型添加方法，返回方法id
	NumMethods() int        //该类型的方法数
	Method(i int) Method    //获取指定id的方法
}

/**************************************
Value: 值，所有可以作为指令参数的对象，都满足该接口
**************************************/
type Value interface {
	// 该值的名字
	// 全局变量、局部变量（含参数）、具名函数的名字与其源代码中的对应标识符保持一直
	// 常量的名字是其字面值
	// 中间变量（虚拟寄存器）的名字为 $t0、$t1 等
	Name() string

	// 该值的类别
	Kind() ValueKind

	// 该值的类型
	Type() Type

	// 该值在源码中的位置
	Pos() int
}

/**************************************
Param: 函数的输入参数，满足 Value 接口
**************************************/
type Param struct {
	name string
	typ  Type
	pos  int
}

func (p *Param) Name() string    { return p.name }
func (p *Param) Kind() ValueKind { return ValueKindLocal }
func (p *Param) Type() Type      { return p.typ }
func (p *Param) Pos() int        { return p.pos }

/**************************************
Const: 常量，满足 Value 接口
**************************************/
type Const struct {
	name string
	typ  Type
	pos  int
}

func (p *Const) Name() string    { return p.name }
func (p *Const) Kind() ValueKind { return ValueKindConst }
func (p *Const) Type() Type      { return p.typ }
func (p *Const) Pos() int        { return p.pos }

/**************************************
FnSig: 函数签名
**************************************/
type FnSig struct {
	Params  []Type //参数类型列表
	Results []Type //返回值类型列表
}

func (p *FnSig) Name() string           { panic("FnSig.Name() is unimplemented") }
func (p *FnSig) Kind() TypeKind         { return TypeKindSignature }
func (p *FnSig) Equal(Type) bool        { panic("FnSig.Equal() is unimplemented") }
func (p *FnSig) AddMethod(m Method) int { panic("FnSig.AddMethod() is unimplemented") }
func (p *FnSig) NumMethods() int        { panic("FnSig.NumMethods() is unimplemented") }
func (p *FnSig) Method(i int) Method    { panic("FnSig.Method() is unimplemented") }

/**************************************
Method: 方法
**************************************/
type Method struct {
	Sig        FnSig  //函数签名
	Name       string //方法名。b.MyMethod() 的方法名为 "MyMethod"
	FullFnName string //方法的函数全路径名（包括包路径、类型名，需要进行名字修饰）
}

/**************************************
FreeVar: 闭包捕获的外部变量
**************************************/
type FreeVar struct {
	name   string
	typ    Type
	pos    int
	object interface{}
	outer  Value // 被捕获的闭包变量
}

func (p *FreeVar) Name() string               { return p.name }
func (p *FreeVar) Kind() ValueKind            { return ValueKindLocal }
func (p *FreeVar) Type() Type                 { return p.typ }
func (p *FreeVar) Pos() int                   { return p.pos }
func (p *FreeVar) Object() interface{}        { return p.object }
func (p *FreeVar) LocationKind() LocationKind { return LocationKindHeap }
func (p *FreeVar) DataType() Type             { return p.typ }

/**************************************
Builtin: 内置函数
**************************************/
type Builtin struct {
	name string
	sig  FnSig
}

func (p *Builtin) Name() string        { return p.name }
func (p *Builtin) Kind() ValueKind     { panic("Builtin.Kind() is unimplemented") }
func (p *Builtin) Type() Type          { return &p.sig }
func (p *Builtin) Pos() int            { panic("Builtin.Pos() is unimplemented") }
func (p *Builtin) Object() interface{} { panic("Builtin.Object() is unimplemented") }
