// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"strings"
)

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
	aStmt
	Label string // 标签
	Stmts []Stmt // 该块所含的指令

	scope   Scope
	objects map[interface{}]*Var // AST 结点 -> 块内变量
	types   *Types               // 该函数所属 Module 的类型库，切勿手动修改
}

// 初始化 Block
func (b *Block) init() {
	b.objects = make(map[interface{}]*Var)
}

// Scope 接口相关
func (b *Block) ScopeKind() ScopeKind { return ScopeKindBlock }
func (b *Block) Lookup(obj interface{}, level VarKind) *Var {
	if v, ok := b.objects[obj]; ok {
		if level > v.kind {
			v.kind = level
		}
		return v
	}

	if b.scope.ScopeKind() == ScopeKindBlock {
		// b 的父域仍是 Block
		return b.scope.Lookup(obj, level)
	}

	parent_fn := b.scope
	if parent_fn.ScopeKind() != ScopeKindFunc {
		panic("Parent of fn-body should be function")
	}

	if parent_fn.ParentScope().ScopeKind() == ScopeKindModule {
		// b 所属的函数是全局函数
		return parent_fn.Lookup(obj, level)
	}

	// b 所属的函数是闭包，需要进行变量捕捉
	//outer := parent_fn.Lookup(obj, LocationKindHeap)
	//v := &FreeVar{
	//	outer: outer,
	//}
	//return v
	panic("Todo")
}
func (b *Block) Format(tab string, sb *strings.Builder) {
	sb.WriteString(tab)
	if len(b.Label) > 0 {
		sb.WriteString("{ : ")
		sb.WriteString(b.Label)
		sb.WriteRune('\n')
	} else {
		sb.WriteString("{\n")
	}

	tab_t := tab + "  "
	for _, v := range b.Stmts {
		v.Format(tab_t, sb)
		sb.WriteString("\n")
	}
	sb.WriteString(tab)
	sb.WriteString("}")
}

func (b *Block) ParentScope() Scope {
	return b.scope
}

// CreateBlock 创建一个 Block 初始化其 scope 等，但并不添加至父 Block 中
func (b *Block) createBlock(label string, pos int) *Block {
	block := &Block{}
	block.Stringer = block
	block.Label = label
	block.pos = pos
	block.objects = make(map[interface{}]*Var)
	block.types = b.types

	block.scope = b
	return block
}

// EmitBlock 在 Block 中添加一个子 Block
func (b *Block) EmitBlock(label string, pos int) *Block {
	block := b.createBlock(label, pos)

	b.emit(block)
	return block
}

// emit 向 Block 中添加一个指令
func (b *Block) emit(stmt Stmt) {
	b.Stmts = append(b.Stmts, stmt)
}
