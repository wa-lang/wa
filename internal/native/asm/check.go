// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package asm

import (
	"fmt"
	"strings"

	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/token"
)

// 检查语义
func CheckFile(fset *token.FileSet, file *ast.File) error {
	p := &_Checker{fset: fset, file: file}
	return p.checkFile()
}

type _Checker struct {
	fset *token.FileSet
	file *ast.File

	// 全局的命名对象
	// *ast.Const, *ast.Global, *ast.Func
	namedObjectMap map[string]interface{}
}

func (p *_Checker) checkFile() error {
	p.namedObjectMap = make(map[string]interface{})

	for _, x := range p.file.Consts {
		if err := p.checkConst(x); err != nil {
			return err
		}
	}
	for _, x := range p.file.Globals {
		if err := p.checkGlobal(x); err != nil {
			return err
		}
	}
	for _, x := range p.file.Funcs {
		if err := p.checkFunc(x); err != nil {
			return err
		}
	}
	return nil
}

func (p *_Checker) namedObjectPos(s string) token.Position {
	switch x := p.namedObjectMap[s].(type) {
	case *ast.Const:
		return p.fset.Position(x.Pos)
	case *ast.Global:
		return p.fset.Position(x.Pos)
	case *ast.Func:
		return p.fset.Position(x.Pos)
	}
	return token.Position{}
}

func (p *_Checker) checkConst(x *ast.Const) error {
	if _, exists := p.namedObjectMap[x.Name]; exists {
		return fmt.Errorf("%v: const %q redecared at %v", p.fset.Position(x.Pos), x.Name, p.namedObjectPos(x.Name))
	}

	// 常量必须是全局的名字, 以  '$' 开头
	if !strings.HasPrefix(x.Name, "$") {
		return fmt.Errorf("%v: const %q must start with '$'", p.fset.Position(x.Pos), x.Name)
	}

	// 验证常量的值
	if err := p.checkConst_value(x); err != nil {
		return err
	}

	// 记录当前命名对象
	p.namedObjectMap[x.Name] = x
	return nil
}

func (p *_Checker) checkConst_value(x *ast.Const) error {
	return nil
}

func (p *_Checker) checkGlobal(g *ast.Global) error {
	return nil
}

func (p *_Checker) checkFunc(fn *ast.Func) error {
	return nil
}
