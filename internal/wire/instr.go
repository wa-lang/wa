// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

/**************************************
Stmt: 指令接口
**************************************/
type Stmt interface {
	// 获取该指令的伪代码
	String() string

	// 获取该指令在源码中的位置
	Pos() int

	// 格式化输出
	Format(tab string, sb *strings.Builder)

	// 获取该指令所属的域
	ParentScope() Scope

	// 设置该指令所属的域
	setScope(Scope)
}

/**************************************
aStmt: 实现 Stmt 接口与 Pos、Scope 相关的方法
包含 aStmt 的对象必须自行实现 Stringer 接口！
**************************************/
type aStmt struct {
	fmt.Stringer
	pos   int
	scope Scope
}

func (v *aStmt) ParentScope() Scope { return v.scope }
func (v *aStmt) setScope(s Scope)   { v.scope = s }
func (v *aStmt) Pos() int           { return v.pos }
func (v *aStmt) Format(tab string, sb *strings.Builder) {
	sb.WriteString(tab)
	sb.WriteString(v.String())
}

/**************************************
Expr: 表达式，所有可以作为指令参数的对象，都满足该接口
**************************************/
type Expr interface {
	// 表达式的名字
	// 变量的名字是其变量名，常量的名字是其字面量，除此外多数表达式的名字是其指令伪代码
	Name() string

	// 该表达式的类型
	Type() Type

	// 表达式在源码中的位置
	Pos() int
}

/**************************************
Alloc: Alloc 指令，分配一个变量位置，该对象同时实现了 Location 接口
**************************************/
type Alloc struct {
	aStmt
	Location LocationKind
	name     string
	dataType Type
	refType  Type
	object   interface{} // 与该值关联的 AST 结点。对凹语言前端，应为 types.Object
}

func (i *Alloc) String() string {
	switch i.Location {
	case LocationKindLocal:
		return fmt.Sprintf("var %s %s", i.dataType.Name(), i.name)

	case LocationKindStack:
		return fmt.Sprintf("var %s %s = alloc.stack(%s)", i.refType.Name(), i.name, i.dataType.Name())

	case LocationKindHeap:
		return fmt.Sprintf("var %s %s = alloc.heap(%s)", i.refType.Name(), i.name, i.dataType.Name())
	}
	panic(fmt.Sprintf("Invalid LocationType: %v", i.Location))
}

func (i *Alloc) Name() string {
	return i.name
}

func (i *Alloc) Type() Type {
	if i.Location == LocationKindLocal {
		return i.dataType
	} else {
		return i.refType
	}
}

// Location 接口相关
func (i *Alloc) LocationKind() LocationKind { return i.Location }
func (i *Alloc) DataType() Type             { return i.dataType }
func (i *Alloc) Object() interface{}        { return i.object }

// AddLocal 在 Block 中分配一个局部变量
func (b *Block) AddLocal(name string, typ Type, pos int, obj interface{}) Location {
	v := &Alloc{}
	v.Stringer = v
	v.name = name
	v.dataType = typ
	v.refType = b.types.GenPtr(typ)
	v.pos = pos
	v.object = obj
	if obj != nil {
		b.objects[obj] = v
	}

	b.emit(v)
	return v
}

/**************************************
Load: Load 指令，装载 Loc 处的变量，Load 实现了 Expr 接口
**************************************/
type Load struct {
	aStmt
	Loc Location
}

func (i *Load) Name() string {
	return i.String()
}

func (i *Load) String() string {
	if i.Loc.LocationKind() == LocationKindLocal {
		return i.Loc.Name()
	} else {
		return "*" + i.Loc.Name()
	}
}

func (i *Load) Type() Type {
	return i.Loc.DataType()
}

// 生成一条 Load 指令
func (b *Block) NewLoad(loc Location, pos int) *Load {
	v := &Load{}
	v.Stringer = v
	v.Loc = loc
	v.pos = pos
	v.setScope(b)

	return v
}

/**************************************
Store: Store 指令，将 Val 存储到 Loc 指定的位置，Store 支持多赋值，该指令应触发 RC+1 动作
 - val 有可能为元组（Tuple），因此 Loc 的长度可能和 Val 的长度不同，此时应将元组展开，完全展开后二者的长度应一致
 - 向 nil 的 loc 赋值是合法的，这等价于向匿名变量 _ 赋值，此时应触发 Drop 动作
**************************************/
type InstStore struct {
	aStmt
	Loc []Location
	Val []Expr
}

func (i *InstStore) String() string {
	var sb strings.Builder
	for i, loc := range i.Loc {
		if i > 0 {
			sb.WriteString(", ")
		}
		if loc.LocationKind() == LocationKindLocal {
			sb.WriteString(loc.Name())
		} else {
			sb.WriteRune('*')
			sb.WriteString(loc.Name())
		}
	}

	sb.WriteString(" = ")
	for i, val := range i.Val {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(val.Name())
	}

	return sb.String()
}

// 在 Block 中添加一条 Store 指令
func (b *Block) EmitStore(loc Location, val Expr, pos int) *InstStore {
	v := &InstStore{}
	v.Stringer = v
	v.Loc = []Location{loc}
	v.Val = []Expr{val}
	v.pos = pos

	b.emit(v)
	return v
}

// Block.EmitStore 的多重赋值版
func (b *Block) EmitStoreN(locs []Location, vals []Expr, pos int) *InstStore {
	v := &InstStore{}
	v.Stringer = v
	v.Loc = locs
	v.Val = vals
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
Extract: Extract 指令，提取元组 Tuple 的第 Index 个元素，Extract 实现了 Expr
**************************************/
type Extract struct {
	aStmt
	X     Expr
	Index int
}

func (i *Extract) Name() string {
	return i.String()
}

func (i *Extract) String() string {
	return fmt.Sprintf("extract(%s, #%d)", i.X.Name(), i.Index)
}

func (i *Extract) Type() Type {
	return i.X.Type().(*Tuple).fields[i.Index]
}

// 生成一条 Extract 指令
func (b *Block) NewExtract(tuple Expr, index int, pos int) *Extract {
	v := &Extract{}
	v.Stringer = v
	v.X = tuple
	v.Index = index
	v.pos = pos
	v.setScope(b)

	return v
}

/**************************************
Return: Return 指令
**************************************/
type Return struct {
	aStmt
	Results []Expr
}

func (i *Return) String() string {
	s := "return"
	for m, r := range i.Results {
		if m > 0 {
			s += ","
		}
		s += " "
		s += r.Name()
	}
	return s
}

// 在 Block 中添加一条 Return 指令
func (b *Block) EmitReturn(results []Expr, pos int) *Return {
	v := &Return{}
	v.Stringer = v
	v.Results = results
	v.pos = pos

	b.emit(v)
	return v
}

/**************************************
OpCode: 运算符
**************************************/
type OpCode int

const (
	NOT OpCode = iota
	NEG
	XOR
	LAND
	LOR
	SHL
	SHR
	ADD
	SUB
	MUL
	QUO
	REM
	AND
	OR
	ANDNOT
	EQL
	NEQ
	GTR
	LSS
	GEQ
	LEQ
	LEG
)

var OpCodeNames = [...]string{
	NOT:    "!",
	NEG:    "-",
	XOR:    "^",
	LAND:   "&&",
	LOR:    "||",
	SHL:    "<<",
	SHR:    ">>",
	ADD:    "+",
	SUB:    "-",
	MUL:    "*",
	QUO:    "/",
	REM:    "%",
	AND:    "&",
	ANDNOT: "&^",
	OR:     "|",
	EQL:    "==",
	NEQ:    "!=",
	GTR:    ">",
	LSS:    "<",
	GEQ:    ">=",
	LEQ:    "<=",
	LEG:    "<=>",
}

/**************************************
Unop: 单目运算，算符范围[NOT, XOR]，Unop 实现了 Expr
**************************************/
type Unop struct {
	aStmt
	X  Expr
	Op OpCode
}

func (i *Unop) Name() string {
	return i.String()
}

func (i *Unop) Type() Type {
	return i.X.Type()
}

func (i *Unop) String() string {
	return fmt.Sprintf("%s%s", OpCodeNames[i.Op], i.X.Name())
}

// 生成一条 Unop 指令
func (b *Block) NewUnop(x Expr, op OpCode, pos int) *Unop {
	v := &Unop{X: x, Op: op}
	v.Stringer = v
	v.pos = pos
	v.setScope(b)

	return v
}

/**************************************
Biop: 双目运算，算符范围[XOR, LEG]，BiOp 实现了 Expr
**************************************/
type Biop struct {
	aStmt
	X, Y Expr
	Op   OpCode
}

func (i *Biop) Name() string {
	return i.String()
}

func (i *Biop) Type() Type {
	switch i.Op {
	case LEG:
		return i.scope.(*Block).types.Int

	case EQL, NEQ, GTR, LSS, GEQ, LEQ:
		return i.scope.(*Block).types.Bool

	default:
		return i.X.Type()
	}
}

func (i *Biop) String() string {
	return fmt.Sprintf("(%s %s %s)", i.X.Name(), OpCodeNames[i.Op], i.Y.Name())
}

// 生成一条 Biop 指令
func (b *Block) NewBiop(x, y Expr, op OpCode, pos int) *Biop {
	v := &Biop{X: x, Y: y, Op: op}
	v.Stringer = v
	v.pos = pos
	v.setScope(b)

	return v
}

///**************************************
//InstCall:  函数调用指令，Call 为 StaticCall、BuiltinCall、MethodCall、InterfaceCall、ClosureCall 之一
//**************************************/
//type InstCall struct {
//	aImv
//	Call Value
//}
//
//func (i *InstCall) Type() Type     { return i.Call.Type() }
//func (i *InstCall) Pos() int       { return i.Call.Pos() }
//func (i *InstCall) String() string { return i.Call.(fmt.Stringer).String() }
//
//// 在 Block 中添加一条 InstCall 指令，Call 为 StaticCall、BuiltinCall、MethodCall、InterfaceCall、ClosureCall 之一
//func (b *Block) EmitInstCall(call Value) *InstCall {
//	v := &InstCall{Call: call}
//	v.Stringer = v
//	v.pos = call.Pos()
//
//	b.emit(v)
//	return v
//}
//
///**************************************
//InstIf:  条件指令
//**************************************/
//type InstIf struct {
//	aImv
//	Cond  Value  // 判断条件
//	True  *Block // 为 true 时的分支，不会为 nil
//	False *Block // 为 false 时的分支，不会为 nil
//}
//
//func (i *InstIf) Type() Type {
//	return i.True.Type()
//}
//
//func (i *InstIf) Format(tab string, sb *strings.Builder) {
//	sb.WriteString(tab)
//	sb.WriteString("if ")
//	sb.WriteString(i.Cond.Name())
//	sb.WriteString("\n")
//
//	i.True.Format(tab, sb)
//
//	sb.WriteString(tab)
//	sb.WriteString("else\n")
//	i.False.Format(tab, sb)
//}
//
//// 在 Block 中添加一条 InstIf 指令
//func (b *Block) EmitInstIf(cond Value, typ Type, pos int) *InstIf {
//	if !cond.Type().Equal(b.types.Bool) {
//		panic("cond must be bool.")
//	}
//
//	v := &InstIf{Cond: cond}
//	v.Stringer = v
//	v.pos = pos
//
//	v.True = b.createBlock("", typ, pos)
//	v.False = b.createBlock("", typ, pos)
//
//	b.emit(v)
//	return v
//}
//
