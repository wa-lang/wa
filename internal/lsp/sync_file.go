// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

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
