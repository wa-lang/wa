// 版权 @2019 凹语言 作者。保留所有权利。

package app

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/wa-lang/wa/internal/waroot"
)

// 命令行选项
type Option struct {
	Debug      bool
	TargetArch string
	TargetOS   string
	Clang      string
	WasmLLC    string
	WasmLD     string
}

// 命令行程序对象
type App struct {
	opt  Option
	path string
	src  string
}

// 构建命令行程序对象
func NewApp(opt *Option) *App {
	return &App{}
}

func (p *App) InitApp(name, pkgpath string, update bool) error {
	if name == "" {
		return fmt.Errorf("init failed: <%s> is empty", name)
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

	fileSys := waroot.GetExampleAppFS()
	return fs.WalkDir(fileSys, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		data, err := fs.ReadFile(fileSys, path)
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
}

func (p *App) Fmt(path string) error {
	panic("TODO")
}

func (p *App) Lex(filename string) error {
	return nil
}

func (p *App) AST(filename string) error {
	panic("TODO")
}

func (p *App) SSA(filename string) error {
	panic("TODO")
}

func (p *App) ASM(filename string) error {
	panic("TODO")
}

func (p *App) Build(filename string, src interface{}, outfile string) (output []byte, err error) {
	panic("TODO")
}

func (p *App) Run(filename string, src interface{}) (data []byte, err error) {
	panic("TODO")
}
