// 版权 @2023 凹语言 作者。保留所有权利。

package apprun

import (
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
	"wa-lang.org/wa/internal/wazero"
)

var CmdRun = &cli.Command{
	Name:  "run",
	Usage: "compile and run Wa program",
	Flags: []cli.Flag{
		appbase.MakeFlag_target(),
		appbase.MakeFlag_tags(),
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
		return err
	}

	// Web 模式启动服务器
	if opt.TargetOS == config.WaOS_js && appbase.IsNativeDir(input) {
		var addr = c.String("http")
		if strings.HasPrefix(addr, ":") {
			addr = "localhost" + addr
		}
		fmt.Printf("listen at http://%s\n", addr)

		go func() {
			time.Sleep(time.Second * 2)
			openBrowser(addr)
		}()

		http.Handle(
			"/", http.FileServer(http.Dir(filepath.Join(input, "output"))),
		)
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	}

	var appArgs []string
	if c.NArg() > 1 {
		appArgs = c.Args().Slice()[1:]
	}

	stdout, stderr, err := wazero.RunWasm(opt.Config(), input, wasmBytes, mainFunc, appArgs...)
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
