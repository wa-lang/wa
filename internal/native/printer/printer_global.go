// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package printer

import (
	"fmt"
)

// global $age: i32 = 5
// global $name = "wa native assembly language"
//
// global $info: 1024 = {
//     5: "abc",    # 从第5字节开始 `abc\0`
//     9: i32(123), # 从第9字节开始
// }

func (p *wsPrinter) printGlobals() error {
	for _, g := range p.f.Globals {
		switch {
		case len(g.Init) == 0:
			assert(g.Size > 0)
			fmt.Fprintf(p.w, "global %s:%d = {}\n", g.Name, g.Size)
		case len(g.Init) == 1 && g.Init[0].Offset == 0:
			xInit := g.Init[0]
			if xInit.Lit != nil {
				if xInit.Lit.LitKind.DefaultNumberType() != xInit.Lit.TypeCast {
					fmt.Fprintf(p.w, "global %s:%d = %v(%s)\n", g.Name, g.Size, xInit.Lit.TypeCast, xInit.Lit.LitString)
				} else {
					fmt.Fprintf(p.w, "global %s:%d = %s\n", g.Name, g.Size, xInit.Lit.LitString)
				}
			} else {
				fmt.Fprintf(p.w, "global %s:%d = %s\n", g.Name, g.Size, xInit.Symbal)
			}
		default:
			// 多个初始化
			fmt.Fprintf(p.w, "global %s:%d = {\n", g.Name, g.Size)
			for _, xInit := range g.Init {
				if xInit.Lit != nil {
					if xInit.Lit.LitKind.DefaultNumberType() != xInit.Lit.TypeCast {
						fmt.Fprintf(p.w, "\t%d: %v(%s),\n", xInit.Offset, xInit.Lit.TypeCast, xInit.Lit.LitString)
					} else {
						fmt.Fprintf(p.w, "\t%d: %s,\n", xInit.Offset, xInit.Lit.LitString)
					}
				} else {
					fmt.Fprintf(p.w, "\t%d: %s,\n", xInit.Offset, xInit.Symbal)
				}
			}
			fmt.Fprintf(p.w, "}\n")
		}
	}
	return nil
}
