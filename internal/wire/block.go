// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
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
**************************************/

type Block struct {
	aStmt
	Label string // 标签
	Stmts []Stmt // 该块所含的指令

	scope   Scope
	objects map[interface{}]*Var // AST 结点 -> 块内变量
	types   *Types               // 该函数所属 Module 的类型库，切勿手动修改

	imvCount int
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

// 判断 Var 是否被赋值
func (b *Block) varStored(v *Var) bool {
	for _, stmt := range b.Stmts {
		if varStoredInStmt(stmt, v) {
			return true
		}
	}

	return false
}

func varStoredInStmt(stmt Stmt, v *Var) bool {
	switch s := stmt.(type) {
	case *Block:
		return s.varStored(v)

	case *Set:
		for _, lh := range s.Lhs {
			if p, _ := lh.(*Var); p == v {
				return true
			}
		}

	case *If:
		if s.True.varStored(v) || s.False.varStored(v) {
			return true
		}

	case *Loop:
		if s.Body.varStored(v) || s.Post.varStored(v) {
			return true
		}

	case *Var:
	case *Get:
	case *Br:
	case *Return:
	case *Unop:
	case *Biop:
	case *Retain:
	case *Drop:
		return false

	default:
		panic(fmt.Sprintf("Todo: %s", s.String()))
	}

	return false
}

// usageRange 使用范围信息
type usageRange struct {
	first int // -1表示未使用
	last  int // -1表示未使用
}

// varUsageRange 分析指定变量在该块中的使用情况，变量声明本身不属于被使用
func (b *Block) varUsageRange(v *Var) usageRange {
	info := usageRange{
		first: -1,
		last:  -1,
	}

	for i, stmt := range b.Stmts {
		// 检查语句是否使用了变量
		if varUsedInStmt(stmt, v) {
			if info.first == -1 {
				info.first = i
			}
			info.last = i
		}
	}

	return info
}

// varUsedInStmt 检查语句是否使用了指定的变量
func varUsedInStmt(stmt Stmt, v *Var) bool {
	switch s := stmt.(type) {
	case *Block:
		for _, ss := range s.Stmts {
			if varUsedInStmt(ss, v) {
				return true
			}
		}

	case *Var:
		return false

	case *Get:
		return exprContainsVar(s.Loc, v)

	case *Set:
		for _, lh := range s.Lhs {
			if exprContainsVar(lh, v) {
				return true
			}
		}
		for _, rh := range s.Rhs {
			if exprContainsVar(rh, v) {
				return true
			}
		}
		return false

	case *Br:
		return false

	case *Return:
		for _, r := range s.Results {
			if exprContainsVar(r, v) {
				return true
			}
		}

	case *Unop:
		return exprContainsVar(s.X, v)

	case *Biop:
		return exprContainsVar(s.X, v) || exprContainsVar(s.Y, v)

	case *Call:
		var call_common *CallCommon
		switch call := s.Callee.(type) {
		case *BuiltinCall:
			call_common = &call.CallCommon

		case *StaticCall:
			call_common = &call.CallCommon

		case *MethodCall:
			if exprContainsVar(call.Recv, v) {
				return true
			}
			call_common = &call.CallCommon

		case *InterfaceCall:
			if exprContainsVar(call.Interface, v) {
				return true
			}
			call_common = &call.CallCommon

		default:
			panic("Todo")
		}

		for _, arg := range call_common.Args {
			if exprContainsVar(arg, v) {
				return true
			}
		}

	case *If:
		return exprContainsVar(s.Cond, v) || varUsedInStmt(s.True, v) || varUsedInStmt(s.False, v)

	case *Loop:
		return exprContainsVar(s.Cond, v) || varUsedInStmt(s.Body, v) || varUsedInStmt(s.Post, v)

	case *Retain:
		return exprContainsVar(s.X, v)

	case *Drop:
		return false

	default:
		panic(fmt.Sprintf("Todo: %s", s.String()))
	}

	return false
}

// exprContainsVar 检查表达式是否包含指定的变量
func exprContainsVar(expr Expr, v *Var) bool {
	if expr == nil {
		return false
	}

	switch e := expr.(type) {
	case *Const:
		return false

	case *Var:
		return e == v

	case *Get:
		return exprContainsVar(e.Loc, v)

	case *Unop:
		return exprContainsVar(e.X, v)

	case *Biop:
		return exprContainsVar(e.X, v) || exprContainsVar(e.Y, v)

	case *Call:
		var call_common *CallCommon
		switch call := e.Callee.(type) {
		case *BuiltinCall:
			call_common = &call.CallCommon

		case *StaticCall:
			call_common = &call.CallCommon

		case *MethodCall:
			if exprContainsVar(call.Recv, v) {
				return true
			}
			call_common = &call.CallCommon

		case *InterfaceCall:
			if exprContainsVar(call.Interface, v) {
				return true
			}
			call_common = &call.CallCommon

		default:
			panic("Todo")
		}

		for _, arg := range call_common.Args {
			if exprContainsVar(arg, v) {
				return true
			}
		}

	case *Retain:
		return exprContainsVar(e.X, v)

	case *DupRef:
		return exprContainsVar(e.X, v)

	default:
		panic(fmt.Sprintf("Todo: %s", e.Name()))
	}
	return false
}
