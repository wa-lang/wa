// 版权 @2021 凹语言 作者。保留所有权利。

package astutil

import (
	"strings"

	"github.com/wa-lang/wa/internal/ast"
)

// 注释中的指令
type CommentInfo struct {
	BuildIgnore bool     // //wa:build ignore
	BuildTags   []string // //wa:build s1 s2 ...

	Inline   bool   // //wa:inline
	Nobounds bool   // //wa:nobounds
	LinkName string // //wa:linkname xxx
}

func ParseCommentInfo(docList ...*ast.CommentGroup) (info CommentInfo) {
	for _, doc := range docList {
		if doc == nil {
			return
		}
		for _, comment := range doc.List {
			if !strings.HasPrefix(comment.Text, "//wa:") {
				continue
			}
			parts := strings.Fields(comment.Text)
			switch parts[0] {
			case "//wa:build":
				if len(parts) >= 2 {
					info.BuildIgnore = parts[1] == "ignore"
				}
				info.BuildTags = parts[1:]
			case "//wa:inline":
				info.Inline = true
			case "//wa:nobounds":
				info.Nobounds = true
			case "//wa:linkname":
				if len(parts) >= 2 {
					info.LinkName = parts[1]
				}
			}
		}
	}
	return
}
