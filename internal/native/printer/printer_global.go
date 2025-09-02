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
			case token.I32, token.I64, token.F32, token.F64:
				fmt.Fprintf(p.w, "global %s:%v = 0\n", g.Name, g.Type)
			default:
				fmt.Fprintf(p.w, "global %s:%d = {}\n", g.Name, g.Size)
			}
		case len(g.Init) == 1:
			xInit := g.Init[0]
			switch {
			case xInit.LitValue != nil:
				if xInit.Type == xInit.LitValue.Type {
					// 不需要转义
					switch xInit.LitValue.Type {
					case token.I32:
						if xInit.Offset != 0 {
							fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, int32(xInit.LitValue.IntValue))
						} else {
							fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, int32(xInit.LitValue.IntValue))
						}
					case token.I64:
						if xInit.Offset != 0 {
							fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, int64(xInit.LitValue.IntValue))
						} else {
							fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, int64(xInit.LitValue.IntValue))
						}
					case token.F32:
						if xInit.Offset != 0 {
							fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, float32(xInit.LitValue.FloatValue))
						} else {
							fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, float32(xInit.LitValue.FloatValue))
						}
					case token.F64:
						if xInit.Offset != 0 {
							fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, float64(xInit.LitValue.FloatValue))
						} else {
							fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, float64(xInit.LitValue.FloatValue))
						}
					default:
						if xInit.LitValue.StrValue != nil {
							if xInit.Offset != 0 {
								fmt.Fprintf(p.w, "global %s:%v = {%d: %q}\n", g.Name, g.Size, xInit.Offset, *xInit.LitValue.StrValue)
							} else {
								fmt.Fprintf(p.w, "global %s:%v = %q\n", g.Name, g.Size, *xInit.LitValue.StrValue)
							}
						} else {
							panic("unreachable")
						}
					}
				} else {
					// 需要转义
					switch xInit.LitValue.Type {
					case token.I32:
						if xInit.Offset != 0 {
							fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, int32(xInit.LitValue.IntValue))
						} else {
							fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, int32(xInit.LitValue.IntValue))
						}
					case token.I64:
						if xInit.Offset != 0 {
							fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, int64(xInit.LitValue.IntValue))
						} else {
							fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, int64(xInit.LitValue.IntValue))
						}
					case token.F32:
						if xInit.Offset != 0 {
							fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, float32(xInit.LitValue.FloatValue))
						} else {
							fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, float32(xInit.LitValue.FloatValue))
						}
					case token.F64:
						if xInit.Offset != 0 {
							fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, float64(xInit.LitValue.FloatValue))
						} else {
							fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, float64(xInit.LitValue.FloatValue))
						}
					default:
						// 字符串不需要转义
						panic("unreachable")
					}
				}
			case xInit.ConstValue != nil:
				if xInit.Type == xInit.ConstValue.Value.Type {
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, xInit.ConstValue.Name)
					} else {
						fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, xInit.ConstValue.Name)
					}
				} else {
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, xInit.ConstValue.Name)
					} else {
						fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, xInit.ConstValue.Name)
					}
				}
			case xInit.GlobalAddr != nil:
				// 地址的初始化强制类型转义, 因为地址的宽度不是稳定的
				if xInit.Offset != 0 {
					fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, xInit.GlobalAddr.Name)
				} else {
					fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, xInit.GlobalAddr.Name)
				}
			default:
				panic("unreachable")
			}
		default:
			// 多个初始化
			fmt.Fprintf(p.w, "global %s:%v = {\n", g.Name, g.Size)
			for _, xInit := range g.Init {
				switch {
				case xInit.LitValue != nil:
					if xInit.Type == xInit.LitValue.Type {
						// 不需要转义
						switch xInit.LitValue.Type {
						case token.I32:
							if xInit.Offset != 0 {
								fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, int32(xInit.LitValue.IntValue))
							} else {
								fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, int32(xInit.LitValue.IntValue))
							}
						case token.I64:
							if xInit.Offset != 0 {
								fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, int64(xInit.LitValue.IntValue))
							} else {
								fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, int64(xInit.LitValue.IntValue))
							}
						case token.F32:
							if xInit.Offset != 0 {
								fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, float32(xInit.LitValue.FloatValue))
							} else {
								fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, float32(xInit.LitValue.FloatValue))
							}
						case token.F64:
							if xInit.Offset != 0 {
								fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, float64(xInit.LitValue.FloatValue))
							} else {
								fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, float64(xInit.LitValue.FloatValue))
							}
						default:
							if xInit.LitValue.StrValue != nil {
								if xInit.Offset != 0 {
									fmt.Fprintf(p.w, "global %s:%v = {%d: %q}\n", g.Name, g.Size, xInit.Offset, *xInit.LitValue.StrValue)
								} else {
									fmt.Fprintf(p.w, "global %s:%v = %q\n", g.Name, g.Size, *xInit.LitValue.StrValue)
								}
							} else {
								panic("unreachable")
							}
						}
					} else {
						// 需要转义
						switch xInit.LitValue.Type {
						case token.I32:
							if xInit.Offset != 0 {
								fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, int32(xInit.LitValue.IntValue))
							} else {
								fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, int32(xInit.LitValue.IntValue))
							}
						case token.I64:
							if xInit.Offset != 0 {
								fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, int64(xInit.LitValue.IntValue))
							} else {
								fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, int64(xInit.LitValue.IntValue))
							}
						case token.F32:
							if xInit.Offset != 0 {
								fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, float32(xInit.LitValue.FloatValue))
							} else {
								fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, float32(xInit.LitValue.FloatValue))
							}
						case token.F64:
							if xInit.Offset != 0 {
								fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, float64(xInit.LitValue.FloatValue))
							} else {
								fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, float64(xInit.LitValue.FloatValue))
							}
						default:
							// 字符串不需要转义
							panic("unreachable")
						}
					}
				case xInit.ConstValue != nil:
					if xInit.Type == xInit.ConstValue.Value.Type {
						if xInit.Offset != 0 {
							fmt.Fprintf(p.w, "global %s:%v = {%d: %v}", g.Name, xInit.Type, xInit.Offset, xInit.ConstValue.Name)
						} else {
							fmt.Fprintf(p.w, "global %s:%v = %v", g.Name, xInit.Type, xInit.ConstValue.Name)
						}
					} else {
						if xInit.Offset != 0 {
							fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, xInit.ConstValue.Name)
						} else {
							fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, xInit.ConstValue.Name)
						}
					}
				case xInit.GlobalAddr != nil:
					// 地址的初始化强制类型转义, 因为地址的宽度不是稳定的
					if xInit.Offset != 0 {
						fmt.Fprintf(p.w, "global %s:%v = {%d: %s(%v)}", g.Name, xInit.Type, xInit.Offset, xInit.Type, xInit.GlobalAddr.Name)
					} else {
						fmt.Fprintf(p.w, "global %s:%v = %s(%v)", g.Name, xInit.Type, xInit.Type, xInit.GlobalAddr.Name)
					}
				default:
					panic("unreachable")
				}
			}
			fmt.Fprintf(p.w, "}\n")
		}
	}
	return nil
}
