package arm

import (
	"wa-lang.org/wa/internal/p9asm/objabi"
)

var UnaryDst = map[objabi.As]bool{
	ASWI:  true,
	AWORD: true,
}
