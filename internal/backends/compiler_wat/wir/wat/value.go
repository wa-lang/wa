// 版权 @2022 凹语言 作者。保留所有权利。

package wat

/**************************************
I32:
**************************************/
type I32 struct {
}

func (t I32) Name() string           { return "i32" }
func (t I32) Equal(u ValueType) bool { _, ok := u.(I32); return ok }

/**************************************
U32:
**************************************/
type U32 struct {
}

func (t U32) Name() string           { return "i32" }
func (t U32) Equal(u ValueType) bool { _, ok := u.(U32); return ok }

/**************************************
I64:
**************************************/
type I64 struct {
}

func (t I64) Name() string           { return "i64" }
func (t I64) Equal(u ValueType) bool { _, ok := u.(I64); return ok }

/**************************************
U64:
**************************************/
type U64 struct {
}

func (t U64) Name() string           { return "i64" }
func (t U64) Equal(u ValueType) bool { _, ok := u.(U64); return ok }

/**************************************
F32:
**************************************/
type F32 struct {
}

func (t F32) Name() string           { return "f32" }
func (t F32) Equal(u ValueType) bool { _, ok := u.(F32); return ok }

/**************************************
F64:
**************************************/
type F64 struct {
}

func (t F64) Name() string           { return "f64" }
func (t F64) Equal(u ValueType) bool { _, ok := u.(F64); return ok }

/**************************************
aVar:
**************************************/
type aVar struct {
	typ  ValueType
	name string
}

func (v *aVar) Name() string    { return v.name }
func (v *aVar) Type() ValueType { return v.typ }

/**************************************
VarI32:
**************************************/
type VarI32 struct {
	aVar
}

func NewVarI32(name string) *VarI32 {
	return &VarI32{aVar: aVar{typ: I32{}, name: name}}
}

/**************************************
VarU32:
**************************************/
type VarU32 struct {
	aVar
}

func NewVarU32(name string) *VarU32 {
	return &VarU32{aVar: aVar{typ: U32{}, name: name}}
}

/**************************************
VarI64:
**************************************/
type VarI64 struct {
	aVar
}

func NewVarI64(name string) *VarI64 {
	return &VarI64{aVar: aVar{typ: I64{}, name: name}}
}

/**************************************
VarU64:
**************************************/
type VarU64 struct {
	aVar
}

func NewVarU64(name string) *VarU64 {
	return &VarU64{aVar: aVar{typ: U64{}, name: name}}
}

/**************************************
VarF32:
**************************************/
type VarF32 struct {
	aVar
}

func NewVarF32(name string) *VarF32 {
	return &VarF32{aVar: aVar{typ: F32{}, name: name}}
}

/**************************************
VarF64:
**************************************/
type VarF64 struct {
	aVar
}

func NewVarF64(name string) *VarF64 {
	return &VarF64{aVar: aVar{typ: F64{}, name: name}}
}
