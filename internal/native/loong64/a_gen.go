// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 生成龙芯英文版指令集表格
// 1. 下载英文版到当前目录 https://loongson.github.io/LoongArch-Documentation/LoongArch-Vol1-EN.pdf
// 2. 执行 `go run a_gen.go`, 生成 `a_out.go` 文件

//go:build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/pdf"
)

var (
	flagPdf     = flag.String("pdf", "LoongArch-Vol1-EN.pdf", "set loong64 spec pdf")
	flagPackage = flag.String("pkg", "loong64", "set package name")
	flagOutput  = flag.String("output", "a_out.go", "set output file")
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("loong64spec: ")
}

func main() {
	flag.Parse()

	if *flagPdf == "" {
		fmt.Fprintf(os.Stderr, "usage: pdf missing\n")
		flag.Usage()
		os.Exit(2)
	}
	if *flagPackage == "" {
		fmt.Fprintf(os.Stderr, "usage: package name missing\n")
		flag.Usage()
		os.Exit(2)
	}
	if *flagOutput == "" {
		fmt.Fprintf(os.Stderr, "usage: output file missing\n")
		flag.Usage()
		os.Exit(2)
	}

	f, err := pdf.Open(*flagPdf)
	if err != nil {
		log.Fatal(err)
	}
	var prologue bytes.Buffer
	prologue.Write([]byte(`// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 注意: 此代码是程序生成, 不要手动修改!!!

package loong64

import "wa-lang.org/wa/internal/native/abi"

`))

	var op_f bytes.Buffer
	op_f.Write([]byte("const (\n\t_ abi.As = iota\n"))

	var opstr_f bytes.Buffer
	opstr_f.Write([]byte("var _Anames = [...]string{\n"))

	var instFormats_f bytes.Buffer
	instFormats_f.Write([]byte("// 指令编码信息表\n"))
	instFormats_f.Write([]byte("var _AOpContextTable = [...]_OpContextType{\n"))

	// Scan document looking for instructions.
	n := f.NumPage()
	var ops []string
	opstrs := map[string]string{}
	instFormatComments := map[string]string{}
	instFormats := map[string]string{}
	var fp int

	mergeMap := func(m1 map[string]string, m2 map[string]string) {
		for k := range m2 {
			m1[k] = m2[k]
		}
	}

	for pageNum := 1; pageNum <= n; pageNum++ {
		p := f.Page(pageNum)
		if fp == 0 {
			if !isFirstPage(p) {
				continue
			}
			fp = pageNum
		}
		cPageOps, cPageOpstrs, cPageInstFormatComments, cPageInstFormats := parsePage(pageNum, p, fp == pageNum)
		ops = append(ops, cPageOps...)
		mergeMap(opstrs, cPageOpstrs)
		mergeMap(instFormatComments, cPageInstFormatComments)
		mergeMap(instFormats, cPageInstFormats)
	}

	sort.Strings(ops)

	for _, op := range ops {
		// 1. 指令枚举值定义, 加 A 前缀, 比如 AADD_W
		op_f.Write([]byte(fmt.Sprintf("\t%s%s %s\n", "A", op, instFormatComments[op])))
		// 2. 指令枚举值对应的名字字符串
		opstr_f.Write([]byte(fmt.Sprintf("\t%s%s\n", "A", opstrs[op])))
		// 3. 指令的参数格式, 枚举值对应列表的下标
		instFormats_f.Write([]byte(fmt.Sprintf("\t%s%s:%s\n", "A", op, instFormats[op])))
	}

	// 增加一个结尾的指令标记
	op_f.Write([]byte("\n\tALAST\n"))

	op_f.Write([]byte(")\n\n"))
	opstr_f.Write([]byte("}\n\n"))
	instFormats_f.Write([]byte("}\n"))

	var fileTables bytes.Buffer
	fileTables.Write(prologue.Bytes())
	fileTables.Write(op_f.Bytes())
	fileTables.Write(opstr_f.Bytes())
	fileTables.Write(instFormats_f.Bytes())

	// 格式化输出代码
	outputCode := fileTables.Bytes()
	if goodCode, err := format.Source(outputCode); err == nil {
		outputCode = goodCode
	}

	// 保存到文件
	if err := os.WriteFile(*flagOutput, outputCode, 0666); err != nil {
		panic(err)
	}
}

func isFirstPage(page pdf.Page) bool {
	content := page.Content()
	appendixb := "AppendixB"
	ct := ""
	for _, t := range content.Text {
		ct += t.S
		if ct == "AppendixB" {
			return true
		}
		if strings.HasPrefix(appendixb, ct) {
			continue
		} else {
			return false
		}
	}
	return false
}

func getArg(name string) (length int, argName string) {
	switch {
	case strings.Contains("arg_fd", name):
		return 5, "arg_fd"
	case strings.Contains("arg_fj", name):
		return 5, "arg_fj"
	case strings.Contains("arg_fk", name):
		return 5, "arg_fk"
	case strings.Contains("arg_fa", name):
		return 5, "arg_fa"
	case strings.Contains("arg_rd", name):
		return 5, "arg_rd"
	case strings.Contains("arg_rj", name) || name == "rj!=0,1":
		return 5, "arg_rj"
	case strings.Contains("arg_rk", name):
		return 5, "arg_rk"
	case name == "csr":
		return 14, "arg_csr_23_10"
	case strings.Contains("arg_cd", name):
		return 5, "arg_cd"
	case strings.Contains("arg_cj", name):
		return 5, "arg_cj"
	case strings.Contains("arg_ca", name):
		return 5, "arg_ca"
	case strings.Contains(name, "sa"):
		length, _ := strconv.Atoi(strings.Split(name, "sa")[1])
		if length == 2 {
			argName = "arg_sa2_16_15"
		} else {
			argName = "arg_sa3_17_15"
		}
		return length, argName
	case strings.Contains("arg_seq_17_10", name):
		return 8, "arg_seq_17_10"
	case strings.Contains("arg_op_4_0", name):
		return 5, "arg_op_4_0"
	case strings.Contains(name, "ui"):
		length, _ := strconv.Atoi(strings.Split(name, "ui")[1])
		if length == 5 {
			argName = "arg_ui5_14_10"
		} else if length == 6 {
			argName = "arg_ui6_15_10"
		} else {
			argName = "arg_ui12_21_10"
		}
		return length, argName
	case strings.Contains("arg_lsbw", name):
		return 5, "arg_lsbw"
	case strings.Contains("arg_msbw", name):
		return 5, "arg_msbw"
	case strings.Contains("arg_lsbd", name):
		return 6, "arg_lsbd"
	case strings.Contains("arg_msbd", name):
		return 6, "arg_msbd"
	case strings.Contains(name, "si"):
		length, _ := strconv.Atoi(strings.Split(name, "si")[1])
		if length == 12 {
			argName = "arg_si12_21_10"
		} else if length == 14 {
			argName = "arg_si14_23_10"
		} else if length == 16 {
			argName = "arg_si16_25_10"
		} else {
			argName = "arg_si20_24_5"
		}
		return length, argName
	case strings.Contains(name, "offs"):
		splitName := strings.Split(name, ":")
		left, _ := strconv.Atoi(strings.Split(splitName[0], "[")[1])
		right, _ := strconv.Atoi(strings.Split(splitName[1], "]")[0])
		return left - right + 1, "offs"
	default:
		return 0, ""
	}
}

func binstrToHex(str string) string {
	rst := 0
	hex := "0x"
	charArray := []byte(str)
	for i := 0; i < 32; {
		rst = 1*(int(charArray[i+3])-48) + 2*(int(charArray[i+2])-48) + 4*(int(charArray[i+1])-48) + 8*(int(charArray[i])-48)
		switch rst {
		case 10:
			hex = hex + "a"
		case 11:
			hex = hex + "b"
		case 12:
			hex = hex + "c"
		case 13:
			hex = hex + "d"
		case 14:
			hex = hex + "e"
		case 15:
			hex = hex + "f"
		default:
			hex += strconv.Itoa(rst)
		}

		i = i + 4
	}
	return hex
}

/*
Here we deal with the instruction FCMP.cond.S/D, which has the following format:

	| 31 - 20 | 19 - 15 | 14 - 10 | 9 - 5 | 4 | 3 | 2 - 0 |
	|---------|---------|---------|-------|---|---|-------|
	|   op    |  cond   |    fk   |   fj  | 0 | 0 |  cd   |

The `cond` field has these possible values:

	"CAF": "00",
	"CUN": "08",
	"CEQ": "04",
	"CUEQ": "0c",
	"CLT": "02",
	"CULT": "0a",
	"CLE": "06",
	"CULE": "0e",
	"CNE": "10",
	"COR": "14",
	"CUNE": "18",
	"SAF": "01",
	"SUN": "09",
	"SEQ": "05",
	"SUEQ": "0d",
	"SLT": "03",
	"SULT": "0b",
	"SLE": "07",
	"SULE": "0f",
	"SNE": "11",
	"SOR": "15",
	"SUNE": "19",

These values are the hexadecimal numbers of bits 19 to 15, the same as
described in the instruction set manual.

The following code defines a map, the values in it represent the hexadecimal
encoding of the cond field in the entire instruction. In this case, the upper
4 bits and the lowest 1 bit are encoded separately, so the encoding is
different from the encoding described above.
*/
func dealWithFcmp(ds string) (fcmpConditions map[string]map[string]string) {
	conds := map[string]string{
		"CAF":  "00",
		"CUN":  "40",
		"CEQ":  "20",
		"CUEQ": "60",
		"CLT":  "10",
		"CULT": "50",
		"CLE":  "30",
		"CULE": "70",
		"CNE":  "80",
		"COR":  "a0",
		"CUNE": "c0",
		"SAF":  "08",
		"SUN":  "48",
		"SEQ":  "28",
		"SUEQ": "68",
		"SLT":  "18",
		"SULT": "58",
		"SLE":  "38",
		"SULE": "78",
		"SNE":  "88",
		"SOR":  "a8",
		"SUNE": "c8",
	}
	fcmpConditions = make(map[string]map[string]string)
	for k, v := range conds {
		op := fmt.Sprintf("FCMP_%s_%s", k, ds)
		opstr := fmt.Sprintf("FCMP_%s_%s:\t\"FCMP.%s.%s\",", k, ds, k, ds)
		instFormatComment := fmt.Sprintf("// FCMP.%s.%s cd, fj, fk", k, ds)
		var instFormat string
		if ds == "D" {
			instFormat = fmt.Sprintf("{mask: 0xffff8018, value: 0x0c2%s000, args: instArgs{arg_cd, arg_fj, arg_fk}},", v)
		} else {
			instFormat = fmt.Sprintf("{mask: 0xffff8018, value: 0x0c1%s000, args: instArgs{arg_cd, arg_fj, arg_fk}},", v)
		}

		fcmpConditions[op] = make(map[string]string)
		fcmpConditions[op]["op"] = op
		fcmpConditions[op]["opstr"] = opstr
		fcmpConditions[op]["instFormatComment"] = instFormatComment
		fcmpConditions[op]["instFormat"] = instFormat
	}
	return
}

func findWords(chars []pdf.Text) (words []pdf.Text) {
	for i := 0; i < len(chars); {
		xRange := []float64{chars[i].X, chars[i].X}
		j := i + 1

		// Find all chars on one line.
		for j < len(chars) && chars[j].Y == chars[i].Y {
			xRange[1] = chars[j].X
			j++
		}

		// we need to note that the word may change line(Y) but belong to one cell. So, after loop over all continued
		// chars whose Y are same, check if the next char's X belong to the range of xRange, if true, means it should
		// be contact to current word, because the next word's X should bigger than current one.
		for j < len(chars) && chars[j].X >= xRange[0] && chars[j].X <= xRange[1] {
			j++
		}

		var end float64
		// Split line into words (really, phrases).
		for k := i; k < j; {
			ck := &chars[k]
			s := ck.S
			end = ck.X + ck.W
			charSpace := ck.FontSize / 6
			wordSpace := ck.FontSize * 2 / 3
			l := k + 1
			for l < j {
				// Grow word.
				cl := &chars[l]

				if math.Abs(cl.FontSize-ck.FontSize) < 0.1 && cl.X <= end+charSpace {
					s += cl.S
					end = cl.X + cl.W
					l++
					continue
				}
				// Add space to phrase before next word.
				if math.Abs(cl.FontSize-ck.FontSize) < 0.1 && cl.X <= end+wordSpace {
					s += " " + cl.S
					end = cl.X + cl.W
					l++
					continue
				}
				break
			}
			f := ck.Font
			words = append(words, pdf.Text{
				Font:     f,
				FontSize: ck.FontSize,
				X:        ck.X,
				Y:        ck.Y,
				W:        end - ck.X,
				S:        s,
			})
			k = l
		}
		i = j
	}

	return words
}

func parsePage(num int, p pdf.Page, isFP bool) (ops []string, opstrs map[string]string, instFormatComments map[string]string, instFormats map[string]string) {
	opstrs = make(map[string]string)
	instFormatComments = make(map[string]string)
	instFormats = make(map[string]string)

	content := p.Content()

	var text []pdf.Text
	for _, t := range content.Text {
		text = append(text, t)
	}

	// table name(70), table header(64), page num(3)
	if isFP {
		text = text[134 : len(text)-3]
	} else {
		text = text[64 : len(text)-3]
	}

	text = findWords(text)

	for i := 0; i < len(text); {
		var fcmpConditions map[string]map[string]string
		if strings.HasPrefix(text[i].S, "FCMP") {
			fcmpConditions = dealWithFcmp(strings.Split(text[i].S, ".")[2])

			for fc, inst := range fcmpConditions {
				ops = append(ops, inst["op"])
				opstrs[fc] = inst["opstr"]
				instFormatComments[fc] = inst["instFormatComment"]
				instFormats[fc] = inst["instFormat"]
			}
			t := i + 1
			for ; text[t].Y == text[i].Y; t++ {
				continue
			}
			i = t
			continue
		}

		op := strings.Replace(text[i].S, ".", "_", -1)
		opstr := fmt.Sprintf("%s:\t\"%s\",", op, text[i].S)
		instFormatComment := ""
		binValue := ""
		binMask := ""
		instArgs := ""
		offs := false
		var offArgs []string

		j := i + 1
		for ; j < len(text) && text[j].Y == text[i].Y; j++ {

			// Some instruction has no arguments, so the next word(text[j].S) is not the arguments string but 0/1 bit, it shouldn't be skipped.
			if res, _ := regexp.MatchString("^\\d+$", text[j].S); j == i+1 && res == false {
				instFormatComment = fmt.Sprintf("// %s %s", text[i].S, strings.Replace(text[j].S, ",", ", ", -1))
				continue
			}
			if text[j].S == "0" || text[j].S == "1" {
				binValue += text[j].S
				binMask += "1"
			} else {
				argLen, argName := getArg(text[j].S)

				// Get argument's length failed, compute it by other arguments.
				if argLen == 0 {
					left := 31 - len(binValue)
					right := 0
					l := j + 1
					if l < len(text) && text[l].Y == text[j].Y {
						for ; text[l].Y == text[j].Y; l++ {
							if text[l].S == "0" || text[l].S == "1" {
								right += 1
							} else {
								tArgLen, _ := getArg(text[l].S)
								if tArgLen == 0 {
									fmt.Fprintf(os.Stderr, "there are more than two args whose length is unknown.\n")
								}
								right += tArgLen
							}
						}
					}
					argLen = left - right + 1
					argName = "arg_" + text[j].S + "_" + strconv.FormatInt(int64(left), 10) + "_" + strconv.FormatInt(int64(right), 10)
				}

				for k := 0; k < argLen; k++ {
					binValue += "0"
					binMask += "0"
				}

				if argName != "offs" {
					if instArgs != "" {
						instArgs = ", " + instArgs
					}
					instArgs = argName + instArgs
				} else {
					offs = true
					offArgs = append(offArgs, text[j].S)
				}
			}
		}

		// The real offset is a combination of two offsets in the binary code of the instruction, for example: BEQZ
		if offs && offArgs != nil {
			var left int
			var right int
			if len(offArgs) == 1 {
				left, _ = strconv.Atoi(strings.Split(strings.Split(offArgs[0], ":")[0], "[")[1])
				right, _ = strconv.Atoi(strings.Split(strings.Split(offArgs[0], ":")[1], "]")[0])
			} else if len(offArgs) == 2 {
				left, _ = strconv.Atoi(strings.Split(strings.Split(offArgs[1], ":")[0], "[")[1])
				right, _ = strconv.Atoi(strings.Split(strings.Split(offArgs[0], ":")[1], "]")[0])
			}

			if instArgs == "" {
				instArgs = fmt.Sprintf("arg_offset_%d_%d", left, right)
			} else {
				instArgs += fmt.Sprintf(", arg_offset_%d_%d", left, right)
			}
		}

		ops = append(ops, op)
		opstrs[op] = opstr
		if instFormatComment == "" {
			instFormatComment = "// " + text[i].S
		} else if strings.HasPrefix(op, "AM") {
			instFormatComment = fmt.Sprintf("// %s rd, rk, rj", text[i].S)
		}
		instFormatComments[op] = instFormatComment
		// The parameter order of some instructions is inconsistent in encoding and syntax, such as BSTRINS.*
		if instArgs != "" {
			args := strings.Split(instFormatComment, " ")[2:]
			tInstArgs := strings.Split(instArgs, ", ")
			newOrderedInstArgs := []string{}
			for _, a := range args {
				a = strings.Split(a, ",")[0]
				for _, aa := range tInstArgs {
					if strings.Contains(aa, a) {
						newOrderedInstArgs = append(newOrderedInstArgs, aa)
						break
					} else if a == "rd" && aa == "arg_fd" {
						newOrderedInstArgs = append(newOrderedInstArgs, "arg_rd")
						break
					}
				}
			}
			instArgs = strings.Join(newOrderedInstArgs, ", ")
		}
		if strings.HasPrefix(op, "AM") {
			instArgs = "arg_rd, arg_rk, arg_rj"
		}
		instFormat := fmt.Sprintf(
			"{mask: %s, value: %s, op: A%s, args: instArgs{%s}},",
			binstrToHex(binMask), binstrToHex(binValue),
			op, instArgs,
		)
		instFormats[op] = instFormat

		i = j // next instruction
	}

	return
}
