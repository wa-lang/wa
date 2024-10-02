package wat2c

import (
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/wat/token"
)

func toCType(typ token.Token) string {
	switch typ {
	case token.I32:
		return "int32_t"
	case token.I64:
		return "int64_t"
	case token.F32:
		return "float"
	case token.F64:
		return "double"
	default:
		return "void"
	}
}

// 生成C语言的标识符
func toCName(name string) string {
	var sb strings.Builder
	for _, c := range ([]rune)(name) {
		switch {
		case c == '_':
			sb.WriteRune(c)
		case 'a' <= c && c <= 'z':
			sb.WriteRune(c)
		case 'A' <= c && c <= 'Z':
			sb.WriteRune(c)
		case c == '$':
			// $ 是保留字符，转义为 _$$_
			sb.WriteRune('_')
			sb.WriteRune('$')
			sb.WriteRune('$')
			sb.WriteRune('_')
		default:
			sb.WriteRune('_')
			sb.WriteRune('$')
			sb.WriteString(strconv.FormatInt(int64(c), 16))
			sb.WriteRune('_')
		}
	}

	return sb.String()
}
