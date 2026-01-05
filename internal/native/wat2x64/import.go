// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

const (
	kImportNamePrefix = "$Import."
)

func (p *wat2X64Worker) buildImport(w io.Writer) error {
	if len(p.m.Imports) == 0 {
		return nil
	}

	if len(p.m.Imports) > 0 {
		p.gasComment(w, "导入函数(由导入文件定义)")
		defer fmt.Fprintln(w)
	}

	// 声明原始的宿主函数
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.FUNC {
			panic(fmt.Sprintf("ERR: import global %s.%s", importSpec.ObjModule, importSpec.ObjName))
		}

		// 检查导入系统调用的函数签名
		p.checkSyscallSig(importSpec)

		// 导入函数有个名字修饰, 避免重名
		// 导入函数属于外部库, 需要通过外部文件单独定义
		absImportName := kImportNamePrefix + importSpec.ObjModule + "." + importSpec.ObjName
		p.gasExtern(w, absImportName)
		p.gasSet(w, kFuncNamePrefix+importSpec.FuncName, absImportName)
	}

	return nil
}
