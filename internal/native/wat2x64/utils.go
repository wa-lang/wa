// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"bytes"
	"fmt"
	"strings"

	nativetok "wa-lang.org/wa/internal/native/token"
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/printer"
	"wa-lang.org/wa/internal/wat/token"
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

func wat2nativeType(typ token.Token) nativetok.Token {
	switch typ {
	case token.I32:
		return nativetok.I32
	case token.I64:
		return nativetok.I64
	case token.F32:
		return nativetok.F32
	case token.F64:
		return nativetok.F64
	default:
		panic("unreachable")
	}
}
