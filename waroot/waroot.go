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

//go:embed src
var _warootFS embed.FS

//go:embed src/base.clang.ws
var baseWsFile_clang string

//go:embed src/base.wat.ws
var baseWsFile_wat string

//go:embed src/base.import.js
var baseImportFile_js string

// 获取汇编基础代码
func GetBaseWsCode(backend string) string {
	switch backend {
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

// 获取宿主基础代码
func GetBaseImportCode(waos string) string {
	switch waos {
	case config.WaOS_js:
		return baseImportFile_js
	case config.WaOS_unknown:
		return ""
	case config.WaOS_wasi:
		return ""
	}
	for _, s := range config.WaBackend_List {
		if s == waos {
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
	var ss []string
	for _, s := range stdPkgs {
		ss = append(ss, s)
	}
	return ss
}

func GetStdTestPkgList() []string {
	var ss []string
	for _, s := range stdPkgs {
		if strings.HasPrefix(s, "syscall") {
			continue
		}
		ss = append(ss, s)
	}
	return ss
}

var stdPkgs = []string{
	"apple",           // 测试已覆盖
	"archive/txtar",   // API 完整
	"bufio",           // API 完整
	"bytes",           // API 完整
	"errors",          // API 完整, 测试已覆盖
	"encoding",        // API 完整
	"encoding/base64", // API 完整
	"encoding/binary", // API 部分
	"encoding/hex",    // API 完整
	"fmt",             // ?
	"hash",            // API 完整
	"hash/crc32",      // API 完整
	"image",           // ?
	"image/bmp",       // ?
	"image/color",     // ?
	"io",              // API 部分
	"math",            // API 部分
	"math/big",        // API 部分
	"math/bits",       // API 完整
	"os",              // API 部分
	"reflect",         // ?
	"regexp",          // API 部分
	"runtime",         //
	"sort",            // API 完整
	"strconv",         // API 完整
	"strings",         // API 完整
	"syscall",         // API 完整
	"syscall/js",      // ?
	"syscall/wasi",    // ?
	"syscall/unknown", // ?
	"text/template",   // 无
	"unicode",         // API 部分
	"unicode/ctypes",  // API 完整, 测试已覆盖
	"unicode/utf8",    // API 完整, 测试已覆盖
}

var wzStdPkgs = []string{
	"书",
}
