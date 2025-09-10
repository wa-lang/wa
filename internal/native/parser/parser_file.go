// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import "wa-lang.org/wa/internal/native/token"

func (p *parser) parseFile() {
	// 解析开头的关联文档
	p.prog.Doc = p.parseCommentGroup()

	// 解析代码主体
	for {
		if p.err != nil {
			return
		}
		if p.tok == token.EOF {
			return
		}

		switch p.tok {
		case token.COMMENT:
			p.prog.Comments = append(p.prog.Comments, p.parseCommentGroup())
		case token.CONST, token.CONST_zh:
			p.prog.Consts = append(p.prog.Consts, p.parseConst())
		case token.GLOBAL, token.GLOBAL_zh:
			p.prog.Globals = append(p.prog.Globals, p.parseGlobal())
		case token.FUNC, token.FUNC_zh:
			p.prog.Funcs = append(p.prog.Funcs, p.parseFunc())
		default:
			p.errorf(p.pos, "unkonw token: %v", p.tok)
		}
	}
}
