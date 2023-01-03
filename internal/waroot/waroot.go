// 版权 @2021 凹语言 作者。保留所有权利。

package waroot

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed _waroot
var _warootFS embed.FS

func GetFS() fs.FS {
	// embed.FS 均采用 Unix 风格路径
	fs, err := fs.Sub(_warootFS, "_waroot/src")
	if err != nil {
		panic(err)
	}
	return fs
}

func IsStdPkg(pkgpath string) bool {
	for _, s := range stdPkgs {
		if s == pkgpath {
			return true
		}
		if strings.HasPrefix(pkgpath, s+"/") {
			return true
		}
	}
	for _, s := range wzStdPkgs {
		if s == pkgpath {
			return true
		}
		if strings.HasPrefix(pkgpath, s+"/") {
			return true
		}
	}
	return false
}

var stdPkgs = []string{
	"arduino",
	"fmt",
	"runtime",
	"syscall",
}

var wzStdPkgs = []string{
	"书",
}
