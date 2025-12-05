// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package parser

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
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

func lookupBuiltinFn(s string) abi.BuiltinFn {
	for fn := abi.BuiltinFn(1); fn < abi.BuiltinFn_Max; fn++ {
		if s == fn.String() {
			return fn
		}
	}
	return 0
}
