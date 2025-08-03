// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import "text/scanner"

var _ TokenReader = (*_Slice)(nil)

// A _Slice reads from a slice of Tokens.
type _Slice struct {
	tokens   []Token
	fileName string
	line     int
	pos      int
}

func newSlice(fileName string, line int, tokens []Token) *_Slice {
	return &_Slice{
		tokens:   tokens,
		fileName: fileName,
		line:     line,
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

func (s *_Slice) File() string {
	return s.fileName
}

func (s *_Slice) Line() int {
	return s.line
}

func (s *_Slice) Col() int {
	// Col is only called when defining a macro, which can't reach here.
	panic("cannot happen: slice col")
}

func (s *_Slice) SetPos(line int, file string) {
	// Cannot happen because we only have slices of already-scanned
	// text, but be prepared.
	s.line = line
	s.fileName = file
}

func (s *_Slice) Close() {
}
