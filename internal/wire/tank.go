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

type Register struct {
	id  int
	typ Type

	_kind  RegisterKind
	_otype Type
}

func (r *Register) Type() Type { return r.typ }

func (r *Register) Kind() RegisterKind { return r._kind }
func (r *Register) OType() Type        { return r._otype }
func (r *Register) String() string {
	if r.id == -1 || r.typ == nil {
		panic(fmt.Sprintf("register is invalid: %d", r.id))
	}
	return fmt.Sprintf("$r%d", r.id)
}

type Tank struct {
	typ      Type
	Register Register
	Member   []*Tank
}

func (t *Tank) Type() Type { return t.typ }

func (t *Tank) String() string {
	var b strings.Builder
	if len(t.Member) > 0 {
		b.WriteByte('[')
		for i, m := range t.Member {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(m.String())
		}
		b.WriteByte(']')
	} else {
		b.WriteString(fmt.Sprintf("$r%d", t.Register.id))
	}
	return b.String()
}

func (t *Tank) Raw() []Register {
	if len(t.Member) == 0 {
		return []Register{t.Register}
	}

	regs := []Register{}
	for _, m := range t.Member {
		regs = append(regs, m.Raw()...)
	}
	return regs
}
