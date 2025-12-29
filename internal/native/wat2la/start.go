// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"
)

// 启动函数
func (p *wat2laWorker) buildStart(w io.Writer) error {
	fmt.Fprintln(w, "func _start {")
	{
		fmt.Fprintf(w, "    bl $wat2la.memory.init")
	}
	fmt.Fprintln(w, "}")

	return nil
}
