// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package format

import (
	"bytes"

	"wa-lang.org/wa/internal/parser/w2parser"
	"wa-lang.org/wa/internal/printer"
	"wa-lang.org/wa/internal/token"
)

var config_zh = printer.Config{Mode: printer.WaZh | printer.UseSpaces | printer.TabIndent, Tabwidth: 8}

func _SourceFile_wz(src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := w2parser.ParseFile(nil, fset, "prog.wz", src, parserMode)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = config_zh.Fprint(&buf, fset, file)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
