// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "strings"

/**************************************
ImpObj:
**************************************/
type ImpObj struct {
	moduleName string
	objName    string
}

func (o *ImpObj) ModuleName() string { return o.moduleName }
func (o *ImpObj) ObjName() string    { return o.objName }

/**************************************
ImpFunc:
**************************************/
type ImpFunc struct {
	ImpObj
	funcName string
	sig      FuncSig
}

func NewImpFunc(moduleName string, objName string, funcName string, sig FuncSig) *ImpFunc {
	return &ImpFunc{ImpObj: ImpObj{moduleName: moduleName, objName: objName}, funcName: funcName, sig: sig}
}
func (o *ImpFunc) Type() ObjType { return ObjTypeFunc }
func (o *ImpFunc) Format(indent string) string {
	var sb strings.Builder
	sb.WriteString(indent)
	sb.WriteString("(import \"")
	sb.WriteString(o.moduleName)
	sb.WriteString("\" \"")
	sb.WriteString(o.objName)
	sb.WriteString("\" (func $")
	sb.WriteString(o.funcName)
	sb.WriteString(o.sig.String())
	sb.WriteString("))")
	return sb.String()
}
