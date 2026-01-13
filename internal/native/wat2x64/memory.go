// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"io"
)

const (
	kMemoryInitFuncName = ".Memory.initFunc"

	kMemoryAddrName     = ".Memory.addr"
	kMemoryPagesName    = ".Memory.pages"
	kMemoryMaxPagesName = ".Memory.maxPages"

	kMemoryDataOffsetPrefix = ".Memory.dataOffset."
	kMemoryDataSizePrefix   = ".Memory.dataSize."
	kMemoryDataPtrPrefix    = ".Memory.dataPtr."
)

func (p *wat2X64Worker) buildMemory(w io.Writer) error {
	if p.m.Memory == nil {
		return nil
	}

	p.gasComment(w, "定义内存")
	p.gasSectionDataStart(w)
	p.gasGlobal(w, kMemoryAddrName)
	p.gasGlobal(w, kMemoryPagesName)
	p.gasGlobal(w, kMemoryMaxPagesName)

	if max := p.m.Memory.MaxPages; max > 0 {
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
		p.gasComment(w, "内存数据")
		p.gasSectionDataStart(w)
		for i, d := range p.m.Data {
			p.gasComment(w, fmt.Sprintf("memcpy(&Memory[%d], data[%d], size)", d.Offset, i))
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
		fmt.Fprintf(w, "    mov  rcx, [rip + %s]\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    shl  rcx, 16\n")
		fmt.Fprintf(w, "    call %s\n", kRuntimeMalloc)
		fmt.Fprintf(w, "    lea  rdx, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov  [rdx], rax\n")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "内存清零")
		fmt.Fprintf(w, "    mov  rcx, [rip + %s]\n", kMemoryAddrName)
		fmt.Fprintf(w, "    mov  rdx, 0\n")
		fmt.Fprintf(w, "    mov  r8, [rip + %s]\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    shl  r8, 16\n")
		fmt.Fprintf(w, "    call %s\n", kRuntimeMemset)
		fmt.Fprintln(w)

		if len(p.m.Data) > 0 {
			p.gasCommentInFunc(w, "初始化内存")
			fmt.Fprintln(w)

			for i, d := range p.m.Data {
				p.gasCommentInFunc(w, fmt.Sprintf("# memcpy(&Memory[%d], data[%d], size)", d.Offset, i))

				fmt.Fprintf(w, "    mov  rax, [rip + %s]\n", kMemoryAddrName)
				fmt.Fprintf(w, "    mov  rcx, [rip + %s%d]\n", kMemoryDataOffsetPrefix, i)
				fmt.Fprintf(w, "    add  rcx, rax\n")
				fmt.Fprintf(w, "    lea  rdx, [rip + %s%d]\n", kMemoryDataPtrPrefix, i)
				fmt.Fprintf(w, "    mov  r8, [rip + %s%d]\n", kMemoryDataSizePrefix, i)
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
