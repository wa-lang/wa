package llenum

// === [ metadata.DIExpressionField ] ==========================================

// IsDIExpressionField ensures that only DIExpression fields can be assigned to
// the metadata.DIExpressionField interface.
func (DwarfAttEncoding) IsDIExpressionField() {}

// IsDIExpressionField ensures that only DIExpression fields can be assigned to
// the metadata.DIExpressionField interface.
func (DwarfOp) IsDIExpressionField() {}

// === [ ir.FuncAttribute ] ====================================================

// IsFuncAttribute ensures that only function attributes can be assigned to the
// ir.FuncAttribute interface.
func (FuncAttr) IsFuncAttribute() {}

// === [ ir.ParamAttribute ] ===================================================

// IsParamAttribute ensures that only parameter attributes can be assigned to
// the ir.ParamAttribute interface.
func (ParamAttr) IsParamAttribute() {}

// === [ ir.ReturnAttribute ] ==================================================

// IsReturnAttribute ensures that only return attributes can be assigned to the
// ir.ReturnAttribute interface.
func (ReturnAttr) IsReturnAttribute() {}
