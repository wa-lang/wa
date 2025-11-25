// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

/**************************************
Instruction: 指令接口
**************************************/
type Instruction interface {
	// 获取该指令的伪代码
	String() string

	// 获取该指令在源码中的位置
	Pos() int

	// 获取该指令所属的域
	ParentScope() Scope

	// 格式化输出
	Format(tab string, sb *strings.Builder)

	// 设置该指令所属的域
	setScope(Scope)
}

/**************************************
InstAlloc: Alloc 指令，分配一个变量位置，该对象同时实现了 Location 接口
**************************************/
type InstAlloc struct {
	aImv
	Location LocationKind
	dataType ValueType
	object   interface{} // 与该值关联的 AST 结点。对凹语言前端，应为 types.Object
}

func (i *InstAlloc) String() string {
	op := "alloc.stack"
	switch i.Location {
	case LocationKindRegister:
		op = "alloc.register"

	case LocationKindHeap:
		op = "alloc.heap"
	}

	return fmt.Sprintf("%s %s", op, i.Type().Name())
}

// Location 接口相关
func (i *InstAlloc) LocationKind() LocationKind { return i.Location }
func (i *InstAlloc) DataType() ValueType        { return i.dataType }
func (i *InstAlloc) Object() interface{}        { return i.object }

/**************************************
InstLoad: Load 指令，装载 Loc 处的变量
**************************************/
type InstLoad struct {
	aImv
	Loc Location
}

func (i *InstLoad) String() string {
	return fmt.Sprintf("*%s", i.Loc.Name())
}

/**************************************
InstStore: Store 指令，将 Val 存储到 Loc 指定的位置
**************************************/
type InstStore struct {
	anInstruction
	Loc Location
	Val Value
}

func (i *InstStore) String() string {
	return fmt.Sprintf("*%s = %s", i.Loc.Name(), i.Val.Name())
}

/**************************************
InstExtract: Extract 指令，提取元组 Tuple 的第 Index 个元素
**************************************/
type InstExtract struct {
	aImv
	Tuple Value
	Index int
}

func (i *InstExtract) String() string {
	return fmt.Sprintf("extract %s #%d", i.Tuple.Name(), i.Index)
}

/**************************************
InstReturn: Return 指令
**************************************/
type InstReturn struct {
	anInstruction
	Results []Value
}

func (i *InstReturn) String() string {
	s := "return"
	for _, r := range i.Results {
		s += " "
		s += r.Name()
	}
	return s
}

/**************************************
InstUnopNot:  一元非指令
**************************************/
type InstUnopNot struct {
	aImv
	X Value
}

func (i *InstUnopNot) String() string {
	return fmt.Sprintf("!%s", i.X.Name())
}

/**************************************
InstUnopSub:  取负指令
**************************************/
type InstUnopSub struct {
	aImv
	X Value
}

func (i *InstUnopSub) String() string {
	return fmt.Sprintf("-%s", i.X.Name())
}

/**************************************
InstUnopXor:  一元异或指令
**************************************/
type InstUnopXor struct {
	aImv
	X Value
}

func (i *InstUnopXor) String() string {
	return fmt.Sprintf("^%s", i.X.Name())
}
