// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import "fmt"

func (p *watPrinter) printTable() error {
	if p.isTableEmpty() {
		return nil
	}

	fmt.Fprint(p.w, p.indent)
	fmt.Fprint(p.w, "(table")
	if s := p.m.Table.Name; s != "" {
		fmt.Fprintf(p.w, " %s", p.identOrIndex(s))
	}
	fmt.Fprint(p.w, " ", p.m.Table.Size)
	if p.m.Table.MaxSize > 0 {
		fmt.Fprint(p.w, " ", p.m.Table.MaxSize)
	}

	fmt.Fprint(p.w, " ", p.m.Table.Type)
	fmt.Fprintln(p.w, ")")
	return nil
}
