// 版权 @2024 凹语言 作者。保留所有权利。

package wat2c

import (
	"fmt"
	"strings"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

// 函数内栈类型和深度计算
type valueTypeStack struct {
	trace              bool            // 是否跟踪
	funcName           string          // 函数名
	funcInstruction    ast.Instruction // 当前指令
	funcInstructionPos int             // 函数指令位置
	stack              []token.Token   // i32/i64/f32/f64/funcref
	maxStackPointer    int
}

func (s *valueTypeStack) Len() int {
	return len(s.stack)
}

func (s *valueTypeStack) MaxDepth() int {
	return s.maxStackPointer
}

func (s *valueTypeStack) NextInstruction(ins ast.Instruction) {
	s.funcInstruction = ins
	s.funcInstructionPos++
	if s.trace {
		fmt.Printf("wat2c.valueTypeStack.Ins : ------ %-*s;; %03d:%v\n",
			40, s.StackString(), s.funcInstructionPos, s.funcInstruction,
		)
	}
}

func (s *valueTypeStack) LastInstruction() ast.Instruction {
	return s.funcInstruction
}

func (s *valueTypeStack) Top(expect token.Token) int {
	if len(s.stack) == 0 {
		panic("unexpected stack empty; # " + s.String())
	}
	idx := len(s.stack) - 1
	if got := s.stack[idx]; got != expect {
		panic(
			"unexpected value type: got " +
				got.String() + ", expect " + expect.String() +
				"; # " + s.String(),
		)
	}
	return idx
}
func (s *valueTypeStack) TopToken() token.Token {
	if len(s.stack) == 0 {
		panic("unexpected stack empty; # " + s.String())
	}
	idx := len(s.stack) - 1
	return s.stack[idx]
}

func (s *valueTypeStack) Push(v token.Token) int {
	if s.trace {
		fmt.Printf("wat2c.valueTypeStack.Push: %v => %-*s;; %03d:%v\n",
			v, 40, s.StackString(), s.funcInstructionPos, s.funcInstruction,
		)
	}
	switch v {
	case token.I32, token.I64, token.F32, token.F64, token.FUNCREF:
	default:
		panic("unexpected value type; # " + s.String())
	}
	s.stack = append(s.stack, v)
	if sp := len(s.stack); sp > s.maxStackPointer {
		s.maxStackPointer = sp
	}
	return len(s.stack) - 1
}

func (s *valueTypeStack) Pop(expect token.Token) int {
	if s.trace {
		fmt.Printf("wat2c.valueTypeStack.Pop : %v <= %-*s;; %03d:%v\n",
			expect, 40, s.StackString(), s.funcInstructionPos, s.funcInstruction,
		)
	}
	switch expect {
	case token.I32, token.I64, token.F32, token.F64, token.FUNCREF:
	default:
		panic("unexpected value type; # " + s.String())
	}

	if len(s.stack) == 0 {
		panic("unexpected stack empty; # " + s.String())
	}
	idx := len(s.stack) - 1
	if got := s.stack[idx]; got != expect {
		panic(
			"unexpected value type: got " +
				got.String() + ", expect " + expect.String() +
				"; # " + s.String(),
		)
	}
	s.stack = s.stack[:len(s.stack)-1]
	return idx
}

func (s *valueTypeStack) DropAny() int {
	if s.trace {
		fmt.Printf("wat2c.valueTypeStack.DropAny: ? <= %-*s;; %03d:%v\n",
			40, s.StackString(), s.funcInstructionPos, s.funcInstruction,
		)
	}
	if len(s.stack) == 0 {
		panic("unexpected stack empty; # " + s.String())
	}
	idx := len(s.stack) - 1
	s.stack = s.stack[:len(s.stack)-1]
	return idx
}

func (s *valueTypeStack) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s:ins[%v@%d]:", s.funcName, s.funcInstruction, s.funcInstructionPos))
	sb.WriteString("stack:[")
	for i, v := range s.stack {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("]")
	return sb.String()
}

func (s *valueTypeStack) StackString() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range s.stack {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("]")
	return sb.String()
}
