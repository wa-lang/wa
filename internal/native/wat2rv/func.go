// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2rv

import (
	watast "wa-lang.org/wa/internal/wat/ast"
)

func (p *wat2rvWorker) buildFuncs() error {
	for _, f := range p.m.Funcs {
		if err := p.buildFunc(f); err != nil {
			return err
		}
	}
	return nil
}

func (p *wat2rvWorker) buildFunc(fn *watast.Func) error {
	p.buildFunc_args(fn)
	p.buildFunc_local(fn)
	p.buildFunc_body(fn)
	return nil
}

func (p *wat2rvWorker) buildFunc_args(fn *watast.Func) error {
	return nil
}

func (p *wat2rvWorker) buildFunc_local(fn *watast.Func) error {
	return nil
}
