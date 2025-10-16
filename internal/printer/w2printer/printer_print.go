// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package w2printer

import (
	"fmt"

	"wa-lang.org/wa/internal/token"
)

// writeIndent writes indentation.
func (p *printer) writeIndent() {
	// use "hard" htabs - indentation columns
	// must not be discarded by the tabwriter
	n := p.indent // include base indentation
	for i := 0; i < n; i++ {
		p.output = append(p.output, '\t')
	}

	// update positions
	p.pos.Offset += n
	p.pos.Column += n
	p.out.Column += n
}

// writeByte writes ch n times to p.output and updates p.pos.
// Only used to write formatting (white space) characters.
func (p *printer) writeByte(ch byte, n int) {
	if p.endAlignment {
		// Ignore any alignment control character;
		// and at the end of the line, break with
		// a formfeed to indicate termination of
		// existing columns.
		switch ch {
		case '\t', '\v':
			ch = ' '
		case '\n', '\f':
			ch = '\f'
			p.endAlignment = false
		}
	}

	if p.out.Column == 1 {
		// no need to write line directives before white space
		p.writeIndent()
	}

	for i := 0; i < n; i++ {
		p.output = append(p.output, ch)
	}

	// update positions
	p.pos.Offset += n
	if ch == '\n' || ch == '\f' {
		p.pos.Line += n
		p.out.Line += n
		p.pos.Column = 1
		p.out.Column = 1
		return
	}
	p.pos.Column += n
	p.out.Column += n
}

// writeString writes the string s to p.output and updates p.pos, p.out,
// and p.last. If isLit is set, s is escaped w/ kEscape characters
// to protect s from being interpreted by the tabwriter.
//
// Note: writeString is only used to write Go tokens, literals, and
// comments, all of which must be written literally. Thus, it is correct
// to always set isLit = true. However, setting it explicitly only when
// needed (i.e., when we don't know that s contains no tabs or line breaks)
// avoids processing extra escape characters and reduces run time of the
// printer benchmark by up to 10%.
func (p *printer) writeString(pos token.Position, s string, isLit bool) {
	if p.out.Column == 1 {
		p.writeIndent()
	}

	if pos.IsValid() {
		// update p.pos (if pos is invalid, continue with existing p.pos)
		// Note: Must do this after handling line beginnings because
		// writeIndent updates p.pos if there's indentation, but p.pos
		// is the position of s.
		p.pos = pos
	}

	if isLit {
		// Protect s such that is passes through the tabwriter
		// unchanged. Note that valid Go programs cannot contain
		// kEscape bytes since they do not appear in legal
		// UTF-8 sequences.
		p.output = append(p.output, kEscape)
	}

	if debug {
		p.output = append(p.output, fmt.Sprintf("/*%s*/", pos)...) // do not update p.pos!
	}
	p.output = append(p.output, s...)

	// update positions
	nlines := 0
	var li int // index of last newline; valid if nlines > 0
	for i := 0; i < len(s); i++ {
		// Raw string literals may contain any character except back quote (`).
		if ch := s[i]; ch == '\n' || ch == '\f' {
			// account for line break
			nlines++
			li = i
			// A line break inside a literal will break whatever column
			// formatting is in place; ignore any further alignment through
			// the end of the line.
			p.endAlignment = true
		}
	}
	p.pos.Offset += len(s)
	if nlines > 0 {
		p.pos.Line += nlines
		p.out.Line += nlines
		c := len(s) - li
		p.pos.Column = c
		p.out.Column = c
	} else {
		p.pos.Column += len(s)
		p.out.Column += len(s)
	}

	if isLit {
		p.output = append(p.output, kEscape)
	}

	p.last = p.pos
}

// whiteWhitespace writes the first n whitespace entries.
func (p *printer) writeWhitespace(n int) {
	// write entries
	for i := 0; i < n; i++ {
		switch ch := p.wsbuf[i]; ch {
		case ignore:
			// ignore!
		case indent:
			p.indent++
		case unindent:
			p.indent--
			if p.indent < 0 {
				p.internalError("negative indentation:", p.indent)
				p.indent = 0
			}
		case newline, formfeed:
			// A line break immediately followed by a "correcting"
			// unindent is swapped with the unindent - this permits
			// proper label positioning. If a comment is between
			// the line break and the label, the unindent is not
			// part of the comment whitespace prefix and the comment
			// will be positioned correctly indented.
			if i+1 < n && p.wsbuf[i+1] == unindent {
				// Use a formfeed to terminate the current section.
				// Otherwise, a long label name on the next line leading
				// to a wide column may increase the indentation column
				// of lines before the label; effectively leading to wrong
				// indentation.
				p.wsbuf[i], p.wsbuf[i+1] = unindent, formfeed
				i-- // do it again
				continue
			}
			fallthrough
		default:
			p.writeByte(byte(ch), 1)
		}
	}

	// shift remaining entries down
	l := copy(p.wsbuf, p.wsbuf[n:])
	p.wsbuf = p.wsbuf[:l]
}
