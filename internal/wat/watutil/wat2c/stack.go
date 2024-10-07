// 版权 @2024 凹语言 作者。保留所有权利。

package wat2c

import (
	"strings"

	"wa-lang.org/wa/internal/wat/token"
)

// 函数内栈深度计算
type valueTypeStack struct {
	stack           []token.Token // i32/i64/f32/f64
	maxStackPointer int
}

func (s *valueTypeStack) String() string {
	var sb strings.Builder
	sb.WriteString("stack: [")
	for i, v := range s.stack {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("]")
	return sb.String()
}

func (s *valueTypeStack) TopIdx() int {
	if len(s.stack) == 0 {
		panic("unexpected stack empty")
	}
	return len(s.stack) - 1
}

func (s *valueTypeStack) Len() int {
	return len(s.stack)
}

func (s *valueTypeStack) MaxDepth() int {
	return s.maxStackPointer
}
func (s *valueTypeStack) Push(v token.Token) {
	switch v {
	case token.I32, token.I64, token.F32, token.F64:
	default:
		panic("unexpected value type")
	}
	s.stack = append(s.stack, v)
	if sp := len(s.stack); sp > s.maxStackPointer {
		s.maxStackPointer = sp
	}
}

func (s *valueTypeStack) Pop(expect token.Token) {
	switch expect {
	case token.I32, token.I64, token.F32, token.F64:
	default:
		panic("unexpected value type")
	}

	if len(s.stack) == 0 {
		panic("unexpected stack empty")
	}
	if got := s.stack[len(s.stack)-1]; got != expect {
		panic("unexpected value type: got " + got.String() + ", expect " + expect.String())
	}
	s.stack = s.stack[:len(s.stack)-1]
	return
}

func (s *valueTypeStack) Drop() {
	if len(s.stack) == 0 {
		panic("unexpected stack empty")
	}
	s.stack = s.stack[:len(s.stack)-1]
	return
}

// todo: 删除
func (s *valueTypeStack) PushN(n int) {
	for i := 0; i < n; i++ {
		s.stack = append(s.stack, 0)
	}
	if sp := len(s.stack); sp > s.maxStackPointer {
		s.maxStackPointer = sp
	}
}

// todo: 删除
func (s *valueTypeStack) PopN(dx int) {
	if len(s.stack) < dx {
		panic("unexpected stack empty")
	}
	s.stack = s.stack[:len(s.stack)-dx]
	return
}

func (s *valueTypeStack) Nop() {
	// 栈不变化
}
