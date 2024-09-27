// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import (
	"wa-lang.org/wa/internal/wat/ast"
)

func (p *watPrinter) isModuleEmpty() bool {
	if len(p.m.Imports) > 0 {
		return false
	}
	if len(p.m.Exports) > 0 {
		return false
	}
	if !p.isMemoyEmpty() {
		return false
	}
	if !p.isTableEmpty() {
		return false
	}
	if !p.isGlobalEmpty() {
		return false
	}
	if !p.isFuncEmpty() {
		return false
	}

	if !p.isTypeEmpty() {
		return false
	}

	return true
}

func (p *watPrinter) isMemoyEmpty() bool {
	if p.m.Memory == nil {
		return true
	}
	if zero := new(ast.Memory); *zero == *p.m.Memory {
		return true
	}
	return false
}

func (p *watPrinter) isTableEmpty() bool {
	if p.m.Table == nil {
		return true
	}
	if zero := new(ast.Table); *zero == *p.m.Table {
		return true
	}
	return false
}

func (p *watPrinter) isGlobalEmpty() bool {
	return len(p.m.Globals) == 0
}

func (p *watPrinter) isFuncEmpty() bool {
	return len(p.m.Funcs) == 0
}

func (p *watPrinter) isTypeEmpty() bool {
	return len(p.m.Types) == 0
}
