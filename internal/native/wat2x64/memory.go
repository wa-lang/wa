// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"io"
)

const (
	kMemoryInitFuncName = "$Memory.initFunc"

	kMemoryAddrName     = "$Memory.addr"
	kMemoryPagesName    = "$Memory.pages"
	kMemoryMaxPagesName = "$Memory.maxPages"

	kMemoryDataOffsetPrefix = "$Memory.dataOffset."
	kMemoryDataSizePrefix   = "$Memory.dataSize."
	kMemoryDataPtrPrefix    = "$Memory.dataPtr."
)

func (p *wat2X64Worker) buildMemory(w io.Writer) error {
	if p.m.Memory == nil {
		return nil
	}

	p.gasComment(w, "定义内存")
	p.gasSectionDataStart(w)

	maxPages := int64(p.m.Memory.Pages)
	if max := p.m.Memory.MaxPages; max > 0 {
		maxPages = int64(max)
		p.gasDefI64(w, kMemoryAddrName, 0)
		p.gasDefI64(w, kMemoryPagesName, int64(p.m.Memory.Pages))
		p.gasDefI64(w, kMemoryMaxPagesName, int64(max))
		fmt.Fprintln(w)
	} else {
		p.gasDefI64(w, kMemoryAddrName, 0)
		p.gasDefI64(w, kMemoryPagesName, int64(p.m.Memory.Pages))
		p.gasDefI64(w, kMemoryMaxPagesName, int64(p.m.Memory.Pages))
		fmt.Fprintln(w)
	}

	// 生成需要填充内存的 data 段数据
	// 需要在程序启动时调用相关函数进行填充
	if len(p.m.Data) > 0 {
		for i, d := range p.m.Data {
			if len(d.Value) > 10 {
				p.gasComment(w, fmt.Sprintf("Memory[%d]: %v...", d.Offset, string(d.Value[:10])))
			} else {
				p.gasComment(w, fmt.Sprintf("Memory[%d]: %v", d.Offset, string(d.Value)))
			}
			p.gasDefI64(w, fmt.Sprintf("%s%d", kMemoryDataOffsetPrefix, i), int64(d.Offset))
			p.gasDefI64(w, fmt.Sprintf("%s%d", kMemoryDataSizePrefix, i), int64(len(d.Value)))
			p.gasDefString(w, fmt.Sprintf("%s%d", kMemoryDataPtrPrefix, i), string(d.Value))
		}
		fmt.Fprintln(w)
	}

	// 生成初始化函数
	{
		p.gasComment(w, "内存初始化函数")
		p.gasSectionTextStart(w)
		p.gasGlobal(w, kMemoryInitFuncName)
		p.gasFuncStart(w, kMemoryInitFuncName)

		p.gasCommentInFunc(w, "影子空间")
		fmt.Fprintln(w, "    sub rsp, 40")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "分配内存")
		fmt.Fprintf(w, "    mov  rcx, %d # %d pages\n", maxPages*(1<<16), maxPages)
		fmt.Fprintf(w, "    call malloc\n")
		fmt.Fprintf(w, "    lea  rdx, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov  [rdx], rax\n")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "内存清零")
		fmt.Fprintf(w, "    lea  rcx, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov  rdx, 0\n")
		fmt.Fprintf(w, "    mov  r8, %d\n", maxPages*(1<<16))
		fmt.Fprintf(w, "    call %s\n", kRuntimeMemset)
		fmt.Fprintln(w)

		if len(p.m.Data) > 0 {
			p.gasCommentInFunc(w, "初始化内存")
			fmt.Fprintln(w)

			for i, d := range p.m.Data {
				if len(d.Value) > 10 {
					p.gasCommentInFunc(w, fmt.Sprintf("Memory[%d]: %v...", d.Offset, string(d.Value[:10])))
				} else {
					p.gasCommentInFunc(w, fmt.Sprintf("Memory[%d]: %v", d.Offset, string(d.Value)))
				}
				fmt.Fprintf(w, "    lea  rcx, [rip + %s]\n", kMemoryAddrName)
				fmt.Fprintf(w, "    add  rcx, %d\n", d.Offset)
				fmt.Fprintf(w, "    mov  rdx, [rip + %s]\n", fmt.Sprintf("%s%d", kMemoryDataOffsetPrefix, i))
				fmt.Fprintf(w, "    mov  r8, %d\n", len(d.Value))
				fmt.Fprintf(w, "    call %s\n", kRuntimeMemcpy)
			}

			fmt.Fprintln(w)
		}

		p.gasCommentInFunc(w, "函数返回")
		fmt.Fprintln(w, "    add rsp, 40")
		fmt.Fprintln(w, "    ret")
		fmt.Fprintln(w)
	}

	return nil
}
