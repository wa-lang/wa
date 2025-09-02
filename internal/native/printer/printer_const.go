// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package printer

import (
	"fmt"

	"wa-lang.org/wa/internal/native/token"
)

// const $UART0       = 0x10000000
// const $EXIT_DEVICE = 0x100000

func (p *wsPrinter) printConsts() error {
	for _, x := range p.f.Consts {
		switch x.Value.Type {
		case token.I32:
			fmt.Fprintf(p.w, "const %s = %v\n", x.Name, int32(x.Value.IntValue))
		case token.I64:
			fmt.Fprintf(p.w, "const %s = %v\n", x.Name, x.Value.IntValue)
		case token.F32:
			fmt.Fprintf(p.w, "const %s = %v\n", x.Name, float32(x.Value.FloatValue))
		case token.F64:
			fmt.Fprintf(p.w, "const %s = %v\n", x.Name, x.Value.FloatValue)
		default:
			if x.Value.StrValue != nil {
				fmt.Fprintf(p.w, "const %s = %q\n", x.Name, *x.Value.StrValue)
			} else {
				panic("unreachable")
			}
		}
	}
	return nil
}
