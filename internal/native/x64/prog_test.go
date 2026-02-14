// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x64

import (
	"fmt"
	"testing"

	"wa-lang.org/wa/internal/native/x64/p9x86"
)

func tAssert(tb testing.TB, ok bool, message ...interface{}) {
	tb.Helper()
	if !ok {
		if len(message) != 0 {
			tb.Fatal(fmt.Sprint(append([]interface{}{"assert failed:"}, message...)...))
		} else {
			tb.Fatal("assert failed")
		}
	}
}

func Test_reg2p9Reg(t *testing.T) {
	tAssert(t, reg2p9Reg(REG_EAX) == p9x86.REG_AX)
	tAssert(t, reg2p9Reg(REG_ECX) == p9x86.REG_CX)
	tAssert(t, reg2p9Reg(REG_EDX) == p9x86.REG_DX)
	tAssert(t, reg2p9Reg(REG_EBX) == p9x86.REG_BX)
	tAssert(t, reg2p9Reg(REG_ESP) == p9x86.REG_SP)
	tAssert(t, reg2p9Reg(REG_EBP) == p9x86.REG_BP)
	tAssert(t, reg2p9Reg(REG_ESI) == p9x86.REG_SI)
	tAssert(t, reg2p9Reg(REG_EDI) == p9x86.REG_DI)
	tAssert(t, reg2p9Reg(REG_R8D) == p9x86.REG_R8)
	tAssert(t, reg2p9Reg(REG_R9D) == p9x86.REG_R9)
	tAssert(t, reg2p9Reg(REG_R10D) == p9x86.REG_R10)
	tAssert(t, reg2p9Reg(REG_R11D) == p9x86.REG_R11)
	tAssert(t, reg2p9Reg(REG_R12D) == p9x86.REG_R12)
	tAssert(t, reg2p9Reg(REG_R13D) == p9x86.REG_R13)
	tAssert(t, reg2p9Reg(REG_R14D) == p9x86.REG_R14)
	tAssert(t, reg2p9Reg(REG_R15D) == p9x86.REG_R15)

	tAssert(t, reg2p9Reg(REG_RAX) == p9x86.REG_AX)
	tAssert(t, reg2p9Reg(REG_RCX) == p9x86.REG_CX)
	tAssert(t, reg2p9Reg(REG_RDX) == p9x86.REG_DX)
	tAssert(t, reg2p9Reg(REG_RBX) == p9x86.REG_BX)
	tAssert(t, reg2p9Reg(REG_RSP) == p9x86.REG_SP)
	tAssert(t, reg2p9Reg(REG_RBP) == p9x86.REG_BP)
	tAssert(t, reg2p9Reg(REG_RSI) == p9x86.REG_SI)
	tAssert(t, reg2p9Reg(REG_RDI) == p9x86.REG_DI)
	tAssert(t, reg2p9Reg(REG_R8) == p9x86.REG_R8)
	tAssert(t, reg2p9Reg(REG_R9) == p9x86.REG_R9)
	tAssert(t, reg2p9Reg(REG_R10) == p9x86.REG_R10)
	tAssert(t, reg2p9Reg(REG_R11) == p9x86.REG_R11)
	tAssert(t, reg2p9Reg(REG_R12) == p9x86.REG_R12)
	tAssert(t, reg2p9Reg(REG_R13) == p9x86.REG_R13)
	tAssert(t, reg2p9Reg(REG_R14) == p9x86.REG_R14)
	tAssert(t, reg2p9Reg(REG_R15) == p9x86.REG_R15)
}
