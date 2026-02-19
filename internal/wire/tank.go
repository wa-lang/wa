// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

type register struct {
	id  int
	typ Type
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
