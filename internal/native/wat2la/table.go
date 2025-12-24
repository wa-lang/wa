// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

// table 是地址数组
// global $info: 1024 = {
//     4:  add,  # 函数
//     12: main, # 函数
// }

func (p *wat2laWorker) buildTable(w io.Writer) error {
	if p.m.Table == nil {
		return nil
	}
	if p.m.Table.Type != token.FUNCREF {
		return fmt.Errorf("unsupported table type: %s", p.m.Table.Type)
	}

	if p.m.Table.Name != "" {
		fmt.Fprintf(w, "// table $%s\n", p.m.Table.Name)
	}
	if max := p.m.Table.MaxSize; max > 0 {
		fmt.Fprintf(w, "uintptr_t %s_table[%d];\n", p.opt.Prefix, max)
		fmt.Fprintf(w, "const int %s_table_init_max_size = %d;\n", p.opt.Prefix, max)
		fmt.Fprintf(w, "int32_t   %s_table_size = %d;\n", p.opt.Prefix, p.m.Table.Size)
	} else {
		fmt.Fprintf(w, "uintptr_t %s_table[%d];\n", p.opt.Prefix, p.m.Table.Size)
		fmt.Fprintf(w, "const int %s_table_init_max_size = %d;\n", p.opt.Prefix, p.m.Table.Size)
		fmt.Fprintf(w, "int32_t   %s_table_size = %d;\n", p.opt.Prefix, p.m.Table.Size)
	}
	fmt.Fprintln(w)

	return nil
}

func (p *wat2laWorker) buildTable_elem(w io.Writer) error {
	return nil
}
