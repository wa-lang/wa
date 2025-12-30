// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	_ "embed"
	"io"

	"wa-lang.org/wa/internal/wat/ast"
)

const (
	kSysWrite = "$syscall.write"
	kSysExit  = "$syscall.exit"
	kSysBrk   = "$syscall.brk"
	kSysMmap  = "$syscall.mmap"
)

//go:embed z_syscall.was
var syscall_was string

//go:embed z_syscall.wzs
var syscall_wzs string

// 生成系统调用代码
func (p *wat2laWorker) buildSyscall(w io.Writer, hostFuncMap map[string]bool) error {
	if len(hostFuncMap) == 0 {
		return nil
	}
	if _, err := w.Write([]byte(syscall_was)); err != nil {
		return err
	}
	if _, err := w.Write([]byte(syscall_wzs)); err != nil {
		return err
	}
	return nil
}

func (p *wat2laWorker) checkSyscallSig(spec *ast.ImportSpec) {
	// TODO: 检查系统调用函数签名类型是否匹配
}
