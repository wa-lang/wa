// 版权 @2024 凹语言 作者。保留所有权利。

//go:build !windows

package dlopen

import (
	"plugin"
)

var _ LibHandle = (*pluginLibHandle)(nil)

type pluginLibHandle struct {
	name   string
	handle *plugin.Plugin
}

type Proc struct {
	x plugin.Symbol
}

func open(path string) (LibHandle, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	h := &pluginLibHandle{
		name:   path,
		handle: p,
	}
	return h, nil
}

func (l *pluginLibHandle) Lookup(symName string) (*Proc, error) {
	x, err := l.handle.Lookup(symName)
	if err != nil {
		return nil, err
	}
	return &Proc{x}, nil
}
