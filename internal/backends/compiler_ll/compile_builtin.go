// 版权 @2021 凹语言 作者。保留所有权利。

package compiler

import (
	"fmt"

	"github.com/wa-lang/wa/internal/3rdparty/llir"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	"github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
	"github.com/wa-lang/wa/internal/config"
	"github.com/wa-lang/wa/internal/logger"
	"github.com/wa-lang/wa/internal/ssa"
)

func (p *Compiler) defineBuiltin() error {
	// syscall: int write(i32 fd, void* buf, int n)
	p.module.NewFunc("write",
		lltypes.I64,
		llir.NewParam("fd", lltypes.I32),
		llir.NewParam("buf", lltypes.NewPointer(lltypes.I8)),
		llir.NewParam("size", lltypes.I64),
	)

	p.module.NewTypeDef(g_typ_complex64_name,
		lltypes.NewStruct(lltypes.Float, lltypes.Float),
	)
	p.module.NewTypeDef(g_typ_complex128_name,
		lltypes.NewStruct(lltypes.Double, lltypes.Double),
	)
	p.module.NewTypeDef(g_typ_string_name,
		lltypes.NewStruct(lltypes.NewPointer(lltypes.I8), lltypes.I64),
	)

	p.module.NewFunc("ugo_print_rune",
		lltypes.I32,
		llir.NewParam("x", lltypes.I32),
	)
	p.module.NewFunc("ugo_print_bool",
		lltypes.I32,
		llir.NewParam("x", lltypes.I8),
	)
	p.module.NewFunc("ugo_print_int",
		lltypes.I32,
		llir.NewParam("x", lltypes.I32),
	)
	p.module.NewFunc("ugo_print_int64",
		lltypes.I32,
		llir.NewParam("x", lltypes.I64),
	)
	p.module.NewFunc("ugo_print_ptr",
		lltypes.I32,
		llir.NewParam("p", lltypes.I8Ptr),
	)
	p.module.NewFunc("ugo_print_cstring",
		lltypes.I32,
		llir.NewParam("s", lltypes.NewPointer(lltypes.I8)),
	)
	p.module.NewFunc("ugo_print_cstring_len",
		lltypes.I32,
		llir.NewParam("s", lltypes.NewPointer(lltypes.I8)),
		llir.NewParam("n", lltypes.I32),
	)

	fn := p.module.NewFunc("ugo_printf",
		lltypes.I32,
		llir.NewParam("format", lltypes.NewPointer(lltypes.I8)),
	)
	fn.Sig.Variadic = true

	p.module.NewFunc("ugo_cstring_join",
		lltypes.NewPointer(lltypes.I8),
		llir.NewParam("s0", lltypes.NewPointer(lltypes.I8)),
		llir.NewParam("s1", lltypes.NewPointer(lltypes.I8)),
	)
	p.module.NewFunc("ugo_cstring_slice",
		lltypes.NewPointer(lltypes.I8),
		llir.NewParam("s", lltypes.NewPointer(lltypes.I8)),
		llir.NewParam("low", lltypes.I32),
		llir.NewParam("high", lltypes.I32),
	)
	p.module.NewFunc("ugo_cstring_index",
		lltypes.I32,
		llir.NewParam("s", lltypes.NewPointer(lltypes.I8)),
		llir.NewParam("idx", lltypes.I32),
	)
	p.module.NewFunc("ugo_cstring_cmp",
		lltypes.I32,
		llir.NewParam("s0", lltypes.NewPointer(lltypes.I8)),
		llir.NewParam("s1", lltypes.NewPointer(lltypes.I8)),
	)

	p.module.NewFunc("ugo_string_new",
		lltypes.NewPointer(
			p.getLLType(g_typ_string_name),
		),
		llir.NewParam("size", lltypes.I32),
		llir.NewParam("data", lltypes.NewPointer(lltypes.I8)),
	)
	p.module.NewFunc("ugo_print_string", lltypes.I32,
		llir.NewParam("s", lltypes.NewPointer(
			p.getLLType(g_typ_string_name),
		)),
	)

	return nil
}

func (p *Compiler) callBuiltin(expr *ssa.Call, builtin *ssa.Builtin) llvalue.Value {
	logger.Tracef(&config.EnableTrace_compiler, "expr=%v, builtin=%v", expr, builtin)

	switch builtin.Name() {
	case "print", "println":
		for i, arg := range expr.Common().Args {
			if i > 0 {
				p.builtin_print_space(expr, builtin)
			}

			arg := p.getValue(arg)
			if p.llIsStrType(arg.Type()) {
				p.builtin_print_str(expr, builtin, arg)
				continue
			}

			switch typ := arg.Type().(type) {
			case *lltypes.IntType:
				switch typ.BitSize {
				case 32:
					p.builtin_print_int(expr, builtin, arg)
				case 64:
					p.builtin_print_int64(expr, builtin, arg)
				}
			default:
				// todo: string => *lltypes.StructType
				fmt.Printf("callBuiltin: %T, %v\n", typ, p.llIsStrType(typ))
			}
		}
		if builtin.Name() == "println" {
			p.builtin_print_newline(expr, builtin)
		}
		return nil
	}
	return nil
}

func (p *Compiler) builtin_print_str(expr *ssa.Call, builtin *ssa.Builtin, arg llvalue.Value) llvalue.Value {
	llirBlock := p.curBlockEntries[expr.Block()]

	switch x := arg.(type) {
	case *llir.Param:
		logger.Panicf("TODO: arg: %T, %v", arg, arg)
		panic("unreachable")

	case *llconstant.Struct:
		data := x.Fields[0].(*llconstant.CharArray)

		_p0 := llirBlock.NewAlloca(data.Typ)
		llirBlock.NewStore(data, _p0)

		_p_data00 := llirBlock.NewGetElementPtr(
			data.Typ, _p0,
			llconstant.NewInt(lltypes.I64, 0),
			llconstant.NewInt(lltypes.I64, 0),
		)

		_size32 := llconstant.NewInt(lltypes.I32, int64(len(data.X)))

		llirBlock.NewCall(p.getFunc("ugo_print_cstring_len"), _p_data00, _size32)
		return nil
	default:
		logger.Panicf("arg: %T, %v", arg, arg)
		panic("unreachable")
	}
}

func (p *Compiler) builtin_print_newline(expr *ssa.Call, builtin *ssa.Builtin) llvalue.Value {
	llirBlock := p.curBlockEntries[expr.Block()]
	llirBlock.NewCall(p.getFunc("ugo_print_rune"), llconstant.NewInt(lltypes.I32, '\n'))
	return nil
}
func (p *Compiler) builtin_print_int(expr *ssa.Call, builtin *ssa.Builtin, arg llvalue.Value) llvalue.Value {
	llirBlock := p.curBlockEntries[expr.Block()]
	llirBlock.NewCall(p.getFunc("ugo_print_int"), arg)
	return nil
}
func (p *Compiler) builtin_print_int64(expr *ssa.Call, builtin *ssa.Builtin, arg llvalue.Value) llvalue.Value {
	llirBlock := p.curBlockEntries[expr.Block()]
	llirBlock.NewCall(p.getFunc("ugo_print_int64"), arg)
	return nil
}

func (p *Compiler) builtin_print_rune(expr *ssa.Call, builtin *ssa.Builtin, arg rune) llvalue.Value {
	llirBlock := p.curBlockEntries[expr.Block()]
	llirBlock.NewCall(p.getFunc("ugo_print_rune"), llconstant.NewInt(lltypes.I32, int64(arg)))
	return nil
}
func (p *Compiler) builtin_print_space(expr *ssa.Call, builtin *ssa.Builtin) llvalue.Value {
	llirBlock := p.curBlockEntries[expr.Block()]
	llirBlock.NewCall(p.getFunc("ugo_print_rune"), llconstant.NewInt(lltypes.I32, ' '))
	return nil
}
