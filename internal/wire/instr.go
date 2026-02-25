// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wire

import (
	"fmt"
	"strings"
)

/**************************************
aStmt: 实现 Stmt 接口与 Pos 相关的方法
包含 aStmt 的对象必须自行实现 Stringer 接口！
**************************************/

type aStmt struct {
	fmt.Stringer
	pos int
}

func (v *aStmt) Pos() int { return v.pos }
func (v *aStmt) Format(tab string, sb *strings.Builder) {
	sb.WriteString(tab)
	sb.WriteString(v.String())
}
