// 版权 @2023 凹语言 作者。保留所有权利。

package appfmt

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"wa-lang.org/wa/internal/format"
)

func Fmt(path string) error {
	if path == "" {
		path = "."
	}

	var waFileList []string
	switch {
	case strings.HasSuffix(path, ".wa"):
		waFileList = append(waFileList, path)
	case strings.HasSuffix(path, ".wz"):
		waFileList = append(waFileList, path)
	case strings.HasSuffix(path, "..."):
		waFileList = getDirWaFileList(
			strings.TrimSuffix(path, "..."),
			true, ".wa", ".wz", // 包含子目录
		)
	default:
		// 不包含子目录
		waFileList = getDirWaFileList(
			path, false, ".wa", ".wz",
		)
	}

	var changedFileList []string
	for _, s := range waFileList {
		changed, err := fmtFile(s)
		if err != nil {
			return err
		}
		if changed {
			changedFileList = append(changedFileList, s)
		}
	}
	for _, s := range changedFileList {
		fmt.Println(s)
	}
	return nil
}

func fmtFile(path string) (changed bool, err error) {
	code, changed, err := format.File(nil, path, nil)
	if err != nil {
		return false, err
	}
	if changed {
		if err = os.WriteFile(path, code, 0666); err != nil {
			return false, err
		}
	}
	return true, nil
}

func getDirWaFileList(dir string, walkSubDir bool, extList ...string) []string {
	var waFileList []string
	if !walkSubDir {
		files, err := os.ReadDir(".")
		if err != nil {
			return nil
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			for _, ext := range extList {
				if strings.HasSuffix(file.Name(), ext) {
					waFileList = append(waFileList, filepath.Join(dir, file.Name()))
				}
			}
		}

		sort.Strings(waFileList)
		return waFileList
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, ext := range extList {
			if strings.HasSuffix(path, ext) {
				waFileList = append(waFileList, path)
				return nil
			}
		}
		return nil
	})
	sort.Strings(waFileList)
	return waFileList
}
