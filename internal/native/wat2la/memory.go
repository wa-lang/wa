// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"
)

const (
	kMemoryInitFuncName = "$wat2la.memory.init"

	kMemoryReg          = "T8"
	kMemoryAddrName     = "$wat2la.memory.addr"
	kMemoryPagesName    = "$wat2la.memory.pages"
	kMemoryMaxPagesName = "$wat2la.memory.maxPages"

	kMemoryDataOffsetPrefix = "$wat2la.memory.data.offset."
	kMemoryDataSizePrefix   = "$wat2la.memory.data.size."
	kMemoryDataPtrPrefix    = "$wat2la.memory.data.ptr."
)

func (p *wat2laWorker) buildMemory(w io.Writer) error {
	if p.m.Memory == nil {
		return nil
	}
	if p.m.Memory.Name != "" {
		fmt.Fprintf(w, "# memory %s\n", p.m.Memory.Name)
	} else {
		fmt.Fprintf(w, "# memory\n")
	}
	if max := p.m.Memory.MaxPages; max > 0 {
		fmt.Fprintf(w, "global %s: i64 = 0\n", kMemoryAddrName)
		fmt.Fprintf(w, "global %s: i64 = %d\n", kMemoryPagesName, p.m.Memory.Pages)
		fmt.Fprintf(w, "global %s: i64 = %d\n", kMemoryMaxPagesName, max)
	} else {
		fmt.Fprintf(w, "global %s: i64 = 0\n", kMemoryAddrName)
		fmt.Fprintf(w, "global %s: i64 = %d\n", kMemoryPagesName, p.m.Memory.Pages)
		fmt.Fprintf(w, "global %s: i64 = %d\n", kMemoryMaxPagesName, p.m.Memory.Pages)
	}
	fmt.Fprintln(w)

	// 生成需要填充内存的 data 段数据
	// 需要在程序启动时调用相关函数进行填充
	for i, d := range p.m.Data {
		fmt.Fprintf(w, "global %s%d: i64 = 0x%08x\n", kMemoryDataOffsetPrefix, i, d.Offset)
		fmt.Fprintf(w, "global %s%d: i64 = %d\n", kMemoryDataSizePrefix, i, len(d.Value))
		fmt.Fprintf(w, "global %s%d: i64 = %q\n", kMemoryDataPtrPrefix, i, d.Value)
	}
	if len(p.m.Data) > 0 {
		fmt.Fprintln(w)
	}

	// 生成初始化函数
	fmt.Fprintf(w, "func %s {\n", kMemoryInitFuncName)
	{
		fmt.Fprintln(w, "    # a0 = 0")
		fmt.Fprintln(w, "    addi.d a0, zero, 0")
		fmt.Fprintln(w)

		fmt.Fprintf(w, "    # a1 = %s * 65536\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s)\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s)\n", kMemoryMaxPagesName)
		fmt.Fprintf(w, "    ld.d      a1, t1, 0\n")
		fmt.Fprintf(w, "    lu12i.w   t0, 0x10   # t0 = 0x10 << 12 = 0x10000, 一页的大小\n")
		fmt.Fprintf(w, "    mul.d     a1, a1, t0\n")
		fmt.Fprintln(w)

		fmt.Fprintln(w, "    # a2 = READ | WRITE = 3")
		fmt.Fprintln(w, "    addi.d a2, zero, 0")
		fmt.Fprintln(w)

		fmt.Fprintln(w, "    # a3 = flags = 0")
		fmt.Fprintln(w, "    addi.d a3, zero, 0")
		fmt.Fprintln(w)

		fmt.Fprintln(w, "    # a4 = offset = -1")
		fmt.Fprintln(w, "    addi.d a4, zero, 0")
		fmt.Fprintln(w, "    # a5 = 0")
		fmt.Fprintln(w, "    addi.d a5, zero, 0")
		fmt.Fprintln(w)

		fmt.Fprintf(w, "    # call %s\n", kSysMmap)
		fmt.Fprintf(w, "    bl %s\n", kSysMmap)
		fmt.Fprintln(w)

		fmt.Fprintf(w, "    # $wat2la.memory.addr = a0\n")
		fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s)\n", kMemoryAddrName)
		fmt.Fprintf(w, "    st.d      a0, t1, 0\n")
		fmt.Fprintln(w)

		fmt.Fprintf(w, "    # %s = a0\n", kMemoryReg)
		fmt.Fprintf(w, "    addi.d %s, a0, 0", kMemoryReg)
		fmt.Fprintln(w)

		for i := range p.m.Data {
			fmt.Fprintf(w, "    # dst = data[%d].offset + base\n", i)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s%d)\n", kMemoryDataOffsetPrefix, i)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s%d)\n", kMemoryDataOffsetPrefix, i)
			fmt.Fprintf(w, "    ld.d      a0, t1, 0\n")
			fmt.Fprintf(w, "    add.d     a0, a0, %s\n", kMemoryReg)
			fmt.Fprintln(w)

			fmt.Fprintf(w, "    # ptr = data[%d].ptr\n", i)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s%d)\n", kMemoryDataPtrPrefix, i)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s%d)\n", kMemoryDataPtrPrefix, i)
			fmt.Fprintf(w, "    ld.d      a1, t1, 0\n")
			fmt.Fprintln(w)

			fmt.Fprintf(w, "    # size = data[%d].size\n", i)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20(%s%08x)\n", kMemoryDataSizePrefix, i)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12(%s%08x)\n", kMemoryDataSizePrefix, i)
			fmt.Fprintf(w, "    ld.d      a2, t1, 0\n")
			fmt.Fprintln(w)

			fmt.Fprintf(w, "    # call %s\n", kBuiltinMemcpy)
			fmt.Fprintf(w, "    bl %s\n", kBuiltinMemcpy)
			fmt.Fprintln(w)
		}

		fmt.Fprintln(w, "    # return")
		fmt.Fprintln(w, "    jirl zero, ra, 0")
	}
	fmt.Fprintln(w, "}")

	return nil
}
