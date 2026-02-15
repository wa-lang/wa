// 版权 @2021 凹语言 作者。保留所有权利。

package astutil

import (
	"strings"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/token"
)

// 注释中的指令
type CommentInfo struct {
	BuildIgnore bool     // #wa:build ignore
	BuildTags   []string // #wa:build s1 s2 ...

	Align         string     // #wa:align xxx
	LinkName      string     // #wa:linkname xxx
	ExportName    string     // #wa:export xxx
	ImportName    [2]string  // #wa:import xxx yyy
	ForceRegister bool       // #wa:force_register
	RuntimeGetter bool       // #wa:runtime_getter
	RuntimeSetter bool       // #wa:runtime_setter
	RuntimeSizer  bool       // #wa:runtime_sizer
	Generic       []string   // #wa:generic xxx yyy
	Operator      [][]string // #wa:operator + xxx yyy
	Embed         string     // #wa:embed filename
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
			if !strContainPrefix(comment.Text, token.K_X_wa, token.K_X_wz) {
				continue
			}
			parts := strings.Fields(comment.Text)
			switch parts[0] {
			case token.K_X_wa_build, token.K_X_wz_build:
				if len(parts) >= 2 {
					switch parts[0] {
					case token.K_X_wa_build:
						if s := parts[1]; s == token.K_X_wa_build_arg_ignore {
							info.BuildIgnore = true
						}
					case token.K_X_wz_build:
						if s := parts[1]; s == token.K_X_wz_build_arg_ignore {
							info.BuildIgnore = true
						}
					default:
						panic("unreachable")
					}
				}
				info.BuildTags = parts[1:]

			case token.K_X_wa_align, token.K_X_wz_align:
				if len(parts) >= 2 {
					info.Align = parts[1]
				}

			case token.K_X_wa_linkname, token.K_X_wz_linkname:
				if len(parts) >= 2 {
					info.LinkName = strings.Join(parts[1:], " ")
				}
			case token.K_X_wa_export, token.K_X_wz_export:
				if len(parts) >= 2 {
					info.ExportName = parts[1]
				}
			case token.K_X_wa_import, token.K_X_wz_import:
				if len(parts) >= 3 {
					info.ImportName[0] = parts[1]
					info.ImportName[1] = parts[2]
				}
			case token.K_X_wa_force_register, token.K_X_wz_force_register:
				info.ForceRegister = true
			case token.K_X_wa_runtime_getter, token.K_X_wz_runtime_getter:
				info.RuntimeGetter = true
			case token.K_X_wa_runtime_setter, token.K_X_wz_runtime_setter:
				info.RuntimeSetter = true
			case token.K_X_wa_runtime_sizer, token.K_X_wz_runtime_sizer:
				info.RuntimeSizer = true

			case token.K_X_wa_generic, token.K_X_wz_generic:
				info.Generic = append(info.Generic, parts[1:]...)
			case token.K_X_wa_operator, token.K_X_wz_operator:
				info.Operator = append(info.Operator, parts[1:])

			case token.K_X_wa_embed, token.K_X_wz_embed:
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
