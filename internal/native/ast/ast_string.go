// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package ast

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/loong64"
	"wa-lang.org/wa/internal/native/riscv"
	"wa-lang.org/wa/internal/native/token"
	"wa-lang.org/wa/internal/printer/tabwriter"
)

func (p *File) String() string {
	var sb strings.Builder

	if p.Doc != nil {
		sb.WriteString(p.Doc.String())
		sb.WriteString("\n")
	}

	if p.IntelSyntax != nil {
		sb.WriteString(p.IntelSyntax.String())
		sb.WriteString("\n")
	}

	// 优先以原始的顺序输出
	var prevObj Object
	for i, obj := range p.Objects {
		if obj.GetDoc() != nil || !isSameType(obj, prevObj) {
			sb.WriteString("\n")
		}

		// 全局变量: 如果是相同的段和对齐只需要打印一次
		if g, _ := obj.(*Global); g != nil {
			if g.Tok != token.GLOBAL_zh {
				gPrev, _ := prevObj.(*Global)
				if gPrev == nil || gPrev.Tok != g.Tok || gPrev.Section != g.Section || gPrev.Align != g.Align {
					sb.WriteString(token.GAS_SECTION.String())
					sb.WriteString(" ")
					sb.WriteString(g.Section)
					sb.WriteString("\n")

					sb.WriteString(token.GAS_ALIGN.String())
					sb.WriteString(" ")
					sb.WriteString(strconv.Itoa(g.Align))
					sb.WriteString("\n")
				}
				if g.ExportName != "" {
					sb.WriteString(token.GAS_GLOBL.String())
					sb.WriteByte(' ')
					sb.WriteString(g.ExportName)
					sb.WriteByte('\n')
				}
			}
		}

		sb.WriteString(obj.String())
		if i < len(p.Objects)-1 {
			sb.WriteString("\n")
		}
		prevObj = obj
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
			sb.WriteString("    ")
		}
		sb.WriteString(c.String())
	}
	return sb.String()
}

func (p *GasIntelSyntaxNoprefix) String() string {
	return fmt.Sprintf("%s %s", token.GAS_X64_INTEL_SYNTAX, token.GAS_X64_NOPREFIX)
}

func (p *Extern) String() string {
	return fmt.Sprintf("%s %s", p.Tok, p.Name)
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
	if p.Tok == token.CONST_zh {
		sb.WriteString(fmt.Sprintf("%v %s = %v", p.Tok, p.Name, p.Value))
	} else {
		sb.WriteString(fmt.Sprintf("%s = %v", p.Name, p.Value))
	}
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

	if p.Tok == token.GLOBAL_zh {
		sb.WriteString(fmt.Sprintf("%v %s: %v = ", p.Tok, p.Name, p.TypeTok))
	} else {
		sb.WriteString(fmt.Sprintf("%s: %v ", p.Name, p.TypeTok))
	}

	if p.Init.Symbal != "" {
		sb.WriteString(p.Init.Symbal)
	} else {
		sb.WriteString(p.Init.Lit.String())
	}

	return sb.String()
}

func (p *InitValue) String() string {
	var sb strings.Builder
	if p.Symbal != "" {
		sb.WriteString(p.Symbal)
	} else if p.Lit != nil {
		sb.WriteString(p.Lit.String())
	} else {
		sb.WriteString("?")
	}
	if p.Comment != nil {
		sb.WriteRune(' ')
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
	if p.Tok == token.FUNC_zh {
		sb.WriteString(p.Tok.String())
		sb.WriteString(" ")
		sb.WriteString(p.Name)
		sb.WriteString(":\n")
	} else {
		if p.Section != "" {
			sb.WriteString(token.GAS_SECTION.String())
			sb.WriteString(" ")
			sb.WriteString(p.Section)
			sb.WriteString("\n")
		}
		if p.ExportName != "" {
			sb.WriteString(token.GAS_GLOBL.String())
			sb.WriteString(" ")
			sb.WriteString(p.ExportName)
			sb.WriteString("\n")
		}
		sb.WriteString(p.Name)
		sb.WriteString(":\n")
	}

	// 打印指令
	// 指令参数和尾部注释对齐
	{
		// 用于格式化指令对齐
		var buf bytes.Buffer
		var w = tabwriter.NewWriter(&buf, 1, 1, 1, ' ', 0)

		for _, obj := range p.Body.Objects {
			if _, ok := obj.(*BlankLine); ok {
				fmt.Fprintln(w)
				w.Flush()
				continue
			}

			// 注释添加缩进
			if comment, ok := obj.(*CommentGroup); ok {
				fmt.Fprintf(w, "    %s\n", comment)
				continue
			}

			instObj, _ := obj.(*Instruction)
			if instObj == nil {
				fmt.Fprintln(w, obj.String())
				continue
			}

			// 指令类型
			if p.Tok == token.FUNC_zh {
				fmt.Fprintln(w, instObj.tabZhString())
			} else {
				fmt.Fprintln(w, instObj.tabString())
			}
		}
		w.Flush()

		// 输出到字符串目标
		sb.Write(buf.Bytes())
	}

	if p.Tok == token.FUNC_zh {
		sb.WriteString(token.END_zh.String())
		sb.WriteString("\n")
	}

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
	sb.WriteString("    ")
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
		sb.WriteString("    ")
		switch p.CPU {
		case abi.LOONG64:
			sb.WriteString(loong64.AsmSyntaxEx(p.As, p.AsName, p.Arg,
				loong64.RegAliasString,
				loong64.AsString,
			))
		case abi.RISCV32, abi.RISCV64:
			sb.WriteString(riscv.AsmSyntaxEx(p.As, p.AsName, p.Arg,
				riscv.RegAliasString,
				riscv.AsString,
			))
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

func (p *Instruction) ZhString() string {
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
		sb.WriteString("    ")
		switch p.CPU {
		case abi.LOONG64:
			sb.WriteString(loong64.AsmSyntaxEx(p.As, p.AsName, p.Arg,
				loong64.ZhRegAliasString,
				loong64.AsString,
			))
		case abi.RISCV32, abi.RISCV64:
			sb.WriteString(riscv.AsmSyntaxEx(p.As, p.AsName, p.Arg,
				riscv.ZhRegAliasString,
				riscv.AsString,
			))
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

func (p *Instruction) tabString() string {
	var sb strings.Builder
	if p.Doc != nil {
		sb.WriteString("    ")
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
		sb.WriteString("    ")
		switch p.CPU {
		case abi.LOONG64:
			s := loong64.AsmSyntaxEx(p.As, p.AsName, p.Arg,
				loong64.RegAliasString,
				loong64.AsString,
			)
			if idx := strings.IndexByte(s, ' '); idx > 0 {
				sb.WriteString(s[:idx])
				sb.WriteByte('\t')
				sb.WriteString(s[idx+1:])
			} else {
				sb.WriteString(s)
			}
		case abi.RISCV32, abi.RISCV64:
			s := riscv.AsmSyntaxEx(p.As, p.AsName, p.Arg,
				riscv.RegAliasString,
				riscv.AsString,
			)
			if idx := strings.IndexByte(s, ' '); idx > 0 {
				sb.WriteString(s[:idx])
				sb.WriteByte('\t')
				sb.WriteString(s[idx+1:])
			} else {
				sb.WriteString(s)
			}
		default:
			panic("unreachable")
		}
	}
	if p.Comment != nil {
		sb.WriteString("\t")
		sb.WriteString(p.Comment.String())
	}
	return sb.String()
}

func (p *Instruction) tabZhString() string {
	var sb strings.Builder
	if p.Doc != nil {
		sb.WriteString("    ")
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
		sb.WriteString("    ")
		switch p.CPU {
		case abi.LOONG64:
			s := loong64.AsmSyntaxEx(p.As, p.AsName, p.Arg,
				loong64.ZhRegAliasString,
				loong64.AsString,
			)
			if idx := strings.IndexByte(s, ' '); idx > 0 {
				sb.WriteString(s[:idx])
				sb.WriteByte('\t')
				sb.WriteString(s[idx+1:])
			} else {
				sb.WriteString(s)
			}
		case abi.RISCV32, abi.RISCV64:
			s := riscv.AsmSyntaxEx(p.As, p.AsName, p.Arg,
				riscv.ZhRegAliasString,
				riscv.AsString,
			)
			if idx := strings.IndexByte(s, ' '); idx > 0 {
				sb.WriteString(s[:idx])
				sb.WriteByte('\t')
				sb.WriteString(s[idx+1:])
			} else {
				sb.WriteString(s)
			}
		default:
			panic("unreachable")
		}
	}
	if p.Comment != nil {
		sb.WriteString("\t")
		sb.WriteString(p.Comment.String())
	}
	return sb.String()
}

func (p *BlankLine) String() string {
	return ""
}
