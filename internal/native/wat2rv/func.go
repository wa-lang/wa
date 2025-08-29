// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2rv

import (
	"fmt"

	watast "wa-lang.org/wa/internal/wat/ast"
	wattoken "wa-lang.org/wa/internal/wat/token"
)

func (p *wat2rvWorker) buildFunc(fn *watast.Func) error {
	p.buildFunc_args(fn)
	p.buildFunc_local(fn)
	p.buildFunc_body(fn)
	panic("TODO")
}

func (p *wat2rvWorker) buildFunc_args(fn *watast.Func) error {
	panic("TODO")
}

func (p *wat2rvWorker) buildFunc_local(fn *watast.Func) error {
	panic("TODO")
}

func (p *wat2rvWorker) buildFunc_body(fn *watast.Func) error {
	for _, ins := range fn.Body.Insts {
		if err := p.buildFunc_ins(fn, ins, 1); err != nil {
			return err
		}
	}
	panic("TODO")
}

func (p *wat2rvWorker) buildFunc_ins(fn *watast.Func, i watast.Instruction, level int) error {
	switch tok := i.Token(); tok {
	default:
		panic(fmt.Sprintf("unreachable: %T", i))
	case wattoken.INS_UNREACHABLE:
		// TODO
	}
	panic("TODO")
}
