// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 不再支持 #define 宏, 作为 inline 的函数处理(可递归), 常量也是

package parser

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"

	"wa-lang.org/wa/internal/p9asm/ast"
	"wa-lang.org/wa/internal/p9asm/objabi"
)

// 解析汇编文件
func ParseFile(vfs fs.FS, fset *objabi.FileSet, filename string, src interface{}) (f *ast.File, err error) {
	if fset == nil {
		panic("parser.ParseFile: no token.FileSet provided (fset == nil)")
	}

	// get source
	text, err := readSource(vfs, filename, src)
	if err != nil {
		return nil, err
	}

	_ = text

	panic("TODO")
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
