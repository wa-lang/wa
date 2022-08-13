// 版权 @2021 凹语言 作者。保留所有权利。

package loader

import (
	"fmt"
	"io/fs"

	"github.com/wa-lang/wa/internal/config"
)

type _LoaderVFS struct {
	vfs fs.FS
	*Program
}

func newLoaderVFS(vfs fs.FS) *_LoaderVFS {
	return &_LoaderVFS{
		vfs: vfs,
		Program: &Program{
			Pkgs: make(map[string]*Package),
		},
	}
}

func (p *_LoaderVFS) LoadProgram(cfg *config.Config, pkgPath string) (*Program, error) {
	return p.Program, fmt.Errorf("TODO")
}
