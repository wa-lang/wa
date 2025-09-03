// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package printer

import (
	"fmt"

	"wa-lang.org/wa/internal/native/token"
)

// const $A  = 0x10000000
// const $BB = 12.5
// const $ID = "name"

func (p *wsPrinter) printConsts() error {
	for _, x := range p.f.Consts {
		switch x.Value.LitKind {
		case token.INT:
			fmt.Fprintf(p.w, "const %s = %v\n", x.Name, x.Value.IntValue)
		case token.FLOAT:
			fmt.Fprintf(p.w, "const %s = %v\n", x.Name, x.Value.FloatValue)
		case token.STRING:
			fmt.Fprintf(p.w, "const %s = %v\n", x.Name, x.Value.StrValue)
		default:
			panic("unreachable")
		}
	}
	return nil
}
