// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"testing"
)

func TestTableAndFunctionPointers(t *testing.T) {
	fooPtr1 := GetFuncPtr_foo()
	fooPtr2 := GetFuncPtr_foo()
	barPtr := GetFuncPtr_bar()

	if fooPtr1 == 0 {
		t.Fatal("获取的 foo 函数指针为空")
	}
	if fooPtr1 != fooPtr2 {
		t.Errorf("多次获取 foo 指针结果不一致: %x vs %x", fooPtr1, fooPtr2)
	}
	if fooPtr1 == barPtr {
		t.Errorf("foo 和 bar 指针碰撞: %x", fooPtr1)
	}

	expectedSize := int32(10)
	if size := TableSize(); size != expectedSize {
		t.Errorf("Table 大小预期 %d, 实际 %d", expectedSize, size)
	}

	TabelSet(0, fooPtr1)
	TabelSet(1, barPtr)
	TabelSet(9, 0xABCDEF)

	t.Run("VerifyTabelGet", func(t *testing.T) {
		if got := TabelGet(0); got != fooPtr1 {
			t.Errorf("Table[0] 读写不匹配: 预期 foo 指针 %x, 得到 %x", fooPtr1, got)
		}
		if got := TabelGet(1); got != barPtr {
			t.Errorf("Table[1] 读写不匹配: 预期 bar 指针 %x, 得到 %x", barPtr, got)
		}
		if got := TabelGet(9); got != 0xABCDEF {
			t.Errorf("Table[9] 读写不匹配: 预期 0xABCDEF, 得到 %x", got)
		}
	})
}

func TestTableInitialization(t *testing.T) {
	for i := int32(0); i < TableSize(); i++ {
		TabelSet(i, 0)
	}

	for i := int32(0); i < TableSize(); i++ {
		if got := TabelGet(i); got != 0 {
			t.Errorf("Table[%d] 初始值不为 0: %x", i, got)
		}
	}
}
