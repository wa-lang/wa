// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	watast "wa-lang.org/wa/internal/wat/ast"
)

func (p *wat2laWorker) buildFuncs() error {
	for _, f := range p.m.Funcs {
		if err := p.buildFunc(f); err != nil {
			return err
		}
	}
	return nil
}

func (p *wat2laWorker) buildFunc(fn *watast.Func) error {
	p.buildFunc_args(fn)
	p.buildFunc_local(fn)
	p.buildFunc_body(fn)
	return nil
}

func (p *wat2laWorker) buildFunc_args(fn *watast.Func) error {
	return nil
}

func (p *wat2laWorker) buildFunc_local(fn *watast.Func) error {
	return nil
}
