// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
)

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
Combo: 组合指令，将多个指令组合成一个指令，实现了 Expr 接口，返回 Result
**************************************/
//
//type Combo struct {
//	aStmt
//	Stmts  []Stmt
//	Result Expr
//}
//
//func (i *Combo) Name() string   { return i.String() }
//func (i *Combo) Type() Type     { return i.Result.Type() }
//func (i *Combo) retained() bool { return i.Result.retained() }
//
//func (i *Combo) String() string {
//	var sb strings.Builder
//	sb.WriteRune('{')
//	for _, stmt := range i.Stmts {
//		sb.WriteString(stmt.String())
//		sb.WriteString("; ")
//	}
//	sb.WriteString(i.Result.Name())
//	sb.WriteRune('}')
//	return sb.String()
//}
//
//func NewCombo(stmts []Stmt, result Expr, pos int) *Combo {
//	v := &Combo{Stmts: stmts, Result: result}
//	v.Stringer = v
//	v.pos = pos
//	return v
//}

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
DupRef: DupRef 指令，引用复制，DupRef 指令实现了 Expr，返回 X
**************************************/

//type DupRef struct {
//	aStmt
//	X   Expr
//	Imv *Var
//}
//
//func (i *DupRef) Name() string   { return i.String() }
//func (i *DupRef) Type() Type     { return i.X.Type() }
//func (i *DupRef) retained() bool { panic("") }
//func (i *DupRef) String() string { return fmt.Sprintf("dupref(%s, %s)", i.X.Name(), i.Imv.Name()) }
//
//// 生成一条 DupRef 指令
//func NewDupRef(x Expr, imvName string, pos int) *DupRef {
//	v := &DupRef{X: x}
//	v.Stringer = v
//	v.pos = pos
//
//	imv := &Var{name: imvName}
//	imv.Stringer = imv
//	imv.dtype = x.Type()
//	v.Imv = imv
//
//	return v
//}
