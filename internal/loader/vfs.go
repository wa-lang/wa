// 版权 @2021 凹语言 作者。保留所有权利。

package loader

import (
	"os"
	"path/filepath"
	"testing/fstest"

	"github.com/wa-lang/wa/internal/config"
	"github.com/wa-lang/wa/internal/logger"
	"github.com/wa-lang/wa/internal/waroot"
)

// 根据路径加载需要的 vfs 和 manifest
func loadProgramMeta(cfg *config.Config, appPath string) (
	vfs *config.PkgVFS,
	manifest *config.Manifest,
	err error,
) {
	logger.Tracef(&config.EnableTrace_loader, "cfg: %+v", cfg)
	logger.Tracef(&config.EnableTrace_loader, "appPath: %s", appPath)

	manifest, err = config.LoadManifest(nil, appPath)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, nil, err
	}

	logger.Tracef(&config.EnableTrace_loader, "manifest: %s", manifest.JSONString())

	vfs = new(config.PkgVFS)
	if vfs.App == nil {
		vfs.App = os.DirFS(filepath.Join(manifest.Root, "src"))
	}

	if vfs.Std == nil {
		if cfg.WaRoot != "" {
			vfs.Std = os.DirFS(filepath.Join(cfg.WaRoot, "src"))
		} else {
			vfs.Std = waroot.GetFS()
		}
	}
	if vfs.Vendor == nil {
		vfs.Vendor = os.DirFS(filepath.Join(manifest.Root, "vendor"))
		if vfs.Vendor == nil {
			vfs.Vendor = make(fstest.MapFS) // empty fs
		}
	}

	return
}
