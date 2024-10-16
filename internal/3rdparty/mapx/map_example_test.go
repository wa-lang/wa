// 版权 @2024 凹语言 作者。保留所有权利。

package mapx_test

import (
	"fmt"

	"wa-lang.org/wa/internal/3rdparty/mapx"
)

func Example() {
	m := mapx.MakeMap()
	m.Update(1, 2)
	m.Update(3, 4)
	m.Update(11, 12)
	m.Update(5, 6)
	m.Update(7, 8)
	m.Update(9, 10)

	for it := mapx.MakeMapIter(m); ; {
		ok, k, v := it.Next()
		if !ok {
			break
		}
		fmt.Println(k, v)
	}

	fmt.Println("===")

	m.Delete(100)
	m.Delete(7)

	for it := mapx.MakeMapIter(m); it.HasNext(); it.Next() {
		k, v := it.KeyValue()
		fmt.Println(k, v)
	}

	// Output:
	// 1 2
	// 3 4
	// 11 12
	// 5 6
	// 7 8
	// 9 10
	// ===
	// 1 2
	// 3 4
	// 11 12
	// 5 6
	// 9 10
}
