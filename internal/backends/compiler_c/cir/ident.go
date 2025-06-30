// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package cir

/**************************************
Ident:
**************************************/
type Ident struct {
	Name string
}

func NewIdent(s string) *Ident {
	return &Ident{Name: s}
}

func (i *Ident) CIRString() string {
	return i.Name
}

func (i *Ident) IsExpr() {}
