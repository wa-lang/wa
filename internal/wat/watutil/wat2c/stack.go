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

func (s *valueTypeStack) maxDepth() int {
	return s.maxStackPointer
}
func (s *valueTypeStack) push(v token.Token) {
	s.stack = append(s.stack, v)
	if sp := len(s.stack); sp > s.maxStackPointer {
		s.maxStackPointer = sp
	}
}

func (s *valueTypeStack) pop() (vt token.Token) {
	vt = s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return
}

func (s *valueTypeStack) pushN(n int) {
	for i := 0; i < n; i++ {
		s.stack = append(s.stack, 0)
	}
	if sp := len(s.stack); sp > s.maxStackPointer {
		s.maxStackPointer = sp
	}
}
func (s *valueTypeStack) popN(dx int) {
	s.stack = s.stack[:len(s.stack)-dx]
	return
}

func (s *valueTypeStack) nop() {
	// 栈不变化
}
