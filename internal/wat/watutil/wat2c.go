// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import (
	"wa-lang.org/wa/internal/wat/watutil/wat2c"
)

func Wat2C(filename string, source []byte) (code, header []byte, err error) {
	return wat2c.Wat2C(filename, source)
}
