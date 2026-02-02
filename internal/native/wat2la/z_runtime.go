// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	_ "embed"
	"fmt"
	"io"
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

// 生成运行时函数
func (p *wat2laWorker) buildRuntimeHead(w io.Writer) error {
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
		p.gasExtern(w, absName)
	}
	fmt.Fprintln(w)
	return nil
}

// 生成运行时函数
func (p *wat2laWorker) buildRuntimeImpl(w io.Writer) error {
	if err := p.buildRuntimeImpl_panic(w); err != nil {
		return err
	}
	return nil
}

func (p *wat2laWorker) buildRuntimeImpl_panic(w io.Writer) error {
	const panicMessage = "panic"
	p.gasSectionDataStart(w)
	p.gasDefString(w, kRuntimePanicMessage, panicMessage)
	p.gasDefI64(w, kRuntimePanicMessageLen, int64(len(panicMessage)))
	fmt.Fprintln(w)

	p.gasSectionTextStart(w)
	p.gasGlobal(w, kRuntimePanic)
	p.gasFuncStart(w, kRuntimePanic)

	fmt.Fprintf(w, "    addi.d  $sp, $sp, -16\n")
	fmt.Fprintf(w, "    st.d    $ra, $sp, 8\n")
	fmt.Fprintf(w, "    st.d    $fp, $sp, 0\n")
	fmt.Fprintf(w, "    addi.d  $fp, $sp, 0\n")
	fmt.Fprintf(w, "    addi.d  $sp, $sp, -32\n")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "    # runtime.write(stderr, panicMessage, size)\n")
	fmt.Fprintf(w, "    addi.d    $a0, $zero, 2\n")
	fmt.Fprintf(w, "    pcalau12i $a1, %%pc_hi20(%s)\n", kRuntimePanicMessage)
	fmt.Fprintf(w, "    addi.d    $a1, $a1, %%pc_lo12(%s)\n", kRuntimePanicMessage)
	fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kRuntimePanicMessageLen)
	fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kRuntimePanicMessageLen)
	fmt.Fprintf(w, "    ld.d      $a2, $t0, 0\n")
	fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kRuntimeWrite)
	fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kRuntimeWrite)
	fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "    # 退出程序\n")
	fmt.Fprintf(w, "    addi.d    $a0, $zero, 1 # 退出码\n")
	fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kRuntimeExit)
	fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kRuntimeExit)
	fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "    # return\n")
	fmt.Fprintf(w, "    addi.d  $sp, $fp, 0\n")
	fmt.Fprintf(w, "    ld.d    $ra, $sp, 8\n")
	fmt.Fprintf(w, "    ld.d    $fp, $sp, 0\n")
	fmt.Fprintf(w, "    addi.d  $sp, $sp, 16\n")
	fmt.Fprintf(w, "    jirl    $zero, $ra, 0\n")
	fmt.Fprintln(w)

	return nil
}
