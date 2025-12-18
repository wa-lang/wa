// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import "strings"

/**************************************
本文件包含了 Block 对象的功能
**************************************/

//-------------------------------------

/**************************************
Block: 指令块，对应于 {...}。不可直接声明该类型对象，必须通过 Module.NewBlock() 创建
Block 本身也满足指令接口，意味着指令块可嵌套
Block 定义了作用域，块内的值无法在块外访问
函数体对应的 Block，其父 Block 应为 nil
Todo: Block 是否满足 Value（既是否可有返回值）待讨论
**************************************/
type Block struct {
	aImv
	typ     Type
	Comment string // 附加注释
	//Locals  []Value       // 该块内定义的局部变量
	Instrs []Instruction // 该块所含的指令

	objects map[interface{}]Location // 关联 AST 结点 -> 块内值
	types   *Types                   // 该函数所属 Module 的类型库，切勿手动修改
}

// 初始化 Block
func (b *Block) init(typ Type) {
	b.typ = typ
	b.objects = make(map[interface{}]Location)
}

func (b *Block) Type() Type {
	return b.typ
}

// Scope 接口相关
func (b *Block) ScopeKind() ScopeKind { return ScopeKindBlock }
func (b *Block) Lookup(obj interface{}, escaping bool) Location {
	if v, ok := b.objects[obj]; ok {
		if alloc, ok := v.(*InstAlloc); ok {
			if escaping {
				alloc.Location = LocationKindHeap
			}
		}
		return v
	}

	if b.scope.ScopeKind() == ScopeKindBlock {
		// b 的父域仍是 Block
		return b.scope.Lookup(obj, escaping)
	}

	parent_fn := b.scope
	if parent_fn.ScopeKind() != ScopeKindFunc {
		panic("Parent of fn-body should be function")
	}

	if parent_fn.ParentScope().ScopeKind() == ScopeKindModule {
		// b 所属的函数是全局函数
		return parent_fn.Lookup(obj, escaping)
	}

	// b 所属的函数是闭包，需要进行变量捕捉
	outer := parent_fn.Lookup(obj, true)
	v := &FreeVar{
		name:   outer.Name(),
		typ:    outer.Type(),
		pos:    outer.Pos(),
		object: outer.Object(),
		outer:  outer,
	}
	return v
}
func (b *Block) Format(tab string, sb *strings.Builder) {
	sb.WriteString(tab)
	sb.WriteString("{\n")

	tab_t := tab + "  "
	for _, v := range b.Instrs {
		v.Format(tab_t, sb)
		sb.WriteString("\n")
	}
	sb.WriteString(tab)
	sb.WriteString("}")
}

// AddLocal 在 Block 中分配一个局部变量，默认分配在栈上
func (b *Block) AddLocal(name string, typ Type, pos int, obj interface{}) Location {
	v := &InstAlloc{}
	v.Stringer = v
	v.refType = b.types.GenPtr(typ)
	v.dataType = typ
	v.pos = pos
	v.object = obj
	if obj != nil {
		b.objects[obj] = v
	}

	b.emit(v)
	return v
}

// CreateBlock 创建一个 Block 初始化其 scope 等，但并不添加至父 Block 中
func (b *Block) createBlock(comment string, typ Type, pos int) *Block {
	block := &Block{}
	block.Stringer = block
	block.Comment = comment
	block.typ = typ
	block.pos = pos
	block.objects = make(map[interface{}]Location)
	block.types = b.types

	block.setScope(b)
	return block
}

// EmitBlock 在 Bloc 中添加一个子 Block
func (b *Block) EmitBlock(comment string, typ Type, pos int) *Block {
	block := b.createBlock(comment, typ, pos)

	b.emit(block)
	return block
}

// emit 向 Block 中添加一个指令，并返回指令对应的 imv
func (b *Block) emit(inst Instruction) Value {
	inst.setScope(b)
	b.Instrs = append(b.Instrs, inst)
	v, _ := inst.(Value)
	return v
}
