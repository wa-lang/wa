// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

const (
	kTableReg          = "T7"
	kTableAddrName     = "$Table.addr"
	kTableSizeName     = "$Table.size"
	kTableMaxSizeName  = "$Table.maxSize"
	kFuncIndexListName = "$Table.funcIndexList"
	kFuncInitFuncName  = "$Table.initFunc"
)

func (p *wat2X64Worker) buildTable(w io.Writer) error {
	if p.m.Table == nil {
		return nil
	}
	if p.m.Table.Type != token.FUNCREF {
		return fmt.Errorf("unsupported table type: %s", p.m.Table.Type)
	}

	const IntSize = 8

	p.gasComment(w, "定义表格")
	p.gasSectionDataStart(w)

	if max := p.m.Table.MaxSize; max > 0 {
		p.gasDefI64(w, kTableSizeName, int64(p.m.Table.Size))
		p.gasDefI64(w, kTableMaxSizeName, int64(p.m.Table.MaxSize))
		p.gasDefArray(w, kTableAddrName, max, IntSize)
		p.gasDefArray(w, kFuncIndexListName, len(p.m.Imports)+len(p.m.Funcs), IntSize)
		fmt.Fprintln(w)
	} else {
		p.gasDefI64(w, kTableSizeName, int64(p.m.Table.Size))
		p.gasDefI64(w, kTableMaxSizeName, int64(p.m.Table.MaxSize))
		p.gasDefArray(w, kTableAddrName, p.m.Table.Size, IntSize)
		p.gasDefArray(w, kFuncIndexListName, len(p.m.Imports)+len(p.m.Funcs), IntSize)
		fmt.Fprintln(w)
	}

	// 定义初始化函数
	{
		p.gasComment(w, "表格初始化函数")
		p.gasSectionTextStart(w)
		p.gasGlobal(w, kFuncInitFuncName)
		p.gasFuncStart(w, kFuncInitFuncName)

		p.gasCommentInFunc(w, "影子空间")
		fmt.Fprintln(w, "    sub rsp, 40")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "初始化全部函数索引列表")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "加载函数索引列表地址")
		fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kFuncIndexListName)

		// 索引从导入函数开始计算
		var importFuncCount int
		for i, x := range p.m.Imports {
			if x.ObjKind == token.FUNC {
				p.gasCommentInFunc(w, fmt.Sprintf("导入[%d] %s", i, x.FuncName))
				fmt.Fprintf(w, "    lea rcx, [rip + %s]\n", kFuncNamePrefix+x.FuncName)
				fmt.Fprintf(w, "    mov [rax + %d], rcx\n", importFuncCount*8)
				importFuncCount++
			}
		}
		for i, x := range p.m.Funcs {
			p.gasCommentInFunc(w, fmt.Sprintf("函数[%d] %s", i, x.Name))
			fmt.Fprintf(w, "    lea rcx, [rip + %s]\n", kFuncNamePrefix+x.Name)
			fmt.Fprintf(w, "    mov [rax + %d], rcx\n", (importFuncCount+i)*8)
			importFuncCount++
		}
		fmt.Fprintln(w)

		if len(p.m.Elem) > 0 {
			p.gasCommentInFunc(w, "初始化表格元素")
			fmt.Fprintln(w)

			p.gasCommentInFunc(w, "加载表格地址")
			fmt.Fprintf(w, "    lea rax, [rip + %s]\n", kTableAddrName)

			for _, elem := range p.m.Elem {
				for i, fnIdxOrName := range elem.Values {
					fnIndex := p.findFuncIndex(fnIdxOrName)
					p.gasCommentInFunc(w, fmt.Sprintf("表格[%d] = %s", i, fnIdxOrName))
					fmt.Fprintf(w, "    mov qword ptr [rax + %d], %d\n", int(elem.Offset)+i*IntSize, fnIndex)
				}
			}
			fmt.Fprintln(w)
		}

		p.gasCommentInFunc(w, "函数返回")
		fmt.Fprintln(w, "    add rsp, 40")
		fmt.Fprintln(w, "    ret")
		fmt.Fprintln(w)
	}

	return nil
}
