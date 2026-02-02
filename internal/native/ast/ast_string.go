// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package ast

import (
	"fmt"
	"strings"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/riscv"
)

func (p *File) String() string {
	var sb strings.Builder

	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteRune('\n')
	}

	if len(p.Objects) > 0 {
		// 优先以原始的顺序输出
		var prevObj Object
		for _, obj := range p.Objects {
			if obj.GetDoc() != nil || !isSameType(obj, prevObj) {
				sb.WriteString("\n")
			}
			sb.WriteString(obj.String())
			sb.WriteString("\n")
			prevObj = obj
		}
	} else {
		// 孤立的注释输出位置将失去上下文相关性
		for _, obj := range p.Comments {
			sb.WriteString(obj.String())
			sb.WriteString("\n\n")
		}

		for _, obj := range p.Consts {
			sb.WriteString(obj.String())
			sb.WriteString("\n\n")
		}
		for _, obj := range p.Globals {
			sb.WriteString(obj.String())
			sb.WriteString("\n\n")
		}
		for _, obj := range p.Funcs {
			sb.WriteString(obj.String())
			sb.WriteString("\n\n")
		}
	}

	return sb.String()
}

func (p *Comment) String() string {
	return p.Text
}

func (p *CommentGroup) String() string {
	var sb strings.Builder
	for i, c := range p.List {
		if i > 0 {
			sb.WriteRune('\n')
		}
		if !p.TopLevel {
			sb.WriteByte('\t')
		}
		sb.WriteString(c.String())
	}
	return sb.String()
}

func (p *BasicLit) String() string {
	return p.LitString
}

func (p *Const) String() string {
	var sb strings.Builder
	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteRune('\n')
	}
	sb.WriteString(fmt.Sprintf("%v %s = %v", p.Tok, p.Name, p.Value))
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

	if p.Type != 0 {
		sb.WriteString(fmt.Sprintf("%v %s:%v = ", p.Tok, p.Name, p.Type))
	} else if p.Size != 0 {
		sb.WriteString(fmt.Sprintf("%v %s:%d = ", p.Tok, p.Name, p.Size))
	} else {
		sb.WriteString(fmt.Sprintf("%v %s = ", p.Tok, p.Name))
	}

	if p.Init != nil {
		sb.WriteString(p.Init.Lit.String())
	} else {
		sb.WriteString(p.Init.Symbal)
	}

	return sb.String()
}

func (p *InitValue) String() string {
	var sb strings.Builder
	if p.Lit != nil {
		sb.WriteString(p.Lit.String())
	} else {
		sb.WriteString(p.Symbal)
	}
	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (p *Func) String() string {
	var sb strings.Builder

	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteRune('\n')
	}
	sb.WriteString(p.Tok.String())
	sb.WriteString(" ")
	sb.WriteString(p.Name)
	sb.WriteString(p.Type.String())
	sb.WriteString(" ")
	sb.WriteString("{\n")

	if len(p.Body.Objects) > 0 {
		var prevObj Object
		for i, obj := range p.Body.Objects {
			insertBlankLine := false
			if i > 0 {
				// 当前语句带文档, 前一个不是 Label, 尽量前面保持分开
				if obj.GetDoc() != nil || !isSameType(obj, prevObj) {
					if inst, ok := prevObj.(*Instruction); ok && inst.Label == "" {
						insertBlankLine = true
					}
				}

				// 当前语句是 Label
				if inst, ok := obj.(*Instruction); ok && inst.Label != "" {
					insertBlankLine = true
				}
			}
			if insertBlankLine {
				sb.WriteString("\n")
			}
			sb.WriteString(obj.String())
			sb.WriteString("\n")
			prevObj = obj
		}
	} else {
		// 孤立的注释输出位置将失去上下文相关性
		for _, obj := range p.Body.Comments {
			sb.WriteString(obj.String())
			sb.WriteString("\n\n")
		}

		for _, obj := range p.Body.Locals {
			sb.WriteString(obj.String())
			sb.WriteString("\n\n")
		}

		for _, obj := range p.Body.Insts {
			sb.WriteString(obj.String())
			sb.WriteString("\n\n")
		}
	}
	sb.WriteString("}\n")

	return sb.String()
}

func (p *FuncType) String() string {
	var sb strings.Builder
	if len(p.Args) > 0 {
		sb.WriteString("(")
		for i, arg := range p.Args {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(arg.Name)
			sb.WriteString(":")
			sb.WriteString(arg.Type.String())
		}
		sb.WriteString(")")
	}
	if len(p.Return) > 0 {
		sb.WriteString(" => ")
		if len(p.Return) == 1 && p.Return[0].Name == "" {
			sb.WriteString(p.Return[0].Type.String())
		} else {
			for i, ret := range p.Return {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(ret.Name)
				sb.WriteString(":")
				sb.WriteString(ret.Type.String())
			}
			sb.WriteString(")")
		}
	}
	return sb.String()
}

func (p *Local) String() string {
	var sb strings.Builder
	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteString("\n")
	}
	sb.WriteString("\t")
	sb.WriteString(p.Tok.String())
	sb.WriteString(" ")
	sb.WriteString(p.Name)
	sb.WriteString(":")
	sb.WriteString(p.Type.String())
	if p.Comment != nil {
		sb.WriteString(" ")
		sb.WriteString(p.Comment.String())
	}
	sb.WriteString("\n")
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

func (p *Instruction) String() string {
	var sb strings.Builder
	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteString("\n")
	}
	if p.Label != "" {
		sb.WriteString(p.Label)
		sb.WriteString(":")
	}
	if p.As != 0 {
		if p.Label != "" {
			sb.WriteString("\n")
		}
		sb.WriteString("\t")
		switch p.CPU {
		case abi.LOONG64:
			sb.WriteString(loong64.AsmSyntax(p.As, p.AsName, p.Arg))
		case abi.RISCV32, abi.RISCV64:
			sb.WriteString(riscv.AsmSyntax(p.As, p.AsName, p.Arg))
		default:
			panic("unreachable")
		}
	}
	if p.Comment != nil {
		sb.WriteString(" ")
		sb.WriteString(p.Comment.String())
	}
	return sb.String()
}
