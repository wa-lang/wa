// 版权 @2019 凹语言 作者。保留所有权利。

// 凹语言打包程序: Windows/Linux/macOS/wasm
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"wa-lang.org/wa/internal/3rdparty/wabt-go"
	"wa-lang.org/wa/internal/version"
	"wa-lang.org/wa/waroot"
)

const (
	macos   = "macos"
	ubuntu  = "ubuntu"
	windows = "windows"
)

var (
	flagOutputDir = flag.String("output", "_output", "set output dir")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	flag.Parse()

	b := &Builder{Output: *flagOutputDir}
	b.GenAll()
}

type Builder struct {
	Output string // 目标的根目录
}

func NewBuilder(outputDir string) *Builder {
	return &Builder{Output: outputDir}
}

func (p *Builder) GenAll() {
	os.RemoveAll(p.getWarootPath(macos))
	os.RemoveAll(p.getWarootPath(ubuntu))
	os.RemoveAll(p.getWarootPath(windows))

	p.genWarootFiles(p.getWarootPath(macos))
	p.genWarootFiles(p.getWarootPath(ubuntu))
	p.genWarootFiles(p.getWarootPath(windows))

	p.genWat2wasmExe()
	p.genWaExe()
}

func (p *Builder) genWaExe() {
	// macos
	{
		waRootPath := p.getWarootPath(macos)
		dstpath := filepath.Join(waRootPath, "bin", "wa")

		cmd := exec.Command("go", "build", "-o", dstpath, "wa-lang.org/wa")
		cmd.Env = append(cmd.Env, os.Environ()...)
		cmd.Env = append(cmd.Env, "GOOS=darwin", "GOARCH=arm64")
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Print(string(output))
			panic(err)
		}
	}

	// ubuntu
	{
		waRootPath := p.getWarootPath(ubuntu)
		dstpath := filepath.Join(waRootPath, "bin", "wa")

		cmd := exec.Command("go", "build", "-o", dstpath, "wa-lang.org/wa")
		cmd.Env = append(cmd.Env, os.Environ()...)
		cmd.Env = append(cmd.Env, "GOOS=linux", "GOARCH=amd64")
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Print(string(output))
			panic(err)
		}
	}

	// windows
	{
		waRootPath := p.getWarootPath(windows)
		dstpath := filepath.Join(waRootPath, "bin", "wa.exe")

		cmd := exec.Command("go", "build", "-o", dstpath, "wa-lang.org/wa")
		cmd.Env = append(cmd.Env, os.Environ()...)
		cmd.Env = append(cmd.Env, "GOOS=windows", "GOARCH=amd64")
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Print(string(output))
			panic(err)
		}
	}
}

func (p *Builder) genWat2wasmExe() {
	// macos
	{
		waRootPath := p.getWarootPath(macos)
		dstpath := filepath.Join(waRootPath, "bin", wabt.Wat2WasmName)
		os.WriteFile(dstpath, wabt.LoadWat2Wasm(), 0777)
	}

	// ubuntu
	{
		waRootPath := p.getWarootPath(ubuntu)
		dstpath := filepath.Join(waRootPath, "bin", wabt.Wat2WasmName)
		os.WriteFile(dstpath, wabt.LoadWat2Wasm(), 0777)
	}

	// windows
	{
		waRootPath := p.getWarootPath(windows)
		dstpath := filepath.Join(waRootPath, "bin", wabt.Wat2WasmName)
		os.WriteFile(dstpath, wabt.LoadWat2Wasm(), 0777)
	}
}

func (p *Builder) genWarootFiles(waRootPath string) error {
	os.MkdirAll(waRootPath, 0777)
	warootfs := waroot.GetRootFS()
	err := fs.WalkDir(waroot.GetRootFS(), ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		data, err := fs.ReadFile(warootfs, path)
		if err != nil {
			return err
		}

		dstpath := filepath.Join(waRootPath, path)
		os.MkdirAll(filepath.Dir(dstpath), 0777)

		if s := filepath.Base(path); s == "_keep" || s == ".keep" {
			return nil
		}

		f, err := os.Create(dstpath)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.Write(data); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Builder) getWarootPath(waos string) string {
	return fmt.Sprintf("%s/wa_%s_%s", p.Output, version.Version, waos)
}
