// 版权 @2025 凹语言 作者。保留所有权利。

package malloc

import "testing"

func TestHeap_init(t *testing.T) {
	h := NewHeap(nil)

	// 全局变量状态
	tAssertEQ(t, h.Global__stack_ptr(), DefaultStackPtr)
	tAssertEQ(t, h.Global__heap_base(), DefaultHeapBase)
	tAssertEQ(t, h.Global__heap_ptr(), DefaultHeapBase+kFreeListHeadSize)
	tAssertEQ(t, h.Global__heap_top(), DefaultMemoryPages*KPageBytes)
	tAssertEQ(t, h.Global__heap_l128_freep(), h.Global__heap_ptr()-8)
	tAssertEQ(t, h.Global__heap_lfixed_cap(), DefaultHeapLFixedCap)

	// 空闲链表头状态
	tAssertEQ(t, h.ReadL24Header(), HeapBlock{0, 0})
	tAssertEQ(t, h.ReadL32Header(), HeapBlock{0, 0})
	tAssertEQ(t, h.ReadL46Header(), HeapBlock{0, 0})
	tAssertEQ(t, h.ReadL80Header(), HeapBlock{0, 0})
	tAssertEQ(t, h.ReadL128Header(), HeapBlock{0, 0})
}

func TestHeap_fixedSize1(t *testing.T) {
	h := NewHeap(nil)
	got1 := h.Malloc(1)
	expect1 := DefaultHeapBase + kFreeListHeadSize + kBlockHeadSize
	tAssertEQ(t, got1, expect1)
	h.Free(got1)

	got2 := h.Malloc(1)
	expect2 := expect1
	tAssertEQ(t, got2, expect2)

	got3 := h.Malloc(1)
	_ = got3
	//h.Free(got3) // TODO: 死循环
	//if got := h.Malloc(1); got != got3 {
	//	t.Fatalf("invalid ptr: expect = %d, got = %d", got3, got3)
	//}
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
