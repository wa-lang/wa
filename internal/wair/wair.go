// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wair

/**************************************
本包定义了凹语言的高级中间表达 wair
**************************************/

//-------------------------------------

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
