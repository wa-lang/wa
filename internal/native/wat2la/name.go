// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import "fmt"

func (p *wat2laWorker) getConstName(x uint64) string {
	return fmt.Sprintf(" $const.0x%x", x)
}

func (p *wat2laWorker) getLableName(label string) string {
	return fmt.Sprintf(" $label.%s", label)
}
