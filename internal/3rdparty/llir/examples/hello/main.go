package main

import (
	"fmt"

	"github.com/wa-lang/wa/internal/3rdparty/llir"
	constant "github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	types "github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	value "github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

func main() {
	m := llir.NewModule()

	fnPuts := m.NewFunc("puts", types.I32, llir.NewParam("", types.NewPointer(types.I8)))

	fnMain := m.NewFunc("main", types.I32)
	{
		entry := fnMain.NewBlock("")
		constString := m.NewGlobalDef("", llStrConstant("hello llir!"))
		entry.NewCall(fnPuts, llToi8Ptr(entry, constString))
		entry.NewRet(constant.NewInt(types.I32, 0))
	}

	fmt.Println(m)
}

func llStrConstant(s string) constant.Constant {
	return constant.NewCharArrayFromString(s + "\x00")
}

func llToi8Ptr(block *llir.Block, src value.Value) *llir.InstGetElementPtr {
	return block.NewGetElementPtr(
		src.Type().(*types.PointerType).ElemType,
		src,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)
}
