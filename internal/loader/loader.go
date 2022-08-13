// 版权 @2021 凹语言 作者。保留所有权利。

package loader

import (
	"io/fs"
	"os"
	pathpkg "path"
	"sort"
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

type _Loader struct {
	cfg  *config.Config
	prog *Program

	appFs    fs.FS
	stdFs    fs.FS
	vednorFs fs.FS
}

func newLoader(cfg *config.Config) *_Loader {
	p := &_Loader{
		cfg: cfg.Clone(),
		prog: &Program{
			Pkgs: make(map[string]*Package),
		},
	}

	return p
}

// 加载程序
func (p *_Loader) LoadProgram(appPath string) (*Program, error) {
	logger.Tracef(&config.EnableTrace_loader, "cfg: %+v", p.cfg)
	logger.Tracef(&config.EnableTrace_loader, "appPath: %s", appPath)

	manifest, err := config.LoadManifest(nil, appPath)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	logger.Tracef(&config.EnableTrace_loader, "manifest: %s", manifest.JSONString())

	p.prog.Cfg = p.cfg
	p.prog.Manifest = manifest
	p.prog.Fset = token.NewFileSet()

	if p.cfg.VFS != nil {
		var err error
		p.appFs, err = fs.Sub(p.cfg.VFS, p.prog.Manifest.Root)
		if err != nil {
			return nil, err
		}
	} else {
		p.appFs = os.DirFS(p.prog.Manifest.Root)
	}

	if p.cfg.WaRoot != "" {
		p.stdFs = os.DirFS(p.cfg.WaRoot)
	} else {
		p.stdFs = waroot.GetFS()
	}

	p.vednorFs, err = fs.Sub(p.appFs, "vendor")
	if err != nil {
		return nil, err
	}

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
	p.prog.SSAProgram = ssa.NewProgram(p.prog.Fset, ssa.SanityCheckFunctions)

	for pkgpath, pkg := range p.prog.Pkgs {
		logger.Tracef(&config.EnableTrace_loader, "build SSA; pkgpath: %v", pkgpath)

		if err := p.buildSSA(pkgpath); err != nil {
			return p.prog, err
		}

		if pkgpath == manifest.MainPkg {
			p.prog.SSAMainPkg = pkg.SSAPkg
		}
	}

	logger.Tracef(&config.EnableTrace_loader, "return ok")
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
		dirpath := "src/" + pkgpath
		logger.Tracef(&config.EnableTrace_loader, "isStdPkg; dirpath: %v", dirpath)

		filenames, datas, err = p.readDirFiles(p.stdFs, dirpath)
		if err != nil {
			return nil, err
		}
	case p.isSelfPkg(pkgpath):
		var dirpath string
		if pkgpath == p.prog.Manifest.Pkg.Pkgpath {
			dirpath = "src"
		} else {
			relpkg := strings.TrimPrefix(pkgpath, p.prog.Manifest.Pkg.Pkgpath)
			dirpath = pathpkg.Join("src", relpkg)
		}

		logger.Tracef(&config.EnableTrace_loader, "isSelfPkg; dirpath: %v", dirpath)

		filenames, datas, err = p.readDirFiles(p.appFs, dirpath)
		if err != nil {
			return nil, err
		}

		logger.Trace(&config.EnableTrace_loader, "isSelfPkg; return ok")

	default: // vendor
		dirpath := strings.TrimPrefix(pkgpath, p.prog.Manifest.Pkg.Pkgpath)
		filenames, datas, err = p.readDirFiles(p.vednorFs, dirpath)
		if err != nil {
			return nil, err
		}
	}

	logger.Tracef(&config.EnableTrace_loader, "filenames: %v", filenames)

	var files []*ast.File
	for i, filename := range filenames {
		f, err := parser.ParseFile(nil, p.prog.Fset, filename, datas[i], parser.AllErrors)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "filename: %v", filename)
			logger.Tracef(&config.EnableTrace_loader, "datas[i]: %s", datas[i])

			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (p *_Loader) readDirFiles(fileSystem fs.FS, path string) (filenames []string, datas [][]byte, err error) {
	logger.Tracef(&config.EnableTrace_loader, "path: %v", path)

	dirEntries, err := fs.ReadDir(fileSystem, path)
	if err != nil {
		return nil, nil, err
	}

	for i, entry := range dirEntries {
		logger.Tracef(&config.EnableTrace_loader, "%d: path=%v", i, entry.Name())

		if entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		if !p.hasExt(entry.Name(), ".go", ".ugo", ".wa") {
			continue
		}

		filenames = append(filenames, entry.Name())
	}

	logger.Tracef(&config.EnableTrace_loader, "filenames=%v", filenames)

	sort.Strings(filenames)
	for _, name := range filenames {
		data, err := fs.ReadFile(fileSystem, path+"/"+name)
		if err != nil {
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
