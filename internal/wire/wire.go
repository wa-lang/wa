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
	Lookup(obj interface{}, level LocationKind) Location

	// 格式化输出
	Format(tab string, sb *strings.Builder)
}

//-------------------------------------

///**************************************
//ValueKind: 值的类别，取值见后续常量
//**************************************/
//type ValueKind int
//
//const (
//	ValueKindLocal ValueKind = iota
//	ValueKindGlobal
//	ValueKindConst
//)

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
	TypeKindChunk
	TypeKindRef
	TypeKindStruct
	TypeKindTuple
	//TypeKindSignature
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

///**************************************
//NamedType: 具名值类型
//**************************************/
//type NamedType interface {
//	Name() string    //类型名，自定义类型应包含包路径，需要进行名字修饰
//	Kind() TypeKind  //该类型的类别
//	Equal(Type) bool //判断该类型与输入类型是否相等，注意该比较仅关心类别和结构，不关心类型名
//
//	AddMethod(m Method) int //为该类型添加方法，返回方法id
//	NumMethods() int        //该类型的方法数
//	Method(i int) Method    //获取指定id的方法
//}

/**************************************
Stmt: 指令接口
**************************************/
type Stmt interface {
	// 获取该指令的伪代码
	String() string

	// 获取该指令在源码中的位置
	Pos() int

	// 格式化输出
	Format(tab string, sb *strings.Builder)

	// 获取该指令所属的域
	ParentScope() Scope

	// 设置该指令所属的域
	setScope(Scope)
}

/**************************************
Expr: 表达式接口，所有可以作为指令参数的对象，都满足该接口
**************************************/
type Expr interface {
	// 表达式的名字
	// 变量的名字是其变量名，常量的名字是其字面量，除此外多数表达式的名字是其指令伪代码
	Name() string

	// 该表达式的类型
	Type() Type

	// 表达式在源码中的位置
	Pos() int

	// 保留
	retained() bool
}

/**************************************
Param: 函数的输入参数，满足 Expr 接口
**************************************/
type Param struct {
	name string
	typ  Type
	pos  int
}

func (p *Param) Name() string   { return p.name }
func (p *Param) Format() string { return p.name }
func (p *Param) Type() Type     { return p.typ }
func (p *Param) Pos() int       { return p.pos }
func (p *Param) retained() bool { return false }

/**************************************
Const: 常量，满足 Expr 接口
**************************************/
type Const struct {
	name string
	typ  Type
	pos  int
}

func (p *Const) Name() string   { return p.name }
func (p *Const) Type() Type     { return p.typ }
func (p *Const) Pos() int       { return p.pos }
func (p *Const) retained() bool { return false }

/**************************************
FnSig: 函数签名
**************************************/
type FnSig struct {
	Params  []Type //参数类型列表
	Results Type   //返回值类型，无返回值: Void，多返回值：Tuple
}

//func (p *FnSig) Name() string    { panic("FnSig.Name() is unimplemented") }
//func (p *FnSig) Kind() TypeKind  { return TypeKindSignature }
func (s *FnSig) Equal(d *FnSig) bool {
	if len(s.Params) != len(d.Params) {
		return false
	}

	for i := range s.Params {
		if !s.Params[i].Equal(d.Params[i]) {
			return false
		}
	}

	return s.Results.Equal(d.Results)
}

/**************************************
Method: 方法
**************************************/
type Method struct {
	Sig      FnSig  //函数签名
	Name     string //方法名。b.MyMethod() 的方法名为 "MyMethod"
	FullName string //方法的函数全路径名（包括包路径、类型名，需要进行名字修饰），接口方法的该属性为空
}

func (m *Method) Equal(d *Method) bool {
	if m.Name != d.Name {
		return false
	}
	return m.Sig.Equal(&d.Sig)
}

/**************************************
FreeVar: 闭包捕获的外部变量
**************************************/
type FreeVar struct {
	name   string
	typ    Type
	pos    int
	object interface{}
	outer  Location // 被捕获的闭包变量
}

func (p *FreeVar) Name() string               { return p.name }
func (p *FreeVar) Type() Type                 { return p.typ }
func (p *FreeVar) Pos() int                   { return p.pos }
func (p *FreeVar) retained() bool             { return false }
func (p *FreeVar) Object() interface{}        { return p.object }
func (p *FreeVar) LocationKind() LocationKind { return LocationKindHeap }
func (p *FreeVar) DataType() Type             { return p.typ }
