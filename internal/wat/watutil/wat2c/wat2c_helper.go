package wat2c

import (
	"fmt"
	"strconv"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

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
