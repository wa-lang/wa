// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package regalloc

import (
	"container/list"
	"io"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/token"
)

type RegAlloctor struct {
	list  *list.List
	table map[string]*list.Element
}

type Handle struct {
	c       *RegAlloctor
	key     string
	value   abi.RegType
	deleter func(key string, value abi.RegType)
	refs    uint32
}

func (h *Handle) Key() string {
	return h.key
}

func (h *Handle) Value() abi.RegType {
	return h.value
}

func (h *Handle) Retain() (handle *Handle) {
	//h.c.addref(h)
	return h
}

func (h *Handle) Close() error {
	//h.c.unref(h)
	return nil
}

func New() *RegAlloctor {
	return &RegAlloctor{}
}

func (p *RegAlloctor) Use(typ token.Token) abi.RegType {
	return 0
}

func (p *RegAlloctor) Insert(key string, value abi.RegType, deleter func(key string, value abi.RegType)) (handle io.Closer) {
	panic("TODO")
}

func (p *RegAlloctor) Lookup(key string) (value abi.RegType, handle io.Closer, ok bool) {
	panic("TODO")
}

func (p *RegAlloctor) Erase(key string) {
	panic("TODO")
}

func (p *RegAlloctor) Close() error {
	panic("TODO")
}
