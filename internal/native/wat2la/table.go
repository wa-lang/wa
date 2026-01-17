// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

const (
	kTableInitFuncName = ".Table.initFunc"

	kTableAddrName    = ".Table.addr"
	kTableSizeName    = ".Table.size"
	kTableMaxSizeName = ".Table.maxSize"

	kTableFuncIndexListName       = ".Table.funcIndexList"
	kTableFuncIndexListElemPrefix = ".Table.funcIndexList."
	kTableFuncIndexListEndName    = ".Table.funcIndexList.end"
)

func (p *wat2laWorker) buildTable(w io.Writer) error {
	if p.m.Table == nil {
		return nil
	}

	if p.m.Table.Type != token.FUNCREF {
		return fmt.Errorf("unsupported table type: %s", p.m.Table.Type)
	}

	const IntSize = 8

	p.gasComment(w, "定义表格")
	p.gasSectionDataStart(w)
	p.gasGlobal(w, kTableAddrName)
	p.gasGlobal(w, kTableSizeName)
	p.gasGlobal(w, kTableMaxSizeName)

	if max := p.m.Table.MaxSize; max > 0 {
		p.gasDefI64(w, kTableAddrName, 0)
		p.gasDefI64(w, kTableSizeName, int64(p.m.Table.Size))
		p.gasDefI64(w, kTableMaxSizeName, int64(p.m.Table.MaxSize))
		fmt.Fprintln(w)
	} else {
		p.gasDefI64(w, kTableAddrName, 0)
		p.gasDefI64(w, kTableSizeName, int64(p.m.Table.Size))
		p.gasDefI64(w, kTableMaxSizeName, int64(p.m.Table.Size))
		fmt.Fprintln(w)
	}

	p.gasComment(w, "函数列表")
	p.gasComment(w, "保持连续并填充全部函数")
	p.gasSectionDataStart(w)
	p.gasFuncLabel(w, kTableFuncIndexListName)
	{
		// 索引从导入函数开始计算
		var importFuncCount int
		for _, x := range p.m.Imports {
			if x.ObjKind == token.FUNC {
				fmt.Fprintf(w, "%s%d: .quad %s\n",
					kTableFuncIndexListElemPrefix, importFuncCount,
					kImportNamePrefix+x.FuncName,
				)
				importFuncCount++
			}
		}
		for _, x := range p.m.Funcs {
			fmt.Fprintf(w, "%s%d: .quad %s\n",
				kTableFuncIndexListElemPrefix, importFuncCount,
				kFuncNamePrefix+x.Name,
			)
			importFuncCount++
		}
	}
	p.gasDefI64(w, kTableFuncIndexListEndName, 0)
	fmt.Fprintln(w)

	// 定义初始化函数
	{
		p.gasComment(w, "表格初始化函数")
		p.gasSectionTextStart(w)
		p.gasGlobal(w, kTableInitFuncName)
		p.gasFuncStart(w, kTableInitFuncName)

		fmt.Fprintf(w, "    addi.d  $sp, $sp, -16\n")
		fmt.Fprintf(w, "    st.d    $ra, $sp, 8\n")
		fmt.Fprintf(w, "    st.d    $fp, $sp, 0\n")
		fmt.Fprintf(w, "    addi.d  $fp, $sp, 0\n")
		fmt.Fprintf(w, "    addi.d  $sp, $sp, -32\n")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "分配表格")
		fmt.Fprintf(w, "    # mov  rcx, [rip + %s]\n", kTableMaxSizeName)
		fmt.Fprintf(w, "    # shl  rcx, 3 # sizeof(i64) == 8\n")
		fmt.Fprintf(w, "    # call %s\n", kRuntimeMalloc)
		fmt.Fprintf(w, "    # mov  [rip + %s], rax\n", kTableAddrName)
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "表格填充 0xFF")
		fmt.Fprintf(w, "    # mov  rcx, [rip + %s]\n", kTableAddrName)
		fmt.Fprintf(w, "    # mov  rdx, 0xFF\n")
		fmt.Fprintf(w, "    # mov  r8, [rip + %s]\n", kTableMaxSizeName)
		fmt.Fprintf(w, "    # shl  r8, 3 # sizeof(i64) == 8\n")
		fmt.Fprintf(w, "    # call %s\n", kRuntimeMemset)
		fmt.Fprintln(w)

		if len(p.m.Elem) > 0 {
			p.gasCommentInFunc(w, "初始化表格")
			fmt.Fprintln(w)

			p.gasCommentInFunc(w, "加载表格地址")
			fmt.Fprintf(w, "    # mov rax, [rip + %s]\n", kTableAddrName)

			for i, elem := range p.m.Elem {
				for j, fnIdxOrName := range elem.Values {
					fnIndex := p.findFuncIndex(fnIdxOrName)
					p.gasCommentInFunc(w, fmt.Sprintf("elem[%d]: table[%d+%d] = %s", i, elem.Offset, j, fnIdxOrName))
					fmt.Fprintf(w, "    # mov qword ptr [rax+%d], %d\n", int(elem.Offset)+j*IntSize, fnIndex)
				}
			}
			fmt.Fprintln(w)
		}

		fmt.Fprintf(w, "    # 函数返回\n")
		fmt.Fprintf(w, "    addi.d  $sp, $fp, 0\n")
		fmt.Fprintf(w, "    ld.d    $ra, $sp, 8\n")
		fmt.Fprintf(w, "    ld.d    $fp, $sp, 0\n")
		fmt.Fprintf(w, "    addi.d  $sp, $sp, 16\n")
		fmt.Fprintf(w, "    jirl    $zero, $ra, 0\n")
		fmt.Fprintln(w)
	}

	return nil
}
