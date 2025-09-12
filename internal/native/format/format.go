// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package format

import (
	"bytes"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/parser"
	"wa-lang.org/wa/internal/native/printer"
	"wa-lang.org/wa/internal/native/token"
)

// 启动调试模式
func SetDebug() {
	parser.DebugMode = true
}

// 格式化汇编代码(丢了注释)
func Format(cpu abi.CPUType, path string, src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	m, err := parser.ParseFile(cpu, fset, path, src)
	if err != nil {
		return nil, err
	}

	// 重新打印
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, m); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
