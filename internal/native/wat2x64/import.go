// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

const (
	kImportNamePrefix = ".Wa.Import."
)

func (p *wat2X64Worker) buildImport(w io.Writer) error {
	if len(p.m.Imports) == 0 {
		return nil
	}

	if len(p.m.Imports) > 0 {
		p.gasComment(w, "导入函数(外部库定义)")
		defer fmt.Fprintln(w)
	}

	// 是否已经导入过
	seenMap := make(map[string]bool)

	// 声明导入函数的别名
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.FUNC {
			panic(fmt.Sprintf("ERR: import global %s.%s", importSpec.ObjModule, importSpec.ObjName))
		}

		absName := kImportNamePrefix + importSpec.ObjModule + "." + importSpec.ObjName
		if seenMap[absName] {
			continue
		}

		seenMap[absName] = true
		p.gasExtern(w, absName)
	}

	return nil
}
