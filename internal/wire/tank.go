// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

type tank struct {
	register string
	typ      Type
	member   []*tank
}

func initTank(typ Type) *tank {
	ut := typ.Underlying()
	var tank tank
	tank.typ = ut

	switch ut := ut.(type) {
	case *Struct:
		for _, m := range ut.member {
			tm := initTank(m.Type)
			tank.member = append(tank.member, tm)
		}
	}
	return &tank
}
