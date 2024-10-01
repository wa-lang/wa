// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import (
	"fmt"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2cWorker) buildHeader() error {
	if len(p.m.Exports) == 0 {
		return nil
	}
	var funcs []*ast.ExportSpec
	for _, e := range p.m.Exports {
		switch e.Kind {
		case token.FUNC:
			funcs = append(funcs, e)
		}
	}
	if len(funcs) == 0 {
		return nil
	}

	fmt.Fprintf(&p.h, "// module %s header file\n", p.m.Name)

	for _, e := range funcs {
		fmt.Fprintf(&p.h, "// extern void %s();\n", e.Name)
	}

	return nil
}
