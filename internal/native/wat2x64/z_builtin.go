// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	_ "embed"
	"io"
)

const (
	kBuiltinMmap   = "$builtin.mmap"
	kBuiltinExit   = "$builtin.exit"
	kBuiltinPanic  = "$builtin.panic"
	kBuiltinWrite  = "$builtin.write"
	kBuiltinMemcpy = "$builtin.memcpy"
	kBuiltinMemset = "$builtin.memset"
)

//go:embed z_builtin-win64.wa.s.txt
var builtin_win64_wa_s string

//go:embed z_builtin-linux.wa.s.txt
var builtin_linux_wa_s string

// 生成系统调用代码
func (p *wat2X64Worker) buildBuiltinLinux(w io.Writer) error {
	if _, err := w.Write([]byte(builtin_linux_wa_s)); err != nil {
		return err
	}
	return nil
}
func (p *wat2X64Worker) buildBuiltinWindows(w io.Writer) error {
	if _, err := w.Write([]byte(builtin_win64_wa_s)); err != nil {
		return err
	}
	return nil
}
