// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"fmt"
	"log"
	"os"
)

func (p *Parser) errorf(format string, args ...interface{}) {
	if p.panicOnError {
		panic(fmt.Errorf(format, args...))
	}
	if p.histLineNum == p.errorLine {
		// Only one error per line.
		return
	}
	p.errorLine = p.histLineNum
	// Put file and line information on head of message.
	format = "%s:%d: " + format + "\n"
	args = append([]interface{}{p.lex.File(), p.lineNum}, args...)
	fmt.Fprintf(os.Stderr, format, args...)
	p.errorCount++
	if p.errorCount > 10 {
		log.Fatal("too many errors")
	}
}
