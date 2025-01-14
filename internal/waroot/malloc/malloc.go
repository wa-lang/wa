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
	DefaultMemoryPages    int32 = 1  // 内存页数
	DefaultMemoryPagesMax int32 = 10 // 内存最大页数

	DefaultStackPtr      int32 = 8 << 12  // 栈指针
	DefaultHeapBase      int32 = 10 << 12 // 实际需要编译根据静态数据大小计算
	DefaultHeapLFixedCap int32 = 100      // 固定大小的空闲链表最大长度
)

// 内部常量
const (
	kFreeListHeadSize = 5 * 8 // 全部空闲链表头大小
	kBlockHeadSize    = 2 * 8 // 块头大小
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

// 初始化获取空闲链表
func (p *Heap) FreeList(size int32) HeapBlock {
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
