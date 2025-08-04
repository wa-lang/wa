// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"text/scanner"

	"wa-lang.org/wa/internal/p9asm/objabi"
)

var _ TokenReader = (*_Slice)(nil)

// A _Slice reads from a slice of Tokens.
type _Slice struct {
	tokens   []Token
	fileName string
	pos      int
}

func newSlice(fileName string, line int, tokens []Token) *_Slice {
	return &_Slice{
		tokens:   tokens,
		fileName: fileName,
		pos:      -1, // Next will advance to zero.
	}
}

func (s *_Slice) Next() ScanToken {
	s.pos++
	if s.pos >= len(s.tokens) {
		return scanner.EOF
	}
	return s.tokens[s.pos].ScanToken
}

func (s *_Slice) Text() string {
	return s.tokens[s.pos].Text
}

func (s *_Slice) Pos() objabi.Pos {
	return s.tokens[s.pos].Pos
}
