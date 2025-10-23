// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appbuild

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/wat/watutil"
	"wa-lang.org/wa/internal/wat/watutil/wat2c"
	"wa-lang.org/wa/internal/wat/watutil/watstrip"
)

//go:embed assets/arduino.ino
var arduino_ino string

//go:embed assets/arduino-host.cpp
var arduino_host_cpp string

//go:embed assets/native-js-host.cpp
var native_js_host_cpp string

//go:embed assets/native.cpp
var native_main_cpp string

//go:embed assets/CMakeLists.txt
var cmakelists_txt string

//go:embed assets/favicon.ico
var favicon_ico string

//go:embed assets/index.html
var w4index_html string

//go:embed assets/wasm4.js
var w4js string

//go:embed assets/wasm4.css
var w4css string

var CmdBuild = &cli.Command{
	Name:  "build",
	Usage: "compile Wa source code",
	Flags: []cli.Flag{
		appbase.MakeFlag_output(),
		appbase.MakeFlag_target(),
		appbase.MakeFlag_wat2c_prefix(),
		appbase.MakeFlag_tags(),
		appbase.MakeFlag_ld_stack_size(),
		appbase.MakeFlag_ld_max_memory(),
		appbase.MakeFlag_optimize(),
		appbase.MakeFlag_wat2c_native(),
	},
	Action: CmdBuildAction,
}

func CmdBuildAction(c *cli.Context) error {
	input := c.Args().First()
	outfile := ""

	if input == "" {
		input, _ = os.Getwd()
	}

	var opt = appbase.BuildOptions(c)
	_, _, _, err := BuildApp(opt, input, outfile)
	return err
}

func BuildApp(opt *appbase.Option, input, outfile string) (mainFunc string, wasmBytes, fsetBytes []byte, err error) {
	// 路径是否存在
	if !appbase.PathExists(input) {
		fmt.Printf("%q not found\n", input)
		os.Exit(1)
	}

	// 输出参数是否合法, 必须是 wasm
	if outfile != "" && !appbase.HasExt(outfile, ".wasm") {
		fmt.Printf("%q is not valid output path\n", outfile)
		os.Exit(1)
	}

	// 已经是 wasm, 直接返回
	if appbase.HasExt(input, ".wasm") {
		wasmBytes, err = os.ReadFile(input)
		return
	}

	// 只编译 wat 文件, 输出路径相同, 后缀名调整
	if appbase.HasExt(input, ".wat") {
		// 设置默认输出目标
		if outfile == "" {
			outfile = appbase.ReplaceExt(input, ".wat", ".wasm")
		}

		watData, err := os.ReadFile(input)
		if err != nil {
			fmt.Printf("read %s failed: %v\n", input, err)
			os.Exit(1)
		}
		wasmBytes, err := watutil.Wat2Wasm(input, watData)
		if err != nil {
			fmt.Printf("wat2wasm %s failed: %v\n", input, err)
			os.Exit(1)
		}

		// 写到文件
		err = os.WriteFile(outfile, wasmBytes, 0666)
		if err != nil {
			fmt.Printf("write %s failed: %v\n", outfile, err)
			os.Exit(1)
		}

		// OK
		return "", wasmBytes, nil, nil
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
			outfile = appbase.ReplaceExt(input, ".wa", ".wasm")
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

		// wat 编译为 wasm
		wasmBytes, err := watutil.Wat2Wasm(watOutfile, watOutput)
		if err != nil {
			fmt.Printf("wat2wasm %s failed: %v\n", input, err)
			os.WriteFile(watOutfile, watOutput, 0666)
			os.Exit(1)
		}

		// wasm 写到文件
		if !opt.RunFileMode {
			err = os.WriteFile(outfile, wasmBytes, 0666)
			if err != nil {
				fmt.Printf("write %s failed: %v\n", outfile, err)
				os.Exit(1)
			}

			// fileset 写到文件中
			fsetfile := appbase.ReplaceExt(outfile, ".wasm", ".fset")
			os.WriteFile(fsetfile, prog.Fset.ToJson(), 0666)
		}

		// OK
		if appbase.HasExt(input, ".wz") {
			return token.K_pkg_主包 + "." + token.K_主控, wasmBytes, prog.Fset.ToJson(), nil
		}

		return token.K_pkg_main + "." + token.K_main, wasmBytes, prog.Fset.ToJson(), nil
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
		if opt.Target != "" {
			manifest.Pkg.Target = opt.Target
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
		prog, compiler, watOutput, err := buildWat(opt, input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 用返回的 prog 替代, 状态可能有更新(W2Mode)
		manifest = prog.Manifest

		if s := manifest.Pkg.Target; opt.Optimize || s == config.WaOS_wasm4 || s == config.WaOS_arduino {
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

		// fileset 写到文件中
		fsetfile := appbase.ReplaceExt(outfile, ".wasm", ".fset")
		fsetBytes := prog.Fset.ToJson()
		os.WriteFile(fsetfile, fsetBytes, 0666)

		// wat 编译为 wasm
		wasmBytes, err := watutil.Wat2Wasm(input, watOutput)
		if err != nil {
			fmt.Printf("wat2wasm %s failed: %v\n", input, err)
			os.Exit(1)
		}

		// 生成 js 胶水代码
		switch manifest.Pkg.Target {
		case config.WaOS_js:
			jsOutfile := appbase.ReplaceExt(outfile, ".wasm", ".js")
			jsOutput := compiler.GenJSBinding(filepath.Base(outfile))
			err = os.WriteFile(jsOutfile, []byte(jsOutput), 0666)
			if err != nil {
				fmt.Printf("write %s failed: %v\n", jsOutfile, err)
				os.Exit(1)
			}

			// 生成 index.html 文件
			indexHtmlPath := filepath.Join(filepath.Dir(outfile), "index.html")
			if !appbase.PathExists(indexHtmlPath) {
				htmlOutput := compiler.GenIndexHtml(filepath.Base(jsOutfile))
				err = os.WriteFile(indexHtmlPath, []byte(htmlOutput), 0666)
				if err != nil {
					fmt.Printf("write %s failed: %v\n", indexHtmlPath, err)
					os.Exit(1)
				}
			}

			// wasm 写到文件
			err = os.WriteFile(outfile, wasmBytes, 0666)
			if err != nil {
				fmt.Printf("write %s failed: %v\n", outfile, err)
				os.Exit(1)
			}

			// 通过 wat2c 生成 native c 代码
			if opt.Wat2CNative {
				nativeDir := filepath.Join(filepath.Dir(outfile), "native")

				mainCppOutfile := filepath.Join(filepath.Dir(outfile), "native", "main.cpp")
				cmakeOutfile := filepath.Join(filepath.Dir(outfile), "native", "CMakeLists.txt")
				hostCppOutfile := filepath.Join(filepath.Dir(outfile), "native", "native-host.cpp")
				waAppCfile := filepath.Join(filepath.Dir(outfile), "native", "wa-app.c")
				waAppHfile := filepath.Join(filepath.Dir(outfile), "native", "wa-app.h")

				os.MkdirAll(nativeDir, 0777)

				// 主文件
				if !appbase.PathExists(mainCppOutfile) {
					os.WriteFile(mainCppOutfile, []byte(native_main_cpp), 0666)
				}
				if !appbase.PathExists(cmakeOutfile) {
					os.WriteFile(cmakeOutfile, []byte(cmakelists_txt), 0666)
				}

				pkgName := manifest.Pkg.Name
				m, code, header, err := watutil.Wat2C("wa-app.wat", watOutput, wat2c.Options{
					Prefix: opt.Wat2CPrefix,
					Exports: map[string]string{
						pkgName + ".Loop": "loop",
					},
				})
				if err != nil {
					os.WriteFile(outfile, code, 0666)
					fmt.Println(err)
					os.Exit(1)
				}

				os.WriteFile(waAppCfile, []byte(code), 0666)
				os.WriteFile(waAppHfile, []byte(header), 0666)

				// 宿主的C代码
				if !appbase.PathExists(hostCppOutfile) {
					sMemoryBytes := "8" // 8 byte
					if m.Memory != nil && m.Memory.Pages > 0 {
						pages := m.Memory.Pages
						if m.Memory.MaxPages > pages {
							pages = m.Memory.MaxPages
						}
						sMemoryBytes = fmt.Sprintf("%d*(1<<16)", pages)
					}

					// 初始化 host 静态内存大小
					host_cpp_code := strings.ReplaceAll(native_js_host_cpp, "{{.MemoryBytes}}", sMemoryBytes)
					os.WriteFile(hostCppOutfile, []byte(host_cpp_code), 0666)
				}
			}

		case config.WaOS_wasm4:
			icoOutfile := filepath.Join(filepath.Dir(outfile), "favicon.ico")
			w4JsOutfile := filepath.Join(filepath.Dir(outfile), "wasm4.js")
			w4CssOutfile := filepath.Join(filepath.Dir(outfile), "wasm4.css")
			w4IndexOutfile := filepath.Join(filepath.Dir(outfile), "index.html")
			w4WasmOutfile := filepath.Join(filepath.Dir(outfile), "cart.wasm")

			os.WriteFile(icoOutfile, []byte(favicon_ico), 0666)
			os.WriteFile(w4JsOutfile, []byte(w4js), 0666)
			os.WriteFile(w4CssOutfile, []byte(w4css), 0666)
			os.WriteFile(w4IndexOutfile, []byte(w4index_html), 0666)
			os.WriteFile(w4WasmOutfile, wasmBytes, 0666)

		case config.WaOS_arduino:
			arduinoDir := filepath.Join(filepath.Dir(outfile), "arduino")
			inoOutfile := filepath.Join(filepath.Dir(outfile), "arduino", "arduino.ino")
			hostCppOutfile := filepath.Join(filepath.Dir(outfile), "arduino", "arduino-host.cpp")
			waAppCfile := filepath.Join(filepath.Dir(outfile), "arduino", "wa-app.c")
			waAppHfile := filepath.Join(filepath.Dir(outfile), "arduino", "wa-app.h")

			os.MkdirAll(arduinoDir, 0777)

			// 主文件
			if !appbase.PathExists(inoOutfile) {
				os.WriteFile(inoOutfile, []byte(arduino_ino), 0666)
			}

			// wasm 写到文件
			err = os.WriteFile(outfile, wasmBytes, 0666)
			if err != nil {
				fmt.Printf("write %s failed: %v\n", outfile, err)
				os.Exit(1)
			}

			// 生成wat转译的C代码
			{
				pkgName := manifest.Pkg.Name
				m, code, header, err := watutil.Wat2C("wa-app.wat", watOutput, wat2c.Options{
					Prefix: opt.Wat2CPrefix,
					Exports: map[string]string{
						pkgName + ".Loop": "loop",
					},
				})
				if err != nil {
					os.WriteFile(outfile, code, 0666)
					fmt.Println(err)
					os.Exit(1)
				}

				os.WriteFile(waAppCfile, []byte(code), 0666)
				os.WriteFile(waAppHfile, []byte(header), 0666)

				// 宿主的C代码
				if !appbase.PathExists(hostCppOutfile) {
					sMemoryBytes := "8" // 8 byte
					if m.Memory != nil && m.Memory.Pages > 0 {
						pages := m.Memory.Pages
						if m.Memory.MaxPages > pages {
							pages = m.Memory.MaxPages
						}
						sMemoryBytes = fmt.Sprintf("%d*(1<<16)", pages)
					}

					// 初始化 host 静态内存大小
					host_cpp_code := strings.ReplaceAll(arduino_host_cpp, "{{.MemoryBytes}}", sMemoryBytes)
					os.WriteFile(hostCppOutfile, []byte(host_cpp_code), 0666)
				}
			}

		default:

			// wasm 写到文件
			err = os.WriteFile(outfile, wasmBytes, 0666)
			if err != nil {
				fmt.Printf("write %s failed: %v\n", outfile, err)
				os.Exit(1)
			}

			// fileset 写到文件中
			fsetfile := appbase.ReplaceExt(outfile, ".wasm", ".fset")
			os.WriteFile(fsetfile, prog.Fset.ToJson(), 0666)
		}

		// 主函数
		var mainFunc string
		if manifest.W2Mode {
			mainFunc = manifest.MainPkg + "." + token.K_主控
		} else {
			mainFunc = manifest.MainPkg + "." + token.K_main
		}

		// OK
		return mainFunc, wasmBytes, fsetBytes, nil
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
