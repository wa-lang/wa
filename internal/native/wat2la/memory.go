// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"
)

// 内存由一个固定的寄存器作为基地址
// 和全局的 global 和 text 段是分开的

func (p *wat2laWorker) buildMemory(w io.Writer) error {
	if p.m.Memory == nil {
		return nil
	}
	if p.m.Memory.Name != "" {
		fmt.Fprintf(w, "# memory $%s\n", p.m.Memory.Name)
	}
	if max := p.m.Memory.MaxPages; max > 0 {
		fmt.Fprintf(w, "global $wat2la.memory.addr: i64 = 0\n")
		fmt.Fprintf(w, "global $wat2la.memory.pages: i64 = %d\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "global $wat2la.memory.maxPages: i64 = %d\n", max)
	} else {
		fmt.Fprintf(w, "global $wat2la.memory.addr: i64 = 0\n")
		fmt.Fprintf(w, "global $wat2la.memory.pages: i64 = %d\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "global $wat2la.memory.maxPages: i64 = %d\n", p.m.Memory.Pages)
	}
	fmt.Fprintln(w)

	// 生成需要填充内存的 data 段数据
	// 需要在程序启动时调用相关函数进行填充
	for i, d := range p.m.Data {
		fmt.Fprintf(w, "# name = %s\n", d.Name)
		fmt.Fprintf(w, "global $wat2la.data.%08x.offset: i64 = 0x%08x\n", i, d.Offset)
		fmt.Fprintf(w, "global $wat2la.data.%08x.size: i64 = %d\n", i, len(d.Value))
		fmt.Fprintf(w, "global $wat2la.data.%08x.data: i64 = %q\n", i, d.Value)
	}

	// 生成初始化函数
	fmt.Fprintln(w, "func $wat2la.memory.init {")
	{
		fmt.Fprintln(w, "    # a0 = 0")
		fmt.Fprintln(w, "    addi.d a0, zero, 0")

		fmt.Fprintf(w, "    # a1 = $wat2la.memory.maxPages * 65536\n")
		fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20($wat2la.memory.maxPages)\n")
		fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12($wat2la.memory.maxPages)\n")
		fmt.Fprintf(w, "    ld.d      a1, t1, 0\n")
		fmt.Fprintf(w, "    lu12i.w   t0, 0x10   # t0 = 0x10 << 12 = 0x10000, 一页的大小\n")
		fmt.Fprintf(w, "    mul.d     a1, a1, t0\n")

		fmt.Fprintln(w, "    # a2 = READ | WRITE = 3")
		fmt.Fprintln(w, "    addi.d a2, zero, 0")

		fmt.Fprintln(w, "    # a3 = flags = 0")
		fmt.Fprintln(w, "    addi.d a3, zero, 0")

		fmt.Fprintln(w, "    # a4 = offset = -1")
		fmt.Fprintln(w, "    addi.d a4, zero, 0")
		fmt.Fprintln(w, "    # a5 = 0")
		fmt.Fprintln(w, "    addi.d a5, zero, 0")

		fmt.Fprintln(w, "    # call $syscall.mmap")
		fmt.Fprintln(w, "    bl $syscall.mmap")

		fmt.Fprintf(w, "    # $wat2la.memory.addr = a0\n")
		fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20($wat2la.memory.addr)\n")
		fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12($wat2la.memory.addr)\n")
		fmt.Fprintf(w, "    st.d      a0, t1, 0\n")

		fmt.Fprintln(w, "    # t8 = a0")
		fmt.Fprintln(w, "    addi.d t0, a8, 0")

		for i := range p.m.Data {
			fmt.Fprintf(w, "    # dst = data[%d].offset\n", i)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20($wat2la.data.%08x.offset)\n", i)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12($wat2la.data.%08x.offset)\n", i)
			fmt.Fprintf(w, "    ld.d      a0, t1, 0\n")
			fmt.Fprintf(w, "    add.d     a0, a0, t8\n")

			fmt.Fprintf(w, "    # ptr = data[%d].data\n", i)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20($wat2la.data.%08x.data)\n", i)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12($wat2la.data.%08x.data)\n", i)
			fmt.Fprintf(w, "    ld.d      a1, t1, 0\n")

			fmt.Fprintf(w, "    # size = data[%d].size\n", i)
			fmt.Fprintf(w, "    pcalau12i t1, %%pc_hi20($wat2la.data.%08x.size)\n", i)
			fmt.Fprintf(w, "    addi.d    t1, t1, %%pc_lo12($wat2la.data.%08x.size)\n", i)
			fmt.Fprintf(w, "    ld.d      a2, t1, 0\n")

			fmt.Fprintf(w, "    # call $builtin.memcpy\n")
			fmt.Fprintf(w, "    bl $builtin.memcpy\n")
		}

		fmt.Fprintln(w, "    # return")
		fmt.Fprintln(w, "    jirl zero, ra, 0")
	}
	fmt.Fprintln(w, "}")

	return nil
}
