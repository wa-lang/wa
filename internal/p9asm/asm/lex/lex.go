// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package lex

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/p9asm/asm/arch"
	"wa-lang.org/wa/internal/p9asm/obj"
)

// Token读取接口
type TokenReader interface {
	Next() ScanToken

	Text() string
	File() string
	Line() int
	Col() int

	SetPos(line int, file string)

	Close()
}

// TODO: 基于 vfs 读文件
func NewTokenReader(name string, ctxt *obj.Link, flags *arch.Flags) (TokenReader, error) {
	linkCtxt = ctxt
	input, err := newInput(name, flags)
	if err != nil {
		return nil, err
	}
	fd, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("asm: %s", err)
	}
	input.Push(newTokenizer(name, fd, fd))
	return input, nil
}

var (
	// TODO(chai2010): 移除全局变量
	// It might be nice if these weren't global.
	linkCtxt *obj.Link     // The link context for all instructions.
	histLine int       = 1 // The cumulative count of lines processed.
)

// TODO(chai2010): 删除
// HistLine reports the cumulative source line number of the token,
// for use in the Prog structure for the linker. (It's always handling the
// instruction from the current lex line.)
// It returns int32 because that's what type ../asm prefers.
func HistLine() int32 {
	return int32(histLine)
}
