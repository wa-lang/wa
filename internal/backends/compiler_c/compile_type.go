// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_c

import (
	"wa-lang.org/wa/internal/backends/compiler_c/cir"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/ssa"
	"wa-lang.org/wa/internal/types"
)

func (p *CompilerC) compileType(t *ssa.Type) {
	switch typ := t.Type().(type) {
	case *types.Named:
		switch utyp := typ.Underlying().(type) {
		case *types.Basic:
			var v []cir.Var
			bt := cir.ToCType(utyp)
			v = append(v, *cir.NewVar("$"+bt.CIRString(), bt))
			p.curScope.AddStructDecl(t.Name(), v)
		case *types.Struct:
			var vs []cir.Var
			for i := 0; i < utyp.NumFields(); i++ {
				f := utyp.Field(i)
				var v cir.Var
				bt := cir.ToCType(f.Type())
				if f.Embedded() {
					v = *cir.NewVar("$"+bt.CIRString(), bt)
				} else {
					v = *cir.NewVar(f.Name(), bt)
				}
				vs = append(vs, v)
			}
			p.curScope.AddStructDecl(t.Name(), vs)
		case *types.Interface:
			logger.Fatalf("Todo: %T", utyp)
		}

	default:
		logger.Fatalf("Todo: %T", typ)
	}
}
