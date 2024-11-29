// 版权 @2021 凹语言 作者。保留所有权利。

package loader

import (
	"fmt"
	"io/fs"
	"os"
	pathpkg "path"
	"path/filepath"
	"sort"
	"strings"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/ast/astutil"
	"wa-lang.org/wa/internal/config"
	wzparser "wa-lang.org/wa/internal/frontend/wz/parser"
	"wa-lang.org/wa/internal/loader/buildtag"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/ssa"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/types"
	"wa-lang.org/wa/internal/wamime"
	wasrc "wa-lang.org/wa/waroot/src"
)

var _loadRuntime bool = true

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

	// 注册 assert 函数
	if p.cfg.UnitTest {
		types.DefPredeclaredTestFuncs()
	}

	p.vfs = *vfs
	p.prog.Cfg = &p.cfg
	p.prog.Manifest = manifest
	p.prog.Fset = token.NewFileSet()

	if p.vfs.Std == nil {
		// pkg/std
		stdPath := filepath.Join(manifest.Root, "pkg", "std")
		if dirPathExists(stdPath) {
			vfs.Std = os.DirFS(stdPath)
		} else {
			vfs.Std = wasrc.GetStdFS()
		}
	}

	// import "runtime"
	if _loadRuntime {
		logger.Trace(&config.EnableTrace_loader, "import runtime")
		if _, err := p.Import("runtime"); err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, err
		}
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
	var filenames []string

	if pkgpath == "unsafe" {
		pkg.Pkg = types.Unsafe
		pkg.Info = &types.Info{
			Types:      make(map[ast.Expr]types.TypeAndValue),
			Defs:       make(map[*ast.Ident]types.Object),
			Uses:       make(map[*ast.Ident]types.Object),
			Implicits:  make(map[ast.Node]types.Object),
			Selections: make(map[*ast.SelectorExpr]*types.Selection),
			Scopes:     make(map[ast.Node]*types.Scope),
		}

		p.prog.Pkgs[pkgpath] = &pkg
		return pkg.Pkg, nil
	}

	// 解析当前包的汇编代码
	pkg.WsFiles, err = p.ParseDir_wsFiles(pkgpath)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	// 解析当前包的宿主代码
	pkg.WImportFiles, err = p.ParseDir_hostImportFiles(pkgpath)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	// 解析当前包的 AST, 隐含了是否测试模式
	filenames, pkg.Files, err = p.ParseDir(pkgpath)
	if err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}

	// main 包隐式导入 runtime
	if _loadRuntime {
		if pkgpath == p.prog.Manifest.MainPkg && pkgpath != "runtime" {
			if len(pkg.Files) > 0 {
				f, err := parser.ParseFile(nil, p.prog.Fset, "_$main$runtime.wa", `import "runtime" => _`, parser.AllErrors)
				if err != nil {
					panic(err)
				}
				pkg.Files[0].Decls = append(f.Decls, pkg.Files[0].Decls...)
			}
		}
	}

	// 过滤 build-tag, main 包忽略
	if pkgpath != p.prog.Manifest.MainPkg || p.prog.Manifest.IsStd {
		var pkgFileNames = make([]string, 0, len(filenames))
		var pkgFiles = make([]*ast.File, 0, len(pkg.Files))
		for i, f := range pkg.Files {
			skiped, err := p.isSkipedAstFile(f)
			if err != nil {
				return nil, err
			}
			if skiped {
				continue
			}
			pkgFileNames = append(pkgFileNames, filenames[i])
			pkgFiles = append(pkgFiles, f)
		}
		filenames = pkgFileNames
		pkg.Files = pkgFiles
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

	// 提取测试信息
	if p.cfg.UnitTest {
		pkg.TestInfo, err = p.parseTestInfo(&pkg, filenames)
		if err != nil {
			return nil, err
		}
	}

	logger.Tracef(&config.EnableTrace_loader, "save pkgpath: %v", pkgpath)

	p.prog.Pkgs[pkgpath] = &pkg
	return pkg.Pkg, nil
}

func (p *_Loader) parseTestInfo(pkg *Package, filenames []string) (*TestInfo, error) {
	tInfo := new(TestInfo)

	for i, filename := range filenames {
		if !p.isTestFile(filename) {
			continue
		}

		file := pkg.Files[i]
		tInfo.Files = append(tInfo.Files, filename)

		// 提取测试/基准/示例函数
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				name := fn.Name.Name

				// 函数参数均为空
				obj := pkg.Pkg.Scope().Lookup(name)
				{
					var bValidFuncType = false
					if fn, ok := obj.(*types.Func); ok {
						if sig, ok := fn.Type().(*types.Signature); ok {
							if sig.Recv() == nil && sig.Params().Len() == 0 && sig.Results().Len() == 0 {
								bValidFuncType = true
							}
						}
					}
					if !bValidFuncType {
						continue // skip
					}
				}

				switch {
				case strings.HasPrefix(name, "Test"):
					output, isPanic := p.parseExampleOutputComment(file, fn)
					tInfo.Tests = append(tInfo.Tests, TestFuncInfo{
						FuncPos:     obj.Pos(),
						Name:        name,
						Output:      output,
						OutputPanic: isPanic,
					})
				case strings.HasPrefix(name, "Bench"):
					tInfo.Benchs = append(tInfo.Benchs, TestFuncInfo{
						FuncPos: obj.Pos(),
						Name:    name,
					})
				case strings.HasPrefix(name, "Example"):
					output, isPanic := p.parseExampleOutputComment(file, fn)
					tInfo.Examples = append(tInfo.Examples, TestFuncInfo{
						FuncPos:     obj.Pos(),
						Name:        name,
						Output:      output,
						OutputPanic: isPanic,
					})
				}
			}
		}
	}

	return tInfo, nil
}

func (p *_Loader) parseExampleOutputComment(f *ast.File, fn *ast.FuncDecl) (output string, isPanic bool) {
	for _, commentGroup := range f.Comments {
		if commentGroup.Pos() <= fn.Body.Pos() {
			continue
		}
		if commentGroup.End() > fn.Body.End() {
			break
		}

		for j, comment := range commentGroup.List {
			switch comment.Text {
			case "// Output:":
				var lineTexts []string
				for _, x := range commentGroup.List[j+1:] {
					if !strings.HasPrefix(x.Text, "//") {
						break
					}
					lineTexts = append(lineTexts, strings.TrimSpace(x.Text[2:]))
				}

				return strings.Join(lineTexts, "\n"), false

			case "// Output(panic):":
				var lineTexts []string
				for _, x := range commentGroup.List[j+1:] {
					if !strings.HasPrefix(x.Text, "//") {
						break
					}
					lineTexts = append(lineTexts, strings.TrimSpace(x.Text[2:]))
				}

				return strings.Join(lineTexts, "\n"), true
			}
		}
	}

	return "", false
}

func (p *_Loader) ParseDir_wsFiles(pkgpath string) (files []*WsFile, err error) {
	logger.Tracef(&config.EnableTrace_loader, "pkgpath: %v", pkgpath)

	if p.cfg.WaBackend == "" {
		panic("unreachable")
	}

	var (
		extNames          = []string{fmt.Sprintf(".%s.ws", p.cfg.WaBackend)}
		unitTestMode bool = false

		filenames []string
		datas     [][]byte
	)

	switch {
	case p.isStdPkg(pkgpath):
		logger.Tracef(&config.EnableTrace_loader, "isStdPkg; pkgpath: %v", pkgpath)

		filenames, datas, err = p.readDirFiles(p.vfs.Std, pkgpath, unitTestMode, extNames)
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

		filenames, datas, err = p.readDirFiles(p.vfs.App, relpkg, unitTestMode, extNames)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, err
		}

		logger.Trace(&config.EnableTrace_loader, "isSelfPkg; return ok")

	default: // vendor
		logger.Tracef(&config.EnableTrace_loader, "vendorPkg; pkgpath: %v", pkgpath)

		filenames, datas, err = p.readDirFiles(p.vfs.Vendor, pkgpath, unitTestMode, extNames)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, err
		}
	}

	for i := 0; i < len(filenames); i++ {
		files = append(files, &WsFile{
			Name: filenames[i],
			Code: string(datas[i]),
		})
	}
	return
}

func (p *_Loader) ParseDir_hostImportFiles(pkgpath string) (files []*WhostFile, err error) {
	logger.Tracef(&config.EnableTrace_loader, "pkgpath: %v", pkgpath)

	if p.cfg.Target == "" && p.prog.Manifest.Pkg.Target == "" {
		p.cfg.Target = config.WaOS_Default
		p.prog.Manifest.Pkg.Target = config.WaOS_Default
	}

	var (
		extNames          = []string{fmt.Sprintf(".import.%s", p.GetTargetOS())}
		unitTestMode bool = false

		filenames []string
		datas     [][]byte
	)

	switch {
	case p.isStdPkg(pkgpath):
		logger.Tracef(&config.EnableTrace_loader, "isStdPkg; pkgpath: %v", pkgpath)

		filenames, datas, err = p.readDirFiles(p.vfs.Std, pkgpath, unitTestMode, extNames)
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

		filenames, datas, err = p.readDirFiles(p.vfs.App, relpkg, unitTestMode, extNames)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, err
		}

		logger.Trace(&config.EnableTrace_loader, "isSelfPkg; return ok")

	default: // vendor
		logger.Tracef(&config.EnableTrace_loader, "vendorPkg; pkgpath: %v", pkgpath)

		filenames, datas, err = p.readDirFiles(p.vfs.Vendor, pkgpath, unitTestMode, extNames)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, err
		}
	}

	for i := 0; i < len(filenames); i++ {
		files = append(files, &WhostFile{
			Name: filenames[i],
			Code: string(datas[i]),
		})
	}
	return
}

func (p *_Loader) ParseDir(pkgpath string) (filenames []string, files []*ast.File, err error) {
	logger.Tracef(&config.EnableTrace_loader, "pkgpath: %v", pkgpath)

	var (
		pkgVFS       fs.FS
		extNames          = []string{".wa", ".wz", ".wa.go"}
		unitTestMode bool = false
		datas        [][]byte
	)
	if pkgpath == p.prog.Manifest.MainPkg && p.cfg.UnitTest {
		unitTestMode = true
	}

	switch {
	case p.isStdPkg(pkgpath):
		logger.Tracef(&config.EnableTrace_loader, "isStdPkg; pkgpath: %v", pkgpath)

		pkgVFS = p.vfs.Std
		filenames, datas, err = p.readDirFiles(p.vfs.Std, pkgpath, unitTestMode, extNames)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, nil, err
		}
	case p.isSelfPkg(pkgpath):
		relpkg := strings.TrimPrefix(pkgpath, p.prog.Manifest.Pkg.Pkgpath)
		if relpkg == "" {
			relpkg = "."
		}

		logger.Tracef(&config.EnableTrace_loader, "isSelfPkg; pkgpath=%v, relpkg=%v", pkgpath, relpkg)

		pkgVFS = p.vfs.App
		filenames, datas, err = p.readDirFiles(p.vfs.App, relpkg, unitTestMode, extNames)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, nil, err
		}

		logger.Trace(&config.EnableTrace_loader, "isSelfPkg; return ok")

	default: // vendor
		logger.Tracef(&config.EnableTrace_loader, "vendorPkg; pkgpath: %v", pkgpath)

		pkgVFS = p.vfs.Vendor
		filenames, datas, err = p.readDirFiles(p.vfs.Vendor, pkgpath, unitTestMode, extNames)
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
			return nil, nil, err
		}
	}

	logger.Tracef(&config.EnableTrace_loader, "filenames: %v", filenames)

	for i, filename := range filenames {
		var f *ast.File
		if wamime.GetCodeMime(filename, datas[i]) == "wz" {
			f, err = wzparser.ParseFile(nil, p.prog.Fset, filename, datas[i], wzparser.AllErrors|wzparser.ParseComments)
		} else {
			f, err = parser.ParseFile(nil, p.prog.Fset, filename, datas[i], parser.AllErrors|parser.ParseComments)
		}
		if err != nil {
			logger.Tracef(&config.EnableTrace_loader, "filename: %v", filename)
			logger.Tracef(&config.EnableTrace_loader, "datas[i]: %s", datas[i])
			logger.Tracef(&config.EnableTrace_loader, "err: %v", err)

			return nil, nil, err
		}
		files = append(files, f)
	}

	// read embed files
	for _, f := range files {
		for _, commentGroup := range f.Comments {
			info := astutil.ParseCommentInfo(commentGroup)
			if info.Embed == "" {
				continue
			}

			if f.EmbedMap == nil {
				f.EmbedMap = make(map[string]string)
			}

			vpath := pathpkg.Join(pkgpath, info.Embed)
			data, err := fs.ReadFile(pkgVFS, vpath)
			if err != nil {
				continue
			}

			f.EmbedMap[info.Embed] = string(data)
		}
	}

	return filenames, files, nil
}

func (p *_Loader) readDirFiles(fileSystem fs.FS, path string, unitTestMode bool, extNames []string) (filenames []string, datas [][]byte, err error) {
	if len(extNames) == 0 {
		panic("unreachable")
	}

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

		if p.isSkipedSouceFile(entry.Name(), unitTestMode, extNames) {
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
	return wasrc.IsStdPkg(pkgpath)
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
		return types.SizesFor(p.GetTargetArch())
	} else {
		return &types.StdSizes{
			WordSize: p.cfg.WaSizes.WordSize,
			MaxAlign: p.cfg.WaSizes.MaxAlign,
		}
	}
}

func (p *_Loader) isSkipedAstFile(f *ast.File) (bool, error) {
	// 1. 找 buildExpr
	var buildExpr *ast.Comment
	if f.Doc != nil {
		for _, x := range f.Doc.List {
			if buildtag.IsWaBuild(x.Text) {
				buildExpr = x
				break
			}
		}
	}
	if buildExpr == nil {
		for _, comment := range f.Comments {
			if buildExpr != nil {
				break
			}
			for _, x := range comment.List {
				if buildtag.IsWaBuild(x.Text) {
					buildExpr = x
					break
				}
			}
		}
	}

	// 没有 build-tag
	if buildExpr == nil {
		return false, nil
	}

	// 解析 build-tag
	expr, err := buildtag.Parse(buildExpr.Text)
	if err != nil {
		err = fmt.Errorf("%v: parsing #wa:build line: %w",
			p.prog.Fset.Position(buildExpr.Slash),
			err,
		)
		return false, err
	}
	ok := expr.Eval(func(tag string) bool {
		if tag == p.GetTargetOS() || tag == p.GetTargetArch() {
			return true
		}
		for _, x := range p.cfg.BuilgTags {
			if x == tag {
				return true
			}
		}
		return false
	})

	return !ok, err
}

func (p *_Loader) isSkipedSouceFile(filename string, unitTestMode bool, extNames []string) bool {
	if len(extNames) == 0 {
		panic("unreachable")
	}
	if strings.HasPrefix(filename, "_") {
		return true
	}
	if !p.hasExt(filename, extNames...) {
		return true
	}

	if !unitTestMode {
		if p.isTestFile(filename) {
			return true
		}
	}

	if p.GetTargetOS() != "" {
		var isTargetFile bool
		for _, ext := range extNames {
			for _, os := range config.WaOS_List {
				if strings.HasSuffix(filename, "_"+os+ext) {
					isTargetFile = true
					break
				}
			}
		}
		if isTargetFile {
			var shouldSkip = true
			for _, ext := range extNames {
				if strings.HasSuffix(filename, "_"+p.GetTargetOS()+ext) {
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

func (p *_Loader) isTestFile(filename string) bool {
	if i := strings.LastIndexAny(filename, "/\\"); i >= 0 {
		filename = filename[i+1:]
	}
	if strings.HasPrefix(filename, "test_") {
		return true
	}
	if strings.HasSuffix(filename, "_test.wa") {
		return true
	}
	if strings.HasSuffix(filename, "_test.wz") {
		return true
	}
	if strings.HasSuffix(filename, "_test.wa.go") {
		return true
	}
	return false
}

func (p *_Loader) GetTargetOS() string {
	if s := p.cfg.Target; s != "" {
		return s
	}
	if s := p.prog.Manifest.Pkg.Target; s != "" {
		return s
	}
	return config.WaOS_Default
}

func (p *_Loader) GetTargetArch() string {
	return config.WaArch_Default
}
