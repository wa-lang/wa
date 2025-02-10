// 版权 @2025 凹语言 作者。保留所有权利。

package malloc

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestHeap_init(t *testing.T) {
	h := NewHeap(nil)

	// 全局变量状态
	tAssertEQ(t, h.Global__stack_ptr(), DefaultStackPtr)
	tAssertEQ(t, h.Global__heap_base(), DefaultHeapBase)
	tAssertEQ(t, h.Global__heap_ptr(), DefaultHeapBase+KFreeListHeadSize)
	tAssertEQ(t, h.Global__heap_top(), DefaultMemoryPages*KPageBytes)
	tAssertEQ(t, h.Global__heap_l128_freep(), h.Global__heap_ptr()-8*2) // 有个 nil 分隔
	tAssertEQ(t, h.Global__heap_lfixed_cap(), DefaultHeapLFixedCap)

	// 空闲链表头状态
	tAssertEQ(t, h.ReadL24Header(), HeapBlock{h.L24HeaderAddr(), 0, 0})
	tAssertEQ(t, h.ReadL32Header(), HeapBlock{h.L32HeaderAddr(), 0, 0})
	tAssertEQ(t, h.ReadL48Header(), HeapBlock{h.L48HeaderAddr(), 0, 0})
	tAssertEQ(t, h.ReadL80Header(), HeapBlock{h.L80HeaderAddr(), 0, 0})
	tAssertEQ(t, h.ReadL128Header(), HeapBlock{h.L128HeaderAddr(), 0, h.L128HeaderAddr()}) // 闭环链表
}

func TestHeap_fixedSize1(t *testing.T) {
	h := NewHeap(nil)
	got1 := h.Malloc(1)
	expect1 := DefaultHeapBase + KFreeListHeadSize + KBlockHeadSize
	tAssertEQ(t, got1, expect1)
	h.Free(got1)

	got2 := h.Malloc(1)
	expect2 := expect1
	tAssertEQ(t, got2, expect2)

	got3 := h.Malloc(1)
	_ = got3
	h.Free(got3)
	if got := h.Malloc(1); got != got3 {
		t.Fatalf("invalid ptr: expect = %d, got = %d", got3, got3)
	}
}

func TestHeap_fixedSize1_100x1000times(t *testing.T) {
	h := NewHeap(nil)
	for i := 0; i < 100; i++ {
		gotList := make([]int32, 100)
		for j := 0; j < len(gotList); j++ {
			gotList[j] = h.Malloc(1)
			if gotList[j] == 0 {
				t.Fatalf("%d:%d: h.Malloc failed", i, j)
			}
		}
		for j := 0; j < len(gotList); j++ {
			h.Free(gotList[j])
		}
	}
}

func TestHeap_largeSize(t *testing.T) {
	cfg := &Config{
		MemoryPages:    1,
		MemoryPagesMax: 2,
		StackPtr:       100,
		HeapBase:       1000,
		HeapLFixedCap:  3,
	}
	h := NewHeap(cfg)

	pOK1 := h.Malloc(KPageBytes - cfg.HeapBase)
	if pOK1 == 0 {
		t.Fatal("failed")
	}
	h.Free(pOK1)

	pOK2 := h.Malloc(KPageBytes - cfg.HeapBase + 100)
	if pOK2 == 0 {
		t.Fatal("failed")
	}

	pFailed := h.Malloc(KPageBytes)
	if pFailed != 0 {
		t.Fatal("failed")
	}
}

func TestHeap_radomMalloc_disableFixed(t *testing.T) {
	cfg := &Config{
		MemoryPages:    1,
		MemoryPagesMax: 2,
		StackPtr:       100,
		HeapBase:       1000,
		HeapLFixedCap:  0, // 禁止fixed策略
	}
	h := NewHeap(cfg)

	rand.Seed(0)

	for i := 0; i < 10; i++ {
		gotList := make([]int32, 100)
		for j := 0; j < len(gotList); j++ {
			size := int32(rand.Intn(2000) + 1)
			gotList[j] = h.Malloc(size)
			if gotList[j] == 0 {
				t.Fatalf("%d:%d: h.Malloc failed", i, j)
			}
		}
		for j := 0; j < len(gotList); j++ {
			if gotList[j] != 0 {
				h.Free(gotList[j])
			}
		}
	}

	got := h.Malloc(100)
	if got == 0 {
		t.Fatal("failed")
	}
}

func TestHeap_radomMalloc(t *testing.T) {
	cfg := &Config{
		MemoryPages:    1,
		MemoryPagesMax: 2,
		StackPtr:       100,
		HeapBase:       1000,
		HeapLFixedCap:  3,
	}
	h := NewHeap(cfg)

	rand.Seed(0)

	for i := 0; i < 10; i++ {
		gotList := make([]int32, 100)
		for j := 0; j < len(gotList); j++ {
			// BUG: j=30 和 j=31 返回相同的 30720
			size := int32(rand.Intn(2000) + 1)
			t.Log(fmt.Sprintf("run: %d %d %d", i, j, size))

			gotList[j] = h.Malloc(size)
			if gotList[j] == 0 {
				t.Fatalf("%d:%d: h.Malloc failed", i, j)
			}
		}
		for j := 0; j < len(gotList); j++ {
			if gotList[j] != 0 {
				h.Free(gotList[j])
			}
		}
	}

	got := h.Malloc(100)
	if got == 0 {
		t.Fatal("failed")
	}
}
func tAssert(tb testing.TB, ok bool) {
	if !ok {
		tb.Helper()
		tb.Fatalf("tAssert failed")
	}
}

func tAssertEQ(tb testing.TB, a, b interface{}) {
	switch a.(type) {
	case int32:
		if a.(int32) != b.(int32) {
			tb.Helper()
			tb.Fatalf("tAssertEQ failed: %d != %d", a, b)
		}
	case HeapBlock:
		if a.(HeapBlock) != b.(HeapBlock) {
			tb.Helper()
			tb.Fatalf("tAssertEQ failed: %v != %v", a, b)
		}
	default:
		panic("unreachable")
	}
}
