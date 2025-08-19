// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wair

/**************************************
本包定义了凹语言的高级中间表达 wair
**************************************/

//-------------------------------------

/**************************************
Module: 定义了一个 wair 模块，对应于 ast.Program
**************************************/
type Module struct {
	Types   Types
	Globals []Value
	Funcs   []*Function
}

// 初始化 Module
func (m *Module) Init() {
	m.Types.Init()
}

// 创建一个 Value。注意该函数仅创建值，并不会将其合并至 Module 的相应位置（如 Globals）
func (m *Module) NewValue(name string, kind ValueKind, typ ValueType, pos int, obj interface{}) Value {
	return m.newValue(name, kind, typ, pos, obj)
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
	TypeKindF32
	TypeKindF64
	TypeKindComplex64
	TypeKindComplex128
	TypeKindRune
	TypeKindPtr
	TypeKindStruct
	TypeKindTuple
	TypeKindString
	TypeKindSlice
	TypeKindArray
	TypeKindMap
	TypeKindInterface
)

/**************************************
ValueType: 值类型
**************************************/
type ValueType interface {
	Name() string         //类型名，自定义类型应包含包路径，需要进行名字修饰
	Size() int            //该类型所占字节数
	Align() int           //该类型对齐字节数
	Kind() TypeKind       //该类型的类别
	Equal(ValueType) bool //判断该类型与输入类型是否相等，注意该比较仅关心类别和结构，不关心类型名

	AddMethod(m Method) int //为该类型添加方法，返回方法id
	NumMethods() int        //该类型的方法数
	Method(i int) Method    //获取指定id的方法
}

/**************************************
FnSig: 函数签名
**************************************/
type FnSig struct {
	Params  []ValueType //参数类型列表
	Results []ValueType //返回值类型列表
}

/**************************************
Method: 方法
**************************************/
type Method struct {
	Sig        FnSig  //函数签名
	Name       string //方法名。b.MyMethod() 的方法名为 "MyMethod"
	FullFnName string //方法的函数全路径名（包括包路径、类型名，需要进行名字修饰）
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
	Type() ValueType

	// 该值在源码中的位置
	Pos() int

	// 与该值关联的 AST 结点。对凹语言前端，应为 types.Object
	Object() interface{}
}

/**************************************
Instruction: 指令
**************************************/
type Instruction interface {
	// 获取该指令的伪代码
	String() string

	// 获取该指令在源码中的位置
	Pos() int

	// 获取该指令所属的指令块
	Parent() *Block

	// 设置该指令所属的指令块
	setParent(*Block)
}

/**************************************
Function: 函数
**************************************/
type Function struct {
	InternalName string      // 函数的内部名称(含包路径)，是其身份标识，应进行名字修饰
	ExternalName string      // 函数的导出名称，非导出函数应为 nil
	Params       []Value     // 参数列表
	Results      []ValueType //返回值列表
	Body         *Block      // 函数体，为 nil 表明该函数为外部导入
}

/**************************************
Block: 指令块，对应于 {...}，指令块本身也满足指令接口，意味着指令块可嵌套
指令块定义了作用域，块内的值无法在块外访问
函数体对应的指令块，其 Parent 应为 nil
Todo: Block 是否满足 Value（既是否可有返回值）待讨论
**************************************/
type Block struct {
	Comment string        // 附加注释
	Locals  []Value       // 该块内定义的局部变量
	Instrs  []Instruction // 该块所含的指令
}
