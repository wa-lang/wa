// 版权 @2024 凹语言 作者。保留所有权利。

package ast

import (
	"bytes"
	"fmt"
	"io"
)

func Fprint(w io.Writer, m *Module) error {
	p := newPrinter(w)
	return p.Print(m)
}

func (m *Module) String() string {
	var buf bytes.Buffer
	p := newPrinter(&buf)
	p.Print(m)
	return buf.String()
}

type printer struct {
	w io.Writer
	m *Module
}

func newPrinter(w io.Writer) *printer {
	return &printer{w: w}
}

func (p *printer) Print(m *Module) error {
	p.m = m

	fmt.Fprint(p.w, "(module")
	if p.m.Name != "" {
		fmt.Fprint(p.w, "", p.m.Name)
	}

	if p.isModuleEmpty() {
		fmt.Fprint(p.w, ")")
		return nil
	}

	p.w.Write([]byte("TODO"))
	return fmt.Errorf("TODO")
}

func (p *printer) isModuleEmpty() bool {
	if len(p.m.Imports) > 0 {
		return false
	}
	if !p.isMemoyEmpty() {
		return false
	}
	if !p.isTableEmpty() {
		return false
	}
	if !p.isGlobalEmpty() {
		return false
	}
	if !p.isFuncEmpty() {
		return false
	}
	return true
}

func (p *printer) isMemoyEmpty() bool {
	return true
}

func (p *printer) isTableEmpty() bool {
	return true
}

func (p *printer) isGlobalEmpty() bool {
	return true
}

func (p *printer) isFuncEmpty() bool {
	return true
}
