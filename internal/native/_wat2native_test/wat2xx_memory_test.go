// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"testing"
)

func TestMemoryManagement(t *testing.T) {
	initSize := MemorySize()

	oldSize := MemoryGrow(2)
	if oldSize != initSize {
		t.Errorf("首次增长返回值预期 %d, 实际 %d", initSize, oldSize)
	}

	if s := MemorySize(); s != initSize+2 {
		t.Errorf("增长后大小预期 %d, 实际 %d", initSize+2, s)
	}

	if x := MemoryGrow(20); x != -1 {
		t.Errorf("期望增长失败, 实际 %d", x)
	}
}

func TestI32LoadStore(t *testing.T) {
	addr := int32(100)
	offset := int32(4)
	val := int32(0x12345678)

	I32Store(addr, offset, val)
	got := I32Load(addr, offset)
	if got != val {
		t.Errorf("I32Store/Load 不匹配: 预期 0x%x, 得到 0x%x", val, got)
	}
}

func TestI64LoadStore(t *testing.T) {
	addr := int32(100)
	offset := int32(4)
	val := int64(0x12345678)

	I64Store(addr, offset, val)
	got := I64Load(addr, offset)
	if got != val {
		t.Errorf("I64Store/Load 不匹配: 预期 0x%x, 得到 0x%x", val, got)
	}
}
