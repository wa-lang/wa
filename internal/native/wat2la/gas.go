// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"
	"strings"
)

func (p *wat2laWorker) gasComment(w io.Writer, text string) {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "#") {
		fmt.Fprintln(w, text)
	} else {
		fmt.Fprintln(w, "#", text)
	}
}

func (p *wat2laWorker) gasCommentInFunc(w io.Writer, text string) {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "#") {
		fmt.Fprintln(w, "    "+text)
	} else {
		fmt.Fprintln(w, "    #", text)
	}
}

func (p *wat2laWorker) gasSectionDataStart(w io.Writer) {
	fmt.Fprintln(w, ".section .data")
	fmt.Fprintln(w, ".align 3") // 2^3 = 8 字节对齐
}

func (p *wat2laWorker) gasSectionTextStart(w io.Writer) {
	fmt.Fprintln(w, ".section .text")
}

func (p *wat2laWorker) gasExtern(w io.Writer, name string) {
	fmt.Fprintln(w, ".extern", name)
}

func (p *wat2laWorker) gasGlobal(w io.Writer, name string) {
	fmt.Fprintln(w, ".globl", name)
}

func (p *wat2laWorker) gasDefI32(w io.Writer, name string, v int32) {
	fmt.Fprintf(w, "%s: .long %d\n", name, v)
}

func (p *wat2laWorker) gasDefI64(w io.Writer, name string, v int64) {
	fmt.Fprintf(w, "%s: .quad %d\n", name, v)
}

func (p *wat2laWorker) gasDefF32(w io.Writer, name string, v float32) {
	fmt.Fprintf(w, "%s: .float %f\n", name, v)
}

func (p *wat2laWorker) gasDefF64(w io.Writer, name string, v float64) {
	fmt.Fprintf(w, "%s: .double %f\n", name, v)
}

func (p *wat2laWorker) gasDefString(w io.Writer, name string, v string) {
	fmt.Fprintf(w, "%s: .asciz \"", name)
	fmt.Fprint(w, v)
	fmt.Fprintln(w, "\"")
}

func (p *wat2laWorker) gasFuncStart(w io.Writer, fnName string) {
	fmt.Fprintf(w, "%s:\n", fnName)
}

func (p *wat2laWorker) gasFuncLabel(w io.Writer, labelName string) {
	fmt.Fprintf(w, "%s:\n", labelName)
}

// 函数名字重定向
func (p *wat2laWorker) gasSet(w io.Writer, src, dst string) {
	fmt.Fprintf(w, ".set %s, %s\n", src, dst)
}
