// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import "fmt"

// (memory $memory 1024)

func (p *watPrinter) printMemory() error {
	if p.isMemoyEmpty() {
		return nil
	}

	fmt.Fprint(p.w, p.indent)
	fmt.Fprint(p.w, "(memory")
	if s := p.m.Memory.Name; s != "" {
		fmt.Fprintf(p.w, " %s", p.identOrIndex(s))
	}
	fmt.Fprint(p.w, " ", p.m.Memory.Pages)
	if p.m.Memory.MaxPages > 0 {
		fmt.Fprint(p.w, " ", p.m.Memory.MaxPages)
	}
	fmt.Fprintln(p.w, ")")

	return nil
}
