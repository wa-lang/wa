// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package spinner

import (
	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/types"
	"wa-lang.org/wa/internal/wire"
)

/**************************************
本包用于将 AST 转为 wire
**************************************/

//-------------------------------------

type Spinner struct {
	Module   wire.Module
	manifest *config.Manifest

	mainPkg *Package

	packages map[*types.Package]*Package

	typeTable map[string]*wire.Type

	funcsDone map[types.Object]*wire.Function
	funcsTodo map[*types.Func]bool

	typesNeedAllMethods map[types.Type]bool
}

func CreateSpinner(manifest *config.Manifest) *Spinner {
	sp := &Spinner{
		Module:              wire.Module{},
		manifest:            manifest,
		packages:            make(map[*types.Package]*Package),
		typeTable:           make(map[string]*wire.Type),
		funcsDone:           make(map[types.Object]*wire.Function),
		funcsTodo:           make(map[*types.Func]bool),
		typesNeedAllMethods: make(map[types.Type]bool),
	}

	return sp
}

func (sp *Spinner) CreatePackage(pkg *types.Package, files []*ast.File, info *types.Info) *Package {
	return sp.createPackage(pkg, files, info)
}

func (sp *Spinner) Build() {
	sp.build()
}
