// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import (
	"fmt"

	"wa-lang.org/wa/internal/wat/ast"
)

func (p *watPrinter) printTypes() error {
	if p.isTypeEmpty() {
		return nil
	}

	for _, typ := range p.m.Types {
		p.printType_func(typ)
	}

	return nil
}

func (p *watPrinter) printType_func(typ *ast.TypeSection) {
	fmt.Fprint(p.w, p.indent)
	fmt.Fprint(p.w, "(type")
	if s := typ.Name; s != "" {
		fmt.Fprintf(p.w, " %s", p.identOrIndex(s))
	}

	fnType := typ.Type
	fmt.Fprint(p.w, " (func")
	if len(fnType.Params) > 0 {
		for _, x := range fnType.Params {
			fmt.Fprintf(p.w, " (param %v)", x.Type)
		}
	}
	if len(fnType.Results) > 0 {
		fmt.Fprintf(p.w, " (result")
		for _, x := range fnType.Results {
			fmt.Fprintf(p.w, " %v", x)
		}
		fmt.Fprint(p.w, ")")
	}

	fmt.Fprintln(p.w, "))")
	return
}
