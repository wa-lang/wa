package wir

import (
	"github.com/wa-lang/wa/internal/backends/compiler_wat/wir/wtypes"
	"github.com/wa-lang/wa/internal/logger"
)

func NewVar(name string, kind ValueKind, typ wtypes.ValueType) Value {
	switch typ.(type) {
	case wtypes.Int32:
		return &VarI32{name: name, kind: kind, typ: typ}

	default:
		logger.Fatalf("Todo: %T", typ)
	}

	return nil
}

/**************************************
VarI32:
**************************************/
type VarI32 struct {
	name string
	kind ValueKind
	typ  wtypes.ValueType
}

func (v *VarI32) Name() string           { return v.name }
func (v *VarI32) Kind() ValueKind        { return v.kind }
func (v *VarI32) Type() wtypes.ValueType { return v.typ }
func (v *VarI32) Raw() []Value           { return append([]Value(nil), v) }
