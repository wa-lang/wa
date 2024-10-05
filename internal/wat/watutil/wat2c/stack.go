// 版权 @2024 凹语言 作者。保留所有权利。

package wat2c

import (
	"wa-lang.org/wa/internal/wat/token"
)

// 函数内栈深度计算
type valueTypeStack struct {
	stack           []token.Token // i32/i64/f32/f64
	maxStackPointer int
}

func (s *valueTypeStack) TopIdx() int {
	return len(s.stack) - 1
}

func (s *valueTypeStack) Len() int {
	return len(s.stack)
}

func (s *valueTypeStack) MaxDepth() int {
	return s.maxStackPointer
}
func (s *valueTypeStack) Push(v token.Token) {
	s.stack = append(s.stack, v)
	if sp := len(s.stack); sp > s.maxStackPointer {
		s.maxStackPointer = sp
	}
}

func (s *valueTypeStack) Pop() (vt token.Token) {
	if len(s.stack) == 0 {
		return // todo: panic
	}
	vt = s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return
}

func (s *valueTypeStack) PushN(n int) {
	for i := 0; i < n; i++ {
		s.stack = append(s.stack, 0)
	}
	if sp := len(s.stack); sp > s.maxStackPointer {
		s.maxStackPointer = sp
	}
}
func (s *valueTypeStack) PopN(dx int) {
	if len(s.stack) < dx {
		return // todo: panic
	}
	s.stack = s.stack[:len(s.stack)-dx]
	return
}

func (s *valueTypeStack) Nop() {
	// 栈不变化
}
