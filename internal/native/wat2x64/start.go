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

	kAppArgcName = ".Wa.App.argc"
	kAppArgvName = ".Wa.App.argv"
	kAppEnvpName = ".Wa.App.envp"
)

// 生成命令行参数
func (p *wat2X64Worker) buildEntryFuncHead(w io.Writer) error {
	p.gasComment(w, "命令行参数")
	p.gasSectionDataStart(w)
	p.gasGlobal(w, kAppArgcName)
	p.gasDefI64(w, kAppArgcName, 0)
	p.gasGlobal(w, kAppArgvName)
	p.gasDefI64(w, kAppArgvName, 0)
	p.gasGlobal(w, kAppEnvpName)
	p.gasDefI64(w, kAppEnvpName, 0)
	fmt.Fprintln(w)
	return nil
}

// 启动函数
func (p *wat2X64Worker) buildEntryFunc(w io.Writer) error {
	switch p.entryFuncName {
	case kFuncStart:
		return p.buildEntryFunc_start(w)
	case kFuncMain, kFuncMain2:
		return p.buildEntryFunc_x(w, kFuncMain)
	default:
		if p.entryFuncName != "" {
			return p.buildEntryFunc_x(w, p.entryFuncName)
		} else {
			return p.buildEntryFunc_start(w)
		}
	}
}

func (p *wat2X64Worker) buildEntryFunc_start(w io.Writer) error {
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

	// 保存命令行参数和 env
	fmt.Fprintf(w, "    mov  rdi, [rsp]\n")
	fmt.Fprintf(w, "    mov  qword ptr [rip+%s], rdi # argc\n", kAppArgcName)
	fmt.Fprintf(w, "    lea  rsi, [rsp+8]\n")
	fmt.Fprintf(w, "    mov  qword ptr [rip+%s], rsi # argv\n", kAppArgvName)
	fmt.Fprintf(w, "    mov  rax, rdi\n")
	fmt.Fprintf(w, "    imul rax, 8\n")
	fmt.Fprintf(w, "    add  rax, rsp # rax = rsp + rdi*8\n")
	fmt.Fprintf(w, "    lea  rax, [rax+16]\n")
	fmt.Fprintf(w, "    mov  qword ptr [rip+%s], rax # envp\n", kAppEnvpName)
	fmt.Fprintf(w, "\n")

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

func (p *wat2X64Worker) buildEntryFunc_x(w io.Writer, entryFuncName string) error {
	p.gasComment(w, "汇编程序入口函数")
	p.gasSectionTextStart(w)

	p.gasGlobal(w, entryFuncName)
	fmt.Fprintf(w, "%s:\n", entryFuncName)

	fmt.Fprintln(w, "    push rbp")
	fmt.Fprintln(w, "    mov  rbp, rsp")
	fmt.Fprintln(w, "    sub  rsp, 32")
	fmt.Fprintln(w)

	// 参数寄存器
	regArg0 := "rcx"
	regArg1 := "rdx"
	regArg2 := "r8"

	if p.cpuType == abi.X64Unix {
		regArg0 = "rdi"
		regArg1 = "rsi"
		regArg2 = "rdx"
	}

	// 如果是 main 函数, 说明不是自定义的 _start, 需要补充获取命令行信息
	if entryFuncName == kFuncMain {
		fmt.Fprintf(w, "    mov  qword ptr [rip+%s], %s # argc\n", kAppArgcName, regArg0)
		fmt.Fprintf(w, "    mov  qword ptr [rip+%s], %s # argv\n", kAppArgvName, regArg1)
		fmt.Fprintf(w, "    mov  qword ptr [rip+%s], %s # rax\n", kAppEnvpName, regArg2)
		fmt.Fprintf(w, "\n")
	}

	if p.m.Memory != nil {
		fmt.Fprintf(w, "    call %s\n", kMemoryInitFuncName)
	}
	if p.m.Table != nil {
		fmt.Fprintf(w, "    call %s\n", kTableInitFuncName)
	}

	if p.m.Start != "" {
		fmt.Fprintf(w, "    call %s\n", kFuncNamePrefix+p.m.Start)
	}
	for _, fn := range p.m.Funcs {
		switch fn.Name {
		case kFuncMain, kFuncMain2:
			fmt.Fprintf(w, "    call %s\n", kFuncNamePrefix+fixName(fn.Name))
		}
	}
	fmt.Fprintln(w)

	p.gasCommentInFunc(w, "return 0")
	fmt.Fprintln(w, "    mov rax, 0")
	fmt.Fprintln(w)

	p.gasCommentInFunc(w, "return")
	fmt.Fprintln(w, "    mov rsp, rbp")
	fmt.Fprintln(w, "    pop rbp")
	fmt.Fprintln(w, "    ret")
	fmt.Fprintln(w)

	return nil
}
