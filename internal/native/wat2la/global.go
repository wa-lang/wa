// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2laWorker) buildGlobal(w io.Writer) error {
	if len(p.m.Globals) == 0 {
		return nil
	}
	for _, g := range p.m.Globals {
		switch g.Type {
		case token.I32:
			fmt.Fprintf(w, "global $%s: %v = %d\n", g.Name, g.Type, g.I32Value)
		case token.I64:
			fmt.Fprintf(w, "global $%s: %v = %d\n", g.Name, g.Type, g.I64Value)
		case token.F32:
			fmt.Fprintf(w, "global $%s: %v = %f\n", g.Name, g.Type, g.F32Value)
		case token.F64:
			fmt.Fprintf(w, "global $%s: %v = %f\n", g.Name, g.Type, g.F64Value)
		default:
			return fmt.Errorf("unsupported global type: %s", g.Type)
		}
	}
	fmt.Fprintln(w)
	return nil
}
