// 版权 @2021 凹语言 作者。保留所有权利。

package astutil

import (
	"strings"

	"wa-lang.org/wa/internal/ast"
)

// 注释中的指令
type CommentInfo struct {
	BuildIgnore bool     // #wa:build ignore
	BuildTags   []string // #wa:build s1 s2 ...

	NeesConstructor bool       // #wa:need-constructor
	Inline          bool       // #wa:inline
	Nobounds        bool       // #wa:nobounds
	LinkName        string     // #wa:linkname xxx
	ExportName      string     // #wa:export xxx
	ImportName      [2]string  // #wa:import xxx yyy
	ForceRegister   bool       // #wa:force_register
	RuntimeGetter   bool       // #wa:runtime_getter
	RuntimeSetter   bool       // #wa:runtime_setter
	RuntimeSizer    bool       // #wa:runtime_sizer
	WasmModule      string     // #wa:wasm-module xxx
	Generic         []string   // #wa:generic xxx yyy
	Operator        [][]string // #wa:operator + xxx yyy
	Embed           string     // #wa:embed filename
}

// 获取节点关联文档
func NodeDoc(n ast.Node) *ast.CommentGroup {
	switch n := n.(type) {
	case *ast.File:
		return n.Doc
	case *ast.ImportSpec:
		return n.Doc
	case *ast.GenDecl:
		return n.Doc
	case *ast.FuncDecl:
		return n.Doc
	case *ast.TypeSpec:
		return n.Doc
	case *ast.ValueSpec:
		return n.Doc
	case *ast.Field:
		return n.Doc
	}
	return nil
}

// 解析注释信息
func ParseCommentInfo(docList ...*ast.CommentGroup) (info CommentInfo) {
	for _, doc := range docList {
		if doc == nil {
			return
		}
		for _, comment := range doc.List {
			if !strContainPrefix(comment.Text, "#wa:", "//wa:") {
				continue
			}
			parts := strings.Fields(comment.Text)
			switch parts[0] {
			case "#wa:build", "//wa:build":
				if len(parts) >= 2 {
					info.BuildIgnore = parts[1] == "ignore"
				}
				info.BuildTags = parts[1:]
			case "#wa:need-constructor", "//wa:need-constructor":
				info.NeesConstructor = true
			case "#wa:inline", "//wa:inline":
				info.Inline = true
			case "#wa:nobounds", "//wa:nobounds":
				info.Nobounds = true
			case "#wa:linkname", "//wa:linkname":
				if len(parts) >= 2 {
					info.LinkName = strings.Join(parts[1:], " ")
				}
			case "#wa:export", "//wa:export":
				if len(parts) >= 2 {
					info.ExportName = parts[1]
				}
			case "#wa:import", "//wa:import":
				if len(parts) >= 3 {
					info.ImportName[0] = parts[1]
					info.ImportName[1] = parts[2]
				}
			case "#wa:force_register", "//wa:force_register":
				info.ForceRegister = true
			case "#wa:runtime_getter", "//wa:runtime_getter":
				info.RuntimeGetter = true
			case "#wa:runtime_setter", "//wa:runtime_setter":
				info.RuntimeSetter = true
			case "#wa:runtime_sizer":
				info.RuntimeSizer = true

			case "#wa:wasm-module", "//wa:wasm-module":
				if len(parts) >= 2 {
					info.WasmModule = parts[1]
				}

			case "#wa:generic":
				info.Generic = append(info.Generic, parts[1:]...)
			case "#wa:operator":
				info.Operator = append(info.Operator, parts[1:])

			case "#wa:embed":
				if len(parts) >= 2 {
					info.Embed = parts[1]
				}
			}
		}
	}
	return
}

func strContainPrefix(s string, prefixList ...string) bool {
	if len(prefixList) == 0 {
		return true
	}
	for _, prefix := range prefixList {
		if !strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
