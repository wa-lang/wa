// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir"
	"github.com/wa-lang/wa/internal/ssa"
	"github.com/wa-lang/wa/internal/types"
)

func GetFnMangleName(v interface{}) string {
	var name string

	switch f := v.(type) {
	case *ssa.Function:
		name = wir.GetPkgMangleName(f.Pkg.Pkg.Path())
		if recv := f.Signature.Recv(); recv != nil {
			switch rt := recv.Type().(type) {
			case *types.Named:
				name += rt.Obj().Name()
			case *types.Pointer:
				btype, ok := rt.Elem().(*types.Named)
				if !ok {
					panic("Unreachable")
				}
				name += btype.Obj().Name()
			}
			name += "."
		}
		name += f.Name()

	case *types.Func:
		name = wir.GetPkgMangleName(f.Pkg().Path())
		sig := f.Type().(*types.Signature)
		if recv := sig.Recv(); recv != nil {
			switch rt := recv.Type().(type) {
			case *types.Named:
				name += rt.Obj().Name()
			case *types.Pointer:
				btype, ok := rt.Elem().(*types.Named)
				if !ok {
					panic("Unreachable")
				}
				name += btype.Obj().Name()
			}
			name += "."
		}
		name += f.Name()
	}

	return name
}
