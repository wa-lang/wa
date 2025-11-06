// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

/**************************************
LocationKind: 位置类别
**************************************/
type LocationKind int

const (
	LocationKindStack LocationKind = iota
	LocationKindRegister
	LocationKindHeap
)

/**************************************
Location: 可被赋值的位置，可能类别为栈地址、寄存器、堆地址
**************************************/
type Location interface {
	Value

	// 变量位置类型
	LocationKind() LocationKind

	// Location指向的对象的类型
	DataType() ValueType

	// 与该位置关联的 AST 结点。对凹语言前端，应为 types.Object
	Object() interface{}
}

/**************************************
Blank: 匿名变量 "_"
**************************************/
type Blank struct{}

func (*Blank) Name() string               { return "_" }
func (*Blank) Kind() ValueKind            { return ValueKindLocal }
func (*Blank) Type() ValueType            { panic("Blank.Type() is unimplemented") }
func (*Blank) Pos() int                   { panic("Blank.Pos() is unimplemented") }
func (*Blank) Object() interface{}        { panic("Blank.Object() is unimplemented") }
func (*Blank) LocationKind() LocationKind { return LocationKindStack }
func (*Blank) DataType() ValueType        { panic("Balck.DataType() is unimplemented") }
