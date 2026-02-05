// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/native/abi"
)

const (
	kMemoryInitFuncName = ".Wa.Memory.initFunc"

	kMemoryAddrName     = ".Wa.Memory.addr"
	kMemoryPagesName    = ".Wa.Memory.pages"
	kMemoryMaxPagesName = ".Wa.Memory.maxPages"

	kMemoryDataOffsetPrefix = ".Wa.Memory.dataOffset."
	kMemoryDataSizePrefix   = ".Wa.Memory.dataSize."
	kMemoryDataPtrPrefix    = ".Wa.Memory.dataPtr."
)

func (p *wat2X64Worker) buildMemory(w io.Writer) error {
	p.gasComment(w, "定义内存")
	p.gasSectionDataStart(w)

	if p.m.Memory == nil {
		p.gasGlobal(w, kMemoryAddrName)
		p.gasDefI64(w, kMemoryAddrName, 0)
		p.gasGlobal(w, kMemoryPagesName)
		p.gasDefI64(w, kMemoryPagesName, 1)
		p.gasGlobal(w, kMemoryMaxPagesName)
		p.gasDefI64(w, kMemoryMaxPagesName, 1)
		fmt.Fprintln(w)
		return nil
	}
	if max := p.m.Memory.MaxPages; max > 0 {
		p.gasGlobal(w, kMemoryAddrName)
		p.gasDefI64(w, kMemoryAddrName, 0)
		p.gasGlobal(w, kMemoryPagesName)
		p.gasDefI64(w, kMemoryPagesName, int64(p.m.Memory.Pages))
		p.gasGlobal(w, kMemoryMaxPagesName)
		p.gasDefI64(w, kMemoryMaxPagesName, int64(max))
		fmt.Fprintln(w)
	} else {
		p.gasGlobal(w, kMemoryAddrName)
		p.gasDefI64(w, kMemoryAddrName, 0)
		p.gasGlobal(w, kMemoryPagesName)
		p.gasDefI64(w, kMemoryPagesName, int64(p.m.Memory.Pages))
		p.gasGlobal(w, kMemoryMaxPagesName)
		p.gasDefI64(w, kMemoryMaxPagesName, int64(p.m.Memory.Pages))
		fmt.Fprintln(w)
	}

	// 生成需要填充内存的 data 段数据
	// 需要在程序启动时调用相关函数进行填充
	if len(p.m.Data) > 0 {
		p.gasComment(w, "内存数据")
		p.gasSectionDataStart(w)
		for i, d := range p.m.Data {
			p.gasComment(w, fmt.Sprintf("memcpy(&Memory[%d], data[%d], size)", d.Offset, i))
			p.gasDefI64(w, fmt.Sprintf("%s%d", kMemoryDataOffsetPrefix, i), int64(d.Offset))
			p.gasDefI64(w, fmt.Sprintf("%s%d", kMemoryDataSizePrefix, i), int64(len(d.Value)))
			p.gasDefString(w, fmt.Sprintf("%s%d", kMemoryDataPtrPrefix, i), toGasString(d.Value))
		}
		fmt.Fprintln(w)
	}

	// 参数寄存器
	regArg0 := "rcx"
	regArg1 := "rdx"
	regArg2 := "r8"

	if p.cpuType == abi.X64Unix {
		regArg0 = "rdi"
		regArg1 = "rsi"
		regArg2 = "rdx"
	}

	// 生成初始化函数
	{
		p.gasComment(w, "内存初始化函数")
		p.gasSectionTextStart(w)
		p.gasGlobal(w, kMemoryInitFuncName)
		p.gasFuncStart(w, kMemoryInitFuncName)

		fmt.Fprintln(w, "    push rbp")
		fmt.Fprintln(w, "    mov  rbp, rsp")
		fmt.Fprintln(w, "    sub  rsp, 32")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "分配内存")
		fmt.Fprintf(w, "    mov  %s, [rip + %s]\n", regArg0, kMemoryMaxPagesName)
		fmt.Fprintf(w, "    shl  %s, 16\n", regArg0)
		fmt.Fprintf(w, "    call %s\n", kRuntimeMalloc)
		fmt.Fprintf(w, "    mov  [rip + %s], rax\n", kMemoryAddrName)
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "内存清零")
		fmt.Fprintf(w, "    mov  %s, [rip + %s]\n", regArg0, kMemoryAddrName)
		fmt.Fprintf(w, "    mov  %s, 0\n", regArg1)
		fmt.Fprintf(w, "    mov  %s, [rip + %s]\n", regArg2, kMemoryMaxPagesName)
		fmt.Fprintf(w, "    shl  %s, 16\n", regArg2)
		fmt.Fprintf(w, "    call %s\n", kRuntimeMemset)
		fmt.Fprintln(w)

		if len(p.m.Data) > 0 {
			p.gasCommentInFunc(w, "初始化内存")
			fmt.Fprintln(w)

			for i, d := range p.m.Data {
				p.gasCommentInFunc(w, fmt.Sprintf("# memcpy(&Memory[%d], data[%d], size)", d.Offset, i))

				fmt.Fprintf(w, "    mov  rax, [rip + %s]\n", kMemoryAddrName)
				fmt.Fprintf(w, "    mov  %s, [rip + %s%d]\n", regArg0, kMemoryDataOffsetPrefix, i)
				fmt.Fprintf(w, "    add  %s, rax\n", regArg0)
				fmt.Fprintf(w, "    lea  %s, [rip + %s%d]\n", regArg1, kMemoryDataPtrPrefix, i)
				fmt.Fprintf(w, "    mov  %s, [rip + %s%d]\n", regArg2, kMemoryDataSizePrefix, i)
				fmt.Fprintf(w, "    call %s\n", kRuntimeMemcpy)
			}

			fmt.Fprintln(w)
		}

		p.gasCommentInFunc(w, "函数返回")
		fmt.Fprintln(w, "    mov rsp, rbp")
		fmt.Fprintln(w, "    pop rbp")
		fmt.Fprintln(w, "    ret")
		fmt.Fprintln(w)
	}

	return nil
}
