// 版权 @2025 凹语言 作者。保留所有权利。

package malloc

import "testing"

func TestHeap(t *testing.T) {
	h := NewHeap(nil)
	got := h.Malloc(1)
	defer h.Free(got)
	expect := DefaultHeapBase + kFreeListHeadSize + kBlockHeadSize
	if got != expect {
		t.Fatalf("invalid ptr: expect = %d, got = %d", expect, got)
	}
}
