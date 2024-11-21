// 版权 @2023 凹语言 作者。保留所有权利。

package apprun

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/app/appbuild"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/wat/watutil"
	"wa-lang.org/wa/internal/wazero"
)

//go:embed favicon.ico
var favicon []byte

var CmdRun = &cli.Command{
	Name:  "run",
	Usage: "compile and run Wa program",
	Flags: []cli.Flag{
		appbase.MakeFlag_target(),
		appbase.MakeFlag_tags(),
		appbase.MakeFlag_optimize(),
		&cli.BoolFlag{
			Name:  "web",
			Usage: "set web mode",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "console",
			Usage: "set console mode",
			Value: false,
		},
		&cli.StringFlag{
			Name:  "http",
			Usage: "set http address",
			Value: ":8000",
		},
	},
	Action: CmdRunAction,
}

func CmdRunAction(c *cli.Context) error {
	input := c.Args().First()

	if input == "" || input == "." {
		input, _ = os.Getwd()
	}

	if appbase.HasExt(input, ".wat") {
		watBytes, err := os.ReadFile(input)
		if err != nil {
			return err
		}
		wasmBytes, err := watutil.Wat2Wasm(input, watBytes)
		if err != nil {
			return err
		}
		var appArgs []string
		if args := c.Args().Slice(); len(args) > 2 {
			appArgs = args[2:]
		}
		return runWasm(input, wasmBytes, appArgs...)
	}

	if appbase.HasExt(input, ".wasm") {
		var appArgs []string
		if args := c.Args().Slice(); len(args) > 2 {
			appArgs = args[2:]
		}
		return runWasm(input, nil, appArgs...)
	}

	var opt = appbase.BuildOptions(c)
	if appbase.HasExt(input, ".wa") {
		// 执行单个 wa 脚本, 避免写磁盘
		opt.RunFileMode = true
	}

	mainFunc, wasmBytes, err := appbuild.BuildApp(opt, input, "")
	if err != nil {
		fmt.Println("appbuild.BuildApp:", err)
		os.Exit(1)
		return nil
	}

	var appArgs []string
	if c.NArg() > 1 {
		appArgs = c.Args().Slice()[1:]
	}

	// 根据导入的宿主类型判断是否为 Web 模式
	wasmImportModuleName, _ := wazero.ReadImportModuleName(wasmBytes)

	// Web 模式启动服务器
	if !c.Bool("console") && !opt.RunFileMode && (c.Bool("web") || wasmImportModuleName != config.WaOS_wasi) {
		var addr = c.String("http")
		if strings.HasPrefix(addr, ":") {
			addr = "localhost" + addr
		}
		fmt.Printf("listen at http://%s\n", addr)

		http.Handle(
			"/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Add("Cache-Control", "no-cache")
				w.Header().Set("content-type", "image/x-icon")
				w.Header().Set("Access-Control-Allow-Origin", "*")

				w.Write(favicon)
			}),
		)

		fileHandler := http.FileServer(http.Dir(filepath.Join(input, "output")))
		http.Handle(
			"/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Add("Cache-Control", "no-cache")
				if strings.HasSuffix(req.URL.Path, ".wasm") {
					w.Header().Set("content-type", "application/wasm")
				}
				w.Header().Set("Access-Control-Allow-Origin", "*")

				fileHandler.ServeHTTP(w, req)
			}),
		)

		go func() {
			time.Sleep(time.Second * 2)
			openBrowser(addr)
		}()

		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	}

	// 单文件执行时删除中间临时文件
	if !c.Bool("debug") {
		if appbase.HasExt(input, ".wa") {
			os.Remove(input + "t")  // *.wat
			os.Remove(input + "sm") // *.wasm
		}
	}

	m, err := wazero.BuildModule(input, wasmBytes, appArgs...)
	if err != nil {
		fmt.Println("wazero.BuildModule:", err)
		os.Exit(1)
		return nil
	}
	defer m.Close()

	stdout, stderr, err := m.RunMain(mainFunc)
	if err != nil {
		if len(stdout) > 0 {
			fmt.Fprint(os.Stdout, string(stdout))
		}
		if len(stderr) > 0 {
			fmt.Fprint(os.Stderr, string(stderr))
		}
		if exitCode, ok := wazero.AsExitError(err); ok {
			os.Exit(exitCode)
		}
		fmt.Println(err)
	} else {
		if len(stdout) > 0 {
			fmt.Fprint(os.Stdout, string(stdout))
		}
		if len(stderr) > 0 {
			fmt.Fprint(os.Stderr, string(stderr))
		}
	}

	return nil
}

func openBrowser(url string) error {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}

func runWasm(input string, wasmBytes []byte, args ...string) error {
	var err error
	if wasmBytes == nil {
		if wasmBytes, err = os.ReadFile(input); err != nil {
			return err
		}
	}

	stdout, stderr, err := wazero.RunWasm(input, wasmBytes, "_main", args...)
	if err != nil {
		if len(stdout) > 0 {
			fmt.Fprint(os.Stdout, string(stdout))
		}
		if len(stderr) > 0 {
			fmt.Fprint(os.Stderr, string(stderr))
		}
		if exitCode, ok := wazero.AsExitError(err); ok {
			os.Exit(exitCode)
		}
		fmt.Println(err)
		return nil
	}
	if len(stdout) > 0 {
		fmt.Fprint(os.Stdout, string(stdout))
	}
	if len(stderr) > 0 {
		fmt.Fprint(os.Stderr, string(stderr))
	}
	return nil
}
