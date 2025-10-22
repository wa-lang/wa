// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package lsp

import (
	"strings"

	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/lsp/protocol"
)

// 一个模块类似一个工作区
type WaModule struct {
	manifest *config.Manifest
	program  *loader.Program
}

func (p *LSPServer) loadModule(uri protocol.DocumentURI) (*WaModule, bool) {
	if !strings.HasSuffix(string(uri), ".wa") {
		return nil, false
	}
	if m, ok := p.waModules[uri]; ok {
		return m, true
	}
	m, err := config.LoadManifest(nil, uri.Path(), false)
	module := &WaModule{manifest: m}
	p.waModules[uri] = module
	return module, err == nil
}

func (p *WaModule) LoadProgram() error {
	if p.program != nil {
		return nil
	}
	prog, err := loader.LoadProgram(config.DefaultConfig(), p.manifest.Root)
	if err != nil {
		return err
	}
	p.program = prog
	return nil
}
