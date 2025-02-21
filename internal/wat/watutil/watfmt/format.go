// 版权 @2025 凹语言 作者。保留所有权利。

package watfmt

import (
	"bytes"

	"wa-lang.org/wa/internal/wat/parser"
	"wa-lang.org/wa/internal/wat/printer"
)

// 格式化Wat代码
func Format(path string, src []byte) ([]byte, error) {
	m, err := parser.ParseModule(path, src)
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
