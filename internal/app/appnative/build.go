// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appnative

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/app/appnative/native_loong64"
	"wa-lang.org/wa/internal/app/appnative/native_x64"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/wat/watutil/watstrip"
)

var CmdNative_Build = &cli.Command{
	Name:  "build",
	Usage: "compile Wa source code in native mode",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set output file",
			Value:   "",
		},
		&cli.StringFlag{
			Name:  "arch",
			Usage: fmt.Sprintf("set arch type (%s|%s)", config.WaArch_x64, config.WaArch_loong64),
			Value: "",
		},
		&cli.StringFlag{
			Name:  "target",
			Usage: fmt.Sprintf("set target type (%s|%s)", config.WaOS_windows, config.WaOS_linux),
			Value: "",
		},
		&cli.StringFlag{
			Name:  "tags",
			Usage: "set build tags",
		},
		&cli.BoolFlag{
			Name:  "optimize",
			Usage: "enable optimize flag",
		},
	},

	Action: CmdBuildAction,
}

func CmdBuildAction(c *cli.Context) error {
	_, err := doCmdBuildAction(c)
	return err
}

func doCmdBuildAction(c *cli.Context) (exePath string, err error) {
	input := c.Args().First()
	outfile := c.String("output")

	if input == "" {
		input, _ = os.Getwd()
	}

	opt := &appbase.Option{
		Debug:     c.Bool("debug"),
		WaBackend: config.WaBackend_Default,
		BuilgTags: strings.Fields(c.String("tags")),
		Optimize:  c.Bool("optimize"),
	}

	targetArch := c.String("arch")
	targetOS := c.String("target")

	if targetArch != "" && !config.CheckWaArch(targetArch) {
		fmt.Printf("unknown arch: %s\n", targetArch)
		os.Exit(1)
	}
	if targetOS != "" && !config.CheckWaOS(targetOS) {
		fmt.Printf("unknown target: %s\n", targetOS)
		os.Exit(1)
	}

	// 处理默认值
	switch {
	case targetArch == "" && targetOS == "":
		switch {
		case runtime.GOARCH == "amd64" && runtime.GOOS == "windows":
			targetArch = config.WaArch_x64
			targetOS = config.WaOS_windows
		case runtime.GOARCH == "amd64" && runtime.GOOS == "linux":
			targetArch = config.WaArch_x64
			targetOS = config.WaOS_linux
		case runtime.GOARCH == "loong64":
			targetArch = config.WaArch_loong64
			targetOS = config.WaOS_linux
		}
	case targetArch == "":
		switch targetOS {
		case "windows":
			targetArch = config.WaArch_x64
		case "linux":
			switch runtime.GOARCH {
			case "amd64":
				targetArch = config.WaArch_x64
			case "loong64":
				targetArch = config.WaArch_loong64
			}
		}
	case targetOS == "":
		switch targetArch {
		case "x64", "amd64":
			switch runtime.GOOS {
			case "windwos":
				targetOS = config.WaOS_windows
			case "linux":
				targetOS = config.WaOS_linux
			}
		case "loong64":
			targetOS = config.WaOS_linux
		}
	}

	// 检查目标类型
	switch {
	case targetOS == config.WaOS_windows && targetArch == config.WaArch_x64:
		opt.TargetArch = config.WaArch_x64
		opt.TargetOS = config.WaOS_windows
	case targetOS == config.WaOS_linux && targetArch == config.WaArch_x64:
		opt.TargetArch = config.WaArch_x64
		opt.TargetOS = config.WaOS_linux
	case targetOS == config.WaOS_linux && targetArch == config.WaArch_loong64:
		opt.TargetArch = config.WaArch_loong64
		opt.TargetOS = config.WaOS_linux
	default:
		fmt.Printf("unsupport target: %s/%s", targetOS, targetArch)
		os.Exit(1)
	}

	return BuildApp(opt, input, outfile)
}

func BuildApp(opt *appbase.Option, input, outfile string) (exePath string, err error) {
	// 路径是否存在
	if !appbase.PathExists(input) {
		fmt.Printf("%q not found\n", input)
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
				outfile = appbase.ReplaceExt(input, ".wz", ".exe")
			} else {
				outfile = appbase.ReplaceExt(input, ".wa", ".exe")
			}
		}

		// 构建本地可执行程序
		switch opt.TargetArch {
		case config.WaArch_loong64:
			// wa build -arch=loong64 -target=linux input.wat
			isZhLang := appbase.HasExt(input, ".wz")
			return native_loong64.BuildApp_wa_wz(opt, input, outfile, prog, watOutput, isZhLang)

		case config.WaArch_x64:
			// wa build -arch=x64 -target=linux input.wat
			// wa build -arch=x64 -target=windows input.wat
			isZhLang := appbase.HasExt(input, ".wz")
			return native_x64.BuildApp_wa_wz(opt, input, outfile, prog, watOutput, isZhLang)

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
				outfile = filepath.Join(manifest.Root, "output", manifest.Pkg.Name) + ".exe"
				os.MkdirAll(filepath.Join(manifest.Root, "output"), 0777)
			} else {
				outfile = "a.out." + manifest.Pkg.Name + ".exe"
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

		// 构建本地可执行程序
		switch opt.TargetArch {
		case config.WaArch_loong64:
			// wa build -arch=loong64 -target=linux input.wat
			return native_loong64.BuildApp_wa_wz(opt, input, outfile, prog, watOutput, manifest.W2Mode)

		case config.WaArch_x64:
			// wa build -arch=x64 -target=linux input.wat
			// wa build -arch=x64 -target=windows input.wat
			return native_x64.BuildApp_wa_wz(opt, input, outfile, prog, watOutput, manifest.W2Mode)

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
