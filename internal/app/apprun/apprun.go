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

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/app/appbuild"
	"wa-lang.org/wa/internal/wazero"
)

//go:embed favicon.ico
var favicon []byte

var CmdRun = &cli.Command{
	Name:  "run",
	Usage: "compile and run Wa program",
	Flags: []cli.Flag{
		appbase.MakeFlag_wabt(),
		appbase.MakeFlag_target(),
		appbase.MakeFlag_tags(),
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
	outfile := ""

	if input == "" || input == "." {
		input, _ = os.Getwd()
	}

	var opt = appbase.BuildOptions(c)
	mainFunc, wasmBytes, err := appbuild.BuildApp(opt, input, outfile)
	if err != nil {
		fmt.Println("appbuild.BuildApp:", err)
		os.Exit(1)
		return nil
	}

	var appArgs []string
	if c.NArg() > 1 {
		appArgs = c.Args().Slice()[1:]
	}

	m, err := wazero.BuildModule(input, wasmBytes, appArgs...)
	if err != nil {
		fmt.Println("wazero.BuildModule:", err)
		os.Exit(1)
		return nil
	}
	defer m.Close()

	// Web 模式启动服务器
	if (m.HasUnknownImportFunc() || c.Bool("web")) && !c.Bool("console") {
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
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	}

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
	}
	if len(stdout) > 0 {
		fmt.Fprint(os.Stdout, string(stdout))
	}
	if len(stderr) > 0 {
		fmt.Fprint(os.Stderr, string(stderr))
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
