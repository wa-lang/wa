package tabwriter

import (
	"unicode"
	"unicode/utf8"
)

// 简化版：计算字符串显示宽度
func strDisplayWidth(b []byte) int {
	width := 0
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if unicode.Is(unicode.Han, r) {
			width += 2
		} else {
			width += 1
		}
		b = b[size:]
	}
	return width
}
