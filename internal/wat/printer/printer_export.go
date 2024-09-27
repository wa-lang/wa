// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import (
	"fmt"

	"wa-lang.org/wa/internal/wat/token"
)

func (p *watPrinter) printExport() error {
	if len(p.m.Exports) == 0 {
		return nil
	}
	for _, e := range p.m.Exports {
		switch e.Kind {
		case token.GLOBAL:
			fmt.Fprintf(p.w, `%s(export "%s" (global %s))`+"\n",
				p.indent, e.Name, p.identOrIndex(e.GlobalIdx),
			)
		case token.FUNC:
			// skip
		case token.MEMORY:
			fmt.Fprintf(p.w, `%s(export "%s" (memory %s))`+"\n",
				p.indent, e.Name, p.identOrIndex(e.MemoryIdx),
			)
		case token.TABLE:
			fmt.Fprintf(p.w, `%s(export "%s" (table %s))`+"\n",
				p.indent, e.Name, p.identOrIndex(e.TableIdx),
			)
		default:
			panic("unreachable")
		}
	}
	return nil
}
