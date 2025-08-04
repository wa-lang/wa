// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 说明: 汇编语言的宏只保留 #define 定义

package lex

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"

	"wa-lang.org/wa/internal/p9asm/asm/arch"
	"wa-lang.org/wa/internal/p9asm/obj"
	"wa-lang.org/wa/internal/p9asm/objabi"
)

// Token读取接口
type TokenReader interface {
	Next() ScanToken
	Text() string
	Pos() objabi.Pos
}

func NewTokenReader(ctxt *obj.Link, vfs fs.FS, filename string, src interface{}, flags *arch.Flags) (TokenReader, error) {
	input, err := newInput(ctxt, filename, flags)
	if err != nil {
		return nil, err
	}

	// get source
	text, err := readSource(vfs, filename, src)
	if err != nil {
		return nil, err
	}

	input.Push(newTokenizer(input.ctxt, filename, text))
	return input, nil
}

func readSource(vfs fs.FS, filename string, src interface{}) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			// is io.Reader, but src is already available in []byte form
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			return io.ReadAll(s)
		}
		return nil, errors.New("invalid source")
	}
	if vfs != nil {
		return fs.ReadFile(vfs, filename)
	}
	return os.ReadFile(filename)
}
