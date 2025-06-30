// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package dlopen

import (
	"fmt"
	"syscall"
	_ "syscall"
)

var _ LibHandle = (*dllLibHandle)(nil)

type dllLibHandle struct {
	handle *syscall.LazyDLL
}

type Proc struct {
	x *syscall.LazyProc
}

func open(path string) (LibHandle, error) {
	dll := syscall.NewLazyDLL(path)
	if err := dll.Load(); err != nil {
		return nil, err
	}
	return &dllLibHandle{dll}, nil
}

func (l *dllLibHandle) Lookup(symName string) (*Proc, error) {
	proc := l.handle.NewProc(symName)
	if proc == nil {
		return nil, fmt.Errorf("error resolving symbol %q", symName)
	}

	return &Proc{proc}, nil
}
