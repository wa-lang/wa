// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/token"
)

func (p *parser) parseFile() {
	// 解析开头的关联文档
	p.prog.Doc = p.parseCommentGroup(true)

	// 特别处理intel语法开关
	x64IntelSyntax := false
	for {
		if p.err != nil {
			return
		}
		if p.tok == token.EOF {
			return
		}

		if p.tok == token.COMMENT {
			commentObj := p.parseCommentGroup(true)
			p.prog.Comments = append(p.prog.Comments, commentObj)
			p.prog.Objects = append(p.prog.Objects, commentObj)
			continue
		}
		if p.tok == token.GAS_X64_INTEL_SYNTAX {
			if p.cpu != abi.X64Unix && p.cpu != abi.X64Windows {
				p.errorf(p.pos, "%v only enabled on X64 CPU", p.tok)
			}
			p.acceptToken(token.GAS_X64_NOPREFIX)
			x64IntelSyntax = true
		}

		// 解析后续代码
		break
	}

	// x64 必须打开 intel 语法开关
	if p.cpu == abi.X64Unix || p.cpu == abi.X64Windows {
		if !x64IntelSyntax {
			p.errorf(p.pos, "%v missing", token.GAS_X64_INTEL_SYNTAX)
		}
	}

	// 解析代码主体
	for {
		if p.err != nil {
			return
		}
		if p.tok == token.EOF {
			break
		}

		switch p.tok {
		case token.COMMENT:
			commentObj := p.parseCommentGroup(true)
			p.prog.Comments = append(p.prog.Comments, commentObj)
			p.prog.Objects = append(p.prog.Objects, commentObj)

		case token.GAS_EXTERN:
			p.parseFile_gasExtern()
		case token.GAS_SECTION:
			p.parseFile_gasSection()

		default:
			p.errorf(p.pos, "unkonw token: %v", p.tok)
		}
	}

	// 收集信息导出符号信息
	for _, g := range p.prog.Globals {
		if p.gasGlobal[g.Name] {
			g.Exported = true
		}
	}
	for _, fn := range p.prog.Funcs {
		if p.gasGlobal[fn.Name] {
			fn.ExportName = fn.Name
		}
	}
}
