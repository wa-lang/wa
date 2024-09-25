// 版权 @2023 凹语言 作者。保留所有权利。

package appinit

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/format"
)

var CmdInit = &cli.Command{
	Name:  "init",
	Usage: "init a sketch Wa module",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "set app name",
			Value:   "hello",
		},
		&cli.StringFlag{
			Name:    "pkgpath",
			Aliases: []string{"p"},
			Usage:   "set pkgpath file",
			Value:   "myapp",
		},
		&cli.BoolFlag{
			Name:  "p5",
			Usage: "p5 example",
		},
		&cli.BoolFlag{
			Name:  "wasm4",
			Usage: "wasm4 game",
		},
		&cli.BoolFlag{
			Name:  "wasi",
			Usage: "wasi example",
		},
		&cli.BoolFlag{
			Name:  "arduino",
			Usage: "arduino nano 33 example",
		},
		&cli.BoolFlag{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update example",
		},
	},

	Action: func(c *cli.Context) error {
		err := InitApp(
			c.String("name"), c.String("pkgpath"),
			c.Bool("p5"), c.Bool("wasm4"), c.Bool("wasi"), c.Bool("arduino"),
			c.Bool("update"),
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

func InitApp(name, pkgpath string, isP5App, isWasm4App, isWasiApp, isArduinoApp, update bool) error {
	if name == "" {
		return fmt.Errorf("init failed: <%s> is empty", name)
	}
	if !appbase.IsValidAppName(name) {
		return fmt.Errorf("init failed: <%s> is invalid name", name)
	}

	if !appbase.IsValidPkgpath(pkgpath) {
		return fmt.Errorf("init failed: <%s> is invalid pkgpath", pkgpath)
	}

	if !update {
		if _, err := os.Lstat(name); err == nil {
			return fmt.Errorf("init failed: <%s> exists", name)
		}
	}

	if isP5App {
		isWasiApp = false
		isWasm4App = false
	}

	os.MkdirAll(name, 0777)
	os.WriteFile(filepath.Join(name, ".gitignore"), []byte(`!/output/index.html
/output/*.wat
/output/*.wasm
/output/*.js
`), 0666)

	var info = struct {
		Name         string
		Pkgpath      string
		Year         int
		IsP5App      bool
		IsWasm4App   bool
		IsWasiApp    bool
		IsArduinoApp bool
	}{
		Name:         name,
		Pkgpath:      pkgpath,
		Year:         time.Now().Year(),
		IsP5App:      isP5App,
		IsWasm4App:   isWasm4App,
		IsWasiApp:    isWasiApp,
		IsArduinoApp: isArduinoApp,
	}

	appFS := waroot_GetExampleAppFS()
	err := fs.WalkDir(appFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		data, err := fs.ReadFile(appFS, path)
		if err != nil {
			return err
		}

		tmpl, err := template.New(path).Parse(string(data))
		if err != nil {
			return err
		}

		dstpath := filepath.Join(name, path)
		os.MkdirAll(filepath.Dir(dstpath), 0777)

		f, err := os.Create(dstpath)
		if err != nil {
			return err
		}
		defer f.Close()

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, &info)
		if err != nil {
			return err
		}

		code, _, err := format.File(nil, path, buf.Bytes())
		if err != nil {
			code = buf.Bytes()
		}

		if _, err := f.Write(code); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 只有默认的 js 生成定制的 index.html
	if isP5App || isWasm4App || isWasiApp || isArduinoApp {
		os.Remove(filepath.Join(name, "output", "index.html"))
	}

	return nil
}
