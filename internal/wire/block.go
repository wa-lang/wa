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
	anInstruction
	Comment string // 附加注释
	//Locals  []Value       // 该块内定义的局部变量
	Instrs []Instruction // 该块所含的指令

	objects map[interface{}]Location // 关联 AST 结点 -> 块内值
	types   *Types                   // 该函数所属 Module 的类型库，切勿手动修改
}

// 初始化 Block
func (b *Block) Init() {
	b.objects = make(map[interface{}]Location)
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
func (b *Block) AddLocal(name string, typ ValueType, pos int, obj interface{}) Location {
	v := &InstAlloc{}
	v.Stringer = v
	v.typ = b.types.GenPtr(typ)
	v.dataType = typ
	v.pos = pos
	v.object = obj
	if obj != nil {
		b.objects[obj] = v
	}

	b.emit(v)
	return v
}

// EmitLoad 在 Block 中添加一条 Load 指令
func (b *Block) EmitLoad(loc Location, typ ValueType, pos int) *InstLoad {
	v := &InstLoad{}
	v.Stringer = v
	v.Loc = loc
	v.typ = typ
	v.pos = pos

	b.emit(v)
	return v
}

// EmitStore 在 Block 中添加一条 Store 指令
func (b *Block) EmitStore(loc Location, val Value, pos int) *InstStore {
	v := &InstStore{}
	v.Stringer = v
	v.Loc = loc
	v.Val = val
	v.pos = pos

	b.emit(v)
	return v
}

// EmitExtract 在 Block 中添加一条 Extract 指令
func (b *Block) EmitExtract(tuple Value, index int, pos int) *InstExtract {
	v := &InstExtract{}
	v.Stringer = v
	v.Tuple = tuple
	v.Index = index
	v.pos = pos

	b.emit(v)
	return v
}

func (b *Block) EmitReturn(results []Value, pos int) *InstReturn {
	v := &InstReturn{}
	v.Stringer = v
	v.Results = results
	v.pos = pos

	b.emit(v)
	return v
}

// Emit
func (b *Block) EmitBlock(comment string, pos int) *Block {
	block := &Block{}
	block.Stringer = block
	block.Comment = comment
	block.pos = pos
	block.objects = make(map[interface{}]Location)
	block.types = b.types

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
