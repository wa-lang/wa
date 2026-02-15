// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appobjdump

import (
	"fmt"

	"wa-lang.org/wa/internal/native/link/pe"
)

func cmdProgdumpPE(filename string) {
	f, err := pe.OpenPE64(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Printf("Class      : PE32+ (64-bit)\n")
	fmt.Printf("Entry      : 0x%X\n", f.OptionalHeader.AddressOfEntryPoint)
	fmt.Printf("Image Base : 0x%X\n", f.OptionalHeader.ImageBase)
	fmt.Printf("Section Num: %d\n", f.FileHeader.NumberOfSections)
	fmt.Println()

	for _, sec := range f.Sections {
		data, err := sec.Data()
		if err != nil {
			fmt.Println("ERR:", err)
			continue
		}

		if sec.Characteristics&pe.IMAGE_SCN_MEM_EXECUTE != 0 {
			textAddr := f.OptionalHeader.ImageBase + uint64(sec.VirtualAddress)
			textData := data
			printProgText_x64(textAddr, textData)
		} else {
			dataAddr := f.OptionalHeader.ImageBase + uint64(sec.VirtualAddress)
			dataData := data
			printProgData(dataAddr, dataData)
		}
	}
}
