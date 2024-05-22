// 版权 @2024 凹语言 作者。保留所有权利。

package loaderx

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/types"
)

var _ types.Importer = (*pkgImporter)(nil)

type pkgImporter struct {
	*Universe
}

func (p *pkgImporter) Import(pkgpath string) (*types.Package, error) {
	if pkg, ok := p.Pkgs[PkgpathURI(pkgpath)]; ok {
		return pkg.Pkg, nil
	}

	pkg, err := p.loadPackage(pkgpath)
	if err != nil {
		return nil, err
	}

	return pkg.Pkg, nil
}

func (p *pkgImporter) loadPackage(pkgpath string) (*Package, error) {
	var err error

	if pkg, ok := p.Pkgs[PkgpathURI(pkgpath)]; ok {
		return pkg, nil
	}

	pkg := &Package{
		Version: 0,
		UriPath: "",
		PkgPath: pkgpath,
	}

	pkg.FileNames, err = p.GlobSources(pkgpath)
	if err != nil {
		return nil, err
	}

	for _, srcPath := range pkg.FileNames {
		f, err := parser.ParseFile(nil, p.Fset, srcPath, nil, parser.AllErrors)
		if err != nil && f == nil {
			return nil, err
		}
		pkg.Files = append(pkg.Files, f)
	}

	conf := types.Config{Importer: p}
	pkg.Info = &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	pkg.Pkg, err = conf.Check(pkgpath, p.Fset, pkg.Files, nil)
	if err != nil {
		return nil, err
	}

	p.Pkgs[PkgpathURI(pkgpath)] = pkg
	return pkg, nil
}

func (p *pkgImporter) GlobSources(pkgpath string) (matches []string, err error) {
	return nil, nil
}

func (p *pkgImporter) GlobMainSources(pkgpath string) (matches []string, err error) {
	return nil, nil
}
