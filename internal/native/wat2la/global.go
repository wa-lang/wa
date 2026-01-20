// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

const (
	kGlobalNamePrefix = ".G."
)

func (p *wat2laWorker) buildGlobal(w io.Writer) error {
	if len(p.m.Globals) == 0 {
		return nil
	}

	p.gasComment(w, "定义全局变量")
	p.gasSectionDataStart(w)
	for _, g := range p.m.Globals {
		switch g.Type {
		case token.I32:
			p.gasDefI32(w, kGlobalNamePrefix+g.Name, g.I32Value)
		case token.I64:
			p.gasDefI64(w, kGlobalNamePrefix+g.Name, g.I64Value)
		case token.F32:
			p.gasDefF32(w, kGlobalNamePrefix+g.Name, g.F32Value)
		case token.F64:
			p.gasDefF64(w, kGlobalNamePrefix+g.Name, g.F64Value)
		default:
			return fmt.Errorf("unsupported global type: %s", g.Type)
		}
	}
	fmt.Fprintln(w)
	return nil
}
