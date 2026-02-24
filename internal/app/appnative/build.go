// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appnative

import (
	"fmt"
	"os"
	"path/filepath"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/app/appnative/native_loong64"
	"wa-lang.org/wa/internal/app/appnative/native_x64"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/wat/watutil/watstrip"
)

func CmdBuildAction(c *cli.Context) error {
	input := c.Args().First()
	outfile := c.String("output")

	if input == "" {
		input, _ = os.Getwd()
	}

	var opt = appbase.BuildOptions(c)
	_, err := BuildApp(opt, input, outfile)
	return err
}

func BuildApp(opt *appbase.Option, input, outfile string) (exePath string, err error) {
	// 路径是否存在
	if !appbase.PathExists(input) {
		fmt.Printf("%q not found\n", input)
		os.Exit(1)
	}

	// 输出参数是否合法, 必须是 wasm
	if outfile != "" && !appbase.HasExt(outfile, ".wasm", ".s") {
		fmt.Printf("%q is not valid output path\n", outfile)
		os.Exit(1)
	}

	// 只编译 wa/wz 文件, 输出路径相同, 后缀名调整
	if appbase.HasExt(input, ".wa", ".wz") {
		prog, _, watOutput, err := buildWat(opt, input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if opt.Optimize {
			watOutput, err = watstrip.WatStrip(input, watOutput)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// 设置默认输出目标
		if outfile == "" {
			if appbase.HasExt(input, ".wz") {
				outfile = appbase.ReplaceExt(input, ".wz", ".wasm")
			} else {
				outfile = appbase.ReplaceExt(input, ".wa", ".wasm")
			}
		}

		// wat 写到文件
		watOutfile := appbase.ReplaceExt(outfile, ".wasm", ".wat")
		if !opt.RunFileMode {
			err = os.WriteFile(watOutfile, watOutput, 0666)
			if err != nil {
				fmt.Printf("write %s failed: %v\n", outfile, err)
				os.Exit(1)
			}
		}

		// 构建本地可执行程序
		switch opt.TargetArch {
		case config.WaArch_loong64:
			// wa build -arch=loong64 -target=linux input.wat
			return native_loong64.BuildApp_wa_wz(opt, input, outfile, prog, watOutput)

		case config.WaArch_x64:
			// wa build -arch=x64 -target=linux input.wat
			// wa build -arch=x64 -target=windows input.wat
			return native_x64.BuildApp_wa_wz(opt, input, outfile, prog, watOutput)

		default:
			err = fmt.Errorf("unknown tarch: %v", opt.TargetArch)
			return "", err
		}
	}

	// 构建目录
	{
		if !appbase.IsNativeDir(input) {
			fmt.Printf("%q is not valid output path\n", outfile)
			os.Exit(1)
		}

		// 尝试读取模块信息
		manifest, err := config.LoadManifest(nil, input, false)
		if err != nil {
			fmt.Printf("%q is invalid wa moudle\n", input)
			os.Exit(1)
		}
		if opt.TargetOS != "" {
			manifest.Pkg.TargetOS = opt.TargetOS
		}
		if err := manifest.Valid(); err != nil {
			fmt.Printf("%q is invalid wa module; %v\n", input, err)
			os.Exit(1)
		}

		if outfile == "" {
			if !manifest.IsStd {
				outfile = filepath.Join(manifest.Root, "output", manifest.Pkg.Name) + ".wasm"
				os.MkdirAll(filepath.Join(manifest.Root, "output"), 0777)
			} else {
				outfile = "a.out." + manifest.Pkg.Name + ".wasm"
			}
		}

		// 编译出 wat 文件
		prog, _, watOutput, err := buildWat(opt, input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 用返回的 prog 替代, 状态可能有更新(W2Mode)
		manifest = prog.Manifest

		if s := manifest.Pkg.TargetOS; opt.Optimize || s == config.WaOS_wasm4 || s == config.WaOS_arduino {
			watOutput, err = watstrip.WatStrip(input, watOutput)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// wat 写到文件
		watOutfile := appbase.ReplaceExt(outfile, ".wasm", ".wat")
		err = os.WriteFile(watOutfile, watOutput, 0666)
		if err != nil {
			fmt.Printf("write %s failed: %v\n", outfile, err)
			os.Exit(1)
		}

		// 构建本地可执行程序
		switch opt.TargetArch {
		case config.WaArch_loong64:
			// wa build -arch=loong64 -target=linux input.wat
			return native_loong64.BuildApp_wa_wz(opt, input, outfile, prog, watOutput)

		case config.WaArch_x64:
			// wa build -arch=x64 -target=linux input.wat
			// wa build -arch=x64 -target=windows input.wat
			return native_x64.BuildApp_wa_wz(opt, input, outfile, prog, watOutput)

		default:
			err = fmt.Errorf("unknown tarch: %v", opt.TargetArch)
			return "", err
		}
	}
}

func buildWat(opt *appbase.Option, filename string) (
	prog *loader.Program, compiler *compiler_wat.Compiler,
	watBytes []byte, err error,
) {
	cfg := opt.Config()
	prog, err = loader.LoadProgram(cfg, filename)
	if err != nil {
		return prog, nil, nil, err
	}

	compiler = compiler_wat.New()
	output, err := compiler.Compile(prog)

	if err != nil {
		return prog, nil, nil, err
	}

	return prog, compiler, []byte(output), nil
}
