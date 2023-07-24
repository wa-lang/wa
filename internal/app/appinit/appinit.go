// 版权 @2023 凹语言 作者。保留所有权利。

package appinit

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/waroot"
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
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update example",
		},
	},

	Action: func(c *cli.Context) error {
		err := InitApp(c.String("name"), c.String("pkgpath"), c.Bool("update"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

func InitApp(name, pkgpath string, update bool) error {
	if name == "" {
		return fmt.Errorf("init failed: <%s> is empty", name)
	}
	if !isValidAppName(name) {
		return fmt.Errorf("init failed: <%s> is invalid name", name)
	}

	if !isValidPkgpath(pkgpath) {
		return fmt.Errorf("init failed: <%s> is invalid pkgpath", pkgpath)
	}

	if !update {
		if _, err := os.Lstat(name); err == nil {
			return fmt.Errorf("init failed: <%s> exists", name)
		}
	}

	var info = struct {
		Name    string
		Pkgpath string
		Year    int
	}{
		Name:    name,
		Pkgpath: pkgpath,
		Year:    time.Now().Year(),
	}

	appFS := waroot.GetExampleAppFS()
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

		err = tmpl.Execute(f, &info)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	vendorFS := waroot.GetExampleVendorFS()
	err = fs.WalkDir(vendorFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		data, err := fs.ReadFile(vendorFS, path)
		if err != nil {
			return err
		}

		tmpl, err := template.New(path).Parse(string(data))
		if err != nil {
			return err
		}

		dstpath := filepath.Join(name, "vendor", path)
		os.MkdirAll(filepath.Dir(dstpath), 0777)

		f, err := os.Create(dstpath)
		if err != nil {
			return err
		}
		defer f.Close()

		err = tmpl.Execute(f, &info)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func isValidAppName(s string) bool {
	if s == "" || s[0] == '_' || (s[0] >= '0' && s[0] <= '9') {
		return false
	}
	for _, c := range []rune(s) {
		if c == '_' || (c >= '0' && c <= '9') || unicode.IsLetter(c) {
			continue
		}
		return false
	}
	return true
}

func isValidPkgpath(s string) bool {
	if s == "" || s[0] == '_' || (s[0] >= '0' && s[0] <= '9') {
		return false
	}
	for _, c := range []rune(s) {
		if c == '_' || c == '.' || c == '/' || (c >= '0' && c <= '9') {
			continue
		}
		if unicode.IsLetter(c) {
			continue
		}
		return false
	}

	var pkgname = s
	if i := strings.LastIndex(s, "/"); i >= 0 {
		pkgname = s[i+1:]
	}
	return isValidAppName(pkgname)
}
