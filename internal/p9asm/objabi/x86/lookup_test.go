// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package x86

import "testing"

func TestRegister(t *testing.T) {
	if len(Register) != int(RegMax)-int(RBase) {
		t.Fatal("invalid Register slice")
	}
	if s := RegString(REG_AL); s != "AL" {
		t.Fatal("invalid Register:", s)
	}
	if s := RegString(REG_TR7); s != "TR7" {
		t.Fatal("invalid Register:", s)
	}
}
