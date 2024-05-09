// 版权 @2023 凹语言 作者。保留所有权利。

// 用于内联的标准库文件

package src

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed base.wat.ws
var baseWsFile_wat string

//go:embed base.import.js
var baseImportFile_js string

//go:embed *
var _stdFS embed.FS

// 获取汇编基础代码
func GetBaseWsCode(backend string) string {
	switch backend {
	case WaBackend_wat:
		return baseWsFile_wat
	}
	for _, s := range WaBackend_List {
		if s == backend {
			return ""
		}
	}
	panic("unreachable")
}

// 获取宿主基础代码
func GetBaseImportCode(waos string) string {
	switch waos {
	case WaOS_js:
		return baseImportFile_js
	case WaOS_unknown:
		return ""
	case WaOS_wasi:
		return ""
	}
	for _, s := range WaBackend_List {
		if s == waos {
			return ""
		}
	}
	panic("unreachable")
}

func GetStdFS() fs.FS {
	return _stdFS
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

func GetStdTestPkgList() []string {
	var ss []string
	for _, s := range stdPkgs {
		if s == "unsafe" || s == "debug" {
			continue
		}
		if strings.HasPrefix(s, "syscall") {
			continue
		}
		if s == "js" || strings.HasPrefix(s, "js/") {
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
	"compress/snappy", // ?
	"container/heap",  // ?
	"container/list",  // ?
	"container/ring",  // ?
	"crypto/md5",      // ?, 测试失败, Skip
	"errors",          // API 完整, 测试已覆盖
	"encoding",        // API 完整
	"encoding/base32", // ?
	"encoding/base64", // API 完整
	"encoding/binary", // API 部分
	"encoding/hex",    // API 完整
	"encoding/pem",    // ?
	"encoding/qrcode", // ?
	"debug",           // ?
	"fmt",             // ?
	"gpu",             // ?
	"hash",            // API 完整
	"hash/adler32",    // ?
	"hash/crc32",      // API 完整
	"hash/fnv",        // ?
	"image",           // ?
	"image/color",     // ?
	"io",              // API 部分
	"js",              // ?
	"js/canvas",       // ?
	"js/p5",           // ?
	"math",            // API 部分
	"math/big",        // API 部分
	"math/bits",       // API 完整
	"math/gf256",      // ?
	"net",             // ?
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
