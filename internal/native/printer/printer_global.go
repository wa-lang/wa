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
			case xInit.LitValue != nil:
				switch xInit.Type {
				case token.I32:
					assert(xInit.LitValue.LitKind == token.INT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: i32(%v)}", g.Name, xInit.Type, xInit.Offset, int32(xInit.LitValue.IntValue))
					} else {
						fmt.Fprintf(p.w, "global %s:%v = i32(%v)", g.Name, xInit.Type, int32(xInit.LitValue.IntValue))
					}
				case token.I64:
					assert(xInit.LitValue.LitKind == token.INT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: i64(%v)}", g.Name, xInit.Type, xInit.Offset, int32(xInit.LitValue.IntValue))
					} else {
						fmt.Fprintf(p.w, "global %s:%v = i64(%v)", g.Name, xInit.Type, int64(xInit.LitValue.IntValue))
					}
				case token.F32:
					assert(xInit.LitValue.LitKind == token.FLOAT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: f32(%v)}", g.Name, xInit.Type, xInit.Offset, float32(xInit.LitValue.FloatValue))
					} else {
						fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, float32(xInit.LitValue.FloatValue))
					}
				case token.F64:
					assert(xInit.LitValue.LitKind == token.FLOAT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: f64(%v)}", g.Name, xInit.Type, xInit.Offset, float32(xInit.LitValue.FloatValue))
					} else {
						fmt.Fprintf(p.w, "global %s:%v = f64(%v)", g.Name, xInit.Type, float64(xInit.LitValue.FloatValue))
					}
				case token.STRING:
					assert(xInit.LitValue.LitKind == token.STRING)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: %q}\n", g.Name, g.Size, xInit.Offset, xInit.LitValue.StrValue)
					} else {
						fmt.Fprintf(p.w, "global %s:%v = %q\n", g.Name, g.Size, xInit.LitValue.StrValue)
					}
				default:
					panic("unreachable")
				}
			case xInit.ConstValue != nil:
				switch xInit.Type {
				case token.I32:
					assert(g.Size >= 4)
					assert(xInit.ConstValue.Value.LitKind == token.INT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: i32(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.ConstValue.Name)
					} else {
						fmt.Fprintf(p.w, "global %s:%v = i32(%v)", g.Name, xInit.Type, xInit.ConstValue.Name)
					}
				case token.I64:
					assert(g.Size >= 8)
					assert(xInit.ConstValue.Value.LitKind == token.INT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: i64(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.ConstValue.Name)
					} else {
						fmt.Fprintf(p.w, "global %s:%v = i64(%v)", g.Name, xInit.Type, xInit.ConstValue.Name)
					}
				case token.F32:
					assert(g.Size >= 4)
					assert(xInit.ConstValue.Value.LitKind == token.FLOAT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: f32(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.ConstValue.Name)
					} else {
						fmt.Fprintf(p.w, "global %s:%v = f32(%v)", g.Name, xInit.Type, xInit.ConstValue.Name)
					}
				case token.F64:
					assert(g.Size >= 8)
					assert(xInit.ConstValue.Value.LitKind == token.FLOAT)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: f64(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.ConstValue.Name)
					} else {
						fmt.Fprintf(p.w, "global %s:%v = f64(%v)", g.Name, xInit.Type, xInit.ConstValue.Name)
					}
				case token.STRING:
					// TODO: 长度要特殊处理转义的字符
					assert(xInit.ConstValue.Value.LitKind == token.STRING)
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, xInit.ConstValue.Name)
					} else {
						fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, xInit.ConstValue.Name)
					}
				}
			case xInit.GlobalAddr != nil:
				// 地址的初始化强制类型转义, 因为地址的宽度不是稳定的
				if xInit.Offset != 0 {
					fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, xInit.GlobalAddr.Name)
				} else {
					fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, xInit.GlobalAddr.Name)
				}
			default:
				panic("unreachable")
			}
		default:
			// 多个初始化
			fmt.Fprintf(p.w, "global %s:%d = {\n", g.Name, g.Size)
			for _, xInit := range g.Init {
				switch {
				case xInit.LitValue != nil:
					switch xInit.Type {
					case token.I32:
						assert(xInit.LitValue.LitKind == token.INT)
						fmt.Fprintf(p.w, "\t%d = i32(%v),", xInit.Offset, int32(xInit.LitValue.IntValue))
					case token.I64:
						assert(xInit.LitValue.LitKind == token.INT)
						fmt.Fprintf(p.w, "\t%d = i64(%v),", xInit.Offset, int64(xInit.LitValue.IntValue))
					case token.F32:
						assert(xInit.LitValue.LitKind == token.FLOAT)
						fmt.Fprintf(p.w, "\t%d = f32(%v),", xInit.Offset, float32(xInit.LitValue.FloatValue))
					case token.F64:
						assert(xInit.LitValue.LitKind == token.FLOAT)
						fmt.Fprintf(p.w, "\t%d = f64(%v),", xInit.Offset, float64(xInit.LitValue.FloatValue))
					case token.STRING:
						assert(xInit.LitValue.LitKind == token.STRING)
						fmt.Fprintf(p.w, "\t%d: %q,\n", xInit.Offset, xInit.LitValue.StrValue)
					default:
						panic("unreachable")
					}
				case xInit.ConstValue != nil:
					switch xInit.Type {
					case token.I32:
						assert(g.Size >= 4)
						assert(xInit.ConstValue.Value.LitKind == token.INT)
						fmt.Fprintf(p.w, "\t%d: i32(%v),", xInit.Offset, xInit.ConstValue.Name)
					case token.I64:
						assert(g.Size >= 8)
						assert(xInit.ConstValue.Value.LitKind == token.INT)
						fmt.Fprintf(p.w, "\t%d: i64(%v),", xInit.Offset, xInit.ConstValue.Name)
					case token.F32:
						assert(g.Size >= 4)
						assert(xInit.ConstValue.Value.LitKind == token.FLOAT)
						fmt.Fprintf(p.w, "\t%d: f32(%v),", xInit.Offset, xInit.ConstValue.Name)
					case token.F64:
						assert(g.Size >= 8)
						assert(xInit.ConstValue.Value.LitKind == token.FLOAT)
						fmt.Fprintf(p.w, "\t%d: f64(%v),", xInit.Offset, xInit.ConstValue.Name)
					case token.STRING:
						// TODO: 长度要特殊处理转义的字符
						assert(xInit.ConstValue.Value.LitKind == token.STRING)
						fmt.Fprintf(p.w, "\t%d: %v,", xInit.Offset, xInit.ConstValue.Name)
					}
				case xInit.GlobalAddr != nil:
					// 全局变量地址本地没有固定的类型(和CPU字长有关)
					fmt.Fprintf(p.w, "\t%d: %v,", xInit.Offset, xInit.GlobalAddr.Name)
				default:
					panic("unreachable")
				}
			}
			fmt.Fprintf(p.w, "}\n")
		}
	}
	return nil
}
