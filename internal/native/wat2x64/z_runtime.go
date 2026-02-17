// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	_ "embed"
	"fmt"
	"io"

	"wa-lang.org/wa/internal/native/abi"
)

const (
	_kRuntimeNamePrefix = ".Wa.Runtime."

	kRuntimeWrite           = _kRuntimeNamePrefix + "write"
	kRuntimeExit            = _kRuntimeNamePrefix + "exit"
	kRuntimeMalloc          = _kRuntimeNamePrefix + "malloc"
	kRuntimeMemcpy          = _kRuntimeNamePrefix + "memcpy"
	kRuntimeMemset          = _kRuntimeNamePrefix + "memset"
	kRuntimeMemmove         = _kRuntimeNamePrefix + "memmove"
	kRuntimePanic           = _kRuntimeNamePrefix + "panic"
	kRuntimePanicMessage    = _kRuntimeNamePrefix + "panic.message"
	kRuntimePanicMessageLen = _kRuntimeNamePrefix + "panic.messageLen"
)

//go:embed assets/native-env-linux-x64.s
var native_env_linux_x64_s string

//go:embed assets/native-env-windows-x64.s
var native_env_windows_x64_s string

// 生成运行时函数
func (p *wat2X64Worker) buildRuntimeHead(w io.Writer) error {
	var list = []string{
		kRuntimeWrite,
		kRuntimeExit,
		kRuntimeMalloc,
		kRuntimeMemcpy,
		kRuntimeMemset,
		kRuntimeMemmove,
		// panic 内部实现
	}

	p.gasComment(w, "运行时函数")
	for _, absName := range list {
		fmt.Fprintf(w, "# .extern %s\n", absName)
	}
	fmt.Fprintln(w)
	return nil
}

// 生成运行时函数
func (p *wat2X64Worker) buildRuntimeImpl(w io.Writer) error {
	switch p.cpuType {
	case abi.X64Unix:
		fmt.Fprintln(w, native_env_linux_x64_s)
	case abi.X64Windows:
		fmt.Fprintln(w, native_env_windows_x64_s)
	default:
		panic("unreachable")
	}

	if err := p.buildRuntimeImpl_panic(w); err != nil {
		return err
	}
	return nil
}

func (p *wat2X64Worker) buildRuntimeImpl_panic(w io.Writer) error {
	const panicMessage = "panic"
	p.gasSectionDataStart(w)
	p.gasDefString(w, kRuntimePanicMessage, panicMessage)
	p.gasDefI64(w, kRuntimePanicMessageLen, int64(len(panicMessage)))
	fmt.Fprintln(w)

	p.gasSectionTextStart(w)
	p.gasGlobal(w, kRuntimePanic)
	p.gasFuncStart(w, kRuntimePanic)

	fmt.Fprintf(w, "    push rbp\n")
	fmt.Fprintf(w, "    mov  rbp, rsp\n")
	fmt.Fprintf(w, "    sub  rsp, 32\n")
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

	fmt.Fprintf(w, "    # runtime.write(stderr, panicMessage, size)\n")
	fmt.Fprintf(w, "    mov  %s, 2 # stderr\n", regArg0)
	fmt.Fprintf(w, "    lea  %s, qword ptr [rip+%s]\n", regArg1, kRuntimePanicMessage)
	fmt.Fprintf(w, "    mov  %s, qword ptr [rip+%s] # size\n", regArg2, kRuntimePanicMessageLen)
	fmt.Fprintf(w, "    call %s\n", kRuntimeWrite)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "    # 退出程序\n")
	fmt.Fprintf(w, "    mov  %s, 1 # 退出码\n", regArg0)
	fmt.Fprintf(w, "    call %s\n", kRuntimeExit)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "    # return\n")
	fmt.Fprintf(w, "    mov rsp, rbp\n")
	fmt.Fprintf(w, "    pop rbp\n")
	fmt.Fprintf(w, "    ret\n")
	fmt.Fprintln(w)

	return nil
}
