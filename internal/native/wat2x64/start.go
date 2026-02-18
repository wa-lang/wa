// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/native/abi"
)

const (
	kFuncStart = "_start"
	kFuncMain  = "main"
	kFuncMain2 = "_main"
)

// 启动函数
func (p *wat2X64Worker) buildStart(w io.Writer) error {
	p.gasComment(w, "汇编程序入口函数")
	p.gasSectionTextStart(w)

	p.gasGlobal(w, kFuncStart)
	fmt.Fprintf(w, "%s:\n", kFuncStart)

	// 参数寄存器
	regArg0 := "rcx"

	if p.cpuType == abi.X64Unix {
		regArg0 = "rdi"
	}

	fmt.Fprintln(w, "    push rbp")
	fmt.Fprintln(w, "    mov  rbp, rsp")
	fmt.Fprintln(w, "    sub  rsp, 32")
	fmt.Fprintln(w)

	if p.m.Memory != nil {
		fmt.Fprintf(w, "    call %s\n", kMemoryInitFuncName)
	}
	if p.m.Table != nil {
		fmt.Fprintf(w, "    call %s\n", kTableInitFuncName)
	}

	if p.m.Start != "" {
		fmt.Fprintf(w, "    call %s\n", kFuncNamePrefix+p.m.Start)
	}
	if p.m.Start != kFuncMain {
		for _, fn := range p.m.Funcs {
			switch fn.Name {
			case kFuncMain, kFuncMain2:
				fmt.Fprintf(w, "    call %s\n", kFuncNamePrefix+fixName(fn.Name))
			}
		}
	}
	fmt.Fprintln(w)

	p.gasCommentInFunc(w, "runtime.exit(0)")
	fmt.Fprintf(w, "    mov  %s, 0\n", regArg0)
	fmt.Fprintf(w, "    call %s\n", kRuntimeExit)
	fmt.Fprintln(w)

	p.gasCommentInFunc(w, "exit 后这里不会被执行, 但是依然保留")
	fmt.Fprintln(w, "    mov rsp, rbp")
	fmt.Fprintln(w, "    pop rbp")
	fmt.Fprintln(w, "    ret")
	fmt.Fprintln(w)

	return nil
}
