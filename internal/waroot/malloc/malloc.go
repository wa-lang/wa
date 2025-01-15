// 版权 @2025 凹语言 作者。保留所有权利。

package malloc

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"sync"

	"wa-lang.org/wa/internal/3rdparty/wazero"
	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	"wa-lang.org/wa/internal/wat/watutil"
)

//go:embed malloc.wat
var malloc_wat string

// 默认值
const (
	KPageBytes = 64 << 10 // 一个内存页大小

	DefaultMemoryPages    int32 = 1  // 内存页数
	DefaultMemoryPagesMax int32 = 10 // 内存最大页数

	DefaultStackPtr      int32 = 8 << 12  // 栈指针
	DefaultHeapBase      int32 = 10 << 12 // 实际需要编译根据静态数据大小计算
	DefaultHeapLFixedCap int32 = 100      // 固定大小的空闲链表最大长度
)

// 内部常量
const (
	kBlockHeadSize    = 8                  // 块头大小
	kFreeListHeadSize = 5 * kBlockHeadSize // 全部空闲链表头大小
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

	wazeroOnce          sync.Once
	wazeroCtx           context.Context
	wazeroConf          wazero.ModuleConfig
	wazeroRuntime       wazero.Runtime
	wazeroCompileModule wazero.CompiledModule
	wazeroModule        api.Module
	wazeroInitErr       error

	fnMalloc api.Function
	fnFree   api.Function
}

// 内存块
type HeapBlock struct {
	Size int32
	Next int32
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
	p.init()
	return p
}

func (p *Heap) init() {
	// 1. 获取 wat 文本
	var buf bytes.Buffer
	buf.WriteString("(module $malloc\n")
	if err := template.Must(template.New("wat").Parse(malloc_wat)).Execute(&buf, p.cfg); err != nil {
		panic(err)
	}
	buf.WriteString("\n)")

	// 2. wat 转为 wasm 字节数组
	wasmBytes, err := watutil.Wat2Wasm("malloc.wat", buf.Bytes())
	if err != nil {
		panic(err)
	}

	// 3. 初始化 wazero 运行时
	p.wazeroCtx = context.Background()
	p.wazeroConf = wazero.NewModuleConfig().WithName("malloc.wat")

	p.wazeroRuntime = wazero.NewRuntime(p.wazeroCtx)
	p.wazeroCompileModule, err = p.wazeroRuntime.CompileModule(p.wazeroCtx, wasmBytes)
	if err != nil {
		p.wazeroInitErr = err
		panic(err)
	}

	p.wazeroModule, p.wazeroInitErr = p.wazeroRuntime.InstantiateModule(
		p.wazeroCtx, p.wazeroCompileModule, p.wazeroConf,
	)

	// 4. 导出函数
	p.fnMalloc = p.wazeroModule.ExportedFunction("wa_malloc")
	if p.fnMalloc == nil {
		err = fmt.Errorf("wazero: func wa_malloc not found")
		return
	}
	p.fnFree = p.wazeroModule.ExportedFunction("wa_free")
	if p.fnFree == nil {
		err = fmt.Errorf("wazero: func wa_free not found")
		return
	}
}

// 获取全局变量
func (p *Heap) Global__stack_ptr() int32 {
	return p.xGlobal("__stack_ptr")
}
func (p *Heap) Global__heap_base() int32 {
	return p.xGlobal("__heap_base")
}
func (p *Heap) Global__heap_ptr() int32 {
	return p.xGlobal("__heap_ptr")
}
func (p *Heap) Global__heap_top() int32 {
	return p.xGlobal("__heap_top")
}
func (p *Heap) Global__heap_l128_freep() int32 {
	return p.xGlobal("__heap_l128_freep")
}
func (p *Heap) Global__heap_lfixed_cap() int32 {
	return p.xGlobal("__heap_lfixed_cap")
}

func (p *Heap) xGlobal(name string) int32 {
	v := p.wazeroModule.ExportedGlobal(name).Get(context.Background())
	return int32(uint32(v))
}

// 去读内存
func (p *Heap) ReadMemoryI32(offset int32) int32 {
	v, _ := p.wazeroModule.Memory().ReadUint32Le(context.Background(), uint32(offset))
	return int32(v)
}

// 读取空闲链表头
func (p *Heap) ReadL24Header() HeapBlock {
	offset := p.Global__heap_base() + 8*0
	return p.ReadBlock(offset)
}
func (p *Heap) ReadL32Header() HeapBlock {
	offset := p.Global__heap_base() + 8*1
	return p.ReadBlock(offset)
}
func (p *Heap) ReadL46Header() HeapBlock {
	offset := p.Global__heap_base() + 8*2
	return p.ReadBlock(offset)
}
func (p *Heap) ReadL80Header() HeapBlock {
	offset := p.Global__heap_base() + 8*3
	return p.ReadBlock(offset)
}
func (p *Heap) ReadL128Header() HeapBlock {
	offset := p.Global__heap_base() + 8*4
	return p.ReadBlock(offset)
}

// 读取 HeapBlock 数据
func (p *Heap) ReadBlock(offset int32) HeapBlock {
	size := p.ReadMemoryI32(offset + 0)
	next := p.ReadMemoryI32(offset + 4)
	return HeapBlock{Size: size, Next: next}
}

// 初始化获取空闲链表
func (p *Heap) ReadFreeListHeader(z int32) HeapBlock {
	return HeapBlock{}
}

// 分配 size 字节的内存, 返回地址 8 字节对齐
func (p *Heap) Malloc(size int32) int32 {
	results, err := p.fnMalloc.Call(p.wazeroCtx, uint64(size))
	if err != nil {
		panic(err)
	}
	if len(results) != 1 {
		panic("unreachable")
	}
	return int32(results[0])
}

// 释放内存
func (p *Heap) Free(ptr int32) {
	results, err := p.fnFree.Call(p.wazeroCtx, uint64(ptr))
	if err != nil {
		panic(err)
	}
	if len(results) != 0 {
		panic("unreachable")
	}
}

// 打印统计信息
func (p *Heap) String() string {
	return "TODO"
}
