// 版权 @2023 凹语言 作者。保留所有权利。

package waroot

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/version"
	"wa-lang.org/wa/internal/wabt"
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

// 本地根目录
var _warootDir string

//go:embed bin/*
//go:embed cplog/*
//go:embed docs/*
//go:embed examples/*
//go:embed misc/*
//go:embed src/*
//go:embed tests/*
//go:embed changelog.md
//go:embed CONTRIBUTORS
//go:embed hello.wa
//go:embed LICENSE
//go:embed logo.png
//go:embed README.md
//go:embed VERSION
var _warootFS embed.FS

//go:embed src/base.wat.ws
var baseWsFile_wat string

//go:embed src/base.import.js
var baseImportFile_js string

func init() {
	if s, _ := os.UserHomeDir(); s != "" {
		_warootDir = filepath.Join(s, "wa")
	}
}

// 获取本地的凹语言根目录路径
func GetWarootPath() string {
	return _warootDir
}

// Waroot 是否有效
func IsWarootValid() bool {
	d, err := os.ReadFile(filepath.Join(_warootDir, "VERSION"))
	if err != nil {
		return false
	}

	ver := string(bytes.TrimSpace(d))
	return ver == version.Version
}

// 初始化Waroot
func InitWarootDir() error {
	if IsWarootValid() {
		return nil
	}

	// 删除旧的waroot
	os.Rename(_warootDir, fmt.Sprintf("%s-%v-bak", _warootDir, time.Now().Format("20060102150405")))

	err := fs.WalkDir(_warootFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		data, err := fs.ReadFile(_warootFS, path)
		if err != nil {
			return err
		}

		dstpath := filepath.Join(_warootDir, path)
		os.MkdirAll(filepath.Dir(dstpath), 0777)

		if filepath.Base(path) == "_keep" {
			return nil
		}

		f, err := os.Create(dstpath)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.Write(data); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 获取当前 wa 命令所在目录
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// exe 已经在目录中
	if isInDir(_warootDir, exePath) {
		return nil
	}

	// 复制 bin 文件
	if exeData, err := os.ReadFile(exePath); err == nil {
		dstpath := filepath.Join(_warootDir, "bin", filepath.Base(exePath))
		os.WriteFile(dstpath, exeData, 0777)
	}

	// 复制 wa.wat2wasm.exe文件
	dstpath := filepath.Join(_warootDir, "bin", wabt.Wat2WasmName)
	os.WriteFile(dstpath, wabt.LoadWat2Wasm(), 0777)

	return nil
}

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

// InDir checks whether path is in the file tree rooted at dir.
// It checks only the lexical form of the file names.
// It does not consider symbolic links.
//
// Copied from go/src/cmd/go/internal/search/search.go.
func isInDir(dir, path string) bool {
	pv := strings.ToUpper(filepath.VolumeName(path))
	dv := strings.ToUpper(filepath.VolumeName(dir))
	path = path[len(pv):]
	dir = dir[len(dv):]
	switch {
	default:
		return false
	case pv != dv:
		return false
	case len(path) == len(dir):
		if path == dir {
			return true
		}
		return false
	case dir == "":
		return path != ""
	case len(path) > len(dir):
		if dir[len(dir)-1] == filepath.Separator {
			if path[:len(dir)] == dir {
				return path[len(dir):] != ""
			}
			return false
		}
		if path[len(dir)] == filepath.Separator && path[:len(dir)] == dir {
			if len(path) == len(dir)+1 {
				return true
			}
			return path[len(dir)+1:] != ""
		}
		return false
	}
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
