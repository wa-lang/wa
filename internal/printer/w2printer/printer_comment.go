// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"strings"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

// commentsHaveNewline reports whether a list of comments belonging to
// an *ast.CommentGroup contains newlines. Because the position information
// may only be partially correct, we also have to read the comment text.
func (p *printer) commentsHaveNewline(list []*ast.Comment) bool {
	// len(list) > 0
	line := p.lineFor(list[0].Pos())
	for i, c := range list {
		if i > 0 && p.lineFor(list[i].Pos()) != line {
			// not all comments on the same line
			return true
		}
		if t := c.Text; len(t) >= 2 && (t[1] == '/' || strings.Contains(t, "\n")) {
			return true
		}
	}
	_ = line
	return false
}

func (p *printer) nextComment() {
	for p.cindex < len(p.comments) {
		c := p.comments[p.cindex]
		p.cindex++
		if list := c.List; len(list) > 0 {
			p.comment = c
			p.commentOffset = p.posFor(list[0].Pos()).Offset
			p.commentNewline = p.commentsHaveNewline(list)
			return
		}
		// we should not reach here (correct ASTs don't have empty
		// ast.CommentGroup nodes), but be conservative and try again
	}
	// no more comments
	p.commentOffset = infinity
}

// commentBefore reports whether the current comment group occurs
// before the next position in the source code and printing it does
// not introduce implicit semicolons.
func (p *printer) commentBefore(next token.Position) bool {
	return p.commentOffset < next.Offset && (!p.impliedSemi || !p.commentNewline)
}

// commentSizeBefore returns the estimated size of the
// comments on the same line before the next position.
func (p *printer) commentSizeBefore(next token.Position) int {
	// save/restore current p.commentInfo (p.nextComment() modifies it)
	defer func(info commentInfo) {
		p.commentInfo = info
	}(p.commentInfo)

	size := 0
	for p.commentBefore(next) {
		for _, c := range p.comment.List {
			size += len(c.Text)
		}
		p.nextComment()
	}
	return size
}

// writeCommentPrefix writes the whitespace before a comment.
// If there is any pending whitespace, it consumes as much of
// it as is likely to help position the comment nicely.
// pos is the comment position, next the position of the item
// after all pending comments, prev is the previous comment in
// a group of comments (or nil), and tok is the next token.
func (p *printer) writeCommentPrefix(pos, next token.Position, prev *ast.Comment, tok token.Token) {
	if len(p.output) == 0 {
		// the comment is the first item to be printed - don't write any whitespace
		return
	}

	if pos.IsValid() && pos.Filename != p.last.Filename {
		// comment in a different file - separate with newlines
		p.writeByte('\f', maxNewlines)
		return
	}

	if pos.Line == p.last.Line && (prev == nil || prev.Text[1] != '/') {
		// comment on the same line as last item:
		// separate with at least one separator
		hasSep := false
		if prev == nil {
			// first comment of a comment group
			j := 0
			for i, ch := range p.wsbuf {
				switch ch {
				case blank:
					// ignore any blanks before a comment
					p.wsbuf[i] = ignore
					continue
				case vtab:
					// respect existing tabs - important
					// for proper formatting of commented structs
					hasSep = true
					continue
				case indent:
					// apply pending indentation
					continue
				}
				j = i
				break
			}
			p.writeWhitespace(j)
		}
		// make sure there is at least one separator
		if !hasSep {
			sep := byte('\t')
			if pos.Line == next.Line {
				// next item is on the same line as the comment
				// (which must be a /*-style comment): separate
				// with a blank instead of a tab
				sep = ' '
			}
			p.writeByte(sep, 1)
		}

	} else {
		// comment on a different line:
		// separate with at least one line break
		droppedLinebreak := false
		j := 0
		for i, ch := range p.wsbuf {
			switch ch {
			case blank, vtab:
				// ignore any horizontal whitespace before line breaks
				p.wsbuf[i] = ignore
				continue
			case indent:
				// apply pending indentation
				continue
			case unindent:
				// if this is not the last unindent, apply it
				// as it is (likely) belonging to the last
				// construct (e.g., a multi-line expression list)
				// and is not part of closing a block
				if i+1 < len(p.wsbuf) && p.wsbuf[i+1] == unindent {
					continue
				}
				// if the next token is not a closing }, apply the unindent
				// if it appears that the comment is aligned with the
				// token; otherwise assume the unindent is part of a
				// closing block and stop (this scenario appears with
				// comments before a case label where the comments
				// apply to the next case instead of the current one)
				if tok != token.RBRACE && pos.Column == next.Column {
					continue
				}
			case newline, formfeed:
				p.wsbuf[i] = ignore
				droppedLinebreak = prev == nil // record only if first comment of a group
			}
			j = i
			break
		}
		p.writeWhitespace(j)

		// determine number of linebreaks before the comment
		n := 0
		if pos.IsValid() && p.last.IsValid() {
			n = pos.Line - p.last.Line
			if n < 0 { // should never happen
				n = 0
			}
		}

		// at the package scope level only (p.indent == 0),
		// add an extra newline if we dropped one before:
		// this preserves a blank line before documentation
		// comments at the package scope level (issue 2570)
		if p.indent == 0 && droppedLinebreak {
			n++
		}

		// make sure there is at least one line break
		// if the previous comment was a line comment
		if n == 0 && prev != nil && prev.Text[1] == '/' {
			n = 1
		}

		if n > 0 {
			// use formfeeds to break columns before a comment;
			// this is analogous to using formfeeds to separate
			// individual lines of /*-style comments
			p.writeByte('\f', nlimit(n))
		}
	}
}

func (p *printer) writeComment(comment *ast.Comment) {
	text := comment.Text
	pos := p.posFor(comment.Pos())

	const linePrefix = "//line "
	if strings.HasPrefix(text, linePrefix) && (!pos.IsValid() || pos.Column == 1) {
		// Possibly a //-style line directive.
		// Suspend indentation temporarily to keep line directive valid.
		defer func(indent int) { p.indent = indent }(p.indent)
		p.indent = 0
	}

	// shortcut common case of #-style comments
	if text[0] == '#' {
		p.writeString(pos, trimRight(text), true)
		return
	}
	// shortcut common case of //-style comments
	if text[0] == '/' && text[1] == '/' {
		p.writeString(pos, trimRight(text), true)
		return
	}

	// for /*-style comments, print line by line and var the
	// write function take care of the proper indentation
	lines := strings.Split(text, "\n")

	// The comment started in the first column but is going
	// to be indented. For an idempotent result, add indentation
	// to all lines such that they look like they were indented
	// before - this will make sure the common prefix computation
	// is the same independent of how many times formatting is
	// applied (was issue 1835).
	if pos.IsValid() && pos.Column == 1 && p.indent > 0 {
		for i, line := range lines[1:] {
			lines[1+i] = "   " + line
		}
	}

	stripCommonPrefix(lines)

	// write comment lines, separated by formfeed,
	// without a line break after the last line
	for i, line := range lines {
		if i > 0 {
			p.writeByte('\f', 1)
			pos = p.pos
		}
		if len(line) > 0 {
			p.writeString(pos, trimRight(line), true)
		}
	}
}

// writeCommentSuffix writes a line break after a comment if indicated
// and processes any leftover indentation information. If a line break
// is needed, the kind of break (newline vs formfeed) depends on the
// pending whitespace. The writeCommentSuffix result indicates if a
// newline was written or if a formfeed was dropped from the whitespace
// buffer.
func (p *printer) writeCommentSuffix(needsLinebreak bool) (wroteNewline, droppedFF bool) {
	for i, ch := range p.wsbuf {
		switch ch {
		case blank, vtab:
			// ignore trailing whitespace
			p.wsbuf[i] = ignore
		case indent, unindent:
			// don't lose indentation information
		case newline, formfeed:
			// if we need a line break, keep exactly one
			// but remember if we dropped any formfeeds
			if needsLinebreak {
				needsLinebreak = false
				wroteNewline = true
			} else {
				if ch == formfeed {
					droppedFF = true
				}
				p.wsbuf[i] = ignore
			}
		}
	}
	p.writeWhitespace(len(p.wsbuf))

	// make sure we have a line break
	if needsLinebreak {
		p.writeByte('\n', 1)
		wroteNewline = true
	}

	return
}

// intersperseComments consumes all comments that appear before the next token
// tok and prints it together with the buffered whitespace (i.e., the whitespace
// that needs to be written before the next token). A heuristic is used to mix
// the comments and whitespace. The intersperseComments result indicates if a
// newline was written or if a formfeed was dropped from the whitespace buffer.
func (p *printer) intersperseComments(next token.Position, tok token.Token) (wroteNewline, droppedFF bool) {
	var last *ast.Comment
	for p.commentBefore(next) {
		for _, c := range p.comment.List {
			p.writeCommentPrefix(p.posFor(c.Pos()), next, last, tok)
			p.writeComment(c)
			last = c
		}
		p.nextComment()
	}

	if last != nil {
		// If the last comment is a /*-style comment and the next item
		// follows on the same line but is not a comma, and not a "closing"
		// token immediately following its corresponding "opening" token,
		// add an extra separator unless explicitly disabled. Use a blank
		// as separator unless we have pending linebreaks, they are not
		// disabled, and we are outside a composite literal, in which case
		// we want a linebreak (issue 15137).
		// TODO(gri) This has become overly complicated. We should be able
		// to track whether we're inside an expression or statement and
		// use that information to decide more directly.
		needsLinebreak := false
		if p.mode&noExtraBlank == 0 &&
			last.Text[1] == '*' && p.lineFor(last.Pos()) == next.Line &&
			tok != token.COMMA &&
			(tok != token.RPAREN || p.prevOpen == token.LPAREN) &&
			(tok != token.RBRACK || p.prevOpen == token.LBRACK) {
			if p.containsLinebreak() && p.mode&noExtraLinebreak == 0 && p.level == 0 {
				needsLinebreak = true
			} else {
				p.writeByte(' ', 1)
			}
		}
		// Ensure that there is a line break after a //-style comment,
		// before EOF, and before a closing '}' unless explicitly disabled.
		if last.Text[1] == '/' ||
			tok == token.EOF ||
			tok == token.RBRACE && p.mode&noExtraLinebreak == 0 {
			needsLinebreak = true
		}
		return p.writeCommentSuffix(needsLinebreak)
	}

	// no comment was written - we should never reach here since
	// intersperseComments should not be called in that case
	p.internalError("intersperseComments called without pending comments")
	return
}

// containsLinebreak reports whether the whitespace buffer contains any line breaks.
func (p *printer) containsLinebreak() bool {
	for _, ch := range p.wsbuf {
		if ch == newline || ch == formfeed {
			return true
		}
	}
	return false
}

// setComment sets g as the next comment if g != nil and if node comments
// are enabled - this mode is used when printing source code fragments such
// as exports only. It assumes that there is no pending comment in p.comments
// and at most one pending comment in the p.comment cache.
func (p *printer) setComment(g *ast.CommentGroup) {
	if g == nil || !p.useNodeComments {
		return
	}
	if p.comments == nil {
		// initialize p.comments lazily
		p.comments = make([]*ast.CommentGroup, 1)
	} else if p.cindex < len(p.comments) {
		// for some reason there are pending comments; this
		// should never happen - handle gracefully and flush
		// all comments up to g, ignore anything after that
		p.flush(p.posFor(g.List[0].Pos()), token.ILLEGAL)
		p.comments = p.comments[0:1]
		// in debug mode, report error
		p.internalError("setComment found pending comments")
	}
	p.comments[0] = g
	p.cindex = 0
	// don't overwrite any pending comment in the p.comment cache
	// (there may be a pending comment when a line comment is
	// immediately followed by a lead comment with no other
	// tokens between)
	if p.commentOffset == infinity {
		p.nextComment() // get comment ready for use
	}
}
