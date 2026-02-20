// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

/**************************************
aStmt: 实现 Stmt 接口与 Pos 相关的方法
包含 aStmt 的对象必须自行实现 Stringer 接口！
**************************************/

type aStmt struct {
	fmt.Stringer
	pos int
}

func (v *aStmt) Pos() int { return v.pos }
func (v *aStmt) Format(tab string, sb *strings.Builder) {
	sb.WriteString(tab)
	sb.WriteString(v.String())
}

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
Var: Var 指令，定义一个变量
**************************************/

type Var struct {
	aStmt
	kind   VarKind
	name   string
	dtype  Type
	rtype  Type
	object interface{} // 与该值关联的 AST 结点。对凹语言前端，应为 types.Object

	tank *tank
}

func (i *Var) Name() string { return i.name }
func (i *Var) Type() Type {
	if i.kind == Register {
		return i.dtype
	} else {
		return i.rtype
	}
}
func (i *Var) retained() bool { return false }
func (i *Var) String() string {
	s := ""
	switch i.kind {
	case Register:
		s = fmt.Sprintf("var %s %s", i.name, i.dtype.Name())

	case Heap:
		s = fmt.Sprintf("var %s %s = alloc.heap(%s)", i.name, i.rtype.Name(), i.dtype.Name())

	default:
		panic(fmt.Sprintf("Todo: VarKind: %v", i.kind))
	}
	if i.tank != nil {
		s += " --- "
		s += i.tank.String()
	}
	return s
}
func (i *Var) DataType() Type { return i.dtype }
func (i *Var) RefType() Type  { return i.rtype }

//func (i *Var) Object() interface{} { return i.object }

// AddLocal 在 Block 中分配一个局部变量，初始时位于 Register
func (b *Block) AddLocal(name string, typ Type, pos int, obj interface{}) *Var {
	v := &Var{}
	v.Stringer = v
	v.name = name
	v.dtype = typ
	v.rtype = b.types.GenRef(typ)

	v.pos = pos
	v.object = obj
	if obj != nil {
		b.objects[obj] = v
	}

	b.emit(v)
	return v
}

/**************************************
Get: Get 指令，获取变量 Loc 的值，Get 实现了 Expr 接口
  - Loc 应为 *Var，或类型为 Ref、Ptr 的 Expr
**************************************/

type Get struct {
	aStmt
	Loc Expr
}

func (i *Get) Name() string { return i.String() }
func (i *Get) Type() Type {
	if v, ok := i.Loc.(*Var); ok {
		return v.DataType()
	}

	switch t := i.Loc.Type().(type) {
	case *Ref:
		return t.Base

	case *Ptr:
		return t.Base

	default:
		panic(fmt.Sprintf("Invalid Loc.Type():%s", i.Loc.Type().Name()))
	}
}
func (i *Get) retained() bool { return false }
func (i *Get) String() string {
	if v, ok := i.Loc.(*Var); ok {
		if v.kind == Register {
			return v.name
		}
	}

	return fmt.Sprintf("*(%s)", i.Loc.Name())
}

// 生成一条 Get 指令
func NewGet(loc Expr, pos int) *Get {
	v := &Get{}
	v.Stringer = v
	v.Loc = loc
	v.pos = pos
	return v
}

/**************************************
Set: Set 指令，将 Val 存储到 Loc 指定的位置，Set 支持多赋值，该指令应触发 RC-1 动作
  - Set 支持多赋值
  - Loc 中的元素应为 *Var，或类型为 Ref、Ptr 的 Expr
  - Rh 可能为元组（Tuple），此时 Rhs 的长度应为 1，Lhs 的长度应为元组长度
  - 向 nil 的 Lh 赋值是合法的，这等价于向匿名变量 _ 赋值，此时若 Rh 已 retain，应触发 release
**************************************/

type Set struct {
	aStmt
	Lhs []Expr
	Rhs []Expr
}

func (i *Set) String() string {
	var sb strings.Builder

	for i, lh := range i.Lhs {
		if i > 0 {
			sb.WriteString(", ")
		}

		if lh == nil {
			sb.WriteRune('_')
			continue
		}

		loc_name := "*(" + lh.Name() + ")"
		if v, ok := lh.(*Var); ok {
			if v.kind == Register {
				loc_name = lh.Name()
			}
		}
		sb.WriteString(loc_name)
	}

	sb.WriteString(" = ")
	for i, rh := range i.Rhs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(rh.Name())
	}

	return sb.String()
}

func NewSet(lhs []Expr, rhs []Expr, pos int) *Set {
	v := &Set{}
	v.Stringer = v
	v.Lhs = lhs
	v.Rhs = rhs
	v.pos = pos
	return v
}

// 在 Block 中添加一条 Set 指令
func (b *Block) EmitSet(lhs Expr, rhs Expr, pos int) *Set {
	v := NewSet([]Expr{lhs}, []Expr{rhs}, pos)
	b.emit(v)
	return v
}

// Block.EmitSet 的多重赋值版
func (b *Block) EmitSetN(lhs []Expr, rhs []Expr, pos int) *Set {
	v := NewSet(lhs, rhs, pos)
	b.emit(v)
	return v
}

/**************************************
Assign: Assign 指令，将 Rhs 赋值给 Lhs
  - Assgin 支持多赋值
  - Lh 必须为 Register 型变量
  - Rh 可能为元组（Tuple），此时 Rhs 的长度应为 1，Lhs 的长度应为元组长度
  - 向 nil 的 Lh 赋值是合法的，这等价于向匿名变量 _ 赋值，此时若 Rh 已 retain，应触发 release
**************************************/

type Assign struct {
	aStmt
	Lhs []*Var
	Rhs []Expr
}

func (i *Assign) String() string {
	var sb strings.Builder

	for i, lh := range i.Lhs {
		if i > 0 {
			sb.WriteString(", ")
		}

		if lh == nil {
			sb.WriteRune('_')
			continue
		}

		sb.WriteString(lh.Name())
	}

	sb.WriteString(" = ")
	for i, rh := range i.Rhs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(rh.Name())
	}

	return sb.String()
}

func NewAssign(lhs []*Var, rhs []Expr, pos int) *Assign {
	v := &Assign{}
	v.Stringer = v
	v.Lhs = lhs
	v.Rhs = rhs
	v.pos = pos
	return v
}

// 在 Block 中添加一条 Assign 指令
func (b *Block) EmitAssign(lhs *Var, rhs Expr, pos int) *Assign {
	v := NewAssign([]*Var{lhs}, []Expr{rhs}, pos)
	b.emit(v)
	return v
}

// EmitAssign 的多重赋值版
func (b *Block) EmitAssignN(lhs []*Var, rhs []Expr, pos int) *Assign {
	v := NewAssign(lhs, rhs, pos)
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

func NewBr(label string, pos int) *Br {
	v := &Br{}
	v.Stringer = v
	v.Label = label
	v.pos = pos
	return v
}

// 在 Block 中添加一条 Br 指令
func (b *Block) EmitBr(label string, pos int) *Br {
	v := NewBr(label, pos)
	b.emit(v)
	return v
}

/**************************************
Return: Return 指令，函数返回，该指令只能出现在 Block 末尾
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
func NewUnop(x Expr, op OpCode, pos int) *Unop {
	v := &Unop{X: x, Op: op}
	v.Stringer = v
	v.pos = pos

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
		return &Int{}

	case EQL, NEQ, GTR, LSS, GEQ, LEQ:
		return &Bool{}

	default:
		return i.X.Type()
	}
}
func (i *Biop) retained() bool { return i.Type().Kind() == TypeKindString }
func (i *Biop) String() string {
	return fmt.Sprintf("(%s %s %s)", i.X.Name(), OpCodeNames[i.Op], i.Y.Name())
}

// 生成一条 Biop 指令
func NewBiop(x, y Expr, op OpCode, pos int) *Biop {
	v := &Biop{X: x, Y: y, Op: op}
	v.Stringer = v
	v.pos = pos
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

/**************************************
Retain: Retain 指令，引用计数 +1，Retain 指令实现了 Expr，返回 X 本身
**************************************/

type Retain struct {
	aStmt
	X Expr
}

func (i *Retain) Name() string   { return i.String() }
func (i *Retain) Type() Type     { return i.X.Type() }
func (i *Retain) retained() bool { panic("") }
func (i *Retain) String() string { return fmt.Sprintf("retain(%s)", i.X.Name()) }

// 生成一条 Retain 指令
func NewRetain(x Expr, pos int) *Retain {
	v := &Retain{X: x}
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
	X *Var
}

func (i *Drop) String() string {
	if i.X.tank != nil {
		return fmt.Sprintf("drop(%s --- %s)", i.X.Name(), i.X.tank.String())
	} else {
		return fmt.Sprintf("drop(%s)", i.X.Name())
	}
}

// 生成一条 Drop 指令
func NewDrop(x *Var, pos int) *Drop {
	v := &Drop{X: x}
	v.Stringer = v
	v.pos = pos
	return v
}

/**************************************
DupRef: DupRef 指令，引用复制，DupRef 指令实现了 Expr，返回 X
**************************************/

type DupRef struct {
	aStmt
	X   Expr
	Imv *Var
}

func (i *DupRef) Name() string   { return i.String() }
func (i *DupRef) Type() Type     { return i.X.Type() }
func (i *DupRef) retained() bool { panic("") }
func (i *DupRef) String() string { return fmt.Sprintf("dupref(%s, %s)", i.X.Name(), i.Imv.Name()) }

// 生成一条 DupRef 指令
func NewDupRef(x Expr, imvName string, pos int) *DupRef {
	v := &DupRef{X: x}
	v.Stringer = v
	v.pos = pos

	imv := &Var{name: imvName}
	imv.Stringer = imv
	imv.dtype = x.Type()
	v.Imv = imv

	return v
}
