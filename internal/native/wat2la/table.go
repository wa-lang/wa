// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2laWorker) buildTable(w io.Writer) error {
	if p.m.Table == nil {
		return nil
	}
	if p.m.Table.Type != token.FUNCREF {
		return fmt.Errorf("unsupported table type: %s", p.m.Table.Type)
	}

	const IntSize = 8

	if p.m.Table.Name != "" {
		fmt.Fprintf(w, "# table $%s\n", p.m.Table.Name)
	}
	if max := p.m.Table.MaxSize; max > 0 {
		fmt.Fprintf(w, "# table.size = %d, table.max = %d\n", p.m.Table.Size, max)
		fmt.Fprintf(w, "global $table: %d = {\n", max*IntSize)
	} else {
		fmt.Fprintf(w, "# table.size = %d\n", p.m.Table.Size)
		fmt.Fprintf(w, "global $table: %d = {\n", p.m.Table.Size*IntSize)
	}
	for _, elem := range p.m.Elem {
		off := elem.Offset
		for i, fn := range elem.Values {
			fmt.Fprintf(w, "    %d: %s,\n", int(off)+i*IntSize, fn)
		}
	}
	fmt.Fprintf(w, "}\n")
	fmt.Fprintln(w)

	return nil
}
