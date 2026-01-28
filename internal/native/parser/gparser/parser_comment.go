// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package gparser

import (
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// 解析多个相邻的注释
func (p *parser) parseCommentGroup(isTopLevel bool) *ast.CommentGroup {
	comment := &ast.CommentGroup{TopLevel: isTopLevel}

	for p.tok == token.COMMENT {
		// 如果注释出现空行, 则结束
		if len(comment.List) > 0 {
			if p.posLine(p.pos) > p.getCommentGroupEndLine(comment)+1 {
				break
			}
		}
		comment.List = append(comment.List, &ast.Comment{
			Pos:  p.pos,
			Text: p.lit,
		})
		p.next()
	}

	// 跳过分号列表
	p.consumeSemicolonList()

	return comment
}

// 解析关联的注释
func (p *parser) parseDocComment(comments *[]*ast.CommentGroup, pos token.Pos) *ast.CommentGroup {
	if lastLine := p.getLastCommentGroupEndLine(*comments); lastLine > 0 && lastLine+1 == p.posLine(pos) {
		doc := (*comments)[len(*comments)-1]
		*comments = (*comments)[:len(*comments)-1]
		return doc
	}
	return nil
}

// 解析同一行的尾部注释
func (p *parser) parseTailComment(pos token.Pos) *ast.Comment {
	p.consumeSemicolonList()
	if p.tok == token.COMMENT && p.posLine(pos) == p.posLine(p.pos) {
		comment := &ast.Comment{Pos: p.pos, Text: p.lit}
		p.acceptToken(token.COMMENT)
		return comment
	}
	return nil
}

// 注释组结束的行号
func (p *parser) getLastCommentGroupEndLine(comments []*ast.CommentGroup) int {
	if n := len(comments); n > 0 {
		return p.getCommentGroupEndLine(comments[n-1])
	}
	return 0
}

func (p *parser) getCommentGroupEndLine(x *ast.CommentGroup) int {
	if len(x.List) > 0 {
		return p.posLine(x.List[len(x.List)-1].Pos)
	}
	return 0
}
