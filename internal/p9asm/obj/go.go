// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package obj

var (
	Framepointer_enabled int
	Fieldtrack_enabled   int
)

func Nopout(p *Prog) {
	p.As = ANOP
	p.Scond = 0
	p.From = Addr{}
	p.From3 = nil
	p.Reg = 0
	p.To = Addr{}
}
