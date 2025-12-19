// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import wattoken "wa-lang.org/wa/internal/wat/token"

// 如果是只读可作为常量

func (p *wat2laWorker) buildGlobal() error {
	for i, g := range p.m.Globals {
		switch g.Type {
		case wattoken.I32:
		case wattoken.I64:
		case wattoken.F32:
		case wattoken.F64:
		default:
			panic("unreachable")
		}
		_ = i
		_ = g
	}
	return nil
}
