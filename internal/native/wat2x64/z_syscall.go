// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	_ "embed"

	"wa-lang.org/wa/internal/wat/ast"
)

func (p *wat2X64Worker) checkSyscallSig(spec *ast.ImportSpec) {
	// TODO: 检查系统调用函数签名类型是否匹配
}
