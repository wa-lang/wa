// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// 删除未使用的对象
type _RemoveUnusedPass struct {
	m *ast.Module

	usedTypesMap  map[string]*ast.TypeSection
	usedImportMap map[string]*ast.ImportSpec
	usedGlobalMap map[string]*ast.Global
	usedFuncMap   map[string]*ast.Func
}

func new_RemoveUnusedPass(m *ast.Module) *_RemoveUnusedPass {
	return &_RemoveUnusedPass{
		m: m,

		usedTypesMap:  make(map[string]*ast.TypeSection),
		usedImportMap: make(map[string]*ast.ImportSpec),
		usedGlobalMap: make(map[string]*ast.Global),
		usedFuncMap:   make(map[string]*ast.Func),
	}
}

func (p *_RemoveUnusedPass) DoPass() *ast.Module {
	p.walkStartFunc()
	p.walkExportFunc()

	m := p.m
	if len(m.Types) != len(p.usedTypesMap) {
		m.Types = make([]*ast.TypeSection, 0, len(p.usedTypesMap))
		for _, x := range p.m.Types {
			if _, ok := p.usedTypesMap[x.Name]; ok {
				m.Types = append(m.Types, x)
			}
		}
	}
	if len(m.Imports) != len(p.usedImportMap) {
		m.Imports = make([]*ast.ImportSpec, 0, len(p.usedImportMap))
		for _, x := range p.m.Imports {
			if x.ObjKind != token.FUNC {
				m.Imports = append(m.Imports, x)
				continue
			}
			if _, ok := p.usedImportMap[x.FuncName]; ok {
				m.Imports = append(m.Imports, x)
			}
		}
	}
	if len(m.Globals) != len(p.usedGlobalMap) {
		m.Globals = make([]*ast.Global, 0, len(p.usedGlobalMap))
		for _, x := range p.m.Globals {
			if _, ok := p.usedGlobalMap[x.Name]; ok {
				m.Globals = append(m.Globals, x)
			}
		}
	}
	if len(m.Funcs) != len(p.usedFuncMap) {
		m.Funcs = make([]*ast.Func, 0, len(p.usedFuncMap))
		for _, x := range p.m.Funcs {
			if _, ok := p.usedTypesMap[x.Name]; ok {
				m.Funcs = append(m.Funcs, x)
			}
		}
	}
	return m
}

func (p *_RemoveUnusedPass) walkStartFunc() {
	// todo
}

func (p *_RemoveUnusedPass) walkExportFunc() {
	// todo
}
