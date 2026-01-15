// 版权 @2026 凹语言 作者。保留所有权利。

package ast

import "strings"

func (p *FuncType) String() string {
	var sb strings.Builder
	if len(p.Params) > 0 {
		sb.WriteString("(")
		for i, arg := range p.Params {
			if i > 0 {
				sb.WriteString(",")
			}
			if arg.Name != "" {
				sb.WriteString(arg.Name)
				sb.WriteString(":")
			}
			sb.WriteString(arg.Type.String())
		}
		sb.WriteString(")")
	}
	if len(p.Results) > 0 {
		sb.WriteString(" => ")
		if len(p.Results) == 1 {
			sb.WriteString(p.Results[0].String())
		} else {
			for i, ret := range p.Results {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(ret.String())
			}
			sb.WriteString(")")
		}
	}
	return sb.String()
}
