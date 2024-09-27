// 版权 @2024 凹语言 作者。保留所有权利。

package ast

import (
	"fmt"
)

func (m *Module) String() string {
	if m.Name != "" {
		return fmt.Sprintf("(module $%s)", m.Name)
	}
	return "(module)"
}
