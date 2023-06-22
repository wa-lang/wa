// 版权 @2019 凹语言 作者。保留所有权利。

package app

import (
	"fmt"
	"os"
	"path/filepath"

	"wa-lang.org/wa/internal/config"
)

func findPkgInfo(workDir string) (pkgroot, pkgpath string, err error) {
	var wd = workDir
	if wd == "" {
		if s, _ := os.Getwd(); s != "" {
			wd = s
		}
	}
	if abs, _ := filepath.Abs(wd); abs != "" {
		wd = abs
	}
	if wd == "" {
		return "", "", fmt.Errorf("pkg root not found")
	}

	pkgroot = wd
	for pkgroot != "" {
		waJsonPath := filepath.Join(pkgroot, config.WaModFile)
		if fi, _ := os.Stat(waJsonPath); fi != nil {
			pkgpath, err = filepath.Rel(pkgroot, wd)
			pkgroot = filepath.ToSlash(pkgroot)
			pkgpath = filepath.ToSlash(pkgpath)
			return // OK
		}
		pkgroot = filepath.Dir(pkgroot)
		if pkgroot == "" || pkgroot == "/" || pkgroot == filepath.Dir(pkgroot) {
			break
		}
	}

	return "", "", fmt.Errorf("pkg root not found")
}
