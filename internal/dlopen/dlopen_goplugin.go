// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

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
