// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 说明: 汇编语言的宏只保留 #define 定义

package lex

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/p9asm/asm/arch"
	"wa-lang.org/wa/internal/p9asm/obj"
	"wa-lang.org/wa/internal/p9asm/objabi"
)

// Token读取接口
type TokenReader interface {
	Next() ScanToken

	Text() string
	File() string
	Line() int
	Col() int

	Pos() objabi.Pos
	SetPos(line int, file string)

	Close()
}

// TODO: 基于 vfs 读文件
func NewTokenReader(ctxt *obj.Link, name string, flags *arch.Flags) (TokenReader, error) {
	input, err := newInput(ctxt, name, flags)
	if err != nil {
		return nil, err
	}
	fd, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("asm: %s", err)
	}
	input.Push(newTokenizer(input.ctxt, name, fd, fd))
	return input, nil
}
