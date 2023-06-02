// 版权 @2021 凹语言 作者。保留所有权利。

package waroot

import (
	"embed"
	"io/fs"
	"strings"

	"wa-lang.org/wa/internal/config"
)

//go:embed _waroot
var _warootFS embed.FS

//go:embed base.clang.ws
var baseWsFile_clang string

//go:embed base.llvm.ws
var baseWsFile_llvm string

//go:embed base.wat.ws
var baseWsFile_wat string

// 获取汇编基础代码
func GetBaseWsCode(backend string) string {
	switch backend {
	case config.WaBackend_clang:
		return baseWsFile_clang
	case config.WaBackend_llvm:
		return baseWsFile_llvm
	case config.WaBackend_wat:
		return baseWsFile_wat
	}
	for _, s := range config.WaBackend_List {
		if s == backend {
			return ""
		}
	}
	panic("unreachable")
}

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
	"errors",
	"fmt",
	"math",
	"os",
	"regexp",
	"runtime",
	"strconv",
	"syscall",
	"unicode/utf8",
}

var wzStdPkgs = []string{
	"书",
}
