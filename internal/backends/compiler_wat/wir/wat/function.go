// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "strings"

func (f *Function) Format(sb *strings.Builder) {
	sb.WriteString("(func $")
	sb.WriteString(f.InternalName)

	if len(f.ExternalName) > 0 {
		sb.WriteString(" (export \"")
		sb.WriteString(f.ExternalName)
		sb.WriteString("\")")
	}

	for _, param := range f.Params {
		sb.WriteString(" (param $")
		sb.WriteString(param.Name())
		sb.WriteByte(' ')
		sb.WriteString(param.Type().Name())
		sb.WriteByte(')')
	}

	if len(f.Results) > 0 {
		sb.WriteString(" (result")
		for _, r := range f.Results {
			sb.WriteByte(' ')
			sb.WriteString(r.Name())
		}
		sb.WriteByte(')')
	}
	sb.WriteByte('\n')

	for _, local := range f.Locals {
		sb.WriteString("  (local $")
		sb.WriteString(local.Name())
		sb.WriteByte(' ')
		sb.WriteString(local.Type().Name())
		sb.WriteByte(')')
		sb.WriteByte('\n')
	}

	indent_t := "  "
	for _, inst := range f.Insts {
		inst.Format(indent_t, sb)
		sb.WriteByte('\n')
	}

	sb.WriteString(") ;;")
	sb.WriteString(f.InternalName)
}

func (sig *FuncSig) String() string {
	var sb strings.Builder

	for _, param := range sig.Params {
		sb.WriteString(" (param ")
		sb.WriteString(param.Name())
		sb.WriteByte(')')
	}

	for _, ret := range sig.Results {
		sb.WriteString(" (result ")
		sb.WriteString(ret.Name())
		sb.WriteByte(')')
	}
	return sb.String()
}
