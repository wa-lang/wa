// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

// 目前宿主函数是固定的

func (p *wat2laWorker) buildImport(w io.Writer) error {
	if len(p.m.Imports) == 0 {
		return nil
	}

	// 同一个对象可能被导入多次
	var hostFuncMap = make(map[string]bool)

	// 声明原始的宿主函数
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.FUNC {
			panic(fmt.Sprintf("ERR: import global %s.%s", importSpec.ObjModule, importSpec.ObjName))
		}

		fnName := importSpec.ObjModule + "." + importSpec.ObjName

		// 已经处理过
		if hostFuncMap[fnName] {
			continue
		}
		hostFuncMap[fnName] = true

		// 检查导入系统调用的函数签名
		p.checkSyscallSig(importSpec)
	}
	if err := p.buildSyscall(w, hostFuncMap); err != nil {
		return err
	}

	return nil
}
