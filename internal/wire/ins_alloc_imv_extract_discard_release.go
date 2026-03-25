// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
)

/**************************************
VarKind: 变量类别
**************************************/
type AllocKind int

const (
	AllocKindRegister AllocKind = iota
	//Stack
	AllocKindHeap
)

/**************************************
Alloc: Alloc 指令，定义一个变量，实现了 Expr、Var 接口
**************************************/
type Alloc struct {
	aStmt
	kind   AllocKind
	name   string
	dtype  Type
	rtype  Type
	object interface{} // 与该值关联的 AST 结点。对凹语言前端，应为 types.Object

	init Expr // 初始值

	tank *tank
}

func (i *Alloc) Name() string { return i.name }
func (i *Alloc) Type() Type {
	if i.kind == AllocKindRegister {
		return i.dtype
	} else {
		return i.rtype
	}
}
func (i *Alloc) retained() bool { return false }
func (i *Alloc) String() string {
	s := ""
	switch i.kind {
	case AllocKindRegister:
		if i.init != nil {
			if rtimp.hasChunk(i.init.Type()) && !i.init.retained() {
				s = fmt.Sprintf("var %s↑ %s = %s", i.name, i.Type().Name(), i.init.Name())
			} else {
				s = fmt.Sprintf("var %s %s = %s", i.name, i.Type().Name(), i.init.Name())
			}
		} else {
			s = fmt.Sprintf("var %s %s = 0", i.name, i.Type().Name())
		}

	case AllocKindHeap:
		//if i.init != nil {
		//panic(fmt.Sprintf("Heap var: %s has init-val", i.name))
		//}
		if i.init != nil {
			s = fmt.Sprintf("var %s %s = alloc.heap(%s:%s)", i.name, i.rtype.Name(), i.init.Name(), i.dtype.Name())
		} else {
			s = fmt.Sprintf("var %s %s = alloc.heap(%s)", i.name, i.rtype.Name(), i.dtype.Name())
		}

	default:
		panic(fmt.Sprintf("Todo: AllocKind: %v", i.kind))
	}

	if i.tank != nil {
		s += " --- "
		s += i.tank.String()
	}
	return s
}
func (i *Alloc) Kind() AllocKind     { return i.kind }
func (i *Alloc) SetKind(t AllocKind) { i.kind = t }
func (i *Alloc) DataType() Type      { return i.dtype }
func (i *Alloc) Tank() *tank         { return i.tank }

// func (i *Alloc) RefType() Type  { return i.rtype }
func (i *Alloc) SetInit(init Expr) { i.init = init }

//func (i *Alloc) Object() interface{} { return i.object }

func (b *Block) NewAlloc(name string, typ Type, pos int, obj interface{}, init Expr) *Alloc {
	v := &Alloc{
		kind:   AllocKindRegister,
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
  - 该指令仅供内部使用，上层高级语法不应直接使用
**************************************/
type Imv struct {
	aStmt
	name string
	val  Expr // 初始值

	tank *tank
}

func (i *Imv) Name() string   { return i.name }
func (i *Imv) Type() Type     { return i.val.Type() }
func (i *Imv) retained() bool { return i.val.retained() }
func (i *Imv) String() string {
	s := fmt.Sprintf("imv %s = %s", i.name, i.val.Name())
	if i.tank != nil {
		s += " --- "
		s += i.tank.String()
	}
	return s
}
func (i *Imv) Kind() AllocKind { return AllocKindRegister }
func (i *Imv) DataType() Type  { return i.Type() }
func (i *Imv) Tank() *tank     { return i.tank }

func newImv(name string, val Expr, pos int) *Imv {
	v := &Imv{
		name: name,
		val:  val,
	}
	v.Stringer = v
	v.pos = pos

	return v
}

/**************************************
Extract: Extract 指令，提取元组变量 X 的第 Index 个元素，Extract 实现了 Var 接口
  - X 应为 Tuple 类型的 Imv
  - 该指令仅供内部使用，上层高级语法不应直接使用
**************************************/
type Extract struct {
	aStmt
	X     *Imv
	Index int
}

func (i *Extract) Name() string   { return i.String() }
func (i *Extract) Type() Type     { return i.X.Type().(*Tuple).members[i.Index] }
func (i *Extract) retained() bool { return i.X.retained() }

func (i *Extract) String() string {
	return fmt.Sprintf("extract(%s, %d)", i.X.Name(), i.Index)
}
func (i *Extract) Kind() AllocKind { return AllocKindRegister }
func (i *Extract) DataType() Type  { return i.Type() }
func (i *Extract) Tank() *tank     { return i.X.Tank().member[i.Index] }

// 生成一条 Extract 指令
func newExtract(x *Imv, index int, pos int) *Extract {
	v := &Extract{}
	v.Stringer = v
	v.X = x
	v.Index = index
	v.pos = pos

	return v
}

/**************************************
Discard: Discard 指令，丢弃 Var，丢弃后它所占用的虚拟寄存器可被重用
  - Discard 一个 chunk 时并不会执行 release
  - 该指令仅供内部使用，上层高级语法不应直接使用
**************************************/
type Discard struct {
	aStmt
	X Var
}

func (i *Discard) String() string {
	tank := i.X.Tank()
	if tank != nil {
		return fmt.Sprintf("discard(%s --- %s)", i.X.Name(), tank.String())
	} else {
		return fmt.Sprintf("discard(%s)", i.X.Name())
	}
}

// 生成一条 Discard 指令
func newDiscard(x Var, pos int) *Discard {
	v := &Discard{X: x}
	v.Stringer = v
	v.pos = pos
	return v
}

/**************************************
Release: 释放一个 chunk
  - 该指令仅供内部使用，上层高级语法不应直接使用
**************************************/
type Release struct {
	aStmt
	X register
}

func (i *Release) String() string {
	return fmt.Sprintf("release($r%d)", i.X.id)
}

// 生成一条 Release 指令
func newRelease(x register, pos int) *Release {
	v := &Release{X: x}
	v.Stringer = v
	v.pos = pos
	return v
}
