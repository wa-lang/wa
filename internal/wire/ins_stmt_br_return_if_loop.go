// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"strings"
)

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

	v.True = b.newBlock("", pos)
	v.False = b.newBlock("", pos)

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
	PreCond []Stmt
	Cond    Expr   // 循环条件
	Label   string //
	Body    *Block // 循环体，不会为 nil
	Post    *Block // 循环后处理，不会为 nil
}

func (i *Loop) Format(tab string, sb *strings.Builder) {
	sb.WriteString(tab)
	sb.WriteString("loop ")
	for j, stmt := range i.PreCond {
		sb.WriteString(stmt.String())
		if j < len(i.PreCond) {
			sb.WriteString("; ")
		}
	}
	sb.WriteString(i.Cond.Name())
	sb.WriteString(" : ")
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

	v.Body = b.newBlock(label+".body", pos)
	v.Post = b.newBlock(label+".post", pos)

	b.emit(v)
	return v
}
