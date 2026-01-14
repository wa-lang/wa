// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"io"
)

const (
	kFuncMain = "main"
)

// 启动函数
func (p *wat2X64Worker) buildStart(w io.Writer) error {
	p.gasComment(w, "汇编程序入口函数")
	p.gasSectionTextStart(w)
	p.gasGlobal(w, kFuncMain)
	fmt.Fprintf(w, "%s:\n", kFuncMain)

	fmt.Fprintln(w, "    push rbp")
	fmt.Fprintln(w, "    mov  rbp, rsp")
	fmt.Fprintln(w, "    sub  rsp, 32")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "    call %s\n", kMemoryInitFuncName)

	if p.m.Start != "" {
		fmt.Fprintf(w, "    call %s\n", kFuncNamePrefix+p.m.Start)
	}
	if p.m.Start != kFuncMain {
		for _, fn := range p.m.Funcs {
			if fn.Name == kFuncMain {
				fmt.Fprintf(w, "    call %s\n", kFuncNamePrefix+kFuncMain)
			}
		}
	}
	fmt.Fprintln(w)

	p.gasCommentInFunc(w, "runtime.exit(0)")
	fmt.Fprintf(w, "    mov  rcx, 0\n")
	fmt.Fprintf(w, "    call %s\n", kRuntimeExit)
	fmt.Fprintln(w)

	p.gasCommentInFunc(w, "exit 后这里不会被执行, 但是依然保留")
	fmt.Fprintln(w, "    mov rsp, rbp")
	fmt.Fprintln(w, "    pop rbp")
	fmt.Fprintln(w, "    ret")
	fmt.Fprintln(w)

	return nil
}
