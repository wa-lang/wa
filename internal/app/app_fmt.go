// 版权 @2023 凹语言 作者。保留所有权利。

package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"wa-lang.org/wa/internal/format"
)

func (p *App) Fmt(path string) error {
	if path == "" {
		path, _ = os.Getwd()
	}

	var waFileList []string
	if strings.HasSuffix(path, "...") {
		waFileList = getDirWaFileList(strings.TrimSuffix(path, "..."))
	}

	var changedFileList []string
	for _, s := range waFileList {
		changed, err := p.fmtFile(s)
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

func (p *App) fmtFile(path string) (changed bool, err error) {
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

func getDirWaFileList(dir string) []string {
	var waFileList []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".wa") {
			waFileList = append(waFileList, path)
		}
		return nil
	})
	sort.Strings(waFileList)
	return waFileList
}
