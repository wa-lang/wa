package wat2c

import (
	"fmt"
	"strconv"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2cWorker) findLocalName(fn *ast.Func, ident string) string {
	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(fn.Type.Params)+len(fn.Body.Locals) {
			panic(fmt.Sprintf("wat2c: unknown local %q", ident))
		}
		return p.localNames[idx]
	}
	for idx, arg := range fn.Type.Params {
		if arg.Name == ident {
			return p.localNames[idx]
		}
	}
	for idx, arg := range fn.Body.Locals {
		if arg.Name == ident {
			return p.localNames[len(fn.Type.Params)+idx]
		}
	}
	panic("unreachable")
}

func (p *wat2cWorker) findType(ident string) *ast.FuncType {
	if idx, err := strconv.Atoi(ident); err == nil {
		if idx < 0 || idx >= len(p.m.Types) {
			panic(fmt.Sprintf("wat2c: unknown type %q", ident))
		}
		return p.m.Types[idx].Type
	}
	for _, x := range p.m.Types {
		if x.Name == ident {
			return x.Type
		}
	}
	panic(fmt.Sprintf("wat2c: unknown type %q", ident))
}

func (p *wat2cWorker) findFuncType(ident string) *ast.FuncType {
	idx := p.findFuncIndex(ident)
	if idx < len(p.m.Imports) {
		return p.m.Imports[idx].FuncType
	}

	return p.m.Funcs[idx-len(p.m.Imports)].Type
}

func (p *wat2cWorker) findFuncIndex(ident string) int {
	if idx, err := strconv.Atoi(ident); err == nil {
		return idx
	}

	var importCount int
	for _, x := range p.m.Imports {
		if x.ObjKind == token.FUNC {
			if x.FuncName == ident {
				return importCount
			}
			importCount++
		}
	}
	for i, fn := range p.m.Funcs {
		if fn.Name == ident {
			return importCount + i
		}
	}
	panic(fmt.Sprintf("wat2c: unknown func %q", ident))
}

func (p *wat2cWorker) findLabelName(label string) string {
	idx := p.findLabelIndex(label)
	if idx < len(p.labelScope) {
		return p.labelScope[len(p.labelScope)-idx-1]
	}
	panic(fmt.Sprintf("wat2c: unknown label %q", label))
}

func (p *wat2cWorker) findLabelIndex(label string) int {
	if idx, err := strconv.Atoi(label); err == nil {
		return idx
	}
	for i := 0; i < len(p.labelScope); i++ {
		if s := p.labelScope[len(p.labelScope)-i-1]; s == label {
			return i
		}
	}
	panic(fmt.Sprintf("wat2c: unknown label %q", label))
}

func (p *wat2cWorker) enterLabelScope(label string) {
	p.labelScope = append(p.labelScope, label)
}
func (p *wat2cWorker) leaveLabelScope() {
	p.labelScope = p.labelScope[:len(p.labelScope)-1]
}
