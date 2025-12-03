// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"sort"
	"strings"
)

var (
	flagPkgname = flag.String("pkg", "main", "set package name")
	flagOutput  = flag.String("out", "a_out2.go.txt", "set output file")
)

func main() {
	flag.Parse()

	if *flagPkgname == "" {
		fmt.Fprintf(os.Stderr, "usage: package name missing\n")
		flag.Usage()
		os.Exit(2)
	}
	if *flagOutput == "" {
		fmt.Fprintf(os.Stderr, "usage: output file missing\n")
		flag.Usage()
		os.Exit(2)
	}

	var out bytes.Buffer

	// 将 OpFormatType 表格化
	genOpFormatType_table(&out, *flagPkgname)

	outputCode := out.Bytes()

	// 修复参数信息的大小写
	outputCode = bytes.ReplaceAll(outputCode, []byte("arg_"), []byte("Arg_"))

	// 格式化输出代码
	if goodCode, err := format.Source(outputCode); err == nil {
		outputCode = goodCode
	}

	// 保存到文件
	if err := os.WriteFile(*flagOutput, outputCode, 0666); err != nil {
		panic(err)
	}
}

func genOpFormatType_table(w *bytes.Buffer, pkgname string) {
	const headCode = `// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 注意: 此代码是程序生成, 不要手动修改!!!

`

	w.Write([]byte(headCode))

	w.WriteString(fmt.Sprintf("package %s\n\n", pkgname))
	w.WriteString(`import "wa-lang.org/wa/internal/native/abi"` + "\n")

	// 根据 op 名字排序
	opKeyList := make([]string, 0, len(AOpContextTable))
	opIdxMap := make(map[string]int)
	for i, op := range AOpContextTable {
		opKeyList = append(opKeyList, op.name)
		opIdxMap[op.name] = i
	}
	sort.Slice(opKeyList, func(i, j int) bool {
		// 点转为下划线之后, 以保持最终的key顺序
		si := strings.ReplaceAll(opKeyList[i], ".", "_")
		sj := strings.ReplaceAll(opKeyList[j], ".", "_")
		return si < sj
	})

	// 1. 生成指令定义列表
	{
		w.WriteString("const (\n")
		w.WriteString("_ abi.As = iota\n")

		for _, key := range opKeyList {
			i := opIdxMap[key]
			op := AOpContextTable[i]
			opName := "A" + strings.ReplaceAll(op.name, ".", "_")
			opComment := AOpContextTable_comment[i]
			w.WriteString(opName)
			w.WriteString(opComment)
			w.WriteRune('\n')
		}

		w.WriteRune('\n')
		w.WriteString("ALAST\n")
		w.WriteString(")\n")
	}

	// 2. 生成指令名字列表
	w.WriteString(`var _Anames = [...]string{` + "\n")
	for _, key := range opKeyList {
		i := opIdxMap[key]
		op := AOpContextTable[i]
		opName := "A" + strings.ReplaceAll(op.name, ".", "_")
		w.WriteString(opName)
		w.WriteString(`:"`)
		w.WriteString(op.name)
		w.WriteString(`",`)
		w.WriteRune('\n')
	}
	w.WriteString("}\n")

	// 3. 生成指令编码元信息列表

	w.WriteString("// 指令编码信息表\n")
	w.WriteString(`var _AOpContextTable = [...]_OpContextType{` + "\n")
	for _, key := range opKeyList {
		i := opIdxMap[key]
		op := AOpContextTable[i]
		opName := "A" + strings.ReplaceAll(op.name, ".", "_")
		fmt.Fprintf(w, "\t%s: {mask: 0x%08x, value: 0x%08x, op: %s, fmt:%s},\n",
			opName, op.mask, op.value, opName,
			OpFormatTypeString(getOpFormatType(i, op.args)),
		)
	}
	w.WriteString("}\n")
}

func getOpFormatType(as int, asArgs instArgs) OpFormatType {
	regCount := 0
	regFCount := 0
	for i, argTyp := range asArgs {
		if argTyp == 0 {
			break
		}
		switch argTyp {
		case Arg_fd, Arg_fj, Arg_fk, Arg_fa:
			regFCount++
		case Arg_rd, Arg_rj, Arg_rk:
			regCount++
		case Arg_ui5_14_10:
			return OpFormatType_2R_ui5
		case Arg_ui6_15_10:
			return OpFormatType_2R_ui6
		case Arg_ui12_21_10:
			return OpFormatType_2R_ui12
		case Arg_si12_21_10:
			return OpFormatType_2R_si12
		case Arg_si14_23_10:
			return OpFormatType_2R_si14
		case Arg_si16_25_10:
			return OpFormatType_2R_si14
		case Arg_si20_24_5:
			return OpFormatType_1R_si20
		case Arg_offset_20_0:
			switch asArgs[0] {
			case Arg_rj:
				return OpFormatType_rj_offset
			case Arg_cj:
				return OpFormatType_cj_offset
			default:
				panic("unreachable")
			}
		case Arg_offset_25_0:
			return OpFormatType_offset
		case Arg_offset_15_0:
			assert(i == 2)
			switch {
			case asArgs[0] == Arg_rd:
				assert(asArgs[1] == Arg_rj)
				return OpFormatType_rd_rj_offset
			case asArgs[0] == Arg_rj:
				assert(asArgs[1] == Arg_rd)
				return OpFormatType_rj_rd_offset
			default:
				panic("unreachable")
			}
		case Arg_sa2_16_15:
			return OpFormatType_3R_s2
		case Arg_sa3_17_15:
			return OpFormatType_3R_s3
		case Arg_code_4_0:
			return OpFormatType_code_1R_si12
		case Arg_code_14_0:
			return OpFormatType_code

		case Arg_lsbw, Arg_msbw:
			return OpFormatType_msbw_lsbw
		case Arg_lsbd, Arg_msbd:
			return OpFormatType_msbd_lsbd

		case Arg_fcsr_4_0:
			if asArgs[0] == Arg_fcsr_4_0 {
				return OpFormatType_fcsr_1R
			} else {
				return OpFormatType_1R_fcsr
			}

		case Arg_fcsr_9_5:
			return OpFormatType_1R_fcsr

		case Arg_cd:
			if arg2 := asArgs[2]; arg2 == Arg_fk {
				return OpFormatType_cd_2R
			} else {
				return OpFormatType_cd_1R
			}

		case Arg_cj:
			if i == 0 {
				assert(asArgs[1] == Arg_offset_20_0)
				return OpFormatType_cj_offset
			} else {
				assert(asArgs[0] == Arg_fd || asArgs[0] == Arg_rd)
				return OpFormatType_1R_cj
			}

		case Arg_csr_23_10:
			if i == 1 {
				return OpFormatType_1R_csr
			} else {
				assert(i == 2)
				return OpFormatType_2R_csr
			}

		case Arg_level_14_0:
			return OpFormatType_level
		case Arg_level_17_10:
			return OpFormatType_2R_level

		case Arg_seq_17_10:
			return OpFormatType_0_1R_seq

		case Arg_ca:
			assert(asArgs[0] == Arg_fd)
			assert(asArgs[1] == Arg_fj)
			assert(asArgs[2] == Arg_fk)
			assert(i == 3)
			return OpFormatType_3R_ca

		case Arg_hint_4_0:
			assert(i == 0)
			if asArgs[2] == Arg_rk {
				return OpFormatType_hint_2R
			} else {
				assert(asArgs[2] == Arg_si12_21_10)
				return OpFormatType_hint_1R_si12
			}
		case Arg_hint_14_0:
			return OpFormatType_hint

		case Arg_op_4_0:
			return OpFormatType_op_2R

		default:
			panic(fmt.Sprintf("argTyp: %d", argTyp))
		}
	}

	switch {
	case regCount == 0 && regFCount == 0:
		return OpFormatType_NULL

	case regCount > 0 && regFCount == 0:
		switch regCount {
		case 2:
			if asArgs[0] != Arg_rd {
				return OpFormatType_0_2R
			}
			return OpFormatType_2R
		case 3:
			return OpFormatType_3R
		default:
			panic("unreachable")
		}
	case regCount == 0 && regFCount > 0:
		switch regFCount {
		case 2:
			return OpFormatType_2F
		case 3:
			return OpFormatType_3F
		case 4:
			return OpFormatType_4F
		default:
			panic("unreachable")
		}
	default:
		switch {
		case regFCount == 1 && regCount == 1:
			switch {
			case asArgs[0] == Arg_fd && asArgs[1] == Arg_rj:
				return OpFormatType_1F_1R
			case asArgs[0] == Arg_rd && asArgs[1] == Arg_fj:
				return OpFormatType_1R_1F
			default:
				panic("unreachable")
			}
		case regFCount == 1 && regCount == 2:
			return OpFormatType_1F_2R
		default:
			panic(fmt.Sprintf("AS: %s", AOpContextTable[as].name))
		}
	}
}

func OpFormatTypeString(x OpFormatType) string {
	switch x {
	default:
		panic(fmt.Sprintf("OpFormatType: %d", int(x)))
	case OpFormatType_NULL:
		return "OpFormatType_NULL"
	case OpFormatType_2R:
		return "OpFormatType_2R"
	case OpFormatType_2F:
		return "OpFormatType_2F"
	case OpFormatType_1F_1R:
		return "OpFormatType_1F_1R"
	case OpFormatType_1R_1F:
		return "OpFormatType_1R_1F"
	case OpFormatType_3R:
		return "OpFormatType_3R"
	case OpFormatType_3F:
		return "OpFormatType_3F"
	case OpFormatType_1F_2R:
		return "OpFormatType_1F_2R"
	case OpFormatType_4F:
		return "OpFormatType_4F"
	case OpFormatType_2R_ui5:
		return "OpFormatType_2R_ui5"
	case OpFormatType_2R_ui6:
		return "OpFormatType_2R_ui6"
	case OpFormatType_2R_si12:
		return "OpFormatType_2R_si12"
	case OpFormatType_2R_ui12:
		return "OpFormatType_2R_ui12"
	case OpFormatType_2R_si14:
		return "OpFormatType_2R_si14"
	case OpFormatType_2R_si16:
		return "OpFormatType_2R_si16"
	case OpFormatType_1R_si20:
		return "OpFormatType_1R_si20"
	case OpFormatType_0_2R:
		return "OpFormatType_0_2R"
	case OpFormatType_3R_s2:
		return "OpFormatType_3R_s2"
	case OpFormatType_3R_s3:
		return "OpFormatType_3R_s3"
	case OpFormatType_code:
		return "OpFormatType_code"
	case OpFormatType_code_1R_si12:
		return "OpFormatType_code_1R_si12"
	case OpFormatType_msbw_lsbw:
		return "OpFormatType_msbw_lsbw"
	case OpFormatType_msbd_lsbd:
		return "OpFormatType_msbd_lsbd"
	case OpFormatType_fcsr_1R:
		return "OpFormatType_fcsr_1R"
	case OpFormatType_1R_fcsr:
		return "OpFormatType_1R_fcsr"
	case OpFormatType_cd_1R:
		return "OpFormatType_cd_1R"
	case OpFormatType_cd_1F:
		return "OpFormatType_cd_1F"
	case OpFormatType_cd_2R:
		return "OpFormatType_cd_2R"
	case OpFormatType_cd_2F:
		return "OpFormatType_cd_2F"
	case OpFormatType_1R_cj:
		return "OpFormatType_1R_cj"
	case OpFormatType_1F_cj:
		return "OpFormatType_1F_cj"
	case OpFormatType_1R_csr:
		return "OpFormatType_1R_csr"
	case OpFormatType_2R_csr:
		return "OpFormatType_2R_csr"
	case OpFormatType_2R_level:
		return "OpFormatType_2R_level"
	case OpFormatType_level:
		return "OpFormatType_level"
	case OpFormatType_0_1R_seq:
		return "OpFormatType_0_1R_seq"
	case OpFormatType_op_2R:
		return "OpFormatType_op_2R"
	case OpFormatType_3R_ca:
		return "OpFormatType_3R_ca"
	case OpFormatType_hint_1R_si12:
		return "OpFormatType_hint_1R_si12"
	case OpFormatType_hint_2R:
		return "OpFormatType_hint_2R"
	case OpFormatType_hint:
		return "OpFormatType_hint"
	case OpFormatType_cj_offset:
		return "OpFormatType_cj_offset"
	case OpFormatType_rj_offset:
		return "OpFormatType_rj_offset"
	case OpFormatType_rj_rd_offset:
		return "OpFormatType_rj_rd_offset"
	case OpFormatType_rd_rj_offset:
		return "OpFormatType_rd_rj_offset"
	case OpFormatType_offset:
		return "OpFormatType_offset"
	}
}

func assert(ok bool) {
	if !ok {
		panic("assert failed")
	}
}
