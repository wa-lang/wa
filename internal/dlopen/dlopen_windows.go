// 版权 @2024 凹语言 作者。保留所有权利。

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
