// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"os"
	"path/filepath"
)

type SyncFile struct {
	RootDir string
}

func (p *SyncFile) SaveFile(path string, text string) {
	if p.RootDir == "" {
		return
	}

	abspath := filepath.Join(p.RootDir, path)
	os.MkdirAll(filepath.Dir(abspath), 0777)
	os.WriteFile(abspath, []byte(text), 0666)
}
