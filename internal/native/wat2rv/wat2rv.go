// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2rv

import (
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/wasm"
	watast "wa-lang.org/wa/internal/wat/ast"
	watparser "wa-lang.org/wa/internal/wat/parser"
	wattoken "wa-lang.org/wa/internal/wat/token"
)

const DebugMode = false

type Options struct {
	Ttext uint64 // 代码段开始地址
	Tdata uint64 // 数据段开始地址
}

// Wat程序转译到 RISCV32
func Wat2rv32(filename string, source []byte, opt Options) (m *watast.Module, f *ast.File, err error) {
	return wat2rv(filename, source, opt, 32)
}

// Wat程序转译到 RISCV64
func Wat2rv64(filename string, source []byte, opt Options) (m *watast.Module, f *ast.File, err error) {
	return wat2rv(filename, source, opt, 64)
}

func wat2rv(filename string, source []byte, opt Options, xlen int) (m *watast.Module, f *ast.File, err error) {
	m, err = watparser.ParseModule(filename, source)
	if err != nil {
		return m, nil, err
	}

	worker := newWat2rvWorker(m, opt)
	f, err = worker.BuildProgram()
	return
}

type wat2rvWorker struct {
	opt Options

	m *watast.Module

	importGlobalCount int // 导入全局只读变量的数目
	importFuncCount   int // 导入函数的数目

	inlinedTypeIndices []*inlinedTypeIndex
	inlinedTypes       []*wasm.FunctionType

	localNames      []string           // 参数和局部变量名
	localTypes      []wattoken.Token   // 参数和局部变量类型
	scopeLabels     []string           // 嵌套的label查询, if/block/loop
	scopeStackBases []int              // if/block/loop, 开始的栈位置
	scopeResults    [][]wattoken.Token // 对应块的返回值数量和类型

	file *ast.File // 当前输出文件
	fn   *ast.Func // 当前编译的函数

	trace bool // 调试开关
}

type inlinedTypeIndex struct {
	section    wasm.SectionID
	idx        wasm.Index
	inlinedIdx wasm.Index
}

func newWat2rvWorker(mWat *watast.Module, opt Options) *wat2rvWorker {
	p := &wat2rvWorker{m: mWat, opt: opt, trace: DebugMode}

	// 统计导入的global和func索引
	p.importGlobalCount = 0
	p.importFuncCount = 0
	for _, importSpec := range p.m.Imports {
		switch importSpec.ObjKind {
		case wattoken.GLOBAL:
			p.importGlobalCount++
		case wattoken.FUNC:
			p.importFuncCount++
		}
	}

	// 如果 start 字段为空, 则尝试用 _start 导出函数替代
	if p.m.Start == "" {
		for _, fn := range p.m.Funcs {
			if fn.ExportName == "_start" {
				p.m.Start = fn.Name
			}
		}
	}

	return p
}

func (p *wat2rvWorker) BuildProgram() (f *ast.File, err error) {
	if err := p.buildImport(); err != nil {
		return nil, err
	}

	if err := p.buildMemory(); err != nil {
		return nil, err
	}
	if err := p.buildTable(); err != nil {
		return nil, err
	}

	if err := p.buildGlobal(); err != nil {
		return nil, err
	}
	if err := p.buildFuncs(); err != nil {
		return nil, err
	}

	if err := p.buildTable_elem(); err != nil {
		return nil, err
	}
	if err := p.buildMemory_data(); err != nil {
		return nil, err
	}

	return p.file, nil
}

func (p *wat2rvWorker) buildImport() error {
	panic("TODO")
}
func (p *wat2rvWorker) buildMemory() error {
	panic("TODO")
}
func (p *wat2rvWorker) buildTable() error {
	panic("TODO")
}

func (p *wat2rvWorker) buildGlobal() error {
	panic("TODO")
}

func (p *wat2rvWorker) buildFuncs() error {
	for _, f := range p.m.Funcs {
		if err := p.buildFunc(f); err != nil {
			return err
		}
	}
	panic("TODO")
}

func (p *wat2rvWorker) buildTable_elem() error {
	panic("TODO")
}

func (p *wat2rvWorker) buildMemory_data() error {
	panic("TODO")
}
