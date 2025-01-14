// 版权 @2025 凹语言 作者。保留所有权利。

package malloc

// 空闲链表
type HeapFreeList struct {
	Fixed  bool
	Header HeapBlock

	heap *Heap
}

func (p *HeapFreeList) Len() int32 {
	return 0
}

func (p *HeapFreeList) Cap() int32 {
	return 0
}

func (p *HeapFreeList) Next(iter *HeapBlock) *HeapBlock {
	return &HeapBlock{}
}

func (p *HeapFreeList) String() string {
	return ""
}
