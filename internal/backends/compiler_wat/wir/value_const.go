package wir

import (
	"strconv"

	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wtypes"
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

func NewConst(t wtypes.ValueType, v interface{}) Const {
	switch t.(type) {
	case wtypes.Int32:
		if v == nil {
			return NewConstI32(0)
		}

		if c, ok := v.(int); ok {
			return NewConstI32(int32(c))
		}
		logger.Fatal("Todo")

	case wtypes.Int64:
		if v == nil {
			return NewConstI64(0)
		}

		if c, ok := v.(int); ok {
			return NewConstI64(int64(c))
		}
		logger.Fatal("Todo")

	default:
		logger.Fatal("Todo")
	}

	return nil
}

/**************************************
ConstInt32:
**************************************/
type ConstI32 struct {
	x int32
}

func NewConstI32(x int32) *ConstI32        { return &ConstI32{x: x} }
func (c *ConstI32) Name() string           { return strconv.FormatInt(int64(c.x), 10) }
func (c *ConstI32) Kind() ValueKind        { return ValueKindConst }
func (c *ConstI32) Type() wtypes.ValueType { return wtypes.Int32{} }
func (c *ConstI32) Raw() []Value           { return append([]Value(nil), c) }

/**************************************
ConstInt64:
**************************************/
type ConstI64 struct {
	x int64
}

func NewConstI64(x int64) *ConstI64        { return &ConstI64{x: x} }
func (c *ConstI64) Name() string           { return strconv.FormatInt(c.x, 10) }
func (c *ConstI64) Kind() ValueKind        { return ValueKindConst }
func (c *ConstI64) Type() wtypes.ValueType { return wtypes.Int64{} }
func (c *ConstI64) Raw() []Value           { return append([]Value(nil), c) }
