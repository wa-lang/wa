// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
)

/**************************************
VarKind: 变量位置类别
**************************************/

type VarKind int

const (
	Register VarKind = iota
	//Stack
	Heap
)

/**************************************
Alloc: Alloc 指令，定义一个变量，实现了 Expr、Var 接口
**************************************/

type Alloc struct {
	aStmt
	kind   VarKind
	name   string
	dtype  Type
	rtype  Type
	object interface{} // 与该值关联的 AST 结点。对凹语言前端，应为 types.Object
	init   Expr        // 初始值
	noinit bool        // 是否不进行 0 值初始化

	tank *tank
}

func (i *Alloc) Name() string { return i.name }
func (i *Alloc) Type() Type {
	if i.kind == Register {
		return i.dtype
	} else {
		return i.rtype
	}
}
func (i *Alloc) retained() bool { return false }
func (i *Alloc) String() string {
	s := ""
	switch i.kind {
	case Register:
		s = fmt.Sprintf("var %s %s", i.name, i.dtype.Name())

	case Heap:
		s = fmt.Sprintf("var %s %s = alloc.heap(%s)", i.name, i.rtype.Name(), i.dtype.Name())

	default:
		panic(fmt.Sprintf("Todo: VarKind: %v", i.kind))
	}

	if i.init != nil {
		s += " = " + i.init.Name()
	} else if !i.noinit {
		s += " = 0"
	}

	if i.tank != nil {
		s += " --- "
		s += i.tank.String()
	}
	return s
}
func (i *Alloc) Kind() VarKind  { return i.kind }
func (i *Alloc) DataType() Type { return i.dtype }
func (i *Alloc) Tank() *tank    { return i.tank }

// func (i *Alloc) RefType() Type  { return i.rtype }
func (i *Alloc) SetInit(init Expr) { i.init = init }

//func (i *Alloc) Object() interface{} { return i.object }

func (b *Block) NewAlloc(name string, typ Type, pos int, obj interface{}, init Expr) *Alloc {
	v := &Alloc{
		kind:   Register,
		name:   name,
		dtype:  typ,
		rtype:  b.types.GenRef(typ),
		init:   init,
		object: obj,
	}
	v.Stringer = v
	v.pos = pos

	return v
}

// AddLocal 在 Block 中分配一个局部变量，初始时位于 Register
func (b *Block) AddLocal(name string, typ Type, pos int, obj interface{}, init Expr) *Alloc {
	v := b.NewAlloc(name, typ, pos, obj, init)
	if obj != nil {
		b.objects[obj] = v
	}

	b.emit(v)
	return v
}

/**************************************
Imv: Imv 指令，定义一个中间变量，实现了 Expr、Var 接口
**************************************/

type Imv struct {
	aStmt
	name string
	val  Expr // 初始值

	tank *tank
}

func (i *Imv) Name() string   { return i.name }
func (i *Imv) Type() Type     { return i.val.Type() }
func (i *Imv) retained() bool { return false }
func (i *Imv) String() string {
	s := fmt.Sprintf("imv %s = %s", i.name, i.val.Name())
	if i.tank != nil {
		s += " --- "
		s += i.tank.String()
	}
	return s
}
func (i *Imv) Kind() VarKind  { return Register }
func (i *Imv) DataType() Type { return i.Type() }
func (i *Imv) Tank() *tank    { return i.tank }

func NewImv(name string, val Expr, pos int) *Imv {
	v := &Imv{
		name: name,
		val:  val,
	}
	v.Stringer = v
	v.pos = pos

	return v
}

/**************************************
Drop: Drop 指令，丢弃 Var，丢弃后它所占用的虚拟寄存器可被重用
  - Drop 一个 chunk 时并不会执行 release
**************************************/

type Drop struct {
	aStmt
	X Var
}

func (i *Drop) String() string {
	tank := i.X.Tank()
	if tank != nil {
		return fmt.Sprintf("drop(%s --- %s)", i.X.Name(), tank.String())
	} else {
		return fmt.Sprintf("drop(%s)", i.X.Name())
	}
}

// 生成一条 Drop 指令
func NewDrop(x Var, pos int) *Drop {
	v := &Drop{X: x}
	v.Stringer = v
	v.pos = pos
	return v
}
