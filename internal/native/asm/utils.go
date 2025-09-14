// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import "fmt"

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
func align(x, a int64) int64 {
	y := x + a - 1
	return y - y%a
}
