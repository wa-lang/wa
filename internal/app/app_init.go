// 版权 @2023 凹语言 作者。保留所有权利。

package app

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"wa-lang.org/wa/waroot"
)

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
