// 版权 @2021 凹语言 作者。保留所有权利。

package waroot

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed _waroot
var _warootFS embed.FS

func GetFS() fs.FS {
	fs, err := fs.Sub(_warootFS, filepath.Join("_waroot", "src"))
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
	return false
}

var stdPkgs = []string{
	"fmt",
	"runtime",
	"syscall",
}
