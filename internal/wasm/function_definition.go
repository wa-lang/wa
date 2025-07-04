package wasm

// BuildFunctionDefinitions generates function metadata that can be parsed from
// the module. This must be called after all validation.
//
// Note: This is exported for tests who don't use wazero.Runtime or
// NewHostModule to compile the module.
func (m *Module) BuildFunctionDefinitions() {
	if len(m.FunctionSection) == 0 {
		return
	}

	var moduleName string
	var functionNames NameMap
	var localNames IndirectNameMap
	if m.NameSection != nil {
		moduleName = m.NameSection.ModuleName
		functionNames = m.NameSection.FunctionNames
		localNames = m.NameSection.LocalNames
	}

	importCount := m.ImportFuncCount()
	m.FunctionDefinitionSection = make([]*FunctionDefinition, 0, importCount+uint32(len(m.FunctionSection)))

	importFuncIdx := Index(0)
	for _, i := range m.ImportSection {
		if i.Type != ExternTypeFunc {
			continue
		}

		m.FunctionDefinitionSection = append(m.FunctionDefinitionSection, &FunctionDefinition{
			importDesc: &[2]string{i.Module, i.Name},
			index:      importFuncIdx,
			funcType:   m.TypeSection[i.DescFunc],
		})
		importFuncIdx++
	}

	for codeIndex, typeIndex := range m.FunctionSection {
		m.FunctionDefinitionSection = append(m.FunctionDefinitionSection, &FunctionDefinition{
			index:    Index(codeIndex) + importCount,
			funcType: m.TypeSection[typeIndex],
		})
	}

	n, nLen := 0, len(functionNames)
	for _, d := range m.FunctionDefinitionSection {
		// The function name section begins with imports, but can be sparse.
		// This keeps track of how far in the name section we've searched.
		funcIdx := d.index
		var funcName string
		for ; n < nLen; n++ {
			next := functionNames[n]
			if next.Index > funcIdx {
				break // we have function names, but starting at a later index.
			} else if next.Index == funcIdx {
				funcName = next.Name
				break
			}
		}

		d.moduleName = moduleName
		d.name = funcName
		d.debugName = wasmdebug_FuncName(moduleName, funcName, funcIdx)
		d.paramNames = paramNames(localNames, funcIdx, len(d.funcType.Params))

		for _, e := range m.ExportSection {
			if e.Type == ExternTypeFunc && e.Index == funcIdx {
				d.exportNames = append(d.exportNames, e.Name)
			}
		}
	}
}

// FunctionDefinition implements api.FunctionDefinition
type FunctionDefinition struct {
	moduleName  string
	index       Index
	name        string
	debugName   string
	funcType    *FunctionType
	importDesc  *[2]string
	exportNames []string
	paramNames  []string
}

// ModuleName implements the same method as documented on api.FunctionDefinition.
func (f *FunctionDefinition) ModuleName() string {
	return f.moduleName
}

// Index implements the same method as documented on api.FunctionDefinition.
func (f *FunctionDefinition) Index() uint32 {
	return f.index
}

// Name implements the same method as documented on api.FunctionDefinition.
func (f *FunctionDefinition) Name() string {
	return f.name
}

// DebugName implements the same method as documented on api.FunctionDefinition.
func (f *FunctionDefinition) DebugName() string {
	return f.debugName
}

// Import implements the same method as documented on api.FunctionDefinition.
func (f *FunctionDefinition) Import() (moduleName, name string, isImport bool) {
	if importDesc := f.importDesc; importDesc != nil {
		moduleName, name, isImport = importDesc[0], importDesc[1], true
	}
	return
}

// ExportNames implements the same method as documented on api.FunctionDefinition.
func (f *FunctionDefinition) ExportNames() []string {
	return f.exportNames
}

// ParamNames implements the same method as documented on api.FunctionDefinition.
func (f *FunctionDefinition) ParamNames() []string {
	return f.paramNames
}

// ParamTypes implements api.FunctionDefinition ParamTypes.
func (f *FunctionDefinition) ParamTypes() []ValueType {
	return f.funcType.Params
}

// ResultTypes implements api.FunctionDefinition ResultTypes.
func (f *FunctionDefinition) ResultTypes() []ValueType {
	return f.funcType.Results
}
