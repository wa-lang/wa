// 版权 @2024 凹语言 作者。保留所有权利。

package types

import (
	"strconv"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/ast/astutil"
	"wa-lang.org/wa/internal/token"
)

// 预处理全局变量文件嵌入
func (check *Checker) processGlobalEmbed() {
	for obj, d := range check.objMap {
		if varObj, ok := obj.(*Var); ok {
			varSpec := varObj.Node().(*ast.ValueSpec)
			varCommentInfo := astutil.ParseCommentInfo(varObj.NodeDoc())

			if varCommentInfo.Embed == "" {
				continue
			}
			if len(varSpec.Names) != 1 {
				check.errorf(varObj.Pos(), "wa:embed only support one global")
				return
			}
			if len(varSpec.Values) != 0 || varSpec.Type == nil {
				check.errorf(varObj.Pos(), "wa:embed donot support init values")
				return
			}

			typeIdent, _ := varSpec.Type.(*ast.Ident)
			if typeIdent == nil {
				check.errorf(varObj.Pos(), "wa:embed must have type")
				return
			}

			scope, obj := check.pkg.scope.LookupParent(typeIdent.Name, varSpec.Pos())
			if scope != Universe || obj.Type() != universeString {
				check.errorf(varObj.Pos(), "wa:embed invalid global type %v", obj)
				continue
			}

			// read file data
			var embedFound bool
			var embedData string
		Loop:
			for _, f := range check.files {
				for k, v := range f.EmbedMap {
					if k == varCommentInfo.Embed {
						embedFound = true
						embedData = v
						break Loop
					}
				}
			}
			if !embedFound {
				check.errorf(varObj.Pos(), "wa:embed file not found")
				continue
			}

			varSpec.Values = []ast.Expr{
				&ast.BasicLit{
					ValuePos: varSpec.Pos(),
					Kind:     token.STRING,
					Value:    strconv.Quote(embedData),
				},
			}

			d.init = varSpec.Values[0]
		}
	}
}
