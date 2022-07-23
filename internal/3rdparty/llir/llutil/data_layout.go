package llutil

import (
	"fmt"
	"strconv"
	"strings"
)

// DataLayout is a structural representation of the datalayout
// string from https://llvm.org/docs/LangRef.html#data-layout
type DataLayout struct {
	IsBigEndian                bool                                   // default - True
	NaturalStackAlignment      uint64                                 // bits in multiple of 8, default: 0
	ProgramMemoryAddressSpace  uint64                                 // default - 0
	GlobalVarAddressSpace      uint64                                 // default - 0
	AllocaAddressSpace         uint64                                 // default - 0
	PointerSizeAlignment       map[uint64]*PointerSizeAlignment       // Key is addressspace
	IntegerSizeAlignment       map[uint64]*IntegerSizeAlignment       // Key is size
	VectorSizeAlignment        map[uint64]*VectorSizeAlignment        // Key is size
	FloatingPointSizeAlignment map[uint64]*FloatingPointSizeAlignment // Key is size
	AggregateAlignment         *AggregateAlignment                    // default - 0:64
	FunctionPointerAlignment   *FunctionPointerAlignment              // default - Not specified in LLVM Ref doc
	Mangling                   ManglingStyle                          // depends on os - "e" for linux, "o" for mac
	NativeIntBitWidths         []uint64                               // default - depends on arch "8:16:32:64" for x86-64
	NonIntegralPointerTypes    []string                               // Unsure of the datatype, so leaving it as string to accommodate anything.
}

// NewDataLayout returns an instance of DataLayout with default values
// taken from https://llvm.org/docs/LangRef.html#data-layout
func NewDataLayout(os string, arch string) *DataLayout {
	dl := &DataLayout{
		IsBigEndian:                true,
		NaturalStackAlignment:      0,
		ProgramMemoryAddressSpace:  0,
		GlobalVarAddressSpace:      0,
		AllocaAddressSpace:         0,
		PointerSizeAlignment:       getDefaultPtrAligns(),
		IntegerSizeAlignment:       getDefaultIntAligns(),
		VectorSizeAlignment:        getDefaultVectorAligns(),
		FloatingPointSizeAlignment: getDefaultFloatAligns(),
		AggregateAlignment:         NewAggregateAlignment(0, 64),
	}
	switch arch {
	case "x86-64":
		dl.NativeIntBitWidths = []uint64{8, 16, 32, 64}
	}

	if os == "linux" {
		dl.Mangling = ELF
	} else if strings.HasPrefix(os, "darwin") || strings.HasPrefix(os, "macosx") {
		dl.Mangling = MachO
	}
	return dl
}

// NewDataLayoutFromString parses a datalayout string and constructs a DataLayout object.
// Default values whereever applicable will be added if options are ommitted.
//
// This API expects the string to be valid as per LLVM spec, so there aren't many validations.
func NewDataLayoutFromString(layoutString, os, arch string) (*DataLayout, error) {
	dl := NewDataLayout(os, arch)
	var err error
	specs := strings.Split(layoutString, "-")
	for _, spec := range specs {
		if spec == "e" {
			dl.IsBigEndian = false
		} else if spec == "E" {
			dl.IsBigEndian = true
		} else if strings.HasPrefix(spec, "S") {
			dl.NaturalStackAlignment, err = strconv.ParseUint(strings.TrimPrefix(spec, "S"), 10, 64)
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(spec, "P") {
			dl.ProgramMemoryAddressSpace, err = strconv.ParseUint(strings.TrimPrefix(spec, "P"), 10, 64)
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(spec, "G") {
			dl.GlobalVarAddressSpace, err = strconv.ParseUint(strings.TrimPrefix(spec, "G"), 10, 64)
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(spec, "A") {
			dl.AllocaAddressSpace, err = strconv.ParseUint(strings.TrimPrefix(spec, "A"), 10, 64)
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(spec, "p") {
			addPtrSizeAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "i") {
			addIntSizeAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "v") {
			addVectorSizeAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "f") {
			addFloatSizeAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "a") {
			addAggAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "F") {
			addFuncPtrAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "m") {
			manglingVals := strings.Split(spec, ":")
			if len(manglingVals) != 2 {
				return nil, fmt.Errorf(`expecting 1 value separated by ":", but got %q`, spec)
			}
			dl.Mangling = ManglingStyle(manglingVals[1])
		} else if strings.HasPrefix(spec, "n") {
			addNativeIntBitWidths(spec, dl)
		} else if strings.HasPrefix(spec, "ni") {
			niValues := strings.Split(strings.TrimPrefix(spec, "ni"), ":")
			if len(niValues) == 0 {
				return nil, fmt.Errorf(`expecting atleast 1 value separated by ":", got %q`, spec)
			}
			dl.NonIntegralPointerTypes = niValues
		}
	}
	return dl, nil
}

func addPtrSizeAlignFromString(spec string, dl *DataLayout) error {
	ptrVals := strings.Split(spec, ":")
	ptrValsLen := len(ptrVals)
	if len(ptrVals) < 3 {
		return fmt.Errorf(`expecting atleast 3 values separated by ":", got %q`, spec)
	}
	ptrAddSpace := uint64(0)
	var err error
	if len(ptrVals[0]) > 1 {
		ptrAddSpace, err = strconv.ParseUint(strings.TrimPrefix(ptrVals[0], "p"), 10, 64)
		if err != nil {
			return err
		}
	}
	size, err := strconv.ParseUint(ptrVals[1], 10, 64)
	if err != nil {
		return err
	}
	abi, err := strconv.ParseUint(ptrVals[2], 10, 64)
	if err != nil {
		return err
	}
	pref := abi
	ind := size
	if ptrValsLen > 3 {
		pref, err = strconv.ParseUint(ptrVals[3], 10, 64)
		if err != nil {
			return err
		}
		if ptrValsLen == 5 {
			ind, err = strconv.ParseUint(ptrVals[4], 10, 64)
			if err != nil {
				return err
			}
		}
	}
	dl.PointerSizeAlignment[ptrAddSpace] = NewPointerSizeAlignment(ptrAddSpace, size, abi, pref, ind)
	return nil
}

func addIntSizeAlignFromString(spec string, dl *DataLayout) error {
	intVals := strings.Split(spec, ":")
	intValsLen := len(intVals)
	if intValsLen < 2 {
		return fmt.Errorf(`expecting atleast 2 values separated by ":", got %q`, spec)
	}
	size, err := strconv.ParseUint(strings.TrimPrefix(intVals[0], "i"), 10, 64)
	if err != nil {
		return err
	}
	abi, err := strconv.ParseUint(intVals[1], 10, 64)
	if err != nil {
		return err
	}
	pref := abi
	if intValsLen == 3 {
		pref, err = strconv.ParseUint(intVals[2], 10, 64)
		if err != nil {
			return err
		}
	}
	dl.IntegerSizeAlignment[size] = NewIntegerSizeAlignment(size, abi, pref)
	return nil
}

func addVectorSizeAlignFromString(spec string, dl *DataLayout) error {
	vectorVals := strings.Split(spec, ":")
	vectorValsLen := len(vectorVals)
	if vectorValsLen < 2 {
		return fmt.Errorf(`expecting atleast 2 values separated by ":", got %q`, spec)
	}
	size, err := strconv.ParseUint(strings.TrimPrefix(vectorVals[0], "v"), 10, 64)
	if err != nil {
		return err
	}
	abi, err := strconv.ParseUint(vectorVals[1], 10, 64)
	if err != nil {
		return err
	}
	pref := abi
	if vectorValsLen == 3 {
		pref, err = strconv.ParseUint(vectorVals[2], 10, 64)
		if err != nil {
			return err
		}
	}
	dl.VectorSizeAlignment[size] = NewVectorSizeAlignment(size, abi, pref)
	return nil
}

func addFloatSizeAlignFromString(spec string, dl *DataLayout) error {
	floatVals := strings.Split(spec, ":")
	floatValsLen := len(floatVals)
	if floatValsLen < 2 {
		return fmt.Errorf(`expecting atleast 2 values separated by ":", got %q`, spec)
	}
	size, err := strconv.ParseUint(strings.TrimPrefix(floatVals[0], "f"), 10, 64)
	if err != nil {
		return err
	}
	abi, err := strconv.ParseUint(floatVals[1], 10, 64)
	if err != nil {
		return err
	}
	pref := abi
	if floatValsLen == 3 {
		pref, err = strconv.ParseUint(floatVals[2], 10, 64)
		if err != nil {
			return err
		}
	}
	dl.FloatingPointSizeAlignment[size] = NewFloatingPointSizeAlignment(size, abi, pref)
	return nil
}

func addAggAlignFromString(spec string, dl *DataLayout) error {
	aggVals := strings.Split(spec, ":")
	aggValsLen := len(aggVals)
	if aggValsLen < 2 {
		return fmt.Errorf(`expecting atleast 2 values separated by ":", got %q`, spec)
	}
	abi, err := strconv.ParseUint(aggVals[1], 10, 64)
	if err != nil {
		return err
	}
	pref := abi
	if aggValsLen == 3 {
		pref, err = strconv.ParseUint(aggVals[2], 10, 64)
		if err != nil {
			return err
		}
	}
	dl.AggregateAlignment = NewAggregateAlignment(abi, pref)
	return nil
}

func addFuncPtrAlignFromString(spec string, dl *DataLayout) error {
	val := strings.TrimPrefix(spec, "F")
	var typ string
	if strings.HasPrefix(val, "i") {
		typ = "i"
	} else if strings.HasPrefix(val, "n") {
		typ = "n"
	} else {
		panic(fmt.Errorf(`invalid function pointer alignment type, expected "i" or "n" prefix, got %q`, spec))
	}
	abi, err := strconv.ParseUint(strings.TrimPrefix(val, typ), 10, 64)
	if err != nil {
		return err
	}
	dl.FunctionPointerAlignment = NewFunctionPointerAlignment(typ == "i", abi)
	return nil
}

func addNativeIntBitWidths(spec string, dl *DataLayout) error {
	// No length validation for widthContents here since "n32"
	// like values are valid and 32 is validated while parsing.
	widthContents := strings.Split(strings.TrimPrefix(spec, "n"), ":")
	widths := make([]uint64, len(widthContents))
	var err error
	for i, width := range widthContents {
		if widths[i], err = strconv.ParseUint(width, 10, 64); err != nil {
			return err
		}
	}
	dl.NativeIntBitWidths = widths
	return nil
}

func (dl *DataLayout) LLString() string {
	layout := &strings.Builder{}
	if dl.IsBigEndian {
		layout.WriteString("E")
	} else {
		layout.WriteString("e")
	}
	if dl.NaturalStackAlignment != 0 {
		layout.WriteString("-S")
		fmt.Fprintf(layout, "%d", dl.NaturalStackAlignment)
	}
	if dl.ProgramMemoryAddressSpace != 0 {
		layout.WriteString("-P")
		fmt.Fprintf(layout, "%d", dl.ProgramMemoryAddressSpace)
	}
	if dl.GlobalVarAddressSpace != 0 {
		layout.WriteString("-G")
		fmt.Fprintf(layout, "%d", dl.GlobalVarAddressSpace)
	}
	if dl.AllocaAddressSpace != 0 {
		layout.WriteString("-A")
		fmt.Fprintf(layout, "%d", dl.AllocaAddressSpace)
	}
	if len(dl.PointerSizeAlignment) > 0 {
		for k, v := range dl.PointerSizeAlignment {
			if k == 0 {
				layout.WriteString("-p")
			} else {
				layout.WriteString("-p")
				fmt.Fprintf(layout, "%d", k)
			}
			layout.WriteString(":")
			fmt.Fprintf(layout, "%d", v.Size)
			layout.WriteString(":")
			fmt.Fprintf(layout, "%d", v.ABIAlignment)
			layout.WriteString(":")
			fmt.Fprintf(layout, "%d", v.PreferredAlignment)
			layout.WriteString(":")
			fmt.Fprintf(layout, "%d", v.AddressCalculationIndex)
		}
	}
	if len(dl.IntegerSizeAlignment) > 0 {
		for k, v := range dl.IntegerSizeAlignment {
			layout.WriteString("-i")
			fmt.Fprintf(layout, "%d", k)
			layout.WriteString(":")
			fmt.Fprintf(layout, "%d", v.ABIAlignment)
			layout.WriteString(":")
			fmt.Fprintf(layout, "%d", v.PreferredAlignment)
		}
	}
	if len(dl.VectorSizeAlignment) > 0 {
		for k, v := range dl.VectorSizeAlignment {
			layout.WriteString("-v")
			fmt.Fprintf(layout, "%d", k)
			layout.WriteString(":")
			fmt.Fprintf(layout, "%d", v.ABIAlignment)
			layout.WriteString(":")
			fmt.Fprintf(layout, "%d", v.PreferredAlignment)
		}
	}
	if len(dl.FloatingPointSizeAlignment) > 0 {
		for k, v := range dl.FloatingPointSizeAlignment {
			layout.WriteString("-f")
			fmt.Fprintf(layout, "%d", k)
			layout.WriteString(":")
			fmt.Fprintf(layout, "%d", v.ABIAlignment)
			layout.WriteString(":")
			fmt.Fprintf(layout, "%d", v.PreferredAlignment)
		}
	}
	if dl.AggregateAlignment != nil {
		layout.WriteString("-a:")
		fmt.Fprintf(layout, "%d", dl.AggregateAlignment.ABIAlignment)
		layout.WriteString(":")
		fmt.Fprintf(layout, "%d", dl.AggregateAlignment.PreferredAlignment)
	}
	if dl.FunctionPointerAlignment != nil {
		if dl.FunctionPointerAlignment.IsIndependant {
			layout.WriteString("-Fi")
			fmt.Fprintf(layout, "%d", dl.FunctionPointerAlignment.ABIAlignment)
		} else {
			layout.WriteString("-Fn")
			fmt.Fprintf(layout, "%d", dl.FunctionPointerAlignment.ABIAlignment)
		}
	}
	layout.WriteString("-m:")
	layout.WriteString(string(dl.Mangling))
	if bitWidths := len(dl.NativeIntBitWidths); bitWidths > 0 {
		layout.WriteString("-n")
		fmt.Fprintf(layout, "%d", dl.NativeIntBitWidths[0])
		for i := 1; i < bitWidths; i++ {
			layout.WriteString(":")
			fmt.Fprintf(layout, "%d", dl.NativeIntBitWidths[i])
		}
	}
	if len(dl.NonIntegralPointerTypes) > 0 {
		layout.WriteString("-ni")
		for _, v := range dl.NonIntegralPointerTypes {
			layout.WriteString(":")
			layout.WriteString(v)
		}
	}
	return layout.String()
}

type PointerSizeAlignment struct { // All sizes are in bits.
	AddressSpace            uint64
	Size                    uint64
	ABIAlignment            uint64
	PreferredAlignment      uint64 // optional and defaults to abi
	AddressCalculationIndex uint64 // default 0
}

func NewPointerSizeAlignment(addSp, size, abiAl, prefAl, addCalInd uint64) *PointerSizeAlignment {
	return &PointerSizeAlignment{AddressSpace: addSp, Size: size, ABIAlignment: abiAl, PreferredAlignment: prefAl, AddressCalculationIndex: addCalInd}
}

func getDefaultPtrAligns() map[uint64]*PointerSizeAlignment {
	ptrAlignMap := make(map[uint64]*PointerSizeAlignment)
	ptrAlignMap[0] = NewPointerSizeAlignment(0, 64, 64, 64, 64)
	return ptrAlignMap
}

type IntegerSizeAlignment struct { // All sizes are in bits.
	Size               uint64
	ABIAlignment       uint64
	PreferredAlignment uint64 // optional and defaults to abi
}

func NewIntegerSizeAlignment(size, abiAl, prefAl uint64) *IntegerSizeAlignment {
	return &IntegerSizeAlignment{Size: size, ABIAlignment: abiAl, PreferredAlignment: prefAl}
}

func getDefaultIntAligns() map[uint64]*IntegerSizeAlignment {
	intAlignMap := make(map[uint64]*IntegerSizeAlignment)
	intAlignMap[1] = NewIntegerSizeAlignment(1, 8, 8)
	intAlignMap[8] = NewIntegerSizeAlignment(8, 8, 8)
	intAlignMap[16] = NewIntegerSizeAlignment(16, 16, 16)
	intAlignMap[32] = NewIntegerSizeAlignment(32, 32, 32)
	intAlignMap[64] = NewIntegerSizeAlignment(64, 32, 64)
	return intAlignMap
}

type VectorSizeAlignment struct { // All sizes are in bits.
	Size               uint64
	ABIAlignment       uint64
	PreferredAlignment uint64 // optional and defaults to abi
}

func NewVectorSizeAlignment(size, abiAl, prefAl uint64) *VectorSizeAlignment {
	return &VectorSizeAlignment{Size: size, ABIAlignment: abiAl, PreferredAlignment: prefAl}
}

func getDefaultVectorAligns() map[uint64]*VectorSizeAlignment {
	vectorAlignMap := make(map[uint64]*VectorSizeAlignment)
	vectorAlignMap[64] = NewVectorSizeAlignment(64, 64, 64)
	vectorAlignMap[128] = NewVectorSizeAlignment(128, 128, 128)
	return vectorAlignMap
}

type FloatingPointSizeAlignment struct { // All sizes are in bits.
	Size               uint64
	ABIAlignment       uint64
	PreferredAlignment uint64 // optional and defaults to abi
}

func NewFloatingPointSizeAlignment(size, abiAl, prefAl uint64) *FloatingPointSizeAlignment {
	return &FloatingPointSizeAlignment{Size: size, ABIAlignment: abiAl, PreferredAlignment: prefAl}
}

func getDefaultFloatAligns() map[uint64]*FloatingPointSizeAlignment {
	floatAlignMap := make(map[uint64]*FloatingPointSizeAlignment)
	floatAlignMap[16] = NewFloatingPointSizeAlignment(16, 16, 16)
	floatAlignMap[32] = NewFloatingPointSizeAlignment(32, 32, 32)
	floatAlignMap[64] = NewFloatingPointSizeAlignment(64, 64, 64)
	floatAlignMap[128] = NewFloatingPointSizeAlignment(128, 128, 128)
	return floatAlignMap
}

type AggregateAlignment struct {
	ABIAlignment       uint64
	PreferredAlignment uint64 // optional and defaults to abi
}

func NewAggregateAlignment(abiAl, prefAl uint64) *AggregateAlignment {
	return &AggregateAlignment{ABIAlignment: abiAl, PreferredAlignment: prefAl}
}

type FunctionPointerAlignment struct {
	IsIndependant bool
	ABIAlignment  uint64
}

func NewFunctionPointerAlignment(isInd bool, abiAl uint64) *FunctionPointerAlignment {
	return &FunctionPointerAlignment{IsIndependant: isInd, ABIAlignment: abiAl}
}

type ManglingStyle string

const (
	ELF        ManglingStyle = "e"
	Mips       ManglingStyle = "m"
	MachO      ManglingStyle = "o"
	WinX86COFF ManglingStyle = "x"
	WinCOFF    ManglingStyle = "w"
	XCOFF      ManglingStyle = "a"
)
