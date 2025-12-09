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
	DataType() Type

	// 与该位置关联的 AST 结点。对凹语言前端，应为 types.Object
	Object() interface{}
}
