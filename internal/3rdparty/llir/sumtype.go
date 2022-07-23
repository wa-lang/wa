package llir

// === [ constant.Constant ] ===================================================

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*Global) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*Func) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*Alias) IsConstant() {}

// IsConstant ensures that only constants can be assigned to the
// constant.Constant interface.
func (*IFunc) IsConstant() {}

// === [ ir.FuncAttribute ] ====================================================

// IsFuncAttribute ensures that only function attributes can be assigned to
// the ir.FuncAttribute interface.
func (AttrString) IsFuncAttribute() {}

// IsFuncAttribute ensures that only function attributes can be assigned to
// the ir.FuncAttribute interface.
func (AttrPair) IsFuncAttribute() {}

// IsFuncAttribute ensures that only function attributes can be assigned to the
// ir.FuncAttribute interface.
func (*AttrGroupDef) IsFuncAttribute() {}

// IsFuncAttribute ensures that only function attributes can be assigned to the
// ir.FuncAttribute interface.
func (Align) IsFuncAttribute() {}

// IsFuncAttribute ensures that only function attributes can be assigned to the
// ir.FuncAttribute interface.
func (AlignStack) IsFuncAttribute() {}

// IsFuncAttribute ensures that only function attributes can be assigned to the
// ir.FuncAttribute interface.
func (AllocSize) IsFuncAttribute() {}

// IsFuncAttribute ensures that only parameter attributes can be assigned to
// the ir.ParamAttribute interface.
func (Preallocated) IsFuncAttribute() {}

// === [ ir.Instruction ] ======================================================

// Binary instructions.
func (*InstFNeg) isInstruction() {}

// Binary instructions.
func (*InstAdd) isInstruction()  {}
func (*InstFAdd) isInstruction() {}
func (*InstSub) isInstruction()  {}
func (*InstFSub) isInstruction() {}
func (*InstMul) isInstruction()  {}
func (*InstFMul) isInstruction() {}
func (*InstUDiv) isInstruction() {}
func (*InstSDiv) isInstruction() {}
func (*InstFDiv) isInstruction() {}
func (*InstURem) isInstruction() {}
func (*InstSRem) isInstruction() {}
func (*InstFRem) isInstruction() {}

// Bitwise instructions.
func (*InstShl) isInstruction()  {}
func (*InstLShr) isInstruction() {}
func (*InstAShr) isInstruction() {}
func (*InstAnd) isInstruction()  {}
func (*InstOr) isInstruction()   {}
func (*InstXor) isInstruction()  {}

// Vector instructions.
func (*InstExtractElement) isInstruction() {}
func (*InstInsertElement) isInstruction()  {}
func (*InstShuffleVector) isInstruction()  {}

// Aggregate instructions.
func (*InstExtractValue) isInstruction() {}
func (*InstInsertValue) isInstruction()  {}

// Memory instructions.
func (*InstAlloca) isInstruction()        {}
func (*InstLoad) isInstruction()          {}
func (*InstStore) isInstruction()         {}
func (*InstFence) isInstruction()         {}
func (*InstCmpXchg) isInstruction()       {}
func (*InstAtomicRMW) isInstruction()     {}
func (*InstGetElementPtr) isInstruction() {}

// Conversion instructions.
func (*InstTrunc) isInstruction()         {}
func (*InstZExt) isInstruction()          {}
func (*InstSExt) isInstruction()          {}
func (*InstFPTrunc) isInstruction()       {}
func (*InstFPExt) isInstruction()         {}
func (*InstFPToUI) isInstruction()        {}
func (*InstFPToSI) isInstruction()        {}
func (*InstUIToFP) isInstruction()        {}
func (*InstSIToFP) isInstruction()        {}
func (*InstPtrToInt) isInstruction()      {}
func (*InstIntToPtr) isInstruction()      {}
func (*InstBitCast) isInstruction()       {}
func (*InstAddrSpaceCast) isInstruction() {}

// Other instructions.
func (*InstICmp) isInstruction()       {}
func (*InstFCmp) isInstruction()       {}
func (*InstPhi) isInstruction()        {}
func (*InstSelect) isInstruction()     {}
func (*InstFreeze) isInstruction()     {}
func (*InstCall) isInstruction()       {}
func (*InstVAArg) isInstruction()      {}
func (*InstLandingPad) isInstruction() {}
func (*InstCatchPad) isInstruction()   {}
func (*InstCleanupPad) isInstruction() {}

// === [ ir.ParamAttribute ] ===================================================

// IsParamAttribute ensures that only parameter attributes can be assigned to
// the ir.ParamAttribute interface.
func (AttrString) IsParamAttribute() {}

// IsParamAttribute ensures that only parameter attributes can be assigned to
// the ir.ParamAttribute interface.
func (AttrPair) IsParamAttribute() {}

// IsParamAttribute ensures that only parameter attributes can be assigned to
// the ir.ParamAttribute interface.
func (Align) IsParamAttribute() {}

// IsParamAttribute ensures that only parameter attributes can be assigned to
// the ir.ParamAttribute interface.
func (Byval) IsParamAttribute() {}

// IsParamAttribute ensures that only parameter attributes can be assigned to
// the ir.ParamAttribute interface.
func (Dereferenceable) IsParamAttribute() {}

// IsParamAttribute ensures that only parameter attributes can be assigned to
// the ir.ParamAttribute interface.
func (Preallocated) IsParamAttribute() {}

// IsParamAttribute ensures that only parameter attributes can be assigned to
// the ir.ParamAttribute interface.
func (SRet) IsParamAttribute() {}

// === [ ir.ReturnAttribute ] ==================================================

// IsReturnAttribute ensures that only return attributes can be assigned to
// the ir.ReturnAttribute interface.
func (AttrString) IsReturnAttribute() {}

// IsReturnAttribute ensures that only return attributes can be assigned to
// the ir.ReturnAttribute interface.
func (AttrPair) IsReturnAttribute() {}

// IsReturnAttribute ensures that only return attributes can be assigned to the
// ir.ReturnAttribute interface.
func (Align) IsReturnAttribute() {}

// IsReturnAttribute ensures that only return attributes can be assigned to
// the ir.ReturnAttribute interface.
func (Dereferenceable) IsReturnAttribute() {}
