// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"testing"

	"wa-lang.org/wa/internal/p9asm/arch"
	"wa-lang.org/wa/internal/p9asm/lex"
	"wa-lang.org/wa/internal/p9asm/obj"
)

// A simple in-out test: Do we print what we parse?

func setArch(waarch string) (*arch.Arch, *obj.Link) {
	architecture := arch.Set(waarch)
	if architecture == nil {
		panic("asm: unrecognized architecture " + waarch)
	}
	return architecture, obj.Linknew(architecture.LinkArch)
}

func newParser(waarch string) *Parser {
	architecture, ctxt := setArch(waarch)
	return NewParser(ctxt, architecture, nil)
}

func testOperandParser(t *testing.T, parser *Parser, tests []operandTest) {
	for _, test := range tests {
		parser.start(lex.Tokenize(test.input))
		addr := obj.Addr{}
		parser.operand(&addr)
		result := obj.Dconv(&emptyProg, &addr)
		if result != test.output {
			t.Errorf("fail at %s: got %s; expected %s\n", test.input, result, test.output)
		}
	}
}

func TestAMD64OperandParser(t *testing.T) {
	parser := newParser("amd64")
	testOperandParser(t, parser, amd64OperandTests)
}

func TestARM64OperandParser(t *testing.T) {
	parser := newParser("arm64")
	testOperandParser(t, parser, arm64OperandTests)
}

type operandTest struct {
	input, output string
}

// Examples collected by scanning all the assembly in the standard repo.

var amd64OperandTests = []operandTest{
	{"$(-1.0)", "$(-1.0)"},
	{"$(0.0)", "$(0.0)"},
	{"$(0x2000000+116)", "$33554548"},
	{"$(0x3F<<7)", "$8064"},
	{"$(112+8)", "$120"},
	{"$(1<<63)", "$-9223372036854775808"},
	{"$-1", "$-1"},
	{"$0", "$0"},
	{"$0-0", "$0"},
	{"$0-16", "$-16"},
	{"$0x000FFFFFFFFFFFFF", "$4503599627370495"},
	{"$0x01", "$1"},
	{"$0x02", "$2"},
	{"$0x04", "$4"},
	{"$0x3FE", "$1022"},
	{"$0x7fffffe00000", "$140737486258176"},
	{"$0xfffffffffffff001", "$-4095"},
	{"$1", "$1"},
	{"$1.0", "$(1.0)"},
	{"$10", "$10"},
	{"$1000", "$1000"},
	{"$1000000", "$1000000"},
	{"$1000000000", "$1000000000"},
	{"$__tsan_func_enter(SB)", "$__tsan_func_enter(SB)"},
	{"$main(SB)", "$main(SB)"},
	{"$masks<>(SB)", "$masks<>(SB)"},
	{"$setg_gcc<>(SB)", "$setg_gcc<>(SB)"},
	{"$shifts<>(SB)", "$shifts<>(SB)"},
	{"$~(1<<63)", "$9223372036854775807"},
	{"$~0x3F", "$-64"},
	{"$~15", "$-16"},
	{"(((8)&0xf)*4)(SP)", "32(SP)"},
	{"(((8-14)&0xf)*4)(SP)", "40(SP)"},
	{"(6+8)(AX)", "14(AX)"},
	{"(8*4)(BP)", "32(BP)"},
	{"(AX)", "(AX)"},
	{"(AX)(CX*8)", "(AX)(CX*8)"},
	{"(BP)(CX*4)", "(BP)(CX*4)"},
	{"(BP)(DX*4)", "(BP)(DX*4)"},
	{"(BP)(R8*4)", "(BP)(R8*4)"},
	{"(BX)", "(BX)"},
	{"(DI)", "(DI)"},
	{"(DI)(BX*1)", "(DI)(BX*1)"},
	{"(DX)", "(DX)"},
	{"(R9)", "(R9)"},
	{"(R9)(BX*8)", "(R9)(BX*8)"},
	{"(SI)", "(SI)"},
	{"(SI)(BX*1)", "(SI)(BX*1)"},
	{"(SI)(DX*1)", "(SI)(DX*1)"},
	{"(SP)", "(SP)"},
	{"+3(PC)", "3(PC)"},
	{"-1(DI)(BX*1)", "-1(DI)(BX*1)"},
	{"-3(PC)", "-3(PC)"},
	{"-64(SI)(BX*1)", "-64(SI)(BX*1)"},
	{"-96(SI)(BX*1)", "-96(SI)(BX*1)"},
	{"AL", "AL"},
	{"AX", "AX"},
	{"BP", "BP"},
	{"BX", "BX"},
	{"CX", "CX"},
	{"DI", "DI"},
	{"DX", "DX"},
	{"R10", "R10"},
	{"R10", "R10"},
	{"R11", "R11"},
	{"R12", "R12"},
	{"R13", "R13"},
	{"R14", "R14"},
	{"R15", "R15"},
	{"R8", "R8"},
	{"R9", "R9"},
	{"SI", "SI"},
	{"SP", "SP"},
	{"X0", "X0"},
	{"X1", "X1"},
	{"X10", "X10"},
	{"X11", "X11"},
	{"X12", "X12"},
	{"X13", "X13"},
	{"X14", "X14"},
	{"X15", "X15"},
	{"X2", "X2"},
	{"X3", "X3"},
	{"X4", "X4"},
	{"X5", "X5"},
	{"X6", "X6"},
	{"X7", "X7"},
	{"X8", "X8"},
	{"X9", "X9"},
	{"_expand_key_128<>(SB)", "_expand_key_128<>(SB)"},
	{"_seek<>(SB)", "_seek<>(SB)"},
	{"a2+16(FP)", "a2+16(FP)"},
	{"addr2+24(FP)", "addr2+24(FP)"},
	{"asmcgocall<>(SB)", "asmcgocall<>(SB)"},
	{"b+24(FP)", "b+24(FP)"},
	{"b_len+32(FP)", "b_len+32(FP)"},
	{"racecall<>(SB)", "racecall<>(SB)"},
	{"rcv_name+20(FP)", "rcv_name+20(FP)"},
	{"retoffset+28(FP)", "retoffset+28(FP)"},
	{"runtime·_GetStdHandle(SB)", "runtime._GetStdHandle(SB)"},
	{"sync\u2215atomic·AddInt64(SB)", "sync/atomic.AddInt64(SB)"},
	{"timeout+20(FP)", "timeout+20(FP)"},
	{"ts+16(FP)", "ts+16(FP)"},
	{"x+24(FP)", "x+24(FP)"},
	{"x·y(SB)", "x.y(SB)"},
	{"x·y(SP)", "x.y(SP)"},
	{"x·y+8(SB)", "x.y+8(SB)"},
	{"x·y+8(SP)", "x.y+8(SP)"},
	{"y+56(FP)", "y+56(FP)"},
	{"·AddUint32(SB", "\"\".AddUint32(SB)"},
	{"·callReflect(SB)", "\"\".callReflect(SB)"},
}

var x86OperandTests = []operandTest{
	{"$(2.928932188134524e-01)", "$(0.29289321881345243)"},
	{"$-1", "$-1"},
	{"$0", "$0"},
	{"$0x00000000", "$0"},
	{"$runtime·badmcall(SB)", "$runtime.badmcall(SB)"},
	{"$setg_gcc<>(SB)", "$setg_gcc<>(SB)"},
	{"$~15", "$-16"},
	{"(-64*1024+104)(SP)", "-65432(SP)"},
	{"(0*4)(BP)", "(BP)"},
	{"(1*4)(DI)", "4(DI)"},
	{"(4*4)(BP)", "16(BP)"},
	{"(AX)", "(AX)"},
	{"(BP)(CX*4)", "(BP)(CX*4)"},
	{"(BP*8)", "0(BP*8)"},
	{"(BX)", "(BX)"},
	{"(SP)", "(SP)"},
	{"*AX", "AX"}, // TODO: Should make * illegal here; a simple alias for JMP AX.
	{"*runtime·_GetStdHandle(SB)", "*runtime._GetStdHandle(SB)"},
	{"-(4+12)(DI)", "-16(DI)"},
	{"-1(DI)(BX*1)", "-1(DI)(BX*1)"},
	{"-96(DI)(BX*1)", "-96(DI)(BX*1)"},
	{"0(AX)", "(AX)"},
	{"0(BP)", "(BP)"},
	{"0(BX)", "(BX)"},
	{"4(AX)", "4(AX)"},
	{"AL", "AL"},
	{"AX", "AX"},
	{"BP", "BP"},
	{"BX", "BX"},
	{"CX", "CX"},
	{"DI", "DI"},
	{"DX", "DX"},
	{"F0", "F0"},
	{"GS", "GS"},
	{"SI", "SI"},
	{"SP", "SP"},
	{"X0", "X0"},
	{"X1", "X1"},
	{"X2", "X2"},
	{"X3", "X3"},
	{"X4", "X4"},
	{"X5", "X5"},
	{"X6", "X6"},
	{"X7", "X7"},
	{"asmcgocall<>(SB)", "asmcgocall<>(SB)"},
	{"ax+4(FP)", "ax+4(FP)"},
	{"ptime-12(SP)", "ptime-12(SP)"},
	{"runtime·_NtWaitForSingleObject(SB)", "runtime._NtWaitForSingleObject(SB)"},
	{"s(FP)", "s(FP)"},
	{"sec+4(FP)", "sec+4(FP)"},
	{"shifts<>(SB)(CX*8)", "shifts<>(SB)(CX*8)"},
	{"x+4(FP)", "x+4(FP)"},
	{"·AddUint32(SB)", "\"\".AddUint32(SB)"},
	{"·reflectcall(SB)", "\"\".reflectcall(SB)"},
}

var arm64OperandTests = []operandTest{
	{"$0", "$0"},
	{"$0.5", "$(0.5)"},
	{"0(R26)", "(R26)"},
	{"0(RSP)", "(RSP)"},
	{"$1", "$1"},
	{"$-1", "$-1"},
	{"$1000", "$1000"},
	{"$1000000000", "$1000000000"},
	{"$0x7fff3c000", "$34358935552"},
	{"$1234", "$1234"},
	{"$~15", "$-16"},
	{"$16", "$16"},
	{"-16(RSP)", "-16(RSP)"},
	{"16(RSP)", "16(RSP)"},
	{"1(R1)", "1(R1)"},
	{"-1(R4)", "-1(R4)"},
	{"18740(R5)", "18740(R5)"},
	{"$2", "$2"},
	{"$-24(R4)", "$-24(R4)"},
	{"-24(RSP)", "-24(RSP)"},
	{"$24(RSP)", "$24(RSP)"},
	{"-32(RSP)", "-32(RSP)"},
	{"$48", "$48"},
	{"$(-64*1024)(R7)", "$-65536(R7)"},
	{"$(8-1)", "$7"},
	{"a+0(FP)", "a(FP)"},
	{"a1+8(FP)", "a1+8(FP)"},
	{"·AddInt32(SB)", `"".AddInt32(SB)`},
	{"runtime·divWVW(SB)", "runtime.divWVW(SB)"},
	{"$argframe+0(FP)", "$argframe(FP)"},
	{"$asmcgocall<>(SB)", "$asmcgocall<>(SB)"},
	{"EQ", "EQ"},
	{"F29", "F29"},
	{"F3", "F3"},
	{"F30", "F30"},
	{"g", "g"},
	{"LR", "R30"},
	{"(LR)", "(R30)"},
	{"R0", "R0"},
	{"R10", "R10"},
	{"R11", "R11"},
	{"$4503601774854144.0", "$(4503601774854144.0)"},
	{"$runtime·badsystemstack(SB)", "$runtime.badsystemstack(SB)"},
	{"ZR", "ZR"},
	{"(ZR)", "(ZR)"},
	{"(R29, RSP)", "(R29, RSP)"},
}
