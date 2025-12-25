// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"bytes"
	"fmt"
	"strings"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/printer"
)

func assert(condition bool, args ...interface{}) {
	if !condition {
		if msg := fmt.Sprint(args...); msg != "" {
			panic(fmt.Sprintf("assert failed, %s", msg))
		} else {
			panic("assert failed")
		}
	}
}

func unreachable() {
	panic("unreachable")
}

// 格式化指令
func insString(i ast.Instruction) string {
	var buf bytes.Buffer
	printer.PrintInstruction(&buf, "", i, 0)
	return strings.TrimSpace(buf.String())
}
