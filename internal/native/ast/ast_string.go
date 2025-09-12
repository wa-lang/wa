// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package ast

import (
	"fmt"
	"strings"

	"wa-lang.org/wa/internal/native/token"
)

func (p *Comment) String() string {
	return fmt.Sprintf("# %s", p.Text)
}

func (p *CommentGroup) String() string {
	var sb strings.Builder
	for i, c := range p.List {
		if i > 0 {
			sb.WriteRune('\n')
		}
		if p.TopLevel {
			sb.WriteByte('\t')
		}
		sb.WriteString(c.String())
	}
	return sb.String()
}

func (p *BasicLit) String() string {
	if p.TypeCast != token.NONE && p.TypeCast != p.LitKind.DefaultNumberType() {
		return fmt.Sprintf("%v(%s)", p.TypeCast, p.LitString)
	}
	return p.LitString
}

func (p *Const) String() string {
	var sb strings.Builder
	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteRune('\n')
	}
	sb.WriteString(fmt.Sprintf("const %s = %v", p.Name, p.Value))
	if p.Comment != nil {
		sb.WriteString(p.Comment.String())
	}
	return sb.String()
}

func (p *Global) String() string {
	var sb strings.Builder
	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteRune('\n')
	}

	if p.Type != token.NONE {
		sb.WriteString(fmt.Sprintf("global %s:%v = ", p.Name, p.Type))
	} else {
		sb.WriteString(fmt.Sprintf("global %s:%d = ", p.Name, p.Size))
	}

	switch {
	case len(p.Init) == 0:
		sb.WriteString("{}")
	case len(p.Init) == 1 && p.Init[0].Doc == nil && p.Init[0].Offset == 0:
		sb.WriteString(p.Init[0].String())
	default:
		sb.WriteString("{")
		for i, xInit := range p.Init {
			if i > 0 {
				sb.WriteByte('\n')
			}
			sb.WriteString(xInit.String())
		}
		sb.WriteString("}")
	}

	return sb.String()
}

func (p *InitValue) String() string {
	var sb strings.Builder
	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteByte('\n')
	}
	if p.Lit != nil {
		sb.WriteByte('\t')
		sb.WriteString(p.Lit.String())
		sb.WriteByte(',')
	} else {
		sb.WriteByte('\t')
		sb.WriteString(p.Symbal)
		sb.WriteByte(',')
	}
	if p.Comment != nil {
		sb.WriteString(p.Comment.String())
	}
	return sb.String()
}

func (p *Func) String() string        { panic("TODO") }
func (p *FuncType) String() string    { panic("TODO") }
func (p *FuncBody) String() string    { panic("TODO") }
func (p *Argument) String() string    { panic("TODO") }
func (p *Local) String() string       { panic("TODO") }
func (p *Instruction) String() string { panic("TODO") }
