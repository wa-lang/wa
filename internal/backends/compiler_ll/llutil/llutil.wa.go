// 版权 @2021 凹语言 作者。保留所有权利。

package llutil

import (
	"github.com/wa-lang/wa/internal/3rdparty/llir"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llconstant"
	"github.com/wa-lang/wa/internal/3rdparty/llir/lltypes"
	"github.com/wa-lang/wa/internal/3rdparty/llir/llvalue"
)

func GetType(m *llir.Module, typName string) lltypes.Type {
	for _, typ := range m.TypeDefs {
		if typ.Name() == typName {
			return typ
		}
	}
	return nil
}

func GetFunc(m *llir.Module, fnName string) *llir.Func {
	for _, x := range m.Funcs {
		if x.Name() == fnName {
			return x
		}
	}
	return nil
}

func PtrElemType(src llvalue.Value) lltypes.Type {
	return src.Type().(*lltypes.PointerType).ElemType
}

func StrConstant(in string) *llconstant.CharArray {
	return llconstant.NewCharArray(append([]byte(in), 0))
}

func StrToi8Ptr(block *llir.Block, src llvalue.Value) *llir.InstGetElementPtr {
	return block.NewGetElementPtr(
		PtrElemType(src), src,
		llconstant.NewInt(lltypes.I32, 0),
		llconstant.NewInt(lltypes.I32, 0),
	)
}

func StrLen(block *llir.Block, src llvalue.Value) llvalue.Value {
	if _, ok := src.Type().(*lltypes.PointerType); ok {
		l := block.NewGetElementPtr(
			PtrElemType(src), src,
			llconstant.NewInt(lltypes.I32, 0),
			llconstant.NewInt(lltypes.I32, 0),
		)
		return block.NewLoad(PtrElemType(l), l)
	}
	return block.NewExtractValue(src, 0)
}

func StrToI8Ptr(block *llir.Block, src llvalue.Value) llvalue.Value {
	if _, ok := src.Type().(*lltypes.PointerType); ok {
		l := block.NewGetElementPtr(
			PtrElemType(src), src,
			llconstant.NewInt(lltypes.I32, 0),
			llconstant.NewInt(lltypes.I32, 1),
		)
		return block.NewLoad(PtrElemType(l), l)
	}
	return block.NewExtractValue(src, 1)
}

func Slice(itemType lltypes.Type) *lltypes.StructType {
	return lltypes.NewStruct(
		lltypes.I32,                  // Len
		lltypes.I32,                  // Cap
		lltypes.I32,                  // Array Offset
		lltypes.NewPointer(itemType), // Content
	)
}

func String() *lltypes.StructType {
	return lltypes.NewStruct(
		lltypes.NewPointer(lltypes.I8), // Content
		lltypes.I32,                    // String length
	)
}

func StringLen(stringType lltypes.Type) *llir.Func {
	param := llir.NewParam("input", stringType)
	res := llir.NewFunc("string_len", lltypes.I64, param)
	block := res.NewBlock("entry")
	block.NewRet(block.NewExtractValue(param, 0))
	return res
}
