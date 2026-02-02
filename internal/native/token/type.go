// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package token

import "fmt"

// 数据类型
type Type int

const (
	I8 Type = iota + 1
	I16
	I32
	I64
	F32
	F64
)

func (t Type) IsIntType() bool {
	switch t {
	case I8, I16, I32, I64:
		return true
	default:
		return false
	}
}

func (t Type) IsFloatType() bool {
	switch t {
	case F32, F64:
		return true
	default:
		return false
	}
}

func (t Type) Size() int {
	switch t {
	case I8:
		return 1
	case I16:
		return 2
	case I32:
		return 4
	case I64:
		return 8
	case F32:
		return 4
	case F64:
		return 8
	default:
		panic("unreachable")
	}
}

func (t Type) String() string {
	switch t {
	case I8:
		return "i8"
	case I16:
		return "i16"
	case I32:
		return "i32"
	case I64:
		return "i64"
	case F32:
		return "f32"
	case F64:
		return "f64"
	default:
		return fmt.Sprintf("token.Type(%d)", int(t))
	}
}
