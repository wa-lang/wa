// 版权 @2023 凹语言 作者。保留所有权利。

package appbuild_mini

import (
	"fmt"
	"os"
	"path/filepath"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
)

var CmdBuild = &cli.Command{
	Name:  "build",
	Usage: "compile Wa source code",
	Flags: []cli.Flag{
		appbase.MakeFlag_output(),
		appbase.MakeFlag_target(),
		appbase.MakeFlag_tags(),
		appbase.MakeFlag_ld_stack_size(),
		appbase.MakeFlag_ld_max_memory(),
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
	if err != nil {
		return err
	}
	return err
}

func BuildApp(opt *appbase.Option, input, outfile string) (mainFunc string, wasmBytes []byte, err error) {
	// 路径是否存在
	if !appbase.PathExists(input) {
		fmt.Printf("%q not found\n", input)
		os.Exit(1)
	}

	// 输出参数是否合法, 必须是 wat
	if outfile != "" && !appbase.HasExt(outfile, ".wat") {
		fmt.Printf("%q is not valid output path\n", outfile)
		os.Exit(1)
	}

	// 只编译 wa/wz 文件, 输出路径相同, 后缀名调整
	if appbase.HasExt(input, ".wa", ".wz") {
		_, _, watOutput, err := buildWat(opt, input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 设置默认输出目标
		if outfile == "" {
			outfile = appbase.ReplaceExt(input, ".wa", ".wasm")
		}

		// wat 写到文件
		watOutfile := appbase.ReplaceExt(outfile, ".wasm", ".wat")
		err = os.WriteFile(watOutfile, watOutput, 0666)
		if err != nil {
			fmt.Printf("write %s failed: %v\n", outfile, err)
			os.Exit(1)
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

		// wat 写到文件
		watOutfile := appbase.ReplaceExt(outfile, ".wasm", ".wat")
		err = os.WriteFile(watOutfile, watOutput, 0666)
		if err != nil {
			fmt.Printf("write %s failed: %v\n", outfile, err)
			os.Exit(1)
		}

		// 生成 js 胶水代码
		if opt.TargetOS == config.WaOS_js {
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
