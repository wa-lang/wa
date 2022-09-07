// 版权 @2022 凹语言 作者。保留所有权利。

package wat

/**************************************
I32:
**************************************/
type I32 struct {
}

func (t I32) Name() string { return "i32" }
func (t I32) Equal(u ValueType) bool {
	if _, ok := u.(I32); ok {
		return true
	}
	return false
}

/**************************************
I64:
**************************************/
type I64 struct {
}

func (t I64) Name() string { return "i64" }
func (t I64) Equal(u ValueType) bool {
	if _, ok := u.(I64); ok {
		return true
	}
	return false
}

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
