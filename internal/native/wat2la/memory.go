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
		fmt.Fprintf(w, "// memory $%s\n", p.m.Memory.Name)
	}
	if max := p.m.Memory.MaxPages; max > 0 {
		fmt.Fprintf(w, "uint8_t*      %s_memory = NULL;\n", p.opt.Prefix)
		fmt.Fprintf(w, "const int32_t %s_memory_init_max_pages = %d;\n", p.opt.Prefix, max)
		fmt.Fprintf(w, "const int32_t %s_memory_init_pages = %d;\n", p.opt.Prefix, p.m.Memory.Pages)
		fmt.Fprintf(w, "int32_t       %s_memory_size = %d;\n", p.opt.Prefix, 0)
	} else {
		fmt.Fprintf(w, "uint8_t*      %s_memory = NULL;\n", p.opt.Prefix)
		fmt.Fprintf(w, "const int32_t %s_memory_init_max_pages = %d;\n", p.opt.Prefix, p.m.Memory.Pages)
		fmt.Fprintf(w, "const int32_t %s_memory_init_pages = %d;\n", p.opt.Prefix, p.m.Memory.Pages)
		fmt.Fprintf(w, "int32_t       %s_memory_size = %d;\n", p.opt.Prefix, 0)
	}
	fmt.Fprintln(w)

	return nil
}

func (p *wat2laWorker) buildMemory_data(w io.Writer) error {
	return nil
}
