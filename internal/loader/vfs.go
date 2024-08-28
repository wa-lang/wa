// 版权 @2021 凹语言 作者。保留所有权利。

package loader

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing/fstest"

	"wa-lang.org/wa/internal/ast/astutil"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/token"
	wasrc "wa-lang.org/wa/waroot/src"
)

// 读取 embed 列表, 提前加载内嵌资源数据
func parseEmbedPathList(filename string, src interface{}) []string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(nil, fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil
	}

	var ss []string
	for _, doc := range f.Comments {
		if info := astutil.ParseCommentInfo(doc); info.Embed != "" {
			ss = append(ss, info.Embed)
		}
	}

	return ss
}

// 根据路径加载需要的 vfs 和 manifest
func loadProgramFileMeta(cfg *config.Config, filename string, src interface{}) (
	vfs *config.PkgVFS,
	manifest *config.Manifest,
	err error,
) {
	logger.Tracef(&config.EnableTrace_loader, "cfg: %+v", cfg)
	logger.Tracef(&config.EnableTrace_loader, "filename: %s", filename)

	// 读取代码
	srcData, err := readSource(filename, src)
	if err != nil {
		return nil, nil, err
	}

	// 尝试加载本地的 manifest
	manifest, err = config.LoadManifest(nil, filepath.Dir(filename))
	if err != nil {
		err = nil // 忽略错误
	}

	// 重新构造 manifest
	if manifest == nil {
		manifest = &config.Manifest{
			Root:    "__main__",
			MainPkg: "__main__",
			Pkg: config.Manifest_package{
				Name:    filepath.Base(filename),
				Pkgpath: "__main__",
			},
		}
	}

	if cfg.Target != "" {
		manifest.Pkg.Target = cfg.Target
	}
	if manifest.Pkg.Target == "" {
		manifest.Pkg.Target = config.WaOS_Default
	}

	if cfg.Target != manifest.Pkg.Target {
		cfg.Target = manifest.Pkg.Target
	}

	logger.Tracef(&config.EnableTrace_loader, "manifest: %s", manifest.JSONString())

	// 构造入口文件
	vfs = new(config.PkgVFS)
	if vfs.App == nil {

		mapFS := fstest.MapFS{
			filepath.Base(filename): &fstest.MapFile{
				Data: srcData,
			},
		}

		// read embed list, and read file data
		embedList := parseEmbedPathList(filename, string(srcData))
		for _, name := range embedList {
			localpath := filepath.Join(filepath.Dir(filename), name)
			if data, err := os.ReadFile(localpath); err == nil {
				mapFS[path.Join(manifest.MainPkg, name)] = &fstest.MapFile{
					Data: data,
				}
			}
		}

		vfs.App = mapFS
	}

	if vfs.Std == nil {
		// pkg/std
		stdPath := filepath.Join(manifest.Root, "pkg", "std")
		if dirPathExists(stdPath) {
			vfs.Std = os.DirFS(stdPath)
		} else {
			vfs.Std = wasrc.GetStdFS()
		}
	}
	if vfs.Vendor == nil {
		if src == nil {
			vendorPath := filepath.Join(manifest.Root, "vendor")
			pkgVendorPath := filepath.Join(manifest.Root, "pkg", "vendor")

			if dirPathExists(vendorPath) {
				vfs.Vendor = os.DirFS(vendorPath)
			} else if dirPathExists(pkgVendorPath) {
				vfs.Vendor = os.DirFS(pkgVendorPath)
			}
		}
		if vfs.Vendor == nil {
			vfs.Vendor = make(fstest.MapFS) // empty fs
		}
	}

	return
}

// 根据路径加载需要的 vfs 和 manifest
func loadProgramMeta(cfg *config.Config, appPath string) (
	vfs *config.PkgVFS,
	manifest *config.Manifest,
	err error,
) {
	logger.Tracef(&config.EnableTrace_loader, "cfg: %+v", cfg)
	logger.Tracef(&config.EnableTrace_loader, "appPath: %s", appPath)

	if wasrc.IsStdPkg(appPath) {
		manifest = &config.Manifest{
			Root:    "",
			MainPkg: appPath,
			IsStd:   true,
			Pkg: config.Manifest_package{
				Name:    appPath[strings.LastIndex(appPath, "/")+1:],
				Pkgpath: appPath,
			},
		}

		vfs = new(config.PkgVFS)

		// pkg/std
		stdPath := filepath.Join(manifest.Root, "pkg", "std")
		if dirPathExists(stdPath) {
			vfs.Std = os.DirFS(stdPath)
		} else {
			vfs.Std = wasrc.GetStdFS()
		}

		return
	}

	manifest, err = config.LoadManifest(nil, appPath)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, nil, err
	}

	if cfg.Target != "" {
		manifest.Pkg.Target = cfg.Target
	}
	if manifest.Pkg.Target == "" {
		manifest.Pkg.Target = config.WaOS_Default
	}

	if cfg.Target != manifest.Pkg.Target {
		cfg.Target = manifest.Pkg.Target
	}

	logger.Tracef(&config.EnableTrace_loader, "manifest: %s", manifest.JSONString())

	vfs = new(config.PkgVFS)
	if vfs.App == nil {
		vfs.App = os.DirFS(filepath.Join(manifest.Root, "src"))
	}

	if vfs.Std == nil {
		// pkg/std
		stdPath := filepath.Join(manifest.Root, "pkg", "std")
		if dirPathExists(stdPath) {
			vfs.Std = os.DirFS(stdPath)
		} else {
			vfs.Std = wasrc.GetStdFS()
		}
	}
	if vfs.Vendor == nil {
		vendorPath := filepath.Join(manifest.Root, "vendor")
		pkgVendorPath := filepath.Join(manifest.Root, "pkg", "vendor")

		if dirPathExists(vendorPath) {
			vfs.Vendor = os.DirFS(vendorPath)
		} else if dirPathExists(pkgVendorPath) {
			vfs.Vendor = os.DirFS(pkgVendorPath)
		}

		if vfs.Vendor == nil {
			vfs.Vendor = make(fstest.MapFS) // empty fs
		}
	}

	return
}

func readSource(filename string, src interface{}) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			d, err := io.ReadAll(s)
			return d, err
		}
		return nil, errors.New("invalid source")
	}

	d, err := os.ReadFile(filename)
	return d, err
}

func isWaFile(path string) bool {
	if fi, err := os.Lstat(path); err == nil && fi.Mode().IsRegular() {
		return strings.HasSuffix(strings.ToLower(path), ".wa")
	}
	return false
}

func isWzFile(path string) bool {
	if fi, err := os.Lstat(path); err == nil && fi.Mode().IsRegular() {
		return strings.HasSuffix(strings.ToLower(path), ".wz")
	}
	return false
}
