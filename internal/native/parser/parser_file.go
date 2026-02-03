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

		case token.GAS_X64_INTEL_SYNTAX:
			// 跳过intel语法开关
			if p.cpu != abi.X64Unix && p.cpu != abi.X64Windows {
				p.errorf(p.pos, "%v only enabled on X64 CPU", p.tok)
			}
			p.acceptToken(token.GAS_X64_INTEL_SYNTAX)
			p.acceptToken(token.GAS_X64_NOPREFIX)

		case token.GAS_EXTERN:
			p.parseFile_gasExtern()
		case token.GAS_SECTION:
			p.parseFile_gasSection()

		case token.CONST_zh:
			constObj := p.parseConst(p.tok)
			p.prog.Consts = append(p.prog.Consts, constObj)
			p.prog.Objects = append(p.prog.Objects, constObj)
		case token.GLOBAL_zh:
			globalObj := p.parseGlobal(p.tok)
			p.prog.Globals = append(p.prog.Globals, globalObj)
			p.prog.Objects = append(p.prog.Objects, globalObj)
		case token.FUNC_zh:
			funcObj := p.parseFunc(p.tok)
			p.prog.Funcs = append(p.prog.Funcs, funcObj)
			p.prog.Objects = append(p.prog.Objects, funcObj)

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
