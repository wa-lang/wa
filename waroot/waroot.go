// 版权 @2023 凹语言 作者。保留所有权利。

package waroot

import (
	"embed"
	"io/fs"
	"strings"

	"wa-lang.org/wa/internal/config"
)

//go:embed VERSION
var _VERSION string

//go:embed misc/_example_app
var _exampleAppFS embed.FS

//go:embed misc/_example_vendor
var _exampleVendorFS embed.FS

// 版本号(dev后缀表示开发版)
func GetVersion() string {
	return _VERSION
}

func GetExampleAppFS() fs.FS {
	fs, err := fs.Sub(_exampleAppFS, "misc/_example_app")
	if err != nil {
		panic(err)
	}
	return fs
}

func GetExampleVendorFS() fs.FS {
	fs, err := fs.Sub(_exampleVendorFS, "misc/_example_vendor")
	if err != nil {
		panic(err)
	}
	return fs
}

//go:embed src
var _warootFS embed.FS

//go:embed src/base.clang.ws
var baseWsFile_clang string

//go:embed src/base.llvm.ws
var baseWsFile_llvm string

//go:embed src/base.wat.ws
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
	fs, err := fs.Sub(_warootFS, "src")
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

func GetStdPkgList() []string {
	return append([]string{}, stdPkgs...)
}

var stdPkgs = []string{
	"apple",
	"archive/txtar",
	"binary",
	"bytes",
	"errors",
	"encoding",
	"fmt",
	"hash",
	"image",
	"image/bmp",
	"image/color",
	"image/color/palette",
	"io",
	"math",
	"math/bits",
	"os",
	"reflect",
	"regexp",
	"runtime",
	"sort",
	"strconv",
	"strings",
	"syscall",
	"text/template",
	"unicode",
	"unicode/utf8",
}

var wzStdPkgs = []string{
	"书",
}
