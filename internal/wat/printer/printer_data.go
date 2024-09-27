// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import (
	"encoding/hex"
	"fmt"
)

func (p *watPrinter) printData() error {
	if len(p.m.Data) == 0 {
		return nil
	}

	for _, d := range p.m.Data {
		fmt.Fprint(p.w, p.indent, "(data")
		if d.Name != "" {
			fmt.Fprint(p.w, p.identOrIndex(d.Name))
		}
		fmt.Fprintf(p.w, " (i32.const %d)", d.Offset)

		fmt.Fprint(p.w, ` "`)
		for i, ch := range hex.EncodeToString(d.Value) {
			if i%2 == 0 {
				fmt.Fprint(p.w, "\\")
			}
			fmt.Fprintf(p.w, "%c", ch)
		}
		fmt.Fprint(p.w, `"`)
		fmt.Fprintln(p.w, ")")
	}

	return nil
}
