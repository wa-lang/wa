// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package ast

import (
	"reflect"

	"wa-lang.org/wa/internal/native/token"
)

// 元素抽象
type Object interface {
	GetDoc() *CommentGroup
	BeginPos() token.Pos
	String() string
}

type ZhStringer interface {
	ZhString() string
}

func (p *File) GetDoc() *CommentGroup                   { return p.Doc }
func (p *Comment) GetDoc() *CommentGroup                { return nil }
func (p *CommentGroup) GetDoc() *CommentGroup           { return nil }
func (p *GasIntelSyntaxNoprefix) GetDoc() *CommentGroup { return nil }
func (p *Extern) GetDoc() *CommentGroup                 { return nil }
func (p *BasicLit) GetDoc() *CommentGroup               { return nil }
func (p *Global) GetDoc() *CommentGroup                 { return p.Doc }
func (p *InitValue) GetDoc() *CommentGroup              { return nil }
func (p *Func) GetDoc() *CommentGroup                   { return p.Doc }
func (p *FuncType) GetDoc() *CommentGroup               { return nil }
func (p *FuncBody) GetDoc() *CommentGroup               { return nil }
func (p *Local) GetDoc() *CommentGroup                  { return p.Doc }
func (p *Instruction) GetDoc() *CommentGroup            { return p.Doc }
func (p *BlankLine) GetDoc() *CommentGroup              { return nil }

func (p *File) BeginPos() token.Pos                   { return p.Pos }
func (p *Comment) BeginPos() token.Pos                { return p.Pos }
func (p *CommentGroup) BeginPos() token.Pos           { return p.List[0].Pos }
func (p *GasIntelSyntaxNoprefix) BeginPos() token.Pos { return p.Pos }
func (p *Extern) BeginPos() token.Pos                 { return p.Pos }
func (p *BasicLit) BeginPos() token.Pos               { return p.Pos }
func (p *Global) BeginPos() token.Pos                 { return p.Pos }
func (p *InitValue) BeginPos() token.Pos              { return p.Pos }
func (p *Func) BeginPos() token.Pos                   { return p.Pos }
func (p *FuncType) BeginPos() token.Pos               { return p.Pos }
func (p *FuncBody) BeginPos() token.Pos               { return p.Pos }
func (p *Local) BeginPos() token.Pos                  { return p.Pos }
func (p *Instruction) BeginPos() token.Pos            { return p.Pos }
func (p *BlankLine) BeginPos() token.Pos              { return p.Pos }

func isSameType(a, b Object) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}
