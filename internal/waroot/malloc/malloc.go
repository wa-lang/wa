// 版权 @2025 凹语言 作者。保留所有权利。

package malloc

import (
	_ "embed"
)

//go:embed malloc.wat
var malloc_wat string

// 默认值
const (
	DefaultMemoryPages    int32 = 1  // 内存页数
	DefaultMemoryPagesMax int32 = 10 // 内存最大页数

	DefaultStackPtr      int32 = 8 << 12  // 栈指针
	DefaultHeapBase      int32 = 10 << 12 // 实际需要编译根据静态数据大小计算
	DefaultHeapLFixedCap int32 = 100      // 固定大小的空闲链表最大长度
)

// Heap配置
type Config struct {
	MemoryPages    int32 // 内存页数
	MemoryPagesMax int32 // 内存最大页数
	StackPtr       int32 // 栈指针(仅用于辅助检查)
	HeapBase       int32 // 启始地址(8字节对齐)
	HeapLFixedCap  int32 // 固定大小空闲链表最大长度
}

// 封装的Heap, 便于测试
type Heap struct {
	cfg *Config
}

// 构造新的Heap
func NewHeap(cfg *Config) *Heap {
	p := &Heap{
		cfg: &Config{
			MemoryPages:    DefaultMemoryPages,
			MemoryPagesMax: DefaultMemoryPagesMax,
			StackPtr:       DefaultStackPtr,
			HeapBase:       DefaultHeapBase,
			HeapLFixedCap:  DefaultHeapLFixedCap,
		},
	}
	if cfg != nil {
		*p.cfg = *cfg
	}
	return p
}

// 分配 size 字节的内存, 返回地址 8 字节对齐
func (p *Heap) Malloc(size int32) int32 {
	return 0
}

// 释放内存
func (p *Heap) Free(ptr int32) {
	return
}

// 打印统计信息
func (p *Heap) String() string {
	return "TODO"
}
