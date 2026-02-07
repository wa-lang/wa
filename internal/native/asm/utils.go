// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"fmt"
)

func assert(ok bool, message ...interface{}) {
	if !ok {
		if len(message) != 0 {
			panic(fmt.Sprint(append([]interface{}{"assert failed:"}, message...)...))
		} else {
			panic("assert failed")
		}
	}
}

// 对齐到 a 的倍数
func align(x0, a int64) (x1 int64, padding int) {
	y := x0 + a - 1
	x1 = y - y%a
	padding = int(x1 - x0)
	return
}
