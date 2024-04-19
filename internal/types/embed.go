// 版权 @2024 凹语言 作者。保留所有权利。

package types

import (
	"strconv"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/ast/astutil"
	"wa-lang.org/wa/internal/constant"
	"wa-lang.org/wa/internal/token"
)

// 预处理全局变量文件嵌入
func (check *Checker) processGlobalEmbed() {
	for obj, d := range check.objMap {
		if constObj, ok := obj.(*Const); ok {
			valueSpec := constObj.Node().(*ast.ValueSpec)
			commentInfo := astutil.ParseCommentInfo(constObj.NodeDoc())

			if commentInfo.Embed == "" {
				continue
			}
			if len(valueSpec.Names) != 1 {
				check.errorf(constObj.Pos(), "wa:embed cannot apply to multiple const")
				return
			}
			if len(valueSpec.Values) != 0 {
				check.errorf(constObj.Pos(), "wa:embed cannot apply to const with initializer")
				return
			}

			typeIdent, _ := valueSpec.Type.(*ast.Ident)
			if typeIdent == nil {
				check.errorf(constObj.Pos(), "wa:embed cannot apply to untyped const")
				return
			}

			scope, obj := check.pkg.scope.LookupParent(typeIdent.Name, valueSpec.Pos())
			if scope != Universe || obj.Type() != universeString {
				check.errorf(constObj.Pos(), "wa:embed invalid global type %v", obj)
				continue
			}

			// read file data
			var embedFound bool
			var embedData string
		Loop:
			for _, f := range check.files {
				for k, v := range f.EmbedMap {
					if k == commentInfo.Embed {
						embedFound = true
						embedData = v
						break Loop
					}
				}
			}
			if !embedFound {
				check.errorf(constObj.Pos(), "wa:embed %s: not found no matching files found", commentInfo.Embed)
				continue
			}

			valueSpec.Values = []ast.Expr{
				&ast.BasicLit{
					ValuePos: valueSpec.Pos(),
					Kind:     token.STRING,
					Value:    strconv.Quote(embedData),
				},
			}

			constObj.val = constant.MakeString(strconv.Quote(embedData))
			d.init = valueSpec.Values[0]
		}
	}
}
