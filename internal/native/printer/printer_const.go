// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package printer

import (
	"fmt"
)

// const $A  = 0x10000000
// const $BB = 12.5
// const $ID = "name"

func (p *wsPrinter) printConsts() error {
	for _, x := range p.f.Consts {
		fmt.Fprintf(p.w, "const %s = %v\n", x.Name, x.Value.LitString)
	}
	return nil
}
