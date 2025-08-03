// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import "text/scanner"

var _ TokenReader = (*_Stack)(nil)

// TokenReader 栈, 用于处理 #include 宏指令
type _Stack struct {
	tr []TokenReader
}

// 主要用于宏展开时压栈
func (s *_Stack) Push(tr TokenReader) {
	s.tr = append(s.tr, tr)
}

// 下一个 Token,
// 当前 TokenReader 结束时自动弹出
func (s *_Stack) Next() ScanToken {
	tos := s.tr[len(s.tr)-1]
	tok := tos.Next()
	for tok == scanner.EOF && len(s.tr) > 1 {
		s.tr = s.tr[:len(s.tr)-1]
		tok = s.Next()
	}
	return tok
}

func (s *_Stack) Text() string {
	return s.tr[len(s.tr)-1].Text()
}

func (s *_Stack) File() string {
	return s.tr[len(s.tr)-1].File()
}

func (s *_Stack) Line() int {
	return s.tr[len(s.tr)-1].Line()
}

func (s *_Stack) Col() int {
	return s.tr[len(s.tr)-1].Col()
}

func (s *_Stack) SetPos(line int, file string) {
	s.tr[len(s.tr)-1].SetPos(line, file)
}

func (s *_Stack) Close() {
}
