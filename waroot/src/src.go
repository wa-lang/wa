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
	"apple",           // 测试已覆盖, wat2wasm ok
	"archive/txtar",   // API 完整, wat2wasm ok
	"bufio",           // API 完整, wat2wasm ok
	"bytes",           // API 完整, wat2wasm failed
	"compress/snappy", // ?
	"container/heap",  // ?
	"container/list",  // ?
	"container/ring",  // ?
	"crypto/md5",      // ?, 测试失败, Skip, wat2wasm ok
	"errors",          // API 完整, 测试已覆盖, wat2wasm ok
	"encoding",        // API 完整, wat2wasm ok
	"encoding/base32", // ?, wat2wasm failed
	"encoding/base64", // API 完整, wat2wasm failed
	"encoding/binary", // API 部分, wat2wasm failed
	"encoding/hex",    // API 完整, wat2wasm ok
	"encoding/pem",    // ?, wat2wasm failed
	"encoding/qrcode", // ?
	"debug",           // ?
	"fmt",             // ?, wat2wasm ok
	"gpu",             // ?
	"hash",            // API 完整, wat2wasm ok
	"hash/adler32",    // ?, wat2wasm ok
	"hash/crc32",      // API 完整, wat2wasm ok
	"hash/fnv",        // ?, wat2wasm failed
	"image",           // ?, wat2wasm ok
	"image/color",     // ?
	"io",              // API 部分, wat2wasm ok
	"js",              // ?
	"js/canvas",       // ?
	"js/p5",           // ?
	"math",            // API 部分
	"math/big",        // API 部分, wat2wasm failed
	"math/bits",       // API 完整, wat2wasm failed
	"math/gf256",      // ?, wat2wasm failed
	"net",             // ?
	"os",              // API 部分, wat2wasm ok
	"reflect",         // ?, wat2wasm failed
	"regexp",          // API 部分
	"runtime",         //
	"sort",            // API 完整, wat2wasm failed
	"strconv",         // API 完整, wat2wasm failed
	"strings",         // API 完整, wat2wasm ok
	"syscall",         // API 完整, wat2wasm failed
	"syscall/js",      //
	"syscall/wasi",    //
	"syscall/unknown", //
	"text/template",   // 无
	"unicode",         // API 部分
	"unicode/ctypes",  // API 完整, 测试已覆盖, wat2wasm ok
	"unicode/utf8",    // API 完整, 测试已覆盖, wat2wasm ok
}

var wzStdPkgs = []string{
	"书",
}
