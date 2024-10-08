// 版权 @2024 凹语言 作者。保留所有权利。

package watutil

import (
	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// 删除未使用的对象
type _RemoveUnusedPass struct {
	m *ast.Module

	funcs map[string]*funcObj
}

type color int

const (
	white color = 0
	black color = -1
)

type funcObj struct {
	isImport bool
	*ast.Func
	color
}

func new_RemoveUnusedPass(m *ast.Module) *_RemoveUnusedPass {
	p := &_RemoveUnusedPass{m: m}

	p.funcs = make(map[string]*funcObj, len(m.Funcs)+len(m.Imports))

	for _, importSpec := range m.Imports {
		if importSpec.ObjKind == token.FUNC {
			p.funcs[importSpec.FuncName] = &funcObj{
				isImport: true,
				Func: &ast.Func{
					Name: importSpec.FuncName,
					Type: importSpec.FuncType,
					Body: &ast.FuncBody{},
				},
			}
		}
	}
	for _, fn := range m.Funcs {
		p.funcs[fn.Name] = &funcObj{Func: fn}
	}

	return p
}

func (p *_RemoveUnusedPass) DoPass() *ast.Module {
	for i := range p.funcs {
		p.funcs[i].color = white
	}

Loop:
	for _, fn := range p.m.Funcs {
		// start
		if fn.Name != "" && fn.Name == p.m.Start {
			p.markFuncReachable(p.funcs[fn.Name])
			continue
		}

		// table elem

		for _, elem := range p.m.Elem {
			for _, elemValue := range elem.Values {
				if fn.Name != "" && fn.Name == elemValue {
					p.markFuncReachable(p.funcs[fn.Name])
					continue Loop
				}
			}
		}

		// export
		for _, exp := range p.m.Exports {
			if exp.Kind == token.FUNC {
				if exp.Name != "" && fn.Name == exp.FuncIdx {
					p.markFuncReachable(p.funcs[fn.Name])
					continue Loop
				}
			}
		}
	}

	m := *p.m
	m.Imports = p.m.Imports[:0]
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind == token.FUNC {
			if fnObj := p.funcs[importSpec.FuncName]; fnObj.color == white {
				continue // skip
			}
		}
		m.Imports = append(m.Imports, importSpec)
	}

	m.Funcs = p.m.Funcs[:0]
	for _, fn := range p.m.Funcs {
		if fnObj := p.funcs[fn.Name]; fnObj.color == black {
			m.Funcs = append(m.Funcs, fn)
		}
	}

	return &m
}

func (p *_RemoveUnusedPass) markFuncReachable(fn *funcObj) {
	fn.color = black
	for _, ins := range fn.Body.Insts {
		p.markFuncReachable_ins(ins)
	}
}

func (p *_RemoveUnusedPass) markFuncReachable_ins(ins ast.Instruction) {
	switch ins := ins.(type) {
	case ast.Ins_Call:
		if xFn := p.funcs[ins.X]; xFn.color == white {
			p.markFuncReachable(xFn)
		}
	case ast.Ins_TableSet:
		if xFn := p.funcs[ins.TableIdx]; xFn.color == white {
			p.markFuncReachable(xFn)
		}
	case ast.Ins_Block:
		for _, x := range ins.List {
			p.markFuncReachable_ins(x)
		}
	case ast.Ins_Loop:
		for _, x := range ins.List {
			p.markFuncReachable_ins(x)
		}
	case ast.Ins_If:
		for _, x := range ins.Body {
			p.markFuncReachable_ins(x)
		}
		for _, x := range ins.Else {
			p.markFuncReachable_ins(x)
		}
	}
}
