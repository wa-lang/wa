// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package native_x64

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/asm"
	"wa-lang.org/wa/internal/native/link"
	"wa-lang.org/wa/internal/native/wat2x64"
	"wa-lang.org/wa/internal/native/wemu/device/dram"
)

func BuildApp_wa_wz(
	opt *appbase.Option, input, outfile string,
	prog *loader.Program, watOutput []byte,
	isZhLang bool,
) (exePath string, err error) {
	// wa build -arch=x64 -target=linux input.wat
	// wa build -arch=x64 -target=windows input.wat

	if outfile == "" {
		panic("unreachable")
	}

	if s := opt.TargetOS; s != "" && s != config.WaOS_linux && s != config.WaOS_windows {
		panic(fmt.Sprintf("x64 donot support %s", s))
	}

	cpuType := abi.X64Unix
	if strings.EqualFold(opt.TargetOS, config.WaOS_windows) {
		cpuType = abi.X64Windows
	}

	// 设置默认输出目标
	var nativeExeFile = outfile
	var nativeAsmFile = nativeExeFile + ".s"
	var clangFilename = nativeExeFile + ".c"
	var gccArgsFilename = nativeExeFile + ".gcc.args.txt"

	// GCC 配置文件
	gccArgcContent := prog.GccArgsCode()
	if len(gccArgcContent) > 0 {
		err = os.WriteFile(gccArgsFilename, []byte(gccArgcContent), 0666)
		if err != nil {
			fmt.Printf("write %s failed: %v\n", gccArgsFilename, err)
			os.Exit(1)
		}
	}

	// 入口函数
	// 如果依然gcc参数, 则通过 main 函数入口
	nativeEntryFuncName := ""
	if len(gccArgcContent) != 0 {
		nativeEntryFuncName = "main"
	}

	// 将 wat 翻译为本地汇编代码
	_, nasmBytes, err := wat2x64.Wat2X64(input, watOutput, cpuType, nativeEntryFuncName)
	if err != nil {
		fmt.Printf("wat2x64 %s failed: %v\n", input, err)
		os.Exit(1)
	}

	// 拼接本地的汇编代码
	nasmCode := prog.NasmCode()
	if len(nasmCode) > 0 {
		nasmBytes = append(nasmBytes, nasmCode...)
	}

	// 本地C代码
	clangCode := prog.ClangCode()
	if len(clangCode) > 0 {
		err = os.WriteFile(clangFilename, []byte(clangCode), 0666)
		if err != nil {
			fmt.Printf("write %s failed: %v\n", clangFilename, err)
			os.Exit(1)
		}
	}

	// 保存生成的汇编语言文件
	err = os.WriteFile(nativeAsmFile, []byte(nasmBytes), 0666)
	if err != nil {
		fmt.Printf("write %s failed: %v\n", nativeAsmFile, err)
		os.Exit(1)
	}

	switch cpuType {
	case abi.X64Unix:
		gccEnabled := true
		if runtime.GOARCH != "amd64" || runtime.GOOS != "linux" {
			gccEnabled = false
		}

		if gccEnabled || len(gccArgcContent) > 0 {
			args := []string{
				nativeAsmFile,
				"-o", nativeExeFile,
				"-static",
				"-z", "noexecstack",
			}
			if len(gccArgcContent) > 0 {
				args = append(args, "@"+gccArgsFilename)
			} else {
				args = append(args, "-nostdlib")
			}
			if len(clangCode) > 0 {
				args = append(args, clangFilename)
			}
			gccCmd := exec.Command("gcc", args...)
			output, err := gccCmd.CombinedOutput()
			if err != nil {
				fmt.Println(string(output))
				fmt.Printf("gcc build failed: %v\n", err)
				os.Exit(1)
			}
		} else {
			// 汇编为 elf 可执行文件
			opt := &abi.LinkOptions{}
			opt.CPU = abi.X64Unix
			opt.DRAMBase = dram.DRAM_BASE_X64_LINUX
			opt.DRAMSize = dram.DRAM_SIZE // 16MB, 临时用于演示

			// 解析汇编程序, 并生成对应cpu的机器码
			prog, err := asm.AssembleFile(nativeAsmFile, nasmBytes, opt)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// 自研工具链保存到ELF格式文件
			elfBytes, err := link.LinkELF(prog)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err := os.WriteFile(nativeExeFile, elfBytes, 0777); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

	case abi.X64Windows:
		gccEnabled := true
		if runtime.GOARCH != "amd64" || runtime.GOOS != "windows" {
			gccEnabled = false
		}

		if gccEnabled || len(gccArgcContent) > 0 {
			var args = []string{
				nativeAsmFile,
				"-o", nativeExeFile,
				"-lkernel32",
				"-luser32",
			}
			if len(gccArgcContent) > 0 {
				args = append(args, "@"+gccArgsFilename)
			} else {
				args = append(args, "-nostdlib", "-Wl,-e,_start")
			}
			if len(clangCode) > 0 {
				args = append(args, clangFilename)
			}
			output, err := exec.Command("gcc", args...).CombinedOutput()
			if err != nil {
				fmt.Println(string(output))
				fmt.Printf("gcc build failed: %v\n", err)
				os.Exit(1)
			}
		} else {
			// 汇编为 pe 可执行文件
			opt := &abi.LinkOptions{}
			opt.CPU = abi.X64Windows
			opt.DRAMBase = dram.DRAM_BASE_X64_WINDOWS
			opt.DRAMSize = dram.DRAM_SIZE // 16MB, 临时用于演示

			// 解析汇编程序, 并生成对应cpu的机器码
			prog, err := asm.AssembleFile(nativeAsmFile, nasmBytes, opt)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// 包存到ELF格式文件
			peBytes, err := link.LinkEXE(prog)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err := os.WriteFile(nativeExeFile, peBytes, 0777); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

	default:
		panic("unreachable")
	}

	// OK
	return nativeExeFile, nil
}
