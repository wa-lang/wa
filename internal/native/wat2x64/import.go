// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

const (
	kImportNamePrefix = ".Import."
)

func (p *wat2X64Worker) buildImport(w io.Writer) error {
	if len(p.m.Imports) == 0 {
		return nil
	}

	if len(p.m.Imports) > 0 {
		p.gasComment(w, "导入函数(由导入文件定义)")
		defer fmt.Fprintln(w)
	}

	// 是否已经导入过
	seenMap := make(map[string]bool)

	// 声明原始的宿主函数
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.FUNC {
			panic(fmt.Sprintf("ERR: import %s.%s", importSpec.ObjModule, importSpec.ObjName))
		}

		if importSpec.ObjModule != "syscall" {
			panic(fmt.Sprintf("ERR: import %s.%s", importSpec.ObjModule, importSpec.ObjName))
		}

		absName := importSpec.ObjModule + "." + importSpec.ObjName
		if seenMap[absName] {
			continue
		}
		seenMap[absName] = true
		p.gasExtern(w, importSpec.ObjName)
	}

	// 定义导入函数的别名
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.FUNC {
			panic(fmt.Sprintf("ERR: import global %s.%s", importSpec.ObjModule, importSpec.ObjName))
		}

		absName := importSpec.ObjModule + "." + importSpec.ObjName
		if !seenMap[absName] {
			continue
		}

		p.gasSet(w, kImportNamePrefix+absName, importSpec.ObjName)
	}

	return nil
}
