// 版权 @2024 凹语言 作者。保留所有权利。

package wat2c

import (
	"fmt"
	"strings"
	"unicode"
)

func assert(condition bool, args ...interface{}) {
	if !condition {
		if msg := fmt.Sprint(args...); msg != "" {
			panic(fmt.Sprintf("assert failed, %s", msg))
		} else {
			panic(fmt.Sprint("assert failed"))
		}
	}
}

func unreachable() {
	panic("unreachable")
}

func toCName(name string) string {
	if name == "" {
		return name
	}

	// pkg 路径可能以数字开头
	// if c := name[0]; c >= '0' && c <= '9' {
	// 	return name
	// }

	var sb strings.Builder
	for _, c := range ([]rune)(name) {
		switch {
		case c == '$':
			sb.WriteRune('_')
		case c == '.':
			sb.WriteRune('_')
		case '0' <= c && c <= '9':
			sb.WriteRune(c)
		case 'a' <= c && c <= 'z':
			sb.WriteRune(c)
		case 'A' <= c && c <= 'Z':
			sb.WriteRune(c)
		case unicode.IsLetter(c):
			sb.WriteRune(c)
		default:
			sb.WriteRune('_')
		}
	}

	return sb.String()
}
