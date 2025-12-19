// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/wasm"
	watast "wa-lang.org/wa/internal/wat/ast"
	watparser "wa-lang.org/wa/internal/wat/parser"
	wattoken "wa-lang.org/wa/internal/wat/token"
)

// TODO: 自动化汇编指令地址回填, 需要一个手动构建指令序列的方式

const DebugMode = false

type Options struct {
	Ttext uint64 // 代码段开始地址
	Tdata uint64 // 数据段开始地址
}

// Wat程序转译到 龙芯汇编
func Wat2LA64(filename string, source []byte) (m *watast.Module, f *ast.File, err error) {
	return wat2la(filename, source)
}

func wat2la(filename string, source []byte) (m *watast.Module, f *ast.File, err error) {
	m, err = watparser.ParseModule(filename, source)
	if err != nil {
		return m, nil, err
	}

	worker := newWat2rvWorker(m)
	f, err = worker.BuildProgram()
	return
}

type wat2laWorker struct {
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

	dataSection []*ast.Global
	textSection []*ast.Func

	trace bool // 调试开关
}

type inlinedTypeIndex struct {
	section    wasm.SectionID
	idx        wasm.Index
	inlinedIdx wasm.Index
}

func newWat2rvWorker(mWat *watast.Module) *wat2laWorker {
	p := &wat2laWorker{m: mWat, trace: DebugMode}

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

func (p *wat2laWorker) BuildProgram() (f *ast.File, err error) {
	p.dataSection = p.dataSection[:0]
	p.textSection = p.textSection[:0]

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

	file := &ast.File{
		Globals: p.dataSection,
		Funcs:   p.textSection,
	}

	file.Doc = &ast.CommentGroup{
		TopLevel: true,
		List: []*ast.Comment{
			{Text: "# 由程序自动生成, 不要直接修改!!!"},
		},
	}

	return file, nil
}
