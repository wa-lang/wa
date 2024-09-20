// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import (
	"fmt"

	"wa-lang.org/wa/internal/wat/token"
)

func (p *watPrinter) printGlobals() error {
	if p.isGlobalEmpty() {
		return nil
	}

	for _, g := range p.m.Globals {
		fmt.Fprint(p.w, p.indent)
		fmt.Fprintf(p.w, "(global")
		if g.Name != "" {
			fmt.Fprintf(p.w, " %s", p.identOrIndex(g.Name))
		}
		if g.Mutable {
			fmt.Fprintf(p.w, " (mut %v)", g.Type)
		} else {
			fmt.Fprint(p.w, " ", g.Type)
		}

		switch g.Type {
		case token.I32:
			fmt.Fprintf(p.w, " (i32.const %v)", g.I32Value)
		case token.I64:
			fmt.Fprintf(p.w, " (i64.const %v)", g.I64Value)
		case token.F32:
			fmt.Fprintf(p.w, " (f32.const %v)", g.F32Value)
		case token.F64:
			fmt.Fprintf(p.w, " (f64.const %v)", g.F64Value)
		default:
			panic("unreachable")
		}

		fmt.Fprint(p.w, ")\n")
	}

	return nil
}
