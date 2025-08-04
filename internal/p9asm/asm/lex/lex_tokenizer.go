// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"bytes"
	"text/scanner"
	"unicode"

	"wa-lang.org/wa/internal/p9asm/obj"
	"wa-lang.org/wa/internal/p9asm/objabi"
)

var _ TokenReader = (*_Tokenizer)(nil)

// 基于 text/scanner.Scanner 包装的词法分析器
type _Tokenizer struct {
	ctxt *obj.Link
	tok  ScanToken
	s    *scanner.Scanner
	file *objabi.File
}

func newTokenizer(ctxt *obj.Link, fileName string, text []byte) *_Tokenizer {
	if ctxt == nil {
		ctxt = new(obj.Link)
	}
	if ctxt.Fset == nil {
		ctxt.Fset = objabi.NewFileSet()
	}

	f := ctxt.Fset.AddFile(fileName, -1, len(text))
	f.SetLinesForContent(text)

	var s scanner.Scanner
	s.Init(bytes.NewReader(text))

	// Newline is like a semicolon; other space characters are fine.
	s.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '

	// Don't skip comments: we need to count newlines.
	s.Mode = scanner.ScanChars |
		scanner.ScanFloats |
		scanner.ScanIdents |
		scanner.ScanInts |
		scanner.ScanStrings |
		scanner.ScanComments

	s.Position.Filename = fileName
	s.IsIdentRune = isIdentRune

	return &_Tokenizer{
		ctxt: ctxt,
		s:    &s,
		file: f,
	}
}

// We want center dot (·) and division slash (∕) to work as identifier characters.
func isIdentRune(ch rune, i int) bool {
	if unicode.IsLetter(ch) {
		return true
	}
	switch ch {
	case '_': // Underscore; traditional.
		return true
	case '\u00B7': // Represents the period in runtime.exit. U+00B7 '·' middle dot
		return true
	case '\u2215': // Represents the slash in runtime/debug.setGCPercent. U+2215 '∕' division slash
		return true
	}
	// Digits are OK only after the first character.
	return i > 0 && unicode.IsDigit(ch)
}

func (t *_Tokenizer) Text() string {
	switch t.tok {
	case LSH:
		return "<<"
	case RSH:
		return ">>"
	case ARR:
		return "->"
	case ROT:
		return "@>"
	}
	return t.s.TokenText()
}

func (in *_Tokenizer) Pos() objabi.Pos {
	return objabi.Pos(in.file.Base() + in.s.Pos().Offset)
}

func (t *_Tokenizer) Next() ScanToken {
	s := t.s
	for {
		t.tok = ScanToken(s.Scan())
		if t.tok != scanner.Comment {
			break
		}
	}
	switch t.tok {
	case '-':
		if s.Peek() == '>' {
			s.Next()
			t.tok = ARR
			return ARR
		}
	case '@':
		if s.Peek() == '>' {
			s.Next()
			t.tok = ROT
			return ROT
		}
	case '<':
		if s.Peek() == '<' {
			s.Next()
			t.tok = LSH
			return LSH
		}
	case '>':
		if s.Peek() == '>' {
			s.Next()
			t.tok = RSH
			return RSH
		}
	}
	return t.tok
}
