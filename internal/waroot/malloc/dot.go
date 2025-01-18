package malloc

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"sort"
)

const tmpl = `
digraph G {
	// 全局设置
	rankdir=LR; // 从左到右的布局方向

	// 定义链表头节点对齐
	subgraph cluster_heads {
		rank=same; // 保证节点在同一层

		root [
			shape=record
			label="{Pages:{{.Pages}}|{HeapBase:{{.HeapBase}}|HeapPtr:{{.HeapPtr}}|HeapTop:{{.HeapTop}}}}"
		];

		L24 [
			shape=record
			label="{L24|{Cap:{{.LfixedCap}}|Size:{{.L24.Size}}|Next:{{.L24.Next}}}}"
		];
		L32 [
			shape=record
			label="{L32|{Cap:{{.LfixedCap}}|Size:{{.L32.Size}}|Next:{{.L32.Next}}}}"
		];
		L48 [
			shape=record
			label="{L48|{Cap:{{.LfixedCap}}|Size:{{.L48.Size}}|Next:{{.L48.Next}}}}"
		];
		L80 [
			shape=record
			label="{L80|{Cap:{{.LfixedCap}}|Size:{{.L80.Size}}|Next:{{.L80.Next}}}}"
		];
		L128 [
			shape=record
			label="{L128|{FreePtr:{{.L128FreePtr}}|Size:{{.L128.Size}}|Next:{{.L128.Next}}}}"
		];

		Used [
			shape=record
			label="{Used|{Count:{{len .UseBlocks}}}}"
		];
	}

	// 已经分配的块
	{{range $block := .UseBlocks}}
	B{{$block.Addr}} [
		shape=record
		label="{Block|{Addr:{{$block.Addr}}|Size:{{$block.Size}}|Next:{{$block.Next}}}}"
	];
	{{end}}

	// 空闲内存块节点
	{{range $block := .FreeBlocks}}
	B{{$block.Addr}} [
		shape=record
		label="{Block|{Addr:{{$block.Addr}}|Size:{{$block.Size}}|Next:{{$block.Next}}}}"
	];
	{{end}}

	// 节点的关系
	{{range $k, $v := .Edges}}
	{{$k}} -> {{$v}};
	{{end}}
}
`

type HeapInfo struct {
	Heap *Heap

	Pages    int32
	HeapBase int32
	HeapPtr  int32
	HeapTop  int32

	// 空闲链表头
	L24, L32, L48, L80, L128 HeapBlock

	// 固定大小的容量
	LfixedCap int32
	// L128当前位置
	L128FreePtr int32

	// 已经分配的块
	UseBlocks []HeapBlock

	// 空闲内存块节点
	FreeBlocks   []HeapBlock
	FreeBlockMap map[int32]HeapBlock

	// 节点的关系
	Edges map[string]string
}

func (p *Heap) DotString() string {
	var info HeapInfo
	info.Heap = p
	info.FreeBlockMap = map[int32]HeapBlock{}
	info.Edges = map[string]string{}

	info.Pages = int32(p.wazeroModule.Memory().Size(context.Background()) / KPageBytes)
	info.HeapBase = p.Global__heap_base()
	info.HeapPtr = p.Global__heap_ptr()
	info.HeapTop = p.Global__heap_top()

	info.L24 = p.ReadL24Header()
	info.L32 = p.ReadL32Header()
	info.L48 = p.ReadL48Header()
	info.L80 = p.ReadL80Header()
	info.L128 = p.ReadL128Header()

	info.LfixedCap = p.Global__heap_lfixed_cap()
	info.L128FreePtr = p.Global__heap_l128_freep()

	// 分配的内存节点
	for _, blk := range p.usedMap {
		info.UseBlocks = append(info.UseBlocks, blk)
	}
	sort.Slice(info.UseBlocks, func(i, j int) bool {
		return info.UseBlocks[i].Addr < info.UseBlocks[j].Addr
	})
	for i, blk := range info.UseBlocks {
		if i == 0 {
			info.Edges["Used"] = fmt.Sprintf("B%d", blk.Addr)
		} else {
			prevAddr := info.UseBlocks[i-1].Addr
			info.Edges[fmt.Sprintf("B%d", prevAddr)] = fmt.Sprintf("B%d", blk.Addr)
		}
	}

	// 收集空闲节点
	info.FreeBlocks = append(info.FreeBlocks, info.LoadFreeBlocks()...)
	sort.Slice(info.FreeBlocks, func(i, j int) bool {
		return info.FreeBlocks[i].Addr < info.FreeBlocks[j].Addr
	})
	for _, blk := range info.FreeBlocks {
		if blk.Addr != 0 {
			info.FreeBlockMap[blk.Addr] = blk
		} else {
			panic(fmt.Sprintf("%+v", blk))
		}
	}

	// 构造空闲链表
	if blk, ok := info.FreeBlockMap[info.L24.Next]; ok {
		info.Edges["L24"] = fmt.Sprintf("B%d", blk.Addr)
	}
	if blk, ok := info.FreeBlockMap[info.L32.Next]; ok {
		info.Edges["L32"] = fmt.Sprintf("B%d", blk.Addr)
	}
	if blk, ok := info.FreeBlockMap[info.L48.Next]; ok {
		info.Edges["L48"] = fmt.Sprintf("B%d", blk.Addr)
	}
	if blk, ok := info.FreeBlockMap[info.L80.Next]; ok {
		info.Edges["L80"] = fmt.Sprintf("B%d", blk.Addr)
	}
	if blk, ok := info.FreeBlockMap[info.L128.Next]; ok {
		info.Edges["L128"] = fmt.Sprintf("B%d", blk.Addr)
	}

	if info.L128.Next == info.L128.Addr {
		info.Edges["L128"] = "L128"
	}

	for _, blk := range info.FreeBlocks {
		if blk.Next != 0 {
			if _, ok := info.FreeBlockMap[blk.Next]; ok {
				info.Edges[fmt.Sprintf("B%d", blk.Addr)] = fmt.Sprintf("B%d", blk.Next)
			}
			if blk.Next == info.L128.Addr {
				info.Edges[fmt.Sprintf("B%d", blk.Addr)] = "L128"
			}
		}
	}

	return info.String()
}

func (p *HeapInfo) LoadFreeBlocks() []HeapBlock {
	var blocks []HeapBlock
	blocks = append(blocks, p.loadFreeBlocks(&p.L24)...)
	blocks = append(blocks, p.loadFreeBlocks(&p.L32)...)
	blocks = append(blocks, p.loadFreeBlocks(&p.L48)...)
	blocks = append(blocks, p.loadFreeBlocks(&p.L80)...)
	blocks = append(blocks, p.loadFreeBlocks(&p.L128)...)
	return blocks
}

func (p *HeapInfo) loadFreeBlocks(head *HeapBlock) []HeapBlock {
	var blocks []HeapBlock
	for addr := head.Next; addr != 0; {
		if addr == p.L128.Addr {
			break // L128 是一个环
		}
		blk := p.Heap.ReadBlock(addr)
		blocks = append(blocks, blk)
		addr = blk.Next
	}
	return blocks
}

func (p *HeapInfo) String() string {
	t, err := template.New("heap.dot").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	var bf bytes.Buffer
	err = t.Execute(&bf, p)
	if err != nil {
		panic(err)
	}
	return bf.String()
}
