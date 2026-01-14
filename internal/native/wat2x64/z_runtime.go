// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	_ "embed"
	"fmt"
	"io"
)

const (
	_kRuntimeNamePrefix = ".Runtime."

	kRuntimeWrite           = _kRuntimeNamePrefix + "_write"
	kRuntimeExit            = _kRuntimeNamePrefix + "_exit"
	kRuntimeMalloc          = _kRuntimeNamePrefix + "malloc"
	kRuntimeMemcpy          = _kRuntimeNamePrefix + "memcpy"
	kRuntimeMemset          = _kRuntimeNamePrefix + "memset"
	kRuntimePanic           = _kRuntimeNamePrefix + "panic"
	kRuntimePanicMessage    = _kRuntimeNamePrefix + "panic.message"
	kRuntimePanicMessageLen = _kRuntimeNamePrefix + "panic.messageLen"
)

// 生成运行时函数
func (p *wat2X64Worker) buildRuntimeHead(w io.Writer) error {
	var list = []string{
		kRuntimeWrite,
		kRuntimeExit,
		kRuntimeMalloc,
		kRuntimeMemcpy,
		kRuntimeMemset,
		// panic 内部实现
	}

	p.gasComment(w, "运行时函数")
	for _, absName := range list {
		baseName := absName[len(_kRuntimeNamePrefix):]
		p.gasExtern(w, baseName)
	}
	for _, absName := range list {
		baseName := absName[len(_kRuntimeNamePrefix):]
		p.gasSet(w, absName, baseName)
	}
	fmt.Fprintln(w)
	return nil
}

// 生成运行时函数
func (p *wat2X64Worker) buildRuntimeImpl(w io.Writer) error {
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

	fmt.Fprintf(w, "    # runtime.write(stderr, panicMessage, size)\n")
	fmt.Fprintf(w, "    mov  rcx, 2 # stderr\n")
	fmt.Fprintf(w, "    mov  rdx, [rip + %s]\n", kRuntimePanicMessage)
	fmt.Fprintf(w, "    mov  r8, [rip + %s] # size\n", kRuntimePanicMessageLen)
	fmt.Fprintf(w, "    call %s\n", kRuntimePanic)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "    # 退出程序\n")
	fmt.Fprintf(w, "    mov  rcx, 1 # 退出码\n")
	fmt.Fprintf(w, "    call %s\n", kRuntimeExit)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "    # return\n")
	fmt.Fprintf(w, "    mov rsp, rbp\n")
	fmt.Fprintf(w, "    pop rbp\n")
	fmt.Fprintf(w, "    ret\n")
	fmt.Fprintln(w)

	return nil
}
