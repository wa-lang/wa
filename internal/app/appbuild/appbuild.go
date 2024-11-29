// 版权 @2023 凹语言 作者。保留所有权利。

package appbuild

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/wat/watutil"
)

//go:embed assets/arduino.ino
var arduino_ino string

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
		appbase.MakeFlag_tags(),
		appbase.MakeFlag_ld_stack_size(),
		appbase.MakeFlag_ld_max_memory(),
		appbase.MakeFlag_optimize(),
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
	_, _, err := BuildApp(opt, input, outfile)
	return err
}

func BuildApp(opt *appbase.Option, input, outfile string) (mainFunc string, wasmBytes []byte, err error) {
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
		return "", wasmBytes, nil
	}

	// 只编译 wa/wz 文件, 输出路径相同, 后缀名调整
	if appbase.HasExt(input, ".wa", ".wz") {
		_, _, watOutput, err := buildWat(opt, input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if opt.Optimize {
			watOutput, err = watutil.WatStrip(input, watOutput)
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
		}

		// OK
		return "__main__.main", wasmBytes, nil
	}

	// 构建目录
	{
		if !appbase.IsNativeDir(input) {
			fmt.Printf("%q is not valid output path\n", outfile)
			os.Exit(1)
		}

		// 尝试读取模块信息
		manifest, err := config.LoadManifest(nil, input)
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
		_, compiler, watOutput, err := buildWat(opt, input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if s := manifest.Pkg.Target; opt.Optimize || s == config.WaOS_wasm4 || s == config.WaOS_arduino {
			watOutput, err = watutil.WatStrip(input, watOutput)
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
			appHeaderOutfile := filepath.Join(filepath.Dir(outfile), "arduino", "app.wasm.h")

			var buf bytes.Buffer
			// unsigned int app_wasm_len = ?;
			// unsigned char app_wasm[] = { 0x00, 0x01, ... };
			fmt.Fprintf(&buf, "// Auto Generate by Wa language. See https://wa-lang.org\n\n")
			fmt.Fprintf(&buf, "unsigned int app_wasm_len = %d;\n\n", len(wasmBytes))
			fmt.Fprintf(&buf, "unsigned char app_wasm[] = {")
			for i, ch := range wasmBytes {
				if i%10 == 0 {
					fmt.Fprintf(&buf, "\n\t0x%02x,", ch)
				} else {
					fmt.Fprintf(&buf, " 0x%02x,", ch)
					if i == len(wasmBytes)-1 {
						fmt.Fprintln(&buf)
					}
				}
			}
			fmt.Fprintf(&buf, "};\n")

			os.MkdirAll(arduinoDir, 0777)
			os.WriteFile(inoOutfile, []byte(arduino_ino), 0666)
			os.WriteFile(appHeaderOutfile, buf.Bytes(), 0666)

			// wasm 写到文件
			err = os.WriteFile(outfile, wasmBytes, 0666)
			if err != nil {
				fmt.Printf("write %s failed: %v\n", outfile, err)
				os.Exit(1)
			}

		default:

			// wasm 写到文件
			err = os.WriteFile(outfile, wasmBytes, 0666)
			if err != nil {
				fmt.Printf("write %s failed: %v\n", outfile, err)
				os.Exit(1)
			}
		}

		// 主函数
		mainFunc := manifest.MainPkg + ".main"

		// OK
		return mainFunc, wasmBytes, nil
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
