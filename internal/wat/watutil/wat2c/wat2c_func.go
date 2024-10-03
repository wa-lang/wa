// 版权 @2024 凹语言 作者。保留所有权利。

package wat2c

import (
	"fmt"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2cWorker) buildFunc_body(fn *ast.Func) error {
	if len(fn.Body.Locals) > 0 {
		for _, x := range fn.Body.Locals {
			fmt.Fprintf(&p.c, "\t// local %s: %s\n", x.Name, x.Type)
			fmt.Fprintf(&p.c, "\t%s var_%s = 0;\n", toCType(x.Type), toCName(x.Name))
		}
		fmt.Fprintln(&p.c)
	}

	for _, ins := range fn.Body.Insts {
		if err := p.buildFunc_ins(fn, ins, 0); err != nil {
			return err
		}
	}
	return nil
}

func (p *wat2cWorker) buildFunc_ins(fn *ast.Func, i ast.Instruction, level int) error {
	switch i.Token() {
	case token.INS_UNREACHABLE:
		fmt.Fprintf(&p.c, "\tabort(); // unreachable\n")
	case token.INS_NOP:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_BLOCK:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_LOOP:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_IF:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)

	case token.INS_ELSE:
		panic("unreachable")
	case token.INS_END:
		panic("unreachable")

	case token.INS_BR:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_BR_IF:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_BR_TABLE:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_RETURN:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_CALL:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_CALL_INDIRECT:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_DROP:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_SELECT:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_TYPED_SELECT:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_LOCAL_GET:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_LOCAL_SET:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_LOCAL_TEE:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_GLOBAL_GET:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_GLOBAL_SET:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_TABLE_GET:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_TABLE_SET:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)

	case token.INS_I32_LOAD:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_I64_LOAD:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_F32_LOAD:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_F64_LOAD:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	case token.INS_I32_LOAD8_S:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)

	default:
		fmt.Fprintf(&p.c, "\t// todo: %T\n", i)
	}
	return nil
}
