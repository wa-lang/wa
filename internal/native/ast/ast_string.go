// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package ast

import (
	"fmt"
	"strings"

	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/token"
)

func (p *File) String() string {
	var sb strings.Builder

	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteRune('\n')
	}

	for _, obj := range p.Objects {
		sb.WriteString(obj.String())
		sb.WriteString("\n\n")
	}

	return sb.String()
}

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

	// TODO: 保留孤立注释的顺序

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

func (p *Func) String() string {
	var sb strings.Builder

	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteRune('\n')
	}
	sb.WriteString("func ")
	sb.WriteString(p.Name)
	sb.WriteString(p.Type.String())
	if len(p.Body.Objects) == 0 {
		sb.WriteString("{")
		for _, obj := range p.Body.Objects {
			sb.WriteString(obj.String())
			sb.WriteString("\n\n")
		}
		sb.WriteString("}\n")
	} else {
		sb.WriteString("{}\n")
	}
	return sb.String()
}

func (p *FuncType) String() string {
	var sb strings.Builder
	if len(p.Args) > 0 {
		sb.WriteString("(")
		for i, arg := range p.Args {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(arg.String())
		}
		sb.WriteString(")")
	}
	if p.Return != token.NONE {
		sb.WriteString(" => ")
		sb.WriteString(p.Return.String())
	}
	return sb.String()
}

func (p *FuncBody) String() string {
	var sb strings.Builder
	for _, obj := range p.Objects {
		sb.WriteString(obj.String())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (p *Argument) String() string {
	var sb strings.Builder
	sb.WriteString(p.Name)
	sb.WriteString(":")
	sb.WriteString(p.Type.String())
	return sb.String()
}

func (p *Local) String() string {
	var sb strings.Builder
	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteString("\n")
	}
	sb.WriteString("\tlocal ")
	sb.WriteString(p.Name)
	sb.WriteString(":")
	sb.WriteString(p.Type.String())
	if p.Comment != nil {
		sb.WriteString(" # ")
		sb.WriteString(p.Comment.String())
	}
	sb.WriteString("\n")
	return sb.String()
}

func (p *Instruction) String() string {
	var sb strings.Builder
	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteString("\n")
	}
	if p.Label != "" {
		sb.WriteString(p.Label)
		sb.WriteString(":\n")
	}
	if p.As != 0 {
		// pc 是否可以省略?
		sb.WriteString(riscv.AsmSyntax(0, p.As, p.Arg))
	}
	if p.Comment != nil {
		sb.WriteString(" # ")
		sb.WriteString(p.Comment.String())
	}
	sb.WriteString("\n")
	return sb.String()
}
