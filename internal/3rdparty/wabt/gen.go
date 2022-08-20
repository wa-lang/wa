// 版权 @2022 凹语言 作者。保留所有权利。

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	c_matches, err := filepath.Glob("./internal/wabt-1.0.29/src/*.c")
	if err != nil {
		panic(err)
	}
	cc_matches, err := filepath.Glob("./internal/wabt-1.0.29/src/*.cc")
	if err != nil {
		panic(err)
	}
	for i, s := range append(c_matches, cc_matches...) {
		fmt.Println(i, "gen", s)
		genCCFile(s)
	}
}

func genCCFile(s string) {
	baseName := "zz_" + strings.TrimPrefix(s, "internal/wabt-1.0.29/src/")
	src := fmt.Sprintf(srcFormat, s)
	err := os.WriteFile(baseName, []byte(src), 0666)
	if err != nil {
		panic(err)
	}
}

const srcFormat = `// 版权 @2022 凹语言 作者。保留所有权利。

// Auto generated, DONOT EDIT!!!

//go:build cgo
// +build cgo

#include %q
`
