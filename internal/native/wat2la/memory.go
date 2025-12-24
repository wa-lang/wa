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
		fmt.Fprintf(w, "global $memory.addr: i64 = 0\n")
		fmt.Fprintf(w, "global $memory.pages: i64 = %d\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "global $memory.maxPages: i64 = %d\n", max)
	} else {
		fmt.Fprintf(w, "global $memory.addr: i64 = 0\n")
		fmt.Fprintf(w, "global $memory.pages: i64 = %d\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "global $memory.maxPages: i64 = %d\n", p.m.Memory.Pages)
	}
	fmt.Fprintln(w)

	// 生成需要填充内存的 data 段数据
	// 需要在程序启动时调用相关函数进行填充
	for i, d := range p.m.Data {
		fmt.Fprintf(w, "# name = %s\n", d.Name)
		fmt.Fprintf(w, "global $data.%08x.offset: i64 = 0x%08x\n", i, d.Offset)
		fmt.Fprintf(w, "global $data.%08x.data: i64 = %q\n", i, d.Value)

	}

	return nil
}
