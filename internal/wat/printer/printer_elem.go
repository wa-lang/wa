// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import "fmt"

func (p *watPrinter) printElem() error {
	if len(p.m.Elem) == 0 {
		return nil
	}
	for _, e := range p.m.Elem {
		fmt.Fprint(p.w, p.indent)
		fmt.Fprintf(p.w, "(elem (i32.const %d)", e.Offset)
		for _, s := range e.Values {
			fmt.Fprint(p.w, " ", p.identOrIndex(s))
		}
		fmt.Fprintln(p.w, ")")
	}
	return nil
}
