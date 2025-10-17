// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"unicode/utf8"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

func (p *printer) fieldList(fields *ast.FieldList, isStruct, isIncomplete bool) {
	lbrace := fields.Opening
	list := fields.List
	rbrace := fields.Closing
	hasComments := isIncomplete || p.commentBefore(p.posFor(rbrace))
	srcIsOneLine := lbrace.IsValid() && rbrace.IsValid() && p.lineFor(lbrace) == p.lineFor(rbrace)

	if !hasComments && srcIsOneLine {
		// possibly a one-line struct/interface
		if len(list) == 0 {
			// no blank between keyword and {} in this case
			p.print(lbrace, token.COLON, rbrace, token.Zh_完毕)
			return
		} else if p.isOneLineFieldList(list) {
			// small enough - print on one line
			// (don't use identList and ignore source line breaks)
			p.print(lbrace, token.COLON, blank)
			f := list[0]
			if isStruct {
				for i, x := range f.Names {
					if i > 0 {
						// no comments so no need for comma position
						p.print(token.COMMA, blank)
					}
					p.expr(x)
				}
				// 匿名嵌入的成员不打印冒号
				if len(f.Names) > 0 {
					p.print(token.COLON)
					p.print(blank)
				}
				p.expr(f.Type)
			} else { // interface
				if ftyp, isFtyp := f.Type.(*ast.FuncType); isFtyp {
					// method
					p.expr(f.Names[0])
					p.signature(ftyp)
				} else {
					// embedded interface
					p.expr(f.Type)
				}
			}
			p.print(blank, rbrace, token.Zh_完毕)
			return
		}
	}
	// hasComments || !srcIsOneLine
	p.print(lbrace, token.COLON, indent)
	if hasComments || len(list) > 0 {
		p.print(formfeed)
	}

	if isStruct {

		sep := vtab
		if len(list) == 1 {
			sep = blank
		}
		var line int
		for i, f := range list {
			if i > 0 {
				p.linebreak(p.lineFor(f.Pos()), 1, ignore, p.linesFrom(line) > 0)
			}
			extraTabs := 0
			p.setComment(f.Doc)
			p.recordLine(&line)
			if len(f.Names) > 0 {
				// named fields
				p.identList(f.Names, false)
				p.print(token.COLON)
				p.print(sep)
				p.expr(f.Type)
				extraTabs = 1
			} else {
				// anonymous field
				p.expr(f.Type)
				extraTabs = 2
			}
			if f.Tag != nil {
				if len(f.Names) > 0 && sep == vtab {
					p.print(sep)
				}
				p.print(sep)
				p.expr(f.Tag)
				extraTabs = 0
			}
			if f.Comment != nil {
				for ; extraTabs > 0; extraTabs-- {
					p.print(sep)
				}
				p.setComment(f.Comment)
			}
		}
		if isIncomplete {
			if len(list) > 0 {
				p.print(formfeed)
			}
			p.flush(p.posFor(rbrace), token.Zh_完毕) // make sure we don't lose the last line comment
			p.setLineComment("// " + filteredMsg)
		}

	} else { // interface

		var line int
		for i, f := range list {
			if i > 0 {
				p.linebreak(p.lineFor(f.Pos()), 1, ignore, p.linesFrom(line) > 0)
			}
			p.setComment(f.Doc)
			p.recordLine(&line)
			if ftyp, isFtyp := f.Type.(*ast.FuncType); isFtyp {
				// method
				p.expr(f.Names[0])
				p.signature(ftyp)
			} else {
				// embedded interface
				p.expr(f.Type)
			}
			p.setComment(f.Comment)
		}
		if isIncomplete {
			if len(list) > 0 {
				p.print(formfeed)
			}
			p.flush(p.posFor(rbrace), token.Zh_完毕) // make sure we don't lose the last line comment
			p.setLineComment("// contains filtered or unexported methods")
		}

	}
	p.print(unindent, formfeed, rbrace, token.Zh_完毕)
}

func (p *printer) isOneLineFieldList(list []*ast.Field) bool {
	if len(list) != 1 {
		return false // allow only one field
	}
	f := list[0]
	if f.Tag != nil || f.Comment != nil {
		return false // don't allow tags or comments
	}
	// only name(s) and type
	const maxSize = 30 // adjust as appropriate, this is an approximate value
	namesSize := identListSize(f.Names, maxSize)
	if namesSize > 0 {
		namesSize = 1 // blank between names and types
	}
	typeSize := p.nodeSize(f.Type, maxSize)
	return namesSize+typeSize <= maxSize
}

func (p *printer) setLineComment(text string) {
	p.setComment(&ast.CommentGroup{List: []*ast.Comment{{Slash: token.NoPos, Text: text}}})
}

func identListSize(list []*ast.Ident, maxSize int) (size int) {
	for i, x := range list {
		if i > 0 {
			size += len(", ")
		}
		size += utf8.RuneCountInString(x.Name)
		if size >= maxSize {
			break
		}
	}
	return
}
