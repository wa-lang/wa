// 版权 @2021 凹语言 作者。保留所有权利。

package loader

import (
	"io/fs"

	"github.com/wa-lang/wa/internal/config"
)

// 从 VFS 加载程序
func LoadProgramVFS(cfg *config.Config, appPath string, vfs fs.FS) (*Program, error) {
	panic("TODO")
}
