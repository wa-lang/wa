// 版权 @2024 凹语言 作者。保留所有权利。

package wasm

import "strconv"

var instructionOpcodes map[string]Opcode

func init() {
	instructionOpcodes = make(map[string]Opcode)

	for i, s := range instructionNames {
		if s != "" {
			instructionOpcodes[s] = Opcode(i)
		}
	}
	for i, s := range miscInstructionNames {
		if s != "" {
			instructionOpcodes[s] = Opcode(i)
		}
	}
	for i, s := range vectorInstructionName {
		if s != "" {
			instructionOpcodes[s] = Opcode(i)
		}
	}
}

func LookupOpcode(s string) (Opcode, bool) {
	v, ok := instructionOpcodes[s]
	return v, ok
}

func OpcodeName(op Opcode) string {
	s := instructionNames[op]
	if s == "" {
		s = "wasm.Opcode(" + strconv.Itoa(int(op)) + ")"
	}
	return s
}
