// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

func (p *parser) parseFunc() *ast.Func {
	fn := &ast.Func{
		Type: new(ast.FuncType),
		Body: new(ast.FuncBody),
	}

	fn.Pos = p.pos
	fn.Doc = p.parseDocComment(&p.prog.Comments, fn.Pos)
	if fn.Doc != nil {
		p.prog.Objects = p.prog.Objects[:len(p.prog.Objects)-1]
	}

	fn.Tok = p.acceptToken(token.FUNC_zh)
	fn.Name = p.parseIdent()

	p.parseFunc_body(fn)
	p.consumeSemicolonList()

	return fn
}

func (p *parser) parseFunc_body(fn *ast.Func) {
	fn.Body.Pos = p.pos

	p.acceptToken(token.COLON)
	defer p.acceptToken(token.END_zh)

	for {
		if p.err != nil {
			break
		}
		if p.tok == token.EOF {
			break
		}
		if p.tok == token.END_zh {
			break
		}

		// 注释
		if p.tok == token.COMMENT {
			// 插入空行
			if len(fn.Body.Objects) > 0 {
				prevObj := fn.Body.Objects[len(fn.Body.Objects)-1]
				prevLine := p.posLine(prevObj.BeginPos())
				curLine := p.posLine(p.pos)
				if curLine-prevLine > 1 {
					fn.Body.Objects = append(fn.Body.Objects, &ast.BlankLine{
						Pos: p.pos - 1,
					})
				}
			}

			commentObj := p.parseCommentGroup(true)
			fn.Body.Comments = append(fn.Body.Comments, commentObj)
			fn.Body.Objects = append(fn.Body.Objects, commentObj)
			continue
		}

		// 解析指令
		if p.tok == token.IDENT || p.tok.IsAs() {
			// 插入空行
			if len(fn.Body.Objects) > 0 {
				prevObj := fn.Body.Objects[len(fn.Body.Objects)-1]
				prevLine := p.posLine(prevObj.BeginPos())
				curLine := p.posLine(p.pos)
				if curLine-prevLine > 1 {
					fn.Body.Objects = append(fn.Body.Objects, &ast.BlankLine{
						Pos: p.pos - 1,
					})
				}
			}

			// 解析指令的关联注释
			instDoc := p.parseDocComment(&fn.Body.Comments, p.pos)
			if instDoc != nil {
				fn.Body.Objects = fn.Body.Objects[:len(fn.Body.Objects)-1]
			}

			inst := p.parseInst(fn)
			inst.Doc = instDoc

			fn.Body.Insts = append(fn.Body.Insts, inst)
			fn.Body.Objects = append(fn.Body.Objects, inst)
			continue
		}

		// 未知 token
		break
	}

	// 删除结尾的空行
	for len(fn.Body.Objects) > 0 {
		lastInstObj := fn.Body.Objects[len(fn.Body.Objects)-1]
		if _, ok := lastInstObj.(*ast.BlankLine); ok {
			fn.Body.Objects = fn.Body.Objects[:len(fn.Body.Objects)-1]
		} else {
			break
		}
	}
}
