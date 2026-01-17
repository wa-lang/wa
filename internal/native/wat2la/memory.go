// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

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

func (p *wat2laWorker) buildMemory(w io.Writer) error {
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

		fmt.Fprintln(w, "    addi.d  $sp, $sp, -16")
		fmt.Fprintln(w, "    st.d    $ra, $sp, 8")
		fmt.Fprintln(w, "    st.d    $fp, $sp, 0")
		fmt.Fprintln(w, "    addi.d  $fp, $sp, 0")
		fmt.Fprintln(w, "    addi.d  $sp, $sp, -32")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "分配内存")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")
		fmt.Fprintf(w, "    slli.d    $a0, $t0, 16\n")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kRuntimeMalloc)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kRuntimeMalloc)
		fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
		fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    st.d      $a0, $t1, 0\n")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "内存清零")
		fmt.Fprintf(w, "    addi.d    $a1, $zero, 0 # a1 = 0\n")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")
		fmt.Fprintf(w, "    slli.d    $a2, $t0, 16\n")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kRuntimeMemset)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kRuntimeMemset)
		fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
		fmt.Fprintln(w)

		if len(p.m.Data) > 0 {
			p.gasCommentInFunc(w, "初始化内存")
			fmt.Fprintln(w)

			for i, d := range p.m.Data {
				p.gasCommentInFunc(w, fmt.Sprintf("# memcpy(&Memory[%d], data[%d], size)", d.Offset, i))

				fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", kMemoryAddrName)
				fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", kMemoryAddrName)
				fmt.Fprintf(w, "    ld.d      $t1, $t1, 0\n")
				fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s%d)\n", kMemoryDataOffsetPrefix, i)
				fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s%d)\n", kMemoryDataOffsetPrefix, i)
				fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")
				fmt.Fprintf(w, "    add.d     $a0, $t1, $t0\n")
				fmt.Fprintf(w, "    pcalau12i $a1, %%pc_hi20(%s%d)\n", kMemoryDataPtrPrefix, i)
				fmt.Fprintf(w, "    addi.d    $a1, $a1, %%pc_lo12(%s%d)\n", kMemoryDataPtrPrefix, i)
				fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s%d)\n", kMemoryDataSizePrefix, i)
				fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s%d)\n", kMemoryDataSizePrefix, i)
				fmt.Fprintf(w, "    ld.d      $a2, $t0, 0\n")
				fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kRuntimeMemcpy)
				fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kRuntimeMemcpy)
				fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
			}

			if len(p.m.Data) > 0 {
				fmt.Fprintln(w)
			}
		}

		p.gasCommentInFunc(w, "函数返回")
		fmt.Fprintln(w, "    addi.d  $sp, $fp, 0")
		fmt.Fprintln(w, "    ld.d    $ra, $sp, 8")
		fmt.Fprintln(w, "    ld.d    $fp, $sp, 0")
		fmt.Fprintln(w, "    addi.d  $sp, $sp, 16")
		fmt.Fprintln(w, "    jirl    $zero, $ra, 0")
		fmt.Fprintln(w)
	}

	return nil
}
