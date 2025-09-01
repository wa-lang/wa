// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package printer

import (
	"fmt"

	"wa-lang.org/wa/internal/native/token"
)

// global $age: i32 = 5
// global $name = "wa native assembly language"
//
// global $info: 1024 = {
//     5: "abc",    // 从第5字节开始 `abc\0`
//     9: i32(123), // 从第9字节开始
// }

func (p *wsPrinter) printGlobals() error {
	for _, g := range p.f.Globals {
		switch {
		case len(g.Init) == 0:
			switch g.Type {
			case token.I32, token.I64, token.F32, token.F64, token.PTR:
				fmt.Fprintf(p.w, "global %s:%v = 0\n", g.Name, g.Type)
			default:
				fmt.Fprintf(p.w, "global %s:%d = {}\n", g.Name, g.Size)
			}
		case len(g.Init) == 1:
			switch {
			case g.Init[0].Symbal != "":
				fmt.Fprintf(p.w, "global %s:%v = %v\n", g.Name, g.Type, g.Init[0].Symbal)
			case g.Init[0].StrValue != nil:
				fmt.Fprintf(p.w, "global %s:%v = %q\n", g.Name, g.Size, *g.Init[0].StrValue)
			default:
				switch g.Type {
				case token.I32:
					fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, g.Type, int32(g.Init[0].IntValue))
				case token.I64:
					fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, g.Type, int64(g.Init[0].IntValue))
				case token.F32:
					fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, g.Type, float32(g.Init[0].FloatValue))
				case token.F64:
					fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, g.Type, float64(g.Init[0].FloatValue))
				case token.PTR:
					fmt.Fprintf(p.w, "global %s:%v = %x", g.Name, g.Type, g.Init[0].IntValue)
				default:
					panic("unreachable")
				}
			}
		default:
			fmt.Fprintf(p.w, "global %s:%v = {\n", g.Name, g.Size)
			for _, xInit := range g.Init {
				switch {
				case xInit.Type == token.I32 || xInit.Type == token.I64:
					fmt.Fprintf(p.w, "%s %d: %v(%v),\n", p.indent, xInit.Offset, xInit.Type, xInit.IntValue)
				case xInit.Type == token.F32 || xInit.Type == token.F64:
					fmt.Fprintf(p.w, "%s %d: %v(%v),\n", p.indent, xInit.Offset, xInit.Type, xInit.FloatValue)
				case xInit.Symbal != "":
					fmt.Fprintf(p.w, "%s %d: %v,\n", p.indent, xInit.Offset, xInit.Symbal) // TODO: 怎么区分常量或全局变量?
				case xInit.StrValue != nil:
					fmt.Fprintf(p.w, "%s %d: %v,\n", p.indent, xInit.Offset, *xInit.StrValue)
				default:
					panic("unreachable")
				}
			}
			fmt.Fprintf(p.w, "}\n")
		}
	}
	return nil
}
