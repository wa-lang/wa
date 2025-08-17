// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package riscv

import "testing"

func TestAnames(t *testing.T) {
	for i, name := range Anames {
		as, ok := LookupAs(name)
		if !ok {
			t.Fatalf("%d: %q not found", i, name)
		}
		if got := AsString(as); got != name {
			t.Fatalf("%d: expect = %q, got = %q", i, name, got)
		}
	}
}

func TestRegisterNames(t *testing.T) {
	for i, name := range Register {
		reg, ok := LookupRegister(name)
		if !ok {
			t.Fatalf("%d: %q not found", i, name)
		}
		if got := RegString(reg); got != name {
			t.Fatalf("%d: expect = %q, got = %q", i, name, got)
		}
	}
}
