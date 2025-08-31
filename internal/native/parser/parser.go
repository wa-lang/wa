// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import "wa-lang.org/wa/internal/native/ast"

type parser struct {
	filename string
	src      []byte

	trace bool
}

func newParser(path string, src []byte) *parser {
	p := &parser{filename: path, src: src}
	return p
}

func (p *parser) ParseFile() (*ast.File, error) {
	panic("TODO")
}
