package wat

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
	return indent + "(import \"" + o.moduleName + "\" \"" + o.objName + "\" (func $" + o.funcName + o.sig.String() + "))"
}
