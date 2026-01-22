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
	const IntSize = 8

	p.gasComment(w, "定义表格")
	p.gasSectionDataStart(w)
	p.gasGlobal(w, kTableAddrName)
	p.gasGlobal(w, kTableSizeName)
	p.gasGlobal(w, kTableMaxSizeName)

	if p.m.Table == nil {
		p.gasDefI64(w, kTableAddrName, 0)
		p.gasDefI64(w, kTableSizeName, 1)
		p.gasDefI64(w, kTableMaxSizeName, 1)
		fmt.Fprintln(w)
		return nil
	}

	if p.m.Table.Type != token.FUNCREF {
		return fmt.Errorf("unsupported table type: %s", p.m.Table.Type)
	}

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
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kTableMaxSizeName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kTableMaxSizeName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")
		fmt.Fprintf(w, "    slli.d    $a0, $t0, 3 # sizeof(i64) == 8\n")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kRuntimeMalloc)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kRuntimeMalloc)
		fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
		fmt.Fprintf(w, "    pcalau12i $t1, %%pc_hi20(%s)\n", kTableAddrName)
		fmt.Fprintf(w, "    addi.d    $t1, $t1, %%pc_lo12(%s)\n", kTableAddrName)
		fmt.Fprintf(w, "    st.d      $a0, $t1, 0\n")
		fmt.Fprintln(w)

		p.gasCommentInFunc(w, "表格填充 0xFF")
		fmt.Fprintf(w, "    addi.d    $a1, $zero, 0xFF # a1 = 0xFF\n")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kTableMaxSizeName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kTableMaxSizeName)
		fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")
		fmt.Fprintf(w, "    slli.d    $a2, $t0, 3 # sizeof(i64) == 8\n")
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kRuntimeMemset)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kRuntimeMemset)
		fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
		fmt.Fprintln(w)

		if len(p.m.Elem) > 0 {
			p.gasCommentInFunc(w, "初始化表格")
			fmt.Fprintln(w)

			p.gasCommentInFunc(w, "加载表格地址")
			fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kTableAddrName)
			fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kTableAddrName)
			fmt.Fprintf(w, "    ld.d      $t0, $t0, 0\n")

			for i, elem := range p.m.Elem {
				for j, fnIdxOrName := range elem.Values {
					fnIndex := p.findFuncIndex(fnIdxOrName)
					tableOffset := int32(int(elem.Offset) + j*IntSize)

					p.gasCommentInFunc(w, fmt.Sprintf("elem[%d]: table[%d+%d] = %s", i, elem.Offset, j, fnIdxOrName))

					// 生成偏移地址常量
					if x := tableOffset; x >= -2048 && x <= 2047 {
						fmt.Fprintf(w, "    addi.d    $t1, $zero, %d # offset\n", x)
					} else {
						hi20 := uint32(x) >> 12
						lo12 := uint32(x) & 0xFFF
						fmt.Fprintf(w, "    lu12i.w   $t1, 0x%X\n # offset", hi20)
						fmt.Fprintf(w, "    ori       $t1, $t1, 0x%X\n", lo12)
					}

					// 生成函数索引常量
					if x := fnIndex; x >= -2048 && x <= 2047 {
						fmt.Fprintf(w, "    addi.d    $t2, $zero, %d # func index\n", x)
					} else {
						hi20 := uint32(x) >> 12
						lo12 := uint32(x) & 0xFFF
						fmt.Fprintf(w, "    lu12i.w   $t2, 0x%X\n # func index", hi20)
						fmt.Fprintf(w, "    ori       $t2, $t2, 0x%X\n", lo12)
					}

					// 计算绝对偏移地址
					fmt.Fprintf(w, "    add.d     $t1, $t1, $t0 # offset\n")
					fmt.Fprintf(w, "    st.d      $t2, $t1, 0\n")
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
