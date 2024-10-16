// 版权 @2024 凹语言 作者。保留所有权利。

package main

import (
	"wa-lang.org/wa/internal/3rdparty/mapx"
)

func main() {
	m := mapx.MakeMap()
	m.Update("three", 3)
	m.Update("one", 1)
	//m.Update("two", 2)

	m.Dump()
}
