package lltypes

import "testing"

func TestIntTypeEqual(t *testing.T) {
	golden := []struct {
		t    *IntType
		u    *IntType
		want bool
	}{
		{
			t:    &IntType{BitSize: 32},
			u:    &IntType{BitSize: 8},
			want: false,
		},
		{
			t:    &IntType{BitSize: 32},
			u:    I32,
			want: true,
		},
		{
			t:    &IntType{BitSize: 32},
			u:    NewInt(32),
			want: true,
		},
		{
			t:    &IntType{TypeName: "foo", BitSize: 32},
			u:    &IntType{TypeName: "foo", BitSize: 8},
			want: false,
		},
		{
			t:    &IntType{TypeName: "foo", BitSize: 32},
			u:    I32,
			want: true,
		},
		{
			t:    &IntType{TypeName: "foo", BitSize: 32},
			u:    NewInt(32),
			want: true,
		},
	}
	for _, g := range golden {
		got := g.t.Equal(g.u)
		if g.want != got {
			t.Errorf("struct equality mismatch between `%s` and `%s`; expected %t, got %t", g.t.LLString(), g.u.LLString(), g.want, got)
		}
	}
}

func TestIsVoid(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{&VoidType{}, true},
		{Void, true},
		{I8, false},
	}
	for _, g := range golden {
		got := IsVoid(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is a void type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestIsFunc(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{&FuncType{RetType: Void}, true},
		{NewFunc(Void), true},
		{NewFunc(I8, I8), true},
		{I8, false},
	}
	for _, g := range golden {
		got := IsFunc(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is a function type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestIsInt(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{t: &IntType{BitSize: 123}, want: true},
		{t: NewInt(123), want: true},
		{t: I1, want: true},
		{t: I8, want: true},
		{t: I16, want: true},
		{t: I32, want: true},
		{t: I64, want: true},
		{t: I128, want: true},
		{t: Void, want: false},
		{t: Double, want: false},
	}
	for _, g := range golden {
		got := IsInt(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is an integer type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestIsFloat(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{t: &FloatType{Kind: FloatKindDouble}, want: true},
		{t: Half, want: true},
		{t: Float, want: true},
		{t: Double, want: true},
		{t: FP128, want: true},
		{t: X86_FP80, want: true},
		{t: PPC_FP128, want: true},
		{t: I8, want: false},
	}
	for _, g := range golden {
		got := IsFloat(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is a floating-point type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestIsMMX(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{t: &MMXType{}, want: true},
		{t: MMX, want: true},
		{t: I8, want: false},
	}
	for _, g := range golden {
		got := IsMMX(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is an MMX type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestIsPointer(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{t: &PointerType{ElemType: I8}, want: true},
		{t: NewPointer(I8), want: true},
		{t: I8, want: false},
	}
	for _, g := range golden {
		got := IsPointer(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is a pointer type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestIsVector(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{t: &VectorType{Len: 5, ElemType: I8}, want: true},
		{t: NewVector(5, I8), want: true},
		{t: I8, want: false},
	}
	for _, g := range golden {
		got := IsVector(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is a vector type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestIsLabel(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{t: &LabelType{}, want: true},
		{t: Label, want: true},
		{t: I8, want: false},
	}
	for _, g := range golden {
		got := IsLabel(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is a label type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestIsToken(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{t: &TokenType{}, want: true},
		{t: Token, want: true},
		{t: I8, want: false},
	}
	for _, g := range golden {
		got := IsToken(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is a token type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestIsMetadata(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{t: &MetadataType{}, want: true},
		{t: Metadata, want: true},
		{t: I8, want: false},
	}
	for _, g := range golden {
		got := IsMetadata(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is a metadata type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestIsArray(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{t: &ArrayType{Len: 5, ElemType: I8}, want: true},
		{t: NewArray(5, I8), want: true},
		{t: I8, want: false},
	}
	for _, g := range golden {
		got := IsArray(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is an array type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestIsStruct(t *testing.T) {
	golden := []struct {
		t    Type
		want bool
	}{
		{t: &StructType{TypeName: "foo", Fields: []Type{I8}}, want: true},
		{t: &StructType{Fields: []Type{I8}}, want: true},
		{t: NewStruct(I8, I32), want: true},
		{t: NewStruct(I8), want: true},
		{t: NewStruct(), want: true},
		{t: I8, want: false},
	}
	for _, g := range golden {
		got := IsStruct(g.t)
		if g.want != got {
			t.Errorf("check if `%s` is an array type mismatch; expected %t, got %t", g.t, g.want, got)
		}
	}
}

func TestEqual(t *testing.T) {
	golden := []struct {
		t    Type
		u    Type
		want bool
	}{
		{t: Void, u: &VoidType{}, want: true},
		{t: Void, u: I8, want: false},
		{t: NewFunc(Void), u: NewFunc(Void), want: true},
		{t: NewFunc(Void, I32), u: NewFunc(Void, I8), want: false},
		{t: &FuncType{RetType: Void, Variadic: true}, u: &FuncType{RetType: Void}, want: false},
		{t: NewFunc(Void), u: I8, want: false},
		{t: I8, u: NewInt(8), want: true},
		{t: I8, u: NewInt(9), want: false},
		{t: I8, u: Double, want: false},
		{t: &FloatType{Kind: FloatKindDouble}, u: Double, want: true},
		{t: Float, u: Double, want: false},
		{t: Float, u: I8, want: false},
		{t: MMX, u: &MMXType{}, want: true},
		{t: MMX, u: I8, want: false},
		{t: NewPointer(I8), u: &PointerType{ElemType: I8}, want: true},
		{t: NewPointer(I8), u: NewPointer(Double), want: false},
		{t: NewPointer(I8), u: I8, want: false},
		{t: NewVector(5, I8), u: &VectorType{Len: 5, ElemType: I8}, want: true},
		{t: NewVector(5, I8), u: NewVector(3, I8), want: false},
		{t: NewVector(5, I8), u: I8, want: false},
		{t: Label, u: &LabelType{}, want: true},
		{t: Label, u: I8, want: false},
		{t: Token, u: &TokenType{}, want: true},
		{t: Token, u: I8, want: false},
		{t: Metadata, u: &MetadataType{}, want: true},
		{t: Metadata, u: I8, want: false},
		{t: NewArray(5, I8), u: &ArrayType{Len: 5, ElemType: I8}, want: true},
		{t: NewArray(5, I8), u: NewArray(3, I8), want: false},
		{t: NewArray(5, I8), u: I8, want: false},
	}
	for _, g := range golden {
		got := Equal(g.t, g.u)
		if g.want != got {
			t.Errorf("equality mismatch between `%s` and `%s`; expected %t, got %t", g.t.LLString(), g.u.LLString(), g.want, got)
		}
	}
}

func TestStructTypeEqual(t *testing.T) {
	// Identified (named) struct types are uniqued by type names, not by
	// structural identity.
	golden := []struct {
		t    *StructType
		u    *StructType
		want bool
	}{
		// Unnamed struct types.
		{
			t:    &StructType{Fields: []Type{I32}},
			u:    NewStruct(I8),
			want: false,
		},
		{
			t:    &StructType{Fields: []Type{I32}, Packed: false},
			u:    &StructType{Fields: []Type{I32}, Packed: true},
			want: false,
		},
		{
			t:    &StructType{Fields: []Type{I32}},
			u:    &StructType{Fields: []Type{I32}},
			want: true,
		},
		{
			t:    &StructType{Fields: []Type{I32}, Packed: true},
			u:    &StructType{Fields: []Type{I32}, Packed: true},
			want: true,
		},
		// Identified struct types.
		{
			t:    &StructType{TypeName: "foo", Fields: []Type{I32}},
			u:    &StructType{TypeName: "bar", Fields: []Type{I8}},
			want: false,
		},
		{
			t:    &StructType{TypeName: "foo", Fields: []Type{I32}, Packed: false},
			u:    &StructType{TypeName: "bar", Fields: []Type{I32}, Packed: true},
			want: false,
		},
		{
			t:    &StructType{TypeName: "foo", Fields: []Type{I32}},
			u:    &StructType{TypeName: "foo", Fields: []Type{I32}},
			want: true,
		},
		{
			t:    &StructType{TypeName: "foo", Fields: []Type{I32}, Packed: true},
			u:    &StructType{TypeName: "foo", Fields: []Type{I32}, Packed: true},
			want: true,
		},
		{
			t:    &StructType{TypeName: "foo", Fields: []Type{I32}},
			u:    &StructType{TypeName: "bar", Fields: []Type{I32}},
			want: false,
		},
		{
			t:    &StructType{TypeName: "foo", Fields: []Type{I32}, Packed: true},
			u:    &StructType{TypeName: "bar", Fields: []Type{I32}, Packed: true},
			want: false,
		},
	}
	for _, g := range golden {
		got := g.t.Equal(g.u)
		if g.want != got {
			t.Errorf("struct equality mismatch between `%s` and `%s`; expected %t, got %t", g.t.LLString(), g.u.LLString(), g.want, got)
		}
	}
}

// Assert that each type implements the types.Type interface.
var (
	_ Type = (*VoidType)(nil)
	_ Type = (*FuncType)(nil)
	_ Type = (*IntType)(nil)
	_ Type = (*FloatType)(nil)
	_ Type = (*MMXType)(nil)
	_ Type = (*PointerType)(nil)
	_ Type = (*VectorType)(nil)
	_ Type = (*LabelType)(nil)
	_ Type = (*TokenType)(nil)
	_ Type = (*MetadataType)(nil)
	_ Type = (*ArrayType)(nil)
	_ Type = (*StructType)(nil)
)
