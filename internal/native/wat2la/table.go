// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

const (
	kTableReg         = "T7"
	kTableAddrName    = "$wat2la.table.addr"
	kTableSizeName    = "$wat2la.table.size"
	kTableMaxSizeName = "$wat2la.table.maxSize"

	kFuncIndexTableName = "$wat2la.func.index.table"
)

func (p *wat2laWorker) buildTable(w io.Writer) error {
	if p.m.Table == nil {
		return nil
	}
	if p.m.Table.Type != token.FUNCREF {
		return fmt.Errorf("unsupported table type: %s", p.m.Table.Type)
	}

	const IntSize = 8

	if p.m.Table.Name != "" {
		fmt.Fprintf(w, "# table %s\n", p.m.Table.Name)
	} else {
		fmt.Fprintf(w, "# table\n")
	}
	if max := p.m.Table.MaxSize; max > 0 {
		fmt.Fprintf(w, "global %s: i64 = %d\n", kTableSizeName, p.m.Table.Size)
		fmt.Fprintf(w, "global %s: i64 = %d\n", kTableMaxSizeName, p.m.Table.MaxSize)
		fmt.Fprintf(w, "global %s: %d = {\n", kTableAddrName, max*IntSize)
	} else {
		fmt.Fprintf(w, "global %s: i64 = %d\n", kTableSizeName, p.m.Table.Size)
		fmt.Fprintf(w, "global %s: i64 = %d\n", kTableMaxSizeName, p.m.Table.Size)
		fmt.Fprintf(w, "global %s: %d = {\n", kTableAddrName, p.m.Table.Size*IntSize)
	}
	{
		// 为了更真实地模拟, 表格保存的是函数的索引!
		for _, elem := range p.m.Elem {
			for i, fnIdxOrName := range elem.Values {
				fnIndex := p.findFuncIndex(fnIdxOrName)
				fmt.Fprintf(w, "    %d: %d,\n", int(elem.Offset)+i*IntSize, fnIndex)
			}
		}
	}
	fmt.Fprintf(w, "}\n")
	fmt.Fprintln(w)

	// 函数索引列表
	if max := p.m.Table.MaxSize; max > 0 {
		fmt.Fprintf(w, "global %s: %d = {\n", kFuncIndexTableName, max*IntSize)
	} else {
		fmt.Fprintf(w, "global %s: %d = {\n", kFuncIndexTableName, p.m.Table.Size*IntSize)
	}
	{
		// 索引从导入函数开始计算
		var importFuncCount int
		for i, x := range p.m.Imports {
			if x.ObjKind == token.FUNC {
				fmt.Fprintf(w, "    %d: %s,\n", i, x.FuncName)
				importFuncCount++
			}
		}
		for i, x := range p.m.Funcs {
			fmt.Fprintf(w, "    %d: %s,\n", importFuncCount+i, x.Name)
		}
	}
	fmt.Fprintf(w, "}\n")
	fmt.Fprintln(w)

	return nil
}
