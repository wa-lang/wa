// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

// 如果是只读可作为常量

func (p *wat2laWorker) buildGlobal(w io.Writer) error {
	if len(p.m.Globals) == 0 {
		return nil
	}
	for _, g := range p.m.Globals {
		fmt.Fprintf(w, "// global $%s: %v\n", g.Name, g.Type)
		switch g.Type {
		case token.I32:
			if g.Mutable {
				fmt.Fprintf(w, "static int32_t %s_%s = %d;\n", p.opt.Prefix, toCName(g.Name), g.I32Value)
			} else {
				fmt.Fprintf(w, "static const int32_t %s_%s = %d;\n", p.opt.Prefix, toCName(g.Name), g.I32Value)
			}
		case token.I64:
			if g.Mutable {
				fmt.Fprintf(w, "static int64_t %s_%s = %d;\n", p.opt.Prefix, toCName(g.Name), g.I64Value)
			} else {
				fmt.Fprintf(w, "static const int64_t %s_%s = %d;\n", p.opt.Prefix, toCName(g.Name), g.I64Value)
			}
		case token.F32:
			if g.Mutable {
				fmt.Fprintf(w, "static float %s_%s = %f;\n", p.opt.Prefix, toCName(g.Name), g.F32Value)
			} else {
				fmt.Fprintf(w, "static const float %s_%s = %f;\n", p.opt.Prefix, toCName(g.Name), g.F32Value)
			}
		case token.F64:
			if g.Mutable {
				fmt.Fprintf(w, "static double %s_%s = %f;\n", p.opt.Prefix, toCName(g.Name), g.F64Value)
			} else {
				fmt.Fprintf(w, "static const double %s_%s = %f;\n", p.opt.Prefix, toCName(g.Name), g.F64Value)
			}
		default:
			return fmt.Errorf("unsupported global type: %s", g.Type)
		}
	}
	fmt.Fprintln(w)
	return nil
}
