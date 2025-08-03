// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package obj

type Plist struct {
	Name    *LSym
	Firstpc *Prog
	Recur   int
	Link    *Plist
}

// start a new Prog list.
func (ctxt *Link) NewPlist() *Plist {
	pl := new(Plist)
	if ctxt.Plist == nil {
		ctxt.Plist = pl
	} else {
		ctxt.Plast.Link = pl
	}
	ctxt.Plast = pl
	return pl
}
