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
//     5: "abc",    # 从第5字节开始 `abc\0`
//     9: i32(123), # 从第9字节开始
// }

func (p *wsPrinter) printGlobals() error {
	for _, g := range p.f.Globals {
		switch {
		case len(g.Init) == 0:
			assert(g.Size > 0)
			fmt.Fprintf(p.w, "global %s:%d = {}\n", g.Name, g.Size)
		case len(g.Init) == 1:
			xInit := g.Init[0]
			switch {
			case xInit.Value != nil:
				switch xInit.Value.TypeDecor {
				case token.I32:
					assert(xInit.Value.LitKind == token.INT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: i32(%v)}", g.Name, xInit.Value.TypeDecor, xInit.Offset, int32(xInit.Value.IntValue))
					} else {
						fmt.Fprintf(p.w, "global %s:%v = i32(%v)", g.Name, xInit.Value.TypeDecor, int32(xInit.Value.IntValue))
					}
				case token.I64:
					assert(xInit.Value.LitKind == token.INT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: i64(%v)}", g.Name, xInit.Value.TypeDecor, xInit.Offset, int32(xInit.Value.IntValue))
					} else {
						fmt.Fprintf(p.w, "global %s:%v = i64(%v)", g.Name, xInit.Value.TypeDecor, int64(xInit.Value.IntValue))
					}
				case token.F32:
					assert(xInit.Value.LitKind == token.FLOAT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: f32(%v)}", g.Name, xInit.Value.TypeDecor, xInit.Offset, float32(xInit.Value.FloatValue))
					} else {
						fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Value.TypeDecor, float32(xInit.Value.FloatValue))
					}
				case token.F64:
					assert(xInit.Value.LitKind == token.FLOAT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: f64(%v)}", g.Name, xInit.Value.TypeDecor, xInit.Offset, float32(xInit.Value.FloatValue))
					} else {
						fmt.Fprintf(p.w, "global %s:%v = f64(%v)", g.Name, xInit.Value.TypeDecor, float64(xInit.Value.FloatValue))
					}
				case token.STRING:
					assert(xInit.Value.LitKind == token.STRING)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: %q}\n", g.Name, g.Size, xInit.Offset, xInit.Value.StrValue)
					} else {
						fmt.Fprintf(p.w, "global %s:%v = %q\n", g.Name, g.Size, xInit.Value.StrValue)
					}
				default:
					panic("unreachable")
				}
			case xInit.Value.Symbal != "":
				if xInit.Offset != 0 {
					fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Value.TypeDecor, xInit.Offset, xInit.Value.Symbal)
				} else {
					fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Value.TypeDecor, xInit.Value.Symbal)
				}
			default:
				panic("unreachable")
			}
		default:
			// 多个初始化
			fmt.Fprintf(p.w, "global %s:%d = {\n", g.Name, g.Size)
			for _, xInit := range g.Init {
				switch {
				case xInit.Value != nil:
					switch xInit.Value.TypeDecor {
					case token.I32:
						assert(xInit.Value.LitKind == token.INT)
						fmt.Fprintf(p.w, "\t%d = i32(%v),", xInit.Offset, int32(xInit.Value.IntValue))
					case token.I64:
						assert(xInit.Value.LitKind == token.INT)
						fmt.Fprintf(p.w, "\t%d = i64(%v),", xInit.Offset, int64(xInit.Value.IntValue))
					case token.F32:
						assert(xInit.Value.LitKind == token.FLOAT)
						fmt.Fprintf(p.w, "\t%d = f32(%v),", xInit.Offset, float32(xInit.Value.FloatValue))
					case token.F64:
						assert(xInit.Value.LitKind == token.FLOAT)
						fmt.Fprintf(p.w, "\t%d = f64(%v),", xInit.Offset, float64(xInit.Value.FloatValue))
					case token.STRING:
						assert(xInit.Value.LitKind == token.STRING)
						fmt.Fprintf(p.w, "\t%d: %q,\n", xInit.Offset, xInit.Value.StrValue)
					default:
						panic("unreachable")
					}
				case xInit.Value.Symbal != "":
					fmt.Fprintf(p.w, "\t%d: %v,", xInit.Offset, xInit.Value.Symbal)
				default:
					panic("unreachable")
				}
			}
			fmt.Fprintf(p.w, "}\n")
		}
	}
	return nil
}
