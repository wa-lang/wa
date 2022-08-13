// 版权 @2021 凹语言 作者。保留所有权利。

package loader

import (
	"io/fs"
	"os"
	pathpkg "path"
	"strings"

	"github.com/wa-lang/wa/internal/ast"
	"github.com/wa-lang/wa/internal/config"
	"github.com/wa-lang/wa/internal/logger"
	"github.com/wa-lang/wa/internal/parser"
	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/token"
	"github.com/wa-lang/wa/internal/types"
	"github.com/wa-lang/wa/internal/waroot"
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

func (p *_LoaderVFS) LoadProgram(cfg *config.Config, appPath string) (*Program, error) {
	logger.Tracef(&config.EnableTrace_loader, "cfg: %+v", cfg)
	logger.Tracef(&config.EnableTrace_loader, "appPath: %s", appPath)

	manifest, err := config.LoadManifestVFS(p.vfs, appPath)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	logger.Tracef(&config.EnableTrace_loader, "manifest: %s", manifest.JSONString())

	p.Cfg = cfg
	p.Manifest = manifest
	p.Fset = token.NewFileSet()

	// import "runtime"
	logger.Trace(&config.EnableTrace_loader, "import runtime")
	if _, err := p.Import("runtime"); err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	// import "main"
	// TODO: 触发递归导入?
	logger.Trace(&config.EnableTrace_loader, "import "+manifest.MainPkg)
	if _, err := p.Import(manifest.MainPkg); err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	// 转为 SSA
	p.SSAProgram = ssa.NewProgram(p.Fset, ssa.SanityCheckFunctions)

	for pkgpath, pkg := range p.Pkgs {
		logger.Tracef(&config.EnableTrace_loader, "build SSA; pkgpath: %v", pkgpath)

		if err := p.buildSSA(pkgpath); err != nil {
			return p.Program, err
		}

		if pkgpath == manifest.MainPkg {
			p.SSAMainPkg = pkg.SSAPkg
		}
	}

	logger.Tracef(&config.EnableTrace_loader, "return ok")
	return p.Program, nil
}

func (p *_LoaderVFS) buildSSA(pkgpath string) error {
	pkg := p.Pkgs[pkgpath]
	if pkg.SSAPkg != nil {
		return nil
	}

	for _, importPkg := range pkg.Pkg.Imports() {
		if p.Pkgs[importPkg.Path()].SSAPkg == nil {
			if err := p.buildSSA(importPkg.Path()); err != nil {
				return err
			}
		}
	}

	pkg.SSAPkg = p.SSAProgram.CreatePackage(pkg.Pkg, pkg.Files, pkg.Info, true)
	pkg.SSAPkg.Build()

	return nil
}

func (p *_LoaderVFS) Import(pkgpath string) (*types.Package, error) {
	logger.Tracef(&config.EnableTrace_loader, "pkgpath: %v", pkgpath)

	if pkg, ok := p.Pkgs[pkgpath]; ok {
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
			if pkgpath == p.Manifest.MainPkg {
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
	pkg.Pkg, err = conf.Check(pkgpath, p.Fset, pkg.Files, pkg.Info)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	logger.Tracef(&config.EnableTrace_loader, "save pkgpath: %v", pkgpath)

	p.Pkgs[pkgpath] = &pkg
	return pkg.Pkg, nil
}

func (p *_LoaderVFS) ParseDir(pkgpath string) ([]*ast.File, error) {
	logger.Tracef(&config.EnableTrace_loader, "pkgpath: %v", pkgpath)

	var (
		filenames []string
		datas     [][]byte
		err       error
	)

	switch {
	case p.isStdPkg(pkgpath):
		dirpath := "src/" + pkgpath
		logger.Tracef(&config.EnableTrace_loader, "isStdPkg; dirpath: %v", dirpath)

		filenames, datas, err = p.readDirFiles(p.getWaRootFS(), dirpath)
		if err != nil {
			return nil, err
		}
	case p.isSelfPkg(pkgpath):
		var dirpath string
		if pkgpath == p.Manifest.Pkg.Pkgpath {
			dirpath = "src"
		} else {
			relpkg := strings.TrimPrefix(pkgpath, p.Manifest.Pkg.Pkgpath)
			dirpath = pathpkg.Join("src", relpkg)
		}

		logger.Tracef(&config.EnableTrace_loader, "isSelfPkg; dirpath: %v", dirpath)

		filenames, datas, err = p.readDirFiles(os.DirFS(p.Manifest.Root), dirpath)
		if err != nil {
			return nil, err
		}

		logger.Trace(&config.EnableTrace_loader, "isSelfPkg; return ok")

	default: // vendor
		dirpath := "vendor/" + strings.TrimPrefix(pkgpath, p.Manifest.Pkg.Pkgpath)
		filenames, datas, err = p.readDirFiles(os.DirFS(p.Manifest.Root), dirpath)
		if err != nil {
			return nil, err
		}
	}

	logger.Tracef(&config.EnableTrace_loader, "filenames: %v", filenames)

	var files []*ast.File
	for i, filename := range filenames {
		f, err := parser.ParseFile(p.vfs, p.Fset, filename, datas[i], parser.AllErrors)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "filename: %v", filename)
			logger.Tracef(&config.EnableTrace_loader, "datas[i]: %s", datas[i])

			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (p *_LoaderVFS) isStdPkg(pkgpath string) bool {
	return waroot.IsStdPkg(pkgpath)
}

func (p *_LoaderVFS) isSelfPkg(pkgpath string) bool {
	if pkgpath == p.Manifest.Pkg.Pkgpath {
		return true
	}
	if strings.HasPrefix(pkgpath, p.Manifest.Pkg.Pkgpath+"/") {
		return true
	}
	return false
}

func (p *_LoaderVFS) readDirFiles(fileSystem fs.FS, path string) (filenames []string, datas [][]byte, err error) {
	logger.Tracef(&config.EnableTrace_loader, "path: %v", path)

	panic("TODO")
}

func (p *_LoaderVFS) getWaRootFS() fs.FS {
	if p.Cfg.WaRoot != "" {
		return os.DirFS(p.Cfg.WaRoot)
	}
	return waroot.GetFS()
}

func (p *_LoaderVFS) getSizes() types.Sizes {
	var zero config.StdSizes
	//types.StdSizes
	if p == nil || p.Cfg.WaSizes == zero {
		return types.SizesFor(p.Cfg.WaArch)
	} else {
		return &types.StdSizes{
			WordSize: p.Cfg.WaSizes.WordSize,
			MaxAlign: p.Cfg.WaSizes.MaxAlign,
		}
	}
}
