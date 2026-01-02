// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package format

import (
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/parser/zparser"
	"wa-lang.org/wa/internal/native/token"
)

// 启动调试模式
func SetDebug() {
	zparser.DebugMode = true
}

// 格式化汇编代码(丢了注释)
func Format(cpu abi.CPUType, path string, src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := zparser.ParseFile(cpu, fset, path, src)
	if err != nil {
		return nil, err
	}

	// 重新打印
	return []byte(f.String()), err
}
