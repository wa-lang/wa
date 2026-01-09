// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

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

func (i *Alloc) Name() string { return i.name }
func (i *Alloc) Type() Type {
	if i.Location == LocationKindLocal {
		return i.dataType
	} else {
		return i.refType
	}
}
func (i *Alloc) retained() bool { return false }
func (i *Alloc) String() string {
	switch i.Location {
	case LocationKindLocal:
		return fmt.Sprintf("var %s %s", i.name, i.dataType.Name())

	case LocationKindStack:
		return fmt.Sprintf("var %s %s = alloc.stack(%s)", i.name, i.refType.Name(), i.dataType.Name())

	case LocationKindHeap:
		return fmt.Sprintf("var %s %s = alloc.heap(%s)", i.name, i.refType.Name(), i.dataType.Name())
	}
	panic(fmt.Sprintf("Invalid LocationType: %v", i.Location))
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

func (i *Load) Name() string   { return i.String() }
func (i *Load) Type() Type     { return i.Loc.DataType() }
func (i *Load) retained() bool { return false }
func (i *Load) String() string {
	if i.Loc.LocationKind() == LocationKindLocal {
		return i.Loc.Name()
	} else {
		return "*" + i.Loc.Name()
	}
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

		if loc == nil {
			sb.WriteRune('_')
			continue
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

///**************************************
//Extract: Extract 指令，提取元组 Tuple 的第 Index 个元素，Extract 实现了 Expr
//**************************************/
//type Extract struct {
//	aStmt
//	X     Expr
//	Index int
//}
//
//func (i *Extract) Name() string {
//	return i.String()
//}
//
//func (i *Extract) String() string {
//	return fmt.Sprintf("extract(%s, #%d)", i.X.Name(), i.Index)
//}
//
//func (i *Extract) Type() Type {
//	return i.X.Type().(*Tuple).fields[i.Index]
//}
//
//// 生成一条 Extract 指令
//func (b *Block) NewExtract(tuple Expr, index int, pos int) *Extract {
//	v := &Extract{}
//	v.Stringer = v
//	v.X = tuple
//	v.Index = index
//	v.pos = pos
//	v.setScope(b)
//
//	return v
//}

/**************************************
Br: Br 指令
**************************************/
type Br struct {
	aStmt
	Label string
}

func (i *Br) String() string {
	s := "br " + i.Label
	return s
}

// 在 Block 中添加一条 Br 指令
func (b *Block) EmitBr(label string, pos int) *Br {
	v := &Br{}
	v.Stringer = v
	v.Label = label
	v.pos = pos

	b.emit(v)
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

func (i *Unop) Name() string   { return i.String() }
func (i *Unop) Type() Type     { return i.X.Type() }
func (i *Unop) retained() bool { return false }
func (i *Unop) String() string { return fmt.Sprintf("%s%s", OpCodeNames[i.Op], i.X.Name()) }

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

func (i *Biop) Name() string { return i.String() }
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
func (i *Biop) retained() bool { return i.Type().Kind() == TypeKindString }
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

/**************************************
If: 条件指令
**************************************/
type If struct {
	aStmt
	Cond  Expr   // 判断条件
	True  *Block // 为 true 时的分支，不会为 nil
	False *Block // 为 false 时的分支，不会为 nil
}

func (i *If) Format(tab string, sb *strings.Builder) {
	sb.WriteString(tab)
	sb.WriteString("if ")
	sb.WriteString(i.Cond.Name())
	sb.WriteString("\n")

	i.True.Format(tab, sb)
	sb.WriteString("\n")

	sb.WriteString(tab)
	sb.WriteString("else\n")
	i.False.Format(tab, sb)
}

// 在 Block 中添加一条 If 指令
func (b *Block) EmitIf(cond Expr, pos int) *If {
	if !cond.Type().Equal(b.types.Bool) {
		panic("cond must be bool.")
	}

	v := &If{Cond: cond}
	v.Stringer = v
	v.pos = pos

	v.True = b.createBlock("", pos)
	v.False = b.createBlock("", pos)

	b.emit(v)
	return v
}

/**************************************
Loop: 循环指令，逻辑如下：
loop $Label {
	if cond_expr $Label.done {
		block $Label.body {
			...body...
		}  // <- continue 转这里
		...post...
		br $Label
	}  // <- break 转这里
}
**************************************/
type Loop struct {
	aStmt
	Cond  Expr   // 循环条件
	Label string //
	Body  *Block // 循环体，不会为 nil
	Post  *Block // 循环后处理，不会为 nil
}

func (i *Loop) Format(tab string, sb *strings.Builder) {
	sb.WriteString(tab)
	sb.WriteString("loop ")
	sb.WriteString(i.Cond.Name())
	sb.WriteString(" $")
	sb.WriteString(i.Label)
	sb.WriteString("\n")

	i.Body.Format(tab, sb)

	sb.WriteString(" post\n")
	i.Post.Format(tab, sb)
}

// 在 Block 中添加一条 Loop 指令
func (b *Block) EmitLoop(cond Expr, label string, pos int) *Loop {
	if !cond.Type().Equal(b.types.Bool) {
		panic("cond must be bool.")
	}

	v := &Loop{Cond: cond}
	v.Stringer = v
	v.pos = pos

	v.Body = b.createBlock(label+".body", pos)
	v.Post = b.createBlock(label+".post", pos)

	b.emit(v)
	return v
}
