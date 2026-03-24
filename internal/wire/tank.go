// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

type RegisterKind int

const (
	RegisterKindUnknonw RegisterKind = iota
	KLocal
	KImv
	KGlobal
)

type register struct {
	id      int
	typ     Type
	_offset int
	_kind   RegisterKind
	_otype  Type
}

func (r *register) String() string {
	if r.id == -1 || r.typ == nil {
		panic(fmt.Sprintf("register is invalid: %d", r.id))
	}
	return fmt.Sprintf("$r%d", r.id)
}

type tank struct {
	typ      Type
	register register
	member   []*tank
}

func (t *tank) String() string {
	var b strings.Builder
	if len(t.member) > 0 {
		b.WriteByte('[')
		for i, m := range t.member {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(m.String())
		}
		b.WriteByte(']')
	} else {
		b.WriteString(fmt.Sprintf("$r%d", t.register.id))
	}
	return b.String()
}

func (t *tank) raw() []register {
	if len(t.member) == 0 {
		return []register{t.register}
	}

	regs := []register{}
	for _, m := range t.member {
		regs = append(regs, m.raw()...)
	}
	return regs
}
