// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	_ "embed"
	"io"
)

const (
	kBuiltinMemcpy = "$builtin.memcpy"
	kBuiltinMemset = "$builtin.memset"
	kBuiltinPanic  = "$builtin.panic"
)

//go:embed z_builtin.was
var builtin_was string

// 生成系统调用代码
func (p *wat2X64Worker) buildBuiltin(w io.Writer) error {
	if _, err := w.Write([]byte(builtin_was)); err != nil {
		return err
	}
	return nil
}
