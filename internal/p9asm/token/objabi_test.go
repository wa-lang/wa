// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import (
	"testing"

	"wa-lang.org/wa/internal/p9asm/objabi"
)

func TestObjabi(t *testing.T) {
	// 确保预留给普通token类型的空间足够
	if objabi_base > instruction_beg {
		t.Fatal("invalid objabi.ABase")
	}
	if register_beg < instruction_end {
		t.Fatal("invalid objabi.RBase")
	}

	// 验证 objabi.Anames 表格成功初始化
	if len(objabi.Anames) != int(objabi.A_ARCHSPECIFIC) {
		t.Fatal("invalid objabi.Anames")
	}
	if objabi.Anames[0] != "XXX" {
		t.Fatal("invalid objabi.Anames[0]")
	}
	for i := 1; i < int(objabi.ABase); i++ {
		if objabi.Anames[i] != "" {
			t.Fatalf("invalid objabi.Anames[%d]", i)
		}
	}
}
