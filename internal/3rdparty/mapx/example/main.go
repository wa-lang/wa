// 版权 @2024 凹语言 作者。保留所有权利。

package main

import (
	"fmt"

	"wa-lang.org/wa/internal/3rdparty/mapx"
)

func main() {
	m := mapx.MakeMap()
	m.Update("three", 3)
	m.Update("one", 1)
	m.Update("two", 2)

	for it := mapx.MakeMapIter(m); ; {
		ok, k, v := it.Next()
		if !ok {
			break
		}
		fmt.Println(k, v)
	}

	fmt.Println("===")
	m.Update("three", 4)

	for it := mapx.MakeMapIter(m); ; {
		ok, k, v := it.Next()
		if !ok {
			break
		}
		fmt.Println(k, v)
	}

}
