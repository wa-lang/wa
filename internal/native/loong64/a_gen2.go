// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 生成龙芯英文版指令集表格
// 1. 下载英文版到当前目录 https://loongson.github.io/LoongArch-Documentation/LoongArch-Vol1-EN.pdf
// 2. 先执行 `go run a_gen.go`, 生成 `a_out.go` 文件
// 3. 删除 `a_out2.go` 文件(生成的代码可能有错误), 重新生成
// 4. 再执行 `go run a_gen2.go`, 生成 `a_out2.go` 文件

//go:build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"strings"

	"wa-lang.org/wa/internal/native/abi"
	. "wa-lang.org/wa/internal/native/loong64"
)

var (
	flagOutput = flag.String("output", "a_out2.go", "set output file")
)

func main() {
	flag.Parse()

	if *flagOutput == "" {
		fmt.Fprintf(os.Stderr, "usage: output file missing\n")
		flag.Usage()
		os.Exit(2)
	}

	var out bytes.Buffer

	// 将 OpFormatType 表格化
	genOpFormatType_table(&out)

	outputCode := out.Bytes()

	// 修复参数信息的大小写
	outputCode = bytes.ReplaceAll(outputCode, []byte("arg_"), []byte("Arg_"))

	// 格式化输出代码
	if goodCode, err := format.Source(outputCode); err == nil {
		outputCode = goodCode
	} else {
		outputCode = append([]byte("//go:build ignore\n\n/*"), outputCode...)
		outputCode = append(outputCode, []byte("\n*/\n")...)
	}

	// 保存到文件
	if err := os.WriteFile(*flagOutput, outputCode, 0666); err != nil {
		panic(err)
	}
}

func genOpFormatType_table(w *bytes.Buffer) {
	const headCode = `// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 注意: 此代码是程序生成, 不要手动修改!!!

package loong64

import "wa-lang.org/wa/internal/native/abi"

// 指令的编码格式
func AsFormatType(as abi.As) OpFormatType {
	if as > 0 && int(as) < len(_AOpContextTable) {
		return _AOpFormatTypeTable[as]
	}
	return OpFormatType_NULL
}

// 指令编码类型表
var _AOpFormatTypeTable = [...]OpFormatType{
`
	const footCode = `}`

	w.Write([]byte(headCode))
	defer w.Write([]byte(footCode))

	for as := abi.As(0); as < ALAST; as++ {
		if !AsValid(as) {
			continue
		}

		asStr := "A" + AsString(as, "")
		asStr = strings.ReplaceAll(asStr, ".", "_")

		fmt.Fprintf(w, "\t%s: %s,\n", asStr, OpFormatTypeString(getOpFormatType(as)))
	}
}

func getOpFormatType(as abi.As) OpFormatType {
	asArgs := AsArgs(as)
	regCount := 0
	for i, argTyp := range asArgs {
		if argTyp == 0 {
			break
		}
		switch argTyp {
		case Arg_fd, Arg_fj, Arg_fk, Arg_fa:
			regCount++
		case Arg_rd, Arg_rj, Arg_rk:
			regCount++
		case Arg_ui5_14_10:
			return OpFormatType_2RI8
		case Arg_ui6_15_10:
			return OpFormatType_2RI8
		case Arg_ui12_21_10:
			return OpFormatType_2RI12
		case Arg_si12_21_10:
			return OpFormatType_2RI12
		case Arg_si14_23_10:
			return OpFormatType_2RI14
		case Arg_si16_25_10:
			return OpFormatType_1RI21
		case Arg_si20_24_5:
			return OpFormatType_1RI20
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
				return OpFormatType_hint_1R_si
			}
		case Arg_hint_14_0:
			return OpFormatType_hint

		case Arg_op_4_0:
			return OpFormatType_op_2R

		default:
			panic(fmt.Sprintf("argTyp: %d", argTyp))
		}
	}

	switch regCount {
	case 0:
		return OpFormatType_NULL
	case 2:
		if asArgs[0] != Arg_rd {
			// assert(opcode.op == AASRTLE_D || opcode.op == AASRTGT_D)
			return OpFormatType_0_2R
		}
		return OpFormatType_2R
	case 3:
		return OpFormatType_3R
	case 4:
		return OpFormatType_4R
	}

	return OpFormatType_NULL
}

func OpFormatTypeString(x OpFormatType) string {
	switch x {
	default:
		panic(fmt.Sprintf("OpFormatType: %d", int(x)))
	case OpFormatType_NULL:
		return "OpFormatType_NULL"
	case OpFormatType_2R:
		return "OpFormatType_2R"
	case OpFormatType_3R:
		return "OpFormatType_3R"
	case OpFormatType_4R:
		return "OpFormatType_4R"
	case OpFormatType_2RI8:
		return "OpFormatType_2RI8"
	case OpFormatType_2RI12:
		return "OpFormatType_2RI12"
	case OpFormatType_2RI14:
		return "OpFormatType_2RI14"
	case OpFormatType_2RI16:
		return "OpFormatType_2RI16"
	case OpFormatType_1RI20:
		return "OpFormatType_1RI20"
	case OpFormatType_1RI21:
		return "OpFormatType_1RI21"
	case OpFormatType_I26:
		return "OpFormatType_I26"
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
	case OpFormatType_cd_2R:
		return "OpFormatType_cd_2R"
	case OpFormatType_1R_cj:
		return "OpFormatType_1R_cj"
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
	case OpFormatType_hint_1R_si:
		return "OpFormatType_hint_1R_si"
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
