// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import (
	"fmt"

	"wa-lang.org/wa/internal/wat/token"
)

// (memory $memory 1024)

func (p *watPrinter) printMemory() error {
	if p.isMemoyEmpty() {
		return nil
	}

	fmt.Fprint(p.w, p.indent)
	fmt.Fprint(p.w, "(memory")
	if s := p.m.Memory.Name; s != "" {
		fmt.Fprintf(p.w, " %s", watPrinter_identOrIndex(s))
	}
	if p.m.Memory.AddrType == token.I64 {
		fmt.Fprintf(p.w, " i64") // 只显式输出 memory64 类型
	}
	fmt.Fprint(p.w, " ", p.m.Memory.Pages)
	if p.m.Memory.MaxPages > 0 {
		fmt.Fprint(p.w, " ", p.m.Memory.MaxPages)
	}
	fmt.Fprintln(p.w, ")")

	return nil
}
