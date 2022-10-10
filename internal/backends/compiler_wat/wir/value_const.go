// 版权 @2022 凹语言 作者。保留所有权利。

package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

type Const interface {
	Value
	isConst()
}

/**************************************
ConstZero:
**************************************/
/*type ConstZero struct {
}

func NewConstZero() *ConstZero              { return &ConstZero{} }
func (c *ConstZero) Name() string           { return "0" }
func (c *ConstZero) Kind() ValueKind        { return ValueKindConst }
func (c *ConstZero) Type() wtypes.ValueType { return wtypes.Void{} }
func (c *ConstZero) Raw() []Value           { return append([]Value(nil), c) }
//*/

func NewConst(t ValueType, lit string) Const {
	switch t.(type) {
	case RUNE:
		return &aConst{typ: t, lit: lit}

	case I32:
		return &aConst{typ: t, lit: lit}

	case U32:
		return &aConst{typ: t, lit: lit}

	case I64:
		return &aConst{typ: t, lit: lit}

	case U64:
		return &aConst{typ: t, lit: lit}

	case F32:
		return &aConst{typ: t, lit: lit}

	case F64:
		return &aConst{typ: t, lit: lit}

	default:
		logger.Fatal("Todo")
	}

	return nil
}

/**************************************
aConst:
**************************************/
type aConst struct {
	typ ValueType
	lit string
}

func (c *aConst) Name() string            { return c.lit }
func (c *aConst) Kind() ValueKind         { return ValueKindConst }
func (c *aConst) Type() ValueType         { return c.typ }
func (c *aConst) raw() []wat.Value        { logger.Fatal("Todo"); return nil }
func (c *aConst) isConst()                {}
func (c *aConst) EmitInit() []wat.Inst    { logger.Fatal("不可0值化常数"); return nil }
func (c *aConst) EmitPop() []wat.Inst     { logger.Fatal("不可Pop至常数"); return nil }
func (c *aConst) EmitRelease() []wat.Inst { logger.Fatal("不可清除常数"); return nil }
func (c *aConst) emitLoadFromAddr(addr Value, offset int) []wat.Inst {
	logger.Fatal("不可Load常数")
	return nil
}
func (c *aConst) EmitPush() []wat.Inst {
	return []wat.Inst{wat.NewInstConst(toWatType(c.Type()), c.lit)}
}
func (c *aConst) emitStoreToAddr(addr Value, offset int) []wat.Inst {
	if !addr.Type().(Pointer).Base.Equal(c.Type()) {
		logger.Fatal("Type not match")
		return nil
	}
	insts := addr.EmitPush()
	insts = append(insts, c.EmitPush()...)
	insts = append(insts, wat.NewInstStore(toWatType(c.Type()), offset, 1))
	return insts
}
