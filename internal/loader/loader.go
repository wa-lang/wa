// 版权 @2021 凹语言 作者。保留所有权利。

package loader

import (
	"io/fs"
	"os"
	pathpkg "path"
	"path/filepath"
	"sort"
	"strings"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/config"
	wzparser "wa-lang.org/wa/internal/frontend/wz/parser"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/ssa"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/types"
	"wa-lang.org/wa/internal/waroot"
)

type _Loader struct {
	cfg  config.Config
	vfs  config.PkgVFS
	prog *Program
}

func newLoader(cfg *config.Config) *_Loader {
	return &_Loader{
		cfg: *cfg.Clone(),
		prog: &Program{
			Pkgs: make(map[string]*Package),
		},
	}
}

func (p *_Loader) LoadProgramFile(filename string, src interface{}) (*Program, error) {
	logger.Tracef(&config.EnableTrace_loader, "cfg: %+v", p.cfg)
	logger.Tracef(&config.EnableTrace_loader, "filename: %v", filename)

	vfs, manifest, err := loadProgramFileMeta(&p.cfg, filename, src)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	return p.loadProgram(vfs, manifest)
}

func (p *_Loader) LoadProgram(appPath string) (*Program, error) {
	logger.Tracef(&config.EnableTrace_loader, "cfg: %+v", p.cfg)
	logger.Tracef(&config.EnableTrace_loader, "appPath: %s", appPath)

	if isWaFile(appPath) || isWzFile(appPath) {
		vfs, manifest, err := loadProgramFileMeta(&p.cfg, appPath, nil)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, err
		}

		return p.loadProgram(vfs, manifest)
	}

	vfs, manifest, err := loadProgramMeta(&p.cfg, appPath)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	return p.loadProgram(vfs, manifest)
}

func (p *_Loader) LoadProgramVFS(vfs *config.PkgVFS, appPath string) (*Program, error) {
	manifest, err := config.LoadManifest(vfs.App, appPath)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	return p.loadProgram(vfs, manifest)
}

// 加载程序
func (p *_Loader) loadProgram(vfs *config.PkgVFS, manifest *config.Manifest) (*Program, error) {
	logger.DumpFS(&config.EnableTrace_loader, "vfs.app", vfs.App, ".")
	logger.Tracef(&config.EnableTrace_loader, "manifest: %s", manifest.JSONString())

	p.vfs = *vfs
	p.prog.Cfg = &p.cfg
	p.prog.Manifest = manifest
	p.prog.Fset = token.NewFileSet()

	if p.vfs.Std == nil {
		if p.cfg.WaRoot != "" {
			p.vfs.Std = os.DirFS(filepath.Join(p.cfg.WaRoot, "src"))
		} else {
			p.vfs.Std = waroot.GetFS()
		}
	}

	// import "runtime"
	logger.Trace(&config.EnableTrace_loader, "import runtime")
	if _, err := p.Import("runtime"); err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	// import "main"
	logger.Trace(&config.EnableTrace_loader, "import "+manifest.MainPkg)
	if _, err := p.Import(manifest.MainPkg); err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	// 转为 SSA
	p.prog.SSAProgram = ssa.NewProgram(p.prog.Fset, ssa.SanityCheckFunctions)

	for pkgpath, pkg := range p.prog.Pkgs {
		logger.Tracef(&config.EnableTrace_loader, "build SSA; pkgpath: %v", pkgpath)

		if err := p.buildSSA(pkgpath); err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return p.prog, err
		}

		if pkgpath == manifest.MainPkg {
			p.prog.SSAMainPkg = pkg.SSAPkg
		}
	}

	logger.Trace(&config.EnableTrace_loader, "return ok")
	return p.prog, nil
}

func (p *_Loader) buildSSA(pkgpath string) error {
	pkg := p.prog.Pkgs[pkgpath]
	if pkg.SSAPkg != nil {
		return nil
	}

	for _, importPkg := range pkg.Pkg.Imports() {
		if p.prog.Pkgs[importPkg.Path()].SSAPkg == nil {
			if err := p.buildSSA(importPkg.Path()); err != nil {
				logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
				return err
			}
		}
	}

	pkg.SSAPkg = p.prog.SSAProgram.CreatePackage(pkg.Pkg, pkg.Files, pkg.Info, true)
	pkg.SSAPkg.Build()

	return nil
}

func (p *_Loader) Import(pkgpath string) (*types.Package, error) {
	logger.Tracef(&config.EnableTrace_loader, "pkgpath: %v", pkgpath)

	if pkg, ok := p.prog.Pkgs[pkgpath]; ok {
		return pkg.Pkg, nil
	}

	var err error
	var pkg Package

	// 解析当前包到 AST
	pkg.Files, err = p.ParseDir(pkgpath)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	// 修复 pkg 名称(wa 是可选)
	for _, f := range pkg.Files {
		if f.Name.Name == "" {
			if pkgpath == p.prog.Manifest.MainPkg {
				f.Name.Name = "main"
			} else {
				pkgname := pkgpath
				if idx := strings.LastIndex(pkgname, "/"); idx != -1 {
					pkgname = pkgname[idx+1:]
				}
				f.Name.Name = pkgname
			}
		}
	}

	pkg.Info = &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	conf := types.Config{
		Importer: p,
		Sizes:    p.getSizes(),
	}
	pkg.Pkg, err = conf.Check(pkgpath, p.prog.Fset, pkg.Files, pkg.Info)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	logger.Tracef(&config.EnableTrace_loader, "save pkgpath: %v", pkgpath)

	p.prog.Pkgs[pkgpath] = &pkg
	return pkg.Pkg, nil
}

func (p *_Loader) ParseDir(pkgpath string) ([]*ast.File, error) {
	logger.Tracef(&config.EnableTrace_loader, "pkgpath: %v", pkgpath)

	var (
		filenames []string
		datas     [][]byte
		err       error
	)

	switch {
	case p.isStdPkg(pkgpath):
		logger.Tracef(&config.EnableTrace_loader, "isStdPkg; pkgpath: %v", pkgpath)

		filenames, datas, err = p.readDirFiles(p.vfs.Std, pkgpath)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, err
		}
	case p.isSelfPkg(pkgpath):
		relpkg := strings.TrimPrefix(pkgpath, p.prog.Manifest.Pkg.Pkgpath)
		if relpkg == "" {
			relpkg = "."
		}

		logger.Tracef(&config.EnableTrace_loader, "isSelfPkg; pkgpath=%v, relpkg=%v", pkgpath, relpkg)

		filenames, datas, err = p.readDirFiles(p.vfs.App, relpkg)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, err
		}

		logger.Trace(&config.EnableTrace_loader, "isSelfPkg; return ok")

	default: // vendor
		logger.Tracef(&config.EnableTrace_loader, "vendorPkg; pkgpath: %v", pkgpath)

		filenames, datas, err = p.readDirFiles(p.vfs.Vendor, pkgpath)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, err
		}
	}

	logger.Tracef(&config.EnableTrace_loader, "filenames: %v", filenames)

	var files []*ast.File
	for i, filename := range filenames {
		var f *ast.File
		if p.hasExt(filename, ".wz") {
			f, err = wzparser.ParseFile(nil, p.prog.Fset, filename, datas[i], wzparser.AllErrors|wzparser.ParseComments)
		} else {
			f, err = parser.ParseFile(nil, p.prog.Fset, filename, datas[i], parser.AllErrors|parser.ParseComments)
		}
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "filename: %v", filename)
			logger.Tracef(&config.EnableTrace_loader, "datas[i]: %s", datas[i])
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)

			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (p *_Loader) readDirFiles(fileSystem fs.FS, path string) (filenames []string, datas [][]byte, err error) {
	path = filepath.ToSlash(path)
	path = strings.TrimPrefix(path, "/")

	logger.Tracef(&config.EnableTrace_loader, "path: %v", path)

	dirEntries, err := fs.ReadDir(fileSystem, path)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "path: %v", path)
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, nil, err
	}

	for i, entry := range dirEntries {
		logger.Tracef(&config.EnableTrace_loader, "%d: path=%v", i, entry.Name())

		if entry.IsDir() {
			continue
		}

		if p.isSkipedSouceFile(entry.Name()) {
			continue
		}

		filenames = append(filenames, entry.Name())
	}

	logger.Tracef(&config.EnableTrace_loader, "filenames=%v", filenames)

	sort.Strings(filenames)
	for _, name := range filenames {
		var fpath string
		if path != "" && path != "." {
			// embed.FS 采用 Unix 风格路径
			fpath = strings.TrimPrefix(pathpkg.Join(path, name), "/")
		} else {
			fpath = strings.TrimPrefix(name, "/")
		}

		data, err := fs.ReadFile(fileSystem, fpath)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "fpath: %v", fpath)
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, nil, err
		}
		datas = append(datas, data)
	}

	logger.Trace(&config.EnableTrace_loader, "return ok")

	return filenames, datas, nil
}

func (p *_Loader) hasExt(name string, extensions ...string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(name, ext) {
			return true
		}
	}
	return false
}

func (p *_Loader) isStdPkg(pkgpath string) bool {
	return waroot.IsStdPkg(pkgpath)
}

func (p *_Loader) isSelfPkg(pkgpath string) bool {
	if pkgpath == p.prog.Manifest.Pkg.Pkgpath {
		return true
	}
	if strings.HasPrefix(pkgpath, p.prog.Manifest.Pkg.Pkgpath+"/") {
		return true
	}
	return false
}

func (p *_Loader) getSizes() types.Sizes {
	var zero config.StdSizes
	if p == nil || p.cfg.WaSizes == zero {
		return types.SizesFor(p.cfg.WaArch)
	} else {
		return &types.StdSizes{
			WordSize: p.cfg.WaSizes.WordSize,
			MaxAlign: p.cfg.WaSizes.MaxAlign,
		}
	}
}

func (p *_Loader) isSkipedSouceFile(filename string) bool {
	if strings.HasPrefix(filename, "_") {
		return true
	}
	if !p.hasExt(filename, ".wa", ".wa.go", ".wz") {
		return true
	}

	if p.cfg.WaOS != "" {
		var isTargetFile bool
		for _, ext := range []string{".wa", ".wa.go", ".wz"} {
			for _, os := range []string{"wasi", "arduino", "chrome"} {
				if strings.HasSuffix(filename, "_"+os+ext) {
					isTargetFile = true
					break
				}
			}
		}
		if isTargetFile {
			var shouldSkip = true
			for _, ext := range []string{".wa", ".wa.go", ".wz"} {
				if strings.HasSuffix(filename, "_"+p.cfg.WaOS+ext) {
					shouldSkip = false
					break
				}
			}
			if shouldSkip {
				return true
			}
		}
	}

	return false
}
