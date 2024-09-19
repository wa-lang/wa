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
	*ast.Func
	color
}

func new_RemoveUnusedPass(m *ast.Module) *_RemoveUnusedPass {
	p := &_RemoveUnusedPass{m: m}
	p.funcs = make(map[string]*funcObj, len(m.Funcs))
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
			if fn.Name != "" && fn.Name == elem.Name {
				p.markFuncReachable(p.funcs[fn.Name])
				continue Loop
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
	m.Funcs = p.m.Funcs[:0]
	for _, fn := range p.funcs {
		switch fn.color {
		case white:
			// skip
		case black:
			m.Funcs = append(m.Funcs, fn.Func)
		}
	}

	return &m
}

func (p *_RemoveUnusedPass) markFuncReachable(fn *funcObj) {
	fn.color = black
	for _, ins := range fn.Body.Insts {
		switch ins := ins.(type) {
		case ast.Ins_Call:
			if xFn := p.funcs[ins.X]; xFn.color == white {
				p.markFuncReachable(xFn)
			}
		case ast.Ins_TableSet:
			if xFn := p.funcs[ins.X]; xFn.color == white {
				p.markFuncReachable(xFn)
			}
		}
	}
}
