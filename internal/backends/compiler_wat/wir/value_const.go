package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wat"
	"github.com/wa-lang/wa/internal/logger"
)

type Const interface {
	Value
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
		return &constRune{aConst: aConst{typ: t, lit: lit}}

	case I32:
		return &constI32{aConst: aConst{typ: t, lit: lit}}

	case U32:
		return &constU32{aConst: aConst{typ: t, lit: lit}}

	case I64:
		return &constI64{aConst: aConst{typ: t, lit: lit}}

	case U64:
		return &constU64{aConst: aConst{typ: t, lit: lit}}

	case F32:
		return &constF32{aConst: aConst{typ: t, lit: lit}}

	case F64:
		return &constF64{aConst: aConst{typ: t, lit: lit}}

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
func (c *aConst) EmitInit() []wat.Inst    { logger.Fatal("不可0值化常数"); return nil }
func (c *aConst) EmitSet() []wat.Inst     { logger.Fatal("不可Pop至常数"); return nil }
func (c *aConst) EmitRelease() []wat.Inst { logger.Fatal("不可清除常数"); return nil }

/**************************************
constRune:
**************************************/
type constRune struct {
	aConst
}

func (c *constRune) EmitGet() []wat.Inst { return []wat.Inst{wat.NewInstConst(wat.I32{}, c.lit)} }

/**************************************
constI32:
**************************************/
type constI32 struct {
	aConst
}

func (c *constI32) EmitGet() []wat.Inst { return []wat.Inst{wat.NewInstConst(wat.I32{}, c.lit)} }

/**************************************
constU32:
**************************************/
type constU32 struct {
	aConst
}

func (c *constU32) EmitGet() []wat.Inst { return []wat.Inst{wat.NewInstConst(wat.U32{}, c.lit)} }

/**************************************
constI64:
**************************************/
type constI64 struct {
	aConst
}

func (c *constI64) EmitGet() []wat.Inst { return []wat.Inst{wat.NewInstConst(wat.I64{}, c.lit)} }

/**************************************
constU64:
**************************************/
type constU64 struct {
	aConst
}

func (c *constU64) EmitGet() []wat.Inst { return []wat.Inst{wat.NewInstConst(wat.U64{}, c.lit)} }

/**************************************
constF32:
**************************************/
type constF32 struct {
	aConst
}

func (c *constF32) EmitGet() []wat.Inst { return []wat.Inst{wat.NewInstConst(wat.F32{}, c.lit)} }

/**************************************
constF64:
**************************************/
type constF64 struct {
	aConst
}

func (c *constF64) EmitGet() []wat.Inst { return []wat.Inst{wat.NewInstConst(wat.F64{}, c.lit)} }
