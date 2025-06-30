// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"wa-lang.org/wa/internal/ast"
)

func (p *Program) GetPkgPathList() []string {
	var ss []string
	for k := range p.Pkgs {
		ss = append(ss, k)
	}
	sort.Strings(ss)
	return ss
}

func (p *Program) DebugString() string {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, "Program.Cfg:", jsonString(p.Cfg))
	fmt.Fprintln(&buf, "Program.Manifest:", jsonString(p.Manifest))

	for _, k := range p.GetPkgPathList() {
		fmt.Fprintln(&buf, "PkgPath:", k)
		for _, f := range p.Pkgs[k].Files {
			fmt.Fprintln(&buf, "File:", f.Name.Name)
			ast.Fprint(&buf, p.Fset, f, ast.NotNilFilter)
		}
	}

	return buf.String()
}

func jsonString(x interface{}) string {
	d, _ := json.MarshalIndent(x, "", "    ")
	return string(d)
}
